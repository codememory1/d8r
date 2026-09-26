package request

import (
	"github.com/codememory1/d8r/internal/application/query"
	httppagination "github.com/codememory1/d8r/internal/presentation/http/pagination"
	corepagination "github.com/codememory1/d8r/pkg/pagination"
	"github.com/codememory1/d8r/pkg/restful"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// listTasksDefaultLimit defines the default page size
	// when the limit query parameter is omitted.
	listTasksDefaultLimit = 20
)

// ListTasks represents query parameters for retrieving a paginated task list.
type ListTasks struct {
	Cursor *string `query:"cursor"`
	Limit  *int    `query:"limit"`
}

// Validate validates list task query parameters.
func (r ListTasks) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(
			&r.Limit,
			validation.Min(query.ListTasksMinLimit).Error("limit must be greater than zero"),
			validation.Max(query.ListTasksMaxLimit).Error("limit cannot exceed 100"),
		),
	)
}

// ToQuery converts HTTP query parameters into an application query.
func (r ListTasks) ToQuery() (query.ListTasks, error) {
	var cursor *corepagination.Cursor

	// Resolve the page size, falling back to the default when omitted.
	limit := listTasksDefaultLimit

	if r.Limit != nil {
		limit = *r.Limit
	}

	// Decode the opaque cursor when pagination continues from a previous page.
	if r.Cursor != nil {
		decodedCursor, err := httppagination.Decode(*r.Cursor)

		if err != nil {
			return query.ListTasks{}, restful.NewError(400, "invalid cursor", err)
		}

		cursor = &decodedCursor
	}

	// Build the application query with normalized pagination parameters.
	return query.ListTasks{
		Cursor: cursor,
		Limit:  limit,
	}, nil
}
