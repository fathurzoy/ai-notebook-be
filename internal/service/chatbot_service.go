package service

import (
	"ai-notetaking-be/internal/contain"
	"ai-notetaking-be/internal/dto"
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IChatbotService interface {
	CreateSession(ctx context.Context) (*dto.CreateSessionResponse, error)
}

type chatbotService struct {
	db                          *pgxpool.Pool
	chatbotSessionRepository    repository.IChatSessionRepository
	chatbotMessageRepository    repository.IChatMessageRepository
	chatbotMessageRawRepository repository.IChatMessageRawRepository
}

func (s *chatbotService) CreateSession(ctx context.Context) (*dto.CreateSessionResponse, error) { // TODO: implement create session// messages to the chatbot.

	now := time.Now()
	chatSession := &entity.ChatSession{
		Id:        uuid.New(),
		Title:     "Unnamed session",
		CreatedAt: now,
	}

	chatMessage := &entity.ChatMessage{
		Id:            uuid.New(),
		Chat:          "Hi, how can I help you ?",
		Role:          contain.ChatMessageRoleModel,
		ChatSessionId: chatSession.Id,
		CreatedAt:     now,
	}

	chatMessageRaw := &entity.ChatMessageRaw{
		Id:            uuid.New(),
		Chat:          contain.ChatMessageRawInitialUserPromptV1,
		Role:          contain.ChatMessageRoleUser,
		ChatSessionId: chatSession.Id,
		CreatedAt:     now,
	}

	chatMessageRawUser := &entity.ChatMessageRaw{
		Id:            uuid.New(),
		Chat:          contain.ChatMessageRawInitialUserPromptV1,
		Role:          contain.ChatMessageRoleUser,
		ChatSessionId: chatSession.Id,
		CreatedAt:     now,
	}

	chatMessageRawModel := &entity.ChatMessageRaw{
		Id:            uuid.New(),
		Chat:          contain.ChatMessageRawInitialModelPromptV1,
		Role:          contain.ChatMessageRoleModel,
		ChatSessionId: chatSession.Id,
		CreatedAt:     now.Add(1 * time.Second),
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	chatSessionRepository := s.chatbotSessionRepository.UsingTx(ctx, tx)
	chatMessageRepository := s.chatbotMessageRepository.UsingTx(ctx, tx)
	chatMessageRawRepository := s.chatbotMessageRawRepository.UsingTx(ctx, tx)

	err = chatSessionRepository.Create(ctx, chatSession)
	if err != nil {
		return nil, err
	}

	err = chatMessageRepository.Create(ctx, chatMessage)
	if err != nil {
		return nil, err
	}

	err = chatMessageRawRepository.Create(ctx, chatMessageRaw)
	if err != nil {
		return nil, err
	}

	err = chatMessageRawRepository.Create(ctx, chatMessageRawUser)
	if err != nil {
		return nil, err
	}

	err = chatMessageRawRepository.Create(ctx, chatMessageRawModel)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: insert ke chat session table

	// TODO: insert default chat ke chat message table

	// TODO: insert default raw chat ke chat message raw table

	return &dto.CreateSessionResponse{
		Id: chatSession.Id,
	}, nil
}

func NewChatbotService(db *pgxpool.Pool, chatbotSessionRepository repository.IChatSessionRepository, chatbotMessageRepository repository.IChatMessageRepository, chatbotMessageRawRepository repository.IChatMessageRawRepository) IChatbotService {
	return &chatbotService{
		db:                          db,
		chatbotSessionRepository:    repository.NewChatSessionRepository(db),
		chatbotMessageRepository:    repository.NewChatMessageRepository(db),
		chatbotMessageRawRepository: repository.NewChatMessageRawRepository(db),
	}
}
