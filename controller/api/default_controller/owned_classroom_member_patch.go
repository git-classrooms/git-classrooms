package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type changeOwnedClassroomMemberRequest struct {
	Role   *database.Role `json:"role" validate:"optional"`
	TeamID *uuid.UUID     `json:"teamId" validate:"optional"`
} //@Name ChangeOwnedClassroomMemberRequest

// @Summary		Update Classroom Members team and or role
// @Description	Update Classroom Members team and or role
// @Id 			ChangeOwnedClassroomMember
// @Tags			member
// @Accept			json
// @Param			classroomId		path	string													true	"Classroom ID"	Format(uuid)
// @Param			memberId		path	int														true	"Member ID"
// @Param			changeClassroom	body	default_controller.changeOwnedClassroomMemberRequest	true	"Update ClassroomMemberRequest"
// @Param			X-Csrf-Token	header	string													true	"Csrf-Token"
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/classrooms/owned/{classroomId}/members/{memberId} [patch]
func (ctrl *DefaultController) ChangeOwnedClassroomMember(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	member := ctx.GetOwnedClassroomMember()
	repo := ctx.GetGitlabRepository()

	requestBody := &changeOwnedClassroomMemberRequest{}
	err = c.BodyParser(requestBody)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = repo.GroupAccessLogin(classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var newTeam *database.Team

	queryTeam := query.Team
	if requestBody.TeamID != nil {
		newTeam, err = queryTeam.
			WithContext(c.Context()).
			Where(queryTeam.ID.Eq(*requestBody.TeamID)).
			First()
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}

	if member.TeamID != nil {
		if newTeam != nil && *member.TeamID != newTeam.ID {
			if err = repo.RemoveUserFromGroup(member.Team.GroupID, member.UserID); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			defer func() {
				if recover() != nil || err != nil {
					if err = repo.AddUserToGroup(member.Team.GroupID, member.UserID, model.DeveloperPermissions); err != nil {
						return
					}
				}
			}()
		}
	}

	if requestBody.Role != nil {
		if *requestBody.Role != member.Role {
			member.Role = *requestBody.Role
			if !classroom.StudentsViewAllProjects && *requestBody.Role == database.Moderator {
				if err = repo.ChangeUserAccessLevelInGroup(classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
				defer func() {
					if recover() != nil || err != nil {
						if err = repo.ChangeUserAccessLevelInGroup(classroom.GroupID, member.UserID, model.GuestPermissions); err != nil {
							return
						}
					}
				}()
			} else if !classroom.StudentsViewAllProjects && *requestBody.Role == database.Student {
				if err = repo.ChangeUserAccessLevelInGroup(classroom.GroupID, member.UserID, model.GuestPermissions); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
				defer func() {
					if recover() != nil || err != nil {
						if err = repo.ChangeUserAccessLevelInGroup(classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
							return
						}
					}
				}()
			}
		}
	}

	if newTeam != nil && (member.TeamID == nil || *member.TeamID != newTeam.ID) {
		if err = repo.AddUserToGroup(newTeam.GroupID, member.UserID, model.DeveloperPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		member.TeamID = &newTeam.ID
		defer func() {
			if recover() != nil || err != nil {
				if err = repo.RemoveUserFromGroup(newTeam.GroupID, member.UserID); err != nil {
					return
				}
			}
		}()
	}

	queryUserClassrooms := query.UserClassrooms
	if err := queryUserClassrooms.
		WithContext(c.Context()).
		Save(member); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
