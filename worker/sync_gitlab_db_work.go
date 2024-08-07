package worker

import (
	"context"
	"log"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
)

type SyncGitlabDbWork struct {
	gitlabConfig gitlabConfig.Config
}

func NewSyncClassroomsWork(config gitlabConfig.Config) *SyncGitlabDbWork {
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

		w.syncClassroom(ctx, *classroom, repo)

		w.syncClassroomMember(ctx, classroom.GroupID, classroom.Member, repo)

		for _, team := range classroom.Teams {
			w.syncTeam(ctx, *team, repo)

			w.syncTeamMember(ctx, team.GroupID, team.Member, repo)
		}
	}
}

func (w *SyncGitlabDbWork) getUnarchivedClassrooms(ctx context.Context) []*database.Classroom {
	classrooms, err := query.Classroom.
		WithContext(ctx).
		Preload(query.Classroom.Member).
		Preload(query.Classroom.Teams).
		Where(query.Classroom.Archived.Is(false)).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching classrooms: %s", err.Error())
		return []*database.Classroom{}
	}

	return classrooms
}

func (w *SyncGitlabDbWork) syncClassroom(ctx context.Context, dbClassroom database.Classroom, repo gitlab.Repository) {
	gitlabClassroom, err := repo.GetGroupById(dbClassroom.GroupID)
	if err != nil {
		log.Default().Printf("Could not retive classroom with id %s, this could indicate a deleted classroom. ErrorMsg: %s", dbClassroom.ID.String(), err.Error())
		// TODO: should we react to the group got deleted via gitlab
		return
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
			_, err := query.UserClassrooms.WithContext(ctx).Updates(dbMember)
			if err != nil {
				log.Default().Printf("Error could not mark user as left of classroom %d: %s", groupId, err.Error())
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
			}
		}
	}

	// TODO: what about new members, which got added to the classroom via gitlab?

	// TODO: should we reacte to changes in access level via gitlab?
}

func (w *SyncGitlabDbWork) syncTeam(ctx context.Context, dbTeam database.Team, repo gitlab.Repository) {
	gitlabTeam, err := repo.GetGroupById(dbTeam.GroupID)
	if err != nil {
		log.Default().Printf("Could not retive team with id %s, this could indicate a deleted team. ErrorMsg: %s", dbTeam.ID.String(), err.Error())

		// TODO: should we react to the group got deleted via gitlab
	}

	if dbTeam.Name != gitlabTeam.Name {
		dbTeam.Name = gitlabTeam.Name
		query.Team.WithContext(ctx).Updates(dbTeam)
	}
}
