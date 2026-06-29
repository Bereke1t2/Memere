import 'quiz_question_entity.dart';

enum QuizAttemptStatus { inProgress, submitted, expired, graded }

class QuizAttemptEntity {
  const QuizAttemptEntity({
    required this.attemptId,
    required this.quizId,
    required this.attemptNumber,
    required this.status,
    required this.startedAt,
    this.expiresAt,
    this.remainingSeconds,
    required this.questions,
  });

  final String attemptId;
  final String quizId;
  final int attemptNumber;
  final QuizAttemptStatus status;
  final DateTime? startedAt;
  final DateTime? expiresAt;
  final int? remainingSeconds;
  final List<QuizQuestionEntity> questions;

  bool get hasTimer => expiresAt != null || remainingSeconds != null;
  bool get isInProgress => status == QuizAttemptStatus.inProgress;
}
