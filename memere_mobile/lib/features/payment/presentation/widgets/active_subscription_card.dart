import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/subscription_entity.dart';
import '../../domain/entities/subscription_plan_entity.dart';

class ActiveSubscriptionCard extends StatelessWidget {
  const ActiveSubscriptionCard({
    super.key,
    required this.subscription,
    required this.onCancel,
    this.isCanceling = false,
  });

  final SubscriptionEntity subscription;
  final VoidCallback onCancel;
  final bool isCanceling;

  @override
  Widget build(BuildContext context) {
    final dateFormat = DateFormat('MMM d, yyyy');
    final color = _statusColor(subscription.status);

    return Container(
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(
                child: Text(
                  '${subscription.plan.label} plan',
                  style: AppTextStyles.headlineSmall,
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSizes.sm,
                  vertical: AppSizes.xs,
                ),
                decoration: BoxDecoration(
                  color: color.withAlpha(28),
                  borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                ),
                child: Text(
                  _statusLabel(subscription.status),
                  style: AppTextStyles.labelSmall.copyWith(color: color),
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSizes.md),
          if (subscription.currentPeriodStart != null)
            _Row(
              label: 'Started',
              value: dateFormat.format(subscription.currentPeriodStart!),
            ),
          if (subscription.currentPeriodEnd != null)
            _Row(
              label: subscription.isCanceled ? 'Access until' : 'Renews',
              value: dateFormat.format(subscription.currentPeriodEnd!),
            ),
          if (subscription.canceledAt != null)
            _Row(
              label: 'Canceled',
              value: dateFormat.format(subscription.canceledAt!),
            ),
          if (!subscription.isCanceled &&
              subscription.status == SubscriptionStatus.active) ...[
            const SizedBox(height: AppSizes.md),
            AppButton(
              label: 'Cancel subscription',
              variant: AppButtonVariant.danger,
              isLoading: isCanceling,
              onPressed: isCanceling ? null : onCancel,
            ),
          ],
        ],
      ),
    );
  }
}

class _Row extends StatelessWidget {
  const _Row({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSizes.xs),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: AppTextStyles.bodySmall),
          Text(value, style: AppTextStyles.bodyMedium),
        ],
      ),
    );
  }
}

Color _statusColor(SubscriptionStatus status) {
  switch (status) {
    case SubscriptionStatus.active:
      return AppColors.success;
    case SubscriptionStatus.pastDue:
      return AppColors.warning;
    case SubscriptionStatus.canceled:
      return AppColors.info;
    case SubscriptionStatus.expired:
      return AppColors.textSecondary;
  }
}

String _statusLabel(SubscriptionStatus status) {
  switch (status) {
    case SubscriptionStatus.active:
      return 'Active';
    case SubscriptionStatus.pastDue:
      return 'Past due';
    case SubscriptionStatus.canceled:
      return 'Canceled';
    case SubscriptionStatus.expired:
      return 'Expired';
  }
}
