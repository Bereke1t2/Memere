class QuestionFeedbackEntity {
  const QuestionFeedbackEntity({
    required this.questionId,
    required this.correct,
    required this.pointsAwarded,
    required this.pointsPossible,
    required this.selectedAnswers,
    required this.correctAnswerIds,
    this.explanation,
  });

  final String questionId;
  final bool correct;
  final int pointsAwarded;
  final int pointsPossible;
  final List<String> selectedAnswers;
  final List<String> correctAnswerIds;
  final String? explanation;
}
