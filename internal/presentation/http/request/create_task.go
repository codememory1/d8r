package request

import (
	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/restful"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ restful.Input = (*CreateTask)(nil)

// CreateTask represents the request body for creating a new task.
type CreateTask struct {
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers"`
	Filename *string           `json:"filename,omitempty"`
	Priority *int              `json:"priority,omitempty"`
}

// Validate validates the task creation request fields.
func (r CreateTask) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.URL, validation.Required.Error("URL is a required field.")),
		validation.Field(
			&r.Filename,
			validation.Skip.When(r.Filename == nil),
			validation.Required.Error("Filename is a required field."),
		),
		validation.Field(
			&r.Priority,
			validation.Skip.When(r.Priority == nil),
			validation.Required.Error("Priority is a required field."),
			validation.Min(-256).Error("The priority cannot be less than -256"),
			validation.Max(256).Error("The priority cannot be greater than 256"),
		),
	)
}

// ToCommand converts the request into an application command.
func (r CreateTask) ToCommand() (command.CreateTask, error) {
	url, err := valueobject.NewURL(r.URL)

	if err != nil {
		return command.CreateTask{}, err
	}

	headers, err := valueobject.NewHeaders(r.Headers)

	if err != nil {
		return command.CreateTask{}, err
	}

	var filename *valueobject.Filename

	if r.Filename != nil {
		vo, err := valueobject.NewFilename(*r.Filename)

		if err != nil {
			return command.CreateTask{}, err
		}

		filename = &vo
	}

	// Use the default priority unless an explicit value is provided.
	priority := valueobject.DefaultPriority()

	if r.Priority != nil {
		vo, err := valueobject.NewPriority(*r.Priority)

		if err != nil {
			return command.CreateTask{}, err
		}

		priority = vo
	}

	// Build the application command using validated domain values.
	return command.CreateTask{
		URL:      url,
		Headers:  headers,
		Filename: filename,
		Priority: priority,
	}, nil
}
