import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:uuid/uuid.dart';

import 'download_store.dart';
import 'offline_attempt_store.dart';
import 'saved_courses_store.dart';

/// Shared Riverpod wiring for the offline-grading + downloads features.
///
/// The stores are stateless persistence wrappers over the (already-open) Hive
/// boxes, so they're plain singletons.

final downloadStoreProvider =
    Provider<DownloadStore>((ref) => const DownloadStore());

final savedCoursesStoreProvider =
    Provider<SavedCoursesStore>((ref) => const SavedCoursesStore());

final offlineAttemptStoreProvider =
    Provider<OfflineAttemptStore>((ref) => const OfflineAttemptStore());

const _uuid = Uuid();

/// A client-generated attempt id for on-device grading. The `local-` prefix is
/// the sole discriminator the attempt notifiers, result providers, and sync
/// service use to route an attempt to local grading instead of the server.
String newLocalAttemptId() => 'local-${_uuid.v4()}';

/// True when [id] denotes an on-device attempt — never submitted to the server
/// as-is; graded locally and resolved from [OfflineAttemptStore].
bool isLocalAttemptId(String id) => id.startsWith('local-');
