package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/groupcache/singleflight"
	authConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"golang.org/x/oauth2"
	"gorm.io/gen/field"
)

type OAuthController struct {
	authConfig   *oauth2.Config
	gitlabConfig gitlabConfig.Config
	g            *singleflight.Group
}

func NewOAuthController(authConfig authConfig.Config,
	gitlabConfig gitlabConfig.Config) *OAuthController {
	g := &singleflight.Group{}
	return &OAuthController{
		authConfig:   authConfig.GetOAuthConfig(),
		gitlabConfig: gitlabConfig,
		g:            g,
	}
}

type authState struct {
	Csrf     string `json:"csrf"`
	Redirect string `json:"redirect"`
}

type authRequest struct {
	Redirect string `form:"redirect"`
}

func (ctrl *OAuthController) SignIn(c *fiber.Ctx) error {
	body := &authRequest{}
	c.BodyParser(body)

	redirect := "/"
	if body.Redirect != "" {
		redirect = body.Redirect
	}

	csrf := c.Locals("csrf").(string)

	stateBytes, err := json.Marshal(authState{csrf, redirect})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	authCodeOption := oauth2.S256ChallengeOption("Challenge")              // here include PKCE-Challenge against csrf-attacks should be random
	url := ctrl.authConfig.AuthCodeURL(string(stateBytes), authCodeOption) // the string state is sent back by the auth-server of gitlab here we could include the redirect url | and or we include a random csrf token that will be validated

	return c.Redirect(url, fiber.StatusSeeOther)
}

// Callback to receive gitlabs' response
func (ctrl *OAuthController) Callback(c *fiber.Ctx) error {
	authCodeOption := oauth2.VerifierOption("Challenge") // this is the validation of the PKCE-Challenge
	token, err := ctrl.authConfig.Exchange(c.Context(), c.FormValue("code"), authCodeOption)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	stateVal := c.FormValue("state") // get the state passed from auth, which was sent by gitlab
	state := &authState{}

	if err := json.Unmarshal([]byte(stateVal), state); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	if state.Csrf == "" || state.Redirect == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid state")
	}

	// Check if the csrf token is valid
	if state.Csrf != c.Locals("csrf").(string) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid csrf token")
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
		Assign(field.Attrs(&database.User{GitlabEmail: gitlabUser.Email, Name: gitlabUser.Name, GitlabUsername: gitlabUser.Username})).
		FirstOrCreate()
	if err != nil {
		// TODO: Use sentry to log errors
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	ua := query.UserAvatar
	_, err = ua.WithContext(c.Context()).
		Where(ua.UserID.Eq(user.ID)).
		Assign(field.Attrs(&database.UserAvatar{
			UserID:            user.ID,
			AvatarURL:         gitlabUser.Avatar.AvatarURL,
			FallbackAvatarURL: gitlabUser.Avatar.FallbackAvatarURL,
		})).FirstOrCreate()

	s := session.Get(c)

	// Save GitLab session in local user session
	s.SetGitlabOauth2Token(token)

	s.SetUserState(session.LoggedIn)
	s.SetUserID(user.ID)

	redirect := state.Redirect

	if err = s.Save(); err != nil {
		// TODO: Use sentry to log errors
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	return c.Redirect(redirect, fiber.StatusSeeOther)
}

func (ctrl *OAuthController) SignOut(c *fiber.Ctx) error {
	body := &authRequest{}
	c.BodyParser(body)

	redirect := "/"
	if body.Redirect != "" {
		redirect = body.Redirect
	}

	err := session.Get(c).Destroy()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Redirect(redirect, fiber.StatusSeeOther)
}

func (ctrl *OAuthController) GetAuth(c *fiber.Ctx) error {
	s := session.Get(c)
	_, err := s.GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// AuthMiddleware to check session for Gitlab Tokens
func (ctrl *OAuthController) AuthMiddleware(c *fiber.Ctx) error {
	sess := session.Get(c)

	userId, err := sess.GetUserID()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	token, err := sess.GetGitlabOauth2Token()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// exp.Add(-20 * time.Minute).After(time.Now())
	// If
	if token.Expiry.Before(time.Now().Add(20 * time.Minute)) {
		// this added to prevent multiple requests from refreshing the token at the same time
		// If 2 refresh requests are sent at the same time, the first one will refresh the token
		// and the second would get an error because the refresh token was already used
		_, err := ctrl.g.Do(fmt.Sprintf("%d", userId), func() (interface{}, error) {
			return nil, ctrl.refreshSession(c.Context(), sess)
		})
		if err != nil {
			return err
		}
		// sess.Save does save the session, which invalidates the pointer and we must get a new one
		sess = session.Get(c)
		token, err = sess.GetGitlabOauth2Token()
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
	}

	repo := gitlabRepo.NewGitlabRepo(ctrl.gitlabConfig)
	if err := repo.Login(token.AccessToken); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Set every variable from the session to the context
	ctx := fiberContext.Get(c)
	ctx.SetGitlabRepository(repo)
	ctx.SetUserID(userId)
	return ctx.Next()
}

// GetCsrf returns a csrf token
//
//	@Summary		Show your csrf-Token
//	@Description	Get your csrf-Token
//	@ID				get-csrf
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	auth.GetCsrf.response
//	@Failure		401	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/api/v1/auth/csrf [get]
func (ctrl *OAuthController) GetCsrf(c *fiber.Ctx) error {
	type response struct {
		Csrf string `json:"csrf"`
	}
	token, ok := c.Locals("csrf").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "There is no csrf token in the context")
	}
	return c.JSON(response{Csrf: token})
}

func (ctrl *OAuthController) refreshSession(c context.Context, sess *session.ClassroomSession) error {
	token, err := sess.GetGitlabOauth2Token()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Set expiry to past to force refresh
	token.Expiry = time.Now().Add(-1 * time.Minute)

	// Refresh token
	newToken, err := ctrl.authConfig.TokenSource(c, token).Token()
	if err != nil {
		oldError := err
		err = sess.Destroy()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, oldError.Error()+"\n"+err.Error())
		}
		return fiber.NewError(fiber.StatusUnauthorized, oldError.Error())
	}

	// Save refreshed token to session
	sess.SetGitlabOauth2Token(newToken)
	if err = sess.Save(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
