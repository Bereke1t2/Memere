import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_constants.dart';
import '../../domain/entities/payment_entity.dart';
import 'course_access_provider.dart';
import 'payment_providers.dart';
import 'purchase_history_provider.dart';
import 'subscription_provider.dart';

class PaymentStatusPollingParams {
  const PaymentStatusPollingParams({
    required this.paymentId,
    required this.courseId,
  });

  final String paymentId;

  /// Empty for subscription checkouts (no single course to refresh).
  final String courseId;

  @override
  bool operator ==(Object other) =>
      other is PaymentStatusPollingParams &&
      other.paymentId == paymentId &&
      other.courseId == courseId;

  @override
  int get hashCode => Object.hash(paymentId, courseId);
}

class PaymentPollingState {
  const PaymentPollingState({
    this.payment,
    this.isPolling = true,
    this.timedOut = false,
    this.error,
  });

  final PaymentEntity? payment;
  final bool isPolling;
  final bool timedOut;
  final String? error;

  PaymentStatus? get status => payment?.status;
  bool get isCompleted => payment?.isCompleted ?? false;
  bool get isFailed => payment?.isFailed ?? false;
  bool get isPending => payment == null || (payment?.isPending ?? false);

  PaymentPollingState copyWith({
    PaymentEntity? payment,
    bool? isPolling,
    bool? timedOut,
    String? error,
    bool clearError = false,
  }) {
    return PaymentPollingState(
      payment: payment ?? this.payment,
      isPolling: isPolling ?? this.isPolling,
      timedOut: timedOut ?? this.timedOut,
      error: clearError ? null : error ?? this.error,
    );
  }
}

final paymentStatusPollingProvider = AsyncNotifierProvider.autoDispose.family<
    PaymentStatusPollingNotifier,
    PaymentPollingState,
    PaymentStatusPollingParams>(PaymentStatusPollingNotifier.new);

/// Polls `GET /payments/:id/status` on a fixed interval until the payment
/// reaches a terminal state or the timeout elapses. On completion it refreshes
/// course access so the unlocked course is reflected everywhere.
class PaymentStatusPollingNotifier extends AutoDisposeFamilyAsyncNotifier<
    PaymentPollingState, PaymentStatusPollingParams> {
  Timer? _timer;
  DateTime? _startedAt;

  @override
  Future<PaymentPollingState> build(PaymentStatusPollingParams arg) async {
    ref.onDispose(() => _timer?.cancel());
    _startedAt = DateTime.now();
    return _pollOnce();
  }

  bool get _timedOut {
    final started = _startedAt;
    if (started == null) return false;
    final elapsed = DateTime.now().difference(started);
    return elapsed.inSeconds >= AppConstants.paymentPollTimeoutSeconds;
  }

  Future<PaymentPollingState> _pollOnce() async {
    final previous = state.valueOrNull;
    final result =
        await ref.read(getPaymentStatusUseCaseProvider)(arg.paymentId);

    return result.fold(
      (failure) {
        // Keep trying through transient errors until the timeout window closes.
        final timedOut = _timedOut;
        if (!timedOut) _scheduleNext();
        return (previous ?? const PaymentPollingState()).copyWith(
          isPolling: !timedOut,
          timedOut: timedOut,
          error: failure.message,
        );
      },
      (payment) {
        if (payment.isCompleted) _onCompleted();
        final terminal = payment.isTerminal;
        final timedOut = !terminal && _timedOut;
        if (!terminal && !timedOut) _scheduleNext();
        return PaymentPollingState(
          payment: payment,
          isPolling: !terminal && !timedOut,
          timedOut: timedOut,
        );
      },
    );
  }

  void _scheduleNext() {
    _timer?.cancel();
    _timer = Timer(
      const Duration(seconds: AppConstants.paymentPollIntervalSeconds),
      () async {
        final next = await _pollOnce();
        state = AsyncData(next);
      },
    );
  }

  /// Manual retry after a pending timeout: reopen the polling window and poll.
  Future<void> checkAgain() async {
    _timer?.cancel();
    _startedAt = DateTime.now();
    final current = state.valueOrNull ?? const PaymentPollingState();
    state = AsyncData(
      current.copyWith(isPolling: true, timedOut: false, clearError: true),
    );
    final next = await _pollOnce();
    state = AsyncData(next);
  }

  void _onCompleted() {
    // Refresh every surface that a completed payment can change.
    ref.invalidate(enrollmentListProvider);
    ref.invalidate(paymentHistoryProvider);
    ref.invalidate(mySubscriptionProvider);
    if (arg.courseId.isNotEmpty) {
      ref.invalidate(courseAccessProvider(arg.courseId));
    }
  }
}
