package api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestPatchClassroomArchive(t *testing.T) {
	restoreDatabase(t)

	owner := factory.User()
	user2 := factory.User()
	user3 := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	members := []*database.UserClassrooms{
		factory.UserClassroom(user2.ID, classroom.ID, database.Student),
		factory.UserClassroom(user3.ID, classroom.ID, database.Student),
	}

	dueDate := time.Now().Add(1 * time.Hour)

	assignment := factory.Assignment(classroom.ID, &dueDate, false)
	team := factory.Team(classroom.ID, members)
	assignmentProject := factory.AssignmentProject(assignment.ID, team.ID)

	app, gitlabRepo, _ := setupApp(t, owner)
	targetRoute := fmt.Sprintf("/api/v1/classrooms/%s/archive", classroom.ID.String())

	t.Run("classroom already archived", func(t *testing.T) {
		classroom.Archived = true
		saveClassroom(t, classroom)

		req := httptest.NewRequest("PATCH", targetRoute, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
		assert.NoError(t, err)
	})

	t.Run("gitlab throws error in changing access level", func(t *testing.T) {
		classroom.Archived = false
		saveClassroom(t, classroom)

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
		classroom.Archived = false
		saveClassroom(t, classroom)

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
