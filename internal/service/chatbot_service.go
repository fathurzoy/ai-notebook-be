package service

import (
	"ai-notetaking-be/internal/constant"
	"ai-notetaking-be/internal/dto"
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/internal/repository"
	"ai-notetaking-be/pkg/chatbot"
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IChatbotService interface {
	CreateSession(ctx context.Context) (*dto.CreateSessionResponse, error)
	GetAllSessions(ctx context.Context) ([]*dto.GetAllSessionsResponse, error)
	GetChatHistory(ctx context.Context, chatSessionId uuid.UUID) ([]*dto.GetChatHistoryResponse, error)
	SendChat(ctx context.Context, request *dto.SendChatRequest) (*dto.SendChatResponse, error)
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
		Role:          constant.ChatMessageRoleModel,
		ChatSessionId: chatSession.Id,
		CreatedAt:     now,
	}

	chatMessageRaw := &entity.ChatMessageRaw{
		Id:            uuid.New(),
		Chat:          constant.ChatMessageRawInitialUserPromptV1,
		Role:          constant.ChatMessageRoleUser,
		ChatSessionId: chatSession.Id,
		CreatedAt:     now,
	}

	chatMessageRawUser := &entity.ChatMessageRaw{
		Id:            uuid.New(),
		Chat:          constant.ChatMessageRawInitialUserPromptV1,
		Role:          constant.ChatMessageRoleUser,
		ChatSessionId: chatSession.Id,
		CreatedAt:     now,
	}

	chatMessageModelRaw := &entity.ChatMessageRaw{
		Id:            uuid.New(),
		Chat:          constant.ChatMessageRawInitialModelPromptV1,
		Role:          constant.ChatMessageRoleModel,
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

	err = chatMessageRawRepository.Create(ctx, chatMessageModelRaw)
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

func (s *chatbotService) GetAllSessions(ctx context.Context) ([]*dto.GetAllSessionsResponse, error) {
	chatSessions, err := s.chatbotSessionRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []*dto.GetAllSessionsResponse
	for _, chatSession := range chatSessions {
		response = append(response, &dto.GetAllSessionsResponse{
			Id:        chatSession.Id,
			Title:     chatSession.Title,
			CreatedAt: chatSession.CreatedAt,
			UpdatedAt: chatSession.UpdatedAt,
			DeletedAt: chatSession.DeletedAt,
		})
	}

	return response, nil

}

func (s *chatbotService) GetChatHistory(ctx context.Context, chatSessionId uuid.UUID) ([]*dto.GetChatHistoryResponse, error) {

	_, err := s.chatbotSessionRepository.GetById(ctx, chatSessionId)
	if err != nil {
		return nil, err
	}

	chatMessages, err := s.chatbotMessageRepository.GetByChatSessionId(ctx, chatSessionId)
	if err != nil {
		return nil, err
	}

	var response []*dto.GetChatHistoryResponse
	for _, chatMessage := range chatMessages {
		response = append(response, &dto.GetChatHistoryResponse{
			Id:        chatMessage.Id,
			Chat:      chatMessage.Chat,
			Role:      chatMessage.Role,
			CreatedAt: chatMessage.CreatedAt,
		})
	}

	return response, nil
}

func (s *chatbotService) SendChat(ctx context.Context, request *dto.SendChatRequest) (*dto.SendChatResponse, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	chatSessionRepository := s.chatbotSessionRepository.UsingTx(ctx, tx)
	chatMessageRepository := s.chatbotMessageRepository.UsingTx(ctx, tx)
	chatMessageRawRepository := s.chatbotMessageRawRepository.UsingTx(ctx, tx)

	chatSession, err := chatSessionRepository.GetById(ctx, request.ChatSessionId)
	if err != nil {
		return nil, err
	}

	existingRawChats, err := chatMessageRawRepository.GetByChatSessionId(ctx, request.ChatSessionId)
	if err != nil {
		return nil, err
	}
	updateSessionTItle := len(existingRawChats) == 2

	now := time.Now()

	// TODO: save user chat ke db
	chatMessage := &entity.ChatMessage{
		Id:            uuid.New(),
		Chat:          request.Chat,
		Role:          constant.ChatMessageRoleUser,
		ChatSessionId: request.ChatSessionId,
		CreatedAt:     now,
	}

	strBuilder := strings.Builder{}
	strBuilder.WriteString("User next question: ")
	strBuilder.WriteString(request.Chat)
	strBuilder.WriteString("\n\n")
	strBuilder.WriteString("Your answer ?")
	// TODO: save user raw chat ke db
	chatMessageRaw := &entity.ChatMessageRaw{
		Id:            uuid.New(),
		Chat:          strBuilder.String(),
		Role:          constant.ChatMessageRoleUser,
		ChatSessionId: request.ChatSessionId,
		CreatedAt:     now,
	}

	existingRawChats = append(existingRawChats, chatMessageRaw)

	geminiReq := make([]*chatbot.ChatHistory, 0)

	for _, existingRawChat := range existingRawChats {
		geminiReq = append(geminiReq, &chatbot.ChatHistory{
			Chat: existingRawChat.Chat,
			Role: existingRawChat.Role,
		})
	}

	reply, err := chatbot.GetGeminiResponse(ctx, os.Getenv("GOOGLE_GEMINI_API_KEY"), geminiReq)
	if err != nil {
		return nil, err
	}

	log.Printf("Reply: %v", reply)

	// TODO: save dummy model reply ke db
	chatMessageModel := &entity.ChatMessage{
		Id: uuid.New(),
		// Chat:          "This is automated dummy response",
		Chat:          reply,
		Role:          constant.ChatMessageRoleModel,
		ChatSessionId: request.ChatSessionId,
		CreatedAt:     now.Add(1 * time.Second),
	}

	// TODO: save dummy model raw reply ke db
	chatMessageModelRaw := &entity.ChatMessageRaw{
		Id: uuid.New(),
		// Chat:          "This is automated dummy response",
		Chat:          reply,
		Role:          constant.ChatMessageRoleModel,
		ChatSessionId: request.ChatSessionId,
		CreatedAt:     now.Add(1 * time.Second),
	}

	err = chatMessageRepository.Create(ctx, chatMessage)
	if err != nil {
		return nil, err
	}

	err = chatMessageRepository.Create(ctx, chatMessageModel)
	if err != nil {
		return nil, err
	}

	err = chatMessageRawRepository.Create(ctx, chatMessageRaw)
	if err != nil {
		return nil, err
	}

	err = chatMessageRawRepository.Create(ctx, chatMessageModelRaw)
	if err != nil {
		return nil, err
	}

	if updateSessionTItle {
		chatSession.Title = request.Chat
		chatSession.UpdatedAt = &now

		err = chatSessionRepository.Update(ctx, chatSession)
		if err != nil {
			return nil, err
		}

		//TODO: create update repo

		err = tx.Commit(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &dto.SendChatResponse{
		ChatSessionId:    chatSession.Id,
		ChatSessionTitle: chatSession.Title,
		Sent: &dto.SendChatResponseChat{
			Id:        chatMessage.Id,
			Chat:      chatMessage.Chat,
			Role:      chatMessage.Role,
			CreatedAt: chatMessage.CreatedAt,
		},
		Reply: &dto.SendChatResponseChat{
			Id:        chatMessageModel.Id,
			Chat:      chatMessageModel.Chat,
			Role:      chatMessageModel.Role,
			CreatedAt: chatMessageModel.CreatedAt,
		},
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
