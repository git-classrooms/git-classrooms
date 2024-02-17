package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
)

func (ctrl *DefaultController) GetClassroomAssignments(c *fiber.Ctx) error {
	userId, err := session.Get(c).GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	classroomId, err := uuid.Parse(c.Params("classroomId"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryClassroom := query.Classroom
	var classroom *database.UserClassrooms
	ownedClassroom, err := queryClassroom.
		WithContext(c.Context()).
		Where(queryClassroom.ID.Eq(classroomId)).
		Where(queryClassroom.OwnerID.Eq(userId)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if len(ownedClassroom) != 1 {
		queryUserClassroom := query.UserClassrooms
		joinedClassroom, err := queryUserClassroom.
			WithContext(c.Context()).
			Preload(queryUserClassroom.Classroom).
			Where(queryUserClassroom.ClassroomID.Eq(classroomId)).
			Where(queryUserClassroom.UserID.Eq(userId)).
			First()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		classroom = joinedClassroom
	} else {
		classroom = &database.UserClassrooms{
			Role:        database.Owner,
			Classroom:   *ownedClassroom[0],
			UserID:      userId,
			ClassroomID: classroomId,
		}
	}

	queryAssignment := query.Assignment
	assignments, err := queryAssignment.
		WithContext(c.Context()).
		Where(queryAssignment.ClassroomID.Eq(classroom.ClassroomID)).
		Find()

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(assignments)
}
