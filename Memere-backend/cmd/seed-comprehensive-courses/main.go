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
	Explanation string
	Options     []Option
}

type QuizData struct {
	Title            string
	TimeLimitSeconds int
	PassPercentage   float64
	Questions        []QuestionData
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

type LessonData struct {
	Title           string
	Type            string // "video", "note", "mixed", "quiz"
	DurationSeconds int
	IsFreePreview   bool
	Content         string
	PdfUrl          string
}

type SectionData struct {
	Title       string
	Description string
	Lessons     []LessonData
}

type CourseData struct {
	Title            string
	Slug             string
	Subject          string
	Grade            int
	Description      string
	ShortDescription string
	ThumbnailUrl     string
	Price            float64
	IsFree           bool
	Level            string
	Sections         []SectionData
	Quizzes          []QuizData
	Exams            []ExamData
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

	// Ensure teacher exists
	var teacherID uuid.UUID
	err = pool.QueryRow(ctx, `SELECT id FROM auth.users WHERE role = 'teacher' AND deleted_at IS NULL LIMIT 1`).Scan(&teacherID)
	if err != nil {
		log.Println("No teacher found, querying any active user...")
		err = pool.QueryRow(ctx, `SELECT id FROM auth.users WHERE deleted_at IS NULL LIMIT 1`).Scan(&teacherID)
		if err != nil {
			log.Fatalf("no user found in auth.users to assign as teacher: %v", err)
		}
	}

	courses := []CourseData{
		buildMathCourse(),
		buildPhysicsCourse(),
		buildBiologyCourse(),
		buildChemistryCourse(),
		buildHistoryCivicsCourse(),
	}

	for _, c := range courses {
		if err := seedCourse(ctx, pool, teacherID, c); err != nil {
			log.Fatalf("failed seeding course %s: %v", c.Title, err)
		}
		fmt.Printf("✅ Successfully seeded course: %s (Sections: %d, Quizzes: %d, Exams: %d)\n",
			c.Title, len(c.Sections), len(c.Quizzes), len(c.Exams))
	}

	fmt.Println("\n🎉 All 5 Grade 12 National Exam Prep Courses seeded successfully with complete sections, quizzes, and exams!")
}

func seedCourse(ctx context.Context, pool *pgxpool.Pool, teacherID uuid.UUID, cd CourseData) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert or Update Course
	var courseID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM courses.courses WHERE slug = $1 AND deleted_at IS NULL LIMIT 1`, cd.Slug).Scan(&courseID)
	if err != nil {
		err = tx.QueryRow(ctx, `
			INSERT INTO courses.courses (
				teacher_id, title, slug, description, short_description, subject, grade, thumbnail_url,
				price, currency, is_free, is_published, language, level, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, 'ETB', $10, true, 'en', $11, NOW(), NOW()
			) RETURNING id
		`, teacherID, cd.Title, cd.Slug, cd.Description, cd.ShortDescription, cd.Subject, cd.Grade,
			cd.ThumbnailUrl, cd.Price, cd.IsFree, cd.Level).Scan(&courseID)
		if err != nil {
			return fmt.Errorf("insert course %s: %w", cd.Title, err)
		}
	} else {
		_, err = tx.Exec(ctx, `
			UPDATE courses.courses SET
				title = $2, description = $3, short_description = $4, subject = $5, grade = $6,
				thumbnail_url = $7, price = $8, is_free = $9, is_published = true, level = $10, updated_at = NOW()
			WHERE id = $1
		`, courseID, cd.Title, cd.Description, cd.ShortDescription, cd.Subject, cd.Grade,
			cd.ThumbnailUrl, cd.Price, cd.IsFree, cd.Level)
		if err != nil {
			return fmt.Errorf("update course %s: %w", cd.Title, err)
		}
	}

	// 2. Sections and Lessons
	totalDuration := 0
	totalLessons := 0

	for sIdx, sec := range cd.Sections {
		var sectionID uuid.UUID
		err = tx.QueryRow(ctx, `
			SELECT id FROM courses.course_sections WHERE course_id = $1 AND title = $2 LIMIT 1
		`, courseID, sec.Title).Scan(&sectionID)

		if err != nil {
			err = tx.QueryRow(ctx, `
				INSERT INTO courses.course_sections (
					course_id, title, description, order_index, is_published, created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, true, NOW(), NOW()
				) RETURNING id
			`, courseID, sec.Title, sec.Description, sIdx).Scan(&sectionID)
			if err != nil {
				return fmt.Errorf("insert section %s: %w", sec.Title, err)
			}
		} else {
			_, err = tx.Exec(ctx, `
				UPDATE courses.course_sections SET
					description = $2, order_index = $3, is_published = true, updated_at = NOW()
				WHERE id = $1
			`, sectionID, sec.Description, sIdx)
			if err != nil {
				return fmt.Errorf("update section %s: %w", sec.Title, err)
			}
		}

		for lIdx, les := range sec.Lessons {
			totalLessons++
			totalDuration += les.DurationSeconds

			var lessonID uuid.UUID
			err = tx.QueryRow(ctx, `
				SELECT id FROM courses.lessons WHERE section_id = $1 AND title = $2 LIMIT 1
			`, sectionID, les.Title).Scan(&lessonID)

			if err != nil {
				err = tx.QueryRow(ctx, `
					INSERT INTO courses.lessons (
						section_id, course_id, title, type, order_index, is_free_preview,
						duration_seconds, is_published, content, pdf_url, created_at, updated_at
					) VALUES (
						$1, $2, $3, $4, $5, $6, $7, true, $8, $9, NOW(), NOW()
					) RETURNING id
				`, sectionID, courseID, les.Title, les.Type, lIdx, les.IsFreePreview,
					les.DurationSeconds, les.Content, les.PdfUrl).Scan(&lessonID)
				if err != nil {
					return fmt.Errorf("insert lesson %s: %w", les.Title, err)
				}
			} else {
				_, err = tx.Exec(ctx, `
					UPDATE courses.lessons SET
						type = $2, order_index = $3, is_free_preview = $4, duration_seconds = $5,
						is_published = true, content = $6, pdf_url = $7, updated_at = NOW()
					WHERE id = $1
				`, lessonID, les.Type, lIdx, les.IsFreePreview, les.DurationSeconds,
					les.Content, les.PdfUrl)
				if err != nil {
					return fmt.Errorf("update lesson %s: %w", les.Title, err)
				}
			}
		}
	}

	// Update course aggregate counts
	_, _ = tx.Exec(ctx, `
		UPDATE courses.courses SET
			total_duration_seconds = $2, total_lessons = $3, rating_avg = 4.9, enrollment_count = 1240, updated_at = NOW()
		WHERE id = $1
	`, courseID, totalDuration, totalLessons)

	// 3. Quizzes
	for _, qz := range cd.Quizzes {
		var quizID uuid.UUID
		err = tx.QueryRow(ctx, `
			SELECT id FROM courses.quizzes WHERE course_id = $1 AND title = $2 LIMIT 1
		`, courseID, qz.Title).Scan(&quizID)

		if err != nil {
			err = tx.QueryRow(ctx, `
				INSERT INTO courses.quizzes (
					course_id, title, time_limit_seconds, pass_percentage, randomize_questions, created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, true, NOW(), NOW()
				) RETURNING id
			`, courseID, qz.Title, qz.TimeLimitSeconds, qz.PassPercentage).Scan(&quizID)
			if err != nil {
				return fmt.Errorf("insert quiz %s: %w", qz.Title, err)
			}
		} else {
			_, err = tx.Exec(ctx, `
				UPDATE courses.quizzes SET
					time_limit_seconds = $2, pass_percentage = $3, randomize_questions = true, updated_at = NOW()
				WHERE id = $1
			`, quizID, qz.TimeLimitSeconds, qz.PassPercentage)
			if err != nil {
				return fmt.Errorf("update quiz %s: %w", qz.Title, err)
			}
			_, _ = tx.Exec(ctx, `DELETE FROM courses.questions WHERE quiz_id = $1`, quizID)
		}

		for qIdx, qd := range qz.Questions {
			var questionID uuid.UUID
			err = tx.QueryRow(ctx, `
				INSERT INTO courses.questions (
					quiz_id, text, type, points, explanation, order_index, created_at, updated_at
				) VALUES (
					$1, $2, 'multiple_choice', 1, $3, $4, NOW(), NOW()
				) RETURNING id
			`, quizID, qd.Text, qd.Explanation, qIdx).Scan(&questionID)
			if err != nil {
				return fmt.Errorf("insert quiz question %d: %w", qIdx+1, err)
			}

			for optIdx, opt := range qd.Options {
				_, err = tx.Exec(ctx, `
					INSERT INTO courses.answers (
						question_id, text, is_correct, order_index, created_at, updated_at
					) VALUES (
						$1, $2, $3, $4, NOW(), NOW()
					)
				`, questionID, opt.Text, opt.IsCorrect, optIdx)
				if err != nil {
					return fmt.Errorf("insert quiz answer %d: %w", optIdx+1, err)
				}
			}
		}
	}

	// 4. Exams
	for _, ed := range cd.Exams {
		var examID uuid.UUID
		err = tx.QueryRow(ctx, `
			SELECT id FROM courses.exams WHERE course_id = $1 AND title = $2 LIMIT 1
		`, courseID, ed.Title).Scan(&examID)

		if err != nil {
			err = tx.QueryRow(ctx, `
				INSERT INTO courses.exams (
					course_id, title, subject, grade, duration_minutes, total_marks, pass_marks,
					instructions, is_published, created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, true, NOW(), NOW()
				) RETURNING id
			`, courseID, ed.Title, ed.Subject, ed.Grade, ed.DurationMinutes, len(ed.Questions),
				ed.PassMarks, ed.Instructions).Scan(&examID)
			if err != nil {
				return fmt.Errorf("insert exam %s: %w", ed.Title, err)
			}
		} else {
			_, err = tx.Exec(ctx, `
				UPDATE courses.exams SET
					subject = $2, grade = $3, duration_minutes = $4, total_marks = $5, pass_marks = $6,
					instructions = $7, is_published = true, updated_at = NOW()
				WHERE id = $1
			`, examID, ed.Subject, ed.Grade, ed.DurationMinutes, len(ed.Questions),
				ed.PassMarks, ed.Instructions)
			if err != nil {
				return fmt.Errorf("update exam %s: %w", ed.Title, err)
			}
			_, _ = tx.Exec(ctx, `DELETE FROM courses.exam_questions WHERE exam_id = $1`, examID)
		}

		for qIdx, qd := range ed.Questions {
			var questionID uuid.UUID
			err = tx.QueryRow(ctx, `
				INSERT INTO courses.questions (
					text, type, points, explanation, order_index, created_at, updated_at
				) VALUES (
					$1, 'multiple_choice', 1, $2, $3, NOW(), NOW()
				) RETURNING id
			`, qd.Text, qd.Explanation, qIdx).Scan(&questionID)
			if err != nil {
				return fmt.Errorf("insert exam question %d: %w", qIdx+1, err)
			}

			for optIdx, opt := range qd.Options {
				_, err = tx.Exec(ctx, `
					INSERT INTO courses.answers (
						question_id, text, is_correct, order_index, created_at, updated_at
					) VALUES (
						$1, $2, $3, $4, NOW(), NOW()
					)
				`, questionID, opt.Text, opt.IsCorrect, optIdx)
				if err != nil {
					return fmt.Errorf("insert exam answer %d: %w", optIdx+1, err)
				}
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO courses.exam_questions (
					exam_id, question_id, order_index, marks
				) VALUES (
					$1, $2, $3, 1
				)
			`, examID, questionID, qIdx)
			if err != nil {
				return fmt.Errorf("link exam question %d: %w", qIdx+1, err)
			}
		}
	}

	return tx.Commit(ctx)
}
