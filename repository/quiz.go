package repository

import (
	"context"
	"errors"
	"fmt"

	"go-openai-server/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertQuizSQL = `
INSERT INTO test_quizzes (id, type, day_num, week_num)
VALUES ($1, $2, $3, $4)
`

const insertQuizQuestionSQL = `
INSERT INTO test_quiz_questions (id, quiz_id, question_text, order_index)
VALUES ($1, $2, $3, $4)
`

const insertQuizChoiceSQL = `
INSERT INTO test_quiz_choices (id, question_id, choice_key, choice_text, is_correct, order_index)
VALUES ($1, $2, $3, $4, $5, $6)
`

const insertQuizExplanationSQL = `
INSERT INTO test_quiz_explanations (id, question_id, language, explanation, order_index)
VALUES ($1, $2, $3, $4, $5)
`

type QuizRepository struct {
	DB *pgxpool.Pool
}

func NewQuizRepository(db *pgxpool.Pool) *QuizRepository {
	return &QuizRepository{DB: db}
}

func (r *QuizRepository) SaveQuizBundle(ctx context.Context, bundle models.QuizBundle) error {
	if r == nil || r.DB == nil {
		return errors.New("quiz repository is not initialized")
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, insertQuizSQL, bundle.Quiz.ID, bundle.Quiz.Type, bundle.Quiz.DayNum, bundle.Quiz.WeekNum); err != nil {
		return fmt.Errorf("insert quiz failed: %w", err)
	}

	questionArgs := make([][]any, 0, len(bundle.Questions))
	for _, question := range bundle.Questions {
		questionArgs = append(questionArgs, []any{
			question.ID,
			question.QuizID,
			question.QuestionText,
			question.OrderIndex,
		})
	}

	if err := execBatch(ctx, tx, insertQuizQuestionSQL, questionArgs, "quiz question"); err != nil {
		return err
	}

	choiceArgs := make([][]any, 0, len(bundle.Choices))
	for _, choice := range bundle.Choices {
		choiceArgs = append(choiceArgs, []any{
			choice.ID,
			choice.QuestionID,
			choice.ChoiceKey,
			choice.ChoiceText,
			choice.IsCorrect,
			choice.OrderIndex,
		})
	}

	if err := execBatch(ctx, tx, insertQuizChoiceSQL, choiceArgs, "quiz choice"); err != nil {
		return err
	}

	explanationArgs := make([][]any, 0, len(bundle.Explanations))
	for _, explanation := range bundle.Explanations {
		explanationArgs = append(explanationArgs, []any{
			explanation.ID,
			explanation.QuestionID,
			explanation.Language,
			explanation.Explanation,
			explanation.OrderIndex,
		})
	}

	if err := execBatch(ctx, tx, insertQuizExplanationSQL, explanationArgs, "quiz explanation"); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

func execBatch(
	ctx context.Context,
	tx pgx.Tx,
	sql string,
	rows [][]any,
	label string,
) error {
	if len(rows) == 0 {
		return nil
	}

	var batch pgx.Batch
	for _, row := range rows {
		batch.Queue(sql, row...)
	}

	br := tx.SendBatch(ctx, &batch)

	for i := range rows {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return fmt.Errorf("insert %s %d failed: %w", label, i, err)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("close %s batch failed: %w", label, err)
	}

	return nil
}
