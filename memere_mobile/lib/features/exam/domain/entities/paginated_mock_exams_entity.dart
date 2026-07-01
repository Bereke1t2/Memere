import 'mock_exam_entity.dart';

class PaginatedMockExamsEntity {
  const PaginatedMockExamsEntity({
    required this.exams,
    required this.nextCursor,
    required this.limit,
  });

  final List<MockExamEntity> exams;
  final String? nextCursor;
  final int limit;

  bool get hasMore => nextCursor != null && nextCursor!.isNotEmpty;
}
