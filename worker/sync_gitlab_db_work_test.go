package worker

import (
	"context"
	"testing"

	"github.com/gofiber/fiber/v2"
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
		ID:       uuid.New(),
		OwnerID:  owner.ID,
		Archived: false,
	}
	testDb.InsertClassroom(&classroom1)

	team1 := database.Team{
		ID:          uuid.New(),
		Name:        "team1",
		ClassroomID: classroom1.ID,
		GroupID:     10,
		Member: []*database.UserClassrooms{
			{
				UserID:      owner.ID,
				ClassroomID: classroom1.ID,
			},
			{
				UserID:      member1.ID,
				ClassroomID: classroom1.ID,
			},
			{
				UserID:      member2.ID,
				ClassroomID: classroom1.ID,
			},
		},
	}
	testDb.InsertTeam(&team1)

	assignment1 := database.Assignment{
		ID:          uuid.New(),
		Name:        "Assignment1",
		ClassroomID: classroom1.ID,
	}
	testDb.InsertAssignment(&assignment1)

	assignment1Project := database.AssignmentProjects{
		ID:           uuid.New(),
		AssignmentID: assignment1.ID,
		Assignment:   assignment1,
		TeamID:       team1.ID,
		ProjectID:    15,
	}
	testDb.InsertAssignmentProjects(&assignment1Project)

	w := NewSyncGitlabDbWork(&gitlabConfig.GitlabConfig{})

	t.Run("getUnarchivedClassrooms", func(t *testing.T) {
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

		dbClassrooms := w.getUnarchivedClassrooms(context.Background())

		if len(dbClassrooms) != 1 {
			t.Errorf("Expected 1 classroom, got %d", len(dbClassrooms))
		}
		assert.Equal(t, classroom1.ID, dbClassrooms[0].ID)
		assert.Len(t, dbClassrooms[0].Member, 3)
		assert.Len(t, dbClassrooms[0].Teams, 1)
	})

	t.Run("syncClassroom", func(t *testing.T) {
		newName := "new name"
		newDescription := "new description"

		repo.EXPECT().
			GetGroupById(classroom1.GroupID).
			Return(&model.Group{
				Name:        newName,
				Description: newDescription,
			}, nil).
			Times(1)

		w.syncClassroom(context.Background(), classroom1, repo)

		repo.AssertExpectations(t)

		dbClassroom1, err := query.Classroom.WithContext(context.Background()).
			Where(query.Classroom.ID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, newName, dbClassroom1.Name)
		assert.Equal(t, newDescription, dbClassroom1.Description)

		// revert changes of db object for next tests
		dbClassroom1.Name = classroom1.Name
		dbClassroom1.Description = classroom1.Description
		query.UserClassrooms.WithContext(context.Background()).Updates(dbClassroom1)
	})

	t.Run("syncClassroomMember", func(t *testing.T) {
		repo.EXPECT().
			GetAllUsersOfGroup(classroom1.GroupID).
			Return([]*model.User{
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
			}, nil).
			Times(1)

		classroom, err := query.Classroom.WithContext(context.Background()).
			Preload(query.Classroom.Member).
			Where(query.Classroom.ID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)

		w.syncClassroomMember(context.Background(), classroom1.GroupID, classroom.Member, repo)

		repo.AssertExpectations(t)

		leftMember, err := query.UserClassrooms.WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(member2.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(classroom1.ID)).
			First()
		assert.Error(t, err)
		assert.Nil(t, leftMember)

		// revert changes for next tests
		testDb.InsertUserClassrooms(&database.UserClassrooms{
			UserID:      member2.ID,
			ClassroomID: classroom1.ID,
		})
	})

	t.Run("syncTeam", func(t *testing.T) {
		newName := "new name"
		repo.EXPECT().
			GetGroupById(team1.GroupID).
			Return(&model.Group{
				Name: newName,
			}, nil).
			Times(1)

		w.syncTeam(context.Background(), team1, repo)

		repo.AssertExpectations(t)

		dbTeam1, err := query.Team.WithContext(context.Background()).
			Where(query.Team.ID.Eq(team1.ID)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, newName, dbTeam1.Name)
	})

	t.Run("syncTeamMember", func(t *testing.T) {
		repo.EXPECT().
			GetAllUsersOfGroup(team1.GroupID).
			Return([]*model.User{
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
			}, nil).
			Times(1)

		team, err := query.Team.WithContext(context.Background()).
			Preload(query.Team.Member).
			Where(query.Team.ID.Eq(team1.ID)).
			First()
		assert.NoError(t, err)

		w.syncTeamMember(context.Background(), team1.GroupID, team.Member, repo)

		repo.AssertExpectations(t)

		leftMember, err := query.UserClassrooms.WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(member2.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		assert.Nil(t, leftMember.TeamID)
		assert.Nil(t, leftMember.Team)

		// revert changes of db object for next
		leftMember, err = query.UserClassrooms.WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(member2.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		leftMember.TeamID = &team1.ID
		leftMember.Team = &team1
		query.UserClassrooms.WithContext(context.Background()).Updates(leftMember)
	})

	t.Run("getAssignmentProjectes", func(t *testing.T) {
		dbAssignments := w.getAssignmentProjects(context.Background(), assignment1.ID)

		if len(dbAssignments) != 1 {
			t.Errorf("Expected 1 assignment, got %d", len(dbAssignments))
		}
		assert.Equal(t, 15, dbAssignments[0].ProjectID)
	})

	t.Run("syncProject", func(t *testing.T) {
		repo.EXPECT().
			GetProjectById(assignment1Project.ProjectID).
			Return(nil, fiber.NewError(404, "404 {message: 404 Project Not Found}")).
			Times(1)

		w.syncProject(context.Background(), assignment1Project, repo)

		repo.AssertExpectations(t)

		deletedProject, err := query.AssignmentProjects.WithContext(context.Background()).
			Where(query.AssignmentProjects.ID.Eq(assignment1Project.ID)).
			First()
		assert.Error(t, err)
		assert.Nil(t, deletedProject)
	})
}
