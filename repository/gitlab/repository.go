package gitlab

import (
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

// Repository defines the operations for interacting with VCS resources.
type Repository interface {
	// Access
	Login(token string) error
	GroupAccessLogin(token string) error
	CreateGroupAccessToken(groupID int, name string, accessLevel model.AccessLevelValue, expiresAt time.Time, scopes ...string) (*model.GroupAccessToken, error)
	GetGroupAccessToken(groupID int, tokenID int) (*model.GroupAccessToken, error)
	RotateGroupAccessToken(groupID int, tokenID int, expiresAt time.Time) (*model.GroupAccessToken, error)

	// Group
	CreateGroup(name string, visibility model.Visibility, description string) (*model.Group, error)
	CreateSubGroup(name string, path string, parentID int, visibility model.Visibility, description string) (*model.Group, error)
	DeleteGroup(id int) error
	ChangeGroupName(id int, name string) (*model.Group, error)
	ChangeGroupDescription(id int, description string) (*model.Group, error)
	AddUserToGroup(groupID int, userID int, accessLevel model.AccessLevelValue) error
	RemoveUserFromGroup(groupID int, userID int) error
	GetGroupByID(id int) (*model.Group, error)
	GetAllGroups() ([]*model.Group, error)
	SearchGroupByExpression(expression string) ([]*model.Group, error)
	CreateGroupInvite(groupID int, email string) error
	GetPendingGroupInvitations(groupID int) ([]*model.PendingInvite, error)
	ChangeUserAccessLevelInGroup(groupID int, userID int, accessLevel model.AccessLevelValue) error
	GetAccessLevelOfUserInGroup(groupID int, userID int) (model.AccessLevelValue, error)

	// User
	GetCurrentUser() (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	GetAllUsersOfGroup(id int) ([]*model.User, error)
	SearchUserByExpression(expression string) ([]*model.User, error)
	SearchUserByExpressionInGroup(expression string, groupID int) ([]*model.User, error)
	SearchUserByExpressionInProject(expression string, projectID int) ([]*model.User, error)
	FindUserIDByEmail(email string) (int, error)

	// Project
	CreateProject(name string, visibility model.Visibility, description string, member []model.User) (*model.Project, error)
	DeleteProject(id int) error
	GetAllProjects(search string) ([]*model.Project, error)
	GetProjectByID(id int) (*model.Project, error)
	GetAllProjectsOfGroup(id int) ([]*model.Project, error)
	SearchProjectByExpression(expression string) ([]*model.Project, error)
	CreateProjectInvite(projectID int, email string) error
	GetPendingProjectInvitations(projectID int) ([]*model.PendingInvite, error)
	DenyPushingToProject(projectID int) error
	AllowPushingToProject(projectID int) error
	ForkProject(projectID int, visibility model.Visibility, namespaceID int, name string, description string) (*model.Project, error)
	ForkProjectWithOnlyDefaultBranch(projectID int, visibility model.Visibility, namespaceID int, name string, description string) (*model.Project, error)
	AddProjectMembers(projectID int, members []model.User) (*model.Project, error)
	AddProjectMember(projectID int, userID int, accessLevel model.AccessLevelValue) error
	RemoveUserFromProject(projectID int, userID int) error
	GetNamespaceOfProject(projectID int) (*string, error)
	ChangeUserAccessLevelInProject(projectID int, userID int, accessLevel model.AccessLevelValue) error
	GetAccessLevelOfUserInProject(projectID int, userID int) (model.AccessLevelValue, error)
	ChangeProjectName(projectID int, name string) (*model.Project, error)
	ChangeProjectDescription(projectID int, description string) (*model.Project, error)
	GetProjectLatestPipeline(projectID int, ref *string) (*model.Pipeline, error)
	GetProjectPipelineTestReportSummary(projectID, pipelineID int) (*model.TestReport, error)
	GetProjectLatestPipelineTestReportSummary(projectID int, ref *string) (*model.TestReport, error)

	// Branches
	CreateBranch(projectID int, branchName string, fromBranch string) (*model.Branch, error)
	ProtectBranch(projectID int, branchName string, accessLevel model.AccessLevelValue) error
	UnprotectBranch(projectID int, branchName string) error
	CreateMergeRequest(projectID int, sourceBranch string, targetBranch string, title string, description string, assigneeID int, reviewerID int) error
	ProtectedBranchExists(projectID int, branchName string) (bool, error)
	BranchExists(projectID int, branchName string) (bool, error)

	// Runners
	GetAvailableRunnersForGitLab() ([]*model.Runner, error)
	GetAvailableRunnersForGroup(groupID int) ([]*model.Runner, error)
	CheckIfFileExistsInProject(projectID int, filePath string) (bool, error)
	GetProjectLanguages(projectID int) (map[string]float32, error)
}
