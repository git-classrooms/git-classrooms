package session

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

func TestClassroomSession_GetGitlabOauth2Token(t *testing.T) {
	// Mock a fiber context
	InitSessionStore("")
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("User Not Logged In", func(t *testing.T) {
		ses := Get(ctx)
		// Ensure the user is not logged in
		ses.Set(userState, Anonymous)

		_, err := ses.GetGitlabOauth2Token()
		assert.NotNil(t, err, "Expected an error for unauthenticated user")
	})

	t.Run("User Logged In With GitLab Access Token", func(t *testing.T) {
		ses := Get(ctx)
		// Set user as logged in and set a GitLab access token
		ses.Set(userState, LoggedIn)
		expectedToken := &oauth2.Token{AccessToken: "gitlab-access-token-value", RefreshToken: "gitlab-refresh-token-value", Expiry: time.Now()}
		ses.Set(gitLabOauth2Token, expectedToken)

		token, err := ses.GetGitlabOauth2Token()
		assert.Nil(t, err, "Expected no error for authenticated user")
		assert.Equal(t, expectedToken, token, "The GitLab access token should match the set value")
	})
}

func TestClassroomSession_GetUserID(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("User Not Logged In", func(t *testing.T) {
		ses := Get(ctx)
		// Set user state to anonymous or ensure it's not logged in
		ses.Set(userState, Anonymous)

		userID, err := ses.GetUserID()
		assert.Equal(t, -1, userID)
		assert.NotNil(t, err)
	})

	t.Run("User Logged In With Valid User ID", func(t *testing.T) {
		ses := Get(ctx)
		// Set user state to logged in and set a valid user ID
		ses.Set(userState, LoggedIn)
		expectedUserID := 123 // Example user ID
		ses.Set(userID, expectedUserID)

		userID, err := ses.GetUserID()
		assert.Equal(t, expectedUserID, userID)
		assert.Nil(t, err)
	})
}

func TestClassroomSession_GetUserState(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	//Tests
	t.Run("Retrieving User State Anonymous", func(t *testing.T) {
		ses := Get(ctx)
		ses.Set(userState, Anonymous)

		state := ses.GetUserState()

		assert.Equal(t, Anonymous, state, "The retrieved user state should be Anonymous")
	})

	t.Run("Retrieving User State LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.Set(userState, LoggedIn)

		state := ses.GetUserState()

		assert.Equal(t, LoggedIn, state, "The retrieved user state should be LoggedIn")
	})
}

func TestClassroomSession_SetGitlabOauth2Token(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Test
	t.Run("Verify Access Token Set", func(t *testing.T) {
		ses := Get(ctx)

		// Define the GitLab access token
		testToken := &oauth2.Token{AccessToken: "test-access-token", RefreshToken: "test-refresh-token", Expiry: time.Now()}

		// Set the access token using the SetGitlabAccessToken method
		ses.SetGitlabOauth2Token(testToken)

		// Retrieve the saved access token from the session
		savedToken := ses.Get(gitLabOauth2Token)

		// Verify that the saved access token matches the provided token
		assert.Equal(t, testToken, savedToken, "GitLab access token should be set correctly")
	})
}

func TestClassroomSession_SetUserID(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Setting User ID", func(t *testing.T) {
		ses := Get(ctx)
		expectedUserID := 123 // Example user ID

		// Set the user ID
		ses.SetUserID(expectedUserID)

		// Retrieve the set user ID from the session
		userID := ses.Get(userID)

		assert.Equal(t, expectedUserID, userID, "The user ID should be correctly set in the session")
	})
}

func TestClassroomSession_SetUserState(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Setting User State to Anonymous", func(t *testing.T) {
		ses := Get(ctx)
		ses.SetUserState(Anonymous)

		// Retrieve the set user state from the session
		state := ses.Get(userState).(UserState)

		assert.Equal(t, Anonymous, state, "The user state should be set to Anonymous")
	})

	t.Run("Setting User State to LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.SetUserState(LoggedIn)

		// Retrieve the set user state from the session
		state := ses.Get(userState).(UserState)

		assert.Equal(t, LoggedIn, state, "The user state should be set to LoggedIn")
	})
}

func TestClassroomSession_checkLogin(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	//Tests
	t.Run("User State is Anonymous", func(t *testing.T) {
		ses := Get(ctx)
		ses.Set(userState, Anonymous)

		err := ses.checkLogin()
		assert.NotNil(t, err)
		assert.Equal(t, ErrorUnauthenticated, err.Error())
	})

	t.Run("User State is LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.Set(userState, LoggedIn)

		err := ses.checkLogin()
		assert.Nil(t, err)
	})
}
