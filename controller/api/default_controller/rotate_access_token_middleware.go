package api

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) RotateAccessTokenMiddleware(c *fiber.Ctx) error {
	ctx := fiberContext.Get(c)
	repo := ctx.GetGitlabRepository()
	classroom := ctx.GetUserClassroom()

	if classroom.Role == database.Owner && !classroom.Classroom.Archived {
		if _, err, _ := ctrl.g.Do(classroom.ClassroomID.String(), func() (interface{}, error) {
			return nil, rotateGroupAccessToken(c.Context(), repo, &classroom.Classroom)
		}); err != nil {
			log.Println(err)
		}
	}

	return c.Next()
}

func rotateGroupAccessToken(ctx context.Context, repo gitlab.Repository, classroom *database.Classroom) error {
	expiresAt := time.Now().AddDate(0, 0, 364)

	// TODO: Remove this after self rotating access tokens
	if len(classroom.Tokens) == 1 {
		accessToken, err := repo.CreateGroupAccessToken(classroom.GroupID, "GitClassrooms", model.OwnerPermissions, expiresAt, "api")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = query.ClassroomToken.WithContext(ctx).Create(&database.ClassroomToken{
			GroupAccessTokenID:        accessToken.ID,
			GroupAccessToken:          accessToken.Token,
			GroupAccessTokenCreatedAt: accessToken.CreatedAt,
			ClassroomID:               classroom.ID,
		}); err != nil {
			return err
		}

		return nil
	}

	token := classroom.Tokens[0]
	if token.GroupAccessTokenCreatedAt.Add(24 * time.Hour).After(time.Now()) {
		return nil
	}

	backupToken := classroom.Tokens[1]

	accessToken, err := repo.RotateGroupAccessToken(classroom.GroupID, backupToken.GroupAccessTokenID, expiresAt)
	if err != nil {
		return err
	}

	log.Println("Rotating access token for classroom", classroom.ID)

	backupToken.GroupAccessTokenID = token.GroupAccessTokenID
	backupToken.GroupAccessToken = token.GroupAccessToken
	backupToken.GroupAccessTokenCreatedAt = token.GroupAccessTokenCreatedAt
	if err = query.ClassroomToken.WithContext(ctx).Save(backupToken); err != nil {
		return err
	}

	token.GroupAccessTokenID = accessToken.ID
	token.GroupAccessToken = accessToken.Token
	token.GroupAccessTokenCreatedAt = accessToken.CreatedAt
	if err = query.ClassroomToken.WithContext(ctx).Save(token); err != nil {
		return err
	}
	return nil
}
