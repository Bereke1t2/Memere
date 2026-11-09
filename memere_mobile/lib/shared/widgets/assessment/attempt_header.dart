import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_sizes.dart';
import '../../../core/constants/app_text_styles.dart';

/// The bar that stays put while a student scrolls a long question.
///
/// An imported exam prompt can run well past a screen, and the facts a student
/// checks most — which question they are on, how much is left, how long they
/// have — were previously at the top of the scroll view and went out of sight
/// exactly when the question got long enough to worry about.
///
/// Progress is measured in questions ANSWERED, not in position. Position tells
/// a student where they are; answered tells them what is left to do, which is
/// the thing that decides whether they can afford to dwell on this question.
class AttemptHeader extends StatelessWidget implements PreferredSizeWidget {
  const AttemptHeader({
    super.key,
    required this.currentIndex,
    required this.total,
    required this.answered,
    required this.timer,
    required this.onOpenPalette,
    required this.onDecreaseText,
    required this.onIncreaseText,
    this.textScaleLabel,
  });

  final int currentIndex;
  final int total;
  final int answered;

  /// The timer badge, supplied by the caller so quiz and exam keep their own.
  final Widget timer;

  final VoidCallback onOpenPalette;

  /// Null disables the control at that end of the range.
  final VoidCallback? onDecreaseText;
  final VoidCallback? onIncreaseText;
  final String? textScaleLabel;

  static const double _height = 104;

  @override
  Size get preferredSize => const Size.fromHeight(_height);

  @override
  Widget build(BuildContext context) {
    final progress = total == 0 ? 0.0 : answered / total;

    return Material(
      color: AppColors.bgPrimary,
      child: SafeArea(
        bottom: false,
        child: SizedBox(
          height: _height,
          child: Padding(
            padding: const EdgeInsets.fromLTRB(
              AppSizes.screenPaddingH,
              AppSizes.sm,
              AppSizes.screenPaddingH,
              AppSizes.sm,
            ),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            'Question ${currentIndex + 1} of $total',
                            style: AppTextStyles.titleMedium.copyWith(
                              color: AppColors.textPrimary,
                              fontWeight: FontWeight.w700,
                            ),
                          ),
                          Text(
                            '$answered answered · ${total - answered} left',
                            style: AppTextStyles.caption.copyWith(
                              color: AppColors.textMuted,
                            ),
                          ),
                        ],
                      ),
                    ),
                    _TextSizeControl(
                      onDecrease: onDecreaseText,
                      onIncrease: onIncreaseText,
                      label: textScaleLabel,
                    ),
                    const SizedBox(width: AppSizes.sm),
                    timer,
                  ],
                ),
                const SizedBox(height: AppSizes.sm),
                Row(
                  children: [
                    Expanded(
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(4),
                        child: LinearProgressIndicator(
                          minHeight: 6,
                          value: progress,
                          backgroundColor: AppColors.bgTertiary,
                          valueColor: const AlwaysStoppedAnimation(
                            AppColors.accentPrimary,
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(width: AppSizes.sm),
                    _PaletteButton(onTap: onOpenPalette),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _TextSizeControl extends StatelessWidget {
  const _TextSizeControl({
    required this.onDecrease,
    required this.onIncrease,
    this.label,
  });

  final VoidCallback? onDecrease;
  final VoidCallback? onIncrease;
  final String? label;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          _ScaleButton(
            icon: Icons.text_decrease_rounded,
            tooltip: 'Smaller text',
            onTap: onDecrease,
          ),
          _ScaleButton(
            icon: Icons.text_increase_rounded,
            tooltip: label == null ? 'Larger text' : 'Larger text ($label)',
            onTap: onIncrease,
          ),
        ],
      ),
    );
  }
}

class _ScaleButton extends StatelessWidget {
  const _ScaleButton({
    required this.icon,
    required this.tooltip,
    required this.onTap,
  });

  final IconData icon;
  final String tooltip;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final enabled = onTap != null;
    return Tooltip(
      message: tooltip,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        child: SizedBox(
          width: 38,
          height: 34,
          child: Icon(
            icon,
            size: AppSizes.iconSm,
            color: enabled ? AppColors.textSecondary : AppColors.textDisabled,
          ),
        ),
      ),
    );
  }
}

class _PaletteButton extends StatelessWidget {
  const _PaletteButton({required this.onTap});

  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: 'All questions',
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(AppSizes.radiusSm),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
          decoration: BoxDecoration(
            color: AppColors.bgTertiary,
            borderRadius: BorderRadius.circular(AppSizes.radiusSm),
            border: Border.all(color: AppColors.borderStrong),
          ),
          child: const Icon(
            Icons.grid_view_rounded,
            size: AppSizes.iconSm,
            color: AppColors.textSecondary,
          ),
        ),
      ),
    );
  }
}
