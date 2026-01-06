package repository

import (
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/pkg/database"
	"context"

	"github.com/google/uuid"
)

type IChatMessageRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) IChatMessageRepository
	Create(ctx context.Context, chatMessage *entity.ChatMessage) error
	GetByChatSessionId(ctx context.Context, chatSessionId uuid.UUID) ([]*entity.ChatMessage, error)
}

type chatMessageRepository struct {
	db database.DatabaseQueryer
}

func (n *chatMessageRepository) UsingTx(ctx context.Context, tx database.DatabaseQueryer) IChatMessageRepository {
	return &chatMessageRepository{
		db: tx,
	}
}

func (cs *chatMessageRepository) Create(ctx context.Context, chatMessage *entity.ChatMessage) error {
	_, err := cs.db.Exec(
		ctx,
		`INSERT INTO chat_message (id, chat, role, chat_session_id, created_at, updated_at, deleted_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		chatMessage.Id, chatMessage.Chat, chatMessage.Role, chatMessage.ChatSessionId, chatMessage.CreatedAt, chatMessage.UpdatedAt, chatMessage.DeletedAt, chatMessage.IsDeleted,
	)

	if err != nil {
		return err
	}

	return nil
}

func (cs *chatMessageRepository) GetByChatSessionId(ctx context.Context, chatSessionId uuid.UUID) ([]*entity.ChatMessage, error) {

	rows, err := cs.db.Query(
		ctx,
		`SELECT id, chat, role, chat_session_id, created_at, updated_at, deleted_at, is_deleted FROM chat_message WHERE chat_session_id = $1 and is_deleted = false ORDER BY created_at ASC`,
		chatSessionId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chatMessages []*entity.ChatMessage
	for rows.Next() {
		var chatMessage entity.ChatMessage
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

func NewChatMessageRepository(db database.DatabaseQueryer) IChatMessageRepository {
	return &chatMessageRepository{
		db: db,
	}
}
