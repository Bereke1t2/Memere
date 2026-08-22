import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:memere_mobile/core/constants/app_constants.dart';
import 'package:memere_mobile/core/offline/offline_attempt_store.dart';
import 'package:memere_mobile/core/storage/hive/models/pending_submission.dart';

/// Exercises [OfflineAttemptStore] against real (plaintext) Hive boxes in a
/// temp dir — the same box names the app opens at boot, so the store's
/// `AppHiveBoxes.*` getters resolve to these. Covers the two persisted shapes:
/// provisional results and the replay queue.
void main() {
  late Directory tempDir;
  const store = OfflineAttemptStore();

  setUpAll(() async {
    tempDir = await Directory.systemTemp.createTemp('offline_attempt_store');
    Hive.init(tempDir.path);
    await Hive.openBox<String>(AppConstants.offlineAttemptResultsBoxKey);
    await Hive.openBox<String>(AppConstants.syncQueueKey);
  });

  tearDown(() async {
    await Hive.box<String>(AppConstants.offlineAttemptResultsBoxKey).clear();
    await Hive.box<String>(AppConstants.syncQueueKey).clear();
  });

  tearDownAll(() async {
    await Hive.close();
    await tempDir.delete(recursive: true);
  });

  group('provisional results', () {
    test('saveResult → getResult round-trips a JSON map', () async {
      final json = {
        'attempt_id': 'local-1',
        'score': 8,
        'passed': true,
        'nested': {'a': 1},
      };
      await store.saveResult('local-1', json);

      expect(store.hasResult('local-1'), isTrue);
      expect(store.getResult('local-1'), equals(json));
    });

    test('getResult is null for a missing id; removeResult deletes', () async {
      expect(store.getResult('nope'), isNull);
      expect(store.hasResult('nope'), isFalse);

      await store.saveResult('local-x', {'k': 'v'});
      await store.removeResult('local-x');
      expect(store.hasResult('local-x'), isFalse);
      expect(store.getResult('local-x'), isNull);
    });
  });

  group('replay queue', () {
    PendingSubmission sub() => PendingSubmission(
          localId: 'local-9',
          kind: SubmissionKind.exam,
          contentId: 'exam-1',
          rawAnswers: {
            'q1': ['a1'],
            'q2': 'free text',
          },
          provisionalScore: 12,
          provisionalPct: 60.0,
          takenAt: DateTime.parse('2026-08-21T10:00:00.000Z'),
        );

    test('enqueue → listQueue → dequeue round-trips a PendingSubmission',
        () async {
      await store.enqueue(sub());

      final queue = store.listQueue();
      expect(queue, hasLength(1));
      final got = queue.single;
      expect(got.localId, 'local-9');
      expect(got.kind, SubmissionKind.exam);
      expect(got.contentId, 'exam-1');
      expect(got.rawAnswers['q1'], ['a1']);
      expect(got.rawAnswers['q2'], 'free text');
      expect(got.provisionalScore, 12);
      expect(got.provisionalPct, 60.0);
      expect(got.takenAt.toUtc(), DateTime.parse('2026-08-21T10:00:00.000Z'));
      expect(got.serverAttemptId, isNull);
      expect(got.attempts, 0);

      await store.dequeue('local-9');
      expect(store.listQueue(), isEmpty);
    });

    test('enqueue overwrites by localId (the copyWith update path)', () async {
      await store.enqueue(sub());
      await store
          .enqueue(sub().copyWith(attempts: 2, serverAttemptId: 'srv-1'));

      final queue = store.listQueue();
      expect(queue, hasLength(1)); // same key → replaced, not duplicated
      expect(queue.single.attempts, 2);
      expect(queue.single.serverAttemptId, 'srv-1');
    });
  });
}
