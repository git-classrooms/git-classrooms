package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/gen/field"
	"time"
)

// @Summary		GetClassroomReport
// @Description	GetClassroomReport
// @Id				GetClassroomReport
// @Tags			report
// @Produce		application/zip
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Success		200			{file}		application/zip
// @Success		200			{array}		utils.ReportDataItem
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/grading/report [get]
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
	if acceptHeader == "application/json" {
		jsonReports, err := utils.GenerateReports(assignments, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(jsonReports)

	}

	c.Set(fiber.HeaderContentType, "application/zip")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=report_%s_%s_%s.zip", time.Now().Format(time.DateOnly), classroom.Classroom.Name, "all"))

	return utils.GenerateCSVReports(c.Response().BodyWriter(), assignments, nil)
}

func assignmentGradingQuery(c *fiber.Ctx, classroomID uuid.UUID) query.IAssignmentDo {
	queryAssignment := query.Assignment
	return queryAssignment.
		WithContext(c.Context()).
		Preload(queryAssignment.Projects).
		Preload(queryAssignment.Projects.Team).
		Preload(queryAssignment.Projects.Team.Member).
		Preload(field.NewRelation("Projects.Team.Member.User", "")).
		Preload(queryAssignment.Projects.GradingManualResults).
		Preload(field.NewRelation("Projects.GradingManualResults.Rubric", "")).
		Preload(queryAssignment.Projects.GradingJUnitTestResult).
		Where(queryAssignment.ClassroomID.Eq(classroomID))
}
