package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

type GitlabRepo struct {
	client *gitlab.Client
	url    string
}

func NewGitlabRepo(gitlabURL string, accessToken string) (*GitlabRepo, error) {
	cli, err := gitlab.NewOAuthClient(accessToken, gitlab.WithBaseURL(gitlabURL))
	if err != nil {
		return nil, err
	}
	return &GitlabRepo{client: cli, url: gitlabURL}, nil
}

func (r *GitlabRepo) Health(ctx context.Context) error {
	res, err := http.Get(fmt.Sprintf("%s/-/health", r.url))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}

func (r *GitlabRepo) CreateUser(ctx context.Context, username, email, name string) (*gitlab.User, error) {
	user, _, err := r.client.Users.CreateUser(&gitlab.CreateUserOptions{
		Email:               &email,
		Username:            &username,
		Name:                &name,
		Password:            Ptr("geheim1234"),
		ResetPassword:       Ptr(false),
		ForceRandomPassword: Ptr(false),
		SkipConfirmation:    Ptr(true),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *GitlabRepo) GetAccessTokenOfUser(ctx context.Context, userID int) (*gitlab.PersonalAccessToken, error) {
	expiresAt := time.Now().AddDate(0, 0, 364)
	token, _, err := r.client.Users.CreatePersonalAccessToken(userID, &gitlab.CreatePersonalAccessTokenOptions{
		Name:      Ptr("root"),
		ExpiresAt: Ptr(gitlab.ISOTime(expiresAt)),
		Scopes:    Ptr([]string{"api", "write_repository"}),
	},
		gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *GitlabRepo) CreateInstanceRunner(ctx context.Context) (*gitlab.UserRunner, error) {
	runner, _, err := r.client.Users.CreateUserRunner(&gitlab.CreateUserRunnerOptions{
		RunnerType:  Ptr("instance_type"),
		RunUntagged: Ptr(true),
	},
		gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return runner, nil
}

func (r *GitlabRepo) CreateApplication(ctx context.Context) (*gitlab.Application, error) {
	application, _, err := r.client.Applications.CreateApplication(&gitlab.CreateApplicationOptions{
		Name:         Ptr("GitClassrooms"),
		Scopes:       Ptr("api openid email profile"),
		Confidential: Ptr(true),
		RedirectURI: Ptr(strings.Join([]string{
			"http://localhost:3000/api/v1/auth/gitlab/callback",
			"http://localhost:5173/api/v1/auth/gitlab/callback",
			"http://localhost:6969/api/v1/auth/gitlab/callback",
		}, "\n")),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (r *GitlabRepo) CreateTemplateProject(ctx context.Context, opts *gitlab.CreateProjectOptions) (*gitlab.Project, error) {
	project, _, err := r.client.Projects.CreateProject(opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *GitlabRepo) CreateClassroom(ctx context.Context, name string) (*gitlab.Group, error) {
	group, _, err := r.client.Groups.CreateGroup(&gitlab.CreateGroupOptions{
		Name:       &name,
		Path:       Ptr(convertToGitLabPath(name)),
		Visibility: Ptr(gitlab.PrivateVisibility),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GitlabRepo) CreateGroupAccessToken(ctx context.Context, groupID int) (*gitlab.GroupAccessToken, error) {
	expiresAt := time.Now().AddDate(0, 0, 364)
	token, _, err := r.client.GroupAccessTokens.CreateGroupAccessToken(groupID, &gitlab.CreateGroupAccessTokenOptions{
		Name:        Ptr("GitClassrooms"),
		Scopes:      Ptr([]string{"api"}),
		ExpiresAt:   Ptr(gitlab.ISOTime(expiresAt)),
		AccessLevel: Ptr(gitlab.OwnerPermissions),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *GitlabRepo) CreateTeam(ctx context.Context, groupID int, name string) (*gitlab.Group, error) {
	group, _, err := r.client.Groups.CreateGroup(&gitlab.CreateGroupOptions{
		ParentID:   &groupID,
		Name:       &name,
		Path:       Ptr(convertToGitLabPath(name)),
		Visibility: Ptr(gitlab.PrivateVisibility),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GitlabRepo) AddUserToGroup(ctx context.Context, groupId, userId int, accessLevel gitlab.AccessLevelValue) error {
	_, _, err := r.client.GroupMembers.AddGroupMember(groupId, &gitlab.AddGroupMemberOptions{
		UserID:      &userId,
		AccessLevel: &accessLevel,
	}, gitlab.WithContext(ctx))
	return err
}

// func (repo *GitlabRepo) AddUserToProject(ctx context.Context, projectId, userId int, accessLevel gitlab.AccessLevelValue) error {
// 	_, _, err := repo.client.ProjectMembers.AddProjectMember(projectId, &gitlab.AddProjectMemberOptions{
// 		UserID:      &userId,
// 		AccessLevel: &accessLevel,
// 	}, gitlab.WithContext(ctx))
// 	return err
// }

func Ptr[T any](v T) *T {
	return &v
}

func convertToGitLabPath(s string) string {
	// Remove unwanted characters
	reg, _ := regexp.Compile("[^a-zA-Z0-9_.-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove leading and trailing special characters
	s = strings.Trim(s, "_.-")

	// Prevent consecutive special characters
	reg, _ = regexp.Compile("[-_.]{2,}")
	s = reg.ReplaceAllString(s, "-")

	// Prevent specific endings
	if strings.HasSuffix(s, ".git") || strings.HasSuffix(s, ".atom") {
		s = s[:len(s)-4]
	}

	// Ensure the path name is at least one character long
	if len(s) == 0 {
		s = "gc_"
	}

	return s
}
