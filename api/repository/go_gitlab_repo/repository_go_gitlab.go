package go_gitlab_repo

import (
	"backend/model"
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type GoGitlabRepo struct {
	client      *gitlab.Client
	isConnected bool
}

func NewGoGitlabRepo() *GoGitlabRepo {
	return &GoGitlabRepo{client: nil, isConnected: false}
}

func (repo *GoGitlabRepo) Login(token string, username string) (*model.User, error) {
	cli, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.hs-flensburg.de"))
	if err != nil {
		return nil, err
	}
	repo.client = cli

	return repo.getUserByUsername(username)
}

func (repo *GoGitlabRepo) DeleteProject(id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Projects.DeleteProject(id)
	return err
}

func (repo *GoGitlabRepo) DeleteGroup(id int) error {
	repo.assertIsConnected()
	_, err := repo.client.Groups.DeleteGroup(id)
	return err
}

func (repo *GoGitlabRepo) ChangeGroupName(id int, name string) (*model.Group, error) {
    repo.assertIsConnected()

    _, _, err := repo.client.Groups.UpdateGroup(id, &gitlab.UpdateGroupOptions{
        Name: gitlab.String(name),
    })
    if err != nil {
        return nil, err
    }

    return repo.GetGroupById(id)
}

func (repo *GoGitlabRepo) AddUserToGroup(groupId int, userId int) (error) {
	repo.assertIsConnected()

    accessLevel := gitlab.DeveloperPermissions

    _, _, err := repo.client.GroupMembers.AddGroupMember(groupId, &gitlab.AddGroupMemberOptions{
        UserID:      &userId,
        AccessLevel: &accessLevel,
    })

	return err
}

func (repo *GoGitlabRepo) RemoveUserFromGroup(groupId int, userId int) (error) {
    repo.assertIsConnected()

    _, err := repo.client.GroupMembers.RemoveGroupMember(groupId, userId, nil)

    return err
}

func (repo *GoGitlabRepo) GetAllProjects() ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GoGitlabRepo) GetProjectById(id int) (*model.Project, error) {
	repo.assertIsConnected()

	gitlabProject, _, err := repo.client.Projects.GetProject(id, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabProject(gitlabProject)
}

func (repo *GoGitlabRepo) GetUserById(id int) (*model.User, error) {
	repo.assertIsConnected()

	gitlabUser, _, err := repo.client.Users.GetUser(id, gitlab.GetUsersOptions{})
	if err != nil {
		return nil, err
	}

	return UserFromGoGitlab(*gitlabUser), nil
}

func (repo *GoGitlabRepo) GetGroupById(id int) (*model.Group, error) {
	repo.assertIsConnected()

	gitlabGroup, _, err := repo.client.Groups.GetGroup(id, &gitlab.GetGroupOptions{})
	if err != nil {
		return nil, err
	}

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

func (repo *GoGitlabRepo) GetAllUsers() ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabUsers, _, err := repo.client.Users.ListUsers(&gitlab.ListUsersOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GoGitlabRepo) GetAllGroups() ([]*model.Group, error) {
	repo.assertIsConnected()

	_, _, err := repo.client.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	if err != nil {
		return nil, err
	}

	// TODO Jannes continue here

	return nil, nil
}

func (repo *GoGitlabRepo) GetAllProjectsOfGroup(id int) ([]*model.Project, error) {
	repo.assertIsConnected()

	gitlabProjects, _, err := repo.client.Groups.ListGroupProjects(id, &gitlab.ListGroupProjectsOptions{})
	if err != nil {
		return nil, err
	}

	return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GoGitlabRepo) GetAllUsersOfGroup(id int) ([]*model.User, error) {
	repo.assertIsConnected()

	gitlabMembers, _, err := repo.client.Groups.ListGroupMembers(id, &gitlab.ListGroupMembersOptions{})
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		users[i] = UserFromGoGitlabGroupMember(*gitlabMember)
	}

	return users, nil
}

func (repo *GoGitlabRepo) SearchProjectByExpression(expression string) ([]*model.Project, error) {
    repo.assertIsConnected()

    gitlabProjects, _, err := repo.client.Search.Projects(expression, &gitlab.SearchOptions{})
    if err != nil {
        return nil, err
    }

    return repo.convertGitlabProjects(gitlabProjects)
}

func (repo *GoGitlabRepo) SearchUserByExpression(expression string) ([]*model.User, error) {
    repo.assertIsConnected()

    gitlabUsers, _, err := repo.client.Search.Users(expression, &gitlab.SearchOptions{})
    if err != nil {
        return nil, err
    }

    return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GoGitlabRepo) SearchUserByExpressionInGroup(expression string, groupId int) ([]*model.User, error) {
    repo.assertIsConnected()

    gitlabUsers, _, err := repo.client.Search.UsersByGroup(groupId, expression, &gitlab.SearchOptions{})
    if err != nil {
        return nil, err
    }

    return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GoGitlabRepo) SearchUserByExpressionInProject(expression string, projectId int) ([]*model.User, error) {
    repo.assertIsConnected()

    gitlabUsers, _, err := repo.client.Search.UsersByProject(projectId, expression, &gitlab.SearchOptions{})
    if err != nil {
        return nil, err
    }

    return repo.convertGitlabUsers(gitlabUsers)
}

func (repo *GoGitlabRepo) SearchGroupByExpression(expression string) ([]*model.Group, error) {
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

func (repo *GoGitlabRepo) GetPendingProjectInvitations(projectId int) ([]*model.PendingInvite, error) {
    repo.assertIsConnected()

    pendingInvites, _, err := repo.client.Invites.ListPendingProjectInvitations(projectId, nil)
    if err != nil {
        return nil, err
    }

    return repo.convertGitlabPendingInvites(pendingInvites)
}

func (repo *GoGitlabRepo) GetPendingGroupInvitations(groupId int) ([]*model.PendingInvite, error) {
    repo.assertIsConnected()

    pendingInvites, _, err := repo.client.Invites.ListPendingGroupInvitations(groupId, nil)
    if err != nil {
        return nil, err
    }

    return repo.convertGitlabPendingInvites(pendingInvites)
}

func (repo *GoGitlabRepo) CreateGroupInvite(groupId int, email string) (error) {
    repo.assertIsConnected()

    _, _, err := repo.client.Invites.GroupInvites(groupId, &gitlab.InvitesOptions{
        Email: &email,
        // Set additional options like AccessLevel, ExpiresAt as needed
    })
    return err
}

func (repo *GoGitlabRepo) CreateProjectInvite(id int, email string) (error) {
    repo.assertIsConnected()

    _, _, err := repo.client.Invites.ProjectInvites(id, &gitlab.InvitesOptions{
        Email: &email,
        // Set additional options like AccessLevel, ExpiresAt as needed
    })
    return err
}


func (repo *GoGitlabRepo) assertIsConnected() {
	if repo.client == nil {
		panic("No connection to Gitlab! Make sure you have executed Login()")
	}
}

func (repo *GoGitlabRepo) getUserByUsername(username string) (*model.User, error) {
	repo.assertIsConnected()

	users, _, err := repo.client.Users.ListUsers(&gitlab.ListUsersOptions{
		Search: gitlab.String(username),
	})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("User with username [%s] not found", username)
	}

	return UserFromGoGitlab(*users[0]), nil
}

func (repo *GoGitlabRepo) convertGitlabUsers(gitlabUsers []*gitlab.User) ([]*model.User, error) {
	users := make([]*model.User, len(gitlabUsers))
	for i, gitlabUser := range gitlabUsers {
		users[i] = UserFromGoGitlab(*gitlabUser)
	}

	return users, nil
}

func (repo *GoGitlabRepo) convertGitlabProjects(gitlabProjects []*gitlab.Project) ([]*model.Project, error) {
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

func (repo *GoGitlabRepo) convertGitlabProject(gitlabProject *gitlab.Project) (*model.Project, error) {
	gitlabMembers, _, err := repo.client.ProjectMembers.ListProjectMembers(gitlabProject.ID, &gitlab.ListProjectMembersOptions{})
	if err != nil {
		return nil, err
	}

	return ProjectFromGoGitlabWithProjectMembers(*gitlabProject, gitlabMembers), nil
}

func (repo *GoGitlabRepo) convertGitlabPendingInvites(gitlabPendingInvites []*gitlab.PendingInvite) ([]*model.PendingInvite, error) {
	pendingInvites := make([]*model.PendingInvite, len(gitlabPendingInvites))
	for i, gitlabPendingInvite := range gitlabPendingInvites {
		pendingInvites[i] = PendingInviteFromGoGitlab(*gitlabPendingInvite)
	}

	return pendingInvites, nil
}
