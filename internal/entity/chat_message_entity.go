package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	Id            uuid.UUID  `json:"id"`
	Chat          string     `json:"chat"`
	Role          string     `json:"role"`
	ChatSessionId uuid.UUID  `json:"chat_session_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	IsDeleted     bool       `json:"is_deleted"`
}
