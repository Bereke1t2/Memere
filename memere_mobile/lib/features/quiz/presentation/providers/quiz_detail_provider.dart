import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/offline/offline_providers.dart';
import '../../domain/entities/quiz_entity.dart';
import 'quiz_providers.dart';

final quizDetailProvider = FutureProvider.family<QuizEntity, String>(
  (ref, quizId) async {
    final useCase = ref.watch(getQuizUseCaseProvider);
    final result = await useCase(quizId);
    return result.fold(
      (failure) {
        // Offline (or server error): fall back to the downloaded copy so a
        // quiz lesson inside a downloaded course can still be opened. The
        // "Start Quiz" button's own offline branch then grades it on-device.
        // attemptsUsed is unknown offline, so 0 (Start stays enabled); the
        // real attempt is rebuilt from the OfflineQuiz, not this entity.
        final offline = ref.read(downloadStoreProvider).getOfflineQuiz(quizId);
        if (offline != null) {
          return QuizEntity(
            id: offline.id,
            courseId: offline.courseId,
            title: offline.title,
            timeLimitSeconds: offline.timeLimitSeconds,
            passPercentage: offline.passPercentage,
            randomizeQuestions: false,
            maxAttempts: offline.maxAttempts,
            questionCount: offline.questions.length,
            attemptsUsed: 0,
          );
        }
        throw failure;
      },
      (quiz) => quiz,
    );
  },
);

final quizStartLoadingProvider =
    StateProvider.autoDispose.family<bool, String>((_, __) => false);
