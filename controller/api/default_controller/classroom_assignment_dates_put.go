package api

import (
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateAssignmentDate struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"duedate"`
}

func (r updateAssignmentDate) isValid() bool {
	return r.DueDate != nil && r.DueDate.After(time.Now()) && r.Description != ""
}

type updateAssignmentDatesRequest struct {
	Dates []updateAssignmentDate `json:"dates"`
} //@Name UpdateAssignmentDateRequest

func (r updateAssignmentDatesRequest) isValid() bool {
	return len(r.Dates) > 0 && utils.All(r.Dates, func(u updateAssignmentDate) bool { return u.isValid() })
}

// @Summary		UpdateAssignmentDates
// @Description	UpdateAssignmentDates
// @Id				UpdateAssignmentDates
// @Tags			assignment
// @Accept			json
// @Param			classroomId		path	string								true	"Classroom ID"	Format(uuid)
// @Param			assignmentId	path	string								true	"Assignment ID"	Format(uuid)
// @Param			assignmentInfo	body	api.updateAssignmentDatesRequest	true	"Assignment Date Update Info"
// @Param			X-Csrf-Token	header	string								true	"Csrf-Token"
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/assignments/{assignmentId}/dates [put]
func (ctrl *DefaultController) UpdateAssignmentDates(c *fiber.Ctx) error {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()
	var err error

	requestBody := &updateAssignmentDatesRequest{}
	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	sortedUpdateDates := slices.SortedFunc(slices.Values(requestBody.Dates), func(a, b updateAssignmentDate) int { return a.DueDate.Compare(*b.DueDate) })

	allClosed := utils.All(assignment.AssignmentDates, func(a *database.AssignmentDate) bool { return a.Closed })

	tx := query.Q.Begin()
	defer tx.Rollback()

	var createdNew bool
	for _, update := range sortedUpdateDates {
		if update.ID == nil {
			assignmentDate := database.AssignmentDate{
				AssignmentID: assignment.ID,
				Description:  update.Description,
				DueDate:      *update.DueDate,
				AssignmentProjectGradingDate: utils.Map(assignment.Projects, func(p *database.AssignmentProjects) *database.AssignmentProjectGradingDate {
					return &database.AssignmentProjectGradingDate{AssignmentProjectID: p.ID}
				}),
			}

			if err = tx.AssignmentDate.WithContext(c.Context()).Create(&assignmentDate); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			createdNew = true
		}

		idx := slices.IndexFunc(assignment.AssignmentDates, func(a *database.AssignmentDate) bool { return *update.ID == a.ID })
		if idx < 0 {
			return fiber.NewError(fiber.StatusNotFound, "assignmentDate not found")
		}

		assignmentDate := assignment.AssignmentDates[idx]
		if assignmentDate.Closed {
			return fiber.NewError(fiber.StatusBadRequest, "assignmentDate already closed")
		}

		assignmentDate.DueDate = *update.DueDate
		assignmentDate.Description = update.Description
		if err = tx.AssignmentDate.WithContext(c.Context()).Save(assignmentDate); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	if allClosed && createdNew {
		ctrl.reopenAssignment(c)
	}

	err = tx.Commit()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *DefaultController) reopenAssignment(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	assignment := ctx.GetAssignment()
	repo := ctx.GetGitlabRepository()

	projects := utils.Filter(assignment.Projects, func(p *database.AssignmentProjects) bool { return p.ProjectStatus == database.Accepted })

	caches := []utils.ProjectAccessLevelCache{}
	defer func() {
		if err != nil {
			for _, cache := range caches {
				repo.ChangeUserAccessLevelInProject(cache.ProjectID, cache.UserID, cache.AccessLevel)
			}
		}
	}()

	for _, project := range projects {
		userClassrooms, err := query.UserClassrooms.
			WithContext(c.Context()).
			Preload(query.UserClassrooms.Classroom).
			Where(query.UserClassrooms.TeamID.Eq(project.TeamID)).
			Find()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		for _, userClassroom := range userClassrooms {
			oldAccessLevel, err := repo.GetAccessLevelOfUserInProject(project.ProjectID, userClassroom.UserID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			if oldAccessLevel == model.OwnerPermissions {
				continue
			}

			err = repo.ChangeUserAccessLevelInProject(project.ProjectID, userClassroom.UserID, model.DeveloperPermissions)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			caches = append(caches, utils.ProjectAccessLevelCache{UserID: userClassroom.UserID, ProjectID: project.ProjectID, AccessLevel: oldAccessLevel})
		}
	}

	return nil
}
