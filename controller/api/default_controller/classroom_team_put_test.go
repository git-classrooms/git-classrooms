package api

import (
	"context"
	"fmt"
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

func TestUpdateTeam(t *testing.T) {
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

	t.Run("TestUpdateTeam", func(t *testing.T) {
		requestBody := updateTeamRequest{
			Name: "New team name",
		}

		gitlabRepo.
			EXPECT().
			GroupAccessLogin(classroom.GroupAccessToken).
			Return(nil).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeGroupName(
				team.GroupID,
				requestBody.Name,
			).
			Return(nil, nil).
			Times(1)

		route := fmt.Sprintf("/api/v1/classrooms/%s/teams/%s", classroom.ID, team.ID)

		req := newPutJsonRequest(route, requestBody)
		resp, err := app.Test(req)
		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)

		assert.NoError(t, err)

		q := query.Team

		updatedTeam, err := q.
			WithContext(context.Background()).
			Where(q.ID.Eq(team.ID)).
			First()

		assert.Equal(t, updatedTeam.ID, team.ID)
		assert.Equal(t, updatedTeam.Name, requestBody.Name)
	})
}
