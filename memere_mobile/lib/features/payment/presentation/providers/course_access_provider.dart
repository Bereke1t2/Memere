import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/enrollment_entity.dart';

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
/// can open lessons for [courseId]. Enrollment check is temporarily disabled.
final courseAccessProvider =
    FutureProvider.family<CourseAccessState, String>((ref, courseId) async {
  // Enrollment checks disabled for now — grant full course access to all students
  return const CourseAccessState(hasAccess: true);
});
