package go_gitlab_repo

import (
	"backend/model"

	"github.com/xanzy/go-gitlab"
)

func ProjectFromGoGitlab(gitlabProject gitlab.Project) *model.Project {
	var owner *model.User = nil
	if gitlabProject.Owner != nil {
		owner = UserFromGoGitlab(*gitlabProject.Owner)
	}

	return &model.Project{
		Name:        gitlabProject.Name,
		ID:          gitlabProject.ID,
		Visibility:  VisibilityFromGoGitlab(gitlabProject.Visibility),
		WebUrl:      gitlabProject.WebURL,
		Description: gitlabProject.Description,
		Owner:       owner,
	}
}

func ProjectFromGoGitlabWithProjectMembers(gitlabProject gitlab.Project, gitlabMembers []*gitlab.ProjectMember) *model.Project {
	var owner *model.User = nil
	if gitlabProject.Owner != nil {
		owner = UserFromGoGitlab(*gitlabProject.Owner)
	}

	members := make([]model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		members[i] = *UserFromGoGitlabProjectMember(*gitlabMember)
	}

	return &model.Project{
		Name:        gitlabProject.Name,
		ID:          gitlabProject.ID,
		Visibility:  VisibilityFromGoGitlab(gitlabProject.Visibility),
		WebUrl:      gitlabProject.WebURL,
		Description: gitlabProject.Description,
		Owner:       owner,
		Member:      members,
	}
}

func ProjectFromGoGitlabWithGroupMembers(gitlabProject gitlab.Project, gitlabMembers []*gitlab.GroupMember) *model.Project {
	var owner *model.User = nil
	if gitlabProject.Owner != nil {
		owner = UserFromGoGitlab(*gitlabProject.Owner)
	}

	members := make([]model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		members[i] = *UserFromGoGitlabGroupMember(*gitlabMember)
	}

	return &model.Project{
		Name:        gitlabProject.Name,
		ID:          gitlabProject.ID,
		Visibility:  VisibilityFromGoGitlab(gitlabProject.Visibility),
		WebUrl:      gitlabProject.WebURL,
		Description: gitlabProject.Description,
		Owner:       owner,
		Member:      members,
	}
}

func VisibilityFromGoGitlab(input gitlab.VisibilityValue) model.Visibility {
	if input == "public" {
		return model.Public
	} else if input == "internal" {
		return model.Internal
	} else if input == "private" {
		return model.Private
	}
	return 0
}

func UserFromGoGitlab(input gitlab.User) *model.User {
	return &model.User{
		ID:       input.ID,
		Username: input.Username,
		Name:     input.Name,
		WebUrl:   input.WebURL,
		Email:    input.Email,
	}
}

func UserFromGoGitlabProjectMember(input gitlab.ProjectMember) *model.User {
	return &model.User{
		ID:       input.ID,
		Username: input.Username,
		Name:     input.Name,
		WebUrl:   input.WebURL,
		Email:    input.Email,
	}
}

func UserFromGoGitlabGroupMember(input gitlab.GroupMember) *model.User {
	return &model.User{
		ID:       input.ID,
		Username: input.Username,
		Name:     input.Name,
		WebUrl:   input.WebURL,
		Email:    input.Email,
	}
}

func GroupFromGoGitlab(input gitlab.Group) *model.Group {
	return &model.Group{
		Name:        input.Name,
		ID:          input.ID,
		Description: input.Description,
		WebUrl:      input.WebURL,
		Visibility:  VisibilityFromGoGitlab(input.Visibility),
	}
}

func GroupFromGoGitlabWithMembersAndProjects(group gitlab.Group, members []model.User, projects []model.Project) *model.Group {

	return &model.Group{
		Name:        group.Name,
		ID:          group.ID,
		Description: group.Description,
		WebUrl:      group.WebURL,
		Visibility:  VisibilityFromGoGitlab(group.Visibility),
		Projects:    projects,
		Member:      members,
	}
}

func ConvertUserPointerSlice(input []*model.User) []model.User {
	output := make([]model.User, len(input))
	for i, ptr := range input {
		output[i] = *ptr
	}
	return output
}

func ConvertProjectPointerSlice(input []*model.Project) []model.Project {
	output := make([]model.Project, len(input))
	for i, ptr := range input {
		output[i] = *ptr
	}
	return output
}
