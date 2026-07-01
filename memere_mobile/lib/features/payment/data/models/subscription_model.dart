import '../../domain/entities/subscription_entity.dart';
import 'payment_model_helpers.dart';

class SubscriptionModel extends SubscriptionEntity {
  const SubscriptionModel({
    required super.id,
    required super.plan,
    required super.status,
    super.currentPeriodStart,
    super.currentPeriodEnd,
    super.canceledAt,
  });

  factory SubscriptionModel.fromJson(Map<String, dynamic> json) {
    return SubscriptionModel(
      id: paymentStringValue(json['id']),
      plan: parseSubscriptionPlanType(paymentStringValue(json['plan'])),
      status: parseSubscriptionStatus(paymentStringValue(json['status'])),
      currentPeriodStart: paymentDateValue(json['current_period_start']),
      currentPeriodEnd: paymentDateValue(json['current_period_end']),
      canceledAt: paymentDateValue(json['canceled_at']),
    );
  }
}
