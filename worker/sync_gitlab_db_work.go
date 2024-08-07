package worker

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
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
		Where(query.Classroom.Archived.Is(false)).
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
		if strings.Contains(err.Error(), "401 {message: 401 Unauthorized}") {
			_, err := query.Classroom.WithContext(ctx).Delete(&dbClassroom)
			if err == nil {
				log.Default().Printf("Classroom %s deleted via gitlab", dbClassroom.Name)
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
	gitlabMember, err := repo.GetAllUsersOfGroup(groupId)
	if err != nil {
		log.Default().Printf("Could not retive members of group with id %d, this could indicate a deleted group. ErrorMsg: %s", groupId, err.Error())
		return
	}

	for _, dbMember := range dbMember {
		found := false

		for _, gitlabMember := range gitlabMember {
			if dbMember.UserID == gitlabMember.ID {
				found = true
				break
			}
		}

		if !found {
			_, err := query.UserClassrooms.WithContext(ctx).Delete(dbMember)
			if err != nil {
				log.Default().Printf("Error could not remove member [%d] from classroom %d: %s", dbMember.UserID, groupId, err.Error())
			} else {
				log.Default().Printf("Removed member %d from classroom %d", dbMember.UserID, groupId)
			}
		}
	}

	// TODO: what about new members, which got added to the classroom via gitlab?

	// TODO: should we reacte to changes in access level via gitlab?
}

func (w *SyncGitlabDbWork) syncTeamMember(ctx context.Context, groupId int, dbMember []*database.UserClassrooms, repo gitlab.Repository) {
	gitlabMember, err := repo.GetAllUsersOfGroup(groupId)
	if err != nil {
		log.Default().Printf("Could not retive members of group with id %d, this could indicate a deleted group. ErrorMsg: %s", groupId, err.Error())
		return
	}

	for _, dbMember := range dbMember {
		found := false

		for _, gitlabMember := range gitlabMember {
			if dbMember.UserID == gitlabMember.ID {
				found = true
				break
			}
		}

		if !found {
			_, err := query.UserClassrooms.WithContext(ctx).Delete(dbMember)
			if err != nil {
				log.Default().Printf("Error could not remove member [%d] from team %d: %s", dbMember.UserID, groupId, err.Error())
				break
			}

			err = query.UserClassrooms.WithContext(ctx).Create(&database.UserClassrooms{
				UserID:      dbMember.UserID,
				ClassroomID: dbMember.ClassroomID,
				Role:        dbMember.Role,
			})
			if err != nil {
				log.Default().Printf("Error could not remove member [%d] from team %d: %s", dbMember.UserID, groupId, err.Error())
			}
		}

	}

	// TODO: what about new members, which got added to the classroom via gitlab?

	// TODO: should we reacte to changes in access level via gitlab?
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
	if err == nil {
		return
	}

	if strings.Contains(err.Error(), "404 {message: 404 Project Not Found}") {
		_, err := query.AssignmentProjects.WithContext(ctx).Delete(&dbProject)
		if err == nil {
			log.Default().Printf("Project with id %d deleted via gitlab", dbProject.ProjectID)
		}
	} else {
		log.Default().Printf("Error while fetching project with id %s. ErrorMsg: %s", dbProject.ID.String(), err.Error())
	}
}
