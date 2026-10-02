-- 20261002221657_add_outbox_events_sequence_key.up.sql

ALTER TABLE outbox_events ADD COLUMN sequence_key VARCHAR(255) NULL;