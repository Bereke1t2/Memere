import 'dart:convert';

import '../storage/hive/hive_boxes.dart';
import '../storage/hive/models/downloaded_item.dart';
import '../storage/hive/models/offline_exam.dart';
import '../storage/hive/models/offline_quiz.dart';

/// Persists downloaded quizzes/exams (WITH answer keys) into the encrypted Hive
/// boxes and records a lightweight [DownloadedItem] row in the unified
/// `downloads_index` manifest.
///
/// Pure persistence — no network, no feature dependencies — so the quiz/exam
/// download controllers and the Saved screen all share it. The heavy payloads
/// live in the encrypted boxes; the manifest row is the queryable index.
class DownloadStore {
  const DownloadStore();

  // ── Quizzes ────────────────────────────────────────────────────────────────

  Future<void> saveQuiz(OfflineQuiz quiz, {String title = ''}) async {
    final payload = jsonEncode(quiz.toJson());
    await AppHiveBoxes.downloadedQuizzes.put(quiz.id, payload);
    await _putManifest(
      DownloadedItem(
        id: quiz.id,
        type: DownloadType.quiz,
        courseId: quiz.courseId,
        title: title.isNotEmpty ? title : quiz.title,
        subtitle: '${quiz.questions.length} questions',
        sizeBytes: payload.length,
        status: DownloadItemStatus.downloaded,
        downloadedAt: quiz.downloadedAt,
      ),
    );
  }

  OfflineQuiz? getOfflineQuiz(String quizId) =>
      _decode(AppHiveBoxes.downloadedQuizzes.get(quizId), OfflineQuiz.fromJson);

  Future<void> removeQuiz(String quizId) async {
    await AppHiveBoxes.downloadedQuizzes.delete(quizId);
    await AppHiveBoxes.downloadsIndex.delete('${DownloadType.quiz.name}:$quizId');
  }

  // ── Exams ────────────────────────────────────────────────────────────────

  Future<void> saveExam(OfflineExam exam, {String title = ''}) async {
    final payload = jsonEncode(exam.toJson());
    await AppHiveBoxes.downloadedExams.put(exam.id, payload);
    await _putManifest(
      DownloadedItem(
        id: exam.id,
        type: DownloadType.exam,
        courseId: exam.courseId ?? '',
        title: title.isNotEmpty ? title : exam.title,
        subtitle: exam.subject,
        sizeBytes: payload.length,
        status: DownloadItemStatus.downloaded,
        downloadedAt: exam.downloadedAt,
      ),
    );
  }

  OfflineExam? getOfflineExam(String examId) =>
      _decode(AppHiveBoxes.downloadedExams.get(examId), OfflineExam.fromJson);

  Future<void> removeExam(String examId) async {
    await AppHiveBoxes.downloadedExams.delete(examId);
    await AppHiveBoxes.downloadsIndex.delete('${DownloadType.exam.name}:$examId');
  }

  // ── Manifest / queries ─────────────────────────────────────────────────────

  /// True when the encrypted payload (answer keys) is present on device — the
  /// precondition for on-device grading.
  bool isDownloaded(DownloadType type, String id) {
    switch (type) {
      case DownloadType.quiz:
        return AppHiveBoxes.downloadedQuizzes.containsKey(id);
      case DownloadType.exam:
        return AppHiveBoxes.downloadedExams.containsKey(id);
      case DownloadType.video:
      case DownloadType.pdf:
        return false; // owned by other stores
    }
  }

  /// Every downloaded asset across kinds, newest first. Malformed rows skipped.
  List<DownloadedItem> listDownloads() {
    final items = <DownloadedItem>[];
    for (final raw in AppHiveBoxes.downloadsIndex.values) {
      final item = _decode(raw, DownloadedItem.fromJson);
      if (item != null) items.add(item);
    }
    items.sort((a, b) => b.downloadedAt.compareTo(a.downloadedAt));
    return items;
  }

  // ── internals ──────────────────────────────────────────────────────────────

  Future<void> _putManifest(DownloadedItem item) =>
      AppHiveBoxes.downloadsIndex.put(item.storageKey, jsonEncode(item.toJson()));

  T? _decode<T>(String? raw, T Function(Map<String, dynamic>) fromJson) {
    if (raw == null || raw.isEmpty) return null;
    try {
      final decoded = jsonDecode(raw);
      if (decoded is Map) return fromJson(Map<String, dynamic>.from(decoded));
    } catch (_) {}
    return null;
  }
}
