import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import 'exam_providers.dart';

/// Exam ids whose download is currently in flight — drives the button spinner
/// without collapsing the whole downloads list into a loading state.
final examDownloadBusyProvider = StateProvider<Set<String>>((ref) => const {});

/// Downloaded exams (manifest rows), newest first. Backed by the shared
/// [DownloadStore]; the heavy payloads live in the encrypted box.
final downloadedExamsProvider =
    AsyncNotifierProvider<DownloadedExamsNotifier, List<DownloadedItem>>(
  DownloadedExamsNotifier.new,
);

class DownloadedExamsNotifier extends AsyncNotifier<List<DownloadedItem>> {
  @override
  Future<List<DownloadedItem>> build() async => _read();

  List<DownloadedItem> _read() => ref
      .read(downloadStoreProvider)
      .listDownloads()
      .where((item) => item.type == DownloadType.exam)
      .toList();

  /// Fetches the exam WITH answer keys and stores it for offline grading.
  /// Downloading is a LOCAL feature — intentionally NOT gated by the account
  /// gate, so guests can download published/free exams. Returns a [Failure]
  /// on error (null on success) so the button can surface it.
  Future<Failure?> download({
    required String examId,
    required String title,
  }) async {
    final busy = ref.read(examDownloadBusyProvider.notifier);
    busy.state = {...busy.state, examId};
    try {
      final exam =
          await ref.read(examRemoteDataSourceProvider).getForDownload(examId);
      await ref.read(downloadStoreProvider).saveExam(exam, title: title);
      state = AsyncData(_read());
      return null;
    } on DioException catch (e) {
      return ServerFailure.fromDioError(e);
    } catch (_) {
      return const UnknownFailure('Could not download this exam.');
    } finally {
      busy.state = {...busy.state}..remove(examId);
    }
  }

  Future<void> remove(String examId) async {
    await ref.read(downloadStoreProvider).removeExam(examId);
    final current = state.valueOrNull ?? const [];
    state = AsyncData(current.where((item) => item.id != examId).toList());
  }
}
