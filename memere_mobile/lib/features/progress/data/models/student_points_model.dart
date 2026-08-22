import '../../domain/entities/student_points_entity.dart';

class StudentPointsModel extends StudentPointsEntity {
  const StudentPointsModel({
    required super.totalPoints,
    required super.quizPoints,
    required super.examPoints,
    required super.avgPercentage,
    required super.quizCount,
    required super.examCount,
  });

  factory StudentPointsModel.fromJson(Map<String, dynamic> json) {
    return StudentPointsModel(
      totalPoints: _doubleValue(json['total_points']),
      quizPoints: _doubleValue(json['quiz_points']),
      examPoints: _doubleValue(json['exam_points']),
      avgPercentage: _doubleValue(json['avg_percentage']),
      quizCount: _intValue(json['quiz_count']),
      examCount: _intValue(json['exam_count']),
    );
  }
}

// The backend sends numeric fields as JSON numbers, but tolerate strings too so
// a serialization change never crashes the profile.
double _doubleValue(Object? value, {double fallback = 0}) {
  if (value is num) return value.toDouble();
  if (value is String) return double.tryParse(value) ?? fallback;
  return fallback;
}

int _intValue(Object? value, {int fallback = 0}) {
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value) ?? fallback;
  return fallback;
}
