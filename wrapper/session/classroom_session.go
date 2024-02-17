package session

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
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
	gitLabAccessToken   = "gitlab-access-token"
	gitLabRefreshToken  = "gitlab-refresh-token"
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
		s.Set(userState, int(Anonymous))
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
	s.Set(userState, int(state))
}

// GetUserState returns the current UserState of this session.
func (s *ClassroomSession) GetUserState() UserState {
	return UserState(s.Get(userState).(int))
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

func (s *ClassroomSession) SetGitlabAccessToken(token string) {
	s.Set(gitLabAccessToken, token)
}

func (s *ClassroomSession) GetGitlabAccessToken() (string, error) {
	if err := s.checkLogin(); err != nil {
		return "", err
	}
	return s.Get(gitLabAccessToken).(string), nil
}

func (s *ClassroomSession) SetGitlabRefreshToken(token string) {
	s.Set(gitLabRefreshToken, token)
}

func (s *ClassroomSession) GetGitlabRefreshToken() (string, error) {
	if err := s.checkLogin(); err != nil {
		return "", err
	}
	return s.Get(gitLabRefreshToken).(string), nil
}

/// Session

// GetExpiry returns the
func (s *ClassroomSession) GetAccessTokenExpiry() (time.Time, error) {
	if err := s.checkLogin(); err != nil {
		return time.Unix(0, 0), err
	}
	value := s.Get(expiresAt).(int64)
	return time.Unix(value, 0), nil
}

// SetExpiry sets a specific expiration for this session. Throws error when failing.
func (s *ClassroomSession) SetAccessTokenExpiry(exp time.Time) {
	s.Set(expiresAt, exp.Unix())
}
