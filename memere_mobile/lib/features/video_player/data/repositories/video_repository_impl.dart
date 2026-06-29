import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/lesson_progress_entity.dart';
import '../../domain/entities/video_download_entity.dart';
import '../../domain/entities/video_status_entity.dart';
import '../../domain/entities/video_stream_entity.dart';
import '../../domain/repositories/video_repository.dart';
import '../datasources/video_remote_datasource.dart';

class VideoRepositoryImpl implements VideoRepository {
  const VideoRepositoryImpl(this._remoteDataSource);
  final VideoRemoteDataSource _remoteDataSource;

  @override
  Future<Either<Failure, VideoStatusEntity>> getVideoStatus(
    String videoId,
  ) async {
    try {
      final status = await _remoteDataSource.getVideoStatus(videoId);
      return Right(status);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, VideoStreamEntity>> getStream(String videoId) async {
    try {
      final stream = await _remoteDataSource.getStream(videoId);
      return Right(stream);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, VideoDownloadEntity>> getDownloadUrl(
    String videoId,
  ) async {
    try {
      final download = await _remoteDataSource.getDownloadUrl(videoId);
      return Right(download);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, LessonProgressEntity>> saveVideoProgress({
    required String lessonId,
    required int positionSeconds,
  }) async {
    try {
      final progress = await _remoteDataSource.saveVideoProgress(
        lessonId: lessonId,
        positionSeconds: positionSeconds,
      );
      return Right(progress);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> markLessonComplete(String lessonId) async {
    try {
      await _remoteDataSource.markLessonComplete(lessonId);
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }
}
