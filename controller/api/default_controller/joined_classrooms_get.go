package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetJoinedClassrooms(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()

	joinedClassrooms, err := joinedClassroomQuery(userID, c).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var joinedClassroomResponses = make([]*getJoinedClassroomResponse, len(joinedClassrooms))
	for i, classroom := range joinedClassrooms {
		joinedClassroomResponses[i] = &getJoinedClassroomResponse{
			UserClassrooms: *classroom,
			GitlabUrl:      fmt.Sprintf("/api/v1/classrooms/joined/%s/gitlab", classroom.ClassroomID.String()),
		}
	}

	return c.JSON(joinedClassroomResponses)
}
