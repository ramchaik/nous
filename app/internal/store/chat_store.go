package store

import (
	"database/sql"

	"nous/internal/models"

	"github.com/google/uuid"
)

type ChatStore interface {
	GetByChatID(string) (*models.Chat, error)
	CreateSession() (string, error)
	GetSession(sessionID string) (*models.Session, error)
	CreateChat(chat *models.Chat) error
	GetChatsBySession(sessionID string) ([]*models.Chat, error)
}

type SQLiteChatStore struct {
	db *sql.DB
}

func NewChatStore(db *sql.DB) ChatStore {
	return &SQLiteChatStore{db: db}
}

func (s *SQLiteChatStore) GetByChatID(chatID string) (*models.Chat, error) {
	chat := &models.Chat{}
	err := s.db.QueryRow("SELECT * FROM chat WHERE chat_id = ?", chatID).
		Scan(&chat.ChatID, &chat.Text, &chat.Type, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *SQLiteChatStore) CreateSession() (string, error) {
	sessionID := GenerateUUID()
	_, err := s.db.Exec("INSERT INTO session (session_id) VALUES (?)", sessionID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (s *SQLiteChatStore) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	err := s.db.QueryRow("SELECT * FROM session WHERE session_id = ?", sessionID).Scan(
		&session.SessionID,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SQLiteChatStore) CreateChat(chat *models.Chat) error {
	_, err := s.db.Exec(
		"INSERT INTO chat (chat_id, session_id, text, type) VALUES (?, ?, ?, ?)",
		chat.ChatID,
		chat.SessionID,
		chat.Text,
		chat.Type,
	)
	return err
}

func (s *SQLiteChatStore) GetChatsBySession(sessionID string) ([]*models.Chat, error) {
	rows, err := s.db.Query("SELECT * FROM chat WHERE session_id = ? ORDER BY created_at ASC", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*models.Chat
	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(
			&chat.MessageID,
			&chat.ChatID,
			&chat.SessionID,
			&chat.Text,
			&chat.Type,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, &chat)
	}
	return chats, nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
