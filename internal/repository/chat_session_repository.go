package repository

import (
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/pkg/database"
	"context"
)

type IChatSessionRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) IChatSessionRepository
	Create(ctx context.Context, chatSession *entity.ChatSession) error
	GetAll(ctx context.Context) ([]*entity.ChatSession, error)
}

type chatSessionRepository struct {
	db database.DatabaseQueryer
}

func (n *chatSessionRepository) UsingTx(ctx context.Context, tx database.DatabaseQueryer) IChatSessionRepository {
	return &chatSessionRepository{
		db: tx,
	}
}

func (cs *chatSessionRepository) Create(ctx context.Context, chatSession *entity.ChatSession) error {
	_, err := cs.db.Exec(
		ctx,
		`INSERT INTO chat_session (id, title, created_at, updated_at, deleted_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)`,
		chatSession.Id, chatSession.Title, chatSession.CreatedAt, chatSession.UpdatedAt, chatSession.DeletedAt, chatSession.IsDeleted,
	)

	if err != nil {
		return err
	}

	return nil

}

func (cs *chatSessionRepository) GetAll(ctx context.Context) ([]*entity.ChatSession, error) {
	rows, err := cs.db.Query(ctx, `SELECT id, title, created_at, updated_at, deleted_at, is_deleted FROM chat_session WHERE is_deleted = false ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chatSessions []*entity.ChatSession
	for rows.Next() {
		var chatSession entity.ChatSession
		err := rows.Scan(&chatSession.Id, &chatSession.Title, &chatSession.CreatedAt, &chatSession.UpdatedAt, &chatSession.DeletedAt, &chatSession.IsDeleted)
		if err != nil {
			return nil, err
		}
		chatSessions = append(chatSessions, &chatSession)
	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return chatSessions, nil
}

func NewChatSessionRepository(db database.DatabaseQueryer) IChatSessionRepository {
	return &chatSessionRepository{
		db: db,
	}
}
