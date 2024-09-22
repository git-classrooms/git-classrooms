// Reference to Go Gitlab Documentation: https://pkg.go.dev/github.com/xanzy/go-gitlab#section-documentation
package gitlab

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	goGitlab "github.com/xanzy/go-gitlab"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

// GitlabRepo manages interactions with the GitLab API.
type GitlabRepo struct {
	client      *goGitlab.Client
	config      gitlabConfig.Config
	isConnected bool
}

// NewGitlabRepo initializes a new GitlabRepo with the given configuration.
func NewGitlabRepo(config gitlabConfig.Config) *GitlabRepo {
	return &GitlabRepo{client: nil, config: config, isConnected: false}
}

// Login authenticates with GitLab using an OAuth token.
func (repo *GitlabRepo) Login(token string) error {
	// With oauth tokens we need the OAuthClient to make requests
	// TODO: But all tests act with a personal token, we just use the normal client for a while
	cli, err := goGitlab.NewOAuthClient(token, goGitlab.WithBaseURL(repo.config.GetURL()))
	if err != nil {
		return err
	}
	repo.client = cli
	return nil
}

// GroupAccessLogin logs in using a group access token.
func (repo *GitlabRepo) GroupAccessLogin(token string) error {
	cli, err := goGitlab.NewClient(token, goGitlab.WithBaseURL(repo.config.GetURL()))
	if err != nil {
		return err
	}
	repo.client = cli
	return nil
}

// GetCurrentUser fetches the current user from GitLab.
func (repo *GitlabRepo) GetCurrentUser() (*model.User, error) {
	repo.assertIsConnected()

	gitlabUser, _, err := repo.client.Users.CurrentUser()
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	classroomUser := UserFromGoGitlab(*gitlabUser)
	classroomUser.Avatar.FallbackAvatarURL, _ = repo.GetPublicAvatarByMail(classroomUser.Email)
	return classroomUser, nil
}

// CreateProject creates a new project with the given name, visibility, and members.
func (repo *GitlabRepo) CreateProject(name string, visibility model.Visibility, description string, members []model.User) (*model.Project, error) {
	repo.assertIsConnected()

	opts := &goGitlab.CreateProjectOptions{
		Name:        goGitlab.String(name),
		Visibility:  goGitlab.Visibility(VisibilityFromModel(visibility)),
		Description: goGitlab.String(description),
	}

	gitlabProject, _, err := repo.client.Projects.CreateProject(opts)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.AddProjectMembers(gitlabProject.ID, members)
}

// ForkProject forks an existing project with the specified parameters.
func (repo *GitlabRepo) ForkProject(projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error) {
	repo.assertIsConnected()

	opts := &goGitlab.ForkProjectOptions{
		Name:                          goGitlab.String(name),
		Path:                          goGitlab.String(convertToGitLabPath(name)),
		NamespaceID:                   goGitlab.Int(namespaceId),
		Visibility:                    goGitlab.Visibility(VisibilityFromModel(visibility)),
		Description:                   goGitlab.String(description),
		MergeRequestDefaultTargetSelf: goGitlab.Bool(true),
	}

	gitlabProject, _, err := repo.client.Projects.ForkProject(projectId, opts)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProject(gitlabProject)
}

// ForkProjectWithOnlyDefaultBranch forks a project with only the default branch.
func (repo *GitlabRepo) ForkProjectWithOnlyDefaultBranch(projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error) {
	repo.assertIsConnected()

	templateProject, _, err := repo.client.Projects.GetProject(projectId, &goGitlab.GetProjectOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	opts := &goGitlab.ForkProjectOptions{
		Name:                          goGitlab.String(name),
		Path:                          goGitlab.String(convertToGitLabPath(name)),
		NamespaceID:                   goGitlab.Int(namespaceId),
		Visibility:                    goGitlab.Visibility(VisibilityFromModel(visibility)),
		Description:                   goGitlab.String(description),
		MergeRequestDefaultTargetSelf: goGitlab.Bool(true),
	}

	gitlabProject, _, err := repo.client.Projects.ForkProject(projectId, opts, func(r *retryablehttp.Request) error {
		query := r.URL.Query()
		query.Add("branches", templateProject.DefaultBranch)
		r.URL.RawQuery = query.Encode()
		return nil
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProject(gitlabProject)
}

// CreateBranch creates a new branch from an existing one.
func (repo *GitlabRepo) CreateBranch(projectId int, branchName string, fromBranch string) (*model.Branch, error) {
	repo.assertIsConnected()

	opts := &goGitlab.CreateBranchOptions{
		Branch: goGitlab.String(branchName),
		Ref:    goGitlab.String(fromBranch),
	}

	branch, _, err := repo.client.Branches.CreateBranch(projectId, opts)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}
	return BranchFromGoGitlab(branch), nil
}

// ProtectBranch protects a branch with the specified access level.
func (repo *GitlabRepo) ProtectBranch(projectId int, branchName string, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	opts := &goGitlab.ProtectRepositoryBranchesOptions{
		Name:             goGitlab.String(branchName),
		AllowForcePush:   goGitlab.Bool(false),
		PushAccessLevel:  goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
		MergeAccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
	}

	_, _, err := repo.client.ProtectedBranches.ProtectRepositoryBranches(projectId, opts)
	return ErrorFromGoGitlab(err)
}

// UnprotectBranch removes the protection from a branch.
func (repo *GitlabRepo) UnprotectBranch(projectId int, branchName string) error {
	repo.assertIsConnected()

	_, err := repo.client.ProtectedBranches.UnprotectRepositoryBranches(projectId, branchName)
	return ErrorFromGoGitlab(err)
}

// CreateMergeRequest creates a merge request between branches.
func (repo *GitlabRepo) CreateMergeRequest(projectId int, sourceBranch string, targetBranch string, title string, description string, assigneeId int, reviewerId int) error {
	repo.assertIsConnected()

	reviewers := []int{reviewerId}

	opts := &goGitlab.CreateMergeRequestOptions{
		Title:        goGitlab.String(title),
		SourceBranch: goGitlab.String(sourceBranch),
		TargetBranch: goGitlab.String(targetBranch),
		Description:  goGitlab.String(description),
		AssigneeID:   goGitlab.Int(assigneeId),
		ReviewerIDs:  &reviewers,
	}

	_, _, err := repo.client.MergeRequests.CreateMergeRequest(projectId, opts)
	return ErrorFromGoGitlab(err)
}

// ProtectedBranchExists checks if a branch is protected in a project.
func (repo *GitlabRepo) ProtectedBranchExists(projectId int, branchName string) (bool, error) {
	repo.assertIsConnected()

	_, response, err := repo.client.ProtectedBranches.GetProtectedBranch(projectId, branchName)
	if err != nil {
		if response.StatusCode == 404 {
			return false, nil
		}
		return false, ErrorFromGoGitlab(err)
	}

	return true, nil
}

// BranchExists checks if a branch exists in a project.
func (repo *GitlabRepo) BranchExists(projectId int, branchName string) (bool, error) {
	repo.assertIsConnected()

	_, response, err := repo.client.Branches.GetBranch(projectId, branchName)
	if err != nil {
		if response.StatusCode == 404 {
			return false, nil
		}
		return false, ErrorFromGoGitlab(err)
	}

	return true, nil
}

// AddProjectMember adds a user to a project with the specified access level.
func (repo *GitlabRepo) AddProjectMember(projectId int, userId int, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	_, _, err := repo.client.ProjectMembers.AddProjectMember(projectId, &goGitlab.AddProjectMemberOptions{
		UserID:      &userId,
		AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
	})
	return ErrorFromGoGitlab(err)
}

// RemoveUserFromProject removes a user from a project.
func (repo *GitlabRepo) RemoveUserFromProject(projectId int, userId int) error {
	repo.assertIsConnected()

	_, err := repo.client.ProjectMembers.DeleteProjectMember(projectId, userId)
	return ErrorFromGoGitlab(err)
}

// AddProjectMembers adds multiple users to a project.
func (repo *GitlabRepo) AddProjectMembers(projectId int, members []model.User) (*model.Project, error) {
	repo.assertIsConnected()

	for _, member := range members {
		_, _, err := repo.client.ProjectMembers.AddProjectMember(projectId, &goGitlab.AddProjectMemberOptions{
			UserID:      &member.ID,
			AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		})
		if err != nil {
			return nil, ErrorFromGoGitlab(err)
		}
	}

	return repo.GetProjectById(projectId)
}

// ChangeUserAccessLevelInProject changes the access level of a user in a project.
func (repo *GitlabRepo) ChangeUserAccessLevelInProject(projectId int, userId int, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	_, _, err := repo.client.ProjectMembers.EditProjectMember(
		projectId,
		userId,
		&goGitlab.EditProjectMemberOptions{AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel))},
	)

	return ErrorFromGoGitlab(err)
}

// GetAccessLevelOfUserInProject retrieves the access level of a user in a project.
func (repo *GitlabRepo) GetAccessLevelOfUserInProject(projectId int, userId int) (model.AccessLevelValue, error) {
	repo.assertIsConnected()

	member, _, err := repo.client.ProjectMembers.GetProjectMember(projectId, userId)
	if err != nil {
		return model.NoPermissions, ErrorFromGoGitlab(err)
	}

	return model.AccessLevelValue(member.AccessLevel), nil
}

// GetNamespaceOfProject fetches the namespace of a project.
func (repo *GitlabRepo) GetNamespaceOfProject(projectId int) (*string, error) {
	repo.assertIsConnected()

	project, _, err := repo.client.Projects.GetProject(projectId, &goGitlab.GetProjectOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return &project.Namespace.Path, nil
}

// CreateGroup creates a new GitLab group.
func (repo *GitlabRepo) CreateGroup(name string, visibility model.Visibility, description string) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabVisibility := VisibilityFromModel(visibility)

	path := convertToGitLabPath(strings.ToLower(name))

	createOpts := &goGitlab.CreateGroupOptions{
		Name:        goGitlab.String(name),
		Path:        goGitlab.String(path),
		Description: goGitlab.String(description),
		Visibility:  &gitlabVisibility,
	}

	gitlabGroup, _, err := repo.client.Groups.CreateGroup(createOpts)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupFromGoGitlab(*gitlabGroup), nil
}

// CreateSubGroup creates a subgroup under a parent group.
func (repo *GitlabRepo) CreateSubGroup(name string, path string, parentId int, visibility model.Visibility, description string) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabVisibility := VisibilityFromModel(visibility)

	path = convertToGitLabPath(strings.ToLower(path))

	createOpts := &goGitlab.CreateGroupOptions{
		Name:        goGitlab.String(name),
		Path:        goGitlab.String(path),
		Description: goGitlab.String(description),
		Visibility:  &gitlabVisibility,
		ParentID:    goGitlab.Int(parentId),
	}

	gitlabGroup, _, err := repo.client.Groups.CreateGroup(createOpts)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupFromGoGitlab(*gitlabGroup), nil
}

// CreateGroupAccessToken creates an access token for a GitLab group.
func (repo *GitlabRepo) CreateGroupAccessToken(groupID int, name string, accessLevel model.AccessLevelValue, expiresAt time.Time, scopes ...string) (*model.GroupAccessToken, error) {
	repo.assertIsConnected()

	gitlabExpiresAt := goGitlab.ISOTime(expiresAt)

	accessToken, _, err := repo.client.GroupAccessTokens.CreateGroupAccessToken(groupID, &goGitlab.CreateGroupAccessTokenOptions{
		Name:        goGitlab.String(name),
		Scopes:      &scopes,
		AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
		ExpiresAt:   &gitlabExpiresAt,
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupAccessTokenFromGoGitlabGroupAccessToken(*accessToken), nil
}

// GetGroupAccessToken retrieves a specific group access token.
func (repo *GitlabRepo) GetGroupAccessToken(groupID int, tokenID int) (*model.GroupAccessToken, error) {
	repo.assertIsConnected()

	accessToken, _, err := repo.client.GroupAccessTokens.GetGroupAccessToken(groupID, tokenID)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupAccessTokenFromGoGitlabGroupAccessToken(*accessToken), nil
}

// RotateGroupAccessToken rotates an existing group access token.
func (repo *GitlabRepo) RotateGroupAccessToken(groupID int, tokenID int, expiresAt time.Time) (*model.GroupAccessToken, error) {
	repo.assertIsConnected()

	accessToken, _, err := repo.client.GroupAccessTokens.RotateGroupAccessToken(groupID, tokenID, func(r *retryablehttp.Request) error {
		return r.SetBody([]byte(fmt.Sprintf(`{"expires_at": "%s"}`, expiresAt.Format(time.DateOnly))))
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupAccessTokenFromGoGitlabGroupAccessToken(*accessToken), nil
}

// DeleteProject deletes a project from GitLab.
func (repo *GitlabRepo) DeleteProject(id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Projects.DeleteProject(id)
	return ErrorFromGoGitlab(err)
}

// DeleteGroup deletes a group from GitLab.
func (repo *GitlabRepo) DeleteGroup(id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Groups.DeleteGroup(id)
	return ErrorFromGoGitlab(err)
}

// ChangeGroupName changes the name of a group.
func (repo *GitlabRepo) ChangeGroupName(id int, name string) (*model.Group, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Groups.UpdateGroup(id, &goGitlab.UpdateGroupOptions{
		Name: goGitlab.String(name),
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetGroupById(id)
}

// ChangeGroupDescription changes the description of a group.
func (repo *GitlabRepo) ChangeGroupDescription(id int, description string) (*model.Group, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Groups.UpdateGroup(id, &goGitlab.UpdateGroupOptions{
		Description: goGitlab.String(description),
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetGroupById(id)
}

// ChangeProjectName changes the name of a project.
func (repo *GitlabRepo) ChangeProjectName(projectId int, name string) (*model.Project, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Projects.EditProject(projectId, &goGitlab.EditProjectOptions{
		Name: goGitlab.String(name),
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetProjectById(projectId)
}

// ChangeProjectDescription changes the description of a project.
func (repo *GitlabRepo) ChangeProjectDescription(projectId int, description string) (*model.Project, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Projects.EditProject(projectId, &goGitlab.EditProjectOptions{
		Description: goGitlab.String(description),
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetProjectById(projectId)
}

// GetProjectPipelineTestReportSummary retrieves the test report summary for a specific pipeline
// in a GitLab project.
//
// Parameters:
// - projectId: The ID of the project.
// - pipelineId: The ID of the pipeline.
//
// Returns:
// - *model.TestReport: The test report summary of the pipeline.
// - error: An error if the retrieval fails.
func (repo *GitlabRepo) GetProjectPipelineTestReportSummary(projectId, pipelineId int) (*model.TestReport, error) {
	repo.assertIsConnected()

	testReport, _, err := repo.client.Pipelines.GetPipelineTestReport(projectId, pipelineId)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return TestReportFromGoGitlabTestReport(testReport), nil
}

// GetProjectLatestPipelineTestReportSummary retrieves the test report summary for the latest pipeline
// in a GitLab project, optionally filtering by a reference (branch or tag).
//
// Parameters:
// - projectId: The ID of the project.
// - ref: An optional reference (branch or tag). If nil, the default branch is used.
//
// Returns:
// - *model.TestReport: The test report summary of the latest pipeline.
// - error: An error if the retrieval fails.
func (repo *GitlabRepo) GetProjectLatestPipeline(projectId int, ref *string) (*model.Pipeline, error) {
	repo.assertIsConnected()

	options := &goGitlab.GetLatestPipelineOptions{}
	if ref != nil {
		options.Ref = ref
	}

	pipeline, _, err := repo.client.Pipelines.GetLatestPipeline(projectId, options)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return PipelineFromGoGitlabPipeline(pipeline), nil
}

// GetProjectLatestPipelineTestReportSummary retrieves the test report summary for the latest pipeline
func (repo *GitlabRepo) GetProjectLatestPipelineTestReportSummary(projectId int, ref *string) (*model.TestReport, error) {
	repo.assertIsConnected()

	pipeline, err := repo.GetProjectLatestPipeline(projectId, ref)
	if err != nil {
		return nil, err
	}

	return repo.GetProjectPipelineTestReportSummary(projectId, pipeline.ID)
}

// AddUserToGroup adds a user to a group with the specified access level.
func (repo *GitlabRepo) AddUserToGroup(groupId int, userId int, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	members, _, err := repo.client.Groups.ListGroupMembers(groupId, &goGitlab.ListGroupMembersOptions{})
	if err != nil {
		return err // Handle error appropriately
	}

	// Check if user is already a member
	for _, member := range members {
		if member.ID == userId {
			return nil
		}
	}

	// User is not a member, proceed to add
	_, _, err = repo.client.GroupMembers.AddGroupMember(groupId, &goGitlab.AddGroupMemberOptions{
		UserID:      goGitlab.Int(userId),
		AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
	})

	return ErrorFromGoGitlab(err)
}

// RemoveUserFromGroup removes a user from a group.
func (repo *GitlabRepo) RemoveUserFromGroup(groupId int, userId int) error {
	repo.assertIsConnected()

	_, err := repo.client.GroupMembers.RemoveGroupMember(groupId, userId, nil)

	return ErrorFromGoGitlab(err)
}

// ChangeUserAccessLevelInGroup changes the access level of a user in a group.
func (repo *GitlabRepo) ChangeUserAccessLevelInGroup(groupId int, userId int, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	_, _, err := repo.client.GroupMembers.EditGroupMember(groupId, userId, &goGitlab.EditGroupMemberOptions{
		AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
	})

	return ErrorFromGoGitlab(err)
}

// GetAccessLevelOfUserInGroup retrieves the access level of a user in a group.
func (repo *GitlabRepo) GetAccessLevelOfUserInGroup(groupId int, userId int) (model.AccessLevelValue, error) {
	repo.assertIsConnected()

	member, _, err := repo.client.GroupMembers.GetGroupMember(groupId, userId)
	if err != nil {
		return model.NoPermissions, ErrorFromGoGitlab(err)
	}

	return model.AccessLevelValue(member.AccessLevel), nil
}

// GetAllProjects fetches all projects from GitLab for a given search term.
func (repo *GitlabRepo) GetAllProjects(search string) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Projects.ListProjects(&goGitlab.ListProjectsOptions{
		Archived:   goGitlab.Bool(false),
		Visibility: goGitlab.Visibility(goGitlab.PublicVisibility),
		Owned:      goGitlab.Bool(true),
		OrderBy:    goGitlab.String("created_at"),
		Search:     goGitlab.String(search),
	}, func(r *retryablehttp.Request) error {
		query := r.URL.Query()
		query.Add("per_page", "100")
		r.URL.RawQuery = query.Encode()
		return nil
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

// GetPublicAvatarByMail fetches the public avatar URL for a given email address.
func (repo *GitlabRepo) GetPublicAvatarByMail(mail string) (url *string, err error) {
	repo.assertIsConnected()

	avatar, response, err := repo.client.Avatar.GetAvatar(&goGitlab.GetAvatarOptions{Email: &mail})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("avatar is not available, status code: %d", response.StatusCode)
	}

	return &avatar.AvatarURL, nil
}

// GetProjectById fetches a project by its ID.
func (repo *GitlabRepo) GetProjectById(id int) (*model.Project, error) {
	repo.assertIsConnected()

	gitlabProject, _, err := repo.client.Projects.GetProject(id, &goGitlab.GetProjectOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProject(gitlabProject)
}

// GetUserById fetches a user by their ID.
func (repo *GitlabRepo) GetUserById(id int) (*model.User, error) {
	repo.assertIsConnected()

	gitlabUser, _, err := repo.client.Users.GetUser(id, goGitlab.GetUsersOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return UserFromGoGitlab(*gitlabUser), nil
}

// GetGroupById fetches a group by its ID.
func (repo *GitlabRepo) GetGroupById(id int) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroup, _, err := repo.client.Groups.GetGroup(id, &goGitlab.GetGroupOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabGroup(gitlabGroup)
}

// GetAllUsers fetches all users from GitLab.
func (repo *GitlabRepo) GetAllUsers() ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Users.ListUsers(&goGitlab.ListUsersOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

// GetAllGroups fetches all groups from GitLab.
func (repo *GitlabRepo) GetAllGroups() ([]*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroups, _, err := repo.client.Groups.ListGroups(&goGitlab.ListGroupsOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabGroups(gitlabGroups)
}

// GetAllProjectsOfGroup fetches all projects of a group by its ID.
func (repo *GitlabRepo) GetAllProjectsOfGroup(id int) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Groups.ListGroupProjects(id, &goGitlab.ListGroupProjectsOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

// GetAllUsersOfGroup fetches all users of a group by its ID.
func (repo *GitlabRepo) GetAllUsersOfGroup(id int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabMembers, _, err := repo.client.Groups.ListGroupMembers(id, &goGitlab.ListGroupMembersOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	users := make([]*model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		users[i] = UserFromGoGitlabGroupMember(*gitlabMember)
	}

	return users, nil
}

// SearchProjectByExpression searches for projects by a given expression.
func (repo *GitlabRepo) SearchProjectByExpression(expression string) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Search.Projects(expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

// SearchUserByExpression searches for users by a given expression.
func (repo *GitlabRepo) SearchUserByExpression(expression string) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.Users(expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

// SearchUserByExpressionInGroup searches for users by a given expression in a group.
func (repo *GitlabRepo) SearchUserByExpressionInGroup(expression string, groupId int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.UsersByGroup(groupId, expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

// SearchUserByExpressionInProject searches for users by a given expression in a project.
func (repo *GitlabRepo) SearchUserByExpressionInProject(expression string, projectId int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.UsersByProject(projectId, expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

// SearchGroupByExpression searches for groups by a given expression.
func (repo *GitlabRepo) SearchGroupByExpression(expression string) ([]*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroups, _, err := repo.client.Groups.SearchGroup(expression)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	groups := make([]*model.Group, len(gitlabGroups))
	for i, gitlabGroup := range gitlabGroups {
		groups[i] = GroupFromGoGitlab(*gitlabGroup)
	}

	return groups, nil
}

// GetPendingProjectInvitations fetches all pending project invitations for a project.
func (repo *GitlabRepo) GetPendingProjectInvitations(projectId int) ([]*model.PendingInvite, error) {
	repo.assertIsConnected()

	pendingInvites, _, err := repo.client.Invites.ListPendingProjectInvitations(projectId, nil)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabPendingInvites(pendingInvites)
}

// GetPendingGroupInvitations fetches all pending group invitations for a group.
func (repo *GitlabRepo) GetPendingGroupInvitations(groupId int) ([]*model.PendingInvite, error) {
	repo.assertIsConnected()

	pendingInvites, _, err := repo.client.Invites.ListPendingGroupInvitations(groupId, nil)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabPendingInvites(pendingInvites)
}

// CreateGroupInvite creates an invitation for a user to join a group.
func (repo *GitlabRepo) CreateGroupInvite(groupId int, email string) error {
	repo.assertIsConnected()

	_, _, err := repo.client.Invites.GroupInvites(groupId, &goGitlab.InvitesOptions{
		Email:       &email,
		AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		// Set additional options like AccessLevel, ExpiresAt as needed
	})
	return ErrorFromGoGitlab(err)
}

// CreateProjectInvite creates an invitation for a user to join a project.
func (repo *GitlabRepo) CreateProjectInvite(id int, email string) error {
	repo.assertIsConnected()

	_, _, err := repo.client.Invites.ProjectInvites(id, &goGitlab.InvitesOptions{
		Email:       &email,
		AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		// Set additional options like AccessLevel, ExpiresAt as needed
	})
	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) DenyPushingToProject(projectId int) error {
	log.Panic("No working option to close an assignment")

	permission := goGitlab.MinimalAccessPermissions

	return repo.changeProjectMemberPermissions(projectId, permission)
}

func (repo *GitlabRepo) AllowPushingToProject(projectId int) error {
	log.Panic("No working option to reopen an assignment")

	permission := goGitlab.DeveloperPermissions

	return repo.changeProjectMemberPermissions(projectId, permission)
}

func (repo *GitlabRepo) changeProjectMemberPermissions(projectId int, accessLevel goGitlab.AccessLevelValue) error {
	repo.assertIsConnected()

	members, _, err := repo.client.ProjectMembers.ListAllProjectMembers(projectId, &goGitlab.ListProjectMembersOptions{})
	if err != nil {
		return ErrorFromGoGitlab(err)
	}

	for _, member := range members {
		if member.AccessLevel == *goGitlab.AccessLevel(goGitlab.OwnerPermissions) {
			continue
		}

		_, _, err := repo.client.ProjectMembers.EditProjectMember(projectId, member.ID, &goGitlab.EditProjectMemberOptions{AccessLevel: &accessLevel})
		if err != nil {
			return ErrorFromGoGitlab(err)
		}
	}

	return nil
}

// GetAvailableRunnersForGitLab fetches all available runners for GitLab.
func (repo *GitlabRepo) GetAvailableRunnersForGitLab() ([]*model.Runner, error) {
	repo.assertIsConnected()

	availableRunners, _, err := repo.client.Runners.ListRunners(
		&goGitlab.ListRunnersOptions{
			Status: goGitlab.String("online"), Paused: goGitlab.Bool(false),
			Type: goGitlab.String("instance_type")})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return convertRunners(availableRunners), nil
}

// GetAvailableRunnersForGroup fetches all available runners for a group.
func (repo *GitlabRepo) GetAvailableRunnersForGroup(groupId int) ([]*model.Runner, error) {
	repo.assertIsConnected()

	runners, _, err := repo.client.Runners.ListGroupsRunners(groupId,
		&goGitlab.ListGroupsRunnersOptions{Status: goGitlab.String("online")})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	var availableRunners []*goGitlab.Runner
	for _, runner := range runners {
		if !runner.Paused {
			availableRunners = append(availableRunners, runner)
		}
	}

	return convertRunners(availableRunners), nil
}

// CheckIfFileExistsInProject checks if a file exists in a project.
func (repo *GitlabRepo) CheckIfFileExistsInProject(projectId int, filepath string) (bool, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.RepositoryFiles.GetFile(projectId, filepath, &goGitlab.GetFileOptions{Ref: goGitlab.String("HEAD")})
	if err != nil {
		var gitlabErr *goGitlab.ErrorResponse
		if errors.As(err, &gitlabErr) &&
			gitlabErr.Response.StatusCode == 404 &&
			strings.Contains(gitlabErr.Message, "404 File Not Found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetProjectLanguages retrieves the languages used in a project.
func (repo *GitlabRepo) GetProjectLanguages(projectId int) (map[string]float32, error) {
	repo.assertIsConnected()

	languages, _, err := repo.client.Projects.GetProjectLanguages(projectId)
	if err != nil {
		return nil, err
	}

	return *languages, nil
}

func (repo *GitlabRepo) assertIsConnected() {
	if repo.client == nil {
		panic("No connection to Gitlab! Make sure you have executed Login()")
	}
}

func (repo *GitlabRepo) getUserByUsername(username string) (*model.User, error) {
	repo.assertIsConnected()

	users, _, err := repo.client.Users.ListUsers(&goGitlab.ListUsersOptions{
		Search: goGitlab.String(username),
	})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with username [%s] not found", username)
	}

	return UserFromGoGitlab(*users[0]), nil
}

// FindUserIDByEmail finds a user ID by their email address.
func (repo *GitlabRepo) FindUserIDByEmail(email string) (int, error) {
	repo.assertIsConnected()

	listUsersOptions := &goGitlab.ListUsersOptions{
		Search: goGitlab.String(email),
	}

	users, _, err := repo.client.Users.ListUsers(listUsersOptions)
	if err != nil {
		return 0, ErrorFromGoGitlab(err)
	}

	if len(users) != 1 {
		return 0, fmt.Errorf("user not found or multiple users found with email: %s", email)
	}

	return users[0].ID, nil
}

func (repo *GitlabRepo) convertGitlabUsers(gitlabUsers []*goGitlab.User) ([]*model.User, error) {
	users := make([]*model.User, len(gitlabUsers))
	for i, gitlabUser := range gitlabUsers {
		users[i] = UserFromGoGitlab(*gitlabUser)
	}

	return users, nil
}

func (repo *GitlabRepo) convertGitlabProjects(gitlabProjects []*goGitlab.Project) ([]*model.Project, error) {
	projects := make([]*model.Project, len(gitlabProjects))
	for i, gitlabProject := range gitlabProjects {
		project, err := repo.convertGitlabProject(gitlabProject)
		if err != nil {
			return nil, ErrorFromGoGitlab(err)
		}

		projects[i] = project
	}

	return projects, nil
}

func (repo *GitlabRepo) convertGitlabProject(gitlabProject *goGitlab.Project) (*model.Project, error) {
	gitlabMembers, _, err := repo.client.ProjectMembers.ListProjectMembers(gitlabProject.ID, &goGitlab.ListProjectMembersOptions{})
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return ProjectFromGoGitlabWithProjectMembers(*gitlabProject, gitlabMembers), nil
}

func (repo *GitlabRepo) convertGitlabGroup(gitlabGroup *goGitlab.Group) (*model.Group, error) {
	Group := GroupFromGoGitlab(*gitlabGroup)

	projects, err := repo.GetAllProjectsOfGroup(gitlabGroup.ID)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}
	Group.Projects = ConvertProjectPointerSlice(projects)

	members, err := repo.GetAllUsersOfGroup(gitlabGroup.ID)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}
	Group.Member = ConvertUserPointerSlice(members)

	return Group, nil
}

func (repo *GitlabRepo) convertGitlabGroups(gitlabGroups []*goGitlab.Group) ([]*model.Group, error) {
	groups := make([]*model.Group, len(gitlabGroups))
	for i, gitlabGroup := range gitlabGroups {
		group, err := repo.convertGitlabGroup(gitlabGroup)
		if err != nil {
			return nil, ErrorFromGoGitlab(err)
		}

		groups[i] = group
	}

	return groups, nil
}

func (repo *GitlabRepo) convertGitlabPendingInvites(gitlabPendingInvites []*goGitlab.PendingInvite) ([]*model.PendingInvite, error) {
	pendingInvites := make([]*model.PendingInvite, len(gitlabPendingInvites))
	for i, gitlabPendingInvite := range gitlabPendingInvites {
		pendingInvites[i] = PendingInviteFromGoGitlab(*gitlabPendingInvite)
	}

	return pendingInvites, nil
}

func convertToGitLabPath(s string) string {
	// Remove unwanted characters
	reg := regexp.MustCompile("[^a-zA-Z0-9_.-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove leading and trailing special characters
	s = strings.Trim(s, "_.-")

	// Prevent consecutive special characters
	reg = regexp.MustCompile("[-_.]{2,}")
	s = reg.ReplaceAllString(s, "-")

	// Prevent specific endings
	if strings.HasSuffix(s, ".git") || strings.HasSuffix(s, ".atom") {
		s = s[:len(s)-4]
	}

	// Ensure the path name is at least one character long
	if len(s) == 0 {
		s = "gc_"
	}

	return s
}

func convertRunners(runners []*goGitlab.Runner) []*model.Runner {
	convertedRunners := make([]*model.Runner, len(runners))
	for i, r := range runners {
		if r != nil {
			rConverted := model.Runner(*r)
			convertedRunners[i] = &rConverted
		}
	}
	return convertedRunners
}
