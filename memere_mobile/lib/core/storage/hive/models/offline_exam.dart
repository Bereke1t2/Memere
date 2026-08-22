import '../hive_json.dart';
import 'offline_question.dart';

/// A fully-downloaded exam WITH answer keys, for offline on-device grading.
/// Parses the `GET /mock-exams/:id/download` payload (dto.ExamDownloadResponse)
/// and round-trips through its own [toJson] when stored in the encrypted
/// `downloaded_exams` Hive box.
///
/// `durationMinutes` drives a display-only client timer offline; `passMarks` is
/// the ABSOLUTE pass threshold (score ≥ passMarks), unlike the quiz percentage.
class OfflineExam {
  const OfflineExam({
    required this.id,
    required this.title,
    required this.subject,
    required this.grade,
    required this.durationMinutes,
    required this.totalMarks,
    required this.passMarks,
    required this.questions,
    required this.downloadedAt,
    this.courseId,
    this.instructions,
  });

  final String id;

  /// Nullable: standalone mock exams are not attached to a course.
  final String? courseId;
  final String title;
  final String subject;
  final int grade;
  final int durationMinutes;
  final int totalMarks;
  final int passMarks;
  final String? instructions;
  final List<OfflineQuestion> questions;
  final DateTime downloadedAt;

  factory OfflineExam.fromJson(Map<String, dynamic> json) => OfflineExam(
        id: hstr(json['id']),
        courseId: hstrOrNull(json['course_id']),
        title: hstr(json['title'], fallback: 'Exam'),
        subject: hstr(json['subject']),
        grade: hint(json['grade']),
        durationMinutes: hint(json['duration_minutes']),
        totalMarks: hint(json['total_marks']),
        passMarks: hint(json['pass_marks']),
        instructions: hstrOrNull(json['instructions']),
        questions: hlist(json['questions'], OfflineQuestion.fromJson),
        downloadedAt: hdateOrNull(json['downloaded_at']) ?? DateTime.now(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'course_id': courseId,
        'title': title,
        'subject': subject,
        'grade': grade,
        'duration_minutes': durationMinutes,
        'total_marks': totalMarks,
        'pass_marks': passMarks,
        'instructions': instructions,
        'questions': questions.map((q) => q.toJson()).toList(),
        'downloaded_at': downloadedAt.toIso8601String(),
      };
}
