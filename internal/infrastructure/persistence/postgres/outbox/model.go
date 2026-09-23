package outbox

import (
	infraoutbox "github.com/codememory1/d8r/internal/infrastructure/outbox"
	"github.com/codememory1/d8r/pkg/ddd"
)

type messageRow struct {
	ID        string `db:"id"`
	EventType string `db:"event_type"`
	Payload   []byte `db:"payload"`
	Version   int64  `db:"version"`
}

func (r messageRow) toMessage() infraoutbox.Message {
	return infraoutbox.Message{
		ID:        r.ID,
		EventType: ddd.EventType(r.EventType),
		Payload:   r.Payload,
		Version:   r.Version,
	}
}
