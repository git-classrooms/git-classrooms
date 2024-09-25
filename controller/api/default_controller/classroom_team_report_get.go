package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTeamReport
// @Description	GetClassroomTeamReport
// @Id				GetClassroomTeamReport
// @Tags			report
// @Produce		text/csv
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{file}		text/csv
// @Success		200			{array}		utils.ReportDataItem
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/teams/{teamId}/grading/report [get]
func (ctrl *DefaultController) GetClassroomTeamReport(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	assignments, err := assignmentGradingQuery(c, classroom.ClassroomID).
		Join(query.AssignmentProjects, query.AssignmentProjects.AssignmentID.EqCol(query.Assignment.ID)).
		Where(query.AssignmentProjects.TeamID.Eq(team.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if len(assignments) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No Assignments found")
	}

	acceptHeader := c.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		jsonReports, err := utils.GenerateReports(assignments, classroom.Classroom.ManualGradingRubrics, &team.ID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(jsonReports)
	}
	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=report_%s_%s_%s.csv", time.Now().Format(time.DateOnly), classroom.Classroom.Name, team.Name))
	return utils.GenerateCSVReports(c.Response().BodyWriter(), assignments, classroom.Classroom.ManualGradingRubrics, &team.ID)
}
