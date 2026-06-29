class QuizAnswerEntity {
  const QuizAnswerEntity({
    required this.id,
    required this.text,
    required this.orderIndex,
  });

  final String id;
  final String text;
  final int orderIndex;
}
