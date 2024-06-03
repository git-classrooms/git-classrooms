package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestJoinAssignment(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	user := factory.User()
	testDB.InsertUser(user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(classroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(assignment)

	team := factory.Team(classroom.ID)
	testDB.InsertTeam(team)

	project := factory.AssignmentProject(assignment.ID, team.ID)

	testDB.InsertAssignmentProjects(project)

	// ------------ END OF SEEDING DATA -----------------

	mailRepo := mailRepoMock.NewMockRepository(t)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetJoinedClassroom(&database.UserClassrooms{ClassroomID: classroom.ID})
		ctx.SetUserID(1)
		ctx.SetGitlabRepository(gitlabRepo)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("Join Assignment", func(t *testing.T) {
		app.Post("/api/v1/classrooms/joined/:classroomId/assignments/:assignmentId/accept", handler.AcceptAssignment)
		route := fmt.Sprintf("/api/v1/classrooms/joined/%s/assignments/%s/accept", classroom.ID.String(), assignment.ID.String())

		req := httptest.NewRequest("POST", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		assignmentProject, err :=
			query.AssignmentProjects.
			WithContext(context.Background()).
			Where(query.AssignmentProjects.AssignmentID.Eq(assignment.ID)).
			Where(query.AssignmentProjects.TeamID.Eq(team.ID)).First()

		assert.NoError(t, err)
		assert.True(t, assignmentProject.AssignmentAccepted)
	})
}
