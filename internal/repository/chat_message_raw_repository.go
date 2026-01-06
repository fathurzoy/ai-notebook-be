package repository

import (
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/pkg/database"
	"context"
)

type IChatMessageRawRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) IChatMessageRawRepository
	Create(ctx context.Context, chatMessageRaw *entity.ChatMessageRaw) error
}

type chatMessageRawRepository struct {
	db database.DatabaseQueryer
}

func (n *chatMessageRawRepository) UsingTx(ctx context.Context, tx database.DatabaseQueryer) IChatMessageRawRepository {
	return &chatMessageRawRepository{
		db: tx,
	}
}

func (cs *chatMessageRawRepository) Create(ctx context.Context, chatMessageRaw *entity.ChatMessageRaw) error {
	_, err := cs.db.Exec(
		ctx,
		`INSERT INTO chat_message_raw (id, chat, role, chat_session_id, created_at) VALUES ($1, $2, $3, $4, $5)`,
		chatMessageRaw.Id, chatMessageRaw.Chat, chatMessageRaw.Role, chatMessageRaw.ChatSessionId, chatMessageRaw.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewChatMessageRawRepository(db database.DatabaseQueryer) IChatMessageRawRepository {
	return &chatMessageRawRepository{
		db: db,
	}
}
