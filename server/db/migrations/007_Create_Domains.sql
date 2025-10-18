CREATE TABLE IF NOT EXISTS domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER NOT NULL,
    domain_name TEXT NOT NULL UNIQUE,
    ssl_status TEXT NOT NULL CHECK(ssl_status IN ('pending', 'active', 'expired')) DEFAULT 'pending', 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
)
