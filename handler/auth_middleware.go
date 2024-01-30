package handler

import (
	"de.hs-flensburg.gitlab/gitlab-classroom/api/repository/go_gitlab_repo"
	"de.hs-flensburg.gitlab/gitlab-classroom/auth"
	"de.hs-flensburg.gitlab/gitlab-classroom/config"
	"de.hs-flensburg.gitlab/gitlab-classroom/context"
	"de.hs-flensburg.gitlab/gitlab-classroom/session"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	applicationConfig, err := config.GetConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, session.ErrorUnauthenticated)
	}

	sess := session.Get(c)

	if sess.GetUserState() == session.Anonymous {
		return fiber.NewError(fiber.StatusUnauthorized, session.ErrorUnauthenticated)
	}

	exp := sess.GetExpiry()

	// exp.Add(-20 * time.Minute).After(time.Now())
	// If
	if exp.Before(time.Now().Add(20 * time.Minute)) {
		refreshToken, err := sess.GetGitlabRefreshToken()
		if err != nil {
			return err
		}

		// Build refresh token from session
		token := new(oauth2.Token)
		token.RefreshToken = refreshToken
		token.Expiry = time.Now().Add(-1 * time.Minute)
		token.TokenType = "Bearer"

		// Refresh token
		token, err = auth.ConfigGitlab(applicationConfig).TokenSource(c.Context(), token).Token()
		if err != nil {
			return err
		}

		// Save refreshed token to session
		sess.SetGitlabAccessToken(token.AccessToken)
		sess.SetGitlabRefreshToken(token.RefreshToken)
		sess.SetExpiry(token.Expiry)
		if err = sess.Save(); err != nil {
			return err
		}
		// sess.Save does save the session, which invalidates the pointer and we must get a new one
		sess = session.Get(c)
	}

	accessToken, err := sess.GetGitlabAccessToken()
	if err != nil {
		return err
	}

	repo := go_gitlab_repo.NewGoGitlabRepo(applicationConfig)
	if err := repo.Login(accessToken); err != nil {
		return err
	}

	context.SetGitlabRepository(c, repo)

	return c.Next()
}
