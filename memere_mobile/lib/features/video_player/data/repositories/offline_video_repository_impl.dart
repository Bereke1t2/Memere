import 'package:fpdart/fpdart.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../domain/entities/offline_video_entity.dart';
import '../../domain/repositories/offline_video_repository.dart';
import '../../domain/repositories/video_repository.dart';
import '../datasources/offline_video_local_datasource.dart';
import '../models/offline_video_model.dart';

class OfflineVideoRepositoryImpl implements OfflineVideoRepository {
  const OfflineVideoRepositoryImpl({
    required OfflineVideoLocalDataSource localDataSource,
    required VideoRepository videoRepository,
  })  : _localDataSource = localDataSource,
        _videoRepository = videoRepository;

  final OfflineVideoLocalDataSource _localDataSource;
  final VideoRepository _videoRepository;

  @override
  Future<Either<Failure, List<OfflineVideoEntity>>> getDownloads() async {
    try {
      await _localDataSource.clearExpiredDownloads();
      final downloads = await _localDataSource.getDownloads();
      return Right(downloads);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, OfflineVideoEntity?>> getDownload(
      String videoId) async {
    try {
      final download = await _localDataSource.getDownload(videoId);
      return Right(download);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, OfflineVideoEntity>> startDownload({
    required String videoId,
    required String lessonId,
    required String courseId,
    required String title,
  }) async {
    if (videoId.trim().isEmpty) {
      return const Left(ValidationFailure('Video ID is required'));
    }

    final downloadUrlResult = await _videoRepository.getDownloadUrl(videoId);
    return downloadUrlResult.fold(
      Left.new,
      (download) async {
        try {
          // Phase 3 stores safe metadata only. Full HLS segment download and
          // encrypted local playback are intentionally deferred.
          final now = DateTime.now();
          final model = OfflineVideoModel(
            videoId: videoId,
            lessonId: lessonId,
            courseId: courseId,
            title: title,
            localPath: '',
            downloadedAt: now,
            expiresAt: now.add(AppConstants.downloadedVideoTtl),
            fileSizeBytes: 0,
            status: OfflineVideoStatus.downloaded,
          );
          await _localDataSource.saveDownload(model);
          return Right(model);
        } catch (e) {
          return Left(CacheFailure(e.toString()));
        }
      },
    );
  }

  @override
  Future<Either<Failure, void>> removeDownload(String videoId) async {
    try {
      await _localDataSource.removeDownload(videoId);
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> clearExpiredDownloads() async {
    try {
      await _localDataSource.clearExpiredDownloads();
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }
}
