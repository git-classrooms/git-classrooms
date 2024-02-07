package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func TestClassroomSession_GetGitlabAccessToken(t *testing.T) {
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
		ses.Set(userState, int(Anonymous))

		_, err := ses.GetGitlabAccessToken()
		assert.NotNil(t, err, "Expected an error for unauthenticated user")
	})

	t.Run("User Logged In With GitLab Access Token", func(t *testing.T) {
		ses := Get(ctx)
		// Set user as logged in and set a GitLab access token
		ses.Set(userState, int(LoggedIn))
		expectedToken := "gitlab-access-token-value"
		ses.Set(gitLabAccessToken, expectedToken)

		token, err := ses.GetGitlabAccessToken()
		assert.Nil(t, err, "Expected no error for authenticated user")
		assert.Equal(t, expectedToken, token, "The GitLab access token should match the set value")
	})
}

func TestClassroomSession_GetAccessTokenExpiry(t *testing.T) {
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
		ses.Set(userState, int(Anonymous))

		_, err := ses.GetAccessTokenExpiry()
		assert.NotNil(t, err, "Expected an error for unauthenticated user")
	})

	t.Run("User Logged In With GitLab Access Token Expiry", func(t *testing.T) {
		ses := Get(ctx)
		// Set user as logged in and set a GitLab access token
		ses.Set(userState, int(LoggedIn))
		expectedTime := time.Unix(1000, 0)
		ses.Set(expiresAt, expectedTime.Unix())

		exp, err := ses.GetAccessTokenExpiry()
		assert.Nil(t, err, "Expected no error for authenticated user")
		assert.Equal(t, expectedTime, exp, "The GitLab access token expiry should match the set value")
	})
}

func TestClassroomSession_GetGitlabRefreshToken(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("User Not Logged In", func(t *testing.T) {
		ses := Get(ctx)
		// Ensure the user is not logged in
		ses.Set(userState, int(Anonymous))

		_, err := ses.GetGitlabRefreshToken()
		assert.NotNil(t, err, "Expected an error for unauthenticated user")
	})

	t.Run("User Logged In With GitLab Refresh Token", func(t *testing.T) {
		ses := Get(ctx)
		// Set user as logged in and set a GitLab refresh token
		ses.Set(userState, int(LoggedIn))
		expectedToken := "gitlab-refresh-token-value"
		ses.Set(gitLabRefreshToken, expectedToken)

		token, err := ses.GetGitlabRefreshToken()
		assert.Nil(t, err, "Expected no error for authenticated user")
		assert.Equal(t, expectedToken, token, "The GitLab refresh token should match the set value")
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
		ses.Set(userState, int(Anonymous))

		userID, err := ses.GetUserID()
		assert.Equal(t, -1, userID)
		assert.NotNil(t, err)
	})

	t.Run("User Logged In With Valid User ID", func(t *testing.T) {
		ses := Get(ctx)
		// Set user state to logged in and set a valid user ID
		ses.Set(userState, int(LoggedIn))
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
		ses.Set(userState, int(Anonymous))

		state := ses.GetUserState()

		assert.Equal(t, Anonymous, state, "The retrieved user state should be Anonymous")
	})

	t.Run("Retrieving User State LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.Set(userState, int(LoggedIn))

		state := ses.GetUserState()

		assert.Equal(t, LoggedIn, state, "The retrieved user state should be LoggedIn")
	})
}

func TestClassroomSession_SetGitlabAccessToken(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Test
	t.Run("Verify Access Token Set", func(t *testing.T) {
		ses := Get(ctx)

		// Define the GitLab access token
		testAccessToken := "test-access-token"

		// Set the access token using the SetGitlabAccessToken method
		ses.SetGitlabAccessToken(testAccessToken)

		// Retrieve the saved access token from the session
		savedAccessToken := ses.Get(gitLabAccessToken)

		// Verify that the saved access token matches the provided token
		assert.Equal(t, testAccessToken, savedAccessToken, "GitLab access token should be set correctly")
	})
}

func TestClassroomSession_SetAccessTokenExpiry(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Test
	t.Run("Verify Access Token Expiry Set", func(t *testing.T) {
		ses := Get(ctx)

		// Define the GitLab access token
		testAccessTokenTime := time.Unix(1000, 0)

		// Set the access token using the SetGitlabAccessToken method
		ses.SetAccessTokenExpiry(testAccessTokenTime)

		// Retrieve the saved access token from the session
		savedAccessTokenTime := ses.Get(expiresAt)

		// Verify that the saved access token matches the provided token
		assert.Equal(t, testAccessTokenTime.Unix(), savedAccessTokenTime, "Access token Expiry should be set correctly")
	})
}

func TestClassroomSession_SetGitlabRefreshToken(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Test
	t.Run("Verify Refresh Token Set", func(t *testing.T) {
		ses := Get(ctx)

		// Define the GitLab refresh token
		testRefreshToken := "test-refresh-token"

		// Set the refresh token using the SetGitlabRefreshToken method
		ses.SetGitlabRefreshToken(testRefreshToken)

		// Retrieve the saved refresh token from the session
		savedRefreshToken := ses.Get(gitLabRefreshToken)

		// Verify that the saved refresh token matches the provided token
		assert.Equal(t, testRefreshToken, savedRefreshToken, "GitLab refresh token should be set correctly")
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
		state := ses.Get(userState).(int)

		assert.Equal(t, int(Anonymous), state, "The user state should be set to Anonymous")
	})

	t.Run("Setting User State to LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.SetUserState(LoggedIn)

		// Retrieve the set user state from the session
		state := ses.Get(userState).(int)

		assert.Equal(t, int(LoggedIn), state, "The user state should be set to LoggedIn")
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
		ses.Set(userState, int(Anonymous))

		err := ses.checkLogin()
		assert.NotNil(t, err)
		assert.Equal(t, ErrorUnauthenticated, err.Error())
	})

	t.Run("User State is LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.Set(userState, int(LoggedIn))

		err := ses.checkLogin()
		assert.Nil(t, err)
	})
}
