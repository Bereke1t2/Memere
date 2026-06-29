import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/quiz_answer_payload.dart';
import '../entities/quiz_attempt_entity.dart';
import '../entities/quiz_entity.dart';
import '../entities/quiz_result_entity.dart';

abstract class QuizRepository {
  Future<Either<Failure, List<QuizEntity>>> listQuizzesByCourse(
    String courseId,
  );

  Future<Either<Failure, QuizEntity>> getQuiz(String quizId);

  Future<Either<Failure, QuizAttemptEntity>> startAttempt(String quizId);

  Future<Either<Failure, void>> saveProgress({
    required String attemptId,
    required QuizAnswerPayload answers,
  });

  Future<Either<Failure, QuizResultEntity>> submitAttempt({
    required String attemptId,
    required QuizAnswerPayload answers,
  });

  Future<Either<Failure, QuizResultEntity>> getResult(String attemptId);
}
