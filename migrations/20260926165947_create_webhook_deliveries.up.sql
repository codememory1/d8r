-- 20260926165947_create_webhook_deliveries.up.sql

CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY,
    webhook_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(32) NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    response_status INTEGER NULL,
    last_error TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    delivered_at TIMESTAMPTZ NULL,

    CONSTRAINT webhook_deliveries_webhook_fk
        FOREIGN KEY (webhook_id)
        REFERENCES webhooks (id)
        ON DELETE CASCADE
);

CREATE INDEX webhook_deliveries_pending_idx
    ON webhook_deliveries (
        next_attempt_at ASC,
        created_at ASC,
        id ASC
    )
    WHERE status IN ('pending', 'failed')