-- -- 20261002195428_add_webhook_deliveries_sequence_key.down.sql

ALTER TABLE webhook_deliveries DROP COLUMN sequence_key;