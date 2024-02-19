package context

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
)

type contextKey string

const (
	gitlabRepoKey                      contextKey = "gitlab-repo"
	classroomKey                       contextKey = "classroom"
	ownedClassroomKey                  contextKey = "owned-classroom"
	ownedClassroomAssignmentKey        contextKey = "owned-classroom-assignment"
	ownedClassroomAssignmentProjectKey contextKey = "owned-classroom-assignment-project"
	joinedClassroomKey                 contextKey = "joined-classroom"
	joinedClassroomAssignmentKey       contextKey = "joined-classroom-assignment"
	userIDKey                          contextKey = "user-id"
)

type FiberContext struct {
	*fiber.Ctx
}

func Get(c *fiber.Ctx) *FiberContext {
	return &FiberContext{Ctx: c}
}

func (c *FiberContext) GetGitlabRepository() gitlab.Repository {
	value, ok := c.Locals(gitlabRepoKey).(gitlab.Repository)
	if !ok {
		return nil
	}
	return value
}

func (c *FiberContext) SetGitlabRepository(repo gitlab.Repository) {
	c.Locals(gitlabRepoKey, repo)
}

func (c *FiberContext) SetClassroom(classroom *database.UserClassrooms) {
	c.Locals(classroomKey, classroom)
}

func (c *FiberContext) GetClassroom() *database.UserClassrooms {
	value, ok := c.Locals(classroomKey).(*database.UserClassrooms)
	if !ok {
		return nil
	}
	return value
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

func (c *FiberContext) GetUserID() int {
	return c.Locals(userIDKey).(int)
}

func (c *FiberContext) SetUserID(userID int) {
	c.Locals(userIDKey, userID)
}
