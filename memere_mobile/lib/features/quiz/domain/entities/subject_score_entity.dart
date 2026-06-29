class SubjectScoreEntity {
  const SubjectScoreEntity({
    required this.earned,
    required this.possible,
  });

  final int earned;
  final int possible;

  double get percentage => possible <= 0 ? 0 : (earned / possible) * 100;
}
