-- 20260908172155_create_tasks.up.sql

CREATE TABLE tasks (
    ID UUID PRIMARY KEY,
    url VARCHAR(1024) NOT NULL,
    headers JSONB NOT NULL DEFAULT '{}'::jsonb,
    filename VARCHAR(255) NULL,
    priority INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);