import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/auth/account_gate.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/connectivity_service.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../../../core/router/app_router.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import '../../../../shared/utils/formatters.dart';
import '../../../../shared/widgets/app_button.dart';
import '../providers/quiz_detail_provider.dart';
import '../providers/quiz_providers.dart';
import '../widgets/quiz_download_button.dart';

class QuizDetailScreen extends ConsumerWidget {
  const QuizDetailScreen({
    super.key,
    required this.quizId,
  });

  final String quizId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizAsync = ref.watch(quizDetailProvider(quizId));
    final isStarting = ref.watch(quizStartLoadingProvider(quizId));
    const emeraldColor = Color(0xFF10B981);

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(
        title: const Text('Interactive Quiz'),
        centerTitle: true,
      ),
      body: SafeArea(
        child: quizAsync.when(
          loading: () => const Center(
            child: CircularProgressIndicator(
              color: emeraldColor,
              strokeWidth: 2,
            ),
          ),
          error: (error, _) => _QuizErrorState(
            title: 'Could not load quiz',
            message: error is Failure ? error.message : 'Please try again.',
            onRetry: () => ref.invalidate(quizDetailProvider(quizId)),
          ),
          data: (quiz) => ListView(
            padding: const EdgeInsets.all(20),
            children: [
              Container(
                padding: const EdgeInsets.all(22),
                decoration: BoxDecoration(
                  color: AppColors.bgSecondary,
                  borderRadius: BorderRadius.circular(24),
                  border: Border.all(
                    color: emeraldColor.withValues(alpha: 0.3),
                    width: 1.2,
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: emeraldColor.withValues(alpha: 0.08),
                      blurRadius: 16,
                      offset: const Offset(0, 6),
                    ),
                  ],
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Container(
                          width: 52,
                          height: 52,
                          alignment: Alignment.center,
                          decoration: BoxDecoration(
                            color: const Color(0x2210B981),
                            borderRadius: BorderRadius.circular(16),
                            border: Border.all(color: const Color(0x5510B981)),
                          ),
                          child: const Icon(
                            Icons.quiz_rounded,
                            color: emeraldColor,
                            size: 26,
                          ),
                        ),
                        const SizedBox(width: 14),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Container(
                                padding: const EdgeInsets.symmetric(
                                    horizontal: 8, vertical: 3),
                                decoration: BoxDecoration(
                                  color: const Color(0x2210B981),
                                  borderRadius: BorderRadius.circular(8),
                                ),
                                child: const Text(
                                  'PRACTICE QUIZ',
                                  style: TextStyle(
                                    fontSize: 10.5,
                                    fontWeight: FontWeight.w800,
                                    color: emeraldColor,
                                    letterSpacing: 0.5,
                                  ),
                                ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                quiz.title,
                                maxLines: 2,
                                overflow: TextOverflow.ellipsis,
                                style: const TextStyle(
                                  fontSize: 18,
                                  fontWeight: FontWeight.w800,
                                  color: Colors.white,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 18),
                    Text(
                      quiz.hasAttemptsRemaining
                          ? 'Test your knowledge on this unit! Answer each question and get real-time graded results with detailed explanations.'
                          : 'You have used all attempts for this practice quiz.',
                      style: const TextStyle(
                        fontSize: 13.5,
                        color: Color(0xFF94A3B8),
                        height: 1.45,
                      ),
                    ),
                    const SizedBox(height: 22),
                    const Divider(color: AppColors.borderStrong, height: 1),
                    const SizedBox(height: 18),
                    _MetaRow(
                      icon: Icons.help_outline_rounded,
                      label: 'Total Questions',
                      value: '${quiz.questionCount} Questions',
                    ),
                    _MetaRow(
                      icon: Icons.percent_rounded,
                      label: 'Passing Score',
                      value: '${quiz.passPercentage.toStringAsFixed(0)}%',
                    ),
                    _MetaRow(
                      icon: Icons.timer_outlined,
                      label: 'Time Limit',
                      value: quiz.timeLimitSeconds == null
                          ? 'Untimed'
                          : formatDurationSeconds(quiz.timeLimitSeconds!),
                    ),
                    _MetaRow(
                      icon: Icons.replay_rounded,
                      label: 'Attempts Used',
                      value: quiz.maxAttempts == null
                          ? '${quiz.attemptsUsed} attempts'
                          : '${quiz.attemptsUsed}/${quiz.maxAttempts}',
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 24),
              AppButton(
                label: quiz.hasAttemptsRemaining
                    ? 'Start Quiz Now'
                    : 'Attempts Exhausted',
                isLoading: isStarting,
                isDisabled: !quiz.hasAttemptsRemaining,
                onPressed: () => _startAttempt(context, ref),
                suffixIcon: Icons.play_arrow_rounded,
              ),
              const SizedBox(height: 12),
              QuizDownloadButton(quizId: quizId, title: quiz.title),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _startAttempt(BuildContext context, WidgetRef ref) async {
    final loading = ref.read(quizStartLoadingProvider(quizId).notifier);
    loading.state = true;

    final downloaded =
        ref.read(downloadStoreProvider).isDownloaded(DownloadType.quiz, quizId);
    final online = await ref.read(connectivityServiceProvider).isOnline();
    // Guests are always graded on-device; so is anyone currently offline. Only
    // an online, signed-in user takes the authoritative server path.
    final wantLocal = isGuest(ref) || !online;

    if (wantLocal) {
      loading.state = false;
      if (!context.mounted) return;
      if (!downloaded) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              online
                  ? 'Download this quiz first to take it as a guest.'
                  : "You're offline. Download this quiz first to take it.",
            ),
          ),
        );
        return;
      }
      _goToLocalAttempt(context);
      return;
    }

    final result = await ref.read(startQuizAttemptUseCaseProvider)(quizId);
    loading.state = false;
    if (!context.mounted) return;

    result.fold(
      (failure) {
        // Started online but the network dropped — grade the downloaded copy.
        if (isOfflineFailure(failure) && downloaded) {
          _goToLocalAttempt(context);
          return;
        }
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(failure.message)),
        );
      },
      (attempt) {
        context.go(
          AppRoutes.quizAttemptPath(
            attemptId: attempt.attemptId,
            quizId: attempt.quizId,
          ),
        );
      },
    );
  }

  void _goToLocalAttempt(BuildContext context) {
    context.go(
      AppRoutes.quizAttemptPath(
        attemptId: newLocalAttemptId(),
        quizId: quizId,
      ),
    );
  }
}

class _MetaRow extends StatelessWidget {
  const _MetaRow({
    required this.icon,
    required this.label,
    required this.value,
  });

  final IconData icon;
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        children: [
          Icon(icon, size: 18, color: const Color(0xFF64748B)),
          const SizedBox(width: 10),
          Expanded(
            child: Text(
              label,
              style: const TextStyle(
                fontSize: 13.5,
                color: Color(0xFF94A3B8),
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              fontSize: 13.5,
              fontWeight: FontWeight.w700,
              color: Colors.white,
            ),
          ),
        ],
      ),
    );
  }
}

class _QuizErrorState extends StatelessWidget {
  const _QuizErrorState({
    required this.title,
    required this.message,
    required this.onRetry,
  });

  final String title;
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
              Icons.error_outline_rounded,
              color: AppColors.error,
              size: AppSizes.iconXl,
            ),
            const SizedBox(height: AppSizes.md),
            Text(title, style: AppTextStyles.headlineSmall),
            const SizedBox(height: AppSizes.sm),
            Text(
              message,
              textAlign: TextAlign.center,
              style: AppTextStyles.bodySmall,
            ),
            const SizedBox(height: AppSizes.lg),
            AppButton(
              label: 'Retry',
              onPressed: onRetry,
              variant: AppButtonVariant.secondary,
            ),
          ],
        ),
      ),
    );
  }
}
