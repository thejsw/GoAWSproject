package repository

import (
	"context"
	"errors"
	"fmt"

	"go-openai-server/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

const insertWordSQL = `
INSERT INTO ai_words (english, korean)
VALUES ($1, $2)
`

type WordRepository struct {
	DB *pgxpool.Pool
}

func NewWordRepository(db *pgxpool.Pool) *WordRepository {
	return &WordRepository{DB: db}
}

func (r *WordRepository) InsertMany(ctx context.Context, words []models.Word) error {
	if r == nil || r.DB == nil {
		return errors.New("word repository is not initialized")
	}

	if len(words) == 0 {
		return nil
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	for i, word := range words {
		if _, err := tx.Exec(ctx, insertWordSQL, word.English, word.Korean); err != nil {
			return fmt.Errorf("insert word %d failed: %w", i, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}
