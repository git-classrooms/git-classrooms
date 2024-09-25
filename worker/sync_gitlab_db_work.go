package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gen/field"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
)

// SyncGitlabDbWork is responsible for synchronizing the GitLab database with the local database.
type SyncGitlabDBWork struct {
	gitlabConfig gitlabConfig.Config
	publicURL    *url.URL
}

// NewSyncGitlabDbWork creates a new instance of SyncGitlabDbWork.
func NewSyncGitlabDBWork(config gitlabConfig.Config, publicURL *url.URL) *SyncGitlabDBWork {
	return &SyncGitlabDBWork{gitlabConfig: config, publicURL: publicURL}
}

// Do synchronizes classrooms, teams, and projects between GitLab and the local database.
func (w *SyncGitlabDBWork) Do(ctx context.Context) {
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
			err = w.syncTeam(ctx, classroom, *team, repo)
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

// getUnarchivedClassrooms retrieves all classrooms that are not archived or deleted.
func (w *SyncGitlabDBWork) getUnarchivedClassrooms(ctx context.Context) []*database.Classroom {
	classrooms, err := query.Classroom.
		WithContext(ctx).
		Preload(query.Classroom.Member).
		Preload(query.Classroom.Member.User).
		Preload(query.Classroom.Teams).
		Preload(field.NewRelation("Teams.Member", "")).
		Preload(field.NewRelation("Teams.Member.User", "")).
		Preload(query.Classroom.Assignments).
		Where(query.Classroom.Archived.Not()).
		Where(query.Classroom.PotentiallyDeleted.Not()).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching classrooms: %s", err.Error())
		return []*database.Classroom{}
	}

	return classrooms
}

// syncClassroom synchronizes the data of a classroom from GitLab with the local database.
func (w *SyncGitlabDBWork) syncClassroom(ctx context.Context, dbClassroom database.Classroom, repo gitlab.Repository) error {
	log.Default().Printf("Syncing classroom %s (ID=%d)", dbClassroom.Name, dbClassroom.GroupID)
	gitlabClassroom, err := repo.GetGroupByID(dbClassroom.GroupID)
	if err != nil {
		var gitLabError *model.GitLabError
		if errors.As(err, &gitLabError) {
			// the following errors are possible:
			// -> classroom deleted -> 403 Forbidden -> after 1 min -> 401 Unauthorized
			// -> access token revoked -> 401 error invalid_token -> after 1 min -> 401 Unauthorized
			if gitLabError.Response.StatusCode == 403 {
				_, err := query.Classroom.WithContext(ctx).Delete(&dbClassroom)
				if err == nil {
					log.Default().Printf("Classroom %s (ID=%d) deleted due to group deletion or member classroom_bot removal via GitLab.", dbClassroom.Name, dbClassroom.GroupID)
				}
			} else if gitLabError.Response.StatusCode == 401 {
				if strings.Contains(gitLabError.Message, "error: invalid_token") {
					dbClassroom.Archived = true
					err := query.Classroom.WithContext(ctx).Save(&dbClassroom)
					if err == nil {
						log.Default().Printf("Classroom %s (ID=%d) archived due to revoked access token", dbClassroom.Name, dbClassroom.GroupID)
					}
				} else if strings.Contains(gitLabError.Message, "message: 401 Unauthorized") {
					dbClassroom.PotentiallyDeleted = true
					err := query.Classroom.WithContext(ctx).Save(&dbClassroom)
					if err == nil {
						log.Default().Printf("Classroom %s (ID=%d) marked as potentially deleted due to 401 Unauthorized. Group access token could be revoked or group could be deleted via GitLab. Clarify on next user access of classroom.", dbClassroom.Name, dbClassroom.GroupID) // Clarify in classroom middleware
					}
				}
			}
		} else {
			log.Default().Printf("Error while fetching group with id %d. ErrorMsg: %s", dbClassroom.GroupID, err.Error())
		}
		return err
	}

	if dbClassroom.Name != gitlabClassroom.Name {
		if _, err := repo.ChangeGroupName(dbClassroom.GroupID, dbClassroom.Name); err != nil {
			log.Default().Printf("Error could not update group name for classroom %d: %s", dbClassroom.GroupID, err.Error())
		}
	}

	shouldDescription := utils.CreateClassroomGitlabDescription(&dbClassroom, w.publicURL)

	if shouldDescription != gitlabClassroom.Description {
		if _, err := repo.ChangeGroupDescription(dbClassroom.GroupID, shouldDescription); err != nil {
			log.Default().Printf("Error could not update group name for classroom %d: %s", dbClassroom.GroupID, err.Error())
		}
	}

	return nil
}

// syncClassroomMember synchronizes the members of a classroom between GitLab and the local database.
func (w *SyncGitlabDBWork) syncClassroomMember(ctx context.Context, groupID int, dbMember []*database.UserClassrooms, repo gitlab.Repository) {
	handleLeftMembers := func(context context.Context, member *database.UserClassrooms, groupID int, repo gitlab.Repository) {
		_, err := query.UserClassrooms.WithContext(ctx).Delete(member)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from classroom %d: %s", member.UserID, groupID, err.Error())
		} else {
			log.Default().Printf("Removed member %d from classroom %d", member.UserID, groupID)
		}
	}

	handleAddedMembers := func(context context.Context, member *model.User, groupID int, repo gitlab.Repository) {
		err := repo.RemoveUserFromGroup(groupID, member.ID)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from gitlab group %d: %s", member.ID, groupID, err.Error())
		} else {
			log.Default().Printf("Removed member %d from gitlab group %d", member.ID, groupID)
		}
	}

	w.syncMember(ctx, groupID, dbMember, repo, handleLeftMembers, handleAddedMembers)
}

// syncTeamMember synchronizes the members of a team between GitLab and the local database.
func (w *SyncGitlabDBWork) syncTeamMember(ctx context.Context, groupID int, dbMember []*database.UserClassrooms, repo gitlab.Repository) {
	// TODO delete Team if teamsize is 1
	handleLeftMembers := func(context context.Context, member *database.UserClassrooms, groupID int, repo gitlab.Repository) {
		member.TeamID = nil
		member.Team = nil
		err := query.UserClassrooms.WithContext(ctx).Save(member)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from team %d: %s", member.UserID, groupID, err.Error())
			return
		}
	}

	handleAddedMembers := func(context context.Context, member *model.User, groupID int, repo gitlab.Repository) {
		err := repo.RemoveUserFromGroup(groupID, member.ID)
		if err != nil {
			log.Default().Printf("Error could not remove member [%d] from gitlab group %d: %s", member.ID, groupID, err.Error())
		} else {
			log.Default().Printf("Removed member %d from gitlab group %d", member.ID, groupID)
		}
	}

	w.syncMember(ctx, groupID, dbMember, repo, handleLeftMembers, handleAddedMembers)
}

// syncMember handles the synchronization of members between GitLab and the local database.
func (w *SyncGitlabDBWork) syncMember(
	ctx context.Context,
	groupID int,
	dbMember []*database.UserClassrooms,
	repo gitlab.Repository,
	handleLeftMembers func(ctx context.Context, member *database.UserClassrooms, groupID int, repo gitlab.Repository),
	handleAddedMembers func(ctx context.Context, member *model.User, groupID int, repo gitlab.Repository),
) {
	gitlabMember, err := repo.GetAllUsersOfGroup(groupID)
	if err != nil {
		log.Default().Printf("Could not retive members of group with id %d. ErrorMsg: %s", groupID, err.Error())
		return
	}

	leftMember := w.leftMembersViaGitlab(dbMember, gitlabMember)
	for _, member := range leftMember {
		handleLeftMembers(ctx, member, groupID, repo)
	}

	addedMember := w.addedMembersViaGitlab(dbMember, gitlabMember, groupID)
	for _, member := range addedMember {
		handleAddedMembers(ctx, member, groupID, repo)
	}
}

// leftMembersViaGitlab finds members who have left the GitLab group but are still present in the local database.
func (w *SyncGitlabDBWork) leftMembersViaGitlab(dbMember []*database.UserClassrooms, gitlabMember []*model.User) []*database.UserClassrooms {
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

// addedMembersViaGitlab finds members who have been added to the GitLab group but are not present in the local database.
func (w *SyncGitlabDBWork) addedMembersViaGitlab(dbMember []*database.UserClassrooms, gitlabMember []*model.User, groupID int) []*model.User {
	addedMember := []*model.User{}

	for _, gitlabMember := range gitlabMember {
		found := false

		for _, dbMember := range dbMember {
			if dbMember.UserID == gitlabMember.ID {
				found = true
				break
			}
		}

		if !found && !w.isGroupBootUser(*gitlabMember, groupID) {
			addedMember = append(addedMember, gitlabMember)
		}
	}

	return addedMember
}

// isGroupBootUser checks if a user is a group bot user in GitLab.
func (w *SyncGitlabDBWork) isGroupBootUser(user model.User, groupID int) bool {
	return strings.Contains(user.Username, fmt.Sprintf("group_%d_bot_", groupID))
}

// syncTeam synchronizes the data of a team from GitLab with the local database.
func (w *SyncGitlabDBWork) syncTeam(ctx context.Context, classroom *database.Classroom, dbTeam database.Team, repo gitlab.Repository) error {
	log.Default().Printf("Syncing team %s (ID=%d)", dbTeam.Name, dbTeam.GroupID)
	gitlabTeam, err := repo.GetGroupByID(dbTeam.GroupID)
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
		if _, err := repo.ChangeGroupName(dbTeam.GroupID, dbTeam.Name); err != nil {
			log.Default().Printf("Error could not update group name for team %d: %s", dbTeam.GroupID, err.Error())
		}
	}

	shouldDescription := utils.CreateTeamGitlabDescription(classroom, &dbTeam, w.publicURL)
	if shouldDescription != gitlabTeam.Description {
		if _, err := repo.ChangeGroupDescription(dbTeam.GroupID, shouldDescription); err != nil {
			log.Default().Printf("Error could not update group name for team %d: %s", dbTeam.GroupID, err.Error())
		}
	}

	return nil
}

// getAssignmentProjects retrieves the list of projects associated with a given assignment.
func (w *SyncGitlabDBWork) getAssignmentProjects(ctx context.Context, assignmentID uuid.UUID) []*database.AssignmentProjects {
	projects, err := query.AssignmentProjects.
		WithContext(ctx).
		Where(query.AssignmentProjects.AssignmentID.Eq(assignmentID)).
		Where(query.AssignmentProjects.ProjectStatus.Eq(string(database.Accepted))).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching classrooms: %s", err.Error())
		return []*database.AssignmentProjects{}
	}

	return projects
}

// syncProject synchronizes the project data between GitLab and the local database.
func (w *SyncGitlabDBWork) syncProject(ctx context.Context, dbProject database.AssignmentProjects, repo gitlab.Repository) {
	_, err := repo.GetProjectByID(dbProject.ProjectID)
	if err == nil || !strings.Contains(err.Error(), "404 {message: 404 Project Not Found}") {
		return
	}

	_, err = query.AssignmentProjects.WithContext(ctx).Delete(&dbProject)
	if err != nil {
		log.Default().Printf("Error while fetching project with id %s. ErrorMsg: %s", dbProject.ID.String(), err.Error())
	}

	log.Default().Printf("Project with id %d deleted via gitlab", dbProject.ProjectID)
}
