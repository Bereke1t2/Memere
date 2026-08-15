import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/exam_question_entity.dart';
import '../providers/exam_attempt_provider.dart';
import '../providers/exam_timer_provider.dart';
import '../widgets/exam_question_card.dart';
import '../widgets/exam_question_palette.dart';
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
          backgroundColor: const Color(0xFF0C0C10),
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
            icon: const Icon(Icons.close_rounded, color: Color(0xFFD4D4D8)),
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
                final progress = state.questionCount > 0
                    ? (state.currentIndex + 1) / state.questionCount
                    : 0.0;
                return LinearProgressIndicator(
                  minHeight: 3,
                  value: progress,
                  backgroundColor: const Color(0xFF181822),
                  valueColor: const AlwaysStoppedAnimation<Color>(
                    Color(0xFF10B981),
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
                    color: Color(0xFF10B981),
                    strokeWidth: 2.5,
                  ),
                  SizedBox(height: 16),
                  Text(
                    'Preparing your exam questions...',
                    style: TextStyle(
                      color: Color(0xFFA1A1AA),
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
          side: const BorderSide(color: Color(0xFF2C2C38)),
        ),
        title: const Row(
          children: [
            Icon(Icons.pause_circle_outline_rounded, color: Color(0xFFF59E0B)),
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
            color: Color(0xFFA1A1AA),
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
              style: TextStyle(color: Color(0xFFA1A1AA), fontWeight: FontWeight.w600),
            ),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFDC2626),
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
        backgroundColor: Color(0xFF10B981),
      ),
    );

    final resultAttemptId =
        await ref.read(examAttemptProvider(_params).notifier).submit();
    if (!mounted) return;
    if (resultAttemptId != null) {
      context.go(AppRoutes.examResultPath(resultAttemptId));
    }
  }

  Future<void> _confirmSubmit(ExamAttemptState state) async {
    final seconds = ref.read(examTimerProvider(_seedFor(state))).valueOrNull;

    final submit = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      backgroundColor: const Color(0xFF111118),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        side: BorderSide(color: Color(0xFF262633)),
      ),
      builder: (sheetContext) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Drag bar
              Center(
                child: Container(
                  width: 36,
                  height: 4,
                  decoration: BoxDecoration(
                    color: const Color(0xFF3F3F46),
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: const Color(0x2210B981),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: const Icon(Icons.assignment_turned_in_rounded,
                        color: Color(0xFF10B981), size: 24),
                  ),
                  const SizedBox(width: 12),
                  const Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Ready to submit?',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 19,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                        SizedBox(height: 2),
                        Text(
                          'Review your completion before finalizing',
                          style: TextStyle(
                            color: Color(0xFF71717A),
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 18),
              // Summary cards
              Row(
                children: [
                  Expanded(
                    child: Container(
                      padding: const EdgeInsets.all(14),
                      decoration: BoxDecoration(
                        color: const Color(0x1810B981),
                        borderRadius: BorderRadius.circular(14),
                        border: Border.all(color: const Color(0xFF10B981).withAlpha(60)),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            '${state.answeredCount}',
                            style: const TextStyle(
                              color: Color(0xFF10B981),
                              fontSize: 22,
                              fontWeight: FontWeight.w800,
                            ),
                          ),
                          const SizedBox(height: 2),
                          const Text(
                            'Answered',
                            style: TextStyle(
                              color: Color(0xFFA1A1AA),
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Container(
                      padding: const EdgeInsets.all(14),
                      decoration: BoxDecoration(
                        color: state.unansweredCount > 0
                            ? const Color(0x18F59E0B)
                            : const Color(0xFF1A1A22),
                        borderRadius: BorderRadius.circular(14),
                        border: Border.all(
                          color: state.unansweredCount > 0
                              ? const Color(0xFFF59E0B).withAlpha(60)
                              : const Color(0xFF2C2C38),
                        ),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            '${state.unansweredCount}',
                            style: TextStyle(
                              color: state.unansweredCount > 0
                                  ? const Color(0xFFF59E0B)
                                  : const Color(0xFFA1A1AA),
                              fontSize: 22,
                              fontWeight: FontWeight.w800,
                            ),
                          ),
                          const SizedBox(height: 2),
                          const Text(
                            'Unanswered',
                            style: TextStyle(
                              color: Color(0xFFA1A1AA),
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  if (seconds != null) ...[
                    const SizedBox(width: 10),
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.all(14),
                        decoration: BoxDecoration(
                          color: const Color(0xFF1A1A22),
                          borderRadius: BorderRadius.circular(14),
                          border: Border.all(color: const Color(0xFF2C2C38)),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              formatExamTimer(seconds),
                              style: const TextStyle(
                                color: Colors.white,
                                fontSize: 16,
                                fontWeight: FontWeight.w700,
                              ),
                            ),
                            const SizedBox(height: 2),
                            const Text(
                              'Time Left',
                              style: TextStyle(
                                color: Color(0xFFA1A1AA),
                                fontSize: 12,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ],
                ],
              ),
              const SizedBox(height: 14),
              if (state.unansweredCount > 0)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                  decoration: BoxDecoration(
                    color: const Color(0x18F59E0B),
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(color: const Color(0xFFF59E0B).withAlpha(60)),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.info_outline_rounded, size: 16, color: Color(0xFFF59E0B)),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          'You have ${state.unansweredCount} unanswered questions. Unanswered questions will receive 0 marks.',
                          style: const TextStyle(
                            color: Color(0xFFFCD34D),
                            fontSize: 12,
                            height: 1.3,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              const SizedBox(height: 20),
              Row(
                children: [
                  Expanded(
                    flex: 1,
                    child: OutlinedButton(
                      onPressed: () => Navigator.of(sheetContext).pop(false),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: const Color(0xFFA1A1AA),
                        side: const BorderSide(color: Color(0xFF323242)),
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                      ),
                      child: const Text('Review Answers', style: TextStyle(fontWeight: FontWeight.w600)),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    flex: 2,
                    child: ElevatedButton(
                      onPressed: () => Navigator.of(sheetContext).pop(true),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF10B981),
                        foregroundColor: Colors.black,
                        elevation: 0,
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                      ),
                      child: const Text(
                        'Submit Exam Now',
                        style: TextStyle(fontWeight: FontWeight.w800, fontSize: 15),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
            ],
          ),
        ),
      ),
    );

    if (submit != true) return;

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

    return ListView(
      physics: const BouncingScrollPhysics(),
      padding: const EdgeInsets.fromLTRB(18, 16, 18, 120),
      children: [
        // Question Card View
        ExamQuestionCard(
          question: currentQuestion,
          selectedAnswer: state.answers[currentQuestion.questionId],
          onSingleAnswerSelected: notifier.selectSingleAnswer,
          onMultiAnswerToggled: notifier.toggleMultiAnswer,
          onShortAnswerChanged: notifier.setShortAnswer,
        ),
        const SizedBox(height: 14),

        // Auto-save Status Indicator
        Align(
          alignment: Alignment.centerLeft,
          child: ExamSaveStatus(
            isSaving: state.isSaving,
            saveError: state.saveError,
            lastSavedAt: state.lastSavedAt,
          ),
        ),
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
        color: Color(0xFF0C0C10),
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
                color: isFirst ? const Color(0xFF111116) : const Color(0xFF1A1A24),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: isFirst ? const Color(0xFF1A1A22) : const Color(0xFF2E2E3C),
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
                    showExamQuestionPaletteSheet(
                      context: context,
                      count: state.questionCount,
                      currentIndex: state.currentIndex,
                      answeredIndexes: answeredIndexes,
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
                          color: Color(0xFF10B981),
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
                          color: Color(0xFFA1A1AA),
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
                      color: const Color(0xFF10B981).withAlpha(50),
                      blurRadius: 10,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                child: ElevatedButton(
                  onPressed: state.isSubmitting ? null : onSubmit,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF10B981),
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
                    colors: [Color(0xFF10B981), Color(0xFF059669)],
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: const Color(0xFF10B981).withAlpha(45),
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
