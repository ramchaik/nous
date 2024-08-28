package models

import "time"

type Chat struct {
	MessageID int       `db:"message_id" json:"message_id"`
	ChatID    string    `db:"chat_id" json:"chat_id"`
	SessionID string    `db:"session_id" json:"session_id"`
	Text      string    `db:"text" json:"text"`
	Type      string    `db:"type" json:"type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
