import '../../features/exam/domain/entities/exam_answer_entity.dart';
import '../../features/exam/domain/entities/exam_attempt_entity.dart';
import '../../features/exam/domain/entities/exam_question_entity.dart';
import '../../features/quiz/domain/entities/quiz_answer_entity.dart';
import '../../features/quiz/domain/entities/quiz_attempt_entity.dart';
import '../../features/quiz/domain/entities/quiz_question_entity.dart';
import '../storage/hive/models/offline_exam.dart';
import '../storage/hive/models/offline_question.dart';
import '../storage/hive/models/offline_quiz.dart';

/// Synthesizes an in-progress attempt entity from downloaded (answer-key-bearing)
/// content — but WITHOUT the key. The produced [QuizAttemptEntity] /
/// [ExamAttemptEntity] is the exact online-play projection the attempt screens
/// already render (answers carry only `id`/`text`/`orderIndex`, never
/// `isCorrect`). The attempt is graded on-device at submit; [attemptId] must be a
/// `local-…` id so the notifiers and result providers route it locally instead of
/// to the server.

QuizAttemptEntity buildLocalQuizAttempt(
  OfflineQuiz quiz, {
  required String attemptId,
  DateTime? startedAt,
}) {
  final questions = _sortedByOrder(quiz.questions);
  return QuizAttemptEntity(
    attemptId: attemptId,
    quizId: quiz.id,
    attemptNumber: 1,
    status: QuizAttemptStatus.inProgress,
    startedAt: startedAt ?? DateTime.now(),
    // Quizzes are untimed offline — there is no server clock to honor.
    expiresAt: null,
    remainingSeconds: null,
    questions: questions.map(_toQuizQuestion).toList(),
  );
}

ExamAttemptEntity buildLocalExamAttempt(
  OfflineExam exam, {
  required String attemptId,
  DateTime? startedAt,
}) {
  final started = startedAt ?? DateTime.now();
  final questions = _sortedByOrder(exam.questions);
  final timed = exam.durationMinutes > 0;
  return ExamAttemptEntity(
    attemptId: attemptId,
    examId: exam.id,
    status: ExamAttemptStatus.inProgress,
    startedAt: started,
    // Display-only countdown offline; expiry just prompts submit — the server
    // does not enforce it for a locally graded attempt.
    expiresAt: timed ? started.add(Duration(minutes: exam.durationMinutes)) : null,
    remainingSeconds: timed ? exam.durationMinutes * 60 : null,
    totalMarks: exam.totalMarks,
    questions: questions.map(_toExamQuestion).toList(),
  );
}

// ── mappers ──────────────────────────────────────────────────────────────────

List<OfflineQuestion> _sortedByOrder(List<OfflineQuestion> questions) =>
    [...questions]..sort((a, b) => a.orderIndex.compareTo(b.orderIndex));

QuizQuestionEntity _toQuizQuestion(OfflineQuestion q) => QuizQuestionEntity(
      id: q.id,
      text: q.text,
      type: _quizType(q.type),
      points: q.points,
      subject: q.subject,
      topic: q.topic,
      answers: _sortedAnswers(q)
          .map((a) => QuizAnswerEntity(
                id: a.id,
                text: a.text,
                orderIndex: a.orderIndex,
              ))
          .toList(),
    );

ExamQuestionEntity _toExamQuestion(OfflineQuestion q) => ExamQuestionEntity(
      questionId: q.id,
      text: q.text,
      type: _examType(q.type),
      marks: q.points,
      subject: q.subject,
      topic: q.topic,
      answers: _sortedAnswers(q)
          .map((a) => ExamAnswerEntity(
                id: a.id,
                text: a.text,
                orderIndex: a.orderIndex,
              ))
          .toList(),
    );

List<OfflineAnswer> _sortedAnswers(OfflineQuestion q) =>
    [...q.answers]..sort((a, b) => a.orderIndex.compareTo(b.orderIndex));

// The backend has no `multiple_select`, so a downloaded multiple-choice maps to
// the single-select projection the screens render; `unknown` degrades the same
// way (a radio list) rather than crashing.
QuizQuestionType _quizType(OfflineQuestionType t) {
  switch (t) {
    case OfflineQuestionType.trueFalse:
      return QuizQuestionType.trueFalse;
    case OfflineQuestionType.shortAnswer:
      return QuizQuestionType.shortAnswer;
    case OfflineQuestionType.multipleChoice:
    case OfflineQuestionType.unknown:
      return QuizQuestionType.multipleChoice;
  }
}

ExamQuestionType _examType(OfflineQuestionType t) {
  switch (t) {
    case OfflineQuestionType.trueFalse:
      return ExamQuestionType.trueFalse;
    case OfflineQuestionType.shortAnswer:
      return ExamQuestionType.shortAnswer;
    case OfflineQuestionType.multipleChoice:
    case OfflineQuestionType.unknown:
      return ExamQuestionType.multipleChoice;
  }
}
