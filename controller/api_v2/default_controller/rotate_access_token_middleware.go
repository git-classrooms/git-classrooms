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
		go rotateGroupAccessToken(c.Context(), repo, &classroom.Classroom)
	}

	return nil
}

func rotateGroupAccessToken(ctx context.Context, repo gitlab.Repository, classroom *database.Classroom) {
	accesstoken, err := repo.GetGroupAccessToken(classroom.GroupID, classroom.GroupAccessTokenID)
	if err != nil {
		log.Println(err)
		return
	}

	if accesstoken.CreatedAt.Add(24 * time.Hour).After(time.Now()) {
		return
	}

	expiresAt := time.Now().AddDate(0, 0, 364)
	accessToken, err := repo.RotateGroupAccessToken(classroom.GroupID, classroom.GroupAccessTokenID, expiresAt)
	if err != nil {
		log.Println(err)
		return
	}

	classroom.GroupAccessTokenID = accessToken.ID
	classroom.GroupAccessToken = accessToken.Token
	if err = query.Classroom.WithContext(ctx).Save(classroom); err != nil {
		log.Println(err)
		return
	}
}
