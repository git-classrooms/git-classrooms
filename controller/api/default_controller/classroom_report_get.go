package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gen/field"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomReport
// @Description	GetClassroomReport
// @Id				GetClassroomReport
// @Tags			report
// @Produce		text/csv
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{file}		text/csv
// @Success		200			{array}		utils.ReportDataItem
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/grading/report [get]
func (ctrl *DefaultController) GetClassroomReport(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()

	assignments, err := assignmentGradingQuery(c, classroom.ClassroomID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if len(assignments) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No assignments found")
	}

	acceptHeader := c.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		jsonReports, err := utils.GenerateReports(assignments, classroom.Classroom.ManualGradingRubrics, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(jsonReports)
	}
	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=report_%s_%s_%s.csv", time.Now().Format(time.DateOnly), classroom.Classroom.Name, "all"))

	return utils.GenerateCSVReports(c.Response().BodyWriter(), assignments, classroom.Classroom.ManualGradingRubrics, nil)
}

func assignmentGradingQuery(c *fiber.Ctx, classroomID uuid.UUID) query.IAssignmentDo {
	queryAssignment := query.Assignment
	return queryAssignment.
		WithContext(c.Context()).
		Preload(queryAssignment.Projects).
		Preload(queryAssignment.GradingManualRubrics).
		Preload(queryAssignment.JUnitTests).
		Preload(queryAssignment.Projects.Team).
		Preload(queryAssignment.Projects.Team.Member).
		Preload(field.NewRelation("Projects.Team.Member.User", "")).
		Preload(queryAssignment.Projects.GradingManualResults).
		Preload(field.NewRelation("Projects.GradingManualResults.Rubric", "")).
		Where(queryAssignment.ClassroomID.Eq(classroomID))
}
