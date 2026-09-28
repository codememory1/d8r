-- 20260926165947_create_webhook_deliveries.down.sql

DROP INDEX IF EXISTS webhook_deliveries_pending_idx;
DROP TABLE IF EXISTS webhook_deliveries;