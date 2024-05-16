package default_controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type updateAssignmentRequest struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

func (ctrl *DefaultController) PutOwnedAssignments(c *fiber.Ctx) error {
	return nil
}
