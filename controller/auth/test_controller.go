package auth

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
)

type TestAuthController struct {
	user       *database.User
	gitlabRepo gitlabRepo.Repository
}

func NewTestAuthController(user *database.User, gitlabRepo gitlabRepo.Repository) *TestAuthController {
	return &TestAuthController{user: user, gitlabRepo: gitlabRepo}
}

func (ctrl *TestAuthController) SignIn(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (ctrl *TestAuthController) Callback(c *fiber.Ctx) error {
	s := session.Get(c)

	s.SetUserState(session.LoggedIn)
	s.SetUserID(ctrl.user.ID)
	s.Save()

	return c.SendStatus(fiber.StatusOK)
}

func (ctrl *TestAuthController) SignOut(c *fiber.Ctx) error {
	err := session.Get(c).Destroy()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (ctrl *TestAuthController) GetAuth(c *fiber.Ctx) error {
	s := session.Get(c)
	_, err := s.GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (ctrl *TestAuthController) AuthMiddleware(c *fiber.Ctx) error {
	ctx := fiberContext.Get(c)
	if ctx == nil {
	ctx.SetGitlabRepository(ctrl.gitlabRepo)
	ctx.SetUserID(ctrl.user.ID)
	return c.Next()
}

func (ctrl *TestAuthController) GetCsrf(c *fiber.Ctx) error {
	type response struct {
		Csrf string `json:"csrf"`
	}
	token, ok := c.Locals("csrf").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "There is no csrf token in the context")
	}
	return c.JSON(response{Csrf: token})
}
