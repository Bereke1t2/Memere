class ExamSubjectScoreEntity {
  const ExamSubjectScoreEntity({
    required this.key,
    required this.earned,
    required this.possible,
  });

  final String key;
  final int earned;
  final int possible;

  double get percentage => possible <= 0 ? 0 : (earned / possible) * 100;
}
