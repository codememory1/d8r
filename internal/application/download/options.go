package download

import (
	"github.com/codememory1/d8r/internal/application/inspect"
	"github.com/codememory1/d8r/internal/domain/valueobject"
)

type Options struct {
	InspectionResult inspect.Result
	Headers          map[string]string
	Filename         *valueobject.Filename
}
