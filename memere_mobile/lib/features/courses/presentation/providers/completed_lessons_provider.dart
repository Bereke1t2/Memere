import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/storage/hive/hive_boxes.dart';

/// Local store provider for completed lesson IDs. Backed by Hive [AppHiveBoxes.savedItems].
final completedLessonsStoreProvider = Provider<CompletedLessonsStore>((ref) {
  return const CompletedLessonsStore();
});

class CompletedLessonsStore {
  const CompletedLessonsStore();

  static String _key(String lessonId) => 'completed_lesson:$lessonId';

  bool isCompleted(String lessonId) {
    if (lessonId.trim().isEmpty) return false;
    return AppHiveBoxes.savedItems.containsKey(_key(lessonId.trim()));
  }

  Future<void> markCompleted(String lessonId) async {
    if (lessonId.trim().isEmpty) return;
    await AppHiveBoxes.savedItems.put(
      _key(lessonId.trim()),
      DateTime.now().toIso8601String(),
    );
  }

  Future<bool> toggleCompleted(String lessonId) async {
    if (lessonId.trim().isEmpty) return false;
    final nowCompleted = !isCompleted(lessonId);
    if (nowCompleted) {
      await markCompleted(lessonId);
    } else {
      await AppHiveBoxes.savedItems.delete(_key(lessonId.trim()));
    }
    return nowCompleted;
  }

  Set<String> getCompletedIds() {
    final set = <String>{};
    for (final key in AppHiveBoxes.savedItems.keys) {
      if (key is String && key.startsWith('completed_lesson:')) {
        set.add(key.replaceFirst('completed_lesson:', ''));
      }
    }
    return set;
  }
}

/// Reactive Riverpod notifier for lesson completion state.
final completedLessonsProvider =
    AsyncNotifierProvider<CompletedLessonsNotifier, Set<String>>(
  CompletedLessonsNotifier.new,
);

class CompletedLessonsNotifier extends AsyncNotifier<Set<String>> {
  @override
  Future<Set<String>> build() async {
    final store = ref.read(completedLessonsStoreProvider);
    return store.getCompletedIds();
  }

  Future<bool> toggle(String lessonId) async {
    final store = ref.read(completedLessonsStoreProvider);
    final nowCompleted = await store.toggleCompleted(lessonId);
    state = AsyncData(store.getCompletedIds());
    return nowCompleted;
  }

  Future<void> markCompleted(String lessonId) async {
    final store = ref.read(completedLessonsStoreProvider);
    await store.markCompleted(lessonId);
    state = AsyncData(store.getCompletedIds());
  }
}
