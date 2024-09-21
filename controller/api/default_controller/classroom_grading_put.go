package api

import (
	"database/sql/driver"

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
	return r.Name != "" && r.MaxScore > 0
}

func rubricRequestIsValid(r gradingManualRubricRequest) bool {
	return r.isValid()
}

type updateGradingRequest struct {
	GradingManualRubrics []gradingManualRubricRequest `json:"gradingManualRubrics"`
} //@Name UpdateGradingRequest

func (r updateGradingRequest) isValid() bool {
	return utils.All(r.GradingManualRubrics, rubricRequestIsValid)
}

// @Summary		UpdateGradingRubrics
// @Description	UpdateGradingRubrics
// @Id				UpdateGradingRubrics
// @Tags			grading
// @Accept			json
// @Param			classroomId		path	string						true	"Classroom ID"	Format(uuid)
// @Param			gradingInfo		body	api.updateGradingRequest	true	"Grading Update Info"
// @Param			X-Csrf-Token	header	string						true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/grading [put]
func (ctrl *DefaultController) UpdateGradingRubrics(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom().Classroom

	var requestBody updateGradingRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Request Body is not valid")
	}

	updateIDs := make([]driver.Valuer, 0)
	for _, e := range requestBody.GradingManualRubrics {
		if e.ID != nil {
			updateIDs = append(updateIDs, *e.ID)
		}
	}

	queryManualGradingRubric := query.ManualGradingRubric
	rubrics, err := queryManualGradingRubric.
		WithContext(c.Context()).
		Where(queryManualGradingRubric.ClassroomID.Eq(classroom.ID)).
		Where(queryManualGradingRubric.ID.In(updateIDs...)).Find()
	if err != nil {
		return err
	}

	if len(rubrics) != len(updateIDs) {
		return fiber.NewError(fiber.StatusBadRequest, "Body includes invalid IDs")
	}

	rubricMap := make(map[uuid.UUID]*database.ManualGradingRubric)
	for _, r := range rubrics {
		rubricMap[r.ID] = r
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		queryManualGradingRubric := tx.ManualGradingRubric
		if _, err := queryManualGradingRubric.
			WithContext(c.Context()).
			Where(queryManualGradingRubric.ClassroomID.Eq(classroom.ID)).
			Not(queryManualGradingRubric.ID.In(updateIDs...)).Delete(); err != nil {
			return err
		}

		for _, e := range requestBody.GradingManualRubrics {
			if e.ID != nil {
				rubric, ok := rubricMap[*e.ID]
				if !ok {
					return fiber.NewError(fiber.StatusBadRequest, "Body includes invalid IDs")
				}

				rubric.Name = e.Name
				rubric.Description = e.Description
				rubric.MaxScore = e.MaxScore
				if err := queryManualGradingRubric.
					WithContext(c.Context()).
					Save(rubric); err != nil {
					return err
				}
			} else {
				if err := queryManualGradingRubric.
					WithContext(c.Context()).
					Create(&database.ManualGradingRubric{
						Name:        e.Name,
						Description: e.Description,
						MaxScore:    e.MaxScore,
						ClassroomID: classroom.ID,
					}); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
