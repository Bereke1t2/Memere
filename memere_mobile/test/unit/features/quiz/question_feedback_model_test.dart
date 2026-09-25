import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/features/quiz/data/models/question_feedback_model.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_result_entity.dart';

void main() {
  group('QuestionFeedbackModel serialization', () {
    test('round-trips through fromJson and toJson including nested answers list', () {
      final json = <String, dynamic>{
        'question_id': 'q-101',
        'correct': true,
        'points_awarded': 2.0,
        'points_possible': 2,
        'selected_answers': ['a-1'],
        'correct_answer_ids': ['a-1'],
        'explanation': 'Correct explanation',
        'answers': [
          {'id': 'a-1', 'text': 'Option A', 'is_correct': true},
          {'id': 'a-2', 'text': 'Option B', 'is_correct': false},
        ],
      };

      final model = QuestionFeedbackModel.fromJson(json);
      expect(model.questionId, equals('q-101'));
      expect(model.correct, isTrue);
      expect(model.pointsAwarded, equals(2.0));
      expect(model.pointsPossible, equals(2));
      expect(model.selectedAnswers, equals(['a-1']));
      expect(model.correctAnswerIds, equals(['a-1']));
      expect(model.explanation, equals('Correct explanation'));
      expect(model.answers.length, equals(2));
      expect(model.answers.first.text, equals('Option A'));

      final encoded = model.toJson();
      expect(encoded['question_id'], equals('q-101'));
      expect(encoded['correct'], isTrue);
      expect(encoded['answers'], isA<List>());
      final encodedAnswers = encoded['answers'] as List;
      expect(encodedAnswers.length, equals(2));
      expect((encodedAnswers[0] as Map)['text'], equals('Option A'));
    });
  });
}
