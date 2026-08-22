import '../hive_json.dart';

/// Question types the backend emits. Mirrors `entity.QuestionType`
/// (`multiple_choice`, `true_false`, `short_answer`) — no others exist.
enum OfflineQuestionType { multipleChoice, trueFalse, shortAnswer, unknown }

OfflineQuestionType offlineQuestionTypeFromString(String value) {
  switch (value.toLowerCase().replaceAll('_', '')) {
    case 'multiplechoice':
      return OfflineQuestionType.multipleChoice;
    case 'truefalse':
      return OfflineQuestionType.trueFalse;
    case 'shortanswer':
      return OfflineQuestionType.shortAnswer;
    default:
      return OfflineQuestionType.unknown;
  }
}

/// The backend wire form (snake_case), so re-serializing round-trips exactly.
String offlineQuestionTypeWire(OfflineQuestionType t) {
  switch (t) {
    case OfflineQuestionType.multipleChoice:
      return 'multiple_choice';
    case OfflineQuestionType.trueFalse:
      return 'true_false';
    case OfflineQuestionType.shortAnswer:
      return 'short_answer';
    case OfflineQuestionType.unknown:
      return 'unknown';
  }
}

/// One answer option, INCLUDING whether it is correct. Present only inside the
/// encrypted download boxes — never in an online-play projection.
class OfflineAnswer {
  const OfflineAnswer({
    required this.id,
    required this.text,
    required this.orderIndex,
    required this.isCorrect,
  });

  final String id;
  final String text;
  final int orderIndex;
  final bool isCorrect;

  factory OfflineAnswer.fromJson(Map<String, dynamic> json) => OfflineAnswer(
        id: hstr(json['id']),
        text: hstr(json['text']),
        orderIndex: hint(json['order_index']),
        isCorrect: hbool(json['is_correct']),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'text': text,
        'order_index': orderIndex,
        'is_correct': isCorrect,
      };
}

/// One question with its full answer key + explanation. `points` carries the
/// per-question weight (quiz points or exam marks — the backend maps both here).
class OfflineQuestion {
  const OfflineQuestion({
    required this.id,
    required this.text,
    required this.type,
    required this.points,
    required this.orderIndex,
    required this.answers,
    this.subject,
    this.topic,
    this.explanation,
  });

  final String id;
  final String text;
  final OfflineQuestionType type;
  final int points;
  final int orderIndex;
  final List<OfflineAnswer> answers;
  final String? subject;
  final String? topic;
  final String? explanation;

  /// The IDs of every option flagged correct (the on-device grading key).
  List<String> get correctAnswerIds =>
      answers.where((a) => a.isCorrect).map((a) => a.id).toList();

  factory OfflineQuestion.fromJson(Map<String, dynamic> json) => OfflineQuestion(
        id: hstr(json['id']),
        text: hstr(json['text']),
        type: offlineQuestionTypeFromString(hstr(json['type'])),
        points: hint(json['points']),
        orderIndex: hint(json['order_index']),
        subject: hstrOrNull(json['subject']),
        topic: hstrOrNull(json['topic']),
        explanation: hstrOrNull(json['explanation']),
        answers: hlist(json['answers'], OfflineAnswer.fromJson),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'text': text,
        'type': offlineQuestionTypeWire(type),
        'points': points,
        'order_index': orderIndex,
        'subject': subject,
        'topic': topic,
        'explanation': explanation,
        'answers': answers.map((a) => a.toJson()).toList(),
      };
}
