package utils

import (
	"fmt"
	"net/url"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

// CreateClassroomGitlabDescription creates a GitLab description for a classroom.
func CreateClassroomGitlabDescription(classroom *database.Classroom, publicURL *url.URL) string {
	return fmt.Sprintf("%s\n\n\n__Managed by [GitClassrooms](%s/classrooms/%s)__", classroom.Description, publicURL, classroom.ID.String())
}

// CreateTeamGitlabDescription creates a GitLab description for a team.
func CreateTeamGitlabDescription(classroom *database.Classroom, team *database.Team, publicURL *url.URL) string {
	return fmt.Sprintf("Team of %s\n\n\n__Managed by [GitClassrooms](%s/classrooms/%s/teams/%s)__", classroom.Name, publicURL, classroom.ID.String(), team.ID.String())
}
