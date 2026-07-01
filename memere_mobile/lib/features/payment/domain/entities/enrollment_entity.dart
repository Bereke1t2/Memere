enum EnrollmentSource { purchase, subscription, free, coupon }

class EnrollmentEntity {
  const EnrollmentEntity({
    required this.id,
    required this.courseId,
    required this.source,
    this.enrolledAt,
    this.expiresAt,
  });

  final String id;
  final String courseId;
  final EnrollmentSource source;
  final DateTime? enrolledAt;
  final DateTime? expiresAt;

  /// An enrollment with no expiry is treated as lifetime access. One with an
  /// expiry is active only while that expiry is in the future.
  bool get isActive {
    final expiry = expiresAt;
    if (expiry == null) return true;
    return expiry.isAfter(DateTime.now());
  }
}
