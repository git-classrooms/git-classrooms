package api

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
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
	if classroom.GroupAccessTokenCreatedAt.Add(24 * time.Hour).After(time.Now()) {
		return nil
	}

	expiresAt := time.Now().AddDate(0, 0, 364)
	accessToken, err := repo.RotateGroupAccessToken(classroom.GroupID, classroom.GroupAccessTokenID, expiresAt)
	if err != nil {
		return err
	}

	log.Println("Rotating access token for classroom", classroom.ID)

	classroom.GroupAccessTokenID = accessToken.ID
	classroom.GroupAccessToken = accessToken.Token
	classroom.GroupAccessTokenCreatedAt = accessToken.CreatedAt
	if err = query.Classroom.WithContext(ctx).Save(classroom); err != nil {
		return err
	}
	return nil
}
