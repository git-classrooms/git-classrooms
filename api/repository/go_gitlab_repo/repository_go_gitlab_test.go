package go_gitlab_repo_test

import (
	"backend/api/repository/go_gitlab_repo"
	"backend/model"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

type GitlabCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ID       int    `yaml:"id"`
	Email    string `yaml:"email"`
	WebUrl   string `yaml:"webUrl"`
	Name     string `yaml:"name"`
	Token    string `yaml:"token"`
}

func LoadCredentialsFromEnv() (*GitlabCredentials, error) {
	_ = godotenv.Load(".env", ".env.local")

	idStr := os.Getenv("GO_GITLAB_TEST_USER_ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	credentials := GitlabCredentials{
		Username: os.Getenv("GO_GITLAB_TEST_USERNAME"),
		Password: os.Getenv("GO_GITLAB_TEST_PASSWORD"),
		ID:       id,
		Email:    os.Getenv("GO_GITLAB_TEST_EMAIL"),
		WebUrl:   os.Getenv("GO_GITLAB_TEST_WEB_URL"),
		Name:     os.Getenv("GO_GITLAB_TEST_NAME"),
		Token:    os.Getenv("GO_GITLAB_TEST_TOKEN"),
	}

	return &credentials, nil
}

func TestGoGitlabRepo(t *testing.T) {
	credentials, err := LoadCredentialsFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	repo := go_gitlab_repo.NewGoGitlabRepo()

	t.Run("LoginByToken", func(t *testing.T) {
		user, err := repo.Login(credentials.Token, credentials.Username)

		assert.NoError(t, err)
		assert.Equal(t, credentials.ID, user.ID)
		assert.Equal(t, credentials.Username, user.Username)
		assert.Equal(t, credentials.Name, user.Name)
		assert.Equal(t, credentials.WebUrl, user.WebUrl)
		// assert.Equal(t, credentials.Email, user.Email) // TODO emails not available with personal access tokens, but should be with session tokens
	})

	_, err = repo.Login(credentials.Token, credentials.Username)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetAllProjects", func(t *testing.T) {
		projects, err := repo.GetAllProjects()

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(projects), 1)
	})

	t.Run("GetProjectById", func(t *testing.T) {
		project, err := repo.GetProjectById(2)

		assert.NoError(t, err)
		assert.Equal(t, 2, project.ID)
		assert.Equal(t, "IntegrationTestsProject2", project.Name)
		assert.Equal(t, "https://hs-flensburg.dev/IntegrationTestsUser2/integrationtestsproject2", project.WebUrl)
		assert.Equal(t, model.Internal, project.Visibility)
	})

	t.Run("GetUserById", func(t *testing.T) {
		user, err := repo.GetUserById(credentials.ID)

		assert.NoError(t, err)
		assert.Equal(t, credentials.ID, user.ID)
		assert.Equal(t, credentials.Username, user.Username)
		assert.Equal(t, credentials.Name, user.Name)
		assert.Equal(t, credentials.WebUrl, user.WebUrl)
		// assert.Equal(t, credentials.Email, user.Email) // TODO no emails available yet
	})

	t.Run("GetGroupById", func(t *testing.T) {
		Group, err := repo.GetGroupById(15)

		assert.NoError(t, err)
		assert.Equal(t, 15, Group.ID)
		assert.Equal(t, "IntegrationsTestGroup1", Group.Name)
		assert.Equal(t, "https://hs-flensburg.dev/groups/integrationstestgroup11", Group.WebUrl)
		assertContainUser(t, model.User{ID: credentials.ID, Name: credentials.Name, Username: credentials.Username, WebUrl: credentials.WebUrl}, Group.Member)
		assertContainProject(t, model.Project{ID: 3, Name: "IntegrationTestsProject3", WebUrl: "https://hs-flensburg.dev/integrationstestgroup11/integrationtestsproject3", Description: ""}, Group.Projects)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		users, err := repo.GetAllUsers()

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
	})

	t.Run("GetAllGroups", func(t *testing.T) {
		groups, err := repo.GetAllGroups()

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(groups), 1)
	})

	t.Run("GetAllProjectsOfGroup", func(t *testing.T) {
		projects, err := repo.GetAllProjectsOfGroup(15)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(projects), 1)
	})

	t.Run("GetAllUsersOfGroup", func(t *testing.T) {
		users, err := repo.GetAllUsersOfGroup(15)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
	})

	/*
		Seems to not work with Push rules
		Other options: https://docs.gitlab.com/ee/api/protected_branches.html#protect-repository-branches
		t.Run("Push rules", func(t *testing.T) {
			client := createTestClient(t, credentials.Token)

			err := repo.DenyPushingToProject(3)
			assert.NoError(t, err)

			assert.True(t, containsPushRule(t, client, 3))

			err = repo.AllowPushingToProject(685)
			assert.NoError(t, err)

			assert.False(t, containsPushRule(t, client, 3))
		})
	*/
}

func createTestClient(t *testing.T, token string) *gitlab.Client {
	cli, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.hs-flensburg.de"))
	if err != nil {
		t.Errorf("Could not create extra client for test")
	}
	return cli
}

func containsPushRule(t *testing.T, client *gitlab.Client, projectId int) bool {
	rule, _, err := client.Projects.GetProjectPushRules(projectId)
	if err != nil {
		t.Errorf("Could not GetProjectPushRules")
	}

	return rule.AuthorEmailRegex == "DenyPushing"
}

func assertContainUser(t *testing.T, expectedUser model.User, users []model.User) {
	for _, user := range users {
		if user == expectedUser {
			return
		}
	}

	t.Errorf("User not found")
}

func assertContainProject(t *testing.T, expectedProject model.Project, projects []model.Project) {
	for _, project := range projects {
		if project.ID == expectedProject.ID && project.Name == expectedProject.Name && project.WebUrl == expectedProject.WebUrl && project.Description == expectedProject.Description {
			return
		}
	}

	t.Errorf("Project not found")
}
