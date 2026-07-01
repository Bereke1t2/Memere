import 'exam_attempt_entity.dart';
import 'exam_question_feedback_entity.dart';
import 'exam_subject_score_entity.dart';

class ExamResultEntity {
  const ExamResultEntity({
    required this.attemptId,
    required this.examId,
    required this.status,
    required this.score,
    required this.totalMarks,
    required this.percentage,
    required this.passMarks,
    required this.passed,
    this.submittedAt,
    required this.feedback,
    required this.subjectBreakdown,
  });

  final String attemptId;
  final String examId;
  final ExamAttemptStatus status;
  final double score;
  final int totalMarks;
  final double percentage;
  final int passMarks;
  final bool passed;
  final DateTime? submittedAt;
  final List<ExamQuestionFeedbackEntity> feedback;
  final Map<String, ExamSubjectScoreEntity> subjectBreakdown;
}
