import '../../domain/entities/payment_initiation_entity.dart';
import 'payment_model_helpers.dart';

class PaymentInitiationModel extends PaymentInitiationEntity {
  const PaymentInitiationModel({
    required super.paymentId,
    required super.redirectUrl,
    required super.amount,
    required super.currency,
  });

  factory PaymentInitiationModel.fromJson(Map<String, dynamic> json) {
    return PaymentInitiationModel(
      paymentId: paymentStringValue(json['payment_id'] ?? json['id']),
      redirectUrl: paymentStringValue(json['redirect_url']),
      amount: paymentMoneyValue(json['amount']),
      currency: paymentStringValue(json['currency']),
    );
  }
}
