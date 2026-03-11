package api

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type createTeachingGroupRequest struct {
} //@Name CreateTeachingGroupRequest

func (r createTeachingGroupRequest) isValid() bool {
	return true
}

// @Summary		CreateClassroomTeachingGroup
// @Description	CreateClassroomTeachingGroup
// @Id				CreateClassroomTeachingGroup
// @Tags			classroom
// @Accept			json
// @Param			classroomId			path	string							true	"Classroom ID"	Format(uuid)
// @Param			teachingGroupInfo	body	api.createTeachingGroupRequest	true	"Teaching-Group Info"
// @Param			X-Csrf-Token		header	string							true	"Csrf-Token"
// @Success		201
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		409	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/teachingGroup [post]
func (ctrl *DefaultController) CreateClassroomTeachingGroup(c *fiber.Ctx) (err error) {
	ctx := fiberContext.Get(c)
	repo := ctx.GetGitlabRepository()
	userClassroom := ctx.GetUserClassroom()

	var requestBody createTeachingGroupRequest
	if err = c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if !requestBody.isValid() {
		return fiber.ErrBadRequest
	}

	if userClassroom.Classroom.TeachingGroupID != nil {
		return fiber.NewError(fiber.StatusConflict, "teaching group already exists")
	}

	// reauthenticate the repo with the group access token
	if err = repo.GroupAccessLogin(userClassroom.Classroom.GroupAccessToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	queryClassroom := query.Classroom
	classroom, err := queryClassroom.
		WithContext(c.Context()).
		Preload(queryClassroom.Member).
		Where(queryClassroom.ID.Eq(userClassroom.ClassroomID)).
		First()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	groupID, cleanup, err := ctrl.createTeachingGroup(c.Context(), repo, classroom)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		cleanup(err)
	}()

	classroom.TeachingGroupID = &groupID

	if err = queryClassroom.WithContext(c.Context()).Save(classroom); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (ctrl *DefaultController) createTeachingGroup(ctx context.Context, repo gitlab.Repository, classroom *database.Classroom) (groupID int, cleanupOnError func(error), err error) {
	group, err := repo.CreateSubGroup(
		"Teaching Material",
		"teaching-material",
		classroom.GroupID,
		model.Private,
		fmt.Sprintf("Teaching Material of classroom %s\n\n\n__Managed by [GitClassrooms](%s/classrooms/%s)__", classroom.Name, ctrl.config.PublicURL, classroom.ID.String()),
	)
	if err != nil {
		return -1, nil, err
	}

	cleanup := func(err error) {
		if err != nil {
			if err := repo.DeleteGroup(group.ID); err != nil {
				log.Println(err.Error())
			}
		}
	}
	defer func() {
		cleanup(err)
	}()

	for _, member := range classroom.Member {
		if member.Role == database.Owner {
			continue
		}
		if err := repo.AddUserToGroup(group.ID, member.UserID, model.ReporterPermissions); err != nil {
			return -1, nil, err
		}
	}

	return group.ID, cleanup, nil
}
