package models

import "time"

type Session struct {
	SessionID string    `db:"session_id" json:"session_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
