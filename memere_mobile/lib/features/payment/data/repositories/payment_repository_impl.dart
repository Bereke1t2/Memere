import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/enrollment_entity.dart';
import '../../domain/entities/payment_entity.dart';
import '../../domain/entities/payment_initiation_entity.dart';
import '../../domain/entities/payment_provider_entity.dart';
import '../../domain/entities/subscription_entity.dart';
import '../../domain/entities/subscription_plan_entity.dart';
import '../../domain/repositories/payment_repository.dart';
import '../datasources/payment_remote_datasource.dart';

class PaymentRepositoryImpl implements PaymentRepository {
  const PaymentRepositoryImpl(this._remoteDataSource);
  final PaymentRemoteDataSource _remoteDataSource;

  @override
  Future<Either<Failure, void>> enrollFree(String courseId) async {
    try {
      await _remoteDataSource.enrollFree(courseId);
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<EnrollmentEntity>>> listEnrollments({
    int limit = 50,
  }) async {
    try {
      final enrollments = await _remoteDataSource.listEnrollments(limit: limit);
      return Right(enrollments);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, PaymentInitiationEntity>> initiateCoursePayment({
    required String courseId,
    required PaymentProvider provider,
    required String idempotencyKey,
    String? couponCode,
  }) async {
    try {
      final initiation = await _remoteDataSource.initiateCoursePayment(
        courseId: courseId,
        provider: provider.apiValue,
        idempotencyKey: idempotencyKey,
        couponCode: couponCode,
      );
      return Right(initiation);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, PaymentEntity>> getPaymentStatus(
    String paymentId,
  ) async {
    try {
      final payment = await _remoteDataSource.getPaymentStatus(paymentId);
      return Right(payment);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<PaymentEntity>>> listPayments({
    int limit = 50,
  }) async {
    try {
      final payments = await _remoteDataSource.listPayments(limit: limit);
      return Right(payments);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<SubscriptionPlanEntity>>>
      listSubscriptionPlans() async {
    try {
      final plans = await _remoteDataSource.listSubscriptionPlans();
      return Right(plans);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, PaymentInitiationEntity>> initiateSubscriptionPayment({
    required SubscriptionPlanType plan,
    required PaymentProvider provider,
    required String idempotencyKey,
    String? couponCode,
  }) async {
    try {
      final initiation = await _remoteDataSource.initiateSubscriptionPayment(
        plan: plan.apiValue,
        provider: provider.apiValue,
        idempotencyKey: idempotencyKey,
        couponCode: couponCode,
      );
      return Right(initiation);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, SubscriptionEntity>> getMySubscription() async {
    try {
      final subscription = await _remoteDataSource.getMySubscription();
      return Right(subscription);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> cancelSubscription(
    String subscriptionId,
  ) async {
    try {
      await _remoteDataSource.cancelSubscription(subscriptionId);
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }
}
