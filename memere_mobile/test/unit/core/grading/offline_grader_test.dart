import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/core/grading/offline_grader.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_exam.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_question.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_quiz.dart';
import 'package:memere_mobile/features/exam/domain/entities/exam_attempt_entity.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_attempt_entity.dart';

/// A single-correct multiple-choice question worth [points], answer "A" correct
/// and "B" wrong — mirrors the backend grading_test.go `mcQuestion` helper.
OfflineQuestion _mc(
  String id,
  int points,
  String subject, {
  required String correctId,
  required String wrongId,
}) {
  return OfflineQuestion(
    id: id,
    text: 'q-$id',
    type: OfflineQuestionType.multipleChoice,
    points: points,
    orderIndex: 0,
    subject: subject,
    answers: [
      OfflineAnswer(id: correctId, text: 'A', orderIndex: 0, isCorrect: true),
      OfflineAnswer(id: wrongId, text: 'B', orderIndex: 1, isCorrect: false),
    ],
  );
}

OfflineQuiz _quiz(List<OfflineQuestion> qs, {double passPct = 50}) => OfflineQuiz(
      id: 'quiz-1',
      courseId: 'course-1',
      title: 'Quiz',
      passPercentage: passPct,
      questions: qs,
      downloadedAt: DateTime(2026, 1, 1),
    );

void main() {
  group('OfflineGrader.gradeQuiz — mirrors grading.Grade', () {
    test('scores, percentage, pass, and subject breakdown', () {
      final q1 = _mc('q1', 5, 'Math', correctId: 'a1', wrongId: 'b1');
      final q2 = _mc('q2', 3, 'Science', correctId: 'a2', wrongId: 'b2');

      final res = OfflineGrader.gradeQuiz(
        quiz: _quiz([q1, q2]),
        submission: {
          'q1': const OfflineAnswerInput(selectedAnswerIds: ['a1']),
          'q2': const OfflineAnswerInput(selectedAnswerIds: ['a2']),
        },
        attemptId: 'local-1',
      );

      expect(res.score, 8);
      expect(res.totalPoints, 8);
      expect(res.percentage, 100);
      expect(res.passed, isTrue);
      expect(res.status, QuizAttemptStatus.graded);
      expect(res.subjectBreakdown['Math']!.earned, 5);
      expect(res.subjectBreakdown['Science']!.earned, 3);
      // Feedback carries the answer key + explanation shape the screen renders.
      expect(res.feedback, hasLength(2));
      expect(res.feedback.first.correctAnswerIds, ['a1']);
      expect(res.feedback.first.answers, hasLength(2));
    });

    test('wrong answer awards zero for that question only', () {
      final q1 = _mc('q1', 5, 'Math', correctId: 'a1', wrongId: 'b1');
      final q2 = _mc('q2', 3, 'Science', correctId: 'a2', wrongId: 'b2');

      final res = OfflineGrader.gradeQuiz(
        quiz: _quiz([q1, q2]),
        submission: {
          'q1': const OfflineAnswerInput(selectedAnswerIds: ['b1']), // wrong
          'q2': const OfflineAnswerInput(selectedAnswerIds: ['a2']),
        },
        attemptId: 'local-1',
      );

      expect(res.score, 3);
      expect(res.percentage, closeTo(37.5, 1e-9));
      expect(res.passed, isFalse); // 37.5% < 50%
    });

    test('short answer is normalized (trim/case/whitespace)', () {
      const q = OfflineQuestion(
        id: 'sa',
        text: 'capital?',
        type: OfflineQuestionType.shortAnswer,
        points: 2,
        orderIndex: 0,
        answers: [
          OfflineAnswer(id: 'x', text: 'Paris', orderIndex: 0, isCorrect: true),
        ],
      );

      final res = OfflineGrader.gradeQuiz(
        quiz: _quiz([q]),
        submission: {'sa': const OfflineAnswerInput(text: '  pARIs  ')},
        attemptId: 'local-1',
      );

      expect(res.score, 2);
    });

    test('multi-correct requires exact set (no partial, no extra)', () {
      OfflineQuestion multi() => const OfflineQuestion(
            id: 'm',
            text: 'pick two',
            type: OfflineQuestionType.multipleChoice,
            points: 4,
            orderIndex: 0,
            answers: [
              OfflineAnswer(id: 'a', text: 'A', orderIndex: 0, isCorrect: true),
              OfflineAnswer(id: 'b', text: 'B', orderIndex: 1, isCorrect: true),
              OfflineAnswer(id: 'c', text: 'C', orderIndex: 2, isCorrect: false),
            ],
          );

      double scoreFor(List<String> selected) => OfflineGrader.gradeQuiz(
            quiz: _quiz([multi()]),
            submission: {'m': OfflineAnswerInput(selectedAnswerIds: selected)},
            attemptId: 'local-1',
          ).score;

      expect(scoreFor(['a']), 0); // partial
      expect(scoreFor(['a', 'b']), 4); // exact
      expect(scoreFor(['a', 'b', 'c']), 0); // over-selection
    });

    test('unanswered question scores zero', () {
      final q1 = _mc('q1', 5, 'Math', correctId: 'a1', wrongId: 'b1');
      final res = OfflineGrader.gradeQuiz(
        quiz: _quiz([q1]),
        submission: const {},
        attemptId: 'local-1',
      );
      expect(res.score, 0);
      expect(res.percentage, 0);
    });
  });

  group('OfflineGrader.gradeExam — absolute pass threshold', () {
    OfflineExam exam(List<OfflineQuestion> qs, {required int passMarks}) => OfflineExam(
          id: 'exam-1',
          title: 'Mock',
          subject: 'Aptitude',
          grade: 12,
          durationMinutes: 60,
          totalMarks: 8,
          passMarks: passMarks,
          questions: qs,
          downloadedAt: DateTime(2026, 1, 1),
        );

    test('passes on score >= passMarks and carries feedback type/subject', () {
      final q1 = _mc('q1', 5, 'Math', correctId: 'a1', wrongId: 'b1');
      final q2 = _mc('q2', 3, 'Science', correctId: 'a2', wrongId: 'b2');

      final res = OfflineGrader.gradeExam(
        exam: exam([q1, q2], passMarks: 5),
        submission: {
          'q1': const OfflineAnswerInput(selectedAnswerIds: ['a1']),
          'q2': const OfflineAnswerInput(selectedAnswerIds: ['a2']),
        },
        attemptId: 'local-1',
      );

      expect(res.score, 8);
      expect(res.totalMarks, 8);
      expect(res.passMarks, 5);
      expect(res.passed, isTrue);
      expect(res.status, ExamAttemptStatus.graded);
      expect(res.feedback.first.type, 'multiple_choice');
      expect(res.feedback.first.subject, 'Math');
      expect(res.subjectBreakdown['Math']!.earned, 5);
    });

    test('fails when score below absolute passMarks even if > 0', () {
      final q1 = _mc('q1', 5, 'Math', correctId: 'a1', wrongId: 'b1');
      final q2 = _mc('q2', 3, 'Science', correctId: 'a2', wrongId: 'b2');

      final res = OfflineGrader.gradeExam(
        exam: exam([q1, q2], passMarks: 5),
        submission: {
          'q1': const OfflineAnswerInput(selectedAnswerIds: ['b1']), // wrong → 3 marks
          'q2': const OfflineAnswerInput(selectedAnswerIds: ['a2']),
        },
        attemptId: 'local-1',
      );

      expect(res.score, 3);
      expect(res.passed, isFalse); // 3 < 5 absolute, despite 37.5%
    });
  });

  group('offlineSubmissionFromRaw', () {
    test('interprets values by question type', () {
      final questions = [
        _mc('q1', 1, 'Math', correctId: 'a1', wrongId: 'b1'),
        const OfflineQuestion(
          id: 'sa',
          text: 'capital?',
          type: OfflineQuestionType.shortAnswer,
          points: 1,
          orderIndex: 1,
          answers: [
            OfflineAnswer(id: 'x', text: 'Paris', orderIndex: 0, isCorrect: true),
          ],
        ),
      ];

      final sub = offlineSubmissionFromRaw(questions, {
        'q1': ['a1'], // list of ids
        'sa': 'Paris', // free text
        'ghost': 'ignored', // unknown question dropped
      });

      expect(sub['q1']!.selectedAnswerIds, ['a1']);
      expect(sub['sa']!.text, 'Paris');
      expect(sub.containsKey('ghost'), isFalse);
    });

    test('accepts a single string id (not just a list)', () {
      final questions = [_mc('q1', 1, 'Math', correctId: 'a1', wrongId: 'b1')];
      final sub = offlineSubmissionFromRaw(questions, {'q1': 'a1'});
      expect(sub['q1']!.selectedAnswerIds, ['a1']);
    });
  });
}
