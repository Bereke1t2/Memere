/// A student's cumulative points across every graded quiz and exam — the best
/// attempt per item, summed on the server, so retakes raise the total only by
/// the improvement. Points are an account feature (guests have none).
class StudentPointsEntity {
  const StudentPointsEntity({
    required this.totalPoints,
    required this.quizPoints,
    required this.examPoints,
    required this.avgPercentage,
    required this.quizCount,
    required this.examCount,
  });

  /// Quiz + exam points combined — the headline figure on the profile.
  final double totalPoints;
  final double quizPoints;
  final double examPoints;

  /// Overall average percentage across the counted quizzes and exams (0–100).
  final double avgPercentage;

  /// Distinct quizzes / exams the student has completed at least once.
  final int quizCount;
  final int examCount;

  /// True when nothing has been completed yet — the UI shows an encouraging
  /// prompt instead of a bare "0".
  bool get isEmpty => quizCount == 0 && examCount == 0;
}
