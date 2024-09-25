package session

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

type UserState int

const (
	Anonymous UserState = iota
	LoggedIn
)

// Session keys
const (
	userState           = "user-state"
	userID              = "user-id"
	gitLabOauth2Token   = "gitlab-oauth2-token"
	expiresAt           = "expires-at"
	oauthRedirectTarget = "oauth-redirect-target"
)

// Error keys
const (
	ErrorUnauthenticated = "user is not authenticated (Anonymous)"
)

type ClassroomSession struct {
	*session.Session
}

func Get(c *fiber.Ctx) *ClassroomSession {
	s, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	// If session is new and unauthenticated set user state to anonymous
	if s.Get(userState) == nil {
		s.Set(userState, Anonymous)
	}

	return &ClassroomSession{Session: s}
}

//// User

// checkLogin returns error if ClassroomSession is unauthenticated
func (s *ClassroomSession) checkLogin() error {
	if s.GetUserState() == Anonymous {
		return errors.New(ErrorUnauthenticated)
	}
	return nil
}

// GetUserID returns the id of the current user. Throws error if user is anonymous.
func (s *ClassroomSession) GetUserID() (int, error) {
	if err := s.checkLogin(); err != nil {
		return -1, err
	}
	return s.Get(userID).(int), nil
}

// SetUserID should be set when user is successfully authenticated
func (s *ClassroomSession) SetUserID(user int) {
	s.Set(userID, user)
}

// SetUserState should be set when user state changes.
func (s *ClassroomSession) SetUserState(state UserState) {
	s.Set(userState, state)
}

// GetUserState returns the current UserState of this session.
func (s *ClassroomSession) GetUserState() UserState {
	return s.Get(userState).(UserState)
}

// SetOAuthRedirectTarget should be set with the initial target user wanted, if interrupted by auth middleware.
func (s *ClassroomSession) SetOAuthRedirectTarget(uri string) {
	s.Set(oauthRedirectTarget, uri)
}

// GetOAuthRedirectTarget returns the last OAuthRedirectTarget of this session.
func (s *ClassroomSession) GetOAuthRedirectTarget() string {
	return s.Get(oauthRedirectTarget).(string)
}

//// GitLab

// SetGitlabOauth2Token
func (s *ClassroomSession) SetGitlabOauth2Token(token *oauth2.Token) {
	s.Set(gitLabOauth2Token, token)
}

func (s *ClassroomSession) GetGitlabOauth2Token() (*oauth2.Token, error) {
	if err := s.checkLogin(); err != nil {
		return nil, err
	}
	return s.Get(gitLabOauth2Token).(*oauth2.Token), nil
}
