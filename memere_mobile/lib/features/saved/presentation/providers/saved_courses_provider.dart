import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/offline/offline_providers.dart';
import '../../../../core/storage/hive/models/saved_item.dart';
import '../../../courses/domain/entities/course_entity.dart';

/// Favorited courses (local, guest-friendly), newest first. Backed by
/// [SavedCoursesStore] over the `saved_items` Hive box. Favoriting is a LOCAL
/// feature — intentionally not account-gated, matching the quiz/exam download
/// precedent.
final savedCoursesProvider =
    AsyncNotifierProvider<SavedCoursesNotifier, List<SavedItem>>(
  SavedCoursesNotifier.new,
);

class SavedCoursesNotifier extends AsyncNotifier<List<SavedItem>> {
  @override
  Future<List<SavedItem>> build() async =>
      ref.read(savedCoursesStoreProvider).listCourses();

  /// Adds the course to favorites if absent, removes it if present. Returns the
  /// new saved state (true = now favorited) so the caller can tailor feedback.
  Future<bool> toggle(CourseEntity course) async {
    final store = ref.read(savedCoursesStoreProvider);
    final nowSaved = !store.isSaved(course.id);
    if (nowSaved) {
      await store.save(
        SavedItem(
          id: course.id,
          type: SavedType.course,
          courseId: course.id,
          title: course.title,
          subtitle: 'Grade ${course.grade} • ${course.subject}',
          savedAt: DateTime.now(),
        ),
      );
    } else {
      await store.remove(course.id);
    }
    state = AsyncData(store.listCourses());
    return nowSaved;
  }

  Future<void> remove(String courseId) async {
    final store = ref.read(savedCoursesStoreProvider);
    await store.remove(courseId);
    state = AsyncData(store.listCourses());
  }
}
