package repository

import (
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/pkg/database"
	"context"

	"github.com/google/uuid"
)

type IChatMessageRawRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) IChatMessageRawRepository
	Create(ctx context.Context, chatMessageRaw *entity.ChatMessageRaw) error
	GetByChatSessionId(ctx context.Context, chatSessionId uuid.UUID) ([]*entity.ChatMessageRaw, error)
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

func (cs *chatMessageRawRepository) GetByChatSessionId(ctx context.Context, chatSessionId uuid.UUID) ([]*entity.ChatMessageRaw, error) {

	rows, err := cs.db.Query(
		ctx,
		`SELECT id, chat, role, chat_session_id, created_at, updated_at, deleted_at, is_deleted FROM chat_message_raw WHERE chat_session_id = $1 and is_deleted = false ORDER BY created_at ASC`,
		chatSessionId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chatMessages []*entity.ChatMessageRaw
	for rows.Next() {
		var chatMessage entity.ChatMessageRaw
		err := rows.Scan(&chatMessage.Id, &chatMessage.Chat, &chatMessage.Role, &chatMessage.ChatSessionId, &chatMessage.CreatedAt, &chatMessage.UpdatedAt, &chatMessage.DeletedAt, &chatMessage.IsDeleted)
		if err != nil {
			return nil, err
		}
		chatMessages = append(chatMessages, &chatMessage)
	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return chatMessages, nil

}

func NewChatMessageRawRepository(db database.DatabaseQueryer) IChatMessageRawRepository {
	return &chatMessageRawRepository{
		db: db,
	}
}
