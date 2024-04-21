package context

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
)

type contextKey string

const (
	gitlabRepoKey                      contextKey = "gitlab-repo"
	gitlabUserIDKey                    contextKey = "gitlab-user-id"
	gitlabGroupIDKey                   contextKey = "gitlab-group-id"
	gitlabProjectIDKey                 contextKey = "gitlab-project-id"
	userIDKey                          contextKey = "user-id"
	ownedClassroomKey                  contextKey = "owned-classroom"
	ownedClassroomAssignmentKey        contextKey = "owned-classroom-assignment"
	ownedClassroomAssignmentProjectKey contextKey = "owned-classroom-assignment-project"
	ownedClassroomMemberKey            contextKey = "owned-classroom-member"
	ownedClassroomTeamKey              contextKey = "owned-classroom-team"
	ownedClassroomTeamMemberKey        contextKey = "owned-classroom-team-member"
	joinedClassroomKey                 contextKey = "joined-classroom"
	joinedClassroomAssignmentKey       contextKey = "joined-classroom-assignment"
	joinedClassroomTeamKey             contextKey = "joined-classroom-team"
)

type FiberContext struct {
	*fiber.Ctx
}

func Get(c *fiber.Ctx) *FiberContext {
	return &FiberContext{Ctx: c}
}

func (c *FiberContext) GetUserID() int {
	return c.Locals(userIDKey).(int)
}

func (c *FiberContext) SetUserID(userID int) {
	c.Locals(userIDKey, userID)
}

func (c *FiberContext) SetGitlabRepository(repo gitlab.Repository) {
	c.Locals(gitlabRepoKey, repo)
}

func (c *FiberContext) GetGitlabRepository() gitlab.Repository {
	value, ok := c.Locals(gitlabRepoKey).(gitlab.Repository)
	if !ok {
		return nil
	}
	return value
}

func (c *FiberContext) SetGitlabUserID(userID int) {
	c.Locals(gitlabUserIDKey, userID)
}

func (c *FiberContext) GetGitlabUserID() int {
	return c.Locals(gitlabUserIDKey).(int)
}

func (c *FiberContext) SetGitlabGroupID(groupID int) {
	c.Locals(gitlabGroupIDKey, groupID)
}

func (c *FiberContext) GetGitlabGroupID() int {
	return c.Locals(gitlabGroupIDKey).(int)
}

func (c *FiberContext) SetGitlabProjectID(projectID int) {
	c.Locals(gitlabProjectIDKey, projectID)
}

func (c *FiberContext) GetGitlabProjectID() int {
	return c.Locals(gitlabProjectIDKey).(int)
}

func (c *FiberContext) SetOwnedClassroom(classroom *database.Classroom) {
	c.Locals(ownedClassroomKey, classroom)
}

func (c *FiberContext) GetOwnedClassroom() *database.Classroom {
	return c.Locals(ownedClassroomKey).(*database.Classroom)
}

func (c *FiberContext) GetOwnedClassroomAssignment() *database.Assignment {
	return c.Locals(ownedClassroomAssignmentKey).(*database.Assignment)
}

func (c *FiberContext) SetOwnedClassroomAssignment(assignment *database.Assignment) {
	c.Locals(ownedClassroomAssignmentKey, assignment)
}

func (c *FiberContext) GetOwnedClassroomAssignmentProject() *database.AssignmentProjects {
	return c.Locals(ownedClassroomAssignmentProjectKey).(*database.AssignmentProjects)
}

func (c *FiberContext) SetOwnedClassroomAssignmentProject(assignmentProject *database.AssignmentProjects) {
	c.Locals(ownedClassroomAssignmentProjectKey, assignmentProject)
}

func (c *FiberContext) SetJoinedClassroom(classroom *database.UserClassrooms) {
	c.Locals(joinedClassroomKey, classroom)
}

func (c *FiberContext) GetJoinedClassroom() *database.UserClassrooms {
	return c.Locals(joinedClassroomKey).(*database.UserClassrooms)
}

func (c *FiberContext) GetJoinedClassroomAssignment() *database.AssignmentProjects {
	return c.Locals(joinedClassroomAssignmentKey).(*database.AssignmentProjects)
}

func (c *FiberContext) SetJoinedClassroomAssignment(assignment *database.AssignmentProjects) {
	c.Locals(joinedClassroomAssignmentKey, assignment)
}

func (c *FiberContext) GetJoinedClassroomTeam() *database.Team {
	return c.Locals(joinedClassroomTeamKey).(*database.Team)
}

func (c *FiberContext) SetJoinedClassroomTeam(team *database.Team) {
	c.Locals(joinedClassroomTeamKey, team)
}

func (c *FiberContext) GetOwnedClassroomMember() *database.UserClassrooms {
	return c.Locals(ownedClassroomMemberKey).(*database.UserClassrooms)
}

func (c *FiberContext) SetOwnedClassroomMember(member *database.UserClassrooms) {
	c.Locals(ownedClassroomMemberKey, member)
}

func (c *FiberContext) GetOwnedClassroomTeam() *database.Team {
	return c.Locals(ownedClassroomTeamKey).(*database.Team)
}

func (c *FiberContext) SetOwnedClassroomTeam(team *database.Team) {
	c.Locals(ownedClassroomTeamKey, team)
}

func (c *FiberContext) GetOwnedClassroomTeamMember() *database.UserClassrooms {
	return c.Locals(ownedClassroomTeamMemberKey).(*database.UserClassrooms)
}

func (c *FiberContext) SetOwnedClassroomTeamMember(member *database.UserClassrooms) {
	c.Locals(ownedClassroomTeamMemberKey, member)
}
