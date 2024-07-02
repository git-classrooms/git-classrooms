package task

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestCloseDueAssignments(t *testing.T) {
	repo := gitlabRepoMock.NewMockRepository(t)

	testDb := db_tests.NewTestDB(t)

	owner := database.User{
		ID:             1,
		GitlabUsername: "owner",
		GitlabEmail:    "owner",
	}
	testDb.InsertUser(&owner)

	student1 := database.User{
		ID:             2,
		GitlabUsername: "student1",
		GitlabEmail:    "student1",
	}
	testDb.InsertUser(&student1)

	student2 := database.User{
		ID:             3,
		GitlabUsername: "student2",
		GitlabEmail:    "student2",
	}
	testDb.InsertUser(&student2)

	classroom := database.Classroom{
		ID:      uuid.New(),
		OwnerID: owner.ID,
	}
	testDb.InsertClassroom(&classroom)

	assignment1 := database.Assignment{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
	}
	testDb.InsertAssignment(&assignment1)

	team1 := database.Team{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
		GroupID:     1,
		Member: []*database.UserClassrooms{
			{
				UserID:      student1.ID,
				ClassroomID: classroom.ID,
				Role:        database.Student,
			},
			{
				UserID:      student2.ID,
				ClassroomID: classroom.ID,
				Role:        database.Student,
			},
		},
	}
	testDb.InsertTeam(&team1)

	assignmentProject1 := database.AssignmentProjects{
		AssignmentID:  assignment1.ID,
		TeamID:        team1.ID,
		ProjectID:     1,
		ProjectStatus: database.Accepted,
	}
	testDb.InsertAssignmentProjects(&assignmentProject1)

	t.Run("All Assignments already closed", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = true
		testDb.SaveAssignment(&assignment1)

		err := CloseDueAssignments(repo, context.Background())
		assert.NoError(t, err)
	})

	t.Run("No Assignments are due", func(t *testing.T) {
		dueDate := time.Now().Add(1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		testDb.SaveAssignment(&assignment1)

		err := CloseDueAssignments(repo, context.Background())
		assert.NoError(t, err)

		assignment1After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment1.ID)).
			First()
		assert.NoError(t, err)
		assert.False(t, assignment1After.Closed)
	})

	t.Run("Close unaccepted Assignment", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		testDb.SaveAssignment(&assignment1)

		assignmentProject1.ProjectStatus = database.Pending
		testDb.SaveAssignmentProjects(&assignmentProject1)

		err := CloseDueAssignments(repo, context.Background())
		assert.NoError(t, err)

		assignment1After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment1.ID)).
			First()
		assert.NoError(t, err)
		assert.True(t, assignment1After.Closed)

		assignmentProject1.ProjectStatus = database.Accepted
		testDb.SaveAssignmentProjects(&assignmentProject1)
	})

	assignmentProject1.ProjectStatus = database.Accepted
	testDb.SaveAssignmentProjects(&assignmentProject1)

	t.Run("repo.GetAccessLevelOfUserInProject throws error -> restore old permissions", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		testDb.SaveAssignment(&assignment1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, student1.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		repo.EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject1.ProjectID, student1.ID, model.ReporterPermissions).
			Return(nil).
			Times(1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, student2.ID).
			Return(model.NoPermissions, assert.AnError).
			Times(1)

		repo.EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject1.ProjectID, student1.ID, model.DeveloperPermissions).
			Return(nil).
			Times(1)

		err := CloseDueAssignments(repo, context.Background())
		assert.Error(t, err)

		repo.AssertExpectations(t)

		assignment1After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment1.ID)).
			First()
		assert.NoError(t, err)
		assert.False(t, assignment1After.Closed)
	})

	t.Run("repo.ChangeUserAccessLevelInProject throws error", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		testDb.SaveAssignment(&assignment1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, student1.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		repo.EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject1.ProjectID, student1.ID, model.ReporterPermissions).
			Return(assert.AnError).
			Times(1)

		err := CloseDueAssignments(repo, context.Background())
		assert.Error(t, err)

		repo.AssertExpectations(t)

		assignment1After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment1.ID)).
			First()
		assert.NoError(t, err)
		assert.False(t, assignment1After.Closed)
	})

	t.Run("Close due Assignment", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		testDb.SaveAssignment(&assignment1)

		dueDate2 := time.Now().Add(1 * time.Hour)
		assignment2 := database.Assignment{
			ID:          uuid.New(),
			ClassroomID: classroom.ID,
			DueDate:     &dueDate2,
			Closed:      false,
		}
		testDb.InsertAssignment(&assignment2)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, student1.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		repo.EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject1.ProjectID, student1.ID, model.ReporterPermissions).
			Return(nil).
			Times(1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, student2.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		repo.EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject1.ProjectID, student2.ID, model.ReporterPermissions).
			Return(nil).
			Times(1)

		err := CloseDueAssignments(repo, context.Background())
		assert.NoError(t, err)

		repo.AssertExpectations(t)

		assignment1After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment1.ID)).
			First()
		assert.NoError(t, err)
		assert.True(t, assignment1After.Closed)

		assignment2After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment2.ID)).
			First()
		assert.NoError(t, err)
		assert.False(t, assignment2After.Closed)
	})
}
