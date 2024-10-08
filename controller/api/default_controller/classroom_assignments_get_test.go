package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetClassroomAssignments(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	dueDate := time.Now().Add(1 * time.Hour)
	assignment := factory.Assignment(classroom.ID, &dueDate, false)

	// ------------ END OF SEEDING DATA -----------------


	app, _,_ := setupApp(t, owner)

	t.Run("GetClassroomAssignments", func(t *testing.T) {
		route := fmt.Sprintf("/api/v1/classrooms/%s/assignments", classroom.ID.String())

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

		var classroomAssignments []*ClassroomAssignmentResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomAssignments)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Len(t, classroomAssignments, 1)

		classroomAssignment := classroomAssignments[0]

		assert.Equal(t, assignment.ID, classroomAssignment.ID)
		assert.Equal(t, assignment.ClassroomID, classroomAssignment.ClassroomID)
		assert.Equal(t, assignment.TemplateProjectID, classroomAssignment.TemplateProjectID)
		assert.Equal(t, assignment.Name, classroomAssignment.Name)
		assert.Equal(t, assignment.Description, classroomAssignment.Description)
	})
}
