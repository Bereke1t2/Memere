import '../../domain/entities/exam_attempt_entity.dart';
import 'exam_model_helpers.dart';
import 'exam_question_model.dart';

class ExamAttemptModel extends ExamAttemptEntity {
  const ExamAttemptModel({
    required super.attemptId,
    required super.examId,
    required super.status,
    super.startedAt,
    super.expiresAt,
    super.remainingSeconds,
    required super.totalMarks,
    required super.questions,
  });

  factory ExamAttemptModel.fromJson(Map<String, dynamic> json) {
    final questionsJson = json['questions'];
    final questions = questionsJson is List
        ? questionsJson
            .whereType<Map<String, dynamic>>()
            .map(ExamQuestionModel.fromJson)
            .toList()
        : <ExamQuestionModel>[];

    return ExamAttemptModel(
      attemptId: examStringValue(json['attempt_id']),
      examId: examStringValue(json['exam_id']),
      status: parseExamAttemptStatus(examStringValue(json['status'])),
      startedAt: examDateValue(json['started_at']),
      expiresAt: examDateValue(json['expires_at']),
      remainingSeconds: examNullableInt(json['remaining_seconds']),
      totalMarks: examIntValue(json['total_marks']),
      questions: questions,
    );
  }

  Map<String, dynamic> toJson() => {
        'attempt_id': attemptId,
        'exam_id': examId,
        'status': examAttemptStatusValue(status),
        'started_at': startedAt?.toIso8601String(),
        'expires_at': expiresAt?.toIso8601String(),
        'remaining_seconds': remainingSeconds,
        'total_marks': totalMarks,
        'questions': questions
            .whereType<ExamQuestionModel>()
            .map((question) => question.toJson())
            .toList(),
      };
}
