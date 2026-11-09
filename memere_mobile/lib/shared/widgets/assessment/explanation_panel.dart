import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_sizes.dart';
import '../../../core/constants/app_text_styles.dart';
import '../math_text.dart';

/// The worked solution shown under a reviewed question.
///
/// Imported explanations are not one-liners: a maths derivation runs to several
/// numbered steps and a displayed equation or two. Printing all of them in full
/// turns a 100-question review into a scroll a student gives up on, so a long
/// one is collapsed behind a fade until it is asked for. Short explanations —
/// the great majority — are never given a toggle they do not need.
///
/// Length is judged on the source text rather than by measuring the rendered
/// height: the content mixes prose with maths widgets, so a measure-then-relayout
/// pass would cost more than it is worth for a threshold this rough.
class ExplanationPanel extends StatefulWidget {
  const ExplanationPanel({
    super.key,
    required this.explanation,
    this.title = 'Explanation',
  });

  final String explanation;
  final String title;

  /// Roughly four lines on a phone. Above this, collapsing earns its keep.
  static const int _collapseThreshold = 220;

  /// Height a collapsed explanation is clipped to.
  static const double _collapsedHeight = 96;

  @override
  State<ExplanationPanel> createState() => _ExplanationPanelState();
}

class _ExplanationPanelState extends State<ExplanationPanel> {
  bool _expanded = false;

  @override
  Widget build(BuildContext context) {
    final body = MathText(
      widget.explanation,
      style: AppTextStyles.bodySmall.copyWith(
        color: AppColors.textSecondary,
        height: 1.45,
      ),
      // Formulas stay brighter than the prose so a derivation reads clearly on
      // the tinted panel.
      mathColor: AppColors.textPrimary,
      selectable: true,
    );

    final isLong = widget.explanation.length > ExplanationPanel._collapseThreshold;
    final collapsed = isLong && !_expanded;

    return Container(
      padding: const EdgeInsets.all(AppSizes.md),
      decoration: BoxDecoration(
        color: AppColors.infoSurface,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(color: AppColors.info.withAlpha(68)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(
                Icons.lightbulb_outline_rounded,
                size: AppSizes.iconSm,
                color: AppColors.info,
              ),
              const SizedBox(width: AppSizes.sm),
              Text(
                widget.title,
                style: AppTextStyles.labelSmall.copyWith(
                  color: AppColors.info,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSizes.sm),
          if (collapsed)
            // Fading the clipped edge says "there is more" without needing a
            // label to explain the cut.
            ShaderMask(
              shaderCallback: (rect) => const LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                colors: [Colors.white, Colors.white, Colors.transparent],
                stops: [0.0, 0.62, 1.0],
              ).createShader(rect),
              blendMode: BlendMode.dstIn,
              child: SizedBox(
                height: ExplanationPanel._collapsedHeight,
                child: SingleChildScrollView(
                  physics: const NeverScrollableScrollPhysics(),
                  child: body,
                ),
              ),
            )
          else
            body,
          if (isLong) ...[
            const SizedBox(height: AppSizes.xs),
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton(
                onPressed: () => setState(() => _expanded = !_expanded),
                style: TextButton.styleFrom(
                  foregroundColor: AppColors.info,
                  padding: const EdgeInsets.symmetric(horizontal: 4),
                  minimumSize: const Size(0, 32),
                  tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                ),
                child: Text(
                  _expanded ? 'Show less' : 'Show full solution',
                  style: AppTextStyles.labelSmall.copyWith(
                    color: AppColors.info,
                    fontWeight: FontWeight.w700,
                  ),
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }
}
