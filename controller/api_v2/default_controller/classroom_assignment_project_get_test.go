package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetClassroomAssignmentProject(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	userClassroom := factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	assignment := factory.Assignment(classroom.ID)
	team := factory.Team(classroom.ID, []*database.UserClassrooms{userClassroom})
	project := factory.AssignmentProject(assignment.ID, team.ID)

	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(classroom)
		ctx.SetUserClassroom(userClassroom)
		ctx.SetAssignment(assignment)
		ctx.SetAssignmentProject(project)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(owner.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})

	t.Run("GetOwnedClassroomAssignmentProject", func(t *testing.T) {
		app.Get("/api/v2/classrooms/:classroomId/assignments/:assignmentId/projects/:projectId", handler.GetClassroomAssignmentProject)
		route := fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s/projects/%s", classroom.ID.String(), assignment.ID.String(), project.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		type ClassroomAssignmentProjectResponse struct {
			ID                 uuid.UUID `json:"id"`
			CreatedAt          time.Time `json:"createdAt"`
			UpdatedAt          time.Time `json:"updatedAt"`
			AssignmentID       uuid.UUID `json:"assignmentId"`
			UserID             int       `json:"userId"`
			AssignmentAccepted bool      `json:"assignmentAccepted"`
			ProjectID          int       `json:"projectId"`
		}

		var returnValue *ClassroomAssignmentProjectResponse

		err = json.NewDecoder(resp.Body).Decode(&returnValue)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, project.ID.String(), returnValue.ID.String())
		//assert.Equal(t, project.AssignmentID, returnValue.AssignmentID)
		//assert.Equal(t, project.Team.ID, returnValue.UserID)
		//assert.Equal(t, project.AssignmentAccepted, returnValue.AssignmentAccepted)
		//assert.Equal(t, project.ProjectID, returnValue.ProjectID)
	})
}
