class QuizAnswerFeedbackEntity {
  const QuizAnswerFeedbackEntity({
    required this.id,
    required this.text,
    required this.isCorrect,
  });

  final String id;
  final String text;
  final bool isCorrect;
}

class QuestionFeedbackEntity {
  const QuestionFeedbackEntity({
    required this.questionId,
    this.questionText,
    required this.correct,
    required this.pointsAwarded,
    required this.pointsPossible,
    required this.selectedAnswers,
    required this.correctAnswerIds,
    this.explanation,
    this.answers = const [],
  });

  final String questionId;
  final String? questionText;
  final bool correct;
  final int pointsAwarded;
  final int pointsPossible;
  final List<String> selectedAnswers;
  final List<String> correctAnswerIds;
  final String? explanation;
  final List<QuizAnswerFeedbackEntity> answers;
}
