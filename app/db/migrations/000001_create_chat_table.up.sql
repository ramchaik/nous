-- Create session table
CREATE TABLE IF NOT EXISTS session (
    session_id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create chat table (if not exists)
CREATE TABLE IF NOT EXISTS chat (
    message_id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    text TEXT NOT NULL,
    type TEXT CHECK(type IN ('agent', 'user')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES session(session_id)
);

-- Create triggers for updating timestamps
CREATE TRIGGER IF NOT EXISTS update_chat_timestamp 
AFTER UPDATE ON chat
BEGIN
    UPDATE chat SET updated_at = CURRENT_TIMESTAMP WHERE chat_id = NEW.chat_id;
END;

CREATE TRIGGER IF NOT EXISTS update_session_timestamp 
AFTER UPDATE ON session
BEGIN
    UPDATE session SET updated_at = CURRENT_TIMESTAMP WHERE session_id = NEW.session_id;
END;