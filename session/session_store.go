package session

import (
	"backend/model/database"
	"backend/model/database/query"
	"errors"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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

var store *session.Store
var instance *ClassroomSession
var once sync.Once

type ClassroomSession struct {
	session *session.Session
	c       *fiber.Ctx
}

func Get(c *fiber.Ctx) *ClassroomSession {
	once.Do(func() {
		store = session.New()
	})

	s, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	// If session is new and unauthenticated set user state to anonymous
	if s.Get(userState) == nil {
		s.Set(userState, int(Anonymous))
	}

	instance = &ClassroomSession{s, c}

	return instance
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
	return s.session.Get(userID).(int), nil
}

// SetUserID should be set when user is successfully authenticated
func (s *ClassroomSession) SetUserID(user int) {
	s.session.Set(userID, user)
}

// GetUser return the object of the current user. Throws error if user is anonymous.
func (s *ClassroomSession) GetUser() (*database.User, error) {
	userId, err := s.GetUserID()
	if err != nil {
		return nil, err
	}
	u := query.User
	user, err := u.WithContext(s.c.Context()).Where(u.ID.Eq(userId)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SetUserState should be set when user state changes.
func (s *ClassroomSession) SetUserState(state UserState) {
	s.session.Set(userState, int(state))
}

// GetUserState returns the current UserState of this session.
func (s *ClassroomSession) GetUserState() UserState {
	return UserState(s.session.Get(userState).(int))
}

// SetOAuthRedirectTarget should be set with the initial target user wanted, if interrupted by auth middleware.
func (s *ClassroomSession) SetOAuthRedirectTarget(uri string) {
	s.session.Set(oauthRedirectTarget, uri)
}

// GetOAuthRedirectTarget returns the last OAuthRedirectTarget of this session.
func (s *ClassroomSession) GetOAuthRedirectTarget() string {
	return s.session.Get(oauthRedirectTarget).(string)
}

//// GitLab

func (s *ClassroomSession) SetGitlabAccessToken(token string) {
	s.session.Set(gitLabAccessToken, token)
}

func (s *ClassroomSession) GetGitlabAccessToken() (string, error) {
	if err := s.checkLogin(); err != nil {
		return "", err
	}
	return s.session.Get(gitLabAccessToken).(string), nil
}

func (s *ClassroomSession) SetGitlabRefreshToken(token string) {
	s.session.Set(gitLabRefreshToken, token)
}

func (s *ClassroomSession) GetGitlabRefreshToken() (string, error) {
	if err := s.checkLogin(); err != nil {
		return "", err
	}
	return s.session.Get(gitLabRefreshToken).(string), nil
}

//// Session

// Fresh is true if the current session is new
func (s *ClassroomSession) Fresh() bool {
	return s.session.Fresh()
}

// ID returns the session id
func (s *ClassroomSession) ID() string {
	return s.session.ID()
}

// Get will return the value
func (s *ClassroomSession) Get(key string) interface{} {
	return s.session.Get(key)
}

// Set will update or create a new key value
func (s *ClassroomSession) Set(key string, val interface{}) {
	s.session.Set(key, val)
}

// Delete will delete the value
func (s *ClassroomSession) Delete(key string) {
	s.session.Delete(key)
}

// Destroy will delete the session from Storage and expire session cookie
func (s *ClassroomSession) Destroy() error {
	return s.session.Destroy()
}

// Regenerate generates a new session id and delete the old one from Storage
func (s *ClassroomSession) Regenerate() error {
	return s.session.Regenerate()
}

// Reset generates a new session id, deletes the old one from storage, and resets the associated data
func (s *ClassroomSession) Reset() error {
	return s.session.Reset()
}

// Save will update the storage and client cookie
func (s *ClassroomSession) Save() error {
	return s.session.Save()
}

// Keys will retrieve all keys in current session
func (s *ClassroomSession) Keys() []string {
	return s.session.Keys()
}

// GetExpiry returns the
func (s *ClassroomSession) GetExpiry() time.Time {
	return time.Unix(s.session.Get(expiresAt).(int64), 0)
}

// SetExpiry sets a specific expiration for this session. Throws error when failing.
func (s *ClassroomSession) SetExpiry(exp time.Time) error {
	s.session.Set(expiresAt, exp.Unix())
	s.session.SetExpiry(time.Until(exp))
	return s.Save()
}
