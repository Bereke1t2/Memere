import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/subscription_entity.dart';
import '../repositories/payment_repository.dart';

class GetMySubscriptionUseCase {
  const GetMySubscriptionUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, SubscriptionEntity>> call() {
    return _repository.getMySubscription();
  }
}
