ALTER TABLE deployments ADD COLUMN stage TEXT DEFAULT 'pending';
ALTER TABLE deployments ADD COLUMN progress INTEGER DEFAULT 0;
ALTER TABLE deployments ADD COLUMN error_message TEXT;
ALTER TABLE deployments ADD COLUMN started_at DATETIME;
ALTER TABLE deployments ADD COLUMN duration INTEGER;

