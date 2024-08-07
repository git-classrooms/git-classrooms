package api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	contextWrapper "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func TestPatchClassroomArchive(t *testing.T) {
	testDb := db_tests.NewTestDB(t)

	user1 := database.User{
		ID:             1,
		GitlabUsername: "user1",
		GitlabEmail:    "user1",
	}
	testDb.InsertUser(&user1)

	user2 := database.User{
		ID:             2,
		GitlabUsername: "user2",
		GitlabEmail:    "user2",
	}
	testDb.InsertUser(&user2)

	user3 := database.User{
		ID:             3,
		GitlabUsername: "user3",
		GitlabEmail:    "user3",
	}
	testDb.InsertUser(&user3)

	members := []*database.UserClassrooms{
		{UserID: user1.ID, Role: database.Owner},
		{UserID: user2.ID, Role: database.Student},
		{UserID: user3.ID, Role: database.Student},
	}

	classroom := database.Classroom{
		ID:       uuid.New(),
		OwnerID:  user1.ID,
		Member:   members,
		Archived: false,
		GroupID:  12,
	}
	testDb.InsertClassroom(&classroom)

	assignment := database.Assignment{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
	}
	testDb.InsertAssignment(&assignment)

	team := database.Team{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
		Member:      members,
	}
	testDb.InsertTeam(&team)

	assignmentProject := database.AssignmentProjects{
		ID:           uuid.New(),
		TeamID:       team.ID,
		AssignmentID: assignment.ID,
		ProjectID:    1,
	}
	testDb.InsertAssignmentProjects(&assignmentProject)

	userClassroom := database.UserClassrooms{
		UserID:      user1.ID,
		User:        user1,
		ClassroomID: classroom.ID,
		Classroom:   classroom,
	}

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := contextWrapper.Get(c)
		ctx.SetUserClassroom(&userClassroom)
		ctx.SetGitlabRepository(gitlabRepo)

		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})
	app.Patch("/api/classrooms/:classroomId/archive", handler.ArchiveClassroom)

	targetRoute := fmt.Sprintf("/api/classrooms/%s/archive", classroom.ID.String())

	t.Run("classroom already archived", func(t *testing.T) {
		userClassroom.Classroom.Archived = true
		req := httptest.NewRequest("PATCH", targetRoute, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		assert.NoError(t, err)
	})

	t.Run("gitlab throws error in changing access level", func(t *testing.T) {
		userClassroom.Classroom.Archived = false
		gitlabRepo.
			EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject.ProjectID, user2.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject.ProjectID, user2.ID, model.ReporterPermissions).
			Return(nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject.ProjectID, user3.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject.ProjectID, user3.ID, model.ReporterPermissions).
			Return(fmt.Errorf("error")).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject.ProjectID, user2.ID, model.DeveloperPermissions).
			Return(nil).
			Times(1)

		req := httptest.NewRequest("PATCH", targetRoute, nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		assert.NoError(t, err)
	})

	t.Run("updates classroom in db", func(t *testing.T) {
		userClassroom.Classroom.Archived = false
		gitlabRepo.
			EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject.ProjectID, user2.ID).
			Return(model.DeveloperPermissions, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject.ProjectID, user2.ID, model.ReporterPermissions).
			Return(nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			GetAccessLevelOfUserInProject(assignmentProject.ProjectID, user3.ID).
			Return(model.NoPermissions, nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeUserAccessLevelInProject(assignmentProject.ProjectID, user3.ID, model.ReporterPermissions).
			Return(nil).
			Times(1)

		req := httptest.NewRequest("PATCH", targetRoute, nil)
		resp, err := app.Test(req)

		gitlabRepo.AssertExpectations(t)

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		dbClassroom, err := query.Classroom.WithContext(context.Background()).Where(query.Classroom.ID.Eq(classroom.ID)).First()
		assert.NoError(t, err)
		assert.Equal(t, true, dbClassroom.Archived)
	})

}
