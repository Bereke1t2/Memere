import '../../domain/entities/subject_score_entity.dart';
import 'quiz_model_helpers.dart';

class SubjectScoreModel extends SubjectScoreEntity {
  const SubjectScoreModel({
    required super.earned,
    required super.possible,
  });

  factory SubjectScoreModel.fromJson(Map<String, dynamic> json) {
    return SubjectScoreModel(
      earned: quizIntValue(json['earned']),
      possible: quizIntValue(json['possible']),
    );
  }

  Map<String, dynamic> toJson() => {
        'earned': earned,
        'possible': possible,
      };
}
