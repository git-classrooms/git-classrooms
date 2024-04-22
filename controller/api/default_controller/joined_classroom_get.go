package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getJoinedClassroomResponse struct {
	database.UserClassrooms
	GitlabURL string `json:"gitlabUrl"`
}

func (ctrl *DefaultController) GetJoinedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()

	response := &getJoinedClassroomResponse{
		UserClassrooms: *classroom,
		GitlabURL:      fmt.Sprintf("/api/v1/classrooms/joined/%s/gitlab", classroom.ClassroomID.String()),
	}

	return c.JSON(response)
}
