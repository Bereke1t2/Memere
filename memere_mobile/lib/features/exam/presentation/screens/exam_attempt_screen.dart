import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/providers/flagged_questions_provider.dart';
import '../../../../shared/providers/reading_scale_provider.dart';
import '../../../../shared/widgets/assessment/question_palette_sheet.dart';
import '../../../../shared/widgets/assessment/question_status.dart';
import '../../../../shared/widgets/assessment/submit_summary_sheet.dart';
import '../../domain/entities/exam_question_entity.dart';
import '../providers/exam_attempt_provider.dart';
import '../providers/exam_timer_provider.dart';
import '../widgets/exam_question_card.dart';
import '../widgets/exam_save_status.dart';
import '../widgets/exam_timer_badge.dart';

class ExamAttemptScreen extends ConsumerStatefulWidget {
  const ExamAttemptScreen({
    super.key,
    required this.attemptId,
    required this.examId,
  });

  final String attemptId;
  final String examId;

  @override
  ConsumerState<ExamAttemptScreen> createState() => _ExamAttemptScreenState();
}

class _ExamAttemptScreenState extends ConsumerState<ExamAttemptScreen> {
  bool _autoSubmitTriggered = false;

  ExamAttemptParams get _params =>
      ExamAttemptParams(attemptId: widget.attemptId, examId: widget.examId);

  void _safeExit() {
    if (context.canPop()) {
      context.pop();
    } else {
      context.go(AppRoutes.mockExams);
    }
  }

  @override
  Widget build(BuildContext context) {
    final attemptAsync = ref.watch(examAttemptProvider(_params));

    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) async {
        if (didPop) return;
        if (await _confirmLeave() && mounted) {
          _safeExit();
        }
      },
      child: Scaffold(
        backgroundColor: const Color(0xFF070709),
        appBar: AppBar(
          backgroundColor: const Color(0xFF0B0E0C),
          elevation: 0,
          scrolledUnderElevation: 0,
          automaticallyImplyLeading: false,
          leading: IconButton(
            tooltip: 'Leave Exam',
            onPressed: () async {
              if (await _confirmLeave() && mounted) {
                _safeExit();
              }
            },
            icon: const Icon(Icons.close_rounded, color: Color(0xFFD7D9D4)),
          ),
          title: attemptAsync.maybeWhen(
            data: (state) => Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              decoration: BoxDecoration(
                color: const Color(0xFF181822),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: const Color(0xFF2C2C3A)),
              ),
              child: Text(
                'Question ${state.currentIndex + 1} of ${state.questionCount}',
                style: AppTextStyles.labelMedium.copyWith(
                  fontWeight: FontWeight.w700,
                  color: Colors.white,
                  letterSpacing: 0.2,
                ),
              ),
            ),
            orElse: () => const Text('Mock Exam'),
          ),
          centerTitle: true,
          actions: [
            attemptAsync.maybeWhen(
              data: (_) => const _TextSizeAction(),
              orElse: () => const SizedBox.shrink(),
            ),
            attemptAsync.maybeWhen(
              data: (state) {
                final seed = _seedFor(state);
                final seconds = ref.watch(examTimerProvider(seed)).valueOrNull;
                return Padding(
                  padding: const EdgeInsets.only(right: 12),
                  child: ExamTimerBadge(seconds: seconds),
                );
              },
              orElse: () => const SizedBox.shrink(),
            ),
          ],
          bottom: PreferredSize(
            preferredSize: const Size.fromHeight(3.0),
            child: attemptAsync.maybeWhen(
              data: (state) {
                // Answered, not position: what is left to do is the fact a
                // student is actually tracking during a timed paper.
                final progress = state.questionCount > 0
                    ? state.answeredCount / state.questionCount
                    : 0.0;
                return LinearProgressIndicator(
                  minHeight: 3,
                  value: progress,
                  backgroundColor: const Color(0xFF181822),
                  valueColor: const AlwaysStoppedAnimation<Color>(
                    Color(0xFF35A862),
                  ),
                );
              },
              orElse: () => const SizedBox(height: 3),
            ),
          ),
        ),
        body: SafeArea(
          top: false,
          child: attemptAsync.when(
            loading: () => const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  CircularProgressIndicator(
                    color: Color(0xFF35A862),
                    strokeWidth: 2.5,
                  ),
                  SizedBox(height: 16),
                  Text(
                    'Preparing your exam questions...',
                    style: TextStyle(
                      color: Color(0xFFA6ABA4),
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              ),
            ),
            error: (error, _) => _AttemptErrorState(
              message: error is Failure
                  ? error.message
                  : 'Could not start exam. Please try again.',
              onRetry: () => ref.invalidate(examAttemptProvider(_params)),
            ),
            data: (state) {
              if (state.attempt.questions.isEmpty) {
                return const _AttemptErrorState(
                  message: 'This exam has no questions yet.',
                );
              }
              return _AttemptBody(
                params: _params,
                state: state,
                onTimerZero: () => _handleExpiry(state),
                onSubmit: () => _confirmSubmit(state),
              );
            },
          ),
        ),
        bottomNavigationBar: attemptAsync.maybeWhen(
          data: (state) => _StickyBottomBar(
            params: _params,
            state: state,
            onSubmit: () => _confirmSubmit(state),
          ),
          orElse: () => null,
        ),
      ),
    );
  }

  Future<bool> _confirmLeave() async {
    final notifier = ref.read(examAttemptProvider(_params).notifier);
    final leave = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        backgroundColor: const Color(0xFF14141B),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
          side: const BorderSide(color: Color(0xFF2A302C)),
        ),
        title: const Row(
          children: [
            Icon(Icons.pause_circle_outline_rounded, color: Color(0xFFD99A2B)),
            SizedBox(width: 8),
            Text(
              'Leave exam?',
              style: TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.w700,
                fontSize: 18,
              ),
            ),
          ],
        ),
        content: const Text(
          'Your progress is automatically saved to the server. '
          'The exam timer will continue running in the background until it expires.',
          style: TextStyle(
            color: Color(0xFFA6ABA4),
            fontSize: 14,
            height: 1.4,
          ),
        ),
        actionsPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text(
              'Stay & Continue',
              style: TextStyle(color: Color(0xFFA6ABA4), fontWeight: FontWeight.w600),
            ),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFC44A41),
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
            ),
            child: const Text('Leave', style: TextStyle(fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
    if (leave == true) {
      await notifier.saveProgress(force: true);
      return true;
    }
    return false;
  }

  Future<void> _handleExpiry(ExamAttemptState state) async {
    if (_autoSubmitTriggered) return;
    _autoSubmitTriggered = true;

    final messenger = ScaffoldMessenger.of(context);
    messenger.showSnackBar(
      const SnackBar(
        content: Text('Time is up! Submitting your answers...'),
        backgroundColor: Color(0xFF35A862),
      ),
    );

    final resultAttemptId =
        await ref.read(examAttemptProvider(_params).notifier).submit();
    if (resultAttemptId != null) {
      await ref
          .read(flaggedQuestionsProvider(widget.attemptId).notifier)
          .clear();
    }
    if (!mounted) return;
    if (resultAttemptId != null) {
      context.go(AppRoutes.examResultPath(resultAttemptId));
    }
  }

  Future<void> _confirmSubmit(ExamAttemptState state) async {
    final flagged = ref.read(flaggedQuestionsProvider(widget.attemptId));
    final statuses = _statuses(state, _answeredIndexes(state), flagged);
    final firstBlank = statuses.indexWhere((s) => !s.isAnswered);
    final seconds = ref.read(examTimerProvider(_seedFor(state))).valueOrNull;

    final submit = await SubmitSummarySheet.show(
      context,
      title: 'Submit exam?',
      total: statuses.length,
      answered: statuses.where((s) => s.isAnswered).length,
      flagged: statuses.where((s) => s.isFlagged).length,
      isTimeUp: seconds == 0,
      onReviewBlank: firstBlank < 0
          ? null
          : () => ref
              .read(examAttemptProvider(_params).notifier)
              .goToQuestion(firstBlank),
    );

    if (!submit) return;

    final resultAttemptId =
        await ref.read(examAttemptProvider(_params).notifier).submit();
    if (!mounted) return;
    if (resultAttemptId == null) {
      final message =
          ref.read(examAttemptProvider(_params)).valueOrNull?.saveError ??
              'Could not submit exam. Please try again.';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(message)),
      );
      return;
    }
    context.go(AppRoutes.examResultPath(resultAttemptId));
  }

  ExamTimerSeed _seedFor(ExamAttemptState state) => ExamTimerSeed(
        attemptId: state.attempt.attemptId,
        expiresAt: state.attempt.expiresAt,
        remainingSeconds: state.attempt.remainingSeconds,
      );
}

class _AttemptBody extends ConsumerWidget {
  const _AttemptBody({
    required this.params,
    required this.state,
    required this.onTimerZero,
    required this.onSubmit,
  });

  final ExamAttemptParams params;
  final ExamAttemptState state;
  final VoidCallback onTimerZero;
  final Future<void> Function() onSubmit;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final notifier = ref.read(examAttemptProvider(params).notifier);
    final seed = ExamTimerSeed(
      attemptId: state.attempt.attemptId,
      expiresAt: state.attempt.expiresAt,
      remainingSeconds: state.attempt.remainingSeconds,
    );

    ref.listen(examTimerProvider(seed), (_, next) {
      if (next.valueOrNull == 0) onTimerZero();
    });

    final currentQuestion = state.attempt.questions[state.currentIndex];

    final scale = ref.watch(readingScaleProvider);
    final attemptId = state.attempt.attemptId;
    final flagged = ref.watch(flaggedQuestionsProvider(attemptId));
    final isFlagged = flagged.contains(currentQuestion.questionId);
    final answeredIndexes = _answeredIndexes(state);
    final firstBlank = List.generate(
      state.questionCount,
      answeredIndexes.contains,
    ).indexWhere((answered) => !answered);

    return ListView(
      physics: const BouncingScrollPhysics(),
      padding: const EdgeInsets.fromLTRB(18, 16, 18, 120),
      children: [
        // The student's chosen size applies to the question content only, so
        // scaling it can never break the chrome around it.
        MediaQuery.withClampedTextScaling(
          minScaleFactor: scale,
          maxScaleFactor: scale,
          child: ExamQuestionCard(
            question: currentQuestion,
            selectedAnswer: state.answers[currentQuestion.questionId],
            onSingleAnswerSelected: notifier.selectSingleAnswer,
            onMultiAnswerToggled: notifier.toggleMultiAnswer,
            onShortAnswerChanged: notifier.setShortAnswer,
          ),
        ),
        const SizedBox(height: 14),

        Row(
          children: [
            _ExamFlagButton(
              isFlagged: isFlagged,
              onTap: () => ref
                  .read(flaggedQuestionsProvider(attemptId).notifier)
                  .toggle(currentQuestion.questionId),
            ),
            const Spacer(),
            ExamSaveStatus(
              isSaving: state.isSaving,
              saveError: state.saveError,
              lastSavedAt: state.lastSavedAt,
            ),
          ],
        ),

        // Paging one at a time through a 100-question paper to find what is
        // left is the slowest part of finishing it.
        if (firstBlank >= 0 && firstBlank != state.currentIndex) ...[
          const SizedBox(height: 14),
          SizedBox(
            width: double.infinity,
            child: TextButton.icon(
              onPressed: () => notifier.goToQuestion(firstBlank),
              icon: const Icon(Icons.skip_next_rounded, size: 18),
              label: Text('Next unanswered (Q${firstBlank + 1})'),
              style: TextButton.styleFrom(
                foregroundColor: AppColors.accentPrimary,
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
        ],
      ],
    );
  }
}

class _StickyBottomBar extends ConsumerWidget {
  const _StickyBottomBar({
    required this.params,
    required this.state,
    required this.onSubmit,
  });

  final ExamAttemptParams params;
  final ExamAttemptState state;
  final VoidCallback onSubmit;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final notifier = ref.read(examAttemptProvider(params).notifier);
    final answeredIndexes = _answeredIndexes(state);
    final isFirst = state.currentIndex == 0;
    final isLast = state.currentIndex == state.questionCount - 1;

    return Container(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 18),
      decoration: const BoxDecoration(
        color: Color(0xFF0B0E0C),
        border: Border(
          top: BorderSide(color: Color(0xFF1E1E28), width: 1.2),
        ),
      ),
      child: SafeArea(
        top: false,
        child: Row(
          children: [
            // Previous Button
            Container(
              height: 46,
              decoration: BoxDecoration(
                color: isFirst ? const Color(0xFF101311) : const Color(0xFF1A1A24),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: isFirst ? const Color(0xFF191D1A) : const Color(0xFF2E2E3C),
                ),
              ),
              child: IconButton(
                onPressed: isFirst ? null : notifier.previousQuestion,
                tooltip: 'Previous Question',
                icon: Icon(
                  Icons.arrow_back_ios_new_rounded,
                  size: 17,
                  color: isFirst ? const Color(0xFF3F3F4E) : Colors.white,
                ),
              ),
            ),
            const SizedBox(width: 10),

            // Question Navigator Pill (Opens 60-question grid modal)
            Expanded(
              child: Material(
                color: Colors.transparent,
                child: InkWell(
                  onTap: () {
                    QuestionPaletteSheet.show(
                      context,
                      statuses: _statuses(
                        state,
                        answeredIndexes,
                        ref.read(
                          flaggedQuestionsProvider(state.attempt.attemptId),
                        ),
                      ),
                      currentIndex: state.currentIndex,
                      onSelected: notifier.goToQuestion,
                    );
                  },
                  borderRadius: BorderRadius.circular(12),
                  child: Container(
                    height: 46,
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                    decoration: BoxDecoration(
                      color: const Color(0xFF13131B),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: const Color(0xFF262634)),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Icon(
                          Icons.grid_view_rounded,
                          size: 16,
                          color: Color(0xFF35A862),
                        ),
                        const SizedBox(width: 7),
                        Text(
                          '${state.answeredCount}/${state.questionCount}',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 13,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                        const SizedBox(width: 4),
                        const Icon(
                          Icons.keyboard_arrow_up_rounded,
                          size: 16,
                          color: Color(0xFFA6ABA4),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 10),

            // Next / Submit Button
            if (isLast)
              Container(
                height: 46,
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(12),
                  boxShadow: [
                    BoxShadow(
                      color: const Color(0xFF35A862).withAlpha(50),
                      blurRadius: 10,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                child: ElevatedButton(
                  onPressed: state.isSubmitting ? null : onSubmit,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF35A862),
                    foregroundColor: Colors.black,
                    elevation: 0,
                    padding: const EdgeInsets.symmetric(horizontal: 18),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                  child: state.isSubmitting
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.black,
                          ),
                        )
                      : const Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              'Submit',
                              style: TextStyle(
                                fontWeight: FontWeight.w800,
                                fontSize: 14,
                              ),
                            ),
                            SizedBox(width: 4),
                            Icon(Icons.check_rounded, size: 17),
                          ],
                        ),
                ),
              )
            else
              Container(
                height: 46,
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(12),
                  gradient: const LinearGradient(
                    colors: [Color(0xFF35A862), Color(0xFF24814A)],
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: const Color(0xFF35A862).withAlpha(45),
                      blurRadius: 10,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                child: ElevatedButton(
                  onPressed: notifier.nextQuestion,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.transparent,
                    shadowColor: Colors.transparent,
                    foregroundColor: Colors.black,
                    padding: const EdgeInsets.symmetric(horizontal: 18),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                  child: const Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        'Next',
                        style: TextStyle(
                          fontWeight: FontWeight.w800,
                          fontSize: 14,
                        ),
                      ),
                      SizedBox(width: 4),
                      Icon(Icons.arrow_forward_ios_rounded, size: 14),
                    ],
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class _AttemptErrorState extends StatelessWidget {
  const _AttemptErrorState({
    required this.message,
    this.onRetry,
  });

  final String message;
  final VoidCallback? onRetry;

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
            if (onRetry != null) ...[
              const SizedBox(height: AppSizes.lg),
              AppButton(
                label: 'Retry',
                onPressed: onRetry,
                variant: AppButtonVariant.secondary,
                width: 120,
                height: AppSizes.buttonHeightSm,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

Set<int> _answeredIndexes(ExamAttemptState state) {
  final indexes = <int>{};
  for (var i = 0; i < state.attempt.questions.length; i++) {
    final question = state.attempt.questions[i];
    final answer = state.answers[question.questionId];
    if (_hasAnswer(question, answer)) indexes.add(i);
  }
  return indexes;
}

bool _hasAnswer(ExamQuestionEntity question, Object? answer) {
  if (answer == null) return false;
  if (question.type == ExamQuestionType.shortAnswer) {
    return answer is String && answer.trim().isNotEmpty;
  }
  if (answer is List) return answer.isNotEmpty;
  return answer is String && answer.trim().isNotEmpty;
}


/// Combines answers and flags into the per-question status the palette and the
/// submit summary both read, so the two can never disagree about what is left.
List<QuestionStatus> _statuses(
  ExamAttemptState state,
  Set<int> answeredIndexes,
  Set<String> flagged,
) {
  return List.generate(state.questionCount, (index) {
    return QuestionStatus.of(
      answered: answeredIndexes.contains(index),
      flagged: flagged.contains(state.attempt.questions[index].questionId),
    );
  });
}

/// A–/A+ in the app bar. Sits next to the timer because both are things a
/// student reaches for mid-question without wanting to leave the page.
class _TextSizeAction extends ConsumerWidget {
  const _TextSizeAction();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.watch(readingScaleProvider);
    final ctrl = ref.read(readingScaleProvider.notifier);
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        IconButton(
          tooltip: 'Smaller text',
          visualDensity: VisualDensity.compact,
          onPressed: ctrl.canDecrease ? ctrl.decrease : null,
          icon: const Icon(Icons.text_decrease_rounded, size: 20),
          color: AppColors.textSecondary,
          disabledColor: AppColors.textDisabled,
        ),
        IconButton(
          tooltip: 'Larger text (${ctrl.label})',
          visualDensity: VisualDensity.compact,
          onPressed: ctrl.canIncrease ? ctrl.increase : null,
          icon: const Icon(Icons.text_increase_rounded, size: 20),
          color: AppColors.textSecondary,
          disabledColor: AppColors.textDisabled,
        ),
      ],
    );
  }
}

class _ExamFlagButton extends StatelessWidget {
  const _ExamFlagButton({required this.isFlagged, required this.onTap});

  final bool isFlagged;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final color = isFlagged ? AppColors.warning : AppColors.textMuted;
    return Semantics(
      button: true,
      selected: isFlagged,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(999),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 7),
          decoration: BoxDecoration(
            color: isFlagged ? AppColors.warningSurface : Colors.transparent,
            borderRadius: BorderRadius.circular(999),
            border: Border.all(
              color: isFlagged ? AppColors.warning : AppColors.borderStrong,
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                isFlagged ? Icons.flag_rounded : Icons.flag_outlined,
                size: 14,
                color: color,
              ),
              const SizedBox(width: 6),
              Text(
                isFlagged ? 'Flagged' : 'Flag for review',
                style: AppTextStyles.caption.copyWith(
                  color: color,
                  fontWeight: isFlagged ? FontWeight.w700 : FontWeight.w500,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
