package utils

import (
	"fmt"
	"net/url"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func CreateClassroomGitlabDescription(classroom *database.Classroom, publicURL *url.URL) string {
	return fmt.Sprintf("%s\n\n\n__Managed by [GitClassrooms](%s/classrooms/%s)__", classroom.Description, publicURL, classroom.ID.String())
}

func CreateTeamGitlabDescription(classroom *database.Classroom, team *database.Team, publicURL *url.URL) string {
	return fmt.Sprintf("Team of %s\n\n\n__Managed by [GitClassrooms](%s/classrooms/%s/teams/%s)__", classroom.Description, publicURL, classroom.ID.String(), team.ID.String())
}
