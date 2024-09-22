package api

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestPostClassroomAssignmentProjects(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	// Seed data
	owner := factory.User()
	user := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)
	userClassroom := factory.UserClassroom(user.ID, classroom.ID, database.Student)

	dueDate := time.Now().Add(1 * time.Hour)
	assignment := factory.Assignment(classroom.ID, &dueDate, false)
	factory.Team(classroom.ID, []*database.UserClassrooms{userClassroom})

	// setup app
	app, _, mockMailRepo := setupApp(t, owner)

	t.Run("PostOwnedClassroomAssignmentProjects", func(t *testing.T) {
		route := fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s/projects", classroom.ID.String(), assignment.ID.String())

		mockMailRepo.
			EXPECT().
			SendAssignmentNotification(
				user.GitlabEmail,
				fmt.Sprintf(`You were invited to a new Assigment "%s"`, classroom.Name),
				mock.Anything,
			).
			RunAndReturn(func(gitlabEmail, mailTitle string, notificationData mailRepo.AssignmentNotificationData) error {
				assert.Equal(t, notificationData.ClassroomName, classroom.Name)
				assert.Equal(t, notificationData.ClassroomOwnerName, owner.Name)
				assert.Equal(t, notificationData.RecipientName, user.Name)
				assert.Equal(t, notificationData.AssignmentName, assignment.Name)

				return nil
			})

		req := httptest.NewRequest("POST", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})
}
