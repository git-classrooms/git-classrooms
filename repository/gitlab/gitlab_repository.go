package gitlab

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"

	goGitlab "github.com/xanzy/go-gitlab"
)

type GitlabRepo struct {
	client      *goGitlab.Client
	config      gitlabConfig.Config
	isConnected bool
}

func NewGitlabRepo(config gitlabConfig.Config) *GitlabRepo {
	return &GitlabRepo{client: nil, config: config, isConnected: false}
}

// Reference to Go Gitlab Documentation: https://pkg.go.dev/github.com/xanzy/go-gitlab#section-documentation

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

func (repo *GitlabRepo) GroupAccessLogin(token string) error {
	cli, err := goGitlab.NewClient(token, goGitlab.WithBaseURL(repo.config.GetURL()))
	if err != nil {
		return err
	}
	repo.client = cli
	return nil
}

func (repo *GitlabRepo) GetCurrentUser(ctx context.Context) (*model.User, error) {
	repo.assertIsConnected()

	gitlabUser, _, err := repo.client.Users.CurrentUser(goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	classroomUser := UserFromGoGitlab(*gitlabUser)
	classroomUser.Avatar.FallbackAvatarURL, _ = repo.GetPublicAvatarByMail(ctx, classroomUser.Email)
	return classroomUser, nil
}

func (repo *GitlabRepo) CreateProject(ctx context.Context, name string, visibility model.Visibility, description string, members []model.User) (*model.Project, error) {
	repo.assertIsConnected()

	opts := &goGitlab.CreateProjectOptions{
		Name:        goGitlab.String(name),
		Visibility:  goGitlab.Visibility(VisibilityFromModel(visibility)),
		Description: goGitlab.String(description),
	}

	gitlabProject, _, err := repo.client.Projects.CreateProject(opts, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.AddProjectMembers(ctx, gitlabProject.ID, members)
}

func (repo *GitlabRepo) ForkProject(ctx context.Context, projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error) {
	repo.assertIsConnected()

	opts := &goGitlab.ForkProjectOptions{
		Name:                          goGitlab.String(name),
		Path:                          goGitlab.String(convertToGitLabPath(name)),
		NamespaceID:                   goGitlab.Int(namespaceId),
		Visibility:                    goGitlab.Visibility(VisibilityFromModel(visibility)),
		Description:                   goGitlab.String(description),
		MergeRequestDefaultTargetSelf: goGitlab.Bool(true),
	}

	gitlabProject, _, err := repo.client.Projects.ForkProject(projectId, opts, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProject(gitlabProject)
}

func (repo *GitlabRepo) ForkProjectWithOnlyDefaultBranch(ctx context.Context, projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error) {
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
	}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProject(gitlabProject)
}

func (repo *GitlabRepo) CreateBranch(ctx context.Context, projectId int, branchName string, fromBranch string) (*model.Branch, error) {
	repo.assertIsConnected()

	opts := &goGitlab.CreateBranchOptions{
		Branch: goGitlab.String(branchName),
		Ref:    goGitlab.String(fromBranch),
	}

	branch, _, err := repo.client.Branches.CreateBranch(projectId, opts, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}
	return BranchFromGoGitlab(branch), nil
}

func (repo *GitlabRepo) ProtectBranch(ctx context.Context, projectId int, branchName string, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	opts := &goGitlab.ProtectRepositoryBranchesOptions{
		Name:             goGitlab.String(branchName),
		AllowForcePush:   goGitlab.Bool(false),
		PushAccessLevel:  goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
		MergeAccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
	}

	_, _, err := repo.client.ProtectedBranches.ProtectRepositoryBranches(projectId, opts, goGitlab.WithContext(ctx))
	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) UnprotectBranch(ctx context.Context, projectId int, branchName string) error {
	repo.assertIsConnected()

	_, err := repo.client.ProtectedBranches.UnprotectRepositoryBranches(projectId, branchName, goGitlab.WithContext(ctx))
	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) CreateMergeRequest(ctx context.Context, projectId int, sourceBranch string, targetBranch string, title string, description string, assigneeId int, reviewerId int) error {
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

	_, _, err := repo.client.MergeRequests.CreateMergeRequest(projectId, opts, goGitlab.WithContext(ctx))
	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) ProtectedBranchExists(ctx context.Context, projectId int, branchName string) (bool, error) {
	repo.assertIsConnected()

	_, response, err := repo.client.ProtectedBranches.GetProtectedBranch(projectId, branchName, goGitlab.WithContext(ctx))
	if err != nil {
		if response.StatusCode == 404 {
			return false, nil
		}
		return false, ErrorFromGoGitlab(err)
	}

	return true, nil
}

func (repo *GitlabRepo) BranchExists(ctx context.Context, projectId int, branchName string) (bool, error) {
	repo.assertIsConnected()

	_, response, err := repo.client.Branches.GetBranch(projectId, branchName, goGitlab.WithContext(ctx))
	if err != nil {
		if response.StatusCode == 404 {
			return false, nil
		}
		return false, ErrorFromGoGitlab(err)
	}

	return true, nil
}

func (repo *GitlabRepo) AddProjectMembers(ctx context.Context, projectId int, members []model.User) (*model.Project, error) {
	repo.assertIsConnected()

	for _, member := range members {
		opts := &goGitlab.AddProjectMemberOptions{
			UserID:      &member.ID,
			AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		}
		_, _, err := repo.client.ProjectMembers.AddProjectMember(projectId, opts, goGitlab.WithContext(ctx))
		if err != nil {
			return nil, ErrorFromGoGitlab(err)
		}
	}

	return repo.GetProjectById(ctx, projectId)
}

func (repo *GitlabRepo) ChangeUserAccessLevelInProject(ctx context.Context, projectId int, userId int, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	_, _, err := repo.client.ProjectMembers.EditProjectMember(
		projectId,
		userId,
		&goGitlab.EditProjectMemberOptions{AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel))},
		goGitlab.WithContext(ctx),
	)

	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) GetAccessLevelOfUserInProject(ctx context.Context, projectId int, userId int) (model.AccessLevelValue, error) {
	repo.assertIsConnected()

	member, _, err := repo.client.ProjectMembers.GetProjectMember(projectId, userId, goGitlab.WithContext(ctx))
	if err != nil {
		return model.NoPermissions, ErrorFromGoGitlab(err)
	}

	return model.AccessLevelValue(member.AccessLevel), nil
}

func (repo *GitlabRepo) GetNamespaceOfProject(ctx context.Context, projectId int) (*string, error) {
	repo.assertIsConnected()

	project, _, err := repo.client.Projects.GetProject(projectId, &goGitlab.GetProjectOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return &project.Namespace.Path, nil
}

func (repo *GitlabRepo) CreateGroup(ctx context.Context, name string, visibility model.Visibility, description string) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabVisibility := VisibilityFromModel(visibility)

	path := convertToGitLabPath(strings.ToLower(name))

	createOpts := &goGitlab.CreateGroupOptions{
		Name:        goGitlab.String(name),
		Path:        goGitlab.String(path),
		Description: goGitlab.String(description),
		Visibility:  &gitlabVisibility,
	}

	gitlabGroup, _, err := repo.client.Groups.CreateGroup(createOpts, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupFromGoGitlab(*gitlabGroup), nil
}

func (repo *GitlabRepo) CreateSubGroup(ctx context.Context, name string, parentId int, visibility model.Visibility, description string) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabVisibility := VisibilityFromModel(visibility)

	path := convertToGitLabPath(strings.ToLower(name))

	createOpts := &goGitlab.CreateGroupOptions{
		Name:        goGitlab.String(name),
		Path:        goGitlab.String(path),
		Description: goGitlab.String(description),
		Visibility:  &gitlabVisibility,
		ParentID:    goGitlab.Int(parentId),
	}

	gitlabGroup, _, err := repo.client.Groups.CreateGroup(createOpts, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupFromGoGitlab(*gitlabGroup), nil
}

func (repo *GitlabRepo) CreateGroupAccessToken(ctx context.Context, groupID int, name string, accessLevel model.AccessLevelValue, expiresAt time.Time, scopes ...string) (*model.GroupAccessToken, error) {
	repo.assertIsConnected()

	gitlabExpiresAt := goGitlab.ISOTime(expiresAt)
	opts := &goGitlab.CreateGroupAccessTokenOptions{
		Name:        goGitlab.String(name),
		Scopes:      &scopes,
		AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
		ExpiresAt:   &gitlabExpiresAt,
	}

	accessToken, _, err := repo.client.GroupAccessTokens.CreateGroupAccessToken(groupID, opts, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupAccessTokenFromGoGitlabGroupAccessToken(*accessToken), nil
}

func (repo *GitlabRepo) RotateGroupAccessToken(ctx context.Context, groupID int, tokenID int, expiresAt time.Time) (*model.GroupAccessToken, error) {
	repo.assertIsConnected()

	accessToken, _, err := repo.client.GroupAccessTokens.RotateGroupAccessToken(groupID, tokenID, func(r *retryablehttp.Request) error {
		return r.SetBody([]byte(fmt.Sprintf(`{"expires_at": "%s"}`, expiresAt.Format(time.DateOnly))))
	}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return GroupAccessTokenFromGoGitlabGroupAccessToken(*accessToken), nil
}

func (repo *GitlabRepo) DeleteProject(ctx context.Context, id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Projects.DeleteProject(id, goGitlab.WithContext(ctx))
	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) DeleteGroup(ctx context.Context, id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Groups.DeleteGroup(id, goGitlab.WithContext(ctx))
	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) ChangeGroupName(ctx context.Context, id int, name string) (*model.Group, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Groups.UpdateGroup(id, &goGitlab.UpdateGroupOptions{Name: goGitlab.String(name)}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetGroupById(ctx, id)
}

func (repo *GitlabRepo) ChangeGroupDescription(ctx context.Context, id int, description string) (*model.Group, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Groups.UpdateGroup(id, &goGitlab.UpdateGroupOptions{Description: goGitlab.String(description)}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetGroupById(ctx, id)
}

func (repo *GitlabRepo) ChangeProjectName(ctx context.Context, projectId int, name string) (*model.Project, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Projects.EditProject(projectId, &goGitlab.EditProjectOptions{Name: goGitlab.String(name)}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetProjectById(ctx, projectId)
}

func (repo *GitlabRepo) ChangeProjectDescription(ctx context.Context, projectId int, description string) (*model.Project, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Projects.EditProject(projectId, &goGitlab.EditProjectOptions{Description: goGitlab.String(description)}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.GetProjectById(ctx, projectId)
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
func (repo *GitlabRepo) GetProjectPipelineTestReportSummary(ctx context.Context, projectId, pipelineId int) (*model.TestReport, error) {
	repo.assertIsConnected()

	testReport, _, err := repo.client.Pipelines.GetPipelineTestReport(projectId, pipelineId, goGitlab.WithContext(ctx))
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
func (repo *GitlabRepo) GetProjectLatestPipeline(ctx context.Context, projectId int, ref *string) (*model.Pipeline, error) {
	repo.assertIsConnected()

	options := &goGitlab.GetLatestPipelineOptions{}
	if ref != nil {
		options.Ref = ref
	}

	pipeline, _, err := repo.client.Pipelines.GetLatestPipeline(projectId, options, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return PipelineFromGoGitlabPipeline(pipeline), nil
}

func (repo *GitlabRepo) GetProjectLatestPipelineTestReportSummary(ctx context.Context, projectId int, ref *string) (*model.TestReport, error) {
	repo.assertIsConnected()

	pipeline, err := repo.GetProjectLatestPipeline(ctx, projectId, ref)
	if err != nil {
		return nil, err
	}

	return repo.GetProjectPipelineTestReportSummary(ctx, projectId, pipeline.ID)
}

func (repo *GitlabRepo) AddUserToGroup(ctx context.Context, groupId int, userId int, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	members, _, err := repo.client.Groups.ListGroupMembers(groupId, &goGitlab.ListGroupMembersOptions{}, goGitlab.WithContext(ctx))
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
	}, goGitlab.WithContext(ctx))

	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) RemoveUserFromGroup(ctx context.Context, groupId int, userId int) error {
	repo.assertIsConnected()

	_, err := repo.client.GroupMembers.RemoveGroupMember(groupId, userId, nil, goGitlab.WithContext(ctx))

	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) ChangeUserAccessLevelInGroup(ctx context.Context, groupId int, userId int, accessLevel model.AccessLevelValue) error {
	repo.assertIsConnected()

	_, _, err := repo.client.GroupMembers.EditGroupMember(groupId, userId, &goGitlab.EditGroupMemberOptions{
		AccessLevel: goGitlab.AccessLevel(AccessLevelFromModel(accessLevel)),
	}, goGitlab.WithContext(ctx))

	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) GetAccessLevelOfUserInGroup(ctx context.Context, groupId int, userId int) (model.AccessLevelValue, error) {
	repo.assertIsConnected()

	member, _, err := repo.client.GroupMembers.GetGroupMember(groupId, userId, goGitlab.WithContext(ctx))
	if err != nil {
		return model.NoPermissions, ErrorFromGoGitlab(err)
	}

	return model.AccessLevelValue(member.AccessLevel), nil
}

func (repo *GitlabRepo) GetAllProjects(ctx context.Context, search string) ([]*model.Project, error) {
	repo.assertIsConnected()

	opts := &goGitlab.ListProjectsOptions{
		Archived:   goGitlab.Bool(false),
		Visibility: goGitlab.Visibility(goGitlab.PublicVisibility),
		Owned:      goGitlab.Bool(true),
		OrderBy:    goGitlab.String("created_at"),
		Search:     goGitlab.String(search),
	}
	gitlabProjects, _, err := repo.client.Projects.ListProjects(opts, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GitlabRepo) GetPublicAvatarByMail(ctx context.Context, mail string) (url *string, err error) {
	repo.assertIsConnected()

	avatar, response, err := repo.client.Avatar.GetAvatar(&goGitlab.GetAvatarOptions{Email: &mail}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("avatar is not available, status code: %d", response.StatusCode)
	}

	return &avatar.AvatarURL, nil
}

func (repo *GitlabRepo) GetProjectById(ctx context.Context, id int) (*model.Project, error) {
	repo.assertIsConnected()

	gitlabProject, _, err := repo.client.Projects.GetProject(id, &goGitlab.GetProjectOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProject(gitlabProject)
}

func (repo *GitlabRepo) GetUserById(ctx context.Context, id int) (*model.User, error) {
	repo.assertIsConnected()

	gitlabUser, _, err := repo.client.Users.GetUser(id, goGitlab.GetUsersOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return UserFromGoGitlab(*gitlabUser), nil
}

func (repo *GitlabRepo) GetGroupById(ctx context.Context, id int) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroup, _, err := repo.client.Groups.GetGroup(id, &goGitlab.GetGroupOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabGroup(ctx, gitlabGroup)
}

func (repo *GitlabRepo) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Users.ListUsers(&goGitlab.ListUsersOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) GetAllGroups(ctx context.Context) ([]*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroups, _, err := repo.client.Groups.ListGroups(&goGitlab.ListGroupsOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabGroups(ctx, gitlabGroups)
}

func (repo *GitlabRepo) GetAllProjectsOfGroup(ctx context.Context, id int) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Groups.ListGroupProjects(id, &goGitlab.ListGroupProjectsOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GitlabRepo) GetAllUsersOfGroup(ctx context.Context, id int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabMembers, _, err := repo.client.Groups.ListGroupMembers(id, &goGitlab.ListGroupMembersOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	users := make([]*model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		users[i] = UserFromGoGitlabGroupMember(*gitlabMember)
	}

	return users, nil
}

func (repo *GitlabRepo) SearchProjectByExpression(ctx context.Context, expression string) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Search.Projects(expression, &goGitlab.SearchOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GitlabRepo) SearchUserByExpression(ctx context.Context, expression string) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.Users(expression, &goGitlab.SearchOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) SearchUserByExpressionInGroup(ctx context.Context, expression string, groupId int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.UsersByGroup(groupId, expression, &goGitlab.SearchOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) SearchUserByExpressionInProject(ctx context.Context, expression string, projectId int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.UsersByProject(projectId, expression, &goGitlab.SearchOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) SearchGroupByExpression(ctx context.Context, expression string) ([]*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroups, _, err := repo.client.Groups.SearchGroup(expression, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	groups := make([]*model.Group, len(gitlabGroups))
	for i, gitlabGroup := range gitlabGroups {
		groups[i] = GroupFromGoGitlab(*gitlabGroup)
	}

	return groups, nil
}

func (repo *GitlabRepo) GetPendingProjectInvitations(ctx context.Context, projectId int) ([]*model.PendingInvite, error) {
	repo.assertIsConnected()

	pendingInvites, _, err := repo.client.Invites.ListPendingProjectInvitations(projectId, nil, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabPendingInvites(pendingInvites)
}

func (repo *GitlabRepo) GetPendingGroupInvitations(ctx context.Context, groupId int) ([]*model.PendingInvite, error) {
	repo.assertIsConnected()

	pendingInvites, _, err := repo.client.Invites.ListPendingGroupInvitations(groupId, nil, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return repo.convertGitlabPendingInvites(pendingInvites)
}

func (repo *GitlabRepo) CreateGroupInvite(ctx context.Context, groupId int, email string) error {
	repo.assertIsConnected()

	_, _, err := repo.client.Invites.GroupInvites(groupId, &goGitlab.InvitesOptions{
		Email:       &email,
		AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		// Set additional options like AccessLevel, ExpiresAt as needed
	}, goGitlab.WithContext(ctx))
	return ErrorFromGoGitlab(err)
}

func (repo *GitlabRepo) CreateProjectInvite(ctx context.Context, id int, email string) error {
	repo.assertIsConnected()

	_, _, err := repo.client.Invites.ProjectInvites(id, &goGitlab.InvitesOptions{
		Email:       &email,
		AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		// Set additional options like AccessLevel, ExpiresAt as needed
	}, goGitlab.WithContext(ctx))
	return ErrorFromGoGitlab(err)
}

/*
TODO:
	Mit personal access tokens ist es bisher nicht möglich ein Assignment zu schließen bzw. das Pushen zu unterbinden (man bekommt bei alle aufgelisteten Möglichkeiten einen 404 zurück)
	- Not with Push Rules
	- Not with Protect Branches
	- Not with change Project Member Access Level
*/

func (repo *GitlabRepo) DenyPushingToProject(ctx context.Context, projectId int) error {
	log.Panic("No working option to close an assignment")

	permission := goGitlab.MinimalAccessPermissions

	return repo.changeProjectMemberPermissions(ctx, projectId, permission)
}

func (repo *GitlabRepo) AllowPushingToProject(ctx context.Context, projectId int) error {
	log.Panic("No working option to reopen an assignment")

	permission := goGitlab.DeveloperPermissions

	return repo.changeProjectMemberPermissions(ctx, projectId, permission)
}

func (repo *GitlabRepo) changeProjectMemberPermissions(ctx context.Context, projectId int, accessLevel goGitlab.AccessLevelValue) error {
	repo.assertIsConnected()

	members, _, err := repo.client.ProjectMembers.ListAllProjectMembers(projectId, &goGitlab.ListProjectMembersOptions{}, goGitlab.WithContext(ctx))
	if err != nil {
		return ErrorFromGoGitlab(err)
	}

	for _, member := range members {
		if member.AccessLevel == *goGitlab.AccessLevel(goGitlab.OwnerPermissions) {
			continue
		}

		_, _, err := repo.client.ProjectMembers.EditProjectMember(projectId, member.ID, &goGitlab.EditProjectMemberOptions{AccessLevel: &accessLevel}, goGitlab.WithContext(ctx))
		if err != nil {
			return ErrorFromGoGitlab(err)
		}
	}

	return nil
}

func (repo *GitlabRepo) GetAvailableRunnersForGitLab(ctx context.Context) ([]*model.Runner, error) {
	repo.assertIsConnected()

	availableRunners, _, err := repo.client.Runners.ListRunners(
		&goGitlab.ListRunnersOptions{
			Status: goGitlab.String("online"), Paused: goGitlab.Bool(false),
			Type: goGitlab.String("instance_type")}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	return convertRunners(availableRunners), nil
}

func (repo *GitlabRepo) GetAvailableRunnersForGroup(ctx context.Context, groupId int) ([]*model.Runner, error) {
	repo.assertIsConnected()

	runners, _, err := repo.client.Runners.ListGroupsRunners(groupId,
		&goGitlab.ListGroupsRunnersOptions{Status: goGitlab.String("online")}, goGitlab.WithContext(ctx))
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

func (repo *GitlabRepo) CheckIfFileExistsInProject(ctx context.Context, projectId int, filepath string) (bool, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.RepositoryFiles.GetFile(projectId, filepath, &goGitlab.GetFileOptions{Ref: goGitlab.String("HEAD")}, goGitlab.WithContext(ctx))
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

func (repo *GitlabRepo) GetProjectLanguages(ctx context.Context, projectId int) (map[string]float32, error) {
	repo.assertIsConnected()

	languages, _, err := repo.client.Projects.GetProjectLanguages(projectId, goGitlab.WithContext(ctx))
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

func (repo *GitlabRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	repo.assertIsConnected()

	users, _, err := repo.client.Users.ListUsers(&goGitlab.ListUsersOptions{
		Search: goGitlab.String(username),
	}, goGitlab.WithContext(ctx))
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with username [%s] not found", username)
	}

	return UserFromGoGitlab(*users[0]), nil
}

func (repo *GitlabRepo) FindUserIDByEmail(ctx context.Context, email string) (int, error) {
	repo.assertIsConnected()

	listUsersOptions := &goGitlab.ListUsersOptions{
		Search: goGitlab.String(email),
	}

	users, _, err := repo.client.Users.ListUsers(listUsersOptions, goGitlab.WithContext(ctx))
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

func (repo *GitlabRepo) convertGitlabGroup(ctx context.Context, gitlabGroup *goGitlab.Group) (*model.Group, error) {
	Group := GroupFromGoGitlab(*gitlabGroup)

	projects, err := repo.GetAllProjectsOfGroup(ctx, gitlabGroup.ID)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}
	Group.Projects = ConvertProjectPointerSlice(projects)

	members, err := repo.GetAllUsersOfGroup(ctx, gitlabGroup.ID)
	if err != nil {
		return nil, ErrorFromGoGitlab(err)
	}
	Group.Member = ConvertUserPointerSlice(members)

	return Group, nil
}

func (repo *GitlabRepo) convertGitlabGroups(ctx context.Context, gitlabGroups []*goGitlab.Group) ([]*model.Group, error) {
	groups := make([]*model.Group, len(gitlabGroups))
	for i, gitlabGroup := range gitlabGroups {
		group, err := repo.convertGitlabGroup(ctx, gitlabGroup)
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
	reg, _ := regexp.Compile("[^a-zA-Z0-9_.-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove leading and trailing special characters
	s = strings.Trim(s, "_.-")

	// Prevent consecutive special characters
	reg, _ = regexp.Compile("[-_.]{2,}")
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
