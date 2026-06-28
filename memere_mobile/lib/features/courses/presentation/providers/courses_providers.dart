import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/datasources/courses_remote_datasource.dart';
import '../../data/repositories/courses_repository_impl.dart';
import '../../domain/repositories/courses_repository.dart';
import '../../domain/usecases/get_course_detail_usecase.dart';
import '../../domain/usecases/list_courses_usecase.dart';

final coursesRemoteDataSourceProvider =
    Provider<CoursesRemoteDataSource>((ref) {
  return CoursesRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final coursesRepositoryProvider = Provider<CoursesRepository>((ref) {
  return CoursesRepositoryImpl(ref.watch(coursesRemoteDataSourceProvider));
});

final listCoursesUseCaseProvider = Provider<ListCoursesUseCase>((ref) {
  return ListCoursesUseCase(ref.watch(coursesRepositoryProvider));
});

final getCourseDetailUseCaseProvider = Provider<GetCourseDetailUseCase>((ref) {
  return GetCourseDetailUseCase(ref.watch(coursesRepositoryProvider));
});
