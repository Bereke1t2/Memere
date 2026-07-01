import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:uuid/uuid.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/payment_initiation_entity.dart';
import '../../domain/entities/payment_provider_entity.dart';
import '../../domain/usecases/initiate_course_payment_usecase.dart';
import 'course_access_provider.dart';
import 'payment_providers.dart';

const _uuid = Uuid();

enum CheckoutStage { idle, working, freeEnrolled, initiated, error }

class CheckoutFlowState {
  const CheckoutFlowState({
    this.selectedProvider,
    this.idempotencyKey,
    this.initiation,
    this.isPolling = false,
    this.stage = CheckoutStage.idle,
    this.error,
    this.errorCode,
  });

  final PaymentProvider? selectedProvider;

  /// Stable key for the in-flight checkout attempt. Reused across retries of the
  /// same attempt and only cleared by [CheckoutFlowNotifier.reset].
  final String? idempotencyKey;
  final PaymentInitiationEntity? initiation;
  final bool isPolling;
  final CheckoutStage stage;
  final String? error;
  final String? errorCode;

  bool get isWorking => stage == CheckoutStage.working;

  CheckoutFlowState copyWith({
    PaymentProvider? selectedProvider,
    String? idempotencyKey,
    PaymentInitiationEntity? initiation,
    bool? isPolling,
    CheckoutStage? stage,
    String? error,
    String? errorCode,
    bool clearError = false,
    bool clearInitiation = false,
  }) {
    return CheckoutFlowState(
      selectedProvider: selectedProvider ?? this.selectedProvider,
      idempotencyKey: idempotencyKey ?? this.idempotencyKey,
      initiation: clearInitiation ? null : initiation ?? this.initiation,
      isPolling: isPolling ?? this.isPolling,
      stage: stage ?? this.stage,
      error: clearError ? null : error ?? this.error,
      errorCode: clearError ? null : errorCode ?? this.errorCode,
    );
  }
}

final checkoutFlowProvider = AsyncNotifierProviderFamily<CheckoutFlowNotifier,
    CheckoutFlowState, String>(CheckoutFlowNotifier.new);

class CheckoutFlowNotifier
    extends FamilyAsyncNotifier<CheckoutFlowState, String> {
  String get _courseId => arg;

  @override
  Future<CheckoutFlowState> build(String arg) async {
    return const CheckoutFlowState();
  }

  CheckoutFlowState get _current =>
      state.valueOrNull ?? const CheckoutFlowState();

  void selectProvider(PaymentProvider provider) {
    state = AsyncData(
      _current.copyWith(selectedProvider: provider, clearError: true),
    );
  }

  /// Free enrollment path. Returns true on success and refreshes access state.
  Future<bool> startFreeEnrollment() async {
    state = AsyncData(
      _current.copyWith(stage: CheckoutStage.working, clearError: true),
    );
    final result = await ref.read(enrollFreeUseCaseProvider)(_courseId);
    return result.fold(
      (failure) {
        state = AsyncData(
          _current.copyWith(
            stage: CheckoutStage.error,
            error: failure.message,
          ),
        );
        return false;
      },
      (_) {
        state = AsyncData(_current.copyWith(stage: CheckoutStage.freeEnrolled));
        _refreshAccess();
        return true;
      },
    );
  }

  /// Paid checkout path. Returns the backend initiation (with redirect URL) on
  /// success, or null on failure. The idempotency key is generated once and
  /// reused on retry so the same tap can never create two charges.
  Future<PaymentInitiationEntity?> startPaidCheckout({
    required PaymentProvider provider,
    String? couponCode,
  }) async {
    final key = _current.idempotencyKey ?? _uuid.v4();
    state = AsyncData(
      _current.copyWith(
        selectedProvider: provider,
        idempotencyKey: key,
        stage: CheckoutStage.working,
        clearError: true,
      ),
    );

    final result = await ref.read(initiateCoursePaymentUseCaseProvider)(
      InitiateCoursePaymentParams(
        courseId: _courseId,
        provider: provider,
        idempotencyKey: key,
        couponCode: couponCode,
      ),
    );

    return result.fold(
      (failure) {
        final code = failure is ServerFailure ? failure.code : null;
        state = AsyncData(
          _current.copyWith(
            stage: CheckoutStage.error,
            error: failure.message,
            errorCode: code,
          ),
        );
        return null;
      },
      (initiation) {
        state = AsyncData(
          _current.copyWith(
            initiation: initiation,
            stage: CheckoutStage.initiated,
          ),
        );
        return initiation;
      },
    );
  }

  /// Clears the checkout attempt, including the idempotency key, so the next
  /// initiate starts a brand-new payment.
  void reset() {
    state = const AsyncData(
      CheckoutFlowState(),
    );
  }

  void _refreshAccess() {
    ref.invalidate(courseAccessProvider(_courseId));
  }
}
