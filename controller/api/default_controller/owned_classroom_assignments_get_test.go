package default_controller

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetOwnedClassroomAssignments(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	// ------------ END OF DB SETUP -----------------

	owner := &database.User{ID: 1, GitlabEmail: "owner@example.com"}
	testDB.InsertUser(owner)

	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	testClassroomAssignments := []*database.Assignment{
		{
			ClassroomID:       classroom.ID,
			TemplateProjectID: 1,
			Name:              "Test Assignment 1",
			Description:       "Description 1",
		},
		{
			ClassroomID:       classroom.ID,
			TemplateProjectID: 2,
			Name:              "Test Assignment 2",
			Description:       "Description 2",
		},
	}

	for _, a := range testClassroomAssignments {
		testDB.InsertAssignment(a)
	}

	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(&classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(owner.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignments", func(t *testing.T) {
		app.Get("/classrooms/owned/:classroomId/assignments", handler.GetOwnedClassroomAssignments)
		route := fmt.Sprintf("/api/classrooms/owned/%s/assignments", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		var classroomAssignments []*database.Assignment
		err = json.NewDecoder(resp.Body).Decode(&classroomAssignments)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Len(t, classroomAssignments, len(testClassroomAssignments))

		for i, assignment := range classroomAssignments {
			assert.Equal(t, testClassroomAssignments[i].ID, assignment.ID)
			assert.Equal(t, testClassroomAssignments[i].ClassroomID, assignment.ClassroomID)
			assert.Equal(t, testClassroomAssignments[i].TemplateProjectID, assignment.TemplateProjectID)
			assert.Equal(t, testClassroomAssignments[i].Name, assignment.Name)
			assert.Equal(t, testClassroomAssignments[i].Description, assignment.Description)
		}
	})
}

