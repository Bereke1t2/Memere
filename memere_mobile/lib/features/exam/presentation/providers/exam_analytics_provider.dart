import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/exam_attempt_analytics_entity.dart';
import 'exam_providers.dart';

final examAnalyticsProvider =
    FutureProvider.family<ExamAttemptAnalyticsEntity, String>(
  (ref, attemptId) async {
    final useCase = ref.watch(getExamAnalyticsUseCaseProvider);
    final result = await useCase(attemptId);
    return result.fold((failure) => throw failure, (analytics) => analytics);
  },
);
