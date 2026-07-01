import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../domain/entities/mock_exam_entity.dart';
import '../providers/exam_attempt_provider.dart';
import '../providers/exam_providers.dart';

class MockExamCard extends ConsumerStatefulWidget {
  const MockExamCard({
    super.key,
    required this.exam,
  });

  final MockExamEntity exam;

  @override
  ConsumerState<MockExamCard> createState() => _MockExamCardState();
}

class _MockExamCardState extends ConsumerState<MockExamCard> {
  bool _isStarting = false;

  @override
  Widget build(BuildContext context) {
    final exam = widget.exam;
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.md),
      gradient: AppColors.cardGradient,
      shadows: AppShadows.md,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const AppIconTile(
                icon: Icons.assignment_outlined,
                color: AppColors.accentTertiary,
                size: 48,
              ),
              const SizedBox(width: AppSizes.md),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      exam.title,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: AppTextStyles.titleMedium,
                    ),
                    const SizedBox(height: AppSizes.xs),
                    Text(
                      '${exam.subject} · Grade ${exam.grade}',
                      style: AppTextStyles.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSizes.md),
          const Divider(color: AppColors.divider),
          const SizedBox(height: AppSizes.md),
          Wrap(
            spacing: AppSizes.sm,
            runSpacing: AppSizes.sm,
            children: [
              _MetaChip(
                icon: Icons.timer_outlined,
                label: exam.durationLabel,
                color: AppColors.accentTertiary,
              ),
              _MetaChip(
                icon: Icons.star_outline_rounded,
                label: '${exam.totalMarks} marks',
                color: AppColors.accentSecondary,
              ),
              _MetaChip(
                icon: Icons.flag_outlined,
                label: 'Pass ${exam.passMarks}',
                color: AppColors.success,
              ),
            ],
          ),
          if (exam.instructions != null) ...[
            const SizedBox(height: AppSizes.md),
            Text(
              exam.instructions!,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: AppTextStyles.bodySmall.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ],
          const SizedBox(height: AppSizes.lg),
          AppButton(
            label: 'Start exam',
            isLoading: _isStarting,
            onPressed: () => _confirmStart(context),
            suffixIcon: Icons.arrow_forward_rounded,
          ),
        ],
      ),
    );
  }

  Future<void> _confirmStart(BuildContext context) async {
    final exam = widget.exam;
    final start = await showModalBottomSheet<bool>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      isScrollControlled: true,
      builder: (sheetContext) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(AppSizes.screenPaddingH),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(exam.title, style: AppTextStyles.headlineSmall),
              const SizedBox(height: AppSizes.md),
              _ConfirmRow(label: 'Duration', value: exam.durationLabel),
              _ConfirmRow(
                label: 'Total marks',
                value: exam.totalMarks.toString(),
              ),
              _ConfirmRow(
                label: 'Pass marks',
                value: exam.passMarks.toString(),
              ),
              if (exam.instructions != null) ...[
                const SizedBox(height: AppSizes.sm),
                Text(
                  exam.instructions!,
                  style: AppTextStyles.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
              const SizedBox(height: AppSizes.md),
              Container(
                padding: const EdgeInsets.all(AppSizes.md),
                decoration: BoxDecoration(
                  color: AppColors.warningSurface,
                  borderRadius: BorderRadius.circular(AppSizes.radiusSm),
                ),
                child: Row(
                  children: [
                    const Icon(
                      Icons.warning_amber_rounded,
                      color: AppColors.warning,
                      size: AppSizes.iconSm,
                    ),
                    const SizedBox(width: AppSizes.sm),
                    Expanded(
                      child: Text(
                        'The timer starts when you begin. Keep the app open '
                        'and submit before time runs out.',
                        style: AppTextStyles.bodySmall.copyWith(
                          color: AppColors.warning,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: AppSizes.lg),
              Row(
                children: [
                  Expanded(
                    child: AppButton(
                      label: 'Cancel',
                      onPressed: () => Navigator.of(sheetContext).pop(false),
                      variant: AppButtonVariant.secondary,
                    ),
                  ),
                  const SizedBox(width: AppSizes.md),
                  Expanded(
                    child: AppButton(
                      label: 'Start exam',
                      onPressed: () => Navigator.of(sheetContext).pop(true),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );

    if (start != true || !context.mounted) return;
    await _startExam(context);
  }

  Future<void> _startExam(BuildContext context) async {
    setState(() => _isStarting = true);
    final result = await ref.read(startExamUseCaseProvider)(widget.exam.id);
    if (!mounted) return;
    setState(() => _isStarting = false);

    result.fold(
      (failure) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(failure.message)),
        );
      },
      (attempt) {
        // Hand the started attempt to the attempt screen so it does not start
        // a second attempt on the server.
        ref.read(pendingExamAttemptProvider.notifier).state = attempt;
        context.go(
          AppRoutes.examAttemptPath(
            attemptId: attempt.attemptId,
            examId: attempt.examId,
          ),
        );
      },
    );
  }
}

class _MetaChip extends StatelessWidget {
  const _MetaChip({
    required this.icon,
    required this.label,
    required this.color,
  });

  final IconData icon;
  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: color.withAlpha(24),
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: color.withAlpha(76)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: AppSizes.iconXs, color: color),
          const SizedBox(width: AppSizes.xs),
          Text(label, style: AppTextStyles.labelSmall),
        ],
      ),
    );
  }
}

class _ConfirmRow extends StatelessWidget {
  const _ConfirmRow({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSizes.sm),
      child: Row(
        children: [
          Expanded(
            child: Text(label, style: AppTextStyles.bodySmall),
          ),
          Text(value, style: AppTextStyles.labelMedium),
        ],
      ),
    );
  }
}
