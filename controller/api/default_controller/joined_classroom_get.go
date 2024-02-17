package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getJoinedClassroomResponse struct {
	database.UserClassrooms
	GitlabUrl string `json:"gitlabUrl"`
}

func (ctrl *DefaultController) GetJoinedClassroom(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetJoinedClassroom()
	repo := ctx.GetGitlabRepository()

	group, err := repo.GetGroupById(classroom.Classroom.GroupID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	response := &getJoinedClassroomResponse{
		UserClassrooms: *classroom,
		GitlabUrl:      group.WebUrl,
	}

	return c.JSON(response)
}
