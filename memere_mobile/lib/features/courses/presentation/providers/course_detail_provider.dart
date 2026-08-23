import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/offline/offline_providers.dart';
import '../../domain/entities/course_detail_entity.dart';
import 'courses_providers.dart';

final courseDetailProvider =
    FutureProvider.family<CourseDetailEntity, String>((ref, courseId) async {
  final useCase = ref.watch(getCourseDetailUseCaseProvider);
  final result = await useCase(courseId);
  return result.fold(
    (failure) {
      // Offline (or server error): fall back to the cached structure of a
      // downloaded course so it can still be opened & studied. Only a course
      // that was downloaded is cached, so its lessons' assets are on device too.
      final cached = ref.read(courseDetailCacheStoreProvider).get(courseId);
      if (cached != null) return cached;
      throw failure;
    },
    (detail) => detail,
  );
});
