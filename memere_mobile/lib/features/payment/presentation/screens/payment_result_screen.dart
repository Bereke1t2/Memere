import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../providers/course_access_provider.dart';
import '../providers/payment_status_polling_provider.dart';
import '../providers/subscription_provider.dart';

/// Shows the backend-verified outcome of a checkout. Course access is only
/// promised once polling reports `completed` and the access provider confirms it.
class PaymentResultScreen extends ConsumerWidget {
  const PaymentResultScreen({
    super.key,
    required this.paymentId,
    required this.courseId,
  });

  final String paymentId;
  final String courseId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final params = PaymentStatusPollingParams(
      paymentId: paymentId,
      courseId: courseId,
    );
    final pollingAsync = ref.watch(paymentStatusPollingProvider(params));
    final state = pollingAsync.valueOrNull;

    final view = _resolveView(state);

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(title: const Text('Payment')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(AppSizes.lg),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Container(
                width: 88,
                height: 88,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: view.color.withAlpha(28),
                  shape: BoxShape.circle,
                ),
                child: view.showSpinner
                    ? const SizedBox(
                        width: 32,
                        height: 32,
                        child: CircularProgressIndicator(
                          color: AppColors.accentPrimary,
                          strokeWidth: 3,
                        ),
                      )
                    : Icon(view.icon, color: view.color, size: AppSizes.iconXl),
              ),
              const SizedBox(height: AppSizes.lg),
              Text(
                view.title,
                textAlign: TextAlign.center,
                style: AppTextStyles.headlineLarge,
              ),
              const SizedBox(height: AppSizes.sm),
              Text(
                view.body,
                textAlign: TextAlign.center,
                style: AppTextStyles.bodyMedium
                    .copyWith(color: AppColors.textSecondary),
              ),
              const SizedBox(height: AppSizes.xl),
              ..._buildActions(context, ref, state, view),
            ],
          ),
        ),
      ),
    );
  }

  List<Widget> _buildActions(
    BuildContext context,
    WidgetRef ref,
    PaymentPollingState? state,
    _ResultView view,
  ) {
    final widgets = <Widget>[];
    final isSubscription = courseId.isEmpty;

    if (view.kind == _ResultKind.completed) {
      widgets.add(
        AppButton(
          label: isSubscription ? 'Start learning' : 'Go to course',
          // Only leave after access state has refreshed to confirm the unlock.
          onPressed: () async {
            if (isSubscription) {
              await ref.read(mySubscriptionProvider.future);
            } else {
              await ref.read(courseAccessProvider(courseId).future);
            }
            if (context.mounted) _leave(context);
          },
        ),
      );
      widgets.add(const SizedBox(height: AppSizes.sm));
      widgets.add(
        AppButton(
          label: 'View purchases',
          variant: AppButtonVariant.outline,
          onPressed: () => _goPurchases(context),
        ),
      );
    } else if (view.kind == _ResultKind.failed) {
      widgets.add(
        AppButton(
          label: 'Back',
          onPressed: () => _leave(context),
        ),
      );
    } else {
      // pending / timed out
      final params = PaymentStatusPollingParams(
        paymentId: paymentId,
        courseId: courseId,
      );
      widgets.add(
        AppButton(
          label: 'Check again',
          onPressed: () => ref
              .read(paymentStatusPollingProvider(params).notifier)
              .checkAgain(),
        ),
      );
      widgets.add(const SizedBox(height: AppSizes.sm));
      widgets.add(
        AppButton(
          label: 'View purchases',
          variant: AppButtonVariant.outline,
          onPressed: () => _goPurchases(context),
        ),
      );
    }

    return widgets;
  }

  /// Leaves the result screen back to whatever launched the checkout (course
  /// detail or the plans screen). Falls back to home for deep links.
  void _leave(BuildContext context) {
    if (context.canPop()) {
      context.pop();
    } else {
      context.go(AppRoutes.home);
    }
  }

  void _goPurchases(BuildContext context) {
    // Replace the transient result screen so back returns to the launcher.
    context.pushReplacement(AppRoutes.purchaseHistory);
  }

  _ResultView _resolveView(PaymentPollingState? state) {
    if (state == null) {
      return const _ResultView(
        kind: _ResultKind.pending,
        title: 'Confirming payment',
        body: 'Please wait while we verify your payment.',
        icon: Icons.hourglass_top_rounded,
        color: AppColors.warning,
        showSpinner: true,
      );
    }

    if (state.isCompleted) {
      final isSubscription = courseId.isEmpty;
      return _ResultView(
        kind: _ResultKind.completed,
        title: isSubscription ? 'Subscription active' : 'Course unlocked',
        body: isSubscription
            ? 'You now have all-access to every course.'
            : 'Your payment was successful. Enjoy the course!',
        icon: Icons.check_circle_rounded,
        color: AppColors.success,
      );
    }

    if (state.isFailed) {
      return const _ResultView(
        kind: _ResultKind.failed,
        title: 'Payment failed',
        body: 'Your payment did not go through. You can try again.',
        icon: Icons.error_rounded,
        color: AppColors.error,
      );
    }

    if (state.timedOut) {
      return const _ResultView(
        kind: _ResultKind.pending,
        title: 'Still pending',
        body:
            'We could not confirm your payment yet. It may take a moment — tap '
            'check again.',
        icon: Icons.schedule_rounded,
        color: AppColors.warning,
      );
    }

    return const _ResultView(
      kind: _ResultKind.pending,
      title: 'Confirming payment',
      body: 'Please wait while we verify your payment.',
      icon: Icons.hourglass_top_rounded,
      color: AppColors.warning,
      showSpinner: true,
    );
  }
}

enum _ResultKind { completed, failed, pending }

class _ResultView {
  const _ResultView({
    required this.kind,
    required this.title,
    required this.body,
    required this.icon,
    required this.color,
    this.showSpinner = false,
  });

  final _ResultKind kind;
  final String title;
  final String body;
  final IconData icon;
  final Color color;
  final bool showSpinner;
}
