package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
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

func (r *GitlabRepo) CreatePersonalAccessTokenForRoot(ctx context.Context) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Jar: jar,
	}

	log.Println("Loging in to Gitlab with Root User")
	err = r.loginGitlab(ctx, client)
	if err != nil {
		return "", err
	}
	log.Println("Creating PAT with Gitlab Session")
	return r.createPersonalAccessToken(ctx, client)
}

var csrfParamRegex = regexp.MustCompile(`<meta name="csrf-param" content="(.*)" />`)
var csrfTokenRegex = regexp.MustCompile(`<meta name="csrf-token" content="(.*)" />`)

type CSRF struct {
	param string
	token string
}

func (r *GitlabRepo) parseCsrfFromBody(rc io.ReadCloser) (CSRF, error) {
	log.Println("Getting CSRF Token from Gitlab")
	defer rc.Close()
	scanner := bufio.NewScanner(rc)
	var paramFound, tokenFound bool
	var csrf CSRF
	for scanner.Scan() {
		line := scanner.Text()

		if paramFound && tokenFound {
			break
		}

		if !paramFound {
			paramMatch := csrfParamRegex.FindStringSubmatch(line)
			if len(paramMatch) > 0 {
				csrf.param = paramMatch[1]
				paramFound = true
				continue
			}
		}
		if !tokenFound {
			tokenMatch := csrfTokenRegex.FindStringSubmatch(line)
			if len(tokenMatch) > 0 {
				csrf.token = tokenMatch[1]
				tokenFound = true
				continue
			}
		}
	}

	if !paramFound || !tokenFound {
		return csrf, errors.New("param or token not found")
	}
	return csrf, nil
}

func (r *GitlabRepo) loginGitlab(ctx context.Context, client *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/sigin_in", r.url), nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("user sign_in not ok")
	}

	csrf, err := r.parseCsrfFromBody(res.Body)
	if err != nil {
		return err
	}

	params := make(url.Values)
	params.Add(csrf.param, csrf.token)
	params.Add("user[login]", "root")
	params.Add("user[password]", "geheim1234")
	params.Add("user[remember_me]", "1")

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/users/sign_in", r.url), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	res, err = client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (r *GitlabRepo) createPersonalAccessToken(ctx context.Context, client *http.Client) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/-/user_settings/personal_access_tokens", r.url), nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New("user_settings personal_access_tokens not ok")
	}

	csrf, err := r.parseCsrfFromBody(res.Body)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err = json.NewEncoder(&body).Encode(map[string]any{
		"name":        "root",
		"description": "",
		"expires_at":  time.Now().AddDate(0, 1, 0).Format(time.DateOnly),
		"scopes":      []string{"api", "admin_mode"},
	}); err != nil {
		return "", err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/-/user_settings/personal_access_tokens", r.url), &body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(body.Len()))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CSRF-Token", csrf.token)

	res, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var resBody struct {
		Token string `json:"token"`
	}
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return "", err
	}

	return resBody.Token, nil
}
