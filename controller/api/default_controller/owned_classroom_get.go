package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomResponse struct {
	database.Classroom
	GitlabURL string `json:"gitlabUrl"`
}

func (ctrl *DefaultController) GetOwnedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	response := &getOwnedClassroomResponse{
		Classroom: *classroom,
		GitlabURL: fmt.Sprintf("/api/v1/classrooms/owned/%s/gitlab", classroom.ID.String()),
	}

	return c.JSON(response)
}
