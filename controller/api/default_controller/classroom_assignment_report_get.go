package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomAssignmentReport
// @Description	GetClassroomAssignmentReport
// @Id				GetClassroomAssignmentReport
// @Tags			report
// @Produce		text/csv
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{file}		text/csv
// @Success		200				{array}		utils.ReportDataItem
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentId}/grading/report [get]
func (ctrl *DefaultController) GetClassroomAssignmentReport(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	assignment := ctx.GetAssignment()

	reportAssignment, err := assignmentGradingQuery(c, classroom.ClassroomID).
		Where(query.Assignment.ID.Eq(assignment.ID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Assignment not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	acceptHeader := c.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		jsonReport, err := utils.GenerateReport(reportAssignment, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(jsonReport)
	}

	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=report_%s_%s_%s.csv", time.Now().Format(time.DateOnly), classroom.Classroom.Name, assignment.Name))

	return utils.GenerateCSVReport(c.Response().BodyWriter(), reportAssignment, reportAssignment.GradingManualRubrics, nil, true)
}
