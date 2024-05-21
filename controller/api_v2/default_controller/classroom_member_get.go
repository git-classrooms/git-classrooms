package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomMember(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	member := ctx.GetClassroomMember()

	response := &UserClassroomResponse{
		UserClassrooms: member,
		WebURL:         fmt.Sprintf("/api/v2/classrooms/%s/members/%d", classroom.ClassroomID.String(), member.UserID),
	}

	return c.JSON(response)
}
