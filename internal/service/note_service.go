package service

import (
	"ai-notetaking-be/internal/dto"
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type INoteService interface {
	Create(ctx context.Context, req *dto.CreateNoteRequest) (*dto.CreateNoteResponse, error)
	Show(ctx context.Context, id uuid.UUID) (*dto.ShowNotebookResponse, error)
}

type noteService struct {
	noteRepository repository.INoteRepository
	db             *pgxpool.Pool
}

func NewNoteService(noteRepository repository.INoteRepository, db *pgxpool.Pool) INoteService {
	return &noteService{
		noteRepository: noteRepository,
		db:             db,
	}
}

func (c *noteService) Create(ctx context.Context, req *dto.CreateNoteRequest) (*dto.CreateNoteResponse, error) {

	note := entity.Note{
		Id:         uuid.New(),
		Title:      req.Title,
		Content:    req.Content,
		NotebookId: req.NotebookId,
		CreatedAt:  time.Now(),
	}

	err := c.noteRepository.Create(ctx, &note)
	if err != nil {
		return nil, err
	}

	return &dto.CreateNoteResponse{
		Id: note.Id,
	}, nil
}

func (c *noteService) Show(ctx context.Context, id uuid.UUID) (*dto.ShowNotebookResponse, error) {

	note, err := c.noteRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	res := dto.ShowNotebookResponse{
		Id:        note.Id,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
	}

	return &res, nil
}
