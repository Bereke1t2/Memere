import 'dart:convert';

import '../storage/hive/hive_boxes.dart';
import '../storage/hive/models/pending_submission.dart';

/// Persists provisional (on-device-graded) results and the replay queue.
///
/// Results are stored as opaque JSON maps keyed by the attempt's **local id**,
/// so this store stays feature-agnostic — the quiz/exam layers (de)serialize
/// their own `QuizResultModel`/`ExamResultModel` via the maps. The queue holds
/// [PendingSubmission]s awaiting server replay (authenticated users only;
/// guests never enqueue).
class OfflineAttemptStore {
  const OfflineAttemptStore();

  // ── Provisional results (offline_attempt_results box) ───────────────────────

  Future<void> saveResult(String attemptId, Map<String, dynamic> resultJson) =>
      AppHiveBoxes.offlineAttemptResults.put(attemptId, jsonEncode(resultJson));

  Map<String, dynamic>? getResult(String attemptId) =>
      _decodeMap(AppHiveBoxes.offlineAttemptResults.get(attemptId));

  bool hasResult(String attemptId) =>
      AppHiveBoxes.offlineAttemptResults.containsKey(attemptId);

  Future<void> removeResult(String attemptId) =>
      AppHiveBoxes.offlineAttemptResults.delete(attemptId);

  // ── Replay queue (sync_queue box) ────────────────────────────────────────────

  /// Inserts or overwrites a queued submission (keyed by its `localId`). Also the
  /// update path after [PendingSubmission.copyWith] bumps `attempts` or records a
  /// `serverAttemptId`.
  Future<void> enqueue(PendingSubmission submission) => AppHiveBoxes.syncQueue
      .put(submission.localId, jsonEncode(submission.toJson()));

  List<PendingSubmission> listQueue() {
    final out = <PendingSubmission>[];
    for (final raw in AppHiveBoxes.syncQueue.values) {
      final map = _decodeMap(raw);
      if (map != null) out.add(PendingSubmission.fromJson(map));
    }
    return out;
  }

  Future<void> dequeue(String localId) =>
      AppHiveBoxes.syncQueue.delete(localId);

  // ── internals ──────────────────────────────────────────────────────────────

  Map<String, dynamic>? _decodeMap(String? raw) {
    if (raw == null || raw.isEmpty) return null;
    try {
      final decoded = jsonDecode(raw);
      if (decoded is Map) return Map<String, dynamic>.from(decoded);
    } catch (_) {}
    return null;
  }
}
