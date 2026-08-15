import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/exam_subject_score_entity.dart';
import '../providers/exam_attempt_provider.dart';
import '../providers/exam_providers.dart';
import '../providers/exam_result_provider.dart';
import '../widgets/exam_question_feedback_tile.dart';
import '../widgets/exam_result_skeleton.dart';
import '../widgets/exam_score_summary.dart';
import '../widgets/exam_subject_breakdown.dart';

class ExamResultScreen extends ConsumerStatefulWidget {
  const ExamResultScreen({
    super.key,
    required this.attemptId,
  });

  final String attemptId;

  @override
  ConsumerState<ExamResultScreen> createState() => _ExamResultScreenState();
}

class _ExamResultScreenState extends ConsumerState<ExamResultScreen> {
  int _filterIndex = 0; // 0: All, 1: Correct, 2: Missed
  bool _isRetaking = false;

  Future<void> _handleRetake(String examId) async {
    if (_isRetaking) return;
    setState(() => _isRetaking = true);

    final messenger = ScaffoldMessenger.of(context);
    try {
      final either = await ref.read(startExamUseCaseProvider)(examId);
      if (!mounted) return;
      either.fold(
        (failure) {
          setState(() => _isRetaking = false);
          messenger.showSnackBar(SnackBar(content: Text(failure.message)));
        },
        (attempt) {
          ref.read(pendingExamAttemptProvider.notifier).state = attempt;
          ref.invalidate(myAllExamAttemptsProvider);
          ref.invalidate(examAttemptsByExamProvider(examId));
          context.push(
            AppRoutes.examAttemptPath(
              attemptId: attempt.attemptId,
              examId: examId,
            ),
          );
        },
      );
    } catch (e) {
      if (!mounted) return;
      setState(() => _isRetaking = false);
      messenger.showSnackBar(SnackBar(content: Text('Could not retake exam: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    final resultAsync = ref.watch(examResultProvider(widget.attemptId));

    return Scaffold(
      backgroundColor: const Color(0xFF070709),
      appBar: AppBar(
        backgroundColor: const Color(0xFF0C0C10),
        elevation: 0,
        scrolledUnderElevation: 0,
        title: const Text(
          'Exam Results & Review',
          style: TextStyle(fontSize: 17, fontWeight: FontWeight.w700, color: Colors.white),
        ),
        leading: IconButton(
          tooltip: 'Back to Mock Exams',
          onPressed: () {
            if (context.canPop()) {
              context.pop();
            } else {
              context.go(AppRoutes.mockExams);
            }
          },
          icon: const Icon(Icons.arrow_back_rounded, color: Color(0xFFD4D4D8)),
        ),
        actions: [
          resultAsync.maybeWhen(
            data: (result) => TextButton.icon(
              onPressed: _isRetaking ? null : () => _handleRetake(result.examId),
              icon: _isRetaking
                  ? const SizedBox(
                      width: 14,
                      height: 14,
                      child: CircularProgressIndicator(strokeWidth: 2, color: Color(0xFF10B981)),
                    )
                  : const Icon(Icons.replay_rounded, size: 16, color: Color(0xFF10B981)),
              label: const Text(
                'Retake',
                style: TextStyle(color: Color(0xFF10B981), fontWeight: FontWeight.w700),
              ),
            ),
            orElse: () => const SizedBox.shrink(),
          ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: resultAsync.when(
          loading: () => const ExamResultSkeleton(),
          error: (error, _) => _ResultErrorState(
            message: error is Failure
                ? error.message
                : 'Could not load exam result. Please try again.',
            onRetry: () => ref.invalidate(examResultProvider(widget.attemptId)),
          ),
          data: (result) {
            final breakdown = result.subjectBreakdown.entries
                .map(
                  (entry) => ExamSubjectScoreEntity(
                    key: entry.key,
                    earned: entry.value.earned,
                    possible: entry.value.possible,
                  ),
                )
                .toList();

            final correctCount = result.feedback.where((f) => f.correct).length;
            final missedCount = result.feedback.length - correctCount;

            final filteredFeedback = result.feedback.asMap().entries.where((entry) {
              if (_filterIndex == 1) return entry.value.correct;
              if (_filterIndex == 2) return !entry.value.correct;
              return true;
            }).toList();

            return ListView(
              physics: const BouncingScrollPhysics(),
              padding: const EdgeInsets.fromLTRB(18, 16, 18, 40),
              children: [
                // Top Score Summary Card
                ExamScoreSummary(result: result),
                const SizedBox(height: 14),

                // Performance Overview Chips
                Row(
                  children: [
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                        decoration: BoxDecoration(
                          color: const Color(0x1810B981),
                          borderRadius: BorderRadius.circular(14),
                          border: Border.all(color: const Color(0xFF10B981).withAlpha(60)),
                        ),
                        child: Row(
                          children: [
                            const Icon(Icons.check_circle_rounded, size: 18, color: Color(0xFF10B981)),
                            const SizedBox(width: 8),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  '$correctCount Correct',
                                  style: const TextStyle(
                                    color: Color(0xFF10B981),
                                    fontWeight: FontWeight.w800,
                                    fontSize: 14,
                                  ),
                                ),
                                Text(
                                  '${result.totalMarks > 0 ? (correctCount / result.feedback.length * 100).toStringAsFixed(0) : 0}% Accuracy',
                                  style: const TextStyle(color: Color(0xFFA1A1AA), fontSize: 11),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                        decoration: BoxDecoration(
                          color: missedCount > 0 ? const Color(0x18EF4444) : const Color(0xFF14141C),
                          borderRadius: BorderRadius.circular(14),
                          border: Border.all(
                            color: missedCount > 0
                                ? const Color(0xFFEF4444).withAlpha(60)
                                : const Color(0xFF262633),
                          ),
                        ),
                        child: Row(
                          children: [
                            Icon(
                              Icons.cancel_rounded,
                              size: 18,
                              color: missedCount > 0 ? const Color(0xFFEF4444) : const Color(0xFFA1A1AA),
                            ),
                            const SizedBox(width: 8),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  '$missedCount Missed',
                                  style: TextStyle(
                                    color: missedCount > 0 ? const Color(0xFFEF4444) : const Color(0xFFA1A1AA),
                                    fontWeight: FontWeight.w800,
                                    fontSize: 14,
                                  ),
                                ),
                                const Text(
                                  'Review below',
                                  style: TextStyle(color: Color(0xFFA1A1AA), fontSize: 11),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // Subject / Topic Breakdown (if available)
                if (breakdown.isNotEmpty) ...[
                  ExamSubjectBreakdown(scores: breakdown),
                  const SizedBox(height: 16),
                ],

                // Filter Tabs: All, Correct, Missed
                Container(
                  padding: const EdgeInsets.all(4),
                  decoration: BoxDecoration(
                    color: const Color(0xFF13131B),
                    borderRadius: BorderRadius.circular(14),
                    border: Border.all(color: const Color(0xFF22222E)),
                  ),
                  child: Row(
                    children: [
                      _buildFilterTab('All (${result.feedback.length})', 0),
                      _buildFilterTab('Correct ($correctCount)', 1),
                      _buildFilterTab('Missed ($missedCount)', 2),
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // Question Feedback Tiles List
                if (filteredFeedback.isEmpty)
                  Container(
                    padding: const EdgeInsets.symmetric(vertical: 40),
                    alignment: Alignment.center,
                    child: Column(
                      children: [
                        Icon(
                          _filterIndex == 1
                              ? Icons.sentiment_dissatisfied_rounded
                              : Icons.celebration_rounded,
                          size: 40,
                          color: const Color(0xFF52525B),
                        ),
                        const SizedBox(height: 10),
                        Text(
                          _filterIndex == 1
                              ? 'No correct answers yet.'
                              : 'Awesome! No missed questions in this section.',
                          style: const TextStyle(color: Color(0xFFA1A1AA), fontSize: 14),
                        ),
                      ],
                    ),
                  )
                else
                  ...filteredFeedback.map(
                    (entry) => ExamQuestionFeedbackTile(
                      feedback: entry.value,
                      index: entry.key,
                    ),
                  ),

                const SizedBox(height: 20),

                // Bottom Action Buttons
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton.icon(
                        onPressed: () => context.go(AppRoutes.mockExams),
                        icon: const Icon(Icons.arrow_back_rounded, size: 16),
                        label: const Text('Back to Exams', style: TextStyle(fontWeight: FontWeight.w700)),
                        style: OutlinedButton.styleFrom(
                          foregroundColor: const Color(0xFFA1A1AA),
                          side: const BorderSide(color: Color(0xFF2C2C3A)),
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton.icon(
                        onPressed: _isRetaking ? null : () => _handleRetake(result.examId),
                        icon: const Icon(Icons.replay_rounded, size: 16),
                        label: const Text(
                          'Retake Exam',
                          style: TextStyle(fontWeight: FontWeight.w800, fontSize: 14),
                        ),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF10B981),
                          foregroundColor: Colors.black,
                          elevation: 0,
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            );
          },
        ),
      ),
    );
  }

  Widget _buildFilterTab(String label, int index) {
    final isSelected = _filterIndex == index;
    return Expanded(
      child: GestureDetector(
        onTap: () => setState(() => _filterIndex = index),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 160),
          padding: const EdgeInsets.symmetric(vertical: 8),
          alignment: Alignment.center,
          decoration: BoxDecoration(
            color: isSelected ? const Color(0xFF1E1E2A) : Colors.transparent,
            borderRadius: BorderRadius.circular(10),
            border: isSelected ? Border.all(color: const Color(0xFF38384A)) : null,
          ),
          child: Text(
            label,
            style: TextStyle(
              color: isSelected ? Colors.white : const Color(0xFF71717A),
              fontWeight: isSelected ? FontWeight.w700 : FontWeight.w500,
              fontSize: 12,
            ),
          ),
        ),
      ),
    );
  }
}

class _ResultErrorState extends StatelessWidget {
  const _ResultErrorState({
    required this.message,
    required this.onRetry,
  });

  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSizes.screenPaddingH),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.assignment_late_outlined,
              color: AppColors.textSecondary,
              size: AppSizes.iconXl,
            ),
            const SizedBox(height: AppSizes.md),
            Text(
              message,
              textAlign: TextAlign.center,
              style: AppTextStyles.bodyMedium,
            ),
            const SizedBox(height: AppSizes.lg),
            AppButton(
              label: 'Retry',
              onPressed: onRetry,
              variant: AppButtonVariant.secondary,
              width: 120,
              height: AppSizes.buttonHeightSm,
            ),
          ],
        ),
      ),
    );
  }
}
