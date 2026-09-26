package outbox

import "github.com/codememory1/d8r/pkg/ddd"

type EventRegistry interface {
	Get(eventType ddd.EventType) (ddd.Event, error)
}
