import '../../domain/entities/quiz_attempt_entity.dart';
import 'quiz_model_helpers.dart';
import 'quiz_question_model.dart';

class QuizAttemptModel extends QuizAttemptEntity {
  const QuizAttemptModel({
    required super.attemptId,
    required super.quizId,
    required super.attemptNumber,
    required super.status,
    required super.startedAt,
    super.expiresAt,
    super.remainingSeconds,
    required super.questions,
  });

  factory QuizAttemptModel.fromJson(Map<String, dynamic> json) {
    final questionsJson = json['questions'];
    final questions = questionsJson is List
        ? questionsJson
            .whereType<Map<String, dynamic>>()
            .map(QuizQuestionModel.fromJson)
            .toList()
        : <QuizQuestionModel>[];

    return QuizAttemptModel(
      attemptId: quizStringValue(json['attempt_id']),
      quizId: quizStringValue(json['quiz_id']),
      attemptNumber: quizIntValue(json['attempt_number'], fallback: 1),
      status: parseQuizAttemptStatus(quizStringValue(json['status'])),
      startedAt: quizDateValue(json['started_at']),
      expiresAt: quizDateValue(json['expires_at']),
      remainingSeconds: quizNullableInt(json['remaining_seconds']),
      questions: questions,
    );
  }

  Map<String, dynamic> toJson() => {
        'attempt_id': attemptId,
        'quiz_id': quizId,
        'attempt_number': attemptNumber,
        'status': quizAttemptStatusValue(status),
        'started_at': startedAt?.toIso8601String(),
        'expires_at': expiresAt?.toIso8601String(),
        'remaining_seconds': remainingSeconds,
        'questions': questions
            .whereType<QuizQuestionModel>()
            .map((question) => question.toJson())
            .toList(),
      };
}

QuizAttemptStatus parseQuizAttemptStatus(String value) {
  switch (value.toLowerCase()) {
    case 'submitted':
      return QuizAttemptStatus.submitted;
    case 'expired':
      return QuizAttemptStatus.expired;
    case 'graded':
      return QuizAttemptStatus.graded;
    default:
      return QuizAttemptStatus.inProgress;
  }
}

String quizAttemptStatusValue(QuizAttemptStatus status) {
  switch (status) {
    case QuizAttemptStatus.inProgress:
      return 'in_progress';
    case QuizAttemptStatus.submitted:
      return 'submitted';
    case QuizAttemptStatus.expired:
      return 'expired';
    case QuizAttemptStatus.graded:
      return 'graded';
  }
}
