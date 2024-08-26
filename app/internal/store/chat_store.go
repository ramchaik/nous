package store

import (
	"database/sql"
	"time"

	"nous/internal/models"

	"github.com/google/uuid"
)

type ChatStore interface {
	Create(*models.Chat) error
	GetByID(string) (*models.Chat, error)
	// TODO: implement CRUDs
	// Update(*models.Chat) error
	// Delete(string) error
}

type SQLiteChatStore struct {
	db *sql.DB
}

func NewChatStore(db *sql.DB) ChatStore {
	return &SQLiteChatStore{db: db}
}

func (s *SQLiteChatStore) Create(chat *models.Chat) error {
	chat.ChatID = uuid.New().String()
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = time.Now()

	_, err := s.db.Exec("INSERT INTO chat (chat_id, text, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		chat.ChatID, chat.Text, chat.Type, chat.CreatedAt, chat.UpdatedAt)
	return err
}

func (s *SQLiteChatStore) GetByID(chatID string) (*models.Chat, error) {
	chat := &models.Chat{}
	err := s.db.QueryRow("SELECT chat_id, text, type, created_at, updated_at FROM chat WHERE chat_id = ?", chatID).
		Scan(&chat.ChatID, &chat.Text, &chat.Type, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
