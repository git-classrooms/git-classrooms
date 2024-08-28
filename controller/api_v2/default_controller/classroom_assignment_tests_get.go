package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/examples"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type assignmentTestResponse struct {
	Activatible bool                       `json:"activatible"`
	PipelineRan bool                       `json:"pipelineRan"`
	Example     examples.LanguageCIExample `json:"example"`
	Report      *model.TestReport          `json:"report"`
} // @Name AssignmentTestResponse

// @Summary		GetClassroomAssignmentTests
// @Description	GetClassroomAssignmentTests
// @Id				GetClassroomAssignmentTests
// @Tags			grading
// @Produce		json
// @Param			classroomId		path		string	true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path		string	true	"Assignment ID"	Format(uuid)
// @Success		200				{object}	AssignmentTestResponse
// @Failure		400				{object}	HTTPError
// @Failure		401				{object}	HTTPError
// @Failure		403				{object}	HTTPError
// @Failure		404				{object}	HTTPError
// @Failure		500				{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/tests [get]
func (ctrl *DefaultController) GetClassroomAssignmentTests(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()
	assignment := ctx.GetAssignment()

	languages, err := repo.GetProjectLanguages(assignment.TemplateProjectID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var maxValue float32
	var maxLanguage string
	for language, percent := range languages {
		if percent > maxValue {
			maxValue = percent
			maxLanguage = language
		}
	}

	var response assignmentTestResponse
	response.Example = examples.GetLanguageCIExample(maxLanguage)
	response.Activatible, err = repo.CheckIfFileExistsInProject(assignment.TemplateProjectID, ".gitlab-ci.yml")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !response.Activatible {
		return c.JSON(response)
	}

	response.Report, err = repo.GetProjectLatestPipelineTestReportSummary(assignment.TemplateProjectID, nil)
	response.PipelineRan = err == nil
	return c.JSON(response)
}
