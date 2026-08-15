-- Allow question bank questions without a mandatory quiz_id
ALTER TABLE courses.questions ALTER COLUMN quiz_id DROP NOT NULL;
ALTER TABLE courses.questions DROP CONSTRAINT IF EXISTS questions_quiz_id_fkey,
    ADD CONSTRAINT questions_quiz_id_fkey FOREIGN KEY (quiz_id) REFERENCES courses.quizzes (id) ON DELETE SET NULL;
