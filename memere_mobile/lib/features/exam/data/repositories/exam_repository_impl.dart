import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/exam_answer_payload.dart';
import '../../domain/entities/exam_attempt_analytics_entity.dart';
import '../../domain/entities/exam_attempt_entity.dart';
import '../../domain/entities/exam_result_entity.dart';
import '../../domain/entities/paginated_mock_exams_entity.dart';
import '../../domain/repositories/exam_repository.dart';
import '../datasources/exam_remote_datasource.dart';

class ExamRepositoryImpl implements ExamRepository {
  const ExamRepositoryImpl(this._remoteDataSource);
  final ExamRemoteDataSource _remoteDataSource;

  @override
  Future<Either<Failure, PaginatedMockExamsEntity>> listMockExams({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  }) async {
    try {
      final page = await _remoteDataSource.listMockExams(
        limit: limit,
        after: after,
        subject: subject,
        grade: grade,
      );
      return Right(page);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, ExamAttemptEntity>> startExam(String examId) async {
    try {
      final attempt = await _remoteDataSource.startExam(examId);
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
    required ExamAnswerPayload answers,
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
  Future<Either<Failure, ExamResultEntity>> submitExam({
    required String attemptId,
    required ExamAnswerPayload answers,
  }) async {
    try {
      final result = await _remoteDataSource.submitExam(
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
  Future<Either<Failure, ExamResultEntity>> getResult(String attemptId) async {
    try {
      final result = await _remoteDataSource.getResult(attemptId);
      return Right(result);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, ExamAttemptAnalyticsEntity>> getAnalytics(
    String attemptId,
  ) async {
    try {
      final analytics = await _remoteDataSource.getAnalytics(attemptId);
      return Right(analytics);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }
}
