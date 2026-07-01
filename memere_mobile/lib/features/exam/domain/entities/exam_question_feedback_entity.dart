class ExamQuestionFeedbackEntity {
  const ExamQuestionFeedbackEntity({
    required this.questionId,
    required this.correct,
    required this.marksAwarded,
    required this.marksPossible,
    required this.selectedAnswers,
    required this.correctAnswerIds,
    this.explanation,
  });

  final String questionId;
  final bool correct;
  final int marksAwarded;
  final int marksPossible;
  final List<String> selectedAnswers;
  final List<String> correctAnswerIds;
  final String? explanation;
}
