import 'package:flutter/foundation.dart';

enum PaymentProvider { chapa, telebirr, stripe, mock }

extension PaymentProviderX on PaymentProvider {
  /// Wire value sent to the backend.
  String get apiValue {
    switch (this) {
      case PaymentProvider.chapa:
        return 'chapa';
      case PaymentProvider.telebirr:
        return 'telebirr';
      case PaymentProvider.stripe:
        return 'stripe';
      case PaymentProvider.mock:
        return 'mock';
    }
  }

  String get label {
    switch (this) {
      case PaymentProvider.chapa:
        return 'Chapa';
      case PaymentProvider.telebirr:
        return 'Telebirr';
      case PaymentProvider.stripe:
        return 'Stripe';
      case PaymentProvider.mock:
        return 'Mock';
    }
  }

  String get description {
    switch (this) {
      case PaymentProvider.chapa:
        return 'Cards, bank & mobile money';
      case PaymentProvider.telebirr:
        return 'Pay with your Telebirr wallet';
      case PaymentProvider.stripe:
        return 'International cards';
      case PaymentProvider.mock:
        return 'Local test payments only';
    }
  }
}

/// Provider selection rules for the checkout sheet.
abstract class PaymentProviders {
  /// Default Ethiopia provider order.
  static const List<PaymentProvider> ethiopiaOrder = [
    PaymentProvider.chapa,
    PaymentProvider.telebirr,
    PaymentProvider.stripe,
  ];

  /// Providers visible to the user. `mock` is only offered in debug/dev builds.
  static List<PaymentProvider> visible() {
    if (kReleaseMode) return ethiopiaOrder;
    return const [...ethiopiaOrder, PaymentProvider.mock];
  }

  static bool isSupported(PaymentProvider provider) {
    return visible().contains(provider);
  }
}
