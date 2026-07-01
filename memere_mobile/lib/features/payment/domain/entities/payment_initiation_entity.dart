class PaymentInitiationEntity {
  const PaymentInitiationEntity({
    required this.paymentId,
    required this.redirectUrl,
    required this.amount,
    required this.currency,
  });

  final String paymentId;

  /// Provider-hosted checkout URL. Only ever opened inside the payment WebView.
  final String redirectUrl;
  final String amount;
  final String currency;

  String get amountLabel {
    if (amount.trim().isEmpty) return currency;
    return '$currency $amount';
  }
}
