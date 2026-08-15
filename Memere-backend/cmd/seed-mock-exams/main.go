// Command seed-mock-exams seeds 4 standard 60-question Grade 12 National Mock Exams:
// Mathematics, Physics, Chemistry, Biology.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/database"
)

type Option struct {
	Text      string
	IsCorrect bool
}

type QuestionData struct {
	Text        string
	Topic       string
	Explanation string
	Options     []Option
}

type ExamData struct {
	Title           string
	Subject         string
	Grade           int
	DurationMinutes int
	PassMarks       int
	Instructions    string
	Questions       []QuestionData
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	exams := []ExamData{
		buildMathExam(),
		buildPhysicsExam(),
		buildChemistryExam(),
		buildBiologyExam(),
	}

	for _, ed := range exams {
		if err := seedExam(ctx, pool, ed); err != nil {
			log.Fatalf("failed seeding exam %s: %v", ed.Title, err)
		}
	}

	fmt.Println("Successfully seeded 4 mock exams with 60 questions each!")
}

func seedExam(ctx context.Context, pool *pgxpool.Pool, ed ExamData) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if exam exists by title or create
	var examID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT id FROM courses.exams WHERE title = $1 AND deleted_at IS NULL LIMIT 1
	`, ed.Title).Scan(&examID)

	if err != nil {
		// Insert exam
		err = tx.QueryRow(ctx, `
			INSERT INTO courses.exams (
				title, subject, grade, duration_minutes, total_marks, pass_marks, instructions, is_published, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, true, NOW(), NOW()
			) RETURNING id
		`, ed.Title, ed.Subject, ed.Grade, ed.DurationMinutes, len(ed.Questions), ed.PassMarks, ed.Instructions).Scan(&examID)
		if err != nil {
			return fmt.Errorf("insert exam: %w", err)
		}
	} else {
		// Update exam properties
		_, err = tx.Exec(ctx, `
			UPDATE courses.exams
			SET subject = $2, grade = $3, duration_minutes = $4, total_marks = $5, pass_marks = $6, instructions = $7, is_published = true, updated_at = NOW()
			WHERE id = $1
		`, examID, ed.Subject, ed.Grade, ed.DurationMinutes, len(ed.Questions), ed.PassMarks, ed.Instructions)
		if err != nil {
			return fmt.Errorf("update exam: %w", err)
		}
		// Clear existing links to avoid duplicate order
		_, _ = tx.Exec(ctx, `DELETE FROM courses.exam_questions WHERE exam_id = $1`, examID)
	}

	for i, qd := range ed.Questions {
		var questionID uuid.UUID
		err = tx.QueryRow(ctx, `
			INSERT INTO courses.questions (
				text, type, points, explanation, order_index, subject, topic, created_at, updated_at
			) VALUES (
				$1, 'multiple_choice', 1, $2, $3, $4, $5, NOW(), NOW()
			) RETURNING id
		`, qd.Text, qd.Explanation, i, ed.Subject, qd.Topic).Scan(&questionID)
		if err != nil {
			return fmt.Errorf("insert question %d: %w", i+1, err)
		}

		for j, opt := range qd.Options {
			_, err = tx.Exec(ctx, `
				INSERT INTO courses.answers (
					question_id, text, is_correct, order_index, created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, NOW(), NOW()
				)
			`, questionID, opt.Text, opt.IsCorrect, j)
			if err != nil {
				return fmt.Errorf("insert answer for question %d option %d: %w", i+1, j+1, err)
			}
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO courses.exam_questions (
				exam_id, question_id, order_index, marks
			) VALUES (
				$1, $2, $3, 1
			)
		`, examID, questionID, i)
		if err != nil {
			return fmt.Errorf("link exam question %d: %w", i+1, err)
		}
	}

	return tx.Commit(ctx)
}
