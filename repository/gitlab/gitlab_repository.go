package gitlab

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"image"
	"io"
	"log"
	"net/http"
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

func (repo *GitlabRepo) GetCurrentUser() (*model.User, error) {
	repo.assertIsConnected()

	gitlabUser, _, err := repo.client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}

	classroomUser := UserFromGoGitlab(*gitlabUser)
	classroomUser.AvatarURL = repo.getValidatedUserAvatarURL(gitlabUser)
	return classroomUser, nil
}

func (repo *GitlabRepo) CreateProject(name string, visibility model.Visibility, description string, members []model.User) (*model.Project, error) {
	repo.assertIsConnected()

	opts := &goGitlab.CreateProjectOptions{
		Name:        goGitlab.String(name),
		Visibility:  goGitlab.Visibility(VisibilityFromModel(visibility)),
		Description: goGitlab.String(description),
	}

	gitlabProject, _, err := repo.client.Projects.CreateProject(opts)
	if err != nil {
		return nil, err
	}

	return repo.AddProjectMembers(gitlabProject.ID, members)
}

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
		return nil, err
	}

	return repo.convertGitlabProject(gitlabProject)
}

func (repo *GitlabRepo) AddProjectMembers(projectId int, members []model.User) (*model.Project, error) {
	repo.assertIsConnected()

	for _, member := range members {
		_, _, err := repo.client.ProjectMembers.AddProjectMember(projectId, &goGitlab.AddProjectMemberOptions{
			UserID:      &member.ID,
			AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		})
		if err != nil {
			return nil, err
		}
	}

	return repo.GetProjectById(projectId)
}

func (repo *GitlabRepo) GetNamespaceOfProject(projectId int) (*string, error) {
	repo.assertIsConnected()

	project, _, err := repo.client.Projects.GetProject(projectId, &goGitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	return &project.Namespace.Path, nil
}

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
		return nil, err
	}

	return GroupFromGoGitlab(*gitlabGroup), nil
}

func (repo *GitlabRepo) CreateSubGroup(name string, parentId int, visibility model.Visibility, description string) (*model.Group, error) {
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

	gitlabGroup, _, err := repo.client.Groups.CreateGroup(createOpts)
	if err != nil {
		return nil, err
	}

	return GroupFromGoGitlab(*gitlabGroup), nil
}

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
		return nil, err
	}

	return GroupAccessTokenFromGoGitlabGroupAccessToken(*accessToken), nil
}

func (repo *GitlabRepo) RotateGroupAccessToken(groupID int, tokenID int, expiresAt time.Time) (*model.GroupAccessToken, error) {
	repo.assertIsConnected()

	accessToken, _, err := repo.client.GroupAccessTokens.RotateGroupAccessToken(groupID, tokenID, func(r *retryablehttp.Request) error {
		return r.SetBody([]byte(fmt.Sprintf(`{"expires_at": "%s"}`, expiresAt.Format(time.DateOnly))))
	})
	if err != nil {
		return nil, err
	}

	return GroupAccessTokenFromGoGitlabGroupAccessToken(*accessToken), nil
}

func (repo *GitlabRepo) DeleteProject(id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Projects.DeleteProject(id)
	return err
}

func (repo *GitlabRepo) DeleteGroup(id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Groups.DeleteGroup(id)
	return err
}

func (repo *GitlabRepo) ChangeGroupName(id int, name string) (*model.Group, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Groups.UpdateGroup(id, &goGitlab.UpdateGroupOptions{
		Name: goGitlab.String(name),
	})
	if err != nil {
		return nil, err
	}

	return repo.GetGroupById(id)
}

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

	return err
}

func (repo *GitlabRepo) RemoveUserFromGroup(groupId int, userId int) error {
	repo.assertIsConnected()

	_, err := repo.client.GroupMembers.RemoveGroupMember(groupId, userId, nil)

	return err
}

func (repo *GitlabRepo) GetAllProjects(search string) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Projects.ListProjects(&goGitlab.ListProjectsOptions{
		Archived:   goGitlab.Bool(false),
		Visibility: goGitlab.Visibility(goGitlab.PublicVisibility),
		Owned:      goGitlab.Bool(true),
		OrderBy:    goGitlab.String("created_at"),
		Search:     goGitlab.String(search),
	})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GitlabRepo) GetPublicAvatarByMail(mail string) (url *string, err error) {
	repo.assertIsConnected()

	avatar, response, err := repo.client.Avatar.GetAvatar(&goGitlab.GetAvatarOptions{Email: &mail})
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("avatar is not available, status code: %d", response.StatusCode)
	}

	return &avatar.AvatarURL, nil
}

func (repo *GitlabRepo) GetProjectById(id int) (*model.Project, error) {
	repo.assertIsConnected()

	gitlabProject, _, err := repo.client.Projects.GetProject(id, &goGitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabProject(gitlabProject)
}

func (repo *GitlabRepo) GetUserById(id int) (*model.User, error) {
	repo.assertIsConnected()

	gitlabUser, _, err := repo.client.Users.GetUser(id, goGitlab.GetUsersOptions{})
	if err != nil {
		return nil, err
	}

	return UserFromGoGitlab(*gitlabUser), nil
}

func (repo *GitlabRepo) GetGroupById(id int) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroup, _, err := repo.client.Groups.GetGroup(id, &goGitlab.GetGroupOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabGroup(gitlabGroup)
}

func (repo *GitlabRepo) GetAllUsers() ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Users.ListUsers(&goGitlab.ListUsersOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) GetAllGroups() ([]*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroups, _, err := repo.client.Groups.ListGroups(&goGitlab.ListGroupsOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabGroups(gitlabGroups)
}

func (repo *GitlabRepo) GetAllProjectsOfGroup(id int) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Groups.ListGroupProjects(id, &goGitlab.ListGroupProjectsOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GitlabRepo) GetAllUsersOfGroup(id int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabMembers, _, err := repo.client.Groups.ListGroupMembers(id, &goGitlab.ListGroupMembersOptions{})
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		users[i] = UserFromGoGitlabGroupMember(*gitlabMember)
	}

	return users, nil
}

func (repo *GitlabRepo) SearchProjectByExpression(expression string) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Search.Projects(expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GitlabRepo) SearchUserByExpression(expression string) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.Users(expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) SearchUserByExpressionInGroup(expression string, groupId int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.UsersByGroup(groupId, expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) SearchUserByExpressionInProject(expression string, projectId int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Search.UsersByProject(projectId, expression, &goGitlab.SearchOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GitlabRepo) SearchGroupByExpression(expression string) ([]*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroups, _, err := repo.client.Groups.SearchGroup(expression)
	if err != nil {
		return nil, err
	}

	groups := make([]*model.Group, len(gitlabGroups))
	for i, gitlabGroup := range gitlabGroups {
		groups[i] = GroupFromGoGitlab(*gitlabGroup)
	}

	return groups, nil
}

func (repo *GitlabRepo) GetPendingProjectInvitations(projectId int) ([]*model.PendingInvite, error) {
	repo.assertIsConnected()

	pendingInvites, _, err := repo.client.Invites.ListPendingProjectInvitations(projectId, nil)
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabPendingInvites(pendingInvites)
}

func (repo *GitlabRepo) GetPendingGroupInvitations(groupId int) ([]*model.PendingInvite, error) {
	repo.assertIsConnected()

	pendingInvites, _, err := repo.client.Invites.ListPendingGroupInvitations(groupId, nil)
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabPendingInvites(pendingInvites)
}

func (repo *GitlabRepo) CreateGroupInvite(groupId int, email string) error {
	repo.assertIsConnected()

	_, _, err := repo.client.Invites.GroupInvites(groupId, &goGitlab.InvitesOptions{
		Email:       &email,
		AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		// Set additional options like AccessLevel, ExpiresAt as needed
	})
	return err
}

func (repo *GitlabRepo) CreateProjectInvite(id int, email string) error {
	repo.assertIsConnected()

	_, _, err := repo.client.Invites.ProjectInvites(id, &goGitlab.InvitesOptions{
		Email:       &email,
		AccessLevel: goGitlab.AccessLevel(goGitlab.DeveloperPermissions),
		// Set additional options like AccessLevel, ExpiresAt as needed
	})
	return err
}

/*
TODO:
	Mit personal access tokens ist es bisher nicht möglich ein Assignment zu schließen bzw. das Pushen zu unterbinden (man bekommt bei alle aufgelisteten Möglichkeiten einen 404 zurück)
	- Not with Push Rules
	- Not with Protect Branches
	- Not with change Project Member Access Level
*/

func (repo *GitlabRepo) DenyPushingToProject(projectId int) error {
	log.Panic("No working option to close an assignment")

	permission := goGitlab.AccessLevelValue(goGitlab.MinimalAccessPermissions)

	return repo.changeProjectMemberPermissions(projectId, permission)
}

func (repo *GitlabRepo) AllowPushingToProject(projectId int) error {
	log.Panic("No working option to reopen an assignment")

	permission := goGitlab.AccessLevelValue(goGitlab.DeveloperPermissions)

	return repo.changeProjectMemberPermissions(projectId, permission)
}

func (repo *GitlabRepo) changeProjectMemberPermissions(projectId int, accessLevel goGitlab.AccessLevelValue) error {
	repo.assertIsConnected()

	members, _, err := repo.client.ProjectMembers.ListAllProjectMembers(projectId, &goGitlab.ListProjectMembersOptions{})
	if err != nil {
		return err
	}

	for _, member := range members {
		if member.AccessLevel == *goGitlab.AccessLevel(goGitlab.OwnerPermissions) {
			continue
		}

		_, _, err := repo.client.ProjectMembers.EditProjectMember(projectId, member.ID, &goGitlab.EditProjectMemberOptions{AccessLevel: &accessLevel})
		if err != nil {
			return err
		}
	}

	return nil
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
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("User with username [%s] not found", username)
	}

	return UserFromGoGitlab(*users[0]), nil
}

func (repo *GitlabRepo) FindUserIDByEmail(email string) (int, error) {
	repo.assertIsConnected()

	listUsersOptions := &goGitlab.ListUsersOptions{
		Search: goGitlab.String(email),
	}

	users, _, err := repo.client.Users.ListUsers(listUsersOptions)
	if err != nil {
		return 0, err
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
			return nil, err
		}

		projects[i] = project
	}

	return projects, nil
}

func (repo *GitlabRepo) convertGitlabProject(gitlabProject *goGitlab.Project) (*model.Project, error) {
	gitlabMembers, _, err := repo.client.ProjectMembers.ListProjectMembers(gitlabProject.ID, &goGitlab.ListProjectMembersOptions{})
	if err != nil {
		return nil, err
	}

	return ProjectFromGoGitlabWithProjectMembers(*gitlabProject, gitlabMembers), nil
}

func (repo *GitlabRepo) convertGitlabGroup(gitlabGroup *goGitlab.Group) (*model.Group, error) {
	Group := GroupFromGoGitlab(*gitlabGroup)

	projects, err := repo.GetAllProjectsOfGroup(gitlabGroup.ID)
	if err != nil {
		return nil, err
	}
	Group.Projects = ConvertProjectPointerSlice(projects)

	members, err := repo.GetAllUsersOfGroup(gitlabGroup.ID)
	if err != nil {
		return nil, err
	}
	Group.Member = ConvertUserPointerSlice(members)

	return Group, nil
}

func (repo *GitlabRepo) convertGitlabGroups(gitlabGroups []*goGitlab.Group) ([]*model.Group, error) {
	groups := make([]*model.Group, len(gitlabGroups))
	for i, gitlabGroup := range gitlabGroups {
		group, err := repo.convertGitlabGroup(gitlabGroup)
		if err != nil {
			return nil, err
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

func (repo *GitlabRepo) getValidatedUserAvatarURL(user *goGitlab.User) (validatedAvatarURL *string) {
	resp, err := http.Get(user.AvatarURL)
	if err != nil {
		log.Printf("Failed to fetch avatar for user %s: %v\n", user.Username, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if isValidImage(resp.Body) {
			return &user.AvatarURL
		}
	} else if resp.StatusCode == http.StatusUnauthorized {
		// Attempt to fetch a public avatar if unauthorized
		if publicAvatarURL, err := repo.GetPublicAvatarByMail(user.Email); err == nil {
			resp, err := http.Get(*publicAvatarURL)
			if err != nil || publicAvatarURL == nil {
				log.Printf("Failed to fetch public avatar for user %s: %v\n", user.Username, err)
				return nil
			}
			defer resp.Body.Close()

			if isValidImage(resp.Body) {
				return publicAvatarURL
			}
		} else {
			log.Printf("Error retrieving public avatar for user %s: %v\n", user.Username, err)
		}
	} else {
		log.Printf("Unexpected status code %d for user %s\n", resp.StatusCode, user.Username)
	}

	return nil
}

func isValidImage(body io.Reader) bool {
	img, _, err := image.Decode(body)
	return err == nil && img != nil
}
