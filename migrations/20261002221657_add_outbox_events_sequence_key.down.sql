-- 20261002221657_add_outbox_events_sequence_key.down.sql

ALTER TABLE outbox_events DROP COLUMN sequence_key;