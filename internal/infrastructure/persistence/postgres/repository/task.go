package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/codememory1/d8r/pkg/optional"
	"github.com/georgysavva/scany/v2/pgxscan"
)

const taskClaimByStatusSQL = `
	WITH filtered_tasks AS (
		SELECT
			id
		FROM tasks
		WHERE status = $1
		ORDER BY 
			priority DESC, 
			id ASC
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	)
	UPDATE tasks
	SET status = $3,
		version = version + 1,
		updated_at = NOW()
	FROM filtered_tasks
	WHERE tasks.id = filtered_tasks.id
	RETURNING tasks.*
`

// taskModel represents the persistence model of a task stored in PostgreSQL.
type taskModel struct {
	ID        string            `db:"id"`
	URL       string            `db:"url"`
	Headers   map[string]string `db:"headers"`
	Filename  *string           `db:"filename"`
	Priority  int               `db:"priority"`
	Status    string            `db:"status"`
	Version   int64             `db:"version"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt *time.Time        `db:"updated_at"`
}

// TaskRepository provides PostgreSQL persistence operations for task entities.
type TaskRepository struct {
	connection *postgres.ConnectionPool
}

// NewTaskRepository creates a PostgreSQL-backed task repository.
func NewTaskRepository(pool *postgres.ConnectionPool) *TaskRepository {
	return &TaskRepository{connection: pool}
}

// GetByID returns a task by its identifier.
func (r *TaskRepository) GetByID(ctx context.Context, id valueobject.ID) (*entity.Task, error) {
	sql, args, err := sq.
		Select("*").
		From("tasks").
		Where(sq.Eq{
			"id": id.String(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var model taskModel

	// Load the persistence model matching the requested identifier.
	if pgxscan.Get(ctx, r.connection, &model, sql, args...) != nil {
		return nil, repository.ErrNotFound
	}

	return model.toDomainEntity()
}

// ClaimReadyToDownload atomically claims tasks that are ready to be downloaded
// and transitions them to the downloading state.
func (r *TaskRepository) ClaimReadyToDownload(ctx context.Context, limit int) ([]*entity.Task, error) {
	var models []taskModel

	err := pgxscan.Select(ctx, r.connection, &models, taskClaimByStatusSQL, []any{
		entity.TaskStatusReadyToDownload,
		limit,
		entity.TaskStatusDownloading,
	}...)

	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Task, len(models))

	for i, model := range models {
		e, err := model.toDomainEntity()

		if err != nil {
			return nil, err
		}

		entities[i] = e
	}

	return entities, nil
}

// ClaimPending atomically claims pending tasks for inspection
// and transitions them to the inspecting state.
func (r *TaskRepository) ClaimPending(ctx context.Context, limit int) ([]*entity.Task, error) {
	var models []taskModel

	err := pgxscan.Select(ctx, r.connection, &models, taskClaimByStatusSQL, []any{
		entity.TaskStatusPending,
		limit,
		entity.TaskStatusInspecting,
	}...)

	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Task, len(models))

	for i, model := range models {
		e, err := model.toDomainEntity()

		if err != nil {
			return nil, err
		}

		entities[i] = e
	}

	return entities, nil
}

// Save persists a task entity in PostgreSQL.
func (r *TaskRepository) Save(ctx context.Context, task *entity.Task) error {
	// Map domain values to their persistence representation.
	sql, args, err := sq.
		Insert("tasks").
		SetMap(map[string]any{
			"id":         task.ID().String(),
			"url":        task.URL().String(),
			"headers":    task.Headers().Map(),
			"filename":   optional.Map(task.Filename(), valueobject.Filename.String),
			"priority":   task.Priority().Int(),
			"status":     string(task.Status()),
			"version":    task.Version(),
			"created_at": task.CreatedAt(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	// Persist the task using the repository connection.
	if _, err = r.connection.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

// Update updates the task entity in PostgreSQL using optimistic locking.
func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	sql, args, err := sq.
		Update("tasks").
		Set("url", task.URL().String()).
		Set("headers", task.Headers().Map()).
		Set("filename", optional.Map(task.Filename(), valueobject.Filename.String)).
		Set("priority", task.Priority().Int()).
		Set("status", string(task.Status())).
		Set("version", task.Version()+1).
		Set("updated_at", time.Now()).
		Where(sq.Eq{
			"id":      task.ID().String(),
			"version": task.Version(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	tag, err := r.connection.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrConcurrentModification
	}

	task.IncrementVersion()

	return nil
}

// toDomainEntity reconstructs a task domain entity from its persistence model.
func (m taskModel) toDomainEntity() (*entity.Task, error) {
	// Restore and validate the task identifier.
	id, err := valueobject.ParseID(m.ID)

	if err != nil {
		return nil, err
	}

	// Restore and validate the task URL.
	url, err := valueobject.NewURL(m.URL)

	if err != nil {
		return nil, err
	}

	// Restore and validate request headers.
	headers, err := valueobject.NewHeaders(m.Headers)

	if err != nil {
		return nil, err
	}

	// Restore the optional filename when present.
	var filename *valueobject.Filename

	if m.Filename != nil {
		v, err := valueobject.NewFilename(*m.Filename)

		if err != nil {
			return nil, err
		}

		filename = &v
	}

	// Restore and validate task priority.
	priority, err := valueobject.NewPriority(m.Priority)

	if err != nil {
		return nil, err
	}

	// Rehydrate the entity using values loaded from persistence.
	return entity.UnmarshalTask(
		id,
		url,
		headers,
		filename,
		priority,
		m.Status,
		m.Version,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}
