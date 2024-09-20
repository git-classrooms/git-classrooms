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
		go func(ctx context.Context, repo gitlab.Repository, classroom database.Classroom) {
			if _, err, _ := ctrl.g.Do(classroom.ID.String(), func() (interface{}, error) {
				return nil, rotateGroupAccessToken(ctx, repo, &classroom)
			}); err != nil {
				log.Println(err)
			}
		}(c.Context(), repo, classroom.Classroom)
	}

	return c.Next()
}

func rotateGroupAccessToken(ctx context.Context, repo gitlab.Repository, classroom *database.Classroom) error {
	accessToken, err := repo.GetGroupAccessToken(classroom.GroupID, classroom.GroupAccessTokenID)
	if err != nil {
		return err
	}

	if accessToken.CreatedAt.Add(24 * time.Hour).After(time.Now()) {
		return nil
	}

	expiresAt := time.Now().AddDate(0, 0, 364)
	accessToken, err = repo.RotateGroupAccessToken(classroom.GroupID, classroom.GroupAccessTokenID, expiresAt)
	if err != nil {
		return err
	}

	queryClassroom := query.Classroom
	if _, err = queryClassroom.
		WithContext(ctx).
		Where(queryClassroom.ID.Eq(classroom.ID)).
		Updates(database.Classroom{GroupAccessTokenID: accessToken.ID, GroupAccessToken: accessToken.Token}); err != nil {
		return err
	}
	return nil
}
