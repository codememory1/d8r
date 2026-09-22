package reader

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/codememory1/d8r/internal/application/query"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/codememory1/d8r/pkg/pagination"
	"github.com/georgysavva/scany/v2/pgxscan"
)

var _ query.TaskReader = (*TaskReader)(nil)

// taskRow represents a task row returned by PostgreSQL.
type taskRow struct {
	ID        string            `db:"id"`
	URL       string            `db:"url"`
	Headers   map[string]string `db:"headers"`
	Filename  *string           `db:"filename"`
	Priority  int               `db:"priority"`
	Status    string            `db:"status"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt *time.Time        `db:"updated_at"`
}

type TaskReader struct {
	connection *postgres.ConnectionPool
}

func NewTaskReader(connection *postgres.ConnectionPool) *TaskReader {
	return &TaskReader{
		connection: connection,
	}
}

func (r *TaskReader) GetAllPaginated(ctx context.Context, cursor *pagination.Cursor, limit int) ([]query.TaskReadModel, error) {
	// Build a deterministic ordering that matches the cursor fields.
	builder := sq.
		Select("*").
		From("tasks").
		Limit(uint64(limit)).
		OrderBy("created_at DESC", "id DESC").
		PlaceholderFormat(sq.Dollar)

	// Continue from the item represented by the cursor when provided.
	if cursor != nil {
		builder = builder.Where(sq.Expr(
			"(created_at, id) < (?, ?)",
			time.UnixMicro(cursor.Timestamp),
			cursor.LastID,
		))
	}

	sql, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	var rows []taskRow

	// Load persistence models from PostgreSQL.
	if err := pgxscan.Select(ctx, r.connection, &rows, sql, args...); err != nil {
		return nil, err
	}

	// Rehydrate domain entities from persistence models.
	readModels := make([]query.TaskReadModel, len(rows))

	for i, row := range rows {
		readModels[i] = row.toReadModel()
	}

	return readModels, nil
}

func (r *TaskReader) GetByID(ctx context.Context, id valueobject.ID) (query.TaskReadModel, error) {
	sql, args, err := sq.
		Select("*").
		From("tasks").
		Where(sq.Eq{
			"id": id.String(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return query.TaskReadModel{}, err
	}

	var row taskRow

	// Load the persistence model matching the requested identifier.
	if pgxscan.Get(ctx, r.connection, &row, sql, args...) != nil {
		return query.TaskReadModel{}, repository.ErrNotFound
	}

	return row.toReadModel(), nil
}

// toReadModel converts a database row into an application read model.
func (r taskRow) toReadModel() query.TaskReadModel {
	return query.TaskReadModel{
		ID:        r.ID,
		URL:       r.URL,
		Headers:   r.Headers,
		Filename:  r.Filename,
		Priority:  r.Priority,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
