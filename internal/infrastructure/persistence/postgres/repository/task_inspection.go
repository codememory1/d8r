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

// taskInspectionModel represents the database model of a task inspection.
type taskInspectionModel struct {
	ID             string     `db:"id"`
	TaskID         string     `db:"task_id"`
	EffectiveURL   string     `db:"effective_url"`
	ContentType    *string    `db:"content_type"`
	Filename       *string    `db:"filename"`
	Size           *int64     `db:"size"`
	Strategy       string     `db:"strategy"`
	ETag           *string    `db:"etag"`
	LastModifiedAt *time.Time `db:"last_modified_at"`
	CreatedAt      time.Time  `db:"created_at"`
}

// TaskInspectionRepository provides access to task inspection persistence.
type TaskInspectionRepository struct {
	connection *postgres.ConnectionPool
}

// NewTaskInspectionRepository creates a new task inspection repository.
func NewTaskInspectionRepository(connection *postgres.ConnectionPool) *TaskInspectionRepository {
	return &TaskInspectionRepository{
		connection: connection,
	}
}

// GetLastByTaskID returns the latest inspection associated with the specified task.
func (r *TaskInspectionRepository) GetLastByTaskID(ctx context.Context, id valueobject.ID) (*entity.TaskInspection, error) {
	sql, args, err := sq.
		Select("*").
		From("task_inspections").
		Where(sq.Eq{
			"task_id": id.String(),
		}).
		OrderBy("created_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var model taskInspectionModel

	if pgxscan.Get(ctx, r.connection, &model, sql, args...) != nil {
		return nil, repository.ErrNotFound
	}

	return model.toDomainEntity()
}

// Save persists a task inspection.
func (r *TaskInspectionRepository) Save(ctx context.Context, entity *entity.TaskInspection) error {
	sql, args, err := sq.
		Insert("task_inspections").
		SetMap(map[string]any{
			"id":               entity.ID().String(),
			"task_id":          entity.TaskID().String(),
			"effective_url":    entity.EffectiveURL().String(),
			"content_type":     optional.Map(entity.ContentType(), valueobject.ContentType.String),
			"filename":         optional.Map(entity.Filename(), valueobject.Filename.String),
			"size":             optional.Map(entity.Size(), valueobject.ByteSize.Int64),
			"strategy":         entity.Strategy().String(),
			"etag":             entity.ETag(),
			"last_modified_at": entity.LastModifiedAt(),
			"created_at":       entity.CreatedAt(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if _, err := r.connection.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

// toDomainEntity restores a task inspection domain entity from its database model.
func (m taskInspectionModel) toDomainEntity() (*entity.TaskInspection, error) {
	id, err := valueobject.ParseID(m.ID)

	if err != nil {
		return nil, err
	}

	taskID, err := valueobject.ParseID(m.TaskID)

	if err != nil {
		return nil, err
	}

	effectiveURL, err := valueobject.NewURL(m.EffectiveURL)

	if err != nil {
		return nil, err
	}

	var contentType *valueobject.ContentType
	var filename *valueobject.Filename
	var size *valueobject.ByteSize

	if m.ContentType != nil {
		v, err := valueobject.NewContentType(*m.ContentType)

		if err != nil {
			return nil, err
		}

		contentType = &v
	}

	if m.Filename != nil {
		v, err := valueobject.NewFilename(*m.Filename)

		if err != nil {
			return nil, err
		}

		filename = &v
	}

	if m.Size != nil {
		v, err := valueobject.NewByteSize(*m.Size)

		if err != nil {
			return nil, err
		}

		size = &v
	}

	strategy, err := valueobject.NewDownloadStrategy(m.Strategy)

	if err != nil {
		return nil, err
	}

	return entity.UnmarshalTaskInspection(
		id,
		taskID,
		effectiveURL,
		contentType,
		filename,
		size,
		strategy,
		m.ETag,
		m.LastModifiedAt,
		m.CreatedAt,
	), nil
}
