import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/errors/failures.dart';
import '../../data/datasources/offline_video_local_datasource.dart';
import '../../data/repositories/offline_video_repository_impl.dart';
import '../../domain/entities/offline_video_entity.dart';
import '../../domain/repositories/offline_video_repository.dart';
import 'video_providers.dart';

final offlineVideoLocalDataSourceProvider =
    Provider<OfflineVideoLocalDataSource>((ref) {
  return const OfflineVideoLocalDataSourceImpl();
});

final offlineVideoRepositoryProvider = Provider<OfflineVideoRepository>((ref) {
  return OfflineVideoRepositoryImpl(
    localDataSource: ref.watch(offlineVideoLocalDataSourceProvider),
    videoRepository: ref.watch(videoRepositoryProvider),
  );
});

final offlineDownloadsProvider =
    AsyncNotifierProvider<OfflineDownloadsNotifier, List<OfflineVideoEntity>>(
  OfflineDownloadsNotifier.new,
);

class OfflineDownloadsNotifier extends AsyncNotifier<List<OfflineVideoEntity>> {
  @override
  Future<List<OfflineVideoEntity>> build() async {
    final result =
        await ref.read(offlineVideoRepositoryProvider).getDownloads();
    return result.fold((failure) => throw failure, (downloads) => downloads);
  }

  Future<void> startDownload({
    required String videoId,
    required String lessonId,
    required String courseId,
    required String title,
  }) async {
    final current = state.valueOrNull ?? const [];
    state = const AsyncLoading<List<OfflineVideoEntity>>().copyWithPrevious(
      state,
    );

    final result = await ref.read(offlineVideoRepositoryProvider).startDownload(
          videoId: videoId,
          lessonId: lessonId,
          courseId: courseId,
          title: title,
        );

    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (download) {
        final next = [
          for (final item in current)
            if (item.videoId != download.videoId) item,
          download,
        ];
        return AsyncData(next);
      },
    );
  }

  Future<void> removeDownload(String videoId) async {
    final result =
        await ref.read(offlineVideoRepositoryProvider).removeDownload(videoId);
    result.fold(
      (failure) => state = AsyncError(failure, StackTrace.current),
      (_) {
        final current = state.valueOrNull ?? const [];
        state = AsyncData(
          current.where((download) => download.videoId != videoId).toList(),
        );
      },
    );
  }

  Future<void> clearExpiredDownloads() async {
    final result =
        await ref.read(offlineVideoRepositoryProvider).clearExpiredDownloads();
    result.fold(
      (failure) => state = AsyncError(failure, StackTrace.current),
      (_) => ref.invalidateSelf(),
    );
  }
}

String offlineFailureMessage(Object error) {
  if (error is Failure) return error.message;
  return 'Could not update downloads.';
}
