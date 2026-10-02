-- 20261002195428_add_webhook_deliveries_sequence_key.up.sql

ALTER TABLE webhook_deliveries ADD COLUMN sequence_key VARCHAR(255) NULL;