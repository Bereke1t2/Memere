import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/payment_initiation_entity.dart';
import '../entities/payment_provider_entity.dart';
import '../entities/subscription_plan_entity.dart';
import '../repositories/payment_repository.dart';

class InitiateSubscriptionPaymentParams {
  const InitiateSubscriptionPaymentParams({
    required this.plan,
    required this.provider,
    required this.idempotencyKey,
    this.couponCode,
  });

  final SubscriptionPlanType plan;
  final PaymentProvider provider;
  final String idempotencyKey;
  final String? couponCode;
}

class InitiateSubscriptionPaymentUseCase {
  const InitiateSubscriptionPaymentUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, PaymentInitiationEntity>> call(
    InitiateSubscriptionPaymentParams params,
  ) {
    if (params.idempotencyKey.trim().isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure(
            'Idempotency key is required',
            field: 'idempotencyKey',
          ),
        ),
      );
    }
    if (!PaymentProviders.isSupported(params.provider)) {
      return Future.value(
        const Left(
          ValidationFailure('Unsupported provider', field: 'provider'),
        ),
      );
    }
    return _repository.initiateSubscriptionPayment(
      plan: params.plan,
      provider: params.provider,
      idempotencyKey: params.idempotencyKey,
      couponCode: params.couponCode,
    );
  }
}
