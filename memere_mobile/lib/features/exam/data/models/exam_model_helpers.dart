import '../../domain/entities/exam_attempt_entity.dart';
import '../../domain/entities/exam_question_entity.dart';

String examStringValue(Object? value, {String fallback = ''}) {
  if (value is String) return value;
  if (value == null) return fallback;
  return value.toString();
}

String? examNullableString(Object? value) {
  if (value is String && value.trim().isNotEmpty) return value;
  return null;
}

int examIntValue(Object? value, {int fallback = 0}) {
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value) ?? fallback;
  return fallback;
}

int? examNullableInt(Object? value) {
  if (value == null) return null;
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value);
  return null;
}

double examDoubleValue(Object? value, {double fallback = 0}) {
  if (value is int) return value.toDouble();
  if (value is double) return value;
  if (value is num) return value.toDouble();
  if (value is String) return double.tryParse(value) ?? fallback;
  return fallback;
}

double? examNullableDouble(Object? value) {
  if (value == null) return null;
  if (value is num) return value.toDouble();
  if (value is String) return double.tryParse(value);
  return null;
}

bool examBoolValue(Object? value, {bool fallback = false}) {
  if (value is bool) return value;
  if (value is String) return value.toLowerCase() == 'true';
  if (value is num) return value != 0;
  return fallback;
}

DateTime? examDateValue(Object? value) {
  if (value is! String || value.isEmpty) return null;
  return DateTime.tryParse(value);
}

List<String> examStringList(Object? value) {
  if (value is List) {
    return value
        .map((item) => examStringValue(item))
        .where((item) => item.trim().isNotEmpty)
        .toList();
  }
  return const [];
}

ExamAttemptStatus parseExamAttemptStatus(String value) {
  switch (value.toLowerCase()) {
    case 'submitted':
      return ExamAttemptStatus.submitted;
    case 'expired':
      return ExamAttemptStatus.expired;
    case 'graded':
      return ExamAttemptStatus.graded;
    default:
      return ExamAttemptStatus.inProgress;
  }
}

String examAttemptStatusValue(ExamAttemptStatus status) {
  switch (status) {
    case ExamAttemptStatus.inProgress:
      return 'in_progress';
    case ExamAttemptStatus.submitted:
      return 'submitted';
    case ExamAttemptStatus.expired:
      return 'expired';
    case ExamAttemptStatus.graded:
      return 'graded';
  }
}

ExamQuestionType parseExamQuestionType(String value) {
  switch (value.toLowerCase()) {
    case 'multiple_select':
      return ExamQuestionType.multipleSelect;
    case 'true_false':
      return ExamQuestionType.trueFalse;
    case 'short_answer':
      return ExamQuestionType.shortAnswer;
    default:
      return ExamQuestionType.multipleChoice;
  }
}

String examQuestionTypeValue(ExamQuestionType type) {
  switch (type) {
    case ExamQuestionType.multipleChoice:
      return 'multiple_choice';
    case ExamQuestionType.multipleSelect:
      return 'multiple_select';
    case ExamQuestionType.trueFalse:
      return 'true_false';
    case ExamQuestionType.shortAnswer:
      return 'short_answer';
  }
}
