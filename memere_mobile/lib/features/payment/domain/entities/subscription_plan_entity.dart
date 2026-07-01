enum SubscriptionPlanType { monthly, annual }

extension SubscriptionPlanTypeX on SubscriptionPlanType {
  /// Wire value sent to the backend.
  String get apiValue {
    switch (this) {
      case SubscriptionPlanType.monthly:
        return 'monthly';
      case SubscriptionPlanType.annual:
        return 'annual';
    }
  }

  String get label {
    switch (this) {
      case SubscriptionPlanType.monthly:
        return 'Monthly';
      case SubscriptionPlanType.annual:
        return 'Annual';
    }
  }
}

class SubscriptionPlanEntity {
  const SubscriptionPlanEntity({
    required this.plan,
    required this.price,
    required this.currency,
    required this.periodDays,
  });

  final SubscriptionPlanType plan;
  final String price;
  final String currency;
  final int periodDays;

  String get priceLabel {
    if (price.trim().isEmpty) return currency;
    return '$currency $price';
  }

  String get periodLabel {
    if (periodDays <= 0) return '';
    if (periodDays % 30 == 0) {
      final months = periodDays ~/ 30;
      return months == 1 ? 'per month' : 'every $months months';
    }
    if (periodDays % 365 == 0) {
      final years = periodDays ~/ 365;
      return years == 1 ? 'per year' : 'every $years years';
    }
    return 'every $periodDays days';
  }
}
