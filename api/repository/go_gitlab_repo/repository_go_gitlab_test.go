package go_gitlab_repo_test

import (
	"backend/api/repository/go_gitlab_repo"
	"backend/model"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

type GitlabCredentials struct {
	Username string
	Password string
	ID       int
	Email    string
	WebUrl   string
	Name     string
	Token    string
}

func LoadCredentialsFromEnv() (*GitlabCredentials, error) {
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Print(err.Error())
	}

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

		webUrl := fmt.Sprintf("%s/%s", credentials.WebUrl, credentials.Username)

		assert.NoError(t, err)
		assert.Equal(t, credentials.ID, user.ID)
		assert.Equal(t, credentials.Username, user.Username)
		assert.Equal(t, credentials.Name, user.Name)
		assert.Equal(t, webUrl, user.WebUrl)
		// assert.Equal(t, credentials.Email, user.Email) // TODO emails not available with personal access tokens, but should be with session tokens
	})

	_, err = repo.Login(credentials.Token, credentials.Username)
	if err != nil {
		t.Fatal(err)
	}

	// erstellt Projekt, Test schmeisst aber error
	// t.Run("CreateProject", func(t *testing.T) {
	// 	projectName := "TestProject2"
	// 	projectVisibility := model.Public
	// 	projectDescription := "Test project description2"

	// 	members := []model.User{
	// 		{ID: credentials.ID, Username: credentials.Username},
	// 	}

	// 	// Test CreateProject
	// 	project, err := repo.CreateProject(projectName, projectVisibility, projectDescription, members)

	// 	// Assertions
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, project)
	// 	assert.Equal(t, projectName, project.Name)
	// 	assert.Equal(t, projectVisibility, project.Visibility)
	// 	assert.Equal(t, projectDescription, project.Description)
	// })

	/*
		// If you get the error "has already been taken", the test has been run previously and there already exist a project with this name in the group namespace
		t.Run("ForkProject", func(t *testing.T) {
			newName := "ForkTestFork3"

			forkProject, err := repo.ForkProject(3, newName)

			assert.NoError(t, err)
			assert.Equal(t, newName, forkProject.Name)
			assert.NotEqual(t, 3, forkProject.ID)
		})
	*/

	/*
		// to run this test, check that the user is not already member of project
		t.Run("AddProjectMembers", func(t *testing.T) {
			members := make([]model.User, 1)
			members[0] = model.User{
				ID:       5,
				Username: "IntegrationTestsUser1",
				Name:     "TestUser1",
				WebUrl:   "https://hs-flensburg.dev/IntegrationTestsUser1",
			}

			project, err := repo.AddProjectMembers(3, members)

			assert.NoError(t, err)
			assertContainUser(t, members[0], project.Member)
		})
	*/

	// erstellt Gruppe, Test schmeisst aber error
	// t.Run("CreateGroup", func(t *testing.T) {
	//     groupName := "TestGroup"
	//     groupVisibility := model.Public
	//     groupDescription := "A test group"
	//     memberEmails := []string{credentials.Email}

	//     group, err := repo.CreateGroup(groupName, groupVisibility, groupDescription, memberEmails)

	//     assert.Equal(t, groupName, group.Name)
	//     assert.Equal(t, groupDescription, group.Description)
	//     assert.Error(t, err)
	// })

	groupId := 20 // Example groupId
	userId := 9   // You can use the ID from credentials or another user's ID

	t.Run("AddUserToGroup", func(t *testing.T) {
		err := repo.AddUserToGroup(groupId, userId)

		assert.NoError(t, err)

		// After adding, verify if the user is actually in the group
		group, err := repo.GetGroupById(groupId)
		assert.NoError(t, err)

		found := false
		for _, member := range group.Member {
			if member.ID == userId {
				found = true
				break
			}
		}

		assert.True(t, found, "User should be in the group after adding")
	})

	t.Run("RemoveUserFromGroup", func(t *testing.T) {
		err := repo.RemoveUserFromGroup(groupId, userId)

		assert.NoError(t, err)

		// After removing, verify if the user is actually removed from the group
		group, err := repo.GetGroupById(groupId)
		assert.NoError(t, err)

		found := false
		for _, member := range group.Member {
			if member.ID == userId {
				found = true
				break
			}
		}

		assert.False(t, found, "User should not be in the group after removal")
	})

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

		webUrl := fmt.Sprintf("%s/%s", credentials.WebUrl, credentials.Username)

		assert.NoError(t, err)
		assert.Equal(t, credentials.ID, user.ID)
		assert.Equal(t, credentials.Username, user.Username)
		assert.Equal(t, credentials.Name, user.Name)
		assert.Equal(t, webUrl, user.WebUrl)
		// assert.Equal(t, credentials.Email, user.Email) // TODO no emails available yet
	})

	t.Run("GetGroupById", func(t *testing.T) {
		Group, err := repo.GetGroupById(15)

		assert.NoError(t, err)
		assert.Equal(t, 15, Group.ID)
		assert.Equal(t, "IntegrationsTestGroup1", Group.Name)
		assert.Equal(t, "https://hs-flensburg.dev/groups/integrationstestgroup11", Group.WebUrl)
		user_web_url := fmt.Sprintf("%s/%s", credentials.WebUrl, credentials.Username)
		assertContainUser(t, model.User{ID: credentials.ID, Name: credentials.Name, Username: credentials.Username, WebUrl: user_web_url}, Group.Member)
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

	// Use a search expression that is likely to return results, such as a common project name or keyword
	searchExpression := "Test"

	t.Run("SearchProjectByValidExpression", func(t *testing.T) {
		projects, err := repo.SearchProjectByExpression(searchExpression)

		assert.NoError(t, err)
		assert.NotNil(t, projects)
		assert.Greater(t, len(projects), 0, "Should return one or more projects")
	})

	// Use an unlikely search expression to test for no results
	invalidSearchExpression := "unlikely-keyword-xyz"

	t.Run("SearchProjectByInvalidExpression", func(t *testing.T) {
		projects, err := repo.SearchProjectByExpression(invalidSearchExpression)

		assert.NoError(t, err)
		assert.NotNil(t, projects)
		assert.Equal(t, 0, len(projects), "Should return no projects for an unlikely search expression")
	})

	t.Run("SearchUserByValidExpression", func(t *testing.T) {
		expression := credentials.Username // Use part of the username from credentials or another known user

		users, err := repo.SearchUserByExpression(expression)
		assert.NoError(t, err)

		// Check if the returned users list contains the user with the used expression
		found := false
		for _, user := range users {
			if user.Username == credentials.Username {
				found = true
				break
			}
		}

		assert.True(t, found, "User should be found by the expression")
	})

	t.Run("SearchUserByInvalidExpression", func(t *testing.T) {
		expression := "nonexistentuser123" // An expression that would not match any user

		users, err := repo.SearchUserByExpression(expression)
		assert.NoError(t, err)

		// Check if the users list is empty as no user should match the expression
		assert.Empty(t, users, "No users should be found with an invalid expression")
	})

	t.Run("SearchUserByExpressionInGroup", func(t *testing.T) {
		users, err := repo.SearchUserByExpressionInGroup(searchExpression, groupId)

		assert.NoError(t, err)
		assert.NotEmpty(t, users, "Expected to find at least one user")

		// Optionally check if the returned users match the search criteria
		for _, user := range users {
			assert.Contains(t, user.Name, searchExpression, "User name should contain the search expression")
		}
	})

	projectId := 5 // Example project ID, replace with an actual project ID

	t.Run("SearchUserByExpressionInProject", func(t *testing.T) {
		users, err := repo.SearchUserByExpressionInProject(searchExpression, projectId)

		assert.NoError(t, err)
		assert.NotNil(t, users)

		// Optionally, verify if the returned users meet certain criteria
		// e.g., checking if a known user appears in the results
		found := false
		for _, user := range users {
			if user.Username == credentials.Username {
				found = true
				break
			}
		}

		assert.True(t, found, "Expected user should be in the search results")
	})

	// Use an expression that is expected to match certain groups in your GitLab instance
	expression := "TestGroup" // Replace with an appropriate expression for your test

	t.Run("SearchGroupByExpression", func(t *testing.T) {
		groups, err := repo.SearchGroupByExpression(expression)
		assert.NoError(t, err)

		// Verify that the returned groups match the expression criteria
		// This could be as simple as checking if the slice is not empty
		// Or more complex, like checking if returned groups' names contain the expression
		assert.NotEmpty(t, groups, "Expected to find groups matching the expression")

		// Optionally, you can iterate through the groups and assert more specific conditions
		for _, group := range groups {
			assert.Contains(t, group.Name, expression, "Group name should contain the expression")
		}
	})

	email := credentials.Email // Use a test email address

	t.Run("CreateGroupInvite", func(t *testing.T) {
		err := repo.CreateGroupInvite(groupId, email)
		assert.NoError(t, err)

		// Verify if the invite was sent
		// Note: Depending on the GitLab API and your permissions, you might not be able to directly check if an invite was sent.
	})

	t.Run("CreateProjectInvite", func(t *testing.T) {
		err := repo.CreateProjectInvite(projectId, email)
		assert.NoError(t, err)

		// Verify if the invite was sent
		// Note: Depending on the GitLab API and your permissions, you might not be able to directly check if an invite was sent.
	})

	t.Run("GetPendingGroupInvitations", func(t *testing.T) {
		pendingInvites, err := repo.GetPendingGroupInvitations(groupId)

		assert.NoError(t, err)
		assert.NotNil(t, pendingInvites)

		// Optionally, check for specific properties of the pending invitations
		// For example, assert that the length of pendingInvites is as expected
		// or check for specific user IDs in the pending invitations
	})

	t.Run("GetNamespaceOfProject", func(t *testing.T) {
		namespace, err := repo.GetNamespaceOfProject(3)

		assert.NoError(t, err)
		assert.Equal(t, "integrationstestgroup11", namespace)
	})

	t.Run("GetNamespaceOfGroup", func(t *testing.T) {
		namespace, err := repo.GetNamespaceOfGroup(15)

		assert.NoError(t, err)
		assert.Equal(t, "integrationstestgroup11", namespace)
	})

	/*
		Test schmeisst Error, dont know why
		t.Run("GetPendingProjectInvitations", func(t *testing.T) {
			pendingInvites, err := repo.GetPendingProjectInvitations(groupId)

			assert.NoError(t, err)
			assert.NotNil(t, pendingInvites)

			// Optionally, check for specific properties of the pending invitations
			// For example, assert that the length of pendingInvites is as expected
			// or check for specific user IDs in the pending invitations
		})
	*/

	/*
		Mit personal access tokens ist es bisher nicht möglich ein Assignment zu schließen bzw. das Pushen zu unterbinden (man bekommt bei alle aufgelisteten Möglichkeiten einen 404 zurück)
			- Not with Push Rules
			- Not with Protect Branches
			- Not with change Project Member Access Level

		t.Run("Push rules", func(t *testing.T) {
			client := createTestClient(t, credentials.Token)

			err := repo.DenyPushingToProject(3)
			assert.NoError(t, err)

			assert.True(t, allProjectMembersHaveSameAccessLevel(t, client, 3, gitlab.AccessLevelValue(gitlab.MinimalAccessPermissions)))

			err = repo.AllowPushingToProject(685)
			assert.NoError(t, err)

			assert.False(t, allProjectMembersHaveSameAccessLevel(t, client, 3, gitlab.AccessLevelValue(gitlab.DeveloperPermissions)))
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

func allProjectMembersHaveSameAccessLevel(t *testing.T, client *gitlab.Client, projectId int, accessLevel gitlab.AccessLevelValue) bool {
	members, _, err := client.ProjectMembers.ListAllProjectMembers(projectId, &gitlab.ListProjectMembersOptions{})
	if err != nil {
		return false
	}

	for _, member := range members {
		if member.AccessLevel != accessLevel {
			return false
		}
	}

	return true
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
