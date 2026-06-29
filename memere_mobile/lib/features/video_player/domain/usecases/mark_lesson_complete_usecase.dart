import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/video_repository.dart';

class MarkLessonCompleteUseCase {
  const MarkLessonCompleteUseCase(this._repository);
  final VideoRepository _repository;

  Future<Either<Failure, void>> call(String lessonId) {
    final trimmedId = lessonId.trim();
    if (trimmedId.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Lesson ID is required', field: 'lessonId'),
        ),
      );
    }
    return _repository.markLessonComplete(trimmedId);
  }
}
