import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../domain/entities/enrollment_entity.dart';

/// Resolved access state for a single course, combining server access status,
/// enrollments and requests.
class CourseAccessState {
  const CourseAccessState({
    required this.hasAccess,
    this.accessLevel = 'none',
    this.status = 'not_enrolled',
    this.requestId,
    this.rejectionReason,
    this.enrollment,
    this.viaSubscription = false,
  });

  const CourseAccessState.none()
      : hasAccess = false,
        accessLevel = 'none',
        status = 'not_enrolled',
        requestId = null,
        rejectionReason = null,
        enrollment = null,
        viaSubscription = false;

  final bool hasAccess;
  final String accessLevel;
  final String status; // 'granted', 'account_pending', 'request_pending', 'not_enrolled', 'rejected', 'anonymous'
  final String? requestId;
  final String? rejectionReason;
  final EnrollmentEntity? enrollment;
  final bool viaSubscription;

  bool get isAccountPending => status == 'account_pending';
  bool get isRequestPending => status == 'request_pending';
  bool get isRejected => status == 'rejected';
  bool get isGranted => hasAccess || status == 'granted';
}

/// Queries real course access status from the backend.
final courseAccessProvider =
    FutureProvider.family<CourseAccessState, String>((ref, courseId) async {
  try {
    final dio = ref.watch(dioClientProvider);
    final res = await dio.get('/courses/$courseId/access-status');
    final data = res.data as Map<String, dynamic>;
    final hasAccess = data['has_access'] as bool? ?? false;
    final accessLevel = data['access_level'] as String? ?? (hasAccess ? 'full' : 'none');
    final status = data['status'] as String? ?? (hasAccess ? 'granted' : 'not_enrolled');
    final requestId = data['request_id'] as String?;
    final rejectionReason = data['rejection_reason'] as String?;

    return CourseAccessState(
      hasAccess: hasAccess,
      accessLevel: accessLevel,
      status: status,
      requestId: requestId,
      rejectionReason: rejectionReason,
    );
  } catch (_) {
    return const CourseAccessState.none();
  }
});
