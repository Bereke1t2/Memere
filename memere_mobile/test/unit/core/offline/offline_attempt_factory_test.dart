import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/core/offline/offline_attempt_factory.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_exam.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_question.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_quiz.dart';
import 'package:memere_mobile/features/exam/domain/entities/exam_attempt_entity.dart';
import 'package:memere_mobile/features/exam/domain/entities/exam_question_entity.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_attempt_entity.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_question_entity.dart';

OfflineAnswer _ans(String id,
        {int order = 0, bool correct = false, String text = 'opt'}) =>
    OfflineAnswer(id: id, text: text, orderIndex: order, isCorrect: correct);

OfflineQuestion _q({
  required String id,
  required OfflineQuestionType type,
  int points = 1,
  int order = 0,
  List<OfflineAnswer> answers = const [],
}) =>
    OfflineQuestion(
      id: id,
      text: 'q-$id',
      type: type,
      points: points,
      orderIndex: order,
      answers: answers,
    );

void main() {
  group('buildLocalQuizAttempt', () {
    final quiz = OfflineQuiz(
      id: 'quiz-1',
      courseId: 'course-1',
      title: 'Q',
      passPercentage: 50,
      downloadedAt: DateTime.parse('2026-08-01T00:00:00.000Z'),
      questions: [
        _q(id: 'b', type: OfflineQuestionType.trueFalse, order: 2, answers: [
          _ans('b2', order: 1, correct: true),
          _ans('b1', order: 0),
        ]),
        _q(
          id: 'a',
          type: OfflineQuestionType.multipleChoice,
          order: 1,
          answers: [_ans('a1', order: 0, correct: true)],
        ),
        _q(id: 'c', type: OfflineQuestionType.shortAnswer, order: 3),
      ],
    );

    test('sets local attempt metadata and is untimed', () {
      final attempt = buildLocalQuizAttempt(quiz, attemptId: 'local-1');
      expect(attempt.attemptId, 'local-1');
      expect(attempt.quizId, 'quiz-1');
      expect(attempt.attemptNumber, 1);
      expect(attempt.status, QuizAttemptStatus.inProgress);
      expect(attempt.startedAt, isNotNull);
      // Quizzes are untimed offline.
      expect(attempt.expiresAt, isNull);
      expect(attempt.remainingSeconds, isNull);
      expect(attempt.hasTimer, isFalse);
    });

    test('orders questions and answers by orderIndex and maps types', () {
      final attempt = buildLocalQuizAttempt(quiz, attemptId: 'local-1');
      expect(attempt.questions.map((q) => q.id).toList(), ['a', 'b', 'c']);
      expect(attempt.questions[0].type, QuizQuestionType.multipleChoice);
      expect(attempt.questions[1].type, QuizQuestionType.trueFalse);
      expect(attempt.questions[2].type, QuizQuestionType.shortAnswer);
      // Answers sorted by orderIndex within a question.
      expect(attempt.questions[1].answers.map((a) => a.id).toList(),
          ['b1', 'b2']);
    });
  });

  group('buildLocalExamAttempt', () {
    OfflineExam examWith({required int duration}) => OfflineExam(
          id: 'exam-1',
          courseId: null,
          title: 'E',
          subject: 'Math',
          grade: 10,
          durationMinutes: duration,
          totalMarks: 30,
          passMarks: 15,
          downloadedAt: DateTime.parse('2026-08-01T00:00:00.000Z'),
          questions: [
            _q(
                id: 'x2',
                type: OfflineQuestionType.unknown,
                order: 5,
                points: 10),
            _q(
                id: 'x1',
                type: OfflineQuestionType.multipleChoice,
                order: 1,
                points: 20),
          ],
        );

    test('timed exam derives expiresAt/remainingSeconds from duration', () {
      final started = DateTime.parse('2026-08-21T09:00:00.000Z');
      final attempt = buildLocalExamAttempt(
        examWith(duration: 30),
        attemptId: 'local-e',
        startedAt: started,
      );
      expect(attempt.attemptId, 'local-e');
      expect(attempt.examId, 'exam-1');
      expect(attempt.status, ExamAttemptStatus.inProgress);
      expect(attempt.totalMarks, 30);
      expect(attempt.remainingSeconds, 30 * 60);
      expect(attempt.expiresAt, started.add(const Duration(minutes: 30)));
      // Questions sorted; unknown type degrades to multiple choice.
      expect(attempt.questions.map((q) => q.questionId).toList(), ['x1', 'x2']);
      expect(attempt.questions[0].marks, 20);
      expect(attempt.questions[0].type, ExamQuestionType.multipleChoice);
      expect(attempt.questions[1].type, ExamQuestionType.multipleChoice);
    });

    test('zero-duration exam is untimed', () {
      final attempt =
          buildLocalExamAttempt(examWith(duration: 0), attemptId: 'local-e');
      expect(attempt.expiresAt, isNull);
      expect(attempt.remainingSeconds, isNull);
      expect(attempt.hasTimer, isFalse);
    });
  });
}
