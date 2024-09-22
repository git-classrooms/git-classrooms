package utils

import "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"

// ProjectAccessLevelCache is a cache for project access levels.
type ProjectAccessLevelCache struct {
	UserID      int
	ProjectID   int
	AccessLevel model.AccessLevelValue
}
