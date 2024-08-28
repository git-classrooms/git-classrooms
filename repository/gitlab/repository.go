package gitlab

import (
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type Repository interface {
	// Access
	Login(token string) error
	GroupAccessLogin(token string) error
	CreateGroupAccessToken(groupID int, name string, accessLevel model.AccessLevelValue, expiresAt time.Time, scopes ...string) (*model.GroupAccessToken, error)
	RotateGroupAccessToken(groupID int, tokenID int, expiresAt time.Time) (*model.GroupAccessToken, error)

	// Group
	CreateGroup(name string, visibility model.Visibility, description string) (*model.Group, error)
	CreateSubGroup(name string, parentId int, visibility model.Visibility, description string) (*model.Group, error)
	DeleteGroup(id int) error
	ChangeGroupName(id int, name string) (*model.Group, error)
	ChangeGroupDescription(id int, description string) (*model.Group, error)
	AddUserToGroup(groupId int, userId int, accessLevel model.AccessLevelValue) error
	RemoveUserFromGroup(groupId int, userId int) error
	GetGroupById(id int) (*model.Group, error)
	GetAllGroups() ([]*model.Group, error)
	SearchGroupByExpression(expression string) ([]*model.Group, error)
	CreateGroupInvite(groupId int, email string) error
	GetPendingGroupInvitations(groupId int) ([]*model.PendingInvite, error)
	ChangeUserAccessLevelInGroup(groupId int, userId int, accessLevel model.AccessLevelValue) error
	GetAccessLevelOfUserInGroup(groupId int, userId int) (model.AccessLevelValue, error)

	// User
	GetCurrentUser() (*model.User, error)
	GetUserById(id int) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	GetAllUsersOfGroup(id int) ([]*model.User, error)
	SearchUserByExpression(expression string) ([]*model.User, error)
	SearchUserByExpressionInGroup(expression string, groupId int) ([]*model.User, error)
	SearchUserByExpressionInProject(expression string, projectId int) ([]*model.User, error)
	FindUserIDByEmail(email string) (int, error)

	// Project
	CreateProject(name string, visibility model.Visibility, description string, member []model.User) (*model.Project, error)
	DeleteProject(id int) error
	GetAllProjects(search string) ([]*model.Project, error)
	GetProjectById(id int) (*model.Project, error)
	GetAllProjectsOfGroup(id int) ([]*model.Project, error)
	SearchProjectByExpression(expression string) ([]*model.Project, error)
	CreateProjectInvite(projectId int, email string) error
	GetPendingProjectInvitations(projectId int) ([]*model.PendingInvite, error)
	DenyPushingToProject(projectId int) error
	AllowPushingToProject(projectId int) error
	ForkProject(projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error)
	ForkProjectWithOnlyDefaultBranch(projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error)
	AddProjectMembers(projectId int, members []model.User) (*model.Project, error)
	GetNamespaceOfProject(projectId int) (*string, error)
	ChangeUserAccessLevelInProject(projectId int, userId int, accessLevel model.AccessLevelValue) error
	GetAccessLevelOfUserInProject(projectId int, userId int) (model.AccessLevelValue, error)
	ChangeProjectName(projectId int, name string) (*model.Project, error)
	ChangeProjectDescription(projectId int, description string) (*model.Project, error)
	GetProjectPipelineTestReportSummary(projectId, pipelineId int) (*model.TestReport, error)
	GetProjectLatestPipelineTestReportSummary(projectId int, ref *string) (*model.TestReport, error)

	// Branches
	CreateBranch(projectId int, branchName string, fromBranch string) (*model.Branch, error)
	ProtectBranch(projectId int, branchName string, accessLevel model.AccessLevelValue) error
	UnprotectBranch(projectId int, branchName string) error
	CreateMergeRequest(projectId int, sourceBranch string, targetBranch string, title string, description string, assigneeId int, recviewerId int) error
	ProtectedBranchExists(projectId int, branchName string) (bool, error)
	BranchExists(projectId int, branchName string) (bool, error)

	// Runners
	GetAvailableRunnersForGitLab() ([]*model.Runner, error)
	GetAvailableRunnersForGroup(groupId int) ([]*model.Runner, error)
	CheckIfFileExistsInProject(projectId int, filePath string) (bool, error)
	GetProjectLanguages(projectId int) (map[string]float32, error)
}
