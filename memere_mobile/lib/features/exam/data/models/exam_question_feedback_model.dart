import '../../domain/entities/exam_question_feedback_entity.dart';
import 'exam_feedback_answer_model.dart';
import 'exam_model_helpers.dart';

class ExamQuestionFeedbackModel extends ExamQuestionFeedbackEntity {
  const ExamQuestionFeedbackModel({
    required super.questionId,
    super.questionText,
    super.type,
    super.subject,
    super.topic,
    required super.correct,
    required super.marksAwarded,
    required super.marksPossible,
    required super.selectedAnswers,
    required super.correctAnswerIds,
    super.explanation,
    super.answers,
  });

  factory ExamQuestionFeedbackModel.fromJson(Map<String, dynamic> json) {
    final answersJson = json['answers'];
    final answers = answersJson is List
        ? answersJson
            .whereType<Map<String, dynamic>>()
            .map(ExamFeedbackAnswerModel.fromJson)
            .toList()
        : <ExamFeedbackAnswerModel>[];

    return ExamQuestionFeedbackModel(
      questionId: examStringValue(json['question_id']),
      questionText: examNullableString(json['question_text']),
      type: examNullableString(json['type']),
      subject: examNullableString(json['subject']),
      topic: examNullableString(json['topic']),
      correct: examBoolValue(json['correct']),
      marksAwarded: examIntValue(json['marks_awarded']),
      marksPossible: examIntValue(json['marks_possible']),
      selectedAnswers: examStringList(json['selected_answers']),
      correctAnswerIds: examStringList(json['correct_answer_ids']),
      explanation: examNullableString(json['explanation']),
      answers: answers,
    );
  }

  Map<String, dynamic> toJson() => {
        'question_id': questionId,
        'question_text': questionText,
        'type': type,
        'subject': subject,
        'topic': topic,
        'correct': correct,
        'marks_awarded': marksAwarded,
        'marks_possible': marksPossible,
        'selected_answers': selectedAnswers,
        'correct_answer_ids': correctAnswerIds,
        'explanation': explanation,
        'answers': answers
            .map((a) => a is ExamFeedbackAnswerModel
                ? a.toJson()
                : {
                    'id': a.id,
                    'text': a.text,
                    'is_correct': a.isCorrect,
                    'order_index': a.orderIndex,
                  })
            .toList(),
      };
}
