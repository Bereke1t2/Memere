import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/enrollment_entity.dart';
import '../entities/payment_entity.dart';
import '../entities/payment_initiation_entity.dart';
import '../entities/payment_provider_entity.dart';
import '../entities/subscription_entity.dart';
import '../entities/subscription_plan_entity.dart';

abstract class PaymentRepository {
  Future<Either<Failure, void>> enrollFree(String courseId);

  Future<Either<Failure, List<EnrollmentEntity>>> listEnrollments({
    int limit = 50,
  });

  Future<Either<Failure, PaymentInitiationEntity>> initiateCoursePayment({
    required String courseId,
    required PaymentProvider provider,
    required String idempotencyKey,
    String? couponCode,
  });

  Future<Either<Failure, PaymentEntity>> getPaymentStatus(String paymentId);

  Future<Either<Failure, List<PaymentEntity>>> listPayments({
    int limit = 50,
  });

  Future<Either<Failure, List<SubscriptionPlanEntity>>> listSubscriptionPlans();

  Future<Either<Failure, PaymentInitiationEntity>> initiateSubscriptionPayment({
    required SubscriptionPlanType plan,
    required PaymentProvider provider,
    required String idempotencyKey,
    String? couponCode,
  });

  Future<Either<Failure, SubscriptionEntity>> getMySubscription();

  Future<Either<Failure, void>> cancelSubscription(String subscriptionId);
}
