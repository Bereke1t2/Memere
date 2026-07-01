import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/subscription_plan_entity.dart';
import '../providers/payment_providers.dart';
import '../providers/subscription_provider.dart';
import '../widgets/active_subscription_card.dart';
import '../widgets/payment_empty_state.dart';
import '../widgets/payment_provider_sheet.dart';
import '../widgets/subscription_plan_card.dart';

class SubscriptionPlansScreen extends ConsumerStatefulWidget {
  const SubscriptionPlansScreen({super.key});

  @override
  ConsumerState<SubscriptionPlansScreen> createState() =>
      _SubscriptionPlansScreenState();
}

class _SubscriptionPlansScreenState
    extends ConsumerState<SubscriptionPlansScreen> {
  SubscriptionPlanType? _selectedPlan;
  bool _isCanceling = false;

  @override
  Widget build(BuildContext context) {
    final plansAsync = ref.watch(subscriptionPlansProvider);
    final subscriptionAsync = ref.watch(mySubscriptionProvider);
    final checkout = ref.watch(subscriptionCheckoutProvider).valueOrNull;
    final activeSubscription = subscriptionAsync.valueOrNull;
    final hasActive = activeSubscription?.grantsAccess ?? false;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(title: const Text('All-access plans')),
      body: SafeArea(
        top: false,
        child: plansAsync.when(
          loading: () => const Center(
            child: CircularProgressIndicator(color: AppColors.accentPrimary),
          ),
          error: (error, _) => PaymentEmptyState(
            icon: Icons.error_outline_rounded,
            title: 'Could not load plans',
            body: error is Failure ? error.message : 'Please try again.',
            buttonLabel: 'Retry',
            onPressed: () => ref.invalidate(subscriptionPlansProvider),
          ),
          data: (plans) {
            return RefreshIndicator(
              color: AppColors.accentPrimary,
              backgroundColor: AppColors.bgSecondary,
              onRefresh: () async {
                ref.invalidate(subscriptionPlansProvider);
                ref.invalidate(mySubscriptionProvider);
                await ref.read(subscriptionPlansProvider.future);
              },
              child: ListView(
                physics: const AlwaysScrollableScrollPhysics(),
                padding: const EdgeInsets.all(AppSizes.md),
                children: [
                  if (hasActive && activeSubscription != null) ...[
                    const Text('Your subscription',
                        style: AppTextStyles.headlineSmall),
                    const SizedBox(height: AppSizes.md),
                    ActiveSubscriptionCard(
                      subscription: activeSubscription,
                      isCanceling: _isCanceling,
                      onCancel: () => _confirmCancel(activeSubscription.id),
                    ),
                    const SizedBox(height: AppSizes.lg),
                  ],
                  Text(
                    hasActive ? 'Other plans' : 'Unlock every course',
                    style: AppTextStyles.headlineSmall,
                  ),
                  const SizedBox(height: AppSizes.sm),
                  const Text(
                    'One subscription gives you access to all paid courses.',
                    style: AppTextStyles.bodySmall,
                  ),
                  const SizedBox(height: AppSizes.md),
                  if (plans.isEmpty)
                    const PaymentEmptyState(
                      icon: Icons.workspace_premium_outlined,
                      title: 'No plans available',
                      body: 'Subscription plans will appear here soon.',
                    )
                  else
                    ...plans.map(
                      (plan) => Padding(
                        padding: const EdgeInsets.only(bottom: AppSizes.sm),
                        child: SubscriptionPlanCard(
                          plan: plan,
                          selected: _selectedPlan == plan.plan,
                          onTap: () =>
                              setState(() => _selectedPlan = plan.plan),
                        ),
                      ),
                    ),
                  if (!hasActive && plans.isNotEmpty) ...[
                    const SizedBox(height: AppSizes.md),
                    AppButton(
                      label: 'Subscribe',
                      isLoading: checkout?.isWorking ?? false,
                      isDisabled: _selectedPlan == null,
                      onPressed: _selectedPlan == null
                          ? null
                          : () => _subscribe(_selectedPlan!),
                    ),
                  ],
                ],
              ),
            );
          },
        ),
      ),
    );
  }

  Future<void> _subscribe(SubscriptionPlanType plan) async {
    final provider = await PaymentProviderSheet.show(context);
    if (provider == null || !mounted) return;

    final initiation = await ref
        .read(subscriptionCheckoutProvider.notifier)
        .startCheckout(plan: plan, provider: provider);
    if (!mounted) return;

    if (initiation == null) {
      final error =
          ref.read(subscriptionCheckoutProvider).valueOrNull?.error;
      _showMessage(error ?? 'Could not start subscription. Please try again.');
      return;
    }

    // Subscription checkout reuses the course payment WebView (no course id).
    context.push(
      AppRoutes.paymentWebViewPath(
        paymentId: initiation.paymentId,
        courseId: '',
        redirectUrl: initiation.redirectUrl,
      ),
    );
  }

  Future<void> _confirmCancel(String subscriptionId) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        backgroundColor: AppColors.bgSecondary,
        title: const Text('Cancel subscription?',
            style: AppTextStyles.headlineSmall),
        content: Text(
          'You will keep access until the end of the current billing period.',
          style: AppTextStyles.bodyMedium
              .copyWith(color: AppColors.textSecondary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text('Keep'),
          ),
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: const Text('Cancel plan'),
          ),
        ],
      ),
    );
    if (confirmed != true || !mounted) return;

    setState(() => _isCanceling = true);
    final result =
        await ref.read(cancelSubscriptionUseCaseProvider)(subscriptionId);
    if (!mounted) return;
    setState(() => _isCanceling = false);

    result.fold(
      (failure) => _showMessage(failure.message),
      (_) {
        // Keep access until current_period_end; just refresh the server view.
        ref.invalidate(mySubscriptionProvider);
        _showMessage('Subscription canceled. Access continues until period end.');
      },
    );
  }

  void _showMessage(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }
}
