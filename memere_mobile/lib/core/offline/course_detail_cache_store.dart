import 'dart:convert';

import '../../features/courses/data/models/course_detail_model.dart';
import '../../features/courses/domain/entities/course_entity.dart';
import '../storage/hive/hive_boxes.dart';

/// Caches full course structure ([CourseDetailModel] JSON) into the
/// `course_details` Hive box, keyed by course id, so a downloaded course can be
/// opened and studied offline. Written only when a whole course is downloaded —
/// a listed course therefore always has its lessons' assets on device too.
///
/// Pure persistence with no network or feature dependencies, mirroring
/// [DownloadStore]/[SavedCoursesStore]. The course-detail provider (offline
/// fallback) and the offline course lists share it.
class CourseDetailCacheStore {
  const CourseDetailCacheStore();

  /// Persists the course structure for offline open. Idempotent by course id.
  Future<void> cache(CourseDetailModel detail) => AppHiveBoxes.courseDetails
      .put(detail.course.id, jsonEncode(detail.toJson()));

  /// The cached structure for [courseId], or null if the course isn't cached.
  CourseDetailModel? get(String courseId) =>
      _decode(AppHiveBoxes.courseDetails.get(courseId));

  /// Every cached (downloaded) course, as its lightweight [CourseEntity],
  /// sorted by title. Malformed rows are skipped.
  List<CourseEntity> listCourses() {
    final courses = <CourseEntity>[];
    for (final raw in AppHiveBoxes.courseDetails.values) {
      final detail = _decode(raw);
      if (detail != null) courses.add(detail.course);
    }
    courses.sort((a, b) => a.title.toLowerCase().compareTo(b.title.toLowerCase()));
    return courses;
  }

  bool isCached(String courseId) =>
      AppHiveBoxes.courseDetails.containsKey(courseId);

  Future<void> remove(String courseId) =>
      AppHiveBoxes.courseDetails.delete(courseId);

  CourseDetailModel? _decode(String? raw) {
    if (raw == null || raw.isEmpty) return null;
    try {
      final decoded = jsonDecode(raw);
      if (decoded is Map) {
        return CourseDetailModel.fromJson(Map<String, dynamic>.from(decoded));
      }
    } catch (_) {}
    return null;
  }
}
