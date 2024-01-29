// Code generated by MockGen. DO NOT EDIT.
// Source: backend/api/repository (interfaces: Repository)

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "backend/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddProjectMembers mocks base method.
func (m *MockRepository) AddProjectMembers(arg0 int, arg1 []model.User) (*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProjectMembers", arg0, arg1)
	ret0, _ := ret[0].(*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProjectMembers indicates an expected call of AddProjectMembers.
func (mr *MockRepositoryMockRecorder) AddProjectMembers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProjectMembers", reflect.TypeOf((*MockRepository)(nil).AddProjectMembers), arg0, arg1)
}

// AddUserToGroup mocks base method.
func (m *MockRepository) AddUserToGroup(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToGroup indicates an expected call of AddUserToGroup.
func (mr *MockRepositoryMockRecorder) AddUserToGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToGroup", reflect.TypeOf((*MockRepository)(nil).AddUserToGroup), arg0, arg1)
}

// AllowPushingToProject mocks base method.
func (m *MockRepository) AllowPushingToProject(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowPushingToProject", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AllowPushingToProject indicates an expected call of AllowPushingToProject.
func (mr *MockRepositoryMockRecorder) AllowPushingToProject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowPushingToProject", reflect.TypeOf((*MockRepository)(nil).AllowPushingToProject), arg0)
}

// ChangeGroupName mocks base method.
func (m *MockRepository) ChangeGroupName(arg0 int, arg1 string) (*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeGroupName", arg0, arg1)
	ret0, _ := ret[0].(*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeGroupName indicates an expected call of ChangeGroupName.
func (mr *MockRepositoryMockRecorder) ChangeGroupName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeGroupName", reflect.TypeOf((*MockRepository)(nil).ChangeGroupName), arg0, arg1)
}

// CreateGroup mocks base method.
func (m *MockRepository) CreateGroup(arg0 string, arg1 model.Visibility, arg2 string, arg3 []string) (*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockRepositoryMockRecorder) CreateGroup(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockRepository)(nil).CreateGroup), arg0, arg1, arg2, arg3)
}

// CreateGroupInvite mocks base method.
func (m *MockRepository) CreateGroupInvite(arg0 int, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroupInvite", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateGroupInvite indicates an expected call of CreateGroupInvite.
func (mr *MockRepositoryMockRecorder) CreateGroupInvite(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroupInvite", reflect.TypeOf((*MockRepository)(nil).CreateGroupInvite), arg0, arg1)
}

// CreateProject mocks base method.
func (m *MockRepository) CreateProject(arg0 string, arg1 model.Visibility, arg2 string, arg3 []model.User) (*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockRepositoryMockRecorder) CreateProject(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockRepository)(nil).CreateProject), arg0, arg1, arg2, arg3)
}

// CreateProjectInvite mocks base method.
func (m *MockRepository) CreateProjectInvite(arg0 int, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProjectInvite", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProjectInvite indicates an expected call of CreateProjectInvite.
func (mr *MockRepositoryMockRecorder) CreateProjectInvite(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProjectInvite", reflect.TypeOf((*MockRepository)(nil).CreateProjectInvite), arg0, arg1)
}

// DeleteGroup mocks base method.
func (m *MockRepository) DeleteGroup(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockRepositoryMockRecorder) DeleteGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockRepository)(nil).DeleteGroup), arg0)
}

// DeleteProject mocks base method.
func (m *MockRepository) DeleteProject(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject.
func (mr *MockRepositoryMockRecorder) DeleteProject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockRepository)(nil).DeleteProject), arg0)
}

// DenyPushingToProject mocks base method.
func (m *MockRepository) DenyPushingToProject(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DenyPushingToProject", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DenyPushingToProject indicates an expected call of DenyPushingToProject.
func (mr *MockRepositoryMockRecorder) DenyPushingToProject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DenyPushingToProject", reflect.TypeOf((*MockRepository)(nil).DenyPushingToProject), arg0)
}

// ForkProject mocks base method.
func (m *MockRepository) ForkProject(arg0 int, arg1 string) (*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForkProject", arg0, arg1)
	ret0, _ := ret[0].(*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForkProject indicates an expected call of ForkProject.
func (mr *MockRepositoryMockRecorder) ForkProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForkProject", reflect.TypeOf((*MockRepository)(nil).ForkProject), arg0, arg1)
}

// GetAllGroups mocks base method.
func (m *MockRepository) GetAllGroups() ([]*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllGroups")
	ret0, _ := ret[0].([]*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllGroups indicates an expected call of GetAllGroups.
func (mr *MockRepositoryMockRecorder) GetAllGroups() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllGroups", reflect.TypeOf((*MockRepository)(nil).GetAllGroups))
}

// GetAllProjects mocks base method.
func (m *MockRepository) GetAllProjects() ([]*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProjects")
	ret0, _ := ret[0].([]*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProjects indicates an expected call of GetAllProjects.
func (mr *MockRepositoryMockRecorder) GetAllProjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProjects", reflect.TypeOf((*MockRepository)(nil).GetAllProjects))
}

// GetAllProjectsOfGroup mocks base method.
func (m *MockRepository) GetAllProjectsOfGroup(arg0 int) ([]*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProjectsOfGroup", arg0)
	ret0, _ := ret[0].([]*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProjectsOfGroup indicates an expected call of GetAllProjectsOfGroup.
func (mr *MockRepositoryMockRecorder) GetAllProjectsOfGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProjectsOfGroup", reflect.TypeOf((*MockRepository)(nil).GetAllProjectsOfGroup), arg0)
}

// GetAllUsers mocks base method.
func (m *MockRepository) GetAllUsers() ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers")
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockRepositoryMockRecorder) GetAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockRepository)(nil).GetAllUsers))
}

// GetAllUsersOfGroup mocks base method.
func (m *MockRepository) GetAllUsersOfGroup(arg0 int) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsersOfGroup", arg0)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsersOfGroup indicates an expected call of GetAllUsersOfGroup.
func (mr *MockRepositoryMockRecorder) GetAllUsersOfGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsersOfGroup", reflect.TypeOf((*MockRepository)(nil).GetAllUsersOfGroup), arg0)
}

// GetGroupById mocks base method.
func (m *MockRepository) GetGroupById(arg0 int) (*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupById", arg0)
	ret0, _ := ret[0].(*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupById indicates an expected call of GetGroupById.
func (mr *MockRepositoryMockRecorder) GetGroupById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupById", reflect.TypeOf((*MockRepository)(nil).GetGroupById), arg0)
}

// GetNamespaceOfGroup mocks base method.
func (m *MockRepository) GetNamespaceOfGroup(arg0 int) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespaceOfGroup", arg0)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespaceOfGroup indicates an expected call of GetNamespaceOfGroup.
func (mr *MockRepositoryMockRecorder) GetNamespaceOfGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespaceOfGroup", reflect.TypeOf((*MockRepository)(nil).GetNamespaceOfGroup), arg0)
}

// GetNamespaceOfProject mocks base method.
func (m *MockRepository) GetNamespaceOfProject(arg0 int) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespaceOfProject", arg0)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespaceOfProject indicates an expected call of GetNamespaceOfProject.
func (mr *MockRepositoryMockRecorder) GetNamespaceOfProject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespaceOfProject", reflect.TypeOf((*MockRepository)(nil).GetNamespaceOfProject), arg0)
}

// GetPendingGroupInvitations mocks base method.
func (m *MockRepository) GetPendingGroupInvitations(arg0 int) ([]*model.PendingInvite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingGroupInvitations", arg0)
	ret0, _ := ret[0].([]*model.PendingInvite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingGroupInvitations indicates an expected call of GetPendingGroupInvitations.
func (mr *MockRepositoryMockRecorder) GetPendingGroupInvitations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingGroupInvitations", reflect.TypeOf((*MockRepository)(nil).GetPendingGroupInvitations), arg0)
}

// GetPendingProjectInvitations mocks base method.
func (m *MockRepository) GetPendingProjectInvitations(arg0 int) ([]*model.PendingInvite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingProjectInvitations", arg0)
	ret0, _ := ret[0].([]*model.PendingInvite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingProjectInvitations indicates an expected call of GetPendingProjectInvitations.
func (mr *MockRepositoryMockRecorder) GetPendingProjectInvitations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingProjectInvitations", reflect.TypeOf((*MockRepository)(nil).GetPendingProjectInvitations), arg0)
}

// GetProjectById mocks base method.
func (m *MockRepository) GetProjectById(arg0 int) (*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectById", arg0)
	ret0, _ := ret[0].(*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectById indicates an expected call of GetProjectById.
func (mr *MockRepositoryMockRecorder) GetProjectById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectById", reflect.TypeOf((*MockRepository)(nil).GetProjectById), arg0)
}

// GetUserById mocks base method.
func (m *MockRepository) GetUserById(arg0 int) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockRepositoryMockRecorder) GetUserById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockRepository)(nil).GetUserById), arg0)
}

// Login mocks base method.
func (m *MockRepository) Login(arg0, arg1 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockRepositoryMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockRepository)(nil).Login), arg0, arg1)
}

// RemoveUserFromGroup mocks base method.
func (m *MockRepository) RemoveUserFromGroup(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromGroup indicates an expected call of RemoveUserFromGroup.
func (mr *MockRepositoryMockRecorder) RemoveUserFromGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromGroup", reflect.TypeOf((*MockRepository)(nil).RemoveUserFromGroup), arg0, arg1)
}

// SearchGroupByExpression mocks base method.
func (m *MockRepository) SearchGroupByExpression(arg0 string) ([]*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchGroupByExpression", arg0)
	ret0, _ := ret[0].([]*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchGroupByExpression indicates an expected call of SearchGroupByExpression.
func (mr *MockRepositoryMockRecorder) SearchGroupByExpression(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchGroupByExpression", reflect.TypeOf((*MockRepository)(nil).SearchGroupByExpression), arg0)
}

// SearchProjectByExpression mocks base method.
func (m *MockRepository) SearchProjectByExpression(arg0 string) ([]*model.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchProjectByExpression", arg0)
	ret0, _ := ret[0].([]*model.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchProjectByExpression indicates an expected call of SearchProjectByExpression.
func (mr *MockRepositoryMockRecorder) SearchProjectByExpression(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchProjectByExpression", reflect.TypeOf((*MockRepository)(nil).SearchProjectByExpression), arg0)
}

// SearchUserByExpression mocks base method.
func (m *MockRepository) SearchUserByExpression(arg0 string) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserByExpression", arg0)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchUserByExpression indicates an expected call of SearchUserByExpression.
func (mr *MockRepositoryMockRecorder) SearchUserByExpression(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserByExpression", reflect.TypeOf((*MockRepository)(nil).SearchUserByExpression), arg0)
}

// SearchUserByExpressionInGroup mocks base method.
func (m *MockRepository) SearchUserByExpressionInGroup(arg0 string, arg1 int) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserByExpressionInGroup", arg0, arg1)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchUserByExpressionInGroup indicates an expected call of SearchUserByExpressionInGroup.
func (mr *MockRepositoryMockRecorder) SearchUserByExpressionInGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserByExpressionInGroup", reflect.TypeOf((*MockRepository)(nil).SearchUserByExpressionInGroup), arg0, arg1)
}

// SearchUserByExpressionInProject mocks base method.
func (m *MockRepository) SearchUserByExpressionInProject(arg0 string, arg1 int) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserByExpressionInProject", arg0, arg1)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchUserByExpressionInProject indicates an expected call of SearchUserByExpressionInProject.
func (mr *MockRepositoryMockRecorder) SearchUserByExpressionInProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserByExpressionInProject", reflect.TypeOf((*MockRepository)(nil).SearchUserByExpressionInProject), arg0, arg1)
}
