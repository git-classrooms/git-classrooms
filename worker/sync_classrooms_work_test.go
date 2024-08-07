package worker

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestSyncClassroomsWork(t *testing.T) {
	repo := gitlabRepoMock.NewMockRepository(t)

	testDb := db_tests.NewTestDB(t)

	owner := database.User{
		ID:             1,
		GitlabUsername: "owner",
		GitlabEmail:    "owner",
	}
	testDb.InsertUser(&owner)

	member1 := database.User{
		ID:             2,
		GitlabUsername: "member1",
		GitlabEmail:    "member1",
	}
	testDb.InsertUser(&member1)

	member2 := database.User{
		ID:             3,
		GitlabUsername: "member2",
		GitlabEmail:    "member2",
	}
	testDb.InsertUser(&member2)

	classroom1 := database.Classroom{
		ID:      uuid.New(),
		OwnerID: owner.ID,
		Member: []*database.UserClassrooms{
			{
				UserID: owner.ID,
				Left:   false,
			},
			{
				UserID: member1.ID,
				Left:   false,
			},
			{
				UserID: member2.ID,
				Left:   false,
			},
		},
		Archived: false,
	}
	testDb.InsertClassroom(&classroom1)

	classroom2 := database.Classroom{
		ID:      uuid.New(),
		OwnerID: owner.ID,
		Member: []*database.UserClassrooms{
			{
				UserID: owner.ID,
			},
		},
		Archived: true,
	}
	testDb.InsertClassroom(&classroom2)

	w := NewSyncClassroomsWork(&gitlabConfig.GitlabConfig{})

	t.Run("get unarchived classrooms", func(t *testing.T) {

		classrooms := w.getUnarchivedClassrooms(context.Background())

		if len(classrooms) != 1 {
			t.Errorf("Expected 1 classroom, got %d", len(classrooms))
		}
		assert.Equal(t, classroom1.ID, classrooms[0].ID)
	})

	t.Run("syncClassroom", func(t *testing.T) {
		newName := "new name"
		newDescription := "new description"

		repo.EXPECT().
			GetGroupById(classroom1.GroupID).
			Return(&model.Group{
				Name:        newName,
				Description: newDescription,
				Member: []model.User{
					{
						ID:       owner.ID,
						Username: owner.GitlabUsername,
						Email:    owner.GitlabEmail,
					},
					{
						ID:       member1.ID,
						Username: member1.GitlabUsername,
						Email:    member1.GitlabEmail,
					},
				},
			}, nil).
			Times(1)

		w.syncClassroom(context.Background(), classroom1, repo)

		repo.AssertExpectations(t)

		dbClassroom1, err := query.Classroom.WithContext(context.Background()).
			Preload(query.Classroom.Member).
			Where(query.Classroom.ID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, "new name", dbClassroom1.Name)
		assert.Equal(t, "new description", dbClassroom1.Description)

		leftMember, err := query.UserClassrooms.WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(member2.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		assert.True(t, leftMember.Left)
	})
}
