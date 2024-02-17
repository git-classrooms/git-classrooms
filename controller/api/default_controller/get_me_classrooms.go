package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"gorm.io/gen/field"
)

type getMeClassroomsClassroom struct {
	database.UserClassrooms
	GitlabUrl string `json:"gitlabUrl"`
}

func (ctrl *DefaultController) GetMeClassrooms(c *fiber.Ctx) error {
	userId, err := session.Get(c).GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryUser := query.User
	user, err := queryUser.WithContext(c.Context()).
		Preload(queryUser.OwnedClassrooms.Owner).
		Preload(queryUser.Classrooms).
		Where(queryUser.ID.Eq(userId)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	queryUserClassrooms := query.UserClassrooms
	joinedClassrooms, err := queryUserClassrooms.
		WithContext(c.Context()).
		Preload(queryUserClassrooms.Classroom).
		Preload(field.NewRelation("Classroom.Owner", "")).
		Where(queryUserClassrooms.UserID.Eq(userId)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	repo := context.Get(c).GetGitlabRepository()
	ownedClassroomResponses := make([]*getMeClassroomsClassroom, len(user.OwnedClassrooms))
	for i, owned := range user.OwnedClassrooms {
		group, err := repo.GetGroupById(owned.GroupID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		ownedClassroomResponses[i] = &getMeClassroomsClassroom{
			UserClassrooms: database.UserClassrooms{
				Role:        database.Owner,
				Classroom:   *owned,
				UserID:      userId,
				ClassroomID: owned.ID,
			},
			GitlabUrl: group.WebUrl,
		}
	}

	joinedClassroomResponses := make([]*getMeClassroomsClassroom, len(joinedClassrooms))
	for i, joinedClassroom := range joinedClassrooms {
		group, err := repo.GetGroupById(joinedClassroom.Classroom.GroupID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		joinedClassroomResponses[i] = &getMeClassroomsClassroom{
			UserClassrooms: *joinedClassroom,
			GitlabUrl:      group.WebUrl,
		}
	}

	return c.JSON(fiber.Map{
		"ownClassrooms":    ownedClassroomResponses,
		"joinedClassrooms": joinedClassroomResponses,
	})
}
