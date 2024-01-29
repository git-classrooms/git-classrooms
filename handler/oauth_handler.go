package handler

import (
	"backend/api/repository/go_gitlab_repo"
	"backend/auth"
	"backend/config"
	"backend/model/database/query"
	"backend/session"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Auth fiber handler
func Auth(c *fiber.Ctx) error {
	applicationConfig, err := config.GetConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, session.ErrorUnauthenticated)
	}

	path := auth.ConfigGitlab(applicationConfig)
	redirect := c.Query("redirect", "/")

	s := session.Get(c)
	s.SetOAuthRedirectTarget(redirect)
	if err := s.Save(); err != nil {
		return err
	}

	url := path.AuthCodeURL("state")
	return c.Redirect(url)
}

// Callback to receive gitlabs' response
func Callback(c *fiber.Ctx) error {
	applicationConfig, err := config.GetConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, session.ErrorUnauthenticated)
	}

	token, err := auth.ConfigGitlab(applicationConfig).Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		return err
	}

	repo := go_gitlab_repo.NewGoGitlabRepo(applicationConfig)
	if err := repo.Login(token.AccessToken); err != nil {
		return err
	}

	// Get user from GitLab
	gitlabUser, err := repo.GetCurrentUser()
	if err != nil {
		return err
	}

	// Save or Update user in DB
	u := query.User
	user, err := u.WithContext(c.Context()).
		Where(u.ID.Eq(gitlabUser.ID)).
		FirstOrCreate()

	if err != nil {
		return nil
	}

	sess := session.Get(c)

	// Save GitLab session in local user session
	sess.SetGitlabAccessToken(token.AccessToken)
	sess.SetGitlabRefreshToken(token.RefreshToken)

	sess.SetUserState(session.LoggedIn)
	sess.SetUserID(user.ID)

	sess.SetExpiry(token.Expiry)
	if err = sess.Save(); err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	s := session.Get(c)
	redirect := s.GetOAuthRedirectTarget()

	return c.Redirect(redirect)
}
