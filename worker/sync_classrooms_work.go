package worker

import (
	"context"
	"log"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type SyncClassroomsWork struct {
	gitlabConfig gitlabConfig.Config
}

func NewSyncClassroomsWork(config gitlabConfig.Config) *SyncClassroomsWork {
	return &SyncClassroomsWork{gitlabConfig: config}
}

func (w *SyncClassroomsWork) Do(ctx context.Context) {
	classrooms := w.getUnarchivedClassrooms(ctx)
	for _, classroom := range classrooms {
		repo, err := GetWorkerRepo(w.gitlabConfig, classroom.GroupAccessToken)
		if err != nil {
			log.Default().Printf("Error occurred while login into gitlab: %s", err.Error())
			continue
		}

		w.syncClassroom(ctx, *classroom, repo)
	}
}

func (w *SyncClassroomsWork) getUnarchivedClassrooms(ctx context.Context) []*database.Classroom {
	classrooms, err := query.Classroom.
		WithContext(ctx).
		Preload(query.Classroom.Member).
		Where(query.Classroom.Archived.Is(false)).
		Find()
	if err != nil {
		log.Default().Printf("Error occurred while fetching classrooms: %s", err.Error())
		return []*database.Classroom{}
	}

	return classrooms
}

func (w *SyncClassroomsWork) syncClassroom(ctx context.Context, dbClassroom database.Classroom, repo gitlab.Repository) {
	gitlabClassroom, err := repo.GetGroupById(dbClassroom.GroupID)
	if err != nil {
		log.Default().Printf("Could not retive classroom with id %s, this could indicate a deleted classroom. ErrorMsg: %s", dbClassroom.ID.String(), err.Error())
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

	w.syncClassroomMember(ctx, dbClassroom.Member, gitlabClassroom.Member)

	if needsUpdate {
		query.Classroom.WithContext(ctx).Updates(dbClassroom)
	}
}

func (w *SyncClassroomsWork) syncClassroomMember(ctx context.Context, dbMember []*database.UserClassrooms, gitlabMember []model.User) {
	newlyLeftMembers := w.newlyLeftMembers(gitlabMember, dbMember)
	for _, newlyLeftMember := range newlyLeftMembers {
		newlyLeftMember.Left = true
		query.UserClassrooms.WithContext(ctx).Updates(newlyLeftMember)
	}

	// TODO: what about new members, which got added to the classroom via gitlab?

	// TODO: should we reacte to changes in access level via gitlab?
}

func (w *SyncClassroomsWork) newlyLeftMembers(gitlabMember []model.User, dbMember []*database.UserClassrooms) []*database.UserClassrooms {
	newlyLeftMembers := []*database.UserClassrooms{}

	for _, dbMember := range dbMember {
		found := false
		for _, gitlabMember := range gitlabMember {
			if dbMember.UserID == gitlabMember.ID && !dbMember.Left {
				found = true
				break
			}
		}

		if !found {
			newlyLeftMembers = append(newlyLeftMembers, dbMember)
		}
	}

	return newlyLeftMembers
}
