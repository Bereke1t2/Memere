import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/subscription_plan_entity.dart';
import '../repositories/payment_repository.dart';

class ListSubscriptionPlansUseCase {
  const ListSubscriptionPlansUseCase(this._repository);
  final PaymentRepository _repository;

  Future<Either<Failure, List<SubscriptionPlanEntity>>> call() {
    return _repository.listSubscriptionPlans();
  }
}
