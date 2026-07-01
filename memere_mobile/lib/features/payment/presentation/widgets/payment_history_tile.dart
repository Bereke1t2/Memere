import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/payment_entity.dart';
import '../providers/purchase_history_provider.dart';

class PaymentHistoryTile extends ConsumerWidget {
  const PaymentHistoryTile({
    super.key,
    required this.payment,
    this.courseLabel,
    this.onTap,
  });

  final PaymentEntity payment;

  /// Explicit label override. When null, the title is resolved from the catalog.
  final String? courseLabel;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final hasCourse =
        payment.courseId != null && payment.courseId!.isNotEmpty;
    final resolvedTitle =
        hasCourse ? ref.watch(courseTitleProvider(payment.courseId!)) : null;
    final label = courseLabel ??
        resolvedTitle ??
        (hasCourse ? 'Course' : 'Subscription');

    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(AppSizes.radiusMd),
      child: Container(
        padding: const EdgeInsets.all(AppSizes.md),
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          border: Border.all(color: AppColors.border),
        ),
        child: Row(
          children: [
            Container(
              width: AppSizes.avatarSm,
              height: AppSizes.avatarSm,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: _statusColor(payment.status).withAlpha(28),
                borderRadius: BorderRadius.circular(AppSizes.radiusSm),
              ),
              child: Icon(
                _statusIcon(payment.status),
                size: AppSizes.iconSm,
                color: _statusColor(payment.status),
              ),
            ),
            const SizedBox(width: AppSizes.md),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    label,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.titleMedium,
                  ),
                  const SizedBox(height: AppSizes.xs),
                  Text(payment.amountLabel, style: AppTextStyles.bodySmall),
                ],
              ),
            ),
            const SizedBox(width: AppSizes.sm),
            _StatusBadge(status: payment.status),
          ],
        ),
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status});

  final PaymentStatus status;

  @override
  Widget build(BuildContext context) {
    final color = _statusColor(status);
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: color.withAlpha(28),
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
      ),
      child: Text(
        _statusLabel(status),
        style: AppTextStyles.labelSmall.copyWith(color: color),
      ),
    );
  }
}

Color _statusColor(PaymentStatus status) {
  switch (status) {
    case PaymentStatus.pending:
      return AppColors.warning;
    case PaymentStatus.completed:
      return AppColors.success;
    case PaymentStatus.failed:
      return AppColors.error;
    case PaymentStatus.refunded:
      return AppColors.info;
  }
}

IconData _statusIcon(PaymentStatus status) {
  switch (status) {
    case PaymentStatus.pending:
      return Icons.schedule_rounded;
    case PaymentStatus.completed:
      return Icons.check_circle_rounded;
    case PaymentStatus.failed:
      return Icons.error_rounded;
    case PaymentStatus.refunded:
      return Icons.replay_rounded;
  }
}

String _statusLabel(PaymentStatus status) {
  switch (status) {
    case PaymentStatus.pending:
      return 'Pending';
    case PaymentStatus.completed:
      return 'Completed';
    case PaymentStatus.failed:
      return 'Failed';
    case PaymentStatus.refunded:
      return 'Refunded';
  }
}
