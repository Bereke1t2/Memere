import '../../domain/entities/subscription_plan_entity.dart';
import 'payment_model_helpers.dart';

class SubscriptionPlanModel extends SubscriptionPlanEntity {
  const SubscriptionPlanModel({
    required super.plan,
    required super.price,
    required super.currency,
    required super.periodDays,
  });

  factory SubscriptionPlanModel.fromJson(Map<String, dynamic> json) {
    return SubscriptionPlanModel(
      plan: parseSubscriptionPlanType(paymentStringValue(json['plan'])),
      price: paymentMoneyValue(json['price']),
      currency: paymentStringValue(json['currency']),
      periodDays: paymentIntValue(json['period_days']),
    );
  }
}
