import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:uuid/uuid.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/payment_initiation_entity.dart';
import '../../domain/entities/payment_provider_entity.dart';
import '../../domain/entities/subscription_entity.dart';
import '../../domain/entities/subscription_plan_entity.dart';
import '../../domain/usecases/initiate_subscription_payment_usecase.dart';
import 'payment_providers.dart';

const _uuid = Uuid();

/// All-access subscription plans.
final subscriptionPlansProvider =
    FutureProvider<List<SubscriptionPlanEntity>>((ref) async {
  final result = await ref.watch(listSubscriptionPlansUseCaseProvider)();
  return result.fold((failure) => throw failure, (plans) => plans);
});

/// The user's active subscription, or null if they have none. A missing
/// subscription (404) resolves to null rather than an error.
final mySubscriptionProvider =
    FutureProvider<SubscriptionEntity?>((ref) async {
  final result = await ref.watch(getMySubscriptionUseCaseProvider)();
  return result.fold(
    (failure) {
      if (failure is ServerFailure && failure.statusCode == 404) return null;
      throw failure;
    },
    (subscription) => subscription,
  );
});

class SubscriptionCheckoutState {
  const SubscriptionCheckoutState({
    this.selectedPlan,
    this.selectedProvider,
    this.idempotencyKey,
    this.initiation,
    this.isWorking = false,
    this.error,
  });

  final SubscriptionPlanType? selectedPlan;
  final PaymentProvider? selectedProvider;
  final String? idempotencyKey;
  final PaymentInitiationEntity? initiation;
  final bool isWorking;
  final String? error;

  SubscriptionCheckoutState copyWith({
    SubscriptionPlanType? selectedPlan,
    PaymentProvider? selectedProvider,
    String? idempotencyKey,
    PaymentInitiationEntity? initiation,
    bool? isWorking,
    String? error,
    bool clearError = false,
  }) {
    return SubscriptionCheckoutState(
      selectedPlan: selectedPlan ?? this.selectedPlan,
      selectedProvider: selectedProvider ?? this.selectedProvider,
      idempotencyKey: idempotencyKey ?? this.idempotencyKey,
      initiation: initiation ?? this.initiation,
      isWorking: isWorking ?? this.isWorking,
      error: clearError ? null : error ?? this.error,
    );
  }
}

final subscriptionCheckoutProvider = AsyncNotifierProvider<
    SubscriptionCheckoutNotifier, SubscriptionCheckoutState>(
  SubscriptionCheckoutNotifier.new,
);

/// Mirrors the course checkout flow but for subscriptions: one stable
/// idempotency key per attempt, reused on retry.
class SubscriptionCheckoutNotifier
    extends AsyncNotifier<SubscriptionCheckoutState> {
  @override
  Future<SubscriptionCheckoutState> build() async {
    return const SubscriptionCheckoutState();
  }

  SubscriptionCheckoutState get _current =>
      state.valueOrNull ?? const SubscriptionCheckoutState();

  void selectPlan(SubscriptionPlanType plan) {
    state = AsyncData(_current.copyWith(selectedPlan: plan, clearError: true));
  }

  Future<PaymentInitiationEntity?> startCheckout({
    required SubscriptionPlanType plan,
    required PaymentProvider provider,
    String? couponCode,
  }) async {
    final key = _current.idempotencyKey ?? _uuid.v4();
    state = AsyncData(
      _current.copyWith(
        selectedPlan: plan,
        selectedProvider: provider,
        idempotencyKey: key,
        isWorking: true,
        clearError: true,
      ),
    );

    final result =
        await ref.read(initiateSubscriptionPaymentUseCaseProvider)(
      InitiateSubscriptionPaymentParams(
        plan: plan,
        provider: provider,
        idempotencyKey: key,
        couponCode: couponCode,
      ),
    );

    return result.fold(
      (failure) {
        state = AsyncData(
          _current.copyWith(isWorking: false, error: failure.message),
        );
        return null;
      },
      (initiation) {
        state = AsyncData(
          _current.copyWith(isWorking: false, initiation: initiation),
        );
        return initiation;
      },
    );
  }

  void reset() {
    state = const AsyncData(SubscriptionCheckoutState());
  }
}
