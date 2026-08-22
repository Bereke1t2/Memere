import 'dart:io';

import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';
import 'package:path_provider/path_provider.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/utils/media_url_helper.dart';
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
      // Filter out files that no longer exist on disk. Apply the same non-trivial
      // size check getDownload() uses, so a stale/aborted or manifest-sized file
      // (e.g. a saved .m3u8) is never surfaced as an available offline video.
      final validDownloads = <OfflineVideoModel>[];
      for (final download in downloads) {
        if (download.localPath.isEmpty) continue;
        final file = File(download.localPath);
        if (file.existsSync() && file.lengthSync() > 1000) {
          validDownloads.add(download);
        }
      }
      return Right(validDownloads);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, OfflineVideoEntity?>> getDownload(
      String videoId) async {
    try {
      final download = await _localDataSource.getDownload(videoId);
      if (download != null && download.localPath.isNotEmpty) {
        final file = File(download.localPath);
        if (await file.exists() && (await file.length()) > 1000) {
          return Right(download);
        }
      }
      return const Right(null);
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
    void Function(int received, int total)? onReceiveProgress,
  }) async {
    if (videoId.trim().isEmpty) {
      return const Left(ValidationFailure('Video ID is required'));
    }

    try {
      // 1. Resolve the offline download URL. The backend serves a real,
      //    self-contained MP4 here (not the HLS manifest), so it is safe to save
      //    and play back offline. We deliberately do NOT fall back to the stream
      //    master URL (an .m3u8 manifest is not a downloadable single file) nor to
      //    any sample clip — a video that isn't ready must surface as an error.
      String? downloadUrl;
      final downloadUrlResult = await _videoRepository.getDownloadUrl(videoId);
      downloadUrlResult.fold(
        (_) {},
        (download) {
          if (download.downloadUrl.isNotEmpty) {
            downloadUrl = fixMediaUrl(download.downloadUrl);
          }
        },
      );

      if (downloadUrl == null || downloadUrl!.isEmpty) {
        return const Left(
          ServerFailure('This video is not ready for offline download yet.'),
        );
      }

      // 2. Ensure destination storage directory exists
      final appDir = await getApplicationDocumentsDirectory();
      final videosDir = Directory('${appDir.path}/offline_videos');
      if (!await videosDir.exists()) {
        await videosDir.create(recursive: true);
      }

      final safeVideoId = videoId.replaceAll(RegExp(r'[^a-zA-Z0-9_-]'), '_');
      final localFilePath = '${videosDir.path}/video_$safeVideoId.mp4';
      final file = File(localFilePath);

      // 3. Download file using Dio with timeout and redirect handling
      final dio = Dio(
        BaseOptions(
          connectTimeout: const Duration(seconds: 30),
          receiveTimeout: const Duration(minutes: 10),
          followRedirects: true,
        ),
      );

      await dio.download(
        downloadUrl!,
        localFilePath,
        // Reports (bytesReceived, totalBytes) as the file streams in. total is
        // -1 when the response carries no Content-Length; the caller guards on
        // total > 0 before computing a percentage.
        onReceiveProgress: onReceiveProgress,
        options: Options(
          validateStatus: (status) => status != null && status < 400,
        ),
      );

      // 4. Validate the bytes are actually a video before recording the download.
      //    Guards against silently saving an HLS manifest (.m3u8) or an HTML error
      //    page under a .mp4 name — either would fail to play offline.
      if (!await _looksLikeVideo(file)) {
        if (await file.exists()) {
          await file.delete();
        }
        return const Left(
          CacheFailure('The downloaded file is not a playable video.'),
        );
      }

      final sizeBytes = await file.exists() ? await file.length() : 0;
      final now = DateTime.now();

      final model = OfflineVideoModel(
        videoId: videoId,
        lessonId: lessonId,
        courseId: courseId,
        title: title,
        localPath: localFilePath,
        downloadedAt: now,
        expiresAt: now.add(AppConstants.downloadedVideoTtl),
        fileSizeBytes: sizeBytes,
        status: OfflineVideoStatus.downloaded,
      );

      await _localDataSource.saveDownload(model);
      return Right(model);
    } catch (e) {
      return Left(CacheFailure('Failed to download video for offline viewing: $e'));
    }
  }

  /// Best-effort sniff that a downloaded file is a real media file, not an HLS
  /// manifest or an HTML error page saved under a .mp4 name. Rejects trivially
  /// small files and the two text signatures we actually see on failure.
  Future<bool> _looksLikeVideo(File file) async {
    if (!await file.exists()) return false;
    if (await file.length() < 1024) return false;
    final head = await file
        .openRead(0, 16)
        .fold<List<int>>(<int>[], (acc, chunk) => acc..addAll(chunk));
    final text = String.fromCharCodes(head).trimLeft();
    if (text.startsWith('#EXTM3U')) return false; // HLS manifest
    if (text.startsWith('<')) return false; // HTML error page
    return true;
  }

  @override
  Future<Either<Failure, void>> removeDownload(String videoId) async {
    try {
      final download = await _localDataSource.getDownload(videoId);
      if (download != null && download.localPath.isNotEmpty) {
        final file = File(download.localPath);
        if (await file.exists()) {
          await file.delete();
        }
      }
      await _localDataSource.removeDownload(videoId);
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> clearExpiredDownloads() async {
    try {
      final downloads = await _localDataSource.getDownloads();
      for (final download in downloads) {
        if (download.isExpired && download.localPath.isNotEmpty) {
          final file = File(download.localPath);
          if (await file.exists()) {
            await file.delete();
          }
        }
      }
      await _localDataSource.clearExpiredDownloads();
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }
}
