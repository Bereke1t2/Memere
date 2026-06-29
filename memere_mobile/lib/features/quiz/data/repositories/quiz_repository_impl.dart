import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/quiz_answer_payload.dart';
import '../../domain/entities/quiz_attempt_entity.dart';
import '../../domain/entities/quiz_entity.dart';
import '../../domain/entities/quiz_result_entity.dart';
import '../../domain/repositories/quiz_repository.dart';
import '../datasources/quiz_remote_datasource.dart';

class QuizRepositoryImpl implements QuizRepository {
  const QuizRepositoryImpl(this._remoteDataSource);
  final QuizRemoteDataSource _remoteDataSource;

  @override
  Future<Either<Failure, List<QuizEntity>>> listQuizzesByCourse(
    String courseId,
  ) async {
    try {
      final quizzes = await _remoteDataSource.listQuizzesByCourse(courseId);
      return Right(quizzes);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, QuizEntity>> getQuiz(String quizId) async {
    try {
      final quiz = await _remoteDataSource.getQuiz(quizId);
      return Right(quiz);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, QuizAttemptEntity>> startAttempt(String quizId) async {
    try {
      final attempt = await _remoteDataSource.startAttempt(quizId);
      return Right(attempt);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> saveProgress({
    required String attemptId,
    required QuizAnswerPayload answers,
  }) async {
    try {
      await _remoteDataSource.saveProgress(
        attemptId: attemptId,
        answers: answers,
      );
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, QuizResultEntity>> submitAttempt({
    required String attemptId,
    required QuizAnswerPayload answers,
  }) async {
    try {
      final result = await _remoteDataSource.submitAttempt(
        attemptId: attemptId,
        answers: answers,
      );
      return Right(result);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, QuizResultEntity>> getResult(String attemptId) async {
    try {
      final result = await _remoteDataSource.getResult(attemptId);
      return Right(result);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }
}
