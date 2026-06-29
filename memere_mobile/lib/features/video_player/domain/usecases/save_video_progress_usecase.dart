import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/lesson_progress_entity.dart';
import '../repositories/video_repository.dart';

class SaveVideoProgressParams {
  const SaveVideoProgressParams({
    required this.lessonId,
    required this.positionSeconds,
  });

  final String lessonId;
  final int positionSeconds;
}

class SaveVideoProgressUseCase {
  const SaveVideoProgressUseCase(this._repository);
  final VideoRepository _repository;

  Future<Either<Failure, LessonProgressEntity>> call(
    SaveVideoProgressParams params,
  ) {
    final lessonId = params.lessonId.trim();
    if (lessonId.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Lesson ID is required', field: 'lessonId'),
        ),
      );
    }
    if (params.positionSeconds < 0) {
      return Future.value(
        const Left(
          ValidationFailure(
            'Progress position cannot be negative',
            field: 'positionSeconds',
          ),
        ),
      );
    }
    return _repository.saveVideoProgress(
      lessonId: lessonId,
      positionSeconds: params.positionSeconds,
    );
  }
}
