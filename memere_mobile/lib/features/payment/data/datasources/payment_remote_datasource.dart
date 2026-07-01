import 'package:dio/dio.dart';

import '../../../../core/network/dio_client.dart';
import '../models/enrollment_model.dart';
import '../models/payment_initiation_model.dart';
import '../models/payment_model.dart';
import '../models/subscription_model.dart';
import '../models/subscription_plan_model.dart';

abstract class PaymentRemoteDataSource {
  Future<void> enrollFree(String courseId);

  Future<List<EnrollmentModel>> listEnrollments({int limit = 50});

  Future<PaymentInitiationModel> initiateCoursePayment({
    required String courseId,
    required String provider,
    required String idempotencyKey,
    String? couponCode,
  });

  Future<PaymentModel> getPaymentStatus(String paymentId);

  Future<List<PaymentModel>> listPayments({int limit = 50});

  Future<List<SubscriptionPlanModel>> listSubscriptionPlans();

  Future<PaymentInitiationModel> initiateSubscriptionPayment({
    required String plan,
    required String provider,
    required String idempotencyKey,
    String? couponCode,
  });

  Future<SubscriptionModel> getMySubscription();

  Future<void> cancelSubscription(String subscriptionId);
}

class PaymentRemoteDataSourceImpl implements PaymentRemoteDataSource {
  const PaymentRemoteDataSourceImpl(this._client);
  final DioClient _client;

  @override
  Future<void> enrollFree(String courseId) async {
    await _client.post('/courses/$courseId/enroll-free');
  }

  @override
  Future<List<EnrollmentModel>> listEnrollments({int limit = 50}) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/me/enrollments',
      queryParameters: {'limit': limit},
    );
    return _parseList(response.data, EnrollmentModel.fromJson);
  }

  @override
  Future<PaymentInitiationModel> initiateCoursePayment({
    required String courseId,
    required String provider,
    required String idempotencyKey,
    String? couponCode,
  }) async {
    final response = await _client.post<Map<String, dynamic>>(
      '/payments/initiate',
      data: {
        'course_id': courseId,
        'provider': provider,
        if (couponCode != null && couponCode.isNotEmpty)
          'coupon_code': couponCode,
      },
      options: Options(headers: {'Idempotency-Key': idempotencyKey}),
    );
    return PaymentInitiationModel.fromJson(_requireBody(response.data));
  }

  @override
  Future<PaymentModel> getPaymentStatus(String paymentId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/payments/$paymentId/status',
    );
    return PaymentModel.fromJson(_requireBody(response.data));
  }

  @override
  Future<List<PaymentModel>> listPayments({int limit = 50}) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/payments',
      queryParameters: {'limit': limit},
    );
    return _parseList(response.data, PaymentModel.fromJson);
  }

  @override
  Future<List<SubscriptionPlanModel>> listSubscriptionPlans() async {
    final response = await _client.get<Map<String, dynamic>>(
      '/subscription-plans',
    );
    return _parseList(response.data, SubscriptionPlanModel.fromJson);
  }

  @override
  Future<PaymentInitiationModel> initiateSubscriptionPayment({
    required String plan,
    required String provider,
    required String idempotencyKey,
    String? couponCode,
  }) async {
    final response = await _client.post<Map<String, dynamic>>(
      '/subscriptions',
      data: {
        'plan': plan,
        'provider': provider,
        if (couponCode != null && couponCode.isNotEmpty)
          'coupon_code': couponCode,
      },
      options: Options(headers: {'Idempotency-Key': idempotencyKey}),
    );
    return PaymentInitiationModel.fromJson(_requireBody(response.data));
  }

  @override
  Future<SubscriptionModel> getMySubscription() async {
    final response = await _client.get<Map<String, dynamic>>(
      '/me/subscription',
    );
    return SubscriptionModel.fromJson(_requireBody(response.data));
  }

  @override
  Future<void> cancelSubscription(String subscriptionId) async {
    await _client.post('/subscriptions/$subscriptionId/cancel');
  }

  Map<String, dynamic> _requireBody(Map<String, dynamic>? data) {
    if (data == null) {
      throw const FormatException('Missing payment response body');
    }
    return data;
  }

  List<T> _parseList<T>(
    Map<String, dynamic>? data,
    T Function(Map<String, dynamic>) fromJson,
  ) {
    final items = data?['data'];
    if (items is! List) return <T>[];
    return items
        .whereType<Map<String, dynamic>>()
        .map(fromJson)
        .toList(growable: false);
  }
}
