package worker

import (
	"context"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/logging"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
)

// GetWorkerRepo logs into the GitLab repository using the provided group access token.
// It returns a GitLab repository object and an error, if any.
func GetWorkerRepo(ctx context.Context, gitlabConfig gitlabConfig.Config, groupAccessToken string) (gitlab.Repository, error) {
	repo := gitlab.NewGitlabRepo(gitlabConfig)
	err := repo.GroupAccessLogin(groupAccessToken, logging.GetLogger(ctx))
	if err != nil {
		return nil, err
	}

	return repo, nil
}
