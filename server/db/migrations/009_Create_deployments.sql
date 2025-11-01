CREATE TABLE IF NOT EXISTS deployments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER NOT NULL,
    commit_hash TEXT NOT NULL,
    commit_message TEXT,
    triggered_by INTEGER,
    logs TEXT,
    container_id TEXT,
    status TEXT NOT NULL CHECK(status IN ('pending', 'building', 'deploying', 'success', 'failed', 'stopped')) DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    FOREIGN KEY (triggered_by) REFERENCES users(id) ON DELETE SET NULL
);
