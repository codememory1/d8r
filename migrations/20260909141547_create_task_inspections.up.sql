-- 20260909141547_create_task_inspections.up.sql

CREATE TABLE task_inspections (
    id UUID PRIMARY KEY,
    task_id uuid NOT NULL,
    effective_url VARCHAR(1024) NOT NULL,
    content_type TEXT NULL,
    filename VARCHAR(255) NULL,
    size BIGINT NULL,
    strategy VARCHAR(16) NOT NULL,
    etag VARCHAR(1024) NULL,
    last_modified_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT task_inspections_task_fk
        FOREIGN KEY (task_id)
        REFERENCES tasks (id)
        ON DELETE CASCADE
);