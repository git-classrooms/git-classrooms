package api

//import (
//	"fmt"
//	"net/http/httptest"
//	"testing"
//
//	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
//	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
//	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
//	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
//
//	"github.com/gofiber/fiber/v2"
//	"github.com/stretchr/testify/assert"
//)

//func TestClassroomAssignmentMiddleware(t *testing.T) {
//	IntegrationTestCleanUp(t)
//
//	owner := factory.User()
//	classroom := factory.Classroom(owner.ID)
//	assignment := factory.Assignment(classroom.ID)
//
//	mailRepo := mailRepoMock.NewMockRepository(t)
//
//	app := fiber.New()
//	app.Use("/api", func(c *fiber.Ctx) error {
//		fctx := fiberContext.Get(c)
//		fctx.SetOwnedClassroom(&classroom)
//
//		s := session.Get(c)
//		s.SetUserState(session.LoggedIn)
//		s.SetUserID(1)
//		s.Save()
//		return c.Next()
//	})
//
//	handler := NewApiV2Controller(mailRepo)
//
//	t.Run("ClassroomAssignmentMiddleware", func(t *testing.T) {
//		app.Use("/api/v2/classrooms/:classroomId/assignments/:assignmentId", handler.ClassroomAssignmentMiddleware)
//
//		app.Get("/api/v2/classrooms/:classroomId/assignments/:assignmentId", func(c *fiber.Ctx) error {
//			ctx := fiberContext.Get(c)
//			handlerAssignment := ctx.GetAssignment()
//
//			assert.Equal(t, assignment.Name, handlerAssignment.Name)
//
//			return c.JSON(nil)
//		})
//
//		route := fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())
//
//		req := httptest.NewRequest("GET", route, nil)
//		resp, err := app.Test(req)
//
//		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
//		assert.NoError(t, err)
//	})
//}
