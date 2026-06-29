import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/datasources/video_remote_datasource.dart';
import '../../data/repositories/video_repository_impl.dart';
import '../../domain/repositories/video_repository.dart';
import '../../domain/usecases/get_video_download_url_usecase.dart';
import '../../domain/usecases/get_video_status_usecase.dart';
import '../../domain/usecases/get_video_stream_usecase.dart';
import '../../domain/usecases/mark_lesson_complete_usecase.dart';
import '../../domain/usecases/save_video_progress_usecase.dart';

final videoRemoteDataSourceProvider = Provider<VideoRemoteDataSource>((ref) {
  return VideoRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final videoRepositoryProvider = Provider<VideoRepository>((ref) {
  return VideoRepositoryImpl(ref.watch(videoRemoteDataSourceProvider));
});

final getVideoStatusUseCaseProvider = Provider<GetVideoStatusUseCase>((ref) {
  return GetVideoStatusUseCase(ref.watch(videoRepositoryProvider));
});

final getVideoStreamUseCaseProvider = Provider<GetVideoStreamUseCase>((ref) {
  return GetVideoStreamUseCase(ref.watch(videoRepositoryProvider));
});

final getVideoDownloadUrlUseCaseProvider =
    Provider<GetVideoDownloadUrlUseCase>((ref) {
  return GetVideoDownloadUrlUseCase(ref.watch(videoRepositoryProvider));
});

final saveVideoProgressUseCaseProvider =
    Provider<SaveVideoProgressUseCase>((ref) {
  return SaveVideoProgressUseCase(ref.watch(videoRepositoryProvider));
});

final markLessonCompleteUseCaseProvider =
    Provider<MarkLessonCompleteUseCase>((ref) {
  return MarkLessonCompleteUseCase(ref.watch(videoRepositoryProvider));
});
