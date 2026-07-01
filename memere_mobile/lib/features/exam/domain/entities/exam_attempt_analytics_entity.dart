import 'exam_subject_score_entity.dart';

class ExamAttemptAnalyticsEntity {
  const ExamAttemptAnalyticsEntity({
    required this.attemptId,
    required this.examId,
    required this.score,
    required this.percentage,
    required this.subjectBreakdown,
    required this.weakAreas,
    this.percentile,
  });

  final String attemptId;
  final String examId;
  final double score;
  final double percentage;
  final List<ExamSubjectScoreEntity> subjectBreakdown;
  final List<ExamSubjectScoreEntity> weakAreas;
  final double? percentile;
}
