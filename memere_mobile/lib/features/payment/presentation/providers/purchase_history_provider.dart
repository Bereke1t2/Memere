import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../courses/presentation/providers/course_list_provider.dart';
import '../../domain/entities/enrollment_entity.dart';
import '../../domain/entities/payment_entity.dart';
import 'payment_providers.dart';

/// Payment history for the signed-in user. Refresh/retry by invalidating.
final paymentHistoryProvider =
    FutureProvider<List<PaymentEntity>>((ref) async {
  final result = await ref.watch(listPaymentsUseCaseProvider)();
  return result.fold((failure) => throw failure, (payments) => payments);
});

/// Enrollment list for the signed-in user.
final enrollmentListProvider =
    FutureProvider<List<EnrollmentEntity>>((ref) async {
  final result = await ref.watch(listEnrollmentsUseCaseProvider)();
  return result.fold(
    (failure) => throw failure,
    (enrollments) => enrollments,
  );
});

/// Resolves a human course title from the already-loaded catalog cache, so
/// history rows show real names instead of raw IDs. Returns null if the course
/// is not in the loaded catalog (no extra network call is made).
final courseTitleProvider = Provider.family<String?, String>((ref, courseId) {
  if (courseId.isEmpty) return null;
  final listState = ref.watch(courseListProvider).valueOrNull;
  if (listState == null) return null;
  for (final course in listState.courses) {
    if (course.id == courseId) return course.title;
  }
  return null;
});

/// Latest payment for a given course, if any — used to surface a pending state
/// on the course detail CTA.
final latestCoursePaymentProvider =
    Provider.family<PaymentEntity?, String>((ref, courseId) {
  final payments = ref.watch(paymentHistoryProvider).valueOrNull;
  if (payments == null) return null;
  for (final payment in payments) {
    if (payment.courseId == courseId) return payment;
  }
  return null;
});
