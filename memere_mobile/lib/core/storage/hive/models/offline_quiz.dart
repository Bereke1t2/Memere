import '../hive_json.dart';
import 'offline_question.dart';

/// A fully-downloaded quiz WITH answer keys, for offline on-device grading.
/// Parses the `GET /quizzes/:id/download` payload (dto.QuizDownloadResponse)
/// and round-trips through its own [toJson] when stored in the encrypted
/// `downloaded_quizzes` Hive box.
class OfflineQuiz {
  const OfflineQuiz({
    required this.id,
    required this.courseId,
    required this.title,
    required this.passPercentage,
    required this.questions,
    required this.downloadedAt,
    this.timeLimitSeconds,
    this.maxAttempts,
  });

  final String id;
  final String courseId;
  final String title;

  /// Pass threshold as a percentage (0–100); grader passes when pct ≥ this.
  final double passPercentage;
  final int? timeLimitSeconds;
  final int? maxAttempts;
  final List<OfflineQuestion> questions;
  final DateTime downloadedAt;

  factory OfflineQuiz.fromJson(Map<String, dynamic> json) => OfflineQuiz(
        id: hstr(json['id']),
        courseId: hstr(json['course_id']),
        title: hstr(json['title'], fallback: 'Quiz'),
        passPercentage: hdouble(json['pass_percentage']),
        timeLimitSeconds: hintOrNull(json['time_limit_seconds']),
        maxAttempts: hintOrNull(json['max_attempts']),
        questions: hlist(json['questions'], OfflineQuestion.fromJson),
        // Not on the wire; stamped on ingest, preserved on reload.
        downloadedAt: hdateOrNull(json['downloaded_at']) ?? DateTime.now(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'course_id': courseId,
        'title': title,
        'pass_percentage': passPercentage,
        'time_limit_seconds': timeLimitSeconds,
        'max_attempts': maxAttempts,
        'questions': questions.map((q) => q.toJson()).toList(),
        'downloaded_at': downloadedAt.toIso8601String(),
      };
}
