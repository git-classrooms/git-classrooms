package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetJoinedClassrooms(c *fiber.Ctx) error {
	ctx := context.Get(c)
	userID := ctx.GetUserID()
	repo := ctx.GetGitlabRepository()

	joinedClassrooms, err := joinedClassroomQuery(userID, c).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var joinedClassroomResponses = make([]*getJoinedClassroomResponse, len(joinedClassrooms))
	for i, classroom := range joinedClassrooms {
		group, err := repo.GetGroupById(classroom.Classroom.GroupID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		joinedClassroomResponses[i] = &getJoinedClassroomResponse{
			UserClassrooms: *classroom,
			GitlabUrl:      group.WebUrl,
		}
	}

	return c.JSON(joinedClassroomResponses)
}
