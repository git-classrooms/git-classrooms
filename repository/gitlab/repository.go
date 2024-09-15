package gitlab

import (
	"context"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

type Repository interface {
	// Access
	Login(token string) error
	GroupAccessLogin(token string) error
	CreateGroupAccessToken(ctx context.Context, groupID int, name string, accessLevel model.AccessLevelValue, expiresAt time.Time, scopes ...string) (*model.GroupAccessToken, error)
	RotateGroupAccessToken(ctx context.Context, groupID int, tokenID int, expiresAt time.Time) (*model.GroupAccessToken, error)

	// Group
	CreateGroup(ctx context.Context, name string, visibility model.Visibility, description string) (*model.Group, error)
	CreateSubGroup(ctx context.Context, name string, parentId int, visibility model.Visibility, description string) (*model.Group, error)
	DeleteGroup(ctx context.Context, id int) error
	ChangeGroupName(ctx context.Context, id int, name string) (*model.Group, error)
	ChangeGroupDescription(ctx context.Context, id int, description string) (*model.Group, error)
	AddUserToGroup(ctx context.Context, groupId int, userId int, accessLevel model.AccessLevelValue) error
	RemoveUserFromGroup(ctx context.Context, groupId int, userId int) error
	GetGroupById(ctx context.Context, id int) (*model.Group, error)
	GetAllGroups(ctx context.Context) ([]*model.Group, error)
	SearchGroupByExpression(ctx context.Context, expression string) ([]*model.Group, error)
	CreateGroupInvite(ctx context.Context, groupId int, email string) error
	GetPendingGroupInvitations(ctx context.Context, groupId int) ([]*model.PendingInvite, error)
	ChangeUserAccessLevelInGroup(ctx context.Context, groupId int, userId int, accessLevel model.AccessLevelValue) error
	GetAccessLevelOfUserInGroup(ctx context.Context, groupId int, userId int) (model.AccessLevelValue, error)

	// User
	GetCurrentUser(ctx context.Context) (*model.User, error)
	GetUserById(ctx context.Context, id int) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetAllUsersOfGroup(ctx context.Context, id int) ([]*model.User, error)
	SearchUserByExpression(ctx context.Context, expression string) ([]*model.User, error)
	SearchUserByExpressionInGroup(ctx context.Context, expression string, groupId int) ([]*model.User, error)
	SearchUserByExpressionInProject(ctx context.Context, expression string, projectId int) ([]*model.User, error)
	FindUserIDByEmail(ctx context.Context, email string) (int, error)

	// Project
	CreateProject(ctx context.Context, name string, visibility model.Visibility, description string, member []model.User) (*model.Project, error)
	DeleteProject(ctx context.Context, id int) error
	GetAllProjects(ctx context.Context, search string) ([]*model.Project, error)
	GetProjectById(ctx context.Context, id int) (*model.Project, error)
	GetAllProjectsOfGroup(ctx context.Context, id int) ([]*model.Project, error)
	SearchProjectByExpression(ctx context.Context, expression string) ([]*model.Project, error)
	CreateProjectInvite(ctx context.Context, projectId int, email string) error
	GetPendingProjectInvitations(ctx context.Context, projectId int) ([]*model.PendingInvite, error)
	DenyPushingToProject(ctx context.Context, projectId int) error
	AllowPushingToProject(ctx context.Context, projectId int) error
	ForkProject(ctx context.Context, projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error)
	ForkProjectWithOnlyDefaultBranch(ctx context.Context, projectId int, visibility model.Visibility, namespaceId int, name string, description string) (*model.Project, error)
	AddProjectMembers(ctx context.Context, projectId int, members []model.User) (*model.Project, error)
	GetNamespaceOfProject(ctx context.Context, projectId int) (*string, error)
	ChangeUserAccessLevelInProject(ctx context.Context, projectId int, userId int, accessLevel model.AccessLevelValue) error
	GetAccessLevelOfUserInProject(ctx context.Context, projectId int, userId int) (model.AccessLevelValue, error)
	ChangeProjectName(ctx context.Context, projectId int, name string) (*model.Project, error)
	ChangeProjectDescription(ctx context.Context, projectId int, description string) (*model.Project, error)
	GetProjectLatestPipeline(ctx context.Context, projectId int, ref *string) (*model.Pipeline, error)
	GetProjectPipelineTestReportSummary(ctx context.Context, projectId, pipelineId int) (*model.TestReport, error)
	GetProjectLatestPipelineTestReportSummary(ctx context.Context, projectId int, ref *string) (*model.TestReport, error)

	// Branches
	CreateBranch(ctx context.Context, projectId int, branchName string, fromBranch string) (*model.Branch, error)
	ProtectBranch(ctx context.Context, projectId int, branchName string, accessLevel model.AccessLevelValue) error
	UnprotectBranch(ctx context.Context, projectId int, branchName string) error
	CreateMergeRequest(ctx context.Context, projectId int, sourceBranch string, targetBranch string, title string, description string, assigneeId int, recviewerId int) error
	ProtectedBranchExists(ctx context.Context, projectId int, branchName string) (bool, error)
	BranchExists(ctx context.Context, projectId int, branchName string) (bool, error)

	// Runners
	GetAvailableRunnersForGitLab(ctx context.Context) ([]*model.Runner, error)
	GetAvailableRunnersForGroup(ctx context.Context, groupId int) ([]*model.Runner, error)
	CheckIfFileExistsInProject(ctx context.Context, projectId int, filePath string) (bool, error)
	GetProjectLanguages(ctx context.Context, projectId int) (map[string]float32, error)
}
