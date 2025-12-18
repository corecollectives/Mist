CREATE TABLE IF NOT EXISTS system_info (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO system_info (key, value) VALUES ('version', '1.0.0');
INSERT OR IGNORE INTO system_info (key, value) VALUES ('build_date', datetime('now'));

-- Update history table
CREATE TABLE IF NOT EXISTS update_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_version TEXT NOT NULL,
    to_version TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'downloading', 'building', 'installing', 'success', 'failed', 'rolled_back')) DEFAULT 'pending',
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    error_message TEXT,
    rollback_available BOOLEAN DEFAULT 0,
    initiated_by INTEGER,
    FOREIGN KEY (initiated_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_update_history_status ON update_history(status);
CREATE INDEX IF NOT EXISTS idx_update_history_started_at ON update_history(started_at DESC);
