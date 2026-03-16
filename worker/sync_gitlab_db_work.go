package worker

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/google/uuid"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/logging"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gorm.io/gen/field"
)

// SyncGitlabDbWork is responsible for synchronizing the GitLab database with the local database.
type SyncGitlabDbWork struct {
	gitlabConfig gitlabConfig.Config
	publicURL    *url.URL
}

// NewSyncGitlabDbWork creates a new instance of SyncGitlabDbWork.
func NewSyncGitlabDbWork(config gitlabConfig.Config, publicUrl *url.URL) *SyncGitlabDbWork {
	return &SyncGitlabDbWork{gitlabConfig: config, publicURL: publicUrl}
}

// Do synchronizes classrooms, teams, and projects between GitLab and the local database.
func (w *SyncGitlabDbWork) Do(ctx context.Context) {
	log := logging.GetLogger(ctx)
	classrooms := w.getUnarchivedClassrooms(ctx)
	for _, classroom := range classrooms {
		log := log.With("classsroom", classroom, "classroomGroupID", classroom.GroupID)
		ctx := logging.SetLogger(ctx, log)

		repo, err := GetWorkerRepo(ctx, w.gitlabConfig, classroom.GroupAccessToken)
		if err != nil {
			log.Error("Error occurred while login into gitlab", "classroom", classroom, "error", err)
			continue
		}

		err = w.syncClassroom(ctx, *classroom, repo)
		if err != nil {
			continue
		}

		classroomMemberCount := len(classroom.Member)
		if classroomMemberCount > 0 {
			log.Info(fmt.Sprintf("fetched %d classroom members for syncing", classroomMemberCount))
		}

		w.syncClassroomMember(ctx, classroom.GroupID, classroom.Member, repo)

		teamsCount := len(classroom.Teams)
		if teamsCount > 0 {
			log.Info(fmt.Sprintf("fetched %d teams for syncing", teamsCount))
		}

		for _, team := range classroom.Teams {
			log := log.With("team", team, "teamGroupID", team.GroupID)
			ctx := logging.SetLogger(ctx, log)

			err = w.syncTeam(ctx, classroom, *team, repo)
			if err != nil {
				continue
			}

			teamMemberCount := len(team.Member)
			if teamMemberCount > 0 {
				log.Info(fmt.Sprintf("fetched %d team members for syncing", teamMemberCount))
			}

			w.syncTeamMember(ctx, team.GroupID, team.Member, classroom, repo)
		}

		assignmentsCount := len(classroom.Assignments)
		if assignmentsCount > 0 {
			log.Info(fmt.Sprintf("fetched %d assignments for syncing", assignmentsCount))
		}

		for _, assignment := range classroom.Assignments {
			log := log.With("assignment", assignment)
			ctx := logging.SetLogger(ctx, log)

			projects := w.getAssignmentProjects(ctx, assignment.ID)
			for _, project := range projects {
				log := log.With("project", project, "projectID", project.ProjectID)
				ctx := logging.SetLogger(ctx, log)

				err = w.syncProject(ctx, *project, repo)
				if err != nil {
					continue
				}
			}
		}
	}
}

// getUnarchivedClassrooms retrieves all classrooms that are not archived or deleted.
func (w *SyncGitlabDbWork) getUnarchivedClassrooms(ctx context.Context) []*database.Classroom {
	log := logging.GetLogger(ctx)
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
		log.Error("error occurred while fetching classrooms", "error", err)
		return []*database.Classroom{}
	}

	classroomsCount := len(classrooms)
	if classroomsCount > 0 {
		log.Info(fmt.Sprintf("fetched %d classrooms for syncing", classroomsCount))
	}

	return classrooms
}

// syncClassroom synchronizes the data of a classroom from GitLab with the local database.
func (w *SyncGitlabDbWork) syncClassroom(ctx context.Context, dbClassroom database.Classroom, repo gitlab.Repository) error {
	log := logging.GetLogger(ctx)
	log.Info("start syncing classroom")

	gitlabClassroom, err := repo.GetGroupById(dbClassroom.GroupID)
	if err != nil {
		var gitLabError *model.GitLabError
		if errors.As(err, &gitLabError) {
			// the following errors are possible:
			// -> classroom deleted -> 403 Forbidden -> after 1 min -> 401 Unauthorized
			// -> access token revoked -> 401 error invalid_token -> after 1 min -> 401 Unauthorized
			if gitLabError.Response.StatusCode == 403 {
				_, err := query.Classroom.WithContext(ctx).Delete(&dbClassroom)
				if err == nil {
					log.Warn("Classroom deleted due to group deletion or member classroom_bot removal via GitLab.")
				}
			} else if gitLabError.Response.StatusCode == 401 {
				if strings.Contains(gitLabError.Message, "error: invalid_token") {
					dbClassroom.Archived = true
					err := query.Classroom.WithContext(ctx).Save(&dbClassroom)
					if err == nil {
						log.Warn("Classroom archived due to revoked access token")
					}
				} else if strings.Contains(gitLabError.Message, "message: 401 Unauthorized") {
					dbClassroom.PotentiallyDeleted = true
					err := query.Classroom.WithContext(ctx).Save(&dbClassroom)
					if err == nil {
						log.Warn("Classroom marked as potentially deleted due to 401 Unauthorized. Group access token could be revoked or group could be deleted via GitLab. Clarify on next user access of classroom.") // Clarify in classroom middleware
					}
				}
			}
		} else {
			log.Error("error while fetching gitlab group", "error", err)
		}
		return err
	}

	if dbClassroom.Name != gitlabClassroom.Name {
		if _, err := repo.ChangeGroupName(dbClassroom.GroupID, dbClassroom.Name); err != nil {
			log.Error("error while updating group name", "error", err)
		} else {
			log.Info("synced classroom name")
		}
	}

	shouldDescription := utils.CreateClassroomGitlabDescription(&dbClassroom, w.publicURL)

	if shouldDescription != gitlabClassroom.Description {
		if _, err := repo.ChangeGroupDescription(dbClassroom.GroupID, shouldDescription); err != nil {
			log.Error("error while updating group description", "error", err)
		} else {
			log.Info("synced classroom description")
		}
	}

	return nil
}

// syncClassroomMember synchronizes the members of a classroom between GitLab and the local database.
func (w *SyncGitlabDbWork) syncClassroomMember(ctx context.Context, groupId int, dbMember []*database.UserClassrooms, repo gitlab.Repository) {
	handleLeftMembers := func(context context.Context, member *database.UserClassrooms, groupId int, repo gitlab.Repository) {
		log := logging.GetLogger(ctx)
		_, err := query.UserClassrooms.WithContext(ctx).Delete(member)
		if err != nil {
			log.Error("error could not remove member from classroom", "error", err)
		} else {
			log.Warn("removed member from classroom because he left the gitlab group")
		}
	}

	handleAddedMembers := func(context context.Context, member *model.User, groupId int, repo gitlab.Repository) {
		log := logging.GetLogger(ctx)
		err := repo.RemoveUserFromGroup(groupId, member.ID)
		if err != nil {
			log.Error("error could not remove member from gitlab group", "error", err)
		} else {
			log.Warn("removed member from gitlab group because he is not registered as a classroom member")
		}
	}

	w.syncMember(ctx, groupId, dbMember, repo, handleLeftMembers, handleAddedMembers)
}

// syncTeamMember synchronizes the members of a team between GitLab and the local database.
func (w *SyncGitlabDbWork) syncTeamMember(ctx context.Context, groupId int, dbMember []*database.UserClassrooms, classroom *database.Classroom, repo gitlab.Repository) {
	handleLeftMembers := func(context context.Context, member *database.UserClassrooms, groupId int, repo gitlab.Repository) {
		log := logging.GetLogger(ctx)

		team := member.Team

		member.TeamID = nil
		member.Team = nil
		if err := query.UserClassrooms.WithContext(ctx).Save(member); err != nil {
			log.Error("error could not remove member from team", "error", err)
			return
		} else {
			log.Warn("Removed member from team")
		}

		if classroom.MaxTeamSize == 1 {
			if _, err := query.Team.WithContext(ctx).Delete(team); err != nil {
				log.Error("error could not delete team", "error", err)
				return
			} else {
				log.Warn("deleted team because teams are not activated")
			}
		}
	}

	handleAddedMembers := func(context context.Context, member *model.User, groupId int, repo gitlab.Repository) {
		log := logging.GetLogger(ctx)

		if err := repo.RemoveUserFromGroup(groupId, member.ID); err != nil {
			log.Error("error could not remove member from gitlab group", "error", err)
			return
		} else {
			log.Warn("Removed member from gitlab group")
		}

		if classroom.MaxTeamSize == 1 {
			if err := repo.DeleteGroup(groupId); err != nil {
				log.Error("error could not delete gitlab group", "error", err)
				return
			} else {
				log.Warn("deleted gitlab group because teams are not activated")
			}
		}
	}

	w.syncMember(ctx, groupId, dbMember, repo, handleLeftMembers, handleAddedMembers)
}

// syncMember handles the synchronization of members between GitLab and the local database.
func (w *SyncGitlabDbWork) syncMember(
	ctx context.Context,
	groupId int,
	dbMember []*database.UserClassrooms,
	repo gitlab.Repository,
	handleLeftMembers func(ctx context.Context, member *database.UserClassrooms, groupId int, repo gitlab.Repository),
	handleAddedMembers func(ctx context.Context, member *model.User, groupId int, repo gitlab.Repository),
) {
	log := logging.GetLogger(ctx)
	log.Info("syncing classroom members")

	gitlabMember, err := repo.GetAllUsersOfGroup(groupId)
	if err != nil {
		log.Error("could not retive members of gitlab group", "error", err)
		return
	}

	leftMember := w.leftMembersViaGitlab(dbMember, gitlabMember)
	leftMemberCount := len(leftMember)
	if leftMemberCount > 0 {
		log.Info(fmt.Sprintf("%d members left via gitlab", leftMemberCount))
	}

	for _, member := range leftMember {
		log := log.With("member", member)
		ctx := logging.SetLogger(ctx, log)
		handleLeftMembers(ctx, member, groupId, repo)
	}

	addedMember := w.addedMembersViaGitlab(dbMember, gitlabMember, groupId)
	addedMemberCount := len(addedMember)
	if addedMemberCount > 0 {
		log.Info(fmt.Sprintf("%d members added via gitlab", addedMemberCount))
	}
	for _, member := range addedMember {
		log := log.With("userID", member.ID)
		ctx := logging.SetLogger(ctx, log)
		handleAddedMembers(ctx, member, groupId, repo)
	}
}

// leftMembersViaGitlab finds members who have left the GitLab group but are still present in the local database.
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

// addedMembersViaGitlab finds members who have been added to the GitLab group but are not present in the local database.
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

// isGroupBootUser checks if a user is a group bot user in GitLab.
func (w *SyncGitlabDbWork) isGroupBootUser(user model.User, groupId int) bool {
	return strings.Contains(user.Username, fmt.Sprintf("group_%d_bot_", groupId))
}

// syncTeam synchronizes the data of a team from GitLab with the local database.
func (w *SyncGitlabDbWork) syncTeam(ctx context.Context, classroom *database.Classroom, dbTeam database.Team, repo gitlab.Repository) error {
	log := logging.GetLogger(ctx)

	log.Info("Syncing team")

	gitlabTeam, err := repo.GetGroupById(dbTeam.GroupID)
	if err != nil {
		if strings.Contains(err.Error(), "404 {message: 404 Group Not Found}") {
			_, err := query.Team.WithContext(ctx).Delete(&dbTeam)
			if err != nil {
				log.Error("error while deleting team", "error", err)
			} else {
				log.Warn("team marked as deleted via gitlab")
			}
		} else {
			log.Error("error while fetching gitlab group", "error", err)
		}

		return err
	}

	if dbTeam.Name != gitlabTeam.Name {
		if _, err := repo.ChangeGroupName(dbTeam.GroupID, dbTeam.Name); err != nil {
			log.Error("error could not update group name for team", "error", err)
		} else {
			log.Info("synced team name")
		}

	}

	shouldDescription := utils.CreateTeamGitlabDescription(classroom, &dbTeam, w.publicURL)
	if shouldDescription != gitlabTeam.Description {
		if _, err := repo.ChangeGroupDescription(dbTeam.GroupID, shouldDescription); err != nil {
			log.Error("error could not update group description for team", "error", err)
		} else {
			log.Info("synced team description")
		}
	}

	return nil
}

// getAssignmentProjects retrieves the list of projects associated with a given assignment.
func (w *SyncGitlabDbWork) getAssignmentProjects(ctx context.Context, assignmentId uuid.UUID) []*database.AssignmentProjects {
	log := logging.GetLogger(ctx)

	projects, err := query.AssignmentProjects.
		WithContext(ctx).
		Preload(query.AssignmentProjects.Team).
		Preload(field.NewRelation("Team.Member", "")).
		Preload(query.AssignmentProjects.Assignment).
		Where(query.AssignmentProjects.AssignmentID.Eq(assignmentId)).
		Where(query.AssignmentProjects.ProjectStatus.Eq(string(database.Accepted))).
		Find()
	if err != nil {
		log.Error("error occurred while fetching projects", "error", err)
		return []*database.AssignmentProjects{}
	}

	projectsCount := len(projects)
	if projectsCount > 0 {
		log.Info(fmt.Sprintf("fetched %d projects for syncing", projectsCount))
	}

	return projects
}

// syncProject synchronizes the project data between GitLab and the local database.
func (w *SyncGitlabDbWork) syncProject(ctx context.Context, dbProject database.AssignmentProjects, repo gitlab.Repository) error {
	log := logging.GetLogger(ctx)

	_, err := repo.GetProjectById(dbProject.ProjectID)
	if err != nil {
		if strings.Contains(err.Error(), "404 {message: 404 Project Not Found}") {
			_, err := query.AssignmentProjects.WithContext(ctx).Delete(&dbProject)
			if err != nil {
				log.Error("error while fetching project", "error", err)
			} else {
				log.Warn("project deleted via gitlab")
			}
		} else {
			log.Error("error while fetching gitlab project", "error", err)
		}
		return err
	}

	userIDs := utils.Map(dbProject.Team.Member, func(member *database.UserClassrooms) int {
		return member.UserID
	})

	gitlabProjectUsers, err := repo.GetAllUsersOfProject(dbProject.ProjectID)
	if err != nil {
		log.Error("can't get users of gitlab group", "error", err)
	}

	targetPermission := model.DeveloperPermissions
	if dbProject.Assignment.Closed {
		targetPermission = model.ReporterPermissions
	}

	for _, userID := range userIDs {
		log := log.With("userID", userID)

		foundIdx := slices.IndexFunc(gitlabProjectUsers, func(g *model.User) bool {
			return g.ID == userID
		})

		if foundIdx == -1 {
			if err := repo.AddProjectMember(dbProject.ProjectID, userID, targetPermission); err != nil {
				log.Error("failed to add member to project", "error", err)
				continue
			}
			log.Warn("add new member to project")
			continue
		}

		gitlabUser := gitlabProjectUsers[foundIdx]
		if *gitlabUser.Permission != targetPermission {
			if err := repo.ChangeUserAccessLevelInProject(dbProject.ProjectID, userID, targetPermission); err != nil {
				log.Error("failed to change permission from user in project", "error", err)
				continue
			}
			log.Warn("change user permission in project")
			continue
		}
	}

	return nil
}
