package auth

import (
	"github.com/gofiber/fiber/v2"
	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	"golang.org/x/oauth2"
	"gorm.io/gen/field"
	"log"
	"time"
)

type OAuthController struct {
	authConfig   authConfig.Config
	gitlabConfig gitlabConfig.Config
}

func NewOAuthController(authConfig authConfig.Config,
	gitlabConfig gitlabConfig.Config) *OAuthController {

	return &OAuthController{
		authConfig:   authConfig,
		gitlabConfig: gitlabConfig,
	}
}

// Auth fiber handler
func (ctrl *OAuthController) Auth(c *fiber.Ctx) error {
	redirect := c.Query("redirect", "/")

	s := session.Get(c)
	s.SetOAuthRedirectTarget(redirect)
	if err := s.Save(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	oauth := ctrl.authConfig.GetOAuthConfig()
	url := oauth.AuthCodeURL("state")
	return c.Redirect(url)
}

// Callback to receive gitlabs' response
func (ctrl *OAuthController) Callback(c *fiber.Ctx) error {
	token, err := ctrl.authConfig.GetOAuthConfig().Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	repo := gitlabRepo.NewGitlabRepo(ctrl.gitlabConfig)
	if err := repo.Login(token.AccessToken); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Get user from GitLab
	gitlabUser, err := repo.GetCurrentUser()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Save or Update user in DB
	u := query.User
	user, err := u.WithContext(c.Context()).
		Where(u.ID.Eq(gitlabUser.ID)).
		Assign(field.Attrs(&database.User{GitlabEmail: gitlabUser.Email, Name: gitlabUser.Name})).
		FirstOrCreate()

	if err != nil {
		// TODO: Use sentry to log errors
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	s := session.Get(c)

	// Save GitLab session in local user session
	s.SetGitlabAccessToken(token.AccessToken)
	s.SetGitlabRefreshToken(token.RefreshToken)

	s.SetUserState(session.LoggedIn)
	s.SetUserID(user.ID)

	s.SetExpiry(token.Expiry)

	redirect := s.GetOAuthRedirectTarget()

	if err = s.Save(); err != nil {
		// TODO: Use sentry to log errors
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	return c.Redirect(redirect)
}

// AuthMiddleware to check session for Gitlab Tokens
func (ctrl *OAuthController) AuthMiddleware(c *fiber.Ctx) error {
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
		token, err = ctrl.authConfig.GetOAuthConfig().TokenSource(c.Context(), token).Token()
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
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	repo := gitlabRepo.NewGitlabRepo(ctrl.gitlabConfig)
	if err := repo.Login(accessToken); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	context.SetGitlabRepository(c, repo)

	return c.Next()
}
