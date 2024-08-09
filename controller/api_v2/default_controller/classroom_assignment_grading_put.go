package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type gradingManualRubricRequest struct {
	ID          *uuid.UUID `json:"id" validate:"optional"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	MaxScore    int        `json:"maxScore"`
} //@Name GradingManualRubricRequest

func (r gradingManualRubricRequest) isValid() bool {
	return r.Name != "" && r.Description != "" && r.MaxScore >= 0
}

func rubricRequestIsValid(r gradingManualRubricRequest) bool {
	return r.isValid()
}

type updateAssignmentGradingRequest struct {
	GradingJUnitAutoGradingActive *bool `json:"gradingJUnitAutoGradingActive"`

	GradingManualRubrics []gradingManualRubricRequest `json:"gradingManualRubrics"`
} //@Name UpdateAssignmentGradingRequest

func (r updateAssignmentGradingRequest) isValid() bool {
	return r.GradingJUnitAutoGradingActive != nil && utils.All(r.GradingManualRubrics, rubricRequestIsValid)
}

// @Summary		UpdateGradingRubrics
// @Description	UpdateGradingRubrics
// @Id				UpdateGradingRubrics
// @Tags			grading
// @Accept			json
// @Param			classroomId		path	string								true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string								true	"Assignment ID"	Format(uuid)
// @Param			assignmentInfo	body	api.updateAssignmentGradingRequest	true	"Grading Update Info"
// @Param			X-Csrf-Token	header	string								true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/grading [put]
func (ctrl *DefaultController) UpdateGradingRubrics(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()

	var requestBody updateAssignmentGradingRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Request Body is not valid")
	}

	updateIDs := make([]uuid.UUID, 0)
	for _, e := range requestBody.GradingManualRubrics {
		if e.ID != nil {
			updateIDs = append(updateIDs, *e.ID)
		}
	}

	queryManualGradingRubric := query.ManualGradingRubric
	rubrics, err := queryManualGradingRubric.
		WithContext(c.Context()).
		FindByAssignmentIDAndInRubricIDs(assignment.ID, updateIDs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(rubrics) != len(updateIDs) {
		return fiber.NewError(fiber.StatusBadRequest, "Body includes invalid IDs")
	}

	updatedRubrics := utils.Map(requestBody.GradingManualRubrics, func(e gradingManualRubricRequest) *database.ManualGradingRubric {
		id := uuid.UUID{}
		if e.ID != nil {
			id = *e.ID
		}

		return &database.ManualGradingRubric{
			ID:           id,
			Name:         e.Name,
			Description:  e.Description,
			AssignmentID: assignment.ID,
			MaxScore:     e.MaxScore,
		}
	})

	if err = queryManualGradingRubric.
		WithContext(c.Context()).
		Save(updatedRubrics...); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newIDs := utils.Map(updatedRubrics, func(e *database.ManualGradingRubric) uuid.UUID {
		return e.ID
	})

	if err = queryManualGradingRubric.
		WithContext(c.Context()).
		DeleteByAssignmentIDAndNotInRubricIDs(assignment.ID, newIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	assignment.GradingJUnitAutoGradingActive = *requestBody.GradingJUnitAutoGradingActive
	if err = query.Assignment.WithContext(c.Context()).Save(assignment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
