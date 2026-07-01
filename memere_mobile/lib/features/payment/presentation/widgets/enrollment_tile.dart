import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../domain/entities/enrollment_entity.dart';
import '../providers/purchase_history_provider.dart';

class EnrollmentTile extends ConsumerWidget {
  const EnrollmentTile({
    super.key,
    required this.enrollment,
    this.courseLabel,
    this.onTap,
  });

  final EnrollmentEntity enrollment;
  final String? courseLabel;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final label = courseLabel ??
        ref.watch(courseTitleProvider(enrollment.courseId)) ??
        'Course';
    final dateFormat = DateFormat('MMM d, yyyy');

    return AppSurface(
      onTap: onTap,
      padding: const EdgeInsets.all(AppSizes.md),
      shadows: AppShadows.sm,
      child: Row(
        children: [
          AppIconTile(
            icon: Icons.play_lesson_outlined,
            color: _sourceColor(enrollment.source),
            size: 44,
            iconSize: AppSizes.iconSm,
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
                Text(
                  [
                    if (enrollment.enrolledAt != null)
                      'Enrolled ${dateFormat.format(enrollment.enrolledAt!)}',
                    if (enrollment.expiresAt != null)
                      'Expires ${dateFormat.format(enrollment.expiresAt!)}',
                  ].join(' • '),
                  style: AppTextStyles.bodySmall,
                ),
              ],
            ),
          ),
          const SizedBox(width: AppSizes.sm),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              _SourceBadge(source: enrollment.source),
              if (onTap != null) ...[
                const SizedBox(height: AppSizes.sm),
                const Icon(
                  Icons.chevron_right_rounded,
                  color: AppColors.textSecondary,
                ),
              ],
            ],
          ),
        ],
      ),
    );
  }
}

class _SourceBadge extends StatelessWidget {
  const _SourceBadge({required this.source});

  final EnrollmentSource source;

  @override
  Widget build(BuildContext context) {
    final color = _sourceColor(source);
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
        _sourceLabel(source),
        style: AppTextStyles.labelSmall.copyWith(color: color),
      ),
    );
  }
}

Color _sourceColor(EnrollmentSource source) {
  switch (source) {
    case EnrollmentSource.purchase:
      return AppColors.accentPrimary;
    case EnrollmentSource.free:
      return AppColors.success;
    case EnrollmentSource.subscription:
      return AppColors.info;
    case EnrollmentSource.coupon:
      return AppColors.warning;
  }
}

String _sourceLabel(EnrollmentSource source) {
  switch (source) {
    case EnrollmentSource.purchase:
      return 'Purchase';
    case EnrollmentSource.free:
      return 'Free';
    case EnrollmentSource.subscription:
      return 'Subscription';
    case EnrollmentSource.coupon:
      return 'Coupon';
  }
}
