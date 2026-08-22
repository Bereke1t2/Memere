import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../data/models/quiz_result_model.dart';
import '../../domain/entities/quiz_result_entity.dart';
import 'quiz_providers.dart';

final quizResultProvider = FutureProvider.family<QuizResultEntity, String>(
  (ref, attemptId) async {
    // On-device-graded attempts resolve from the local store; a `local-…` id
    // never maps to a server attempt, so it must not hit the network.
    if (isLocalAttemptId(attemptId)) {
      final json = ref.read(offlineAttemptStoreProvider).getResult(attemptId);
      if (json == null) {
        throw const CacheFailure('This offline result is no longer available.');
      }
      return QuizResultModel.fromJson(json);
    }

    final useCase = ref.watch(getQuizResultUseCaseProvider);
    final result = await useCase(attemptId);
    return result.fold((failure) => throw failure, (quizResult) => quizResult);
  },
);
