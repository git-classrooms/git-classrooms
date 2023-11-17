package repository

import (
	"backend/model"
)

type Repository interface {
	// Groups als Gitlab Group realisieren?

	Login(token string, username string) (*model.User, error)                                                                //j
	CreateProject(name string, visibility model.Visibility, description string, member []model.User) (*model.Project, error) //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L735 https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L466
	CreateGroup(name string, visibility model.Visibility, description string, memberEmails []string) (*model.Group, error)   //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L366

	DeleteProject(id int) error                                               //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1110
	DeleteGroup(id int) error                                                 //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L566
	ChangeGroupName(id int, name string) (*model.Group, error)                //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L495
	AddUserToGroup(groupId int, userId int) error                             //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/group_members.go#L237
	RemoveUserFromGroup(groupId int, userId int) error                        //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/group_members.go#L349
	GetAllProjects() ([]*model.Project, error)                                //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L373
	GetProjectById(id int) (*model.Project, error)                            //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L577
	GetUserById(id int) (*model.User, error)                                  //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/users.go#L169
	GetGroupById(id int) (*model.Group, error)                                //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L288
	GetAllUsers() ([]*model.User, error)                                      //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/users.go#L144
	GetAllGroups() ([]*model.Group, error)                                    //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L149
	GetAllProjectsOfGroup(id int) ([]*model.Project, error)                   //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L255
	GetAllUsersOfGroup(id int) ([]*model.User, error)                         //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/group_members.go#L77
	SearchProjectByExpression(expression string) ([]*model.Project, error)    //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L49
	SearchUserByExpression(expression string) ([]*model.User, error)          //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L290
	SearchUserByExpressionInGroup(expression string) ([]*model.User, error)   //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L300
	SearchUserByExpressionInProject(expression string) ([]*model.User, error) //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L310
	SearchGroupByExpression(expression string) ([]*model.Group, error)        //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L609
	GetPendingProjectInvitations(id int) (*string, error)                     //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L85
	GetPendingGroupInvitations(id int) (*string, error)                       //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L60
	CreateGroupInvite(groupId int, email string) error                        //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L132
	CreateProjectInvite(projectId int, email string) error                    //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L157
	DenyPushingToProject(projectId int) error                                 //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1608
	AllowPushingToProject(projectId int) error                                //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1679 https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1564

	// Function without associated implementation option
	Logout() error
	JoinGroup(groupId int) error
	CreateProjectByTemplate(templateUrl string) (*model.Project, error)
}
