import '../../domain/entities/exam_subject_score_entity.dart';
import 'exam_model_helpers.dart';

class ExamSubjectScoreModel extends ExamSubjectScoreEntity {
  const ExamSubjectScoreModel({
    required super.key,
    required super.earned,
    required super.possible,
  });

  /// Parses an entry that carries its own `key` (analytics list shape).
  factory ExamSubjectScoreModel.fromJson(Map<String, dynamic> json) {
    return ExamSubjectScoreModel(
      key: examStringValue(json['key']),
      earned: examIntValue(json['earned']),
      possible: examIntValue(json['possible']),
    );
  }

  /// Parses a value where the `key` comes from the enclosing map (result shape).
  factory ExamSubjectScoreModel.fromMapEntry(
    String key,
    Map<String, dynamic> json,
  ) {
    return ExamSubjectScoreModel(
      key: key,
      earned: examIntValue(json['earned']),
      possible: examIntValue(json['possible']),
    );
  }

  Map<String, dynamic> toJson() => {
        'key': key,
        'earned': earned,
        'possible': possible,
      };
}
