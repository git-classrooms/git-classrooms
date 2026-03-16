package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/logging"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) RotateAccessTokenMiddleware(c *fiber.Ctx) error {
	ctx := fiberContext.Get(c)
	log := ctx.GetLoggerForHandler("RotateAccessTokenMiddleware")
	repo := ctx.GetGitlabRepository()
	classroom := ctx.GetUserClassroom()

	logCtx := logging.SetLogger(c.Context(), log)

	if classroom.Role == database.Owner && !classroom.Classroom.Archived {
		if _, err, _ := ctrl.g.Do(classroom.ClassroomID.String(), func() (interface{}, error) {
			return nil, rotateGroupAccessToken(logCtx, repo, &classroom.Classroom)
		}); err != nil {
			log.Error("error while rotating group access token", "error", err)
		}
	}

	return c.Next()
}

func rotateGroupAccessToken(ctx context.Context, repo gitlab.Repository, classroom *database.Classroom) error {
	log := logging.GetLogger(ctx)
	if classroom.GroupAccessTokenCreatedAt.Add(24 * time.Hour).After(time.Now()) {
		return nil
	}

	expiresAt := time.Now().AddDate(0, 0, 364)
	accessToken, err := repo.RotateGroupAccessToken(classroom.GroupID, classroom.GroupAccessTokenID, expiresAt)
	if err != nil {
		return err
	}

	log.Info("Rotating access token for classroom", "classroomID", classroom.ID)

	classroom.GroupAccessTokenID = accessToken.ID
	classroom.GroupAccessToken = accessToken.Token
	classroom.GroupAccessTokenCreatedAt = accessToken.CreatedAt
	if err = query.Classroom.WithContext(ctx).Save(classroom); err != nil {
		return err
	}
	return nil
}
