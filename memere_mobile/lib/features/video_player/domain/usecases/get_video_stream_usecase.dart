import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/video_stream_entity.dart';
import '../repositories/video_repository.dart';

class GetVideoStreamUseCase {
  const GetVideoStreamUseCase(this._repository);
  final VideoRepository _repository;

  Future<Either<Failure, VideoStreamEntity>> call(String videoId) {
    final trimmedId = videoId.trim();
    if (trimmedId.isEmpty) {
      return Future.value(
        const Left(ValidationFailure('Video ID is required', field: 'videoId')),
      );
    }
    return _repository.getStream(trimmedId);
  }
}
