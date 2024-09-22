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

	userClassroomKey     contextKey = "user-classroom"
	assignmentKey        contextKey = "assignment"
	assignmentProjectKey contextKey = "assignment-project"
	classroomMember      contextKey = "classroom-member"
	teamKey              contextKey = "team"
)

// FiberContext wraps the fiber.Ctx to provide additional methods.
type FiberContext struct {
	*fiber.Ctx
}

// Get returns a FiberContext from the fiber.Ctx.
func Get(c *fiber.Ctx) *FiberContext {
	return &FiberContext{Ctx: c}
}

// GetUserID returns the user ID from the context.
func (c *FiberContext) GetUserID() int {
	return c.Locals(userIDKey).(int)
}

// SetUserID sets the user ID in the context.
func (c *FiberContext) SetUserID(userID int) {
	c.Locals(userIDKey, userID)
}

// SetGitlabRepository sets the GitLab repository in the context.
func (c *FiberContext) SetGitlabRepository(repo gitlab.Repository) {
	c.Locals(gitlabRepoKey, repo)
}

// GetGitlabRepository returns the GitLab repository from the context.
func (c *FiberContext) GetGitlabRepository() gitlab.Repository {
	value, ok := c.Locals(gitlabRepoKey).(gitlab.Repository)
	if !ok {
		return nil
	}
	return value
}

// SetGitlabUserID sets the GitLab user ID in the context.
func (c *FiberContext) SetGitlabUserID(userID int) {
	c.Locals(gitlabUserIDKey, userID)
}

// GetGitlabUserID returns the GitLab user ID from the context.
func (c *FiberContext) GetGitlabUserID() int {
	return c.Locals(gitlabUserIDKey).(int)
}

// SetGitlabGroupID sets the GitLab group ID in the context.
func (c *FiberContext) SetGitlabGroupID(groupID int) {
	c.Locals(gitlabGroupIDKey, groupID)
}

// GetGitlabGroupID returns the GitLab group ID from the context.
func (c *FiberContext) GetGitlabGroupID() int {
	return c.Locals(gitlabGroupIDKey).(int)
}

// SetGitlabProjectID sets the GitLab project ID in the context.
func (c *FiberContext) SetGitlabProjectID(projectID int) {
	c.Locals(gitlabProjectIDKey, projectID)
}

// GetGitlabProjectID returns the GitLab project ID from the context.
func (c *FiberContext) GetGitlabProjectID() int {
	return c.Locals(gitlabProjectIDKey).(int)
}

// GetUserClassroom returns the user classroom from the context.
func (c *FiberContext) GetUserClassroom() *database.UserClassrooms {
	return c.Locals(userClassroomKey).(*database.UserClassrooms)
}

// SetUserClassroom sets the user classroom in the context.
func (c *FiberContext) SetUserClassroom(classroom *database.UserClassrooms) {
	c.Locals(userClassroomKey, classroom)
}

// GetAssignment returns the assignment from the context.
func (c *FiberContext) GetAssignment() *database.Assignment {
	return c.Locals(assignmentKey).(*database.Assignment)
}

// SetAssignment sets the assignment in the context.
func (c *FiberContext) SetAssignment(assignment *database.Assignment) {
	c.Locals(assignmentKey, assignment)
}

// GetAssignmentProject returns the assignment project from the context.
func (c *FiberContext) GetAssignmentProject() *database.AssignmentProjects {
	return c.Locals(assignmentProjectKey).(*database.AssignmentProjects)
}

// SetAssignmentProject sets the assignment project in the context.
func (c *FiberContext) SetAssignmentProject(project *database.AssignmentProjects) {
	c.Locals(assignmentProjectKey, project)
}

// GetClassroomMember returns the user classroom from the context.
func (c *FiberContext) GetClassroomMember() *database.UserClassrooms {
	return c.Locals(classroomMember).(*database.UserClassrooms)
}

// SetClassroomMember sets the user classroom in the context.
func (c *FiberContext) SetClassroomMember(member *database.UserClassrooms) {
	c.Locals(classroomMember, member)
}

// GetTeam returns the team from the context.
func (c *FiberContext) GetTeam() *database.Team {
	return c.Locals(teamKey).(*database.Team)
}

// SetTeam sets the team in the context.
func (c *FiberContext) SetTeam(team *database.Team) {
	c.Locals(teamKey, team)
}
