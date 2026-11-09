import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_sizes.dart';
import '../../../core/constants/app_text_styles.dart';
import '../app_button.dart';

/// The last thing a student sees before an attempt becomes final.
///
/// Submission cannot be undone, so this states plainly what is about to be
/// handed in. Blank questions are called out loudest and the confirm button
/// names the consequence, because the failure this exists to prevent is a
/// student submitting a 100-question paper with ten questions they never saw.
class SubmitSummarySheet extends StatelessWidget {
  const SubmitSummarySheet({
    super.key,
    required this.title,
    required this.total,
    required this.answered,
    required this.flagged,
    required this.isTimeUp,
    this.onReviewBlank,
  });

  final String title;
  final int total;
  final int answered;
  final int flagged;

  /// When the timer has run out there is nothing left to review, so the sheet
  /// drops the "go back" affordance and simply confirms.
  final bool isTimeUp;

  /// Jumps to the first unanswered question. Null when there are none.
  final VoidCallback? onReviewBlank;

  int get blank => total - answered;

  /// Returns true when the student confirmed submission.
  static Future<bool> show(
    BuildContext context, {
    required String title,
    required int total,
    required int answered,
    required int flagged,
    bool isTimeUp = false,
    VoidCallback? onReviewBlank,
  }) async {
    final result = await showModalBottomSheet<bool>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      isScrollControlled: true,
      showDragHandle: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSizes.radiusLg)),
      ),
      builder: (context) => SubmitSummarySheet(
        title: title,
        total: total,
        answered: answered,
        flagged: flagged,
        isTimeUp: isTimeUp,
        onReviewBlank: onReviewBlank == null
            ? null
            : () {
                Navigator.of(context).pop(false);
                onReviewBlank();
              },
      ),
    );
    return result ?? false;
  }

  @override
  Widget build(BuildContext context) {
    final hasBlank = blank > 0;

    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(
          AppSizes.screenPaddingH,
          0,
          AppSizes.screenPaddingH,
          AppSizes.lg,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(title, style: AppTextStyles.headlineSmall),
            const SizedBox(height: AppSizes.lg),
            Row(
              children: [
                Expanded(
                  child: _Stat(
                    value: '$answered',
                    label: 'answered',
                    color: AppColors.accentPrimary,
                  ),
                ),
                Expanded(
                  child: _Stat(
                    value: '$blank',
                    label: hasBlank ? 'still blank' : 'blank',
                    color: hasBlank ? AppColors.warning : AppColors.textMuted,
                    emphasise: hasBlank,
                  ),
                ),
                Expanded(
                  child: _Stat(
                    value: '$flagged',
                    label: 'flagged',
                    color: flagged > 0 ? AppColors.warning : AppColors.textMuted,
                  ),
                ),
              ],
            ),
            const SizedBox(height: AppSizes.lg),
            if (hasBlank && !isTimeUp) ...[
              Container(
                padding: const EdgeInsets.all(AppSizes.md),
                decoration: BoxDecoration(
                  color: AppColors.warningSurface,
                  borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  border: Border.all(color: AppColors.warning.withAlpha(90)),
                ),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Icon(
                      Icons.warning_amber_rounded,
                      size: AppSizes.iconSm,
                      color: AppColors.warning,
                    ),
                    const SizedBox(width: AppSizes.sm),
                    Expanded(
                      child: Text(
                        blank == 1
                            ? '1 question has no answer. It will be marked wrong.'
                            : '$blank questions have no answer. They will be marked wrong.',
                        style: AppTextStyles.bodySmall.copyWith(
                          color: AppColors.warning,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: AppSizes.md),
            ],
            if (isTimeUp) ...[
              Text(
                'Time is up. Your answers so far are what will be graded.',
                style: AppTextStyles.bodySmall.copyWith(color: AppColors.warning),
              ),
              const SizedBox(height: AppSizes.md),
            ] else
              Text(
                'Submitting is final — you cannot change your answers afterwards.',
                style: AppTextStyles.bodySmall.copyWith(color: AppColors.textMuted),
              ),
            const SizedBox(height: AppSizes.lg),
            if (hasBlank && onReviewBlank != null && !isTimeUp) ...[
              SizedBox(
                width: double.infinity,
                child: AppButton(
                  label: 'Go to first unanswered',
                  onPressed: onReviewBlank,
                  variant: AppButtonVariant.secondary,
                ),
              ),
              const SizedBox(height: AppSizes.sm),
            ],
            Row(
              children: [
                Expanded(
                  child: AppButton(
                    label: 'Keep working',
                    onPressed: () => Navigator.of(context).pop(false),
                    variant: AppButtonVariant.secondary,
                  ),
                ),
                const SizedBox(width: AppSizes.md),
                Expanded(
                  child: AppButton(
                    label: 'Submit',
                    onPressed: () => Navigator.of(context).pop(true),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _Stat extends StatelessWidget {
  const _Stat({
    required this.value,
    required this.label,
    required this.color,
    this.emphasise = false,
  });

  final String value;
  final String label;
  final Color color;
  final bool emphasise;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          value,
          style: AppTextStyles.displayMedium.copyWith(
            color: color,
            fontWeight: emphasise ? FontWeight.w800 : FontWeight.w700,
          ),
        ),
        Text(
          label,
          style: AppTextStyles.caption.copyWith(color: AppColors.textMuted),
        ),
      ],
    );
  }
}
