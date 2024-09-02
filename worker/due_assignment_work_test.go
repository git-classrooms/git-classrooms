package worker

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestDueAssignmentWorker(t *testing.T) {
	pg, err := db_tests.StartPostgres()
	if err != nil {
		t.Fatalf("Failed to start postgres container: %s", err.Error())
	}

	dbURL, err := pg.ConnectionString(context.Background())
	if err != nil {
		t.Fatalf("Failed to obtain connection string: %s", err.Error())
	}

	db, err := gorm.Open(postgres.Open(dbURL))
	if err != nil {
		t.Fatal(err)
	}

	// 1. Migrate database
	err = utils.MigrateDatabase(db)
	if err != nil {
		t.Fatalf("could not migrate database: %s", err.Error())
	}

	db, err := gorm.Open(postgres.Open(dbURL))
	if err != nil {
		t.Fatal(err)
	}

	// 1. Migrate database
	err = utils.MigrateDatabase(db)
	if err != nil {
		t.Fatalf("could not migrate database: %s", err.Error())
	}

	query.SetDefault(db)
	repo := gitlabRepoMock.NewMockRepository(t)

	owner := factory.User()
	student1 := factory.User()
	student2 := factory.User()
	classroom := factory.Classroom(owner.ID)

	dueDate := time.Now().Add(1 * time.Hour)
	assignment1 := factory.Assignment(classroom.ID, &dueDate)

	members := []*database.UserClassrooms{
		factory.UserClassroom(student1.ID, classroom.ID, database.Student),
		factory.UserClassroom(student2.ID, classroom.ID, database.Student),
	}

	team1 := factory.Team(classroom.ID, members)

	assignmentProject1 := factory.AssignmentProject(assignment1.ID, team1.ID)

	work := NewDueAssignmentWork(&gitlab.GitlabConfig{})

	t.Run("All Assignments already closed", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = true
		SaveAssignment(t, assignment1)

		assignments := work.getAssignments2Close(context.Background())
		assert.Empty(t, assignments)
	})

	t.Run("No Assignments are due", func(t *testing.T) {
		dueDate := time.Now().Add(1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		SaveAssignment(t, assignment1)

		assignments := work.getAssignments2Close(context.Background())
		assert.Empty(t, assignments)
	})

	t.Run("Fetches due Assignments", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		SaveAssignment(t, assignment1)

		assignments := work.getAssignments2Close(context.Background())
		assert.Len(t, assignments, 1)
		assert.Equal(t, assignment1.ID, assignments[0].ID)
	})

	t.Run("Close unaccepted Assignment", func(t *testing.T) {
		dueDate := time.Now().Add(-1 * time.Hour)
		assignment1.DueDate = &dueDate
		assignment1.Closed = false
		SaveAssignment(t, assignment1)

		assignmentProject1.ProjectStatus = database.Pending
		err := query.AssignmentProjects.WithContext(context.Background()).Save(assignmentProject1)
		if err != nil {
			t.Fatalf("could not update assignment project: %s", err.Error())
		}

		err = work.closeAssignment(context.Background(), assignment1, repo)
		assert.NoError(t, err)

		assignment1After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment1.ID)).
			First()
		assert.NoError(t, err)
		assert.True(t, assignment1After.Closed)
	})

	assignmentProject1.ProjectStatus = database.Accepted
	assignmentProject1.Team = *team1
	assignment1.Projects = []*database.AssignmentProjects{assignmentProject1}

	t.Run("repo.GetAccessLevelOfUserInProject throws error -> restore old permissions", func(t *testing.T) {
		assignment1.Closed = false
		SaveAssignment(t, assignment1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, owner.ID).
			Return(model.OwnerPermissions, nil).
			Times(1)

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

		err := work.closeAssignment(context.Background(), assignment1, repo)
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
		assignment1.Closed = false
		SaveAssignment(t, assignment1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, owner.ID).
			Return(model.OwnerPermissions, nil).
			Times(1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, student1.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		repo.EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject1.ProjectID, student1.ID, model.ReporterPermissions).
			Return(assert.AnError).
			Times(1)

		err := work.closeAssignment(context.Background(), assignment1, repo)
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
		assignment1.Closed = false
		SaveAssignment(t, assignment1)

		repo.EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject1.ProjectID, owner.ID).
			Return(model.OwnerPermissions, nil).
			Times(1)

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

		err := work.closeAssignment(context.Background(), assignment1, repo)
		assert.NoError(t, err)

		repo.AssertExpectations(t)

		assignment1After, err := query.Assignment.
			WithContext(context.Background()).
			Where(query.Assignment.ID.Eq(assignment1.ID)).
			First()
		assert.NoError(t, err)
		assert.True(t, assignment1After.Closed)
	})
}

func SaveAssignment(t *testing.T, assignment *database.Assignment) {
	err := query.Assignment.WithContext(context.Background()).Save(assignment)
	if err != nil {
		t.Fatalf("could not update assignment: %s", err.Error())
	}
}

func SaveAssignmentProjects(t *testing.T, project *database.AssignmentProjects) {
	err := query.AssignmentProjects.WithContext(context.Background()).Save(project)
	if err != nil {
		t.Fatalf("could not update assignment project: %s", err.Error())
	}
}
