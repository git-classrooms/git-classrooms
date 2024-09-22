package api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestRemoveMemberFromTeam(t *testing.T) {
	// setup database
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	member := factory.User()
	removeMember := factory.User()

	dueDate := time.Now().Add(1 * time.Hour)

	members := []*database.UserClassrooms{
		factory.UserClassroom(member.ID, classroom.ID, database.Student),
		factory.UserClassroom(removeMember.ID, classroom.ID, database.Student),
	}

	factory.Assignment(classroom.ID, &dueDate, false)
	team := factory.Team(classroom.ID, members)

	app, gitlabRepo, _ := setupApp(t, owner)
	route := fmt.Sprintf("/api/v1/classrooms/%s/teams/%s/members/%d", classroom.ID, team.ID, removeMember.ID)

	t.Run("TestRemoveMemberFromTeam", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GroupAccessLogin(classroom.GroupAccessToken).
			Return(nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			RemoveUserFromGroup(
				team.GroupID,
				removeMember.ID,
			).
			Return(nil).
			Times(1)

		req := httptest.NewRequest("DELETE", route, nil)
		resp, err := app.Test(req)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		assert.NoError(t, err)

		q := query.UserClassrooms

		updatedUserClassoom, err := q.
			WithContext(context.Background()).
			Where(q.UserID.Eq(removeMember.ID)).
			Where(q.ClassroomID.Eq(classroom.ID)).
			First()

		assert.Equal(t, updatedUserClassoom.UserID, removeMember.ID)
		assert.Nil(t, updatedUserClassoom.TeamID)
	})
}
