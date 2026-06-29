import '../../domain/entities/quiz_entity.dart';
import 'quiz_model_helpers.dart';

class QuizModel extends QuizEntity {
  const QuizModel({
    required super.id,
    required super.courseId,
    required super.title,
    super.timeLimitSeconds,
    required super.passPercentage,
    required super.randomizeQuestions,
    super.maxAttempts,
    required super.questionCount,
    required super.attemptsUsed,
  });

  factory QuizModel.fromJson(Map<String, dynamic> json) {
    return QuizModel(
      id: quizStringValue(json['id']),
      courseId: quizStringValue(json['course_id']),
      title: quizStringValue(json['title'], fallback: 'Untitled quiz'),
      timeLimitSeconds: quizNullableInt(json['time_limit_seconds']),
      passPercentage: quizDoubleValue(json['pass_percentage']),
      randomizeQuestions: quizBoolValue(json['randomize_questions']),
      maxAttempts: quizNullableInt(json['max_attempts']),
      questionCount: quizIntValue(json['question_count']),
      attemptsUsed: quizIntValue(json['attempts_used']),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'course_id': courseId,
        'title': title,
        'time_limit_seconds': timeLimitSeconds,
        'pass_percentage': passPercentage,
        'randomize_questions': randomizeQuestions,
        'max_attempts': maxAttempts,
        'question_count': questionCount,
        'attempts_used': attemptsUsed,
      };
}
