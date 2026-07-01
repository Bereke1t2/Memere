import '../../domain/entities/enrollment_entity.dart';
import '../../domain/entities/payment_entity.dart';
import '../../domain/entities/subscription_entity.dart';
import '../../domain/entities/subscription_plan_entity.dart';

String paymentStringValue(Object? value, {String fallback = ''}) {
  if (value is String) return value;
  if (value == null) return fallback;
  return value.toString();
}

String? paymentNullableString(Object? value) {
  if (value is String && value.trim().isNotEmpty) return value;
  return null;
}

/// Money values arrive as strings; numeric values are normalised to a string so
/// the rest of the app never loses precision by parsing into a double.
String paymentMoneyValue(Object? value, {String fallback = ''}) {
  if (value is String) return value;
  if (value is num) {
    return value % 1 == 0 ? value.toInt().toString() : value.toString();
  }
  if (value == null) return fallback;
  return value.toString();
}

int paymentIntValue(Object? value, {int fallback = 0}) {
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value) ?? fallback;
  return fallback;
}

DateTime? paymentDateValue(Object? value) {
  if (value is! String || value.isEmpty) return null;
  return DateTime.tryParse(value);
}

PaymentStatus parsePaymentStatus(String value) {
  switch (value.toLowerCase()) {
    case 'completed':
    case 'succeeded':
    case 'success':
      return PaymentStatus.completed;
    case 'failed':
    case 'canceled':
    case 'cancelled':
      return PaymentStatus.failed;
    case 'refunded':
      return PaymentStatus.refunded;
    default:
      return PaymentStatus.pending;
  }
}

EnrollmentSource parseEnrollmentSource(String value) {
  switch (value.toLowerCase()) {
    case 'subscription':
      return EnrollmentSource.subscription;
    case 'free':
      return EnrollmentSource.free;
    case 'coupon':
      return EnrollmentSource.coupon;
    default:
      return EnrollmentSource.purchase;
  }
}

SubscriptionStatus parseSubscriptionStatus(String value) {
  switch (value.toLowerCase()) {
    case 'active':
      return SubscriptionStatus.active;
    case 'past_due':
    case 'pastdue':
      return SubscriptionStatus.pastDue;
    case 'canceled':
    case 'cancelled':
      return SubscriptionStatus.canceled;
    default:
      return SubscriptionStatus.expired;
  }
}

SubscriptionPlanType parseSubscriptionPlanType(String value) {
  switch (value.toLowerCase()) {
    case 'annual':
    case 'yearly':
      return SubscriptionPlanType.annual;
    default:
      return SubscriptionPlanType.monthly;
  }
}
