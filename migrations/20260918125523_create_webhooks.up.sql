-- 20260918125523_create_webhooks.up.sql

CREATE TABLE webhooks (
    id UUID PRIMARY KEY,
    url VARCHAR(1024) NOT NULL,
    headers JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);

CREATE TABLE webhook_events (
    id UUID PRIMARY KEY,
    webhook_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT webhook_events_webhook_fk
        FOREIGN KEY (webhook_id)
        REFERENCES webhooks (id)
        ON DELETE CASCADE
);