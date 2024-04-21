package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomMemberResponse struct {
	*database.UserClassrooms
	GitlabURL string `json:"gitlabUrl"`
}

func (ctrl *DefaultController) GetOwnedClassroomMember(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	member := ctx.GetOwnedClassroomMember()
	response := &getOwnedClassroomMemberResponse{
		UserClassrooms: member,
		GitlabURL:      fmt.Sprintf("/api/v1/classrooms/owned/%s/members/%d/gitlab", classroom.ID.String(), member.UserID),
	}

	return c.JSON(response)
}
