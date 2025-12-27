package dto

import (
	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	Title      string    `json:"title" validate:"required"`
	Content    string    `json:"content" validate:"required"`
	NotebookId uuid.UUID `json:"notebook_id" validate:"required"`
}

type CreateNoteResponse struct {
	Id uuid.UUID `json:"id"`
}

type UpdateNoteRequest struct {
	Id      uuid.UUID `json:"id"`
	Title   string    `json:"title" validate:"required"`
	Content string    `json:"content" validate:"required"`
}

type UpdateNoteResponse struct {
	Id uuid.UUID `json:"id"`
}
