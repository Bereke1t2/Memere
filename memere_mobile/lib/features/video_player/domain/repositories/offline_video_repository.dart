import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/offline_video_entity.dart';

abstract class OfflineVideoRepository {
  Future<Either<Failure, List<OfflineVideoEntity>>> getDownloads();

  Future<Either<Failure, OfflineVideoEntity?>> getDownload(String videoId);

  Future<Either<Failure, OfflineVideoEntity>> startDownload({
    required String videoId,
    required String lessonId,
    required String courseId,
    required String title,
  });

  Future<Either<Failure, void>> removeDownload(String videoId);

  Future<Either<Failure, void>> clearExpiredDownloads();
}
