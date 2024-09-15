package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	gitlabModel "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type updateMemberRoleRequest struct {
	Role *database.Role `json:"role"`
} //@Name UpdateMemberRoleRequest

func (r updateMemberRoleRequest) isValid() bool {
	return r.Role != nil
}

// @Summary		Update Classroom Members role
// @Description	Update Classroom Members role
// @Id				UpdateMemberRole
// @Tags			member
// @Accept			json
// @Param			classroomId			path	string						true	"Classroom ID"	Format(uuid)
// @Param			memberId			path	int							true	"Member ID"
// @Param			updateMemberRole	body	api.updateMemberRoleRequest	true	"Update Member Role"
// @Param			X-Csrf-Token		header	string						true	"Csrf-Token"
// @Success		202
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/members/{memberId}/role [patch]
func (ctrl *DefaultController) UpdateMemberRole(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	member := ctx.GetClassroomMember()
	repo := ctx.GetGitlabRepository()

	var requestBody updateMemberRoleRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	if err = repo.GroupAccessLogin(classroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if *requestBody.Role == member.Role {
		return c.SendStatus(fiber.StatusNoContent)
	}

	if classroom.Classroom.OwnerID == member.UserID {
		return fiber.NewError(fiber.StatusForbidden, "The Role of the Creator of the classroom cannot be changed.")
	}

	if *requestBody.Role == database.Owner && classroom.Classroom.OwnerID != classroom.UserID {
		return fiber.NewError(fiber.StatusForbidden, "Only the Creator of the classroom can assign the owner role.")
	}

	if member.Role == database.Owner && classroom.Classroom.OwnerID != classroom.UserID {
		return fiber.NewError(fiber.StatusForbidden, "Only the Creator of the classroom can remove the owner role.")
	}

	oldRole := member.Role
	viewOtherProjects := classroom.Classroom.StudentsViewAllProjects
	member.Role = *requestBody.Role

	switch {
	case oldRole == database.Owner && *requestBody.Role == database.Student && viewOtherProjects:
		fallthrough
	case oldRole == database.Owner && *requestBody.Role == database.Moderator:
		if err = repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.OwnerPermissions); err != nil {
					log.Println(err)
				}
			}
		}()

	case oldRole == database.Owner && *requestBody.Role == database.Student && !viewOtherProjects:
		if err = repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.GuestPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.OwnerPermissions); err != nil {
					log.Println(err)
				}
			}
		}()

	case oldRole == database.Moderator && *requestBody.Role == database.Student && !viewOtherProjects:
		if err = repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.GuestPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
					log.Println(err)
				}
			}
		}()

	case oldRole == database.Moderator && *requestBody.Role == database.Student && viewOtherProjects:
		// The Permission donẗ change
		break

	case oldRole == database.Moderator && *requestBody.Role == database.Owner:
		if err = repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.OwnerPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
					log.Println(err)
				}
			}
		}()

	case oldRole == database.Student && *requestBody.Role == database.Moderator && !viewOtherProjects:
		if err = repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.GuestPermissions); err != nil {
					log.Println(err)
				}
			}
		}()

	case oldRole == database.Student && *requestBody.Role == database.Moderator && viewOtherProjects:
		// The Permission donẗ change

	case oldRole == database.Student && *requestBody.Role == database.Owner && !viewOtherProjects:
		if err = repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.OwnerPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.GuestPermissions); err != nil {
					log.Println(err)
				}
			}
		}()

	case oldRole == database.Student && *requestBody.Role == database.Owner && viewOtherProjects:
		if err = repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.OwnerPermissions); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			if recover() != nil || err != nil {
				if err := repo.ChangeUserAccessLevelInGroup(classroom.Classroom.GroupID, member.UserID, model.ReporterPermissions); err != nil {
					log.Println(err)
				}
			}
		}()
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		queryTeam := tx.Team
		if member.Role != database.Student {
			if classroom.Classroom.MaxTeamSize == 1 {
				if _, err := queryTeam.WithContext(c.Context()).Delete(member.Team); err != nil {
					return err
				}

				if err := repo.DeleteGroup(member.Team.GroupID); err != nil {
					return err
				}
			}
			member.Team = nil
			member.TeamID = nil
		} else if classroom.Classroom.MaxTeamSize == 1 {
			subgroup, err := repo.CreateSubGroup(
				member.User.Name,
				member.User.GitlabUsername,
				classroom.Classroom.GroupID,
				gitlabModel.Private,
				fmt.Sprintf("Team %s of classroom %s", member.User.Name, classroom.Classroom.Name),
			)
			if err != nil {
				return err
			}
			defer func() {
				if recover() != nil || err != nil {
					repo.DeleteGroup(subgroup.ID)
				}
			}()

			team := &database.Team{
				ClassroomID: classroom.ClassroomID,
				GroupID:     subgroup.ID,
				Name:        member.User.Name,
			}
			err = queryTeam.WithContext(c.Context()).Create(team)
			if err != nil {
				return err
			}
			member.TeamID = &team.ID
			member.Team = team

			if err = repo.AddUserToGroup(team.GroupID, member.UserID, gitlabModel.ReporterPermissions); err != nil {
				return err
			}
		}

		return tx.UserClassrooms.
			WithContext(c.Context()).
			Save(member)
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusAccepted)
}
