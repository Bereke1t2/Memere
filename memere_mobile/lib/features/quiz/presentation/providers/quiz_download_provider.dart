import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import 'quiz_providers.dart';

/// Quiz ids whose download is currently in flight — drives the button spinner
/// without collapsing the whole downloads list into a loading state.
final quizDownloadBusyProvider = StateProvider<Set<String>>((ref) => const {});

/// Downloaded quizzes (manifest rows), newest first. Backed by the shared
/// [DownloadStore]; the heavy payloads live in the encrypted box.
final downloadedQuizzesProvider =
    AsyncNotifierProvider<DownloadedQuizzesNotifier, List<DownloadedItem>>(
  DownloadedQuizzesNotifier.new,
);

class DownloadedQuizzesNotifier extends AsyncNotifier<List<DownloadedItem>> {
  @override
  Future<List<DownloadedItem>> build() async => _read();

  List<DownloadedItem> _read() => ref
      .read(downloadStoreProvider)
      .listDownloads()
      .where((item) => item.type == DownloadType.quiz)
      .toList();

  /// Fetches the quiz WITH answer keys and stores it for offline grading.
  /// Downloading is a LOCAL feature — intentionally NOT gated by the account
  /// gate, so guests can download published/free quizzes. Returns a [Failure]
  /// on error (null on success) so the button can surface it.
  Future<Failure?> download({
    required String quizId,
    required String title,
  }) async {
    final busy = ref.read(quizDownloadBusyProvider.notifier);
    busy.state = {...busy.state, quizId};
    try {
      final quiz =
          await ref.read(quizRemoteDataSourceProvider).getForDownload(quizId);
      await ref.read(downloadStoreProvider).saveQuiz(quiz, title: title);
      state = AsyncData(_read());
      return null;
    } on DioException catch (e) {
      return ServerFailure.fromDioError(e);
    } catch (_) {
      return const UnknownFailure('Could not download this quiz.');
    } finally {
      busy.state = {...busy.state}..remove(quizId);
    }
  }

  Future<void> remove(String quizId) async {
    await ref.read(downloadStoreProvider).removeQuiz(quizId);
    final current = state.valueOrNull ?? const [];
    state = AsyncData(current.where((item) => item.id != quizId).toList());
  }
}
