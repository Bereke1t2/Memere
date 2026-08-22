import 'dart:convert';

import '../storage/hive/hive_boxes.dart';
import '../storage/hive/models/saved_item.dart';

/// Persists user-favorited courses into the `saved_items` Hive box as
/// lightweight [SavedItem] rows. Purely local — guest favorites are never
/// synced to the server (the confirmed "keep local" decision).
///
/// Pure persistence with no network or feature dependencies, mirroring
/// [DownloadStore]; the favorites provider and the My Courses screen share it.
class SavedCoursesStore {
  const SavedCoursesStore();

  /// Favorited courses, newest first. Malformed rows are skipped.
  List<SavedItem> listCourses() {
    final items = <SavedItem>[];
    for (final raw in AppHiveBoxes.savedItems.values) {
      final item = _decode(raw);
      if (item != null && item.type == SavedType.course) items.add(item);
    }
    items.sort((a, b) => b.savedAt.compareTo(a.savedAt));
    return items;
  }

  bool isSaved(String courseId) =>
      AppHiveBoxes.savedItems.containsKey('${SavedType.course.name}:$courseId');

  Future<void> save(SavedItem item) =>
      AppHiveBoxes.savedItems.put(item.storageKey, jsonEncode(item.toJson()));

  Future<void> remove(String courseId) =>
      AppHiveBoxes.savedItems.delete('${SavedType.course.name}:$courseId');

  SavedItem? _decode(String? raw) {
    if (raw == null || raw.isEmpty) return null;
    try {
      final decoded = jsonDecode(raw);
      if (decoded is Map) {
        return SavedItem.fromJson(Map<String, dynamic>.from(decoded));
      }
    } catch (_) {}
    return null;
  }
}
