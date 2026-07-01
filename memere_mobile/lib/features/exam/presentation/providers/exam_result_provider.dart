import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/exam_result_entity.dart';
import 'exam_providers.dart';

final examResultProvider = FutureProvider.family<ExamResultEntity, String>(
  (ref, attemptId) async {
    final useCase = ref.watch(getExamResultUseCaseProvider);
    final result = await useCase(attemptId);
    return result.fold((failure) => throw failure, (examResult) => examResult);
  },
);
