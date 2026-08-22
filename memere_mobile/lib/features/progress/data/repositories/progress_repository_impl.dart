import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/student_points_entity.dart';
import '../../domain/repositories/progress_repository.dart';
import '../datasources/progress_remote_datasource.dart';

class ProgressRepositoryImpl implements ProgressRepository {
  const ProgressRepositoryImpl(this._remoteDataSource);
  final ProgressRemoteDataSource _remoteDataSource;

  @override
  Future<Either<Failure, StudentPointsEntity>> getMyPoints() async {
    try {
      final points = await _remoteDataSource.getMyPoints();
      return Right(points);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }
}
