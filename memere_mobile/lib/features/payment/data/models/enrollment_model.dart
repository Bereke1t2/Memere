import '../../domain/entities/enrollment_entity.dart';
import 'payment_model_helpers.dart';

class EnrollmentModel extends EnrollmentEntity {
  const EnrollmentModel({
    required super.id,
    required super.courseId,
    required super.source,
    super.enrolledAt,
    super.expiresAt,
  });

  factory EnrollmentModel.fromJson(Map<String, dynamic> json) {
    return EnrollmentModel(
      id: paymentStringValue(json['id']),
      courseId: paymentStringValue(json['course_id']),
      source: parseEnrollmentSource(paymentStringValue(json['source'])),
      enrolledAt: paymentDateValue(json['enrolled_at']),
      expiresAt: paymentDateValue(json['expires_at']),
    );
  }
}
