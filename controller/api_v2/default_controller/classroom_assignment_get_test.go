package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
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

func TestGetClassroomAssignment(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	assignment := factory.Assignment(classroom.ID)

	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(classroom)
		ctx.SetAssignment(assignment)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(owner.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})

	t.Run("GetClassroomAssignment", func(t *testing.T) {
		app.Get("/api/v2/classrooms/:classroomId/assignments/:assignmentId", handler.GetClassroomAssignment)
		route := fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		type ClassroomAssignmentResponse struct {
			ID                uuid.UUID  `json:"id"`
			CreatedAt         time.Time  `json:"createdAt"`
			UpdatedAt         time.Time  `json:"updatedAt"`
			ClassroomID       uuid.UUID  `json:"classroomId"`
			TemplateProjectID int        `json:"templateProjectId"`
			Name              string     `json:"name"`
			Description       string     `json:"description"`
			DueDate           *time.Time `json:"dueDate"`
		}

		var classroomAssignment *ClassroomAssignmentResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomAssignment)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, assignment.ID, classroomAssignment.ID)
		assert.Equal(t, assignment.ClassroomID, classroomAssignment.ClassroomID)
		assert.Equal(t, assignment.TemplateProjectID, classroomAssignment.TemplateProjectID)
		assert.Equal(t, assignment.Name, classroomAssignment.Name)
		assert.Equal(t, assignment.Description, classroomAssignment.Description)
	})
}
