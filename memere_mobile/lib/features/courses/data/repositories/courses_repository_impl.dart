import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/course_detail_entity.dart';
import '../../domain/entities/paginated_courses_entity.dart';
import '../../domain/repositories/courses_repository.dart';
import '../datasources/courses_remote_datasource.dart';

class CoursesRepositoryImpl implements CoursesRepository {
  const CoursesRepositoryImpl(this._remoteDataSource);
  final CoursesRemoteDataSource _remoteDataSource;

  @override
  Future<Either<Failure, PaginatedCoursesEntity>> listCourses({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  }) async {
    try {
      final result = await _remoteDataSource.listCourses(
        limit: limit,
        after: after,
        subject: subject,
        grade: grade,
      );
      return Right(result);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, CourseDetailEntity>> getCourseDetail(
    String courseId,
  ) async {
    try {
      final result = await _remoteDataSource.getCourseDetail(courseId);
      return Right(result);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }
}
