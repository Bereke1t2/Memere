import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/core/grading/offline_grader.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_question.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_quiz.dart';
import 'package:memere_mobile/features/quiz/data/models/quiz_result_model.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_attempt_entity.dart';

/// Guards the persistence seam for quiz results: an on-device-graded
/// [QuizResultModel] must survive `toJson` → box → `fromJson` intact so it
/// renders on the result screen with answer choices and texts intact.
void main() {
  final quiz = OfflineQuiz(
    id: 'quiz-1',
    courseId: 'course-1',
    title: 'Algebra Basics',
    passPercentage: 70,
    timeLimitSeconds: 600,
    maxAttempts: 3,
    downloadedAt: DateTime.parse('2026-08-01T00:00:00.000Z'),
    questions: const [
      OfflineQuestion(
        id: 'q1',
        text: 'What is 3 + 5?',
        type: OfflineQuestionType.multipleChoice,
        points: 10,
        orderIndex: 0,
        subject: 'Math',
        explanation: '3 + 5 = 8.',
        answers: [
          OfflineAnswer(id: 'a1', text: '8', orderIndex: 0, isCorrect: true),
          OfflineAnswer(id: 'a2', text: '9', orderIndex: 1, isCorrect: false),
          OfflineAnswer(id: 'a3', text: '7', orderIndex: 2, isCorrect: false),
        ],
      ),
      OfflineQuestion(
        id: 'q2',
        text: 'Is 2 an even prime number?',
        type: OfflineQuestionType.trueFalse,
        points: 5,
        orderIndex: 1,
        subject: 'Math',
        explanation: '2 is the only even prime.',
        answers: [
          OfflineAnswer(id: 'a4', text: 'True', orderIndex: 0, isCorrect: true),
          OfflineAnswer(id: 'a5', text: 'False', orderIndex: 1, isCorrect: false),
        ],
      ),
    ],
  );

  test('QuizResultModel.toJson round-trips through fromJson with answers preserved', () {
    final result = OfflineGrader.gradeQuiz(
      quiz: quiz,
      submission: {
        'q1': const OfflineAnswerInput(selectedAnswerIds: ['a1']), // correct
        'q2': const OfflineAnswerInput(selectedAnswerIds: ['a5']), // wrong (False selected)
      },
      attemptId: 'local-quiz-1',
      takenAt: DateTime.parse('2026-08-25T14:00:00.000Z'),
    );

    final json = result.toJson();
    final restored = QuizResultModel.fromJson(json);

    expect(restored.attemptId, 'local-quiz-1');
    expect(restored.quizId, 'quiz-1');
    expect(restored.status, QuizAttemptStatus.graded);
    expect(restored.score, 10.0);
    expect(restored.totalPoints, 15);
    expect(restored.percentage, closeTo(66.666, 0.01));
    expect(restored.passed, isFalse); // 66.66% < 70%
    expect(restored.submittedAt!.toUtc(), DateTime.parse('2026-08-25T14:00:00.000Z'));

    // Question 1 feedback
    expect(restored.feedback, hasLength(2));
    final q1 = restored.feedback.firstWhere((f) => f.questionId == 'q1');
    expect(q1.questionText, 'What is 3 + 5?');
    expect(q1.correct, isTrue);
    expect(q1.pointsAwarded, 10);
    expect(q1.pointsPossible, 10);
    expect(q1.selectedAnswers, ['a1']);
    expect(q1.correctAnswerIds, ['a1']);
    expect(q1.explanation, '3 + 5 = 8.');
    expect(q1.answers, hasLength(3));
    expect(q1.answers[0].id, 'a1');
    expect(q1.answers[0].text, '8');
    expect(q1.answers[0].isCorrect, isTrue);
    expect(q1.answers[1].id, 'a2');
    expect(q1.answers[1].text, '9');
    expect(q1.answers[1].isCorrect, isFalse);

    // Question 2 feedback
    final q2 = restored.feedback.firstWhere((f) => f.questionId == 'q2');
    expect(q2.questionText, 'Is 2 an even prime number?');
    expect(q2.correct, isFalse);
    expect(q2.pointsAwarded, 0);
    expect(q2.pointsPossible, 5);
    expect(q2.selectedAnswers, ['a5']);
    expect(q2.correctAnswerIds, ['a4']);
    expect(q2.answers, hasLength(2));
    expect(q2.answers[0].text, 'True');
    expect(q2.answers[0].isCorrect, isTrue);
    expect(q2.answers[1].text, 'False');
    expect(q2.answers[1].isCorrect, isFalse);

    // Subject breakdown
    expect(restored.subjectBreakdown.containsKey('Math'), isTrue);
    expect(restored.subjectBreakdown['Math']!.earned, 10);
    expect(restored.subjectBreakdown['Math']!.possible, 15);
  });
}
