package repository

import (
	"ai-notetaking-be/internal/entity"
	"ai-notetaking-be/internal/pkg/serverutils"
	"ai-notetaking-be/pkg/database"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type INoteRepository interface {
	UsingTx(ctx context.Context, tx database.DatabaseQueryer) INoteRepository
	Create(ctx context.Context, note *entity.Note) error
	GetById(ctx context.Context, id uuid.UUID) (*entity.Note, error)
	GetByNotebookIds(ctx context.Context, notebookIds []uuid.UUID) ([]*entity.Note, error)
	Update(ctx context.Context, note *entity.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByNotebookId(ctx context.Context, uuid uuid.UUID) error
}

type noteRepository struct {
	db database.DatabaseQueryer
}

func (n *noteRepository) UsingTx(ctx context.Context, tx database.DatabaseQueryer) INoteRepository {
	return &noteRepository{
		db: tx,
	}
}

func NewNoteRepository(db *pgxpool.Pool) INoteRepository {
	return &noteRepository{
		db: db,
	}
}

func (n *noteRepository) Create(ctx context.Context, note *entity.Note) error {

	_, err := n.db.Exec(
		ctx,
		`INSERT INTO note (id, title, content, notebook_id, created_at, updated_at, deleted_at, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		note.Id, note.Title, note.Content, note.NotebookId, note.CreatedAt, note.UpdatedAt, note.DeletedAt, note.IsDeleted,
	)

	if err != nil {
		return err
	}

	return nil
}

func (n *noteRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Note, error) {

	row := n.db.QueryRow(
		ctx,
		`SELECT id, title, content, notebook_id, created_at, updated_at, deleted_at, is_deleted FROM note WHERE id = $1`,
		id,
	)
	var note entity.Note
	err := row.Scan(
		&note.Id, &note.Title, &note.Content, &note.NotebookId, &note.CreatedAt, &note.UpdatedAt, &note.DeletedAt, &note.IsDeleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serverutils.ErrNotFound
		}
		return nil, err
	}

	return &note, nil
}

func (n *noteRepository) Update(ctx context.Context, note *entity.Note) error {

	_, err := n.db.Exec(
		ctx,
		`UPDATE note SET title = $2, content = $3, updated_at = $4, notebook_id = $5 WHERE id = $1`,
		note.Id, note.Title, note.Content, note.UpdatedAt, note.NotebookId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (n *noteRepository) Delete(ctx context.Context, id uuid.UUID) error {

	_, err := n.db.Exec(
		ctx,
		`UPDATE note SET is_deleted = $2, deleted_at = $3 WHERE id = $1`,
		id, true, time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (n *noteRepository) DeleteByNotebookId(ctx context.Context, id uuid.UUID) error {

	_, err := n.db.Exec(
		ctx,
		`UPDATE note SET is_deleted = $2, deleted_at = $3 WHERE notebook_id = $1`,
		id, true, time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (n *noteRepository) GetByNotebookIds(ctx context.Context, notebookIds []uuid.UUID) ([]*entity.Note, error) {
	if len(notebookIds) == 0 {
		return make([]*entity.Note, 0), nil
	}

	idStr := make([]string, len(notebookIds))
	for i, id := range notebookIds {
		idStr[i] = id.String()
	}

	rows, err := n.db.Query(
		ctx,
		`SELECT id, title, content, notebook_id, created_at, updated_at, deleted_at, is_deleted FROM note WHERE notebook_id = ANY($1) and is_deleted = false`,
		idStr,
	)
	if err != nil {
		return nil, err
	}

	var notes []*entity.Note
	for rows.Next() {
		var note entity.Note
		err := rows.Scan(
			&note.Id, &note.Title, &note.Content, &note.NotebookId, &note.CreatedAt, &note.UpdatedAt, &note.DeletedAt, &note.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}
