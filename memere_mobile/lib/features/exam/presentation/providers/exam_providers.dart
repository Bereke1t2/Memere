import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/datasources/exam_remote_datasource.dart';
import '../../data/repositories/exam_repository_impl.dart';
import '../../domain/entities/exam_attempt_history_entity.dart';
import '../../domain/entities/mock_exam_entity.dart';
import '../../domain/repositories/exam_repository.dart';
import '../../domain/usecases/get_exam_analytics_usecase.dart';
import '../../domain/usecases/get_exam_result_usecase.dart';
import '../../domain/usecases/list_exams_by_course_usecase.dart';
import '../../domain/usecases/list_mock_exams_usecase.dart';
import '../../domain/usecases/list_my_exam_attempts_usecase.dart';
import '../../domain/usecases/save_exam_progress_usecase.dart';
import '../../domain/usecases/start_exam_usecase.dart';
import '../../domain/usecases/submit_exam_usecase.dart';

final listExamsByCourseUseCaseProvider =
    Provider<ListExamsByCourseUseCase>((ref) {
  return ListExamsByCourseUseCase(ref.watch(examRepositoryProvider));
});

final courseExamsProvider =
    FutureProvider.family<List<MockExamEntity>, String>((ref, courseId) async {
  final useCase = ref.watch(listExamsByCourseUseCaseProvider);
  final result = await useCase(courseId);
  return result.fold(
    (failure) => <MockExamEntity>[],
    (exams) => exams,
  );
});

final examRemoteDataSourceProvider = Provider<ExamRemoteDataSource>((ref) {
  return ExamRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final examRepositoryProvider = Provider<ExamRepository>((ref) {
  return ExamRepositoryImpl(ref.watch(examRemoteDataSourceProvider));
});

final listMockExamsUseCaseProvider = Provider<ListMockExamsUseCase>((ref) {
  return ListMockExamsUseCase(ref.watch(examRepositoryProvider));
});

final startExamUseCaseProvider = Provider<StartExamUseCase>((ref) {
  return StartExamUseCase(ref.watch(examRepositoryProvider));
});

final saveExamProgressUseCaseProvider =
    Provider<SaveExamProgressUseCase>((ref) {
  return SaveExamProgressUseCase(ref.watch(examRepositoryProvider));
});

final submitExamUseCaseProvider = Provider<SubmitExamUseCase>((ref) {
  return SubmitExamUseCase(ref.watch(examRepositoryProvider));
});

final getExamResultUseCaseProvider = Provider<GetExamResultUseCase>((ref) {
  return GetExamResultUseCase(ref.watch(examRepositoryProvider));
});

final getExamAnalyticsUseCaseProvider =
    Provider<GetExamAnalyticsUseCase>((ref) {
  return GetExamAnalyticsUseCase(ref.watch(examRepositoryProvider));
});

final listMyExamAttemptsUseCaseProvider =
    Provider<ListMyExamAttemptsUseCase>((ref) {
  return ListMyExamAttemptsUseCase(ref.watch(examRepositoryProvider));
});

/// Fetches all attempts made by the current user across all mock exams.
final myAllExamAttemptsProvider =
    FutureProvider<List<ExamAttemptHistoryEntity>>((ref) async {
  final useCase = ref.watch(listMyExamAttemptsUseCaseProvider);
  final result = await useCase();
  return result.fold(
    (failure) => <ExamAttemptHistoryEntity>[],
    (attempts) => attempts,
  );
});

/// Fetches attempts made by the user for a specific exam.
final examAttemptsByExamProvider =
    FutureProvider.family<List<ExamAttemptHistoryEntity>, String>((ref, examId) async {
  final useCase = ref.watch(listMyExamAttemptsUseCaseProvider);
  final result = await useCase(examId: examId);
  return result.fold(
    (failure) => <ExamAttemptHistoryEntity>[],
    (attempts) => attempts,
  );
});
