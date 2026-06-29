import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/quiz_result_entity.dart';
import 'quiz_providers.dart';

final quizResultProvider = FutureProvider.family<QuizResultEntity, String>(
  (ref, attemptId) async {
    final useCase = ref.watch(getQuizResultUseCaseProvider);
    final result = await useCase(attemptId);
    return result.fold((failure) => throw failure, (quizResult) => quizResult);
  },
);
