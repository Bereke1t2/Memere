import 'question_feedback_entity.dart';
import 'quiz_attempt_entity.dart';
import 'subject_score_entity.dart';

class QuizResultEntity {
  const QuizResultEntity({
    required this.attemptId,
    required this.quizId,
    required this.attemptNumber,
    required this.status,
    required this.score,
    required this.totalPoints,
    required this.percentage,
    required this.passed,
    this.submittedAt,
    required this.feedback,
    required this.subjectBreakdown,
  });

  final String attemptId;
  final String quizId;
  final int attemptNumber;
  final QuizAttemptStatus status;
  final double score;
  final int totalPoints;
  final double percentage;
  final bool passed;
  final DateTime? submittedAt;
  final List<QuestionFeedbackEntity> feedback;
  final Map<String, SubjectScoreEntity> subjectBreakdown;
}
