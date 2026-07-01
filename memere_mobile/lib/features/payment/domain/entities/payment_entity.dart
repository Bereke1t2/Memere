enum PaymentStatus { pending, completed, failed, refunded }

class PaymentEntity {
  const PaymentEntity({
    required this.paymentId,
    required this.status,
    required this.amount,
    required this.currency,
    required this.courseId,
  });

  final String paymentId;
  final PaymentStatus status;

  /// Money value kept as a string exactly as the backend returns it.
  final String amount;
  final String currency;
  final String? courseId;

  bool get isPending => status == PaymentStatus.pending;
  bool get isCompleted => status == PaymentStatus.completed;
  bool get isFailed => status == PaymentStatus.failed;
  bool get isRefunded => status == PaymentStatus.refunded;

  /// Terminal states no longer need polling.
  bool get isTerminal => isCompleted || isFailed || isRefunded;

  String get amountLabel {
    if (amount.trim().isEmpty) return currency;
    return '$currency $amount';
  }
}
