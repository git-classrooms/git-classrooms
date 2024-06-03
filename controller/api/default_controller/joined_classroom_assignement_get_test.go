package default_controller

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetJoinedClassroomAssignment(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	owner := &database.User{ID: 1, GitlabEmail: "owner@example.com"}
	testDB.InsertUser(owner)

	member := &database.User{ID: 2, GitlabEmail: "member@example.com"}
	testDB.InsertUser(member)

	classroom := factory.Classroom(map[string]any{"OwnerID": owner.ID})
	testDB.InsertClassroom(classroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(assignment)

	team := factory.Team(classroom.ID)
	testDB.InsertTeam(team)

	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(member.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetJoinedClassroomAssignment", func(t *testing.T) {
		app.Get("classrooms/joined/:classroomId/assignments/:assignmentId", handler.GetJoinedClassroomAssignment)
		route := fmt.Sprintf("/api/classrooms/joined/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

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
