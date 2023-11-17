package go_gitlab_repo_test

import (
	"backend/api/repository/go_gitlab_repo"
	"backend/model"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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

func LoadCredentialsFromFile(path string) (*GitlabCredentials, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var credentials GitlabCredentials
	if err := yaml.NewDecoder(f).Decode(&credentials); err != nil {
		return nil, err
	}

	return &credentials, nil
}

func TestGoGitlabRepo(t *testing.T) {
	credentialsPath := flag.String("credentials", "credentials.yml", "The path to your gitlab credentials file")
	flag.Parse()
	credentials, err := LoadCredentialsFromFile(*credentialsPath)
	if err != nil {
		t.Fatal(err)
	}

	repo := go_gitlab_repo.NewGoGitlabRepo()

	// hs-flensburg.dev
	// user
	// Anwendungsschlüssel
	//

	t.Run("LoginByToken", func(t *testing.T) {
		user, err := repo.Login(credentials.Token, credentials.Username)

		assert.NoError(t, err)
		assert.Equal(t, credentials.ID, user.ID)
		assert.Equal(t, credentials.Username, user.Username)
		assert.Equal(t, credentials.Name, user.Name)
		assert.Equal(t, credentials.WebUrl, user.WebUrl)
		// assert.Equal(t, credentials.Email, user.Email) // TODO no emails available yet
	})

	_, err = repo.Login(credentials.Token, credentials.Username)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetProjectById", func(t *testing.T) {
		project, err := repo.GetProjectById(685)

		assert.NoError(t, err)
		assert.Equal(t, 685, project.ID)
		assert.Equal(t, "MyTestProject", project.Name)
		assert.Equal(t, "https://gitlab.hs-flensburg.de/mytestgroup/mytestproject", project.WebUrl)
		assert.Equal(t, model.Private, project.Visibility)
	})

	t.Run("GetProjectById", func(t *testing.T) {
		projects, err := repo.GetAllProjects()

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(projects), 1)
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
		Group, err := repo.GetGroupById(1237)

		assert.NoError(t, err)
		assert.Equal(t, 1237, Group.ID)
		assert.Equal(t, "MyTestGroup", Group.Name)
		assert.Equal(t, "https://gitlab.hs-flensburg.de/groups/mytestgroup", Group.WebUrl)
		assertContainUser(t, model.User{ID: credentials.ID, Name: credentials.Name, Username: credentials.Username, WebUrl: credentials.WebUrl}, Group.Member)
		assertContainProject(t, model.Project{ID: 685, Name: "MyTestProject", WebUrl: "https://gitlab.hs-flensburg.de/mytestgroup/mytestproject", Description: ""}, Group.Projects)
	})

	t.Run("GetAllProjectsOfGroup", func(t *testing.T) {
		projects, err := repo.GetAllProjectsOfGroup(1237)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(projects), 1)
	})

	t.Run("GetAllUsersOfGroup", func(t *testing.T) {
		users, err := repo.GetAllUsersOfGroup(1237)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		users, err := repo.GetAllUsers()

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
	})
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
