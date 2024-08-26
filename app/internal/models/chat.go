package models

import (
	"time"
)

type Chat struct {
	ChatID    string    `json:"chat_id"`
	Text      string    `json:"text"`
	Type      *string   `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
