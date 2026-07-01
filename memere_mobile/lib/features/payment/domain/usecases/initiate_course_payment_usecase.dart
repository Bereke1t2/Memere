import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/payment_initiation_entity.dart';
import '../entities/payment_provider_entity.dart';
import '../repositories/payment_repository.dart';

class InitiateCoursePaymentParams {
  const InitiateCoursePaymentParams({
    required this.courseId,
    required this.provider,
    required this.idempotencyKey,
    this.couponCode,
  });

  final String courseId;
  final PaymentProvider provider;
  final String idempotencyKey;
  final String? couponCode;
}

class InitiateCoursePaymentUseCase {
  const InitiateCoursePaymentUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, PaymentInitiationEntity>> call(
    InitiateCoursePaymentParams params,
  ) {
    final courseId = params.courseId.trim();
    if (courseId.isEmpty) {
      return Future.value(
        const Left(
          ValidationFailure('Course ID is required', field: 'courseId'),
        ),
      );
    }
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
    return _repository.initiateCoursePayment(
      courseId: courseId,
      provider: params.provider,
      idempotencyKey: params.idempotencyKey,
      couponCode: params.couponCode,
    );
  }
}
