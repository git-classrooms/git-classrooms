package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestJoinTeam(t *testing.T) {
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

	dueDate := time.Now().Add(1 * time.Hour)

	members := []*database.UserClassrooms{
		factory.UserClassroom(member.ID, classroom.ID, database.Student),
	}

	newMember := factory.User()
	factory.UserClassroom(newMember.ID, classroom.ID, database.Student)

	factory.Assignment(classroom.ID, &dueDate, false)
	team := factory.Team(classroom.ID, members)

	app, gitlabRepo, _ := setupApp(t, newMember)

	t.Run("TestJoinTeam", func(t *testing.T) {
		gitlabRepo.
			EXPECT().
			GroupAccessLogin(classroom.GroupAccessToken).
			Return(nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			AddUserToGroup(
				team.GroupID,
				newMember.ID,
				model.ReporterPermissions,
			).
			Return(nil).
			Times(1)

		route := fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/join", classroom.ID, team.ID)

		req := newPostJsonRequest(route, nil)
		resp, err := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assert.NoError(t, err)

		q := query.UserClassrooms

		updatedUserClassoom, err := q.
			WithContext(context.Background()).
			Where(q.UserID.Eq(newMember.ID)).
			Where(q.ClassroomID.Eq(classroom.ID)).
			First()

		assert.Equal(t, updatedUserClassoom.UserID, newMember.ID)
		assert.Equal(t, updatedUserClassoom.TeamID, &team.ID)
	})
}
