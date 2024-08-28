package api

import (
	"fmt"
	"slices"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/gen/field"
)

type assignmentTestRequest struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
} //@Name AssignmentTestRequest

func (r assignmentTestRequest) isValid() bool {
	return r.Name != "" && r.Score > 0
}

func assignmentTestRequestIsValid(r assignmentTestRequest) bool {
	return r.isValid()
}

type updateAssignmentTestRequest struct {
	AssignmentTests []assignmentTestRequest `json:"assignmentTests"`
} //@Name UpdateAssignmentTestRequest

func (r updateAssignmentTestRequest) isValid() bool {
	return utils.All(r.AssignmentTests, assignmentTestRequestIsValid)
}

// @Summary		UpdateAssignmentTests
// @Description	UpdateAssignmentTests
// @Id				UpdateAssignmentTests
// @Tags			grading
// @Accept			json
// @Param			classroomId			path	string						true	"Classroom ID"	Format(uuid)
// @Param			assignmentId		path	string						true	"Assignment ID"	Format(uuid)
// @Param			assignmentTestInfo	body	UpdateAssignmentTestRequest	true	"Assignment Test Update Info"
// @Param			X-Csrf-Token		header	string						true	"Csrf-Token"
// @Success		202
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/assignments/{assignmentId}/tests [put]
func (ctrl *DefaultController) UpdateAssignmentTests(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	repo := ctx.GetGitlabRepository()
	assignment := ctx.GetAssignment()

	var requestBody updateAssignmentTestRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.NewError(fiber.StatusBadRequest, "Request Body is not valid")
	}

	report, err := repo.GetProjectLatestPipelineTestReportSummary(assignment.TemplateProjectID, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	testNames := utils.FlatMap(report.TestSuites, func(ts model.TestReportTestSuite) []string {
		return utils.Map(ts.TestCases, func(tc model.TestReportTestCase) string {
			return fmt.Sprintf("%s/%s", ts.Name, tc.Name)
		})
	})

	names := utils.Map(requestBody.AssignmentTests, func(e assignmentTestRequest) string { return e.Name })

	if !utils.All(names, func(e string) bool { return slices.Contains(testNames, e) }) {
		return fiber.NewError(fiber.StatusBadRequest, "Body includes invalid test names")
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		queryAssignmentJunitTest := tx.AssignmentJunitTest

		if _, err := queryAssignmentJunitTest.
			WithContext(c.Context()).
			Where(queryAssignmentJunitTest.AssignmentID.Eq(assignment.ID)).
			Not(queryAssignmentJunitTest.Name.In(names...)).
			Delete(); err != nil {
			return err
		}

		for _, e := range requestBody.AssignmentTests {
			if _, err := queryAssignmentJunitTest.
				WithContext(c.Context()).
				Assign(field.Attrs(&database.AssignmentJunitTest{Score: e.Score})).
				Where(queryAssignmentJunitTest.AssignmentID.Eq(assignment.ID)).
				Where(queryAssignmentJunitTest.Name.Eq(e.Name)).
				FirstOrCreate(); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
