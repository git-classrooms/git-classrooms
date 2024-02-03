package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
)

func (ctrl *DefaultController) GetMeClassrooms(c *fiber.Ctx) error {
	userId, err := session.Get(c).GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryUser := query.User
	user, err := queryUser.WithContext(c.Context()).Preload().Preload(queryUser.OwnedClassrooms).Where(queryUser.ID.Eq(userId)).First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	queryClassrooms := query.UserClassrooms
	joinedClassrooms, err := query.UserClassrooms.WithContext(c.Context()).Preload(queryClassrooms.Classroom).Where(queryClassrooms.UserID.Eq(userId)).Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ownedClassrooms := utils.Map(user.OwnedClassrooms, func(classroom *database.Classroom) *database.UserClassrooms {
		return &database.UserClassrooms{
			Role:        database.Owner,
			Classroom:   *classroom,
			UserID:      userId,
			ClassroomID: classroom.ID,
		}
	})

	return c.JSON(fiber.Map{
		"ownClassrooms":    ownedClassrooms,
		"joinedClassrooms": joinedClassrooms,
	})
}
