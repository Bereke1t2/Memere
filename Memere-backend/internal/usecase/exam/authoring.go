package exam

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/validator"
)

const (
	maxExamTitleLen   = 200
	maxInstructionLen = 5000
	minGrade          = 1
	maxGrade          = 12
)

// CreateExamInput is the teacher-facing exam create request after decoding.
type CreateExamInput struct {
	CourseID        *uuid.UUID
	Title           string
	Subject         string
	Grade           int
	DurationMinutes int
	PassMarks       int
	Instructions    *string
}

// UpdateExamInput carries partial exam updates; nil fields are unchanged.
// TotalMarks is intentionally absent — it is derived from the question set.
type UpdateExamInput struct {
	Title           *string
	Subject         *string
	Grade           *int
	DurationMinutes *int
	PassMarks       *int
	Instructions    *string
}

// CreateExam creates an exam, optionally under a course the actor owns. total_marks
// starts at 0 and is maintained from the question set as questions are added.
func (s *Service) CreateExam(ctx context.Context, actor *Actor, in CreateExamInput) (*entity.Exam, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may create exams", nil)
	}
	if in.CourseID != nil {
		if _, err := s.loadOwnedCourse(ctx, actor, *in.CourseID); err != nil {
			return nil, err
		}
	}
	if details := validateExamInput(in.Title, in.Subject, in.Grade, in.DurationMinutes, in.PassMarks, 0, in.Instructions); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}

	exam := &entity.Exam{
		CourseID:        in.CourseID,
		Title:           in.Title,
		Subject:         in.Subject,
		Grade:           in.Grade,
		DurationMinutes: in.DurationMinutes,
		TotalMarks:      0,
		PassMarks:       in.PassMarks,
		Instructions:    in.Instructions,
		IsPublished:     false,
	}
	if err := s.exams.Create(ctx, exam); err != nil {
		return nil, err
	}
	return exam, nil
}

// AddExamQuestion links a bank question to an exam and recomputes the exam's
// total_marks from the sum of its question marks, transactionally, so the
// header total never drifts from the parts (skill decision: total is derived).
func (s *Service) AddExamQuestion(ctx context.Context, actor *Actor, examID, questionID uuid.UUID, marks, orderIndex int) (*entity.Exam, error) {
	exam, err := s.loadOwnedExam(ctx, actor, examID)
	if err != nil {
		return nil, err
	}
	if marks <= 0 {
		return nil, apperror.Validation(map[string]any{"marks": "must be positive"}, nil)
	}

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		eq := &entity.ExamQuestion{ExamID: examID, QuestionID: questionID, OrderIndex: orderIndex, Marks: marks}
		if err := s.exams.AddQuestion(ctx, eq); err != nil {
			return err
		}
		links, err := s.exams.ListQuestions(ctx, examID)
		if err != nil {
			return err
		}
		total := 0
		for _, l := range links {
			total += l.Marks
		}
		exam.TotalMarks = total
		return s.exams.Update(ctx, exam)
	})
	if err != nil {
		return nil, err
	}
	return exam, nil
}

// ListExamQuestions returns an exam's question links (ownership-checked).
func (s *Service) ListExamQuestions(ctx context.Context, actor *Actor, examID uuid.UUID) ([]*entity.ExamQuestion, error) {
	if _, err := s.loadOwnedExam(ctx, actor, examID); err != nil {
		return nil, err
	}
	return s.exams.ListQuestions(ctx, examID)
}

// UpdateExam applies partial updates to an exam the actor owns.
func (s *Service) UpdateExam(ctx context.Context, actor *Actor, examID uuid.UUID, in UpdateExamInput) (*entity.Exam, error) {
	exam, err := s.loadOwnedExam(ctx, actor, examID)
	if err != nil {
		return nil, err
	}

	if in.Title != nil {
		exam.Title = *in.Title
	}
	if in.Subject != nil {
		exam.Subject = *in.Subject
	}
	if in.Grade != nil {
		exam.Grade = *in.Grade
	}
	if in.DurationMinutes != nil {
		exam.DurationMinutes = *in.DurationMinutes
	}
	if in.PassMarks != nil {
		exam.PassMarks = *in.PassMarks
	}
	if in.Instructions != nil {
		exam.Instructions = in.Instructions
	}

	if details := validateExamInput(exam.Title, exam.Subject, exam.Grade, exam.DurationMinutes, exam.PassMarks, exam.TotalMarks, exam.Instructions); len(details) > 0 {
		return nil, apperror.Validation(details, nil)
	}
	if err := s.exams.Update(ctx, exam); err != nil {
		return nil, err
	}
	return exam, nil
}

// PublishExam makes an exam visible to students (ownership-checked).
func (s *Service) PublishExam(ctx context.Context, actor *Actor, examID uuid.UUID) (*entity.Exam, error) {
	return s.setPublished(ctx, actor, examID, true)
}

// UnpublishExam hides an exam from students (ownership-checked).
func (s *Service) UnpublishExam(ctx context.Context, actor *Actor, examID uuid.UUID) (*entity.Exam, error) {
	return s.setPublished(ctx, actor, examID, false)
}

func (s *Service) setPublished(ctx context.Context, actor *Actor, examID uuid.UUID, published bool) (*entity.Exam, error) {
	exam, err := s.loadOwnedExam(ctx, actor, examID)
	if err != nil {
		return nil, err
	}
	exam.IsPublished = published
	if err := s.exams.Update(ctx, exam); err != nil {
		return nil, err
	}
	return exam, nil
}

// DeleteExam soft-deletes an exam the actor owns.
func (s *Service) DeleteExam(ctx context.Context, actor *Actor, examID uuid.UUID) error {
	if _, err := s.loadOwnedExam(ctx, actor, examID); err != nil {
		return err
	}
	return s.exams.Delete(ctx, examID)
}

// ListExams returns exams matching the filter. Non-privileged callers
// (anonymous/student) see published exams only; teachers/admins may filter
// freely (including their unpublished exams).
func (s *Service) ListExams(ctx context.Context, actor *Actor, filter repository.ExamFilter, cursor *pagination.Cursor, limit int) ([]*entity.Exam, *pagination.Cursor, error) {
	if !actor.isTeacherOrAdmin() {
		published := true
		filter.IsPublished = &published
	}
	return s.exams.List(ctx, filter, cursor, limit)
}

// loadOwnedCourse loads a course and asserts the actor may author under it.
func (s *Service) loadOwnedCourse(ctx context.Context, actor *Actor, courseID uuid.UUID) (*entity.Course, error) {
	course, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if !actor.ownsCourse(course) {
		return nil, apperror.Forbidden("you do not own this course", nil)
	}
	return course, nil
}

// loadOwnedExam loads an exam and asserts the actor may author it: an admin, or
// the teacher who owns the exam's parent course. A course-less exam is editable
// by admins only.
func (s *Service) loadOwnedExam(ctx context.Context, actor *Actor, examID uuid.UUID) (*entity.Exam, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may author exams", nil)
	}
	exam, err := s.exams.FindByID(ctx, examID)
	if err != nil {
		return nil, err
	}
	if actor.Role == entity.RoleAdmin {
		return exam, nil
	}
	if exam.CourseID == nil {
		return nil, apperror.Forbidden("only an admin may edit a standalone exam", nil)
	}
	if _, err := s.loadOwnedCourse(ctx, actor, *exam.CourseID); err != nil {
		return nil, err
	}
	return exam, nil
}

func validateExamInput(title, subject string, grade, duration, passMarks, totalMarks int, instructions *string) map[string]any {
	v := validator.New()
	v.Required("title", title)
	v.MaxLen("title", title, maxExamTitleLen)
	v.Required("subject", subject)
	v.InRange("grade", grade, minGrade, maxGrade)
	if duration <= 0 {
		v.Add("duration_minutes", "must be positive")
	}
	if passMarks < 0 {
		v.Add("pass_marks", "must not be negative")
	}
	// pass_marks must be reachable. total_marks is 0 until questions are added, so
	// only enforce the relation once the exam has marks.
	if totalMarks > 0 && passMarks > totalMarks {
		v.Add("pass_marks", "must not exceed total_marks")
	}
	if instructions != nil {
		v.MaxLen("instructions", *instructions, maxInstructionLen)
	}
	return v.Map()
}
