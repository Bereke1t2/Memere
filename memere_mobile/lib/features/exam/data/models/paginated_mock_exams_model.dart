import '../../domain/entities/paginated_mock_exams_entity.dart';
import 'exam_model_helpers.dart';
import 'mock_exam_model.dart';

class PaginatedMockExamsModel extends PaginatedMockExamsEntity {
  const PaginatedMockExamsModel({
    required super.exams,
    required super.nextCursor,
    required super.limit,
  });

  factory PaginatedMockExamsModel.fromJson(Map<String, dynamic> json) {
    final data = json['data'];
    final exams = data is List
        ? data
            .whereType<Map<String, dynamic>>()
            .map(MockExamModel.fromJson)
            .toList()
        : <MockExamModel>[];

    final nextCursor = json['next_cursor'];
    return PaginatedMockExamsModel(
      exams: exams,
      nextCursor:
          nextCursor is String && nextCursor.isNotEmpty ? nextCursor : null,
      limit: examIntValue(json['limit'], fallback: 20),
    );
  }
}
