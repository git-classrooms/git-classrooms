package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassrooms(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	ownedClassrooms, err := ownedClassroomQuery(userID, c).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var ownedClassroomResponses = make([]*getOwnedClassroomResponse, len(ownedClassrooms))
	for i, classroom := range ownedClassrooms {
		ownedClassroomResponses[i] = &getOwnedClassroomResponse{
			Classroom: *classroom,
			GitlabUrl: fmt.Sprintf("/api/v1/classrooms/owned/%s/gitlab", classroom.ID.String()),
		}
	}

	return c.JSON(ownedClassroomResponses)
}
