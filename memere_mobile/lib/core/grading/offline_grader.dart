import '../../features/exam/data/models/exam_feedback_answer_model.dart';
import '../../features/exam/data/models/exam_question_feedback_model.dart';
import '../../features/exam/data/models/exam_result_model.dart';
import '../../features/exam/data/models/exam_subject_score_model.dart';
import '../../features/exam/domain/entities/exam_attempt_entity.dart';
import '../../features/quiz/data/models/question_feedback_model.dart';
import '../../features/quiz/data/models/quiz_result_model.dart';
import '../../features/quiz/data/models/subject_score_model.dart';
import '../../features/quiz/domain/entities/quiz_attempt_entity.dart';
import '../storage/hive/models/offline_exam.dart';
import '../storage/hive/models/offline_question.dart';
import '../storage/hive/models/offline_quiz.dart';

/// A single question's response, normalized for grading: the selected option
/// IDs (choice/true-false) and/or the free text (short answer).
class OfflineAnswerInput {
  const OfflineAnswerInput({
    this.selectedAnswerIds = const <String>[],
    this.text,
  });

  final List<String> selectedAnswerIds;
  final String? text;
}

/// On-device port of the backend `internal/usecase/grading/grading.go`.
///
/// It reproduces the authoritative semantics EXACTLY so an offline score is a
/// faithful preview of what the server will return on reconnect:
///   - choice / true-false: full marks only when the SET of selected option IDs
///     equals the SET of correct option IDs (no missing, no extra);
///   - short answer: normalized (trim → lowercase → collapse whitespace) text
///     matches any option flagged correct;
///   - a question is awarded its full marks or zero — no partial credit;
///   - quiz pass = percentage ≥ passPercentage; exam pass = score ≥ passMarks
///     (absolute), matching the split the backend leaves to its callers.
///
/// It emits the very same [QuizResultModel] / [ExamResultModel] the online path
/// builds from the server response, so the result screens render offline with
/// zero change. Scores are provisional until the server re-grades on sync.
class OfflineGrader {
  const OfflineGrader._();

  static QuizResultModel gradeQuiz({
    required OfflineQuiz quiz,
    required Map<String, OfflineAnswerInput> submission,
    required String attemptId,
    DateTime? takenAt,
  }) {
    final scored = _grade(quiz.questions, submission);
    final pct = scored.total > 0 ? scored.earned / scored.total * 100 : 0.0;

    final feedback = scored.outcomes
        .map((o) => QuestionFeedbackModel(
              questionId: o.question.id,
              questionText: o.question.text,
              correct: o.correct,
              pointsAwarded: o.awarded,
              pointsPossible: o.question.points,
              selectedAnswers: o.selected,
              correctAnswerIds: o.question.correctAnswerIds,
              explanation: o.question.explanation,
              answers: o.question.answers
                  .map((a) => QuizAnswerFeedbackModel(
                        id: a.id,
                        text: a.text,
                        isCorrect: a.isCorrect,
                      ))
                  .toList(),
            ))
        .toList();

    final breakdown = scored.breakdown.map(
      (k, v) => MapEntry(k, SubjectScoreModel(earned: v.earned, possible: v.possible)),
    );

    return QuizResultModel(
      attemptId: attemptId,
      quizId: quiz.id,
      attemptNumber: 1,
      status: QuizAttemptStatus.graded,
      score: scored.earned.toDouble(),
      totalPoints: scored.total,
      percentage: pct,
      passed: pct >= quiz.passPercentage,
      submittedAt: takenAt ?? DateTime.now(),
      feedback: feedback,
      subjectBreakdown: breakdown,
    );
  }

  static ExamResultModel gradeExam({
    required OfflineExam exam,
    required Map<String, OfflineAnswerInput> submission,
    required String attemptId,
    DateTime? takenAt,
  }) {
    final scored = _grade(exam.questions, submission);
    final pct = scored.total > 0 ? scored.earned / scored.total * 100 : 0.0;

    final feedback = scored.outcomes
        .map((o) => ExamQuestionFeedbackModel(
              questionId: o.question.id,
              questionText: o.question.text,
              type: offlineQuestionTypeWire(o.question.type),
              subject: o.question.subject,
              topic: o.question.topic,
              correct: o.correct,
              marksAwarded: o.awarded,
              marksPossible: o.question.points,
              selectedAnswers: o.selected,
              correctAnswerIds: o.question.correctAnswerIds,
              explanation: o.question.explanation,
              answers: o.question.answers
                  .map((a) => ExamFeedbackAnswerModel(
                        id: a.id,
                        text: a.text,
                        isCorrect: a.isCorrect,
                        orderIndex: a.orderIndex,
                      ))
                  .toList(),
            ))
        .toList();

    final breakdown = scored.breakdown.map(
      (k, v) => MapEntry(k, ExamSubjectScoreModel(key: k, earned: v.earned, possible: v.possible)),
    );

    return ExamResultModel(
      attemptId: attemptId,
      examId: exam.id,
      status: ExamAttemptStatus.graded,
      score: scored.earned.toDouble(),
      totalMarks: scored.total,
      percentage: pct,
      passMarks: exam.passMarks,
      passed: scored.earned >= exam.passMarks,
      submittedAt: takenAt ?? DateTime.now(),
      feedback: feedback,
      subjectBreakdown: breakdown,
    );
  }

  // ── core (mirrors grading.Grade) ───────────────────────────────────────────

  static _Scored _grade(
    List<OfflineQuestion> questions,
    Map<String, OfflineAnswerInput> sub,
  ) {
    var earned = 0;
    var total = 0;
    final outcomes = <_Outcome>[];
    final breakdown = <String, _Tally>{};

    for (final q in questions) {
      total += q.points;
      final input = sub[q.id] ?? const OfflineAnswerInput();
      final correct = _isResponseCorrect(q, input);
      final awarded = correct ? q.points : 0;
      if (correct) earned += q.points;

      outcomes.add(_Outcome(
        question: q,
        correct: correct,
        awarded: awarded,
        selected: input.selectedAnswerIds,
      ));

      final key = _breakdownKey(q);
      if (key.isNotEmpty) {
        final tally = breakdown[key] ?? _Tally();
        tally.earned += awarded;
        tally.possible += q.points;
        breakdown[key] = tally;
      }
    }

    return _Scored(earned: earned, total: total, outcomes: outcomes, breakdown: breakdown);
  }

  static bool _isResponseCorrect(OfflineQuestion q, OfflineAnswerInput input) {
    if (q.type == OfflineQuestionType.shortAnswer) {
      final got = _normalize(input.text ?? '');
      if (got.isEmpty) return false;
      for (final a in q.answers) {
        if (a.isCorrect && _normalize(a.text) == got) return true;
      }
      return false;
    }
    // multiple_choice, true_false (and unknown): exact set equality.
    final want = q.answers.where((a) => a.isCorrect).map((a) => a.id).toSet();
    final got = input.selectedAnswerIds.toSet();
    if (want.isEmpty || want.length != got.length) return false;
    return want.every(got.contains);
  }

  static String _breakdownKey(OfflineQuestion q) {
    if (q.subject != null && q.subject!.isNotEmpty) return q.subject!;
    if (q.topic != null && q.topic!.isNotEmpty) return q.topic!;
    return '';
  }

  /// trims, lowercases, and collapses internal whitespace runs to a single
  /// space — matching Go's `strings.Join(strings.Fields(strings.ToLower(s)), " ")`.
  static String _normalize(String s) =>
      s.toLowerCase().split(RegExp(r'\s+')).where((e) => e.isNotEmpty).join(' ');
}

/// Interprets a persisted raw-answer map (questionId → value) using each
/// question's type: short-answer values are text, everything else is option
/// ID(s). Lets a queued [rawAnswers] payload be graded and re-submitted from one
/// canonical shape.
Map<String, OfflineAnswerInput> offlineSubmissionFromRaw(
  List<OfflineQuestion> questions,
  Map<String, dynamic> raw,
) {
  final byId = {for (final q in questions) q.id: q};
  final out = <String, OfflineAnswerInput>{};
  raw.forEach((qid, value) {
    final q = byId[qid];
    if (q == null) return;
    if (q.type == OfflineQuestionType.shortAnswer) {
      out[qid] = OfflineAnswerInput(text: value?.toString());
    } else {
      out[qid] = OfflineAnswerInput(selectedAnswerIds: _asIdList(value));
    }
  });
  return out;
}

List<String> _asIdList(dynamic value) {
  if (value == null) return const <String>[];
  if (value is List) {
    return value.map((e) => e.toString()).where((e) => e.isNotEmpty).toList();
  }
  final s = value.toString();
  return s.isEmpty ? const <String>[] : <String>[s];
}

class _Outcome {
  _Outcome({
    required this.question,
    required this.correct,
    required this.awarded,
    required this.selected,
  });

  final OfflineQuestion question;
  final bool correct;
  final int awarded;
  final List<String> selected;
}

class _Tally {
  int earned = 0;
  int possible = 0;
}

class _Scored {
  _Scored({
    required this.earned,
    required this.total,
    required this.outcomes,
    required this.breakdown,
  });

  final int earned;
  final int total;
  final List<_Outcome> outcomes;
  final Map<String, _Tally> breakdown;
}
