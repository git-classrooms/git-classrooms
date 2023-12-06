package handler

import (
	"backend/api/repository/go_gitlab_repo"
	"backend/auth"
	"backend/session"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"time"
)

func AuthMiddleware(c *fiber.Ctx) error {
	sess := session.Get(c)

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
		token.TokenType = "bearer"

		// Refresh token
		token, err = auth.ConfigGitlab().TokenSource(c.Context(), token).Token()
		if err != nil {
			return err
		}

		// Save refreshed token to session
		sess.SetGitlabAccessToken(token.AccessToken)
		sess.SetGitlabRefreshToken(token.RefreshToken)
		err = sess.SetExpiry(token.Expiry)
		if err != nil {
			return err
		}
	}

	accessToken, err := sess.GetGitlabAccessToken()
	if err != nil {
		return err
	}

	repo := go_gitlab_repo.NewGoGitlabRepo()
	if err := repo.Login(accessToken); err != nil {
		return err
	}

	c.Locals("gitlab-repo", repo)

	return c.Next()
}
