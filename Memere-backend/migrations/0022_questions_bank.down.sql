-- Rollback: restore NOT NULL on quiz_id (if all rows have quiz_id)
ALTER TABLE courses.questions ALTER COLUMN quiz_id SET NOT NULL;
