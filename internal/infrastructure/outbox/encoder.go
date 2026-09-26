package outbox

import "github.com/codememory1/d8r/pkg/ddd"

type Encoder interface {
	Encode(event ddd.Event) ([]byte, error)
}
