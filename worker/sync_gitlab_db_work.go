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
	log.Default().Println("Starting sync gitlab db work")

	classrooms := w.getUnarchivedClassrooms(ctx)
	for _, classroom := range classrooms {
		log.Default().Printf("Syncing classroom %s", classroom.Name)

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
			if team.Deleted {
				continue
			}

			log.Default().Printf("Syncing team %s", team.Name)

			err = w.syncTeam(ctx, *team, repo)
			if err != nil {
				continue
			}

			w.syncTeamMember(ctx, team.GroupID, team.Member, repo)
		}

		for _, assignment := range classroom.Assignments {
			projects := w.getAssignmentProjects(ctx, assignment.ID)
			for _, project := range projects {
				log.Default().Printf("Syncing assignment %s, project %d", assignment.ID.String(), project.ProjectID)

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
			dbClassroom.Deleted = true
			_, err := query.Classroom.WithContext(ctx).Updates(dbClassroom)
			if err == nil {
				log.Default().Printf("Classroom %s marked as deleted via gitlab", dbClassroom.Name)
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

		if !found && !dbMember.LeftClassroom {
			dbMember.LeftClassroom = true
			dbMember.LeftTeam = true
			_, err := query.UserClassrooms.WithContext(ctx).Updates(dbMember)
			if err != nil {
				log.Default().Printf("Error could not mark user as left of classroom %d: %s", groupId, err.Error())
			} else {
				log.Default().Printf("Marked user %d as left of classroom %d", dbMember.UserID, groupId)
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

		if !found && !dbMember.LeftTeam {
			dbMember.LeftTeam = true
			_, err := query.UserClassrooms.WithContext(ctx).Updates(dbMember)
			if err != nil {
				log.Default().Printf("Error could not mark user as left of classroom %d: %s", groupId, err.Error())
			} else {
				log.Default().Printf("Marked user %d as left of team %d", dbMember.UserID, groupId)
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
			dbTeam.Deleted = true
			_, err := query.Team.WithContext(ctx).Updates(dbTeam)
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
		_, err := query.AssignmentProjects.WithContext(ctx).Updates(dbProject)
		if err == nil {
			log.Default().Printf("Project with id %d marked as deleted via gitlab", dbProject.ProjectID)
		}
	} else {
		log.Default().Printf("Error while fetching project with id %s. ErrorMsg: %s", dbProject.ID.String(), err.Error())
	}
}
