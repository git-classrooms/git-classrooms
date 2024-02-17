package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomResponse struct {
	database.Classroom
	GitlabUrl string `json:"gitlabUrl"`
}

func (ctrl *DefaultController) GetOwnedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	repo := ctx.GetGitlabRepository()

	group, err := repo.GetGroupById(classroom.GroupID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	response := &getOwnedClassroomResponse{
		Classroom: *classroom,
		GitlabUrl: group.WebUrl,
	}

	return c.JSON(response)
}
