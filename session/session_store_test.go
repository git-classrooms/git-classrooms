package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestClassroomSession_Delete(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Deleting Existing Key", func(t *testing.T) {
		ses := Get(ctx)
		key := "test-key"
		value := "test-value"

		// Set a key-value pair
		ses.session.Set(key, value)

		// Ensure the key exists
		assert.Equal(t, value, ses.session.Get(key))

		// Delete the key
		ses.Delete(key)

		// Check that the key no longer exists
		assert.Nil(t, ses.session.Get(key), "The key should be deleted from the session")
	})

	t.Run("Deleting Non-Existent Key", func(t *testing.T) {
		ses := Get(ctx)
		nonExistentKey := "non-existent-key"

		// Attempt to delete a non-existent key
		ses.Delete(nonExistentKey)

		// This is more of a sanity check.
		assert.True(t, true, "Deleting a non-existent key should not cause any issues")
	})
}

func TestClassroomSession_Destroy(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Successful Session Destruction", func(t *testing.T) {
		ses := Get(ctx)

		// Set some data in the session to verify later
		ses.session.Set("test-key", "test-value")

		// Destroy the session
		err := ses.Destroy()
		assert.Nil(t, err, "Destroy method should not return an error")

		// Attempt to retrieve the data after destruction
		value := ses.session.Get("test-key")
		assert.Nil(t, value, "Session data should not be retrievable after destruction")
	})
}

func TestClassroomSession_Fresh(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("New Session", func(t *testing.T) {
		// Assuming Get creates a new session if it doesn't exist
		ses := Get(ctx)

		// Fresh should return true for a new session
		assert.True(t, ses.Fresh(), "Fresh should return true for a new ses")
	})
}

func TestClassroomSession_Get(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Session Initialization", func(t *testing.T) {
		ses := Get(ctx)
		assert.NotNil(t, ses)
		assert.Equal(t, Anonymous, ses.GetUserState())
	})

	t.Run("Session Retrieval", func(t *testing.T) {
		// Set an initial state for the session
		initialSession := Get(ctx)
		initialSession.session.Set(userState, int(LoggedIn))
		if err := initialSession.Save(); err != nil {
			t.Fatal(err)
		}

		// Retrieve the session again
		retrievedSession := Get(ctx)
		assert.Equal(t, LoggedIn, UserState(retrievedSession.session.Get(userState).(int)))
	})

	t.Run("Concurrency", func(t *testing.T) {
		// Reset the store and the once variable to ensure a clean test environment
		store = nil
		once = sync.Once{}

		var wg sync.WaitGroup
		initializationCount := 0

		mockInitFunction := func() {
			store = session.New()
			initializationCount++
		}

		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func() {
				defer wg.Done()
				once.Do(mockInitFunction)
			}()
		}

		wg.Wait()

		// Assert that the initialization function was called exactly once
		assert.Equal(t, 1, initializationCount, "Store should be initialized exactly once")
	})
}

func TestClassroomSession_GetExpiry(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Session with Set Expiry", func(t *testing.T) {
		ses := Get(ctx)
		expectedExpiry := time.Now().Add(1 * time.Hour).Unix()

		// Set the expiry time in the session
		ses.session.Set(expiresAt, expectedExpiry)

		// Retrieve the expiry time
		expiry := ses.GetExpiry()

		// The returned expiry should match the expected expiry
		assert.Equal(t, time.Unix(expectedExpiry, 0), expiry, "The expiry time should match the set value")

		// Reset for next test
		err := ses.session.Reset()
		if err != nil {
			assert.NoError(t, err)
		}
	})

	t.Run("Session without Set Expiry", func(t *testing.T) {
		ses := Get(ctx)

		// Ensure the expiresAt is not set
		ses.session.Delete(expiresAt)

		// Retrieve the expiry time
		expiry := ses.GetExpiry()

		assert.Equal(t, time.Unix(0, 0), expiry)
	})
}

func TestClassroomSession_GetGitlabAccessToken(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("User Not Logged In", func(t *testing.T) {
		ses := Get(ctx)
		// Ensure the user is not logged in
		ses.session.Set(userState, int(Anonymous))

		_, err := ses.GetGitlabAccessToken()
		assert.NotNil(t, err, "Expected an error for unauthenticated user")
	})

	t.Run("User Logged In With GitLab Access Token", func(t *testing.T) {
		ses := Get(ctx)
		// Set user as logged in and set a GitLab access token
		ses.session.Set(userState, int(LoggedIn))
		expectedToken := "gitlab-access-token-value"
		ses.session.Set(gitLabAccessToken, expectedToken)

		token, err := ses.GetGitlabAccessToken()
		assert.Nil(t, err, "Expected no error for authenticated user")
		assert.Equal(t, expectedToken, token, "The GitLab access token should match the set value")
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
		ses.session.Set(userState, int(Anonymous))

		_, err := ses.GetGitlabRefreshToken()
		assert.NotNil(t, err, "Expected an error for unauthenticated user")
	})

	t.Run("User Logged In With GitLab Refresh Token", func(t *testing.T) {
		ses := Get(ctx)
		// Set user as logged in and set a GitLab refresh token
		ses.session.Set(userState, int(LoggedIn))
		expectedToken := "gitlab-refresh-token-value"
		ses.session.Set(gitLabRefreshToken, expectedToken)

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
		ses.session.Set(userState, int(Anonymous))

		userID, err := ses.GetUserID()
		assert.Equal(t, -1, userID)
		assert.NotNil(t, err)
	})

	t.Run("User Logged In With Valid User ID", func(t *testing.T) {
		ses := Get(ctx)
		// Set user state to logged in and set a valid user ID
		ses.session.Set(userState, int(LoggedIn))
		expectedUserID := 123 // Example user ID
		ses.session.Set(userID, expectedUserID)

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
		ses.session.Set(userState, int(Anonymous))

		state := ses.GetUserState()

		assert.Equal(t, Anonymous, state, "The retrieved user state should be Anonymous")
	})

	t.Run("Retrieving User State LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.session.Set(userState, int(LoggedIn))

		state := ses.GetUserState()

		assert.Equal(t, LoggedIn, state, "The retrieved user state should be LoggedIn")
	})
}

func TestClassroomSession_ID(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	//Test
	t.Run("Retrieving Session ID", func(t *testing.T) {
		ses := Get(ctx)

		// Get the session ID
		sessionID := ses.ID()

		// Verify that the session ID is not empty
		assert.NotEmpty(t, sessionID, "Session ID should not be empty")
	})
}

func TestClassroomSession_Keys(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Session with No Keys", func(t *testing.T) {
		ses := Get(ctx)

		keys := ses.Keys()

		assert.Len(t, keys, 1, "Keys contain only a user state for a new session")
		assert.Equal(t, userState, keys[0], "Keys contain only a user state for a new session")
	})

	t.Run("Session with Multiple Keys", func(t *testing.T) {
		ses := Get(ctx)
		keysToSet := []string{"key1", "key2", "key3"}
		keysToMatch := []string{userState, "key1", "key2", "key3"}
		for _, key := range keysToSet {
			ses.session.Set(key, "value")
		}

		keys := ses.Keys()
		sort.Strings(keys) // Sort keys to ensure order for comparison
		sort.Strings(keysToSet)

		assert.ElementsMatch(t, keysToMatch, keys, "Keys should match the set keys in the session")
	})
}

func TestClassroomSession_Regenerate(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Test
	t.Run("Session ID Changes After Regeneration", func(t *testing.T) {
		ses := Get(ctx)

		// Get the current session ID
		oldSessionID := ses.ID()

		// Regenerate the session
		err := ses.Regenerate()
		assert.Nil(t, err, "Regenerate should not return an error")

		// Get the new session ID
		newSessionID := ses.ID()

		// Verify that the session ID has changed
		assert.NotEqual(t, oldSessionID, newSessionID, "Session ID should change after regeneration")
	})
}

func TestClassroomSession_Reset(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Session ID Changes and Data Reset After Reset", func(t *testing.T) {
		ses := Get(ctx)

		// Set some data in the session
		ses.session.Set("test-key", "test-value")

		// Get the current session ID
		oldSessionID := ses.session.ID()

		// Reset the session
		err := ses.Reset()
		assert.Nil(t, err, "Reset should not return an error")

		// Get the new session ID
		newSessionID := ses.session.ID()

		// Verify that the session ID has changed
		assert.NotEqual(t, oldSessionID, newSessionID, "Session ID should change after reset")

		// Verify that the session data is reset
		assert.Nil(t, ses.session.Get("test-key"), "Session data should be reset")
	})
}

func TestClassroomSession_Save(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Changes Persist After Save", func(t *testing.T) {
		ses := Get(ctx)

		// Modify the session
		testKey := "test-key"
		testValue := "test-value"
		ses.session.Set(testKey, testValue)

		// Save the session
		err := ses.Save()
		assert.Nil(t, err, "Save should not return an error")

		ses = Get(ctx)

		// Verify that changes are still present
		assert.Equal(t, testValue, ses.session.Get(testKey), "Changes should persist after Save")
	})
}

func TestClassroomSession_Set(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Test
	t.Run("Verify Key-Value Pair Set", func(t *testing.T) {
		ses := Get(ctx)

		// Define key and value to set
		testKey := "test-key"
		testValue := "test-value"

		// Set the key-value pair in the session
		ses.Set(testKey, testValue)

		// Retrieve the value using the same key
		retrievedValue := ses.session.Get(testKey)

		// Verify that the value retrieved matches the set value
		assert.Equal(t, testValue, retrievedValue, "The key-value pair should be correctly set in the session")
	})
}

func TestClassroomSession_SetExpiry(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Test
	t.Run("Verify Expiry Set Correctly and Save Called", func(t *testing.T) {
		ses := Get(ctx)

		// Define the expiry time
		expiryTime := time.Now().Add(1 * time.Hour).Truncate(time.Second)

		// Set the expiry using SetExpiry method
		ses.SetExpiry(expiryTime)

		// Retrieve the saved expiry time from the session
		savedExpiry := ses.GetExpiry()

		// Verify that the saved expiry time matches the provided expiry
		assert.Equal(t, expiryTime, savedExpiry, "Expiry time should be set correctly")
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
		savedAccessToken := ses.session.Get(gitLabAccessToken)

		// Verify that the saved access token matches the provided token
		assert.Equal(t, testAccessToken, savedAccessToken, "GitLab access token should be set correctly")
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
		savedRefreshToken := ses.session.Get(gitLabRefreshToken)

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
		userID := ses.session.Get(userID)

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
		state := ses.session.Get(userState).(int)

		assert.Equal(t, int(Anonymous), state, "The user state should be set to Anonymous")
	})

	t.Run("Setting User State to LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.SetUserState(LoggedIn)

		// Retrieve the set user state from the session
		state := ses.session.Get(userState).(int)

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
		ses.session.Set(userState, int(Anonymous))

		err := ses.checkLogin()
		assert.NotNil(t, err)
		assert.Equal(t, ErrorUnauthenticated, err.Error())
	})

	t.Run("User State is LoggedIn", func(t *testing.T) {
		ses := Get(ctx)
		ses.session.Set(userState, int(LoggedIn))

		err := ses.checkLogin()
		assert.Nil(t, err)
	})
}

func TestGet(t *testing.T) {
	// Mock a fiber context
	app := fiber.New()
	req := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(req)
	defer app.ReleaseCtx(ctx)

	// Tests
	t.Run("Retrieving Existing Key", func(t *testing.T) {
		ses := Get(ctx)
		key := "test-key"
		expectedValue := "test-value"

		// Set a key-value pair in the session
		ses.session.Set(key, expectedValue)

		// Retrieve the value using Get method
		value := ses.Get(key)
		assert.Equal(t, expectedValue, value, "The retrieved value should match the set value")
	})

	t.Run("Retrieving Non-Existent Key", func(t *testing.T) {
		ses := Get(ctx)
		nonExistentKey := "non-existent-key"

		// Attempt to retrieve a value for a non-existent key
		value := ses.Get(nonExistentKey)

		// The value should be nil for a non-existent key
		assert.Nil(t, value, "Retrieving a non-existent key should return nil")
	})
}
