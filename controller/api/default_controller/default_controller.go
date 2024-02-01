package default_controller

import (
	mailRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
)

type DefaultController struct {
	mailRepo mailRepo.Repository
}

func NewApiController(mailRepo mailRepo.Repository) *DefaultController {
	return &DefaultController{mailRepo: mailRepo}
}
