import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/datasources/payment_remote_datasource.dart';
import '../../data/repositories/payment_repository_impl.dart';
import '../../domain/repositories/payment_repository.dart';
import '../../domain/usecases/cancel_subscription_usecase.dart';
import '../../domain/usecases/enroll_free_usecase.dart';
import '../../domain/usecases/get_my_subscription_usecase.dart';
import '../../domain/usecases/get_payment_status_usecase.dart';
import '../../domain/usecases/initiate_course_payment_usecase.dart';
import '../../domain/usecases/initiate_subscription_payment_usecase.dart';
import '../../domain/usecases/list_enrollments_usecase.dart';
import '../../domain/usecases/list_payments_usecase.dart';
import '../../domain/usecases/list_subscription_plans_usecase.dart';

final paymentRemoteDataSourceProvider =
    Provider<PaymentRemoteDataSource>((ref) {
  return PaymentRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final paymentRepositoryProvider = Provider<PaymentRepository>((ref) {
  return PaymentRepositoryImpl(ref.watch(paymentRemoteDataSourceProvider));
});

final enrollFreeUseCaseProvider = Provider<EnrollFreeUseCase>((ref) {
  return EnrollFreeUseCase(ref.watch(paymentRepositoryProvider));
});

final listEnrollmentsUseCaseProvider =
    Provider<ListEnrollmentsUseCase>((ref) {
  return ListEnrollmentsUseCase(ref.watch(paymentRepositoryProvider));
});

final initiateCoursePaymentUseCaseProvider =
    Provider<InitiateCoursePaymentUseCase>((ref) {
  return InitiateCoursePaymentUseCase(ref.watch(paymentRepositoryProvider));
});

final getPaymentStatusUseCaseProvider =
    Provider<GetPaymentStatusUseCase>((ref) {
  return GetPaymentStatusUseCase(ref.watch(paymentRepositoryProvider));
});

final listPaymentsUseCaseProvider = Provider<ListPaymentsUseCase>((ref) {
  return ListPaymentsUseCase(ref.watch(paymentRepositoryProvider));
});

final listSubscriptionPlansUseCaseProvider =
    Provider<ListSubscriptionPlansUseCase>((ref) {
  return ListSubscriptionPlansUseCase(ref.watch(paymentRepositoryProvider));
});

final initiateSubscriptionPaymentUseCaseProvider =
    Provider<InitiateSubscriptionPaymentUseCase>((ref) {
  return InitiateSubscriptionPaymentUseCase(
    ref.watch(paymentRepositoryProvider),
  );
});

final getMySubscriptionUseCaseProvider =
    Provider<GetMySubscriptionUseCase>((ref) {
  return GetMySubscriptionUseCase(ref.watch(paymentRepositoryProvider));
});

final cancelSubscriptionUseCaseProvider =
    Provider<CancelSubscriptionUseCase>((ref) {
  return CancelSubscriptionUseCase(ref.watch(paymentRepositoryProvider));
});
