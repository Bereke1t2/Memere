class ExamFeedbackAnswerEntity {
  const ExamFeedbackAnswerEntity({
    required this.id,
    required this.text,
    required this.isCorrect,
    this.orderIndex = 0,
  });

  final String id;
  final String text;
  final bool isCorrect;
  final int orderIndex;
}
