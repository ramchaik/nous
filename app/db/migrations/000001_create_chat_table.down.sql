DROP TRIGGER IF EXISTS update_chat_timestamp;
DROP TRIGGER IF EXISTS update_session_timestamp;
ALTER TABLE chat DROP COLUMN session_id;
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS chat;