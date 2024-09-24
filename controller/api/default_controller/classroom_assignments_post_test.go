package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestPostClassroomAssignment(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	// Seed data
	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	dueDate := time.Now().Add(1 * time.Hour)

	// setup app
	app, gitlabRepo, _ := setupApp(t, owner)
	route := fmt.Sprintf("/api/v1/classrooms/%s/assignments", classroom.ID.String())

	t.Run("PostClassroomAssignment", func(t *testing.T) {
		requestBody := createAssignmentRequest{
			Name:              gofakeit.Name(),
			Description:       gofakeit.EmojiDescription(),
			TemplateProjectID: gofakeit.Int(),
			DueDate:           &dueDate,
		}

		gitlabRepo.
			EXPECT().GetProjectByID(requestBody.TemplateProjectID).Return(nil, nil)

		req := tests.NewPostJSONRequest(route, requestBody)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assignment, err := query.Assignment.WithContext(context.Background()).Where(query.Assignment.ClassroomID.Eq(classroom.ID)).First()
		assert.NoError(t, err)

		assert.NotNil(t, assignment)

		assert.Equal(t, assignment.Name, requestBody.Name)
		assert.Equal(t, assignment.Description, requestBody.Description)
		assert.Equal(t, assignment.TemplateProjectID, requestBody.TemplateProjectID)
		assert.WithinDuration(t, *assignment.DueDate, *requestBody.DueDate, 1*time.Minute)
	})
}
