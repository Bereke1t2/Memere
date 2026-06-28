import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/course_detail_entity.dart';
import 'courses_providers.dart';

final courseDetailProvider =
    FutureProvider.family<CourseDetailEntity, String>((ref, courseId) async {
  final useCase = ref.watch(getCourseDetailUseCaseProvider);
  final result = await useCase(courseId);
  return result.fold(
    (failure) => throw failure,
    (detail) => detail,
  );
});
