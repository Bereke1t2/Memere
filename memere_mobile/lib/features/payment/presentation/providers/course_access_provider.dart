import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/enrollment_entity.dart';
import 'payment_providers.dart';

/// Resolved access state for a single course, combining enrollments and the
/// active subscription. Refresh by invalidating [courseAccessProvider].
class CourseAccessState {
  const CourseAccessState({
    required this.hasAccess,
    this.enrollment,
    this.viaSubscription = false,
  });

  const CourseAccessState.none()
      : hasAccess = false,
        enrollment = null,
        viaSubscription = false;

  final bool hasAccess;
  final EnrollmentEntity? enrollment;
  final bool viaSubscription;
}

/// Loads enrollments (and the active subscription) to decide whether the user
/// can open lessons for [courseId]. A missing/expired subscription is treated as
/// "no subscription access" rather than an error.
final courseAccessProvider =
    FutureProvider.family<CourseAccessState, String>((ref, courseId) async {
  final enrollmentsResult =
      await ref.watch(listEnrollmentsUseCaseProvider)();

  final enrollment = enrollmentsResult.fold<EnrollmentEntity?>(
    (failure) => null,
    (enrollments) {
      for (final e in enrollments) {
        if (e.courseId == courseId && e.isActive) return e;
      }
      return null;
    },
  );

  if (enrollment != null) {
    return CourseAccessState(hasAccess: true, enrollment: enrollment);
  }

  // Fall back to an all-access subscription if one is active.
  final subscriptionResult =
      await ref.watch(getMySubscriptionUseCaseProvider)();
  final viaSubscription = subscriptionResult.fold(
    (failure) => false,
    (subscription) => subscription.grantsAccess,
  );

  return CourseAccessState(hasAccess: viaSubscription, viaSubscription: viaSubscription);
});
