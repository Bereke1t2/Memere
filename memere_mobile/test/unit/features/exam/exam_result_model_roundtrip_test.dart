import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/core/grading/offline_grader.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_exam.dart';
import 'package:memere_mobile/core/storage/hive/models/offline_question.dart';
import 'package:memere_mobile/features/exam/data/models/exam_result_model.dart';
import 'package:memere_mobile/features/exam/domain/entities/exam_attempt_entity.dart';

/// Guards the persistence seam added for offline exams: an on-device-graded
/// [ExamResultModel] must survive `toJson` → box → `fromJson` intact so it
/// renders on the existing result screen exactly as the server-graded one would.
void main() {
  final exam = OfflineExam(
    id: 'exam-1',
    courseId: 'course-1',
    title: 'Mock',
    subject: 'Math',
    grade: 10,
    durationMinutes: 45,
    totalMarks: 30,
    passMarks: 15,
    downloadedAt: DateTime.parse('2026-08-01T00:00:00.000Z'),
    questions: const [
      OfflineQuestion(
        id: 'q1',
        text: '2+2',
        type: OfflineQuestionType.multipleChoice,
        points: 20,
        orderIndex: 0,
        subject: 'Arithmetic',
        answers: [
          OfflineAnswer(id: 'a1', text: '4', orderIndex: 0, isCorrect: true),
          OfflineAnswer(id: 'a2', text: '5', orderIndex: 1, isCorrect: false),
        ],
      ),
      OfflineQuestion(
        id: 'q2',
        text: 'capital of Ethiopia?',
        type: OfflineQuestionType.shortAnswer,
        points: 10,
        orderIndex: 1,
        subject: 'Geography',
        answers: [
          OfflineAnswer(id: 'a3', text: 'Addis', orderIndex: 0, isCorrect: true),
        ],
      ),
    ],
  );

  test('ExamResultModel.toJson round-trips through fromJson', () {
    final result = OfflineGrader.gradeExam(
      exam: exam,
      submission: {
        'q1': const OfflineAnswerInput(selectedAnswerIds: ['a1']), // correct
        'q2': const OfflineAnswerInput(text: 'wrong'), // wrong
      },
      attemptId: 'local-1',
      takenAt: DateTime.parse('2026-08-21T10:00:00.000Z'),
    );

    final restored = ExamResultModel.fromJson(result.toJson());

    expect(restored.attemptId, 'local-1');
    expect(restored.examId, 'exam-1');
    expect(restored.status, ExamAttemptStatus.graded);
    expect(restored.score, result.score); // 20 of 30
    expect(restored.totalMarks, 30);
    expect(restored.percentage, result.percentage);
    expect(restored.passMarks, 15);
    expect(restored.passed, result.passed); // 20 >= 15
    expect(restored.submittedAt!.toUtc(),
        DateTime.parse('2026-08-21T10:00:00.000Z'));

    // Per-question feedback preserved, including type + mark fields.
    expect(restored.feedback, hasLength(2));
    final q1 = restored.feedback.firstWhere((f) => f.questionId == 'q1');
    expect(q1.correct, isTrue);
    expect(q1.marksAwarded, 20);
    expect(q1.type, 'multiple_choice');
    final q2 = restored.feedback.firstWhere((f) => f.questionId == 'q2');
    expect(q2.correct, isFalse);
    expect(q2.marksAwarded, 0);
    expect(q2.type, 'short_answer');

    // Subject breakdown preserved by key (result-shape map → fromMapEntry).
    expect(restored.subjectBreakdown.keys.toSet(), {'Arithmetic', 'Geography'});
    expect(restored.subjectBreakdown['Arithmetic']!.earned, 20);
    expect(restored.subjectBreakdown['Arithmetic']!.possible, 20);
    expect(restored.subjectBreakdown['Geography']!.earned, 0);
    expect(restored.subjectBreakdown['Geography']!.possible, 10);
  });
}
