package outbox

import "github.com/codememory1/d8r/pkg/ddd"

type Decoder interface {
	Decode(payload []byte, target ddd.Event) error
}
