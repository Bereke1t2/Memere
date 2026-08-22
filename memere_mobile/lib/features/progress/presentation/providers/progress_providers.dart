import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/datasources/progress_remote_datasource.dart';
import '../../data/repositories/progress_repository_impl.dart';
import '../../domain/entities/student_points_entity.dart';
import '../../domain/repositories/progress_repository.dart';
import '../../domain/usecases/get_my_points_usecase.dart';

final progressRemoteDataSourceProvider =
    Provider<ProgressRemoteDataSource>((ref) {
  return ProgressRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final progressRepositoryProvider = Provider<ProgressRepository>((ref) {
  return ProgressRepositoryImpl(ref.watch(progressRemoteDataSourceProvider));
});

final getMyPointsUseCaseProvider = Provider<GetMyPointsUseCase>((ref) {
  return GetMyPointsUseCase(ref.watch(progressRepositoryProvider));
});

/// The signed-in student's cumulative quiz + exam points. Refresh/retry by
/// invalidating. Any failure — including the endpoint not being deployed yet —
/// surfaces as an AsyncError, which the profile renders by hiding the card.
final studentPointsProvider =
    FutureProvider<StudentPointsEntity>((ref) async {
  final result = await ref.watch(getMyPointsUseCaseProvider)();
  return result.fold((failure) => throw failure, (points) => points);
});
