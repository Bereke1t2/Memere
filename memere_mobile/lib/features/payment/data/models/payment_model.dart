import '../../domain/entities/payment_entity.dart';
import 'payment_model_helpers.dart';

class PaymentModel extends PaymentEntity {
  const PaymentModel({
    required super.paymentId,
    required super.status,
    required super.amount,
    required super.currency,
    required super.courseId,
  });

  factory PaymentModel.fromJson(Map<String, dynamic> json) {
    return PaymentModel(
      paymentId: paymentStringValue(json['payment_id'] ?? json['id']),
      status: parsePaymentStatus(paymentStringValue(json['status'])),
      amount: paymentMoneyValue(json['amount']),
      currency: paymentStringValue(json['currency']),
      courseId: paymentNullableString(json['course_id']),
    );
  }

  Map<String, dynamic> toJson() => {
        'payment_id': paymentId,
        'status': status.name,
        'amount': amount,
        'currency': currency,
        'course_id': courseId,
      };
}
