package outbox

import "github.com/codememory1/d8r/pkg/ddd"

type Decoder interface {
	Decode(eventType ddd.EventType, payload []byte) (ddd.Event, error)
}
