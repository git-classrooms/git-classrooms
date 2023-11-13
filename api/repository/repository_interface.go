package repository

import (
	"backend/model"
)

type Repository interface {
	// Classrooms als Gitlab Group realisieren?

	Login(username string, password string) (*model.User, error)                                                                   //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/gitlab.go#L256
	CreateProject(name string, visibility model.Visibility, description string, member []model.User) (*model.Project, error)       //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L735 https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L466
	CreateClassroom(name string, visibility model.Visibility, description string, memberEmails []string) (*model.Classroom, error) //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L366

	DeleteProject(id int) error                                                 //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1110
	DeleteClassroom(id int) error                                               //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L566
	ChangeClassroomName(id int, name string) (*model.Classroom, error)          //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L495
	AddUserToClassroom(groupId int, userId int) error                           //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/group_members.go#L237
	RemoveUserFromClassroom(groupId int, userId int) error                      //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/group_members.go#L349
	GetAllProjects() ([]*model.Project, error)                                  //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L373
	GetProjectById(id int) (*model.Project, error)                              //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L577
	GetUserById(id int) (*model.User, error)                                    //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/users.go#L169
	GetClassroomById(id int) (*model.Classroom, error)                          //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L288
	GetAllUsers() ([]*model.User, error)                                        //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/users.go#L144
	GetAllClassrooms() ([]*model.Classroom, error)                              //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L149
	GetAllProjectsOfClassroom(id int) ([]*model.Project, error)                 //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L255
	GetAllUsersOfClassroom(id int) ([]*model.User, error)                       //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/group_members.go#L77
	SearchProjectByExpression(expression string) ([]*model.Project, error)      //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L49
	SearchUserByExpression(expression string) ([]*model.User, error)            //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L290
	SearchUserByExpressionInClassroom(expression string) ([]*model.User, error) //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L300
	SearchUserByExpressionInProject(expression string) ([]*model.User, error)   //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/search.go#L310
	SearchClassroomByExpression(expression string) ([]*model.Classroom, error)  //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/groups.go#L609
	GetPendingProjectInvitations(id int) (*string, error)                       //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L85
	GetPendingClassroomInvitations(id int) (*string, error)                     //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L60
	CreateGroupInvite(groupId int, email string) error                          //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L132
	CreateProjectInvite(projectId int, email string) error                      //c https://github.com/xanzy/go-gitlab/blob/v0.93.2/invites.go#L157
	DenyPushingToProject(projectId int) error                                   //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1608
	AllowPushingToProject(projectId int) error                                  //j https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1679 https://github.com/xanzy/go-gitlab/blob/v0.93.2/projects.go#L1564

	// Function without associated implementation option
	Logout() error
	JoinClassroom(groupId int) error
	CreateProjectByTemplate(templateUrl string) (*model.Project, error)
}
