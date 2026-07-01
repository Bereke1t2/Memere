import 'subscription_plan_entity.dart';

enum SubscriptionStatus { active, pastDue, canceled, expired }

class SubscriptionEntity {
  const SubscriptionEntity({
    required this.id,
    required this.plan,
    required this.status,
    this.currentPeriodStart,
    this.currentPeriodEnd,
    this.canceledAt,
  });

  final String id;
  final SubscriptionPlanType plan;
  final SubscriptionStatus status;
  final DateTime? currentPeriodStart;
  final DateTime? currentPeriodEnd;
  final DateTime? canceledAt;

  bool get isCanceled => canceledAt != null || status == SubscriptionStatus.canceled;

  /// Whether the subscription currently grants access. A canceled subscription
  /// still grants access until the current period ends.
  bool get grantsAccess {
    final end = currentPeriodEnd;
    final withinPeriod = end == null || end.isAfter(DateTime.now());
    switch (status) {
      case SubscriptionStatus.active:
      case SubscriptionStatus.pastDue:
        return true;
      case SubscriptionStatus.canceled:
        return withinPeriod;
      case SubscriptionStatus.expired:
        return false;
    }
  }
}
