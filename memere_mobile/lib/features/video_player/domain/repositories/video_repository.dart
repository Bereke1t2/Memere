import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/lesson_progress_entity.dart';
import '../entities/video_download_entity.dart';
import '../entities/video_status_entity.dart';
import '../entities/video_stream_entity.dart';

abstract class VideoRepository {
  Future<Either<Failure, VideoStatusEntity>> getVideoStatus(String videoId);

  Future<Either<Failure, VideoStreamEntity>> getStream(String videoId);

  Future<Either<Failure, VideoDownloadEntity>> getDownloadUrl(String videoId);

  Future<Either<Failure, LessonProgressEntity>> saveVideoProgress({
    required String lessonId,
    required int positionSeconds,
  });

  Future<Either<Failure, void>> markLessonComplete(String lessonId);
}
