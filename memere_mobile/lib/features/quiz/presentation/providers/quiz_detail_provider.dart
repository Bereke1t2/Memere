import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/quiz_entity.dart';
import 'quiz_providers.dart';

final quizDetailProvider = FutureProvider.family<QuizEntity, String>(
  (ref, quizId) async {
    final useCase = ref.watch(getQuizUseCaseProvider);
    final result = await useCase(quizId);
    return result.fold((failure) => throw failure, (quiz) => quiz);
  },
);

final quizStartLoadingProvider =
    StateProvider.autoDispose.family<bool, String>((_, __) => false);
