package default_controller

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
)

type testJoinedClassroomAssignmentResponse struct {
	AssignmentProjects *database.AssignmentProjects `json:"assignmentProjects"`
	ProjectPath        string                       `json:"projectPath"`
}

// func testJoinedClassroomAssignmentQuery(classroomID uuid.UUID, c *fiber.Ctx) *gorm.DB {
// 	// Mock implementation of the query function.
// 	db := query.DB()
// 	return db.Where("classroom_id = ?", classroomID)
// }

func TestGetJoinedClassroomAssignments(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	// Create test user
	user := factory.User()
	testDB.InsertUser(&user)

	// Create test classroom
	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(&assignment)

	project := factory.AssignmentProject(assignment.ID, classroom.ID)
	testDB.InsertAssignmentProjects(&project)

	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetJoinedClassroom(&database.UserClassrooms{Classroom: classroom})

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetJoinedClassroomAssignments", func(t *testing.T) {
		app.Get("/api/classrooms/joined/:classroomId/assignments", handler.GetJoinedClassroomAssignments)
		route := fmt.Sprintf("/api/classrooms/joined/%s/assignments", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		var responses []*testJoinedClassroomAssignmentResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		assert.NoError(t, err)

		assert.Len(t, responses, 1)
		assert.Equal(t, project.ID, responses[0].AssignmentProjects.ID)
		assert.Equal(t, "/api/v1/classrooms/owned/1/gitlab", responses[0].ProjectPath)
	})
}
