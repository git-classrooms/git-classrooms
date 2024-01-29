package context

import (
	"backend/api/repository/go_gitlab_repo"
	"backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestClassroomSession_Delete(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()

	t.Run("Get GitlabRepo", func(t *testing.T) {
		t.Run("Should return nil if not set", func(t *testing.T) {
			// given
			req := new(fasthttp.RequestCtx)
			ctx := app.AcquireCtx(req)
			defer app.ReleaseCtx(ctx)

			// when
			repo := GetGitlabRepository(ctx)

			//then
			assert.Nil(t, repo)
		})

		t.Run("Should return nil if false value is set", func(t *testing.T) {
			// given
			req := new(fasthttp.RequestCtx)
			ctx := app.AcquireCtx(req)
			defer app.ReleaseCtx(ctx)
			ctx.Locals(gitlabRepoKey, "test")

			// when
			repo := GetGitlabRepository(ctx)

			// then
			assert.Nil(t, repo)
		})

		t.Run("Should return repo if value is set", func(t *testing.T) {
			// given
			req := new(fasthttp.RequestCtx)
			ctx := app.AcquireCtx(req)
			defer app.ReleaseCtx(ctx)
			shouldRepo := go_gitlab_repo.NewGoGitlabRepo(&config.Config{})
			ctx.Locals(gitlabRepoKey, shouldRepo)

			// when
			repo := GetGitlabRepository(ctx)

			// then
			assert.Equal(t, shouldRepo, repo)
		})
	})

	t.Run("Set GitlabRepo", func(t *testing.T) {
		t.Run("Should be empty if not set", func(t *testing.T) {
			// given
			req := new(fasthttp.RequestCtx)
			ctx := app.AcquireCtx(req)
			defer app.ReleaseCtx(ctx)

			// when
			SetGitlabRepository(ctx, nil)

			// then
			repo := ctx.Locals(gitlabRepoKey)
			assert.Nil(t, repo)
		})

		t.Run("Should set repo", func(t *testing.T) {
			// given
			req := new(fasthttp.RequestCtx)
			ctx := app.AcquireCtx(req)
			defer app.ReleaseCtx(ctx)
			shouldRepo := go_gitlab_repo.NewGoGitlabRepo(&config.Config{})

			// when
			SetGitlabRepository(ctx, shouldRepo)

			// then
			repo := ctx.Locals(gitlabRepoKey)
			assert.Equal(t, shouldRepo, repo)
		})
	})
}
