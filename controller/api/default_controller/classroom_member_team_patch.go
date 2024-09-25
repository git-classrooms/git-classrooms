package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateMemberTeamRequest struct {
	TeamID *uuid.UUID `json:"teamId"`
} // @Name UpdateMemberTeamRequest

func (r updateMemberTeamRequest) isValid() bool {
	return r.TeamID != nil
}

// @Summary		Update Classroom Members team
// @Description	Update Classroom Members team
// @Id				UpdateMemberTeam
// @Tags			member
// @Accept			json
// @Param			classroomId			path	string						true	"Classroom ID"	Format(uuid)
// @Param			memberId			path	int							true	"Member ID"
// @Param			updateMemberTeam	body	api.updateMemberTeamRequest	true	"Update Member Team"
// @Param			X-Csrf-Token		header	string						true	"Csrf-Token"
// @Success		202
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/members/{memberId}/team [patch]
func (ctrl *DefaultController) UpdateMemberTeam(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	member := ctx.GetClassroomMember()
	repo := ctx.GetGitlabRepository()

	if classroom.Classroom.MaxTeamSize == 1 {
		return fiber.NewError(fiber.StatusForbidden, "Teams are disabled for this classroom.")
	}

	if member.Role != database.Student {
		return fiber.NewError(fiber.StatusForbidden, "Only students can be assigned to a team.")
	}

	var requestBody updateMemberTeamRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	if err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryTeam := query.Team
	newTeam, err := queryTeam.
		WithContext(c.Context()).
		Where(queryTeam.ID.Eq(*requestBody.TeamID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if member.TeamID != nil {
		if *member.TeamID == newTeam.ID {
			return c.SendStatus(fiber.StatusNoContent)
		}
		if err = repo.RemoveUserFromGroup(member.Team.GroupID, member.UserID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.AddUserToGroup(member.Team.GroupID, member.UserID, model.ReporterPermissions); err != nil {
					log.Println(err)
				}
			}
		}()

		queryAssignmentProjects := query.AssignmentProjects
		projects, err := queryAssignmentProjects.
			WithContext(c.Context()).
			Preload(queryAssignmentProjects.Assignment).
			Where(queryAssignmentProjects.TeamID.Eq(*member.TeamID)).
			Find()

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		for _, project := range projects {
			if err = repo.RemoveUserFromProject(project.ProjectID, member.UserID); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}
		defer func() {
			if recover() != nil || err != nil {
				for _, project := range projects {
					accessLevel := model.DeveloperPermissions
					if project.Assignment.Closed {
						accessLevel = model.ReporterPermissions
					}
					if err := repo.AddProjectMember(project.ProjectID, member.UserID, accessLevel); err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}

	if err = repo.AddUserToGroup(newTeam.GroupID, member.UserID, model.ReporterPermissions); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	member.TeamID = &newTeam.ID
	defer func() {
		if recover() != nil || err != nil {
			if err := repo.RemoveUserFromGroup(newTeam.GroupID, member.UserID); err != nil {
				log.Println(err)
			}
		}
	}()

	queryAssignmentProjects := query.AssignmentProjects
	projects, err := queryAssignmentProjects.
		WithContext(c.Context()).
		Preload(queryAssignmentProjects.Assignment).
		Where(queryAssignmentProjects.TeamID.Eq(*member.TeamID)).
		Find()

	for _, project := range projects {
		accessLevel := model.DeveloperPermissions
		if project.Assignment.Closed {
			accessLevel = model.ReporterPermissions
		}
		if err := repo.AddProjectMember(project.ProjectID, member.UserID, accessLevel); err != nil {
			log.Println(err)
		}
	}
	defer func() {
		if recover() != nil || err != nil {
			for _, project := range projects {
				if err = repo.RemoveUserFromProject(project.ProjectID, member.UserID); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	queryUserClassrooms := query.UserClassrooms
	if err = queryUserClassrooms.
		WithContext(c.Context()).
		Save(member); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
