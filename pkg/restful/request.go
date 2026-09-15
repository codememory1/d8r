package restful

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-playground/form/v4"
)

type Input interface {
	validation.Validatable
}

var queryDecoder = form.NewDecoder()

func init() {
	queryDecoder.SetTagName("query")
}

// DecodeBody decodes the request body into the specified input type
// and validates the decoded data.
func DecodeBody[T Input](r io.Reader) (T, error) {
	var input T

	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return input, NewError(http.StatusBadRequest, "invalid request", err)
	}

	return validate(input)
}

func DecodeQuery[T Input](values url.Values) (T, error) {
	var input T

	if err := queryDecoder.Decode(&input, values); err != nil {
		return input, NewError(http.StatusBadRequest, "invalid request", err)
	}

	return validate(input)
}

func validate[T Input](input T) (T, error) {
	if err := input.Validate(); err != nil {
		if validationErrors, ok := errors.AsType[validation.Errors](err); ok {
			for _, validationError := range validationErrors {
				return input, NewError(http.StatusUnprocessableEntity, validationError.Error(), err)
			}
		}

		return input, err
	}

	return input, nil
}
