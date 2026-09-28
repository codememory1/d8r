-- 20260928141648_add_webhook_deliveries_version.up.sql

ALTER TABLE webhook_deliveries ADD COLUMN version BIGINT NOT NULL DEFAULT 1;