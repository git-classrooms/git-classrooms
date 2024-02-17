package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassrooms(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	repo := ctx.GetGitlabRepository()

	ownedClassrooms, err := ownedClassroomQuery(userID, c).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var ownedClassroomResponses = make([]*getOwnedClassroomResponse, len(ownedClassrooms))
	for i, classroom := range ownedClassrooms {
		group, err := repo.GetGroupById(classroom.GroupID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		ownedClassroomResponses[i] = &getOwnedClassroomResponse{
			Classroom: *classroom,
			GitlabUrl: group.WebUrl,
		}
	}

	return c.JSON(ownedClassroomResponses)
}
