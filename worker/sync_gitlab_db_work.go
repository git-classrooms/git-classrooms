package worker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type SyncGitlabDbWork struct {
	gitlabConfig gitlabConfig.Config
}

func NewSyncGitlabDbWork(config gitlabConfig.Config) *SyncGitlabDbWork {
	return &SyncGitlabDbWork{gitlabConfig: config}
}

func (w *SyncGitlabDbWork) Do(ctx context.Context) {
	classrooms := w.getUnarchivedClassrooms(ctx)
	for _, classroom := range classrooms {
		repo, err := GetWorkerRepo(w.gitlabConfig, classroom.GroupAccessToken)
		if err != nil {
			log.Default().Printf("Error occurred while login into gitlab: %s", err.Error())
			continue
		}

		err = w.syncClassroom(ctx, *classroom, repo)
		if err != nil {
			continue
		}

		w.syncClassroomMember(ctx, classroom.GroupID, classroom.Member, repo)

		for _, team := range classroom.Teams {
			err = w.syncTeam(ctx, *team, repo)
			if err != nil {
				continue
			}

			w.syncTeamMember(ctx, team.GroupID, team.Member, repo)
		}

		for _, assignment := range classroom.Assignments {
			projects := w.getAssignmentProjects(ctx, assignment.ID)
			for _, project := range projects {
				w.syncProject(ctx, *project, repo)
			}
		}
	}
}

func (w *SyncGitlabDbWork) getUnarchivedClassrooms(ctx context.Context) []*database.Classroom {
	classrooms, err := query.Classroom.
		WithContext(ctx).
		Preload(query.Classroom.Member).
		Preload(query.Classroom.Teams).
		Preload(query.Classroom.Assignments).
		Where(query.Classroom.Archived.Not()).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching classrooms: %s", err.Error())
		return []*database.Classroom{}
	}

	return classrooms
}

func (w *SyncGitlabDbWork) syncClassroom(ctx context.Context, dbClassroom database.Classroom, repo gitlab.Repository) error {
	gitlabClassroom, err := repo.GetGroupById(dbClassroom.GroupID)
	if err != nil {
		if strings.Contains(err.Error(), "403 {message: 403 Forbidden}") {
			_, err := query.Classroom.WithContext(ctx).Delete(&dbClassroom)
			if err == nil {
				log.Default().Printf("Classroom %s (ID=%d) deleted due to group deletion or member classroom_bot removal via GitLab.", dbClassroom.Name, dbClassroom.GroupID)
			}
		} else if strings.Contains(err.Error(), "401 {message: 401 Unauthorized}") || strings.Contains(err.Error(), "401 {error: invalid_token}") {
			dbClassroom.Archived = true
			_, err := query.Classroom.WithContext(ctx).Updates(dbClassroom)
			if err == nil {
				log.Default().Printf("Classroom %s (ID=%d) archived due to revoked access token", dbClassroom.Name, dbClassroom.GroupID)
			}
		} else {
			log.Default().Printf("Error while fetching group with id %d. ErrorMsg: %s", dbClassroom.GroupID, err.Error())
		}
		return err
	}

	needsUpdate := false

	if dbClassroom.Name != gitlabClassroom.Name {
		dbClassroom.Name = gitlabClassroom.Name
		needsUpdate = true
	}

	if dbClassroom.Description != gitlabClassroom.Description {
		dbClassroom.Description = gitlabClassroom.Description
		needsUpdate = true
	}

	if needsUpdate {
		query.Classroom.WithContext(ctx).Updates(dbClassroom)
	}

	return nil
}

func (w *SyncGitlabDbWork) syncClassroomMember(ctx context.Context, groupId int, dbMember []*database.UserClassrooms, repo gitlab.Repository) {
	handleLeftMembers := func(context context.Context, member *database.UserClassrooms, groupId int, repo gitlab.Repository) {
		_, err := query.UserClassrooms.WithContext(ctx).Delete(member)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from classroom %d: %s", member.UserID, groupId, err.Error())
		} else {
			log.Default().Printf("Removed member %d from classroom %d", member.UserID, groupId)
		}
	}

	handleAddedMembers := func(context context.Context, member *model.User, groupId int, repo gitlab.Repository) {
		err := repo.RemoveUserFromGroup(groupId, member.ID)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from gitlab group %d: %s", member.ID, groupId, err.Error())
		} else {
			log.Default().Printf("Removed member %d from gitlab group %d", member.ID, groupId)
		}
	}

	w.syncMember(ctx, groupId, dbMember, repo, handleLeftMembers, handleAddedMembers)
}

func (w *SyncGitlabDbWork) syncTeamMember(ctx context.Context, groupId int, dbMember []*database.UserClassrooms, repo gitlab.Repository) {
	handleLeftMembers := func(context context.Context, member *database.UserClassrooms, groupId int, repo gitlab.Repository) {
		member.TeamID = nil
		member.Team = nil
		err := query.UserClassrooms.WithContext(ctx).Save(member)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from team %d: %s", member.UserID, groupId, err.Error())
			return
		}
	}

	handleAddedMembers := func(context context.Context, member *model.User, groupId int, repo gitlab.Repository) {
		err := repo.RemoveUserFromGroup(groupId, member.ID)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from gitlab group %d: %s", member.ID, groupId, err.Error())
		} else {
			log.Default().Printf("Removed member %d from gitlab group %d", member.ID, groupId)
		}
	}

	w.syncMember(ctx, groupId, dbMember, repo, handleLeftMembers, handleAddedMembers)
}

func (w *SyncGitlabDbWork) syncMember(
	ctx context.Context,
	groupId int,
	dbMember []*database.UserClassrooms,
	repo gitlab.Repository,
	handleLeftMembers func(ctx context.Context, member *database.UserClassrooms, groupId int, repo gitlab.Repository),
	handleAddedMembers func(ctx context.Context, member *model.User, groupId int, repo gitlab.Repository),
) {
	gitlabMember, err := repo.GetAllUsersOfGroup(groupId)
	if err != nil {
		log.Default().Printf("Could not retive members of group with id %d. ErrorMsg: %s", groupId, err.Error())
		return
	}

	leftMember := w.leftMembersViaGitlab(dbMember, gitlabMember)
	for _, member := range leftMember {
		handleLeftMembers(ctx, member, groupId, repo)
	}

	addedMember := w.addedMembersViaGitlab(dbMember, gitlabMember, groupId)
	for _, member := range addedMember {
		handleAddedMembers(ctx, member, groupId, repo)
	}
}

func (w *SyncGitlabDbWork) leftMembersViaGitlab(dbMember []*database.UserClassrooms, gitlabMember []*model.User) []*database.UserClassrooms {
	leftMember := []*database.UserClassrooms{}

	for _, dbMember := range dbMember {
		found := false

		for _, gitlabMember := range gitlabMember {
			if dbMember.UserID == gitlabMember.ID {
				found = true
				break
			}
		}

		if !found {
			leftMember = append(leftMember, dbMember)
		}
	}

	return leftMember
}

func (w *SyncGitlabDbWork) addedMembersViaGitlab(dbMember []*database.UserClassrooms, gitlabMember []*model.User, groupId int) []*model.User {
	addedMember := []*model.User{}

	for _, gitlabMember := range gitlabMember {
		found := false

		for _, dbMember := range dbMember {
			if dbMember.UserID == gitlabMember.ID {
				found = true
				break
			}
		}

		if !found && !w.isGroupBootUser(*gitlabMember, groupId) {
			addedMember = append(addedMember, gitlabMember)
		}
	}

	return addedMember
}

func (w *SyncGitlabDbWork) isGroupBootUser(user model.User, groupId int) bool {
	return strings.Contains(user.Username, fmt.Sprintf("group_%d_bot_", groupId))
}

func (w *SyncGitlabDbWork) syncTeam(ctx context.Context, dbTeam database.Team, repo gitlab.Repository) error {
	gitlabTeam, err := repo.GetGroupById(dbTeam.GroupID)
	if err != nil {
		if strings.Contains(err.Error(), "404 {message: 404 Group Not Found}") {
			_, err := query.Team.WithContext(ctx).Delete(&dbTeam)
			if err == nil {
				log.Default().Printf("Team %s marked as deleted via gitlab", dbTeam.Name)
			}
		} else {
			log.Default().Printf("Error while fetching group with id %d. ErrorMsg: %s", dbTeam.GroupID, err.Error())
		}

		return err
	}

	if dbTeam.Name != gitlabTeam.Name {
		dbTeam.Name = gitlabTeam.Name
		query.Team.WithContext(ctx).Updates(dbTeam)
	}

	return nil
}

func (w *SyncGitlabDbWork) getAssignmentProjects(ctx context.Context, assignmentId uuid.UUID) []*database.AssignmentProjects {
	projects, err := query.AssignmentProjects.
		WithContext(ctx).
		Preload(query.AssignmentProjects.GradingJUnitTestResult).
		Where(query.AssignmentProjects.AssignmentID.Eq(assignmentId)).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching classrooms: %s", err.Error())
		return []*database.AssignmentProjects{}
	}

	return projects
}

func (w *SyncGitlabDbWork) syncProject(ctx context.Context, dbProject database.AssignmentProjects, repo gitlab.Repository) {
	_, err := repo.GetProjectById(dbProject.ProjectID)
	if err == nil || !strings.Contains(err.Error(), "404 {message: 404 Project Not Found}") {
		return
	}

	if dbProject.GradingJUnitTestResult != nil {
		_, err := query.JUnitTestResult.WithContext(ctx).Delete(dbProject.GradingJUnitTestResult)
		if err != nil {
			log.Default().Printf("Project with id %d deleted via gitlab, but failed to delete GradingJUnitTestResults of project", dbProject.ProjectID)
			return
		}
	}

	_, err = query.AssignmentProjects.WithContext(ctx).Delete(&dbProject)
	if err != nil {
		log.Default().Printf("Error while fetching project with id %s. ErrorMsg: %s", dbProject.ID.String(), err.Error())
	}

	log.Default().Printf("Project with id %d deleted via gitlab", dbProject.ProjectID)
}
