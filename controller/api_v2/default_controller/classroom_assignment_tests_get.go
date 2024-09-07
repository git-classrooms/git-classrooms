package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/examples"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type assignmentTestReport struct {
	Name      string `json:"name"`
	TestName  string `json:"testName"`
	TestSuite string `json:"testSuite"`
} // @Name AssignmentTestReport

type assignmentTestResponse struct {
	Activatible   bool                            `json:"activatible"`
	Example       examples.LanguageCIExample      `json:"example"`
	Report        []*assignmentTestReport         `json:"report"`
	SelectedTests []*database.AssignmentJunitTest `json:"selectedTests"`
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

	response := &assignmentTestResponse{
		SelectedTests: make([]*database.AssignmentJunitTest, 0),
		Report:        make([]*assignmentTestReport, 0),
		Example:       examples.GetLanguageCIExample(maxLanguage),
	}
	response.Activatible, err = repo.CheckIfFileExistsInProject(assignment.TemplateProjectID, ".gitlab-ci.yml")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !response.Activatible {
		return c.JSON(response)
	}

	report, err := repo.GetProjectLatestPipelineTestReportSummary(assignment.TemplateProjectID, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response.Report = utils.FlatMap(report.TestSuites, func(ts model.TestReportTestSuite) []*assignmentTestReport {
		return utils.Map(ts.TestCases, func(tc model.TestReportTestCase) *assignmentTestReport {
			return &assignmentTestReport{
				Name:      fmt.Sprintf("%s/%s", ts.Name, tc.Name),
				TestName:  tc.Name,
				TestSuite: ts.Name,
			}
		})
	})

	response.SelectedTests = assignment.JUnitTests

	return c.JSON(response)
}
