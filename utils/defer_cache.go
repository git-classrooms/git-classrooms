package utils

import "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"

type ProjectAccessLevelCache struct {
	UserID      int
	ProjectID   int
	AccessLevel model.AccessLevelValue
}
