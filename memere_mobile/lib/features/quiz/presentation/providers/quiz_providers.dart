import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/datasources/quiz_remote_datasource.dart';
import '../../data/repositories/quiz_repository_impl.dart';
import '../../domain/entities/quiz_entity.dart';
import '../../domain/repositories/quiz_repository.dart';
import '../../domain/usecases/get_quiz_result_usecase.dart';
import '../../domain/usecases/get_quiz_usecase.dart';
import '../../domain/usecases/list_quizzes_by_course_usecase.dart';
import '../../domain/usecases/save_quiz_progress_usecase.dart';
import '../../domain/usecases/start_quiz_attempt_usecase.dart';
import '../../domain/usecases/submit_quiz_attempt_usecase.dart';

final quizRemoteDataSourceProvider = Provider<QuizRemoteDataSource>((ref) {
  return QuizRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final quizRepositoryProvider = Provider<QuizRepository>((ref) {
  return QuizRepositoryImpl(ref.watch(quizRemoteDataSourceProvider));
});

final listQuizzesByCourseUseCaseProvider =
    Provider<ListQuizzesByCourseUseCase>((ref) {
  return ListQuizzesByCourseUseCase(ref.watch(quizRepositoryProvider));
});

final getQuizUseCaseProvider = Provider<GetQuizUseCase>((ref) {
  return GetQuizUseCase(ref.watch(quizRepositoryProvider));
});

final startQuizAttemptUseCaseProvider =
    Provider<StartQuizAttemptUseCase>((ref) {
  return StartQuizAttemptUseCase(ref.watch(quizRepositoryProvider));
});

final saveQuizProgressUseCaseProvider =
    Provider<SaveQuizProgressUseCase>((ref) {
  return SaveQuizProgressUseCase(ref.watch(quizRepositoryProvider));
});

final submitQuizAttemptUseCaseProvider =
    Provider<SubmitQuizAttemptUseCase>((ref) {
  return SubmitQuizAttemptUseCase(ref.watch(quizRepositoryProvider));
});

final getQuizResultUseCaseProvider = Provider<GetQuizResultUseCase>((ref) {
  return GetQuizResultUseCase(ref.watch(quizRepositoryProvider));
});

final courseQuizzesProvider =
    FutureProvider.family<List<QuizEntity>, String>((ref, courseId) async {
  final useCase = ref.watch(listQuizzesByCourseUseCaseProvider);
  final result = await useCase(courseId);
  return result.fold(
    (failure) => <QuizEntity>[],
    (quizzes) => quizzes,
  );
});
