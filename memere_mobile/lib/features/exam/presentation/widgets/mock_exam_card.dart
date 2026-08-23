import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/auth/account_gate.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/connectivity_service.dart';
import '../../../../core/offline/offline_attempt_factory.dart';
import '../../../../core/offline/offline_providers.dart';
import '../../../../core/router/app_router.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import '../../../../shared/widgets/ai_robot_mascot.dart';
import '../../domain/entities/mock_exam_entity.dart';
import '../providers/exam_attempt_provider.dart';
import '../providers/exam_providers.dart';
import 'exam_attempt_history_sheet.dart';
import 'exam_download_button.dart';

/// Rich Vibrant Gradient 2-Column Exam Card with attempt tracking and history drawer.
class MockExamCard extends ConsumerStatefulWidget {
  const MockExamCard({
    super.key,
    required this.exam,
    this.cardIndex = 0,
  });

  final MockExamEntity exam;
  final int cardIndex;

  @override
  ConsumerState<MockExamCard> createState() => _MockExamCardState();
}

class _MockExamCardState extends ConsumerState<MockExamCard> {
  bool _isStarting = false;
  bool _isPressed = false;

  static const List<LinearGradient> _cardGradients = [
    LinearGradient(
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
      colors: [Color(0xFF0F261E), Color(0xFF141926)],
    ),
    LinearGradient(
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
      colors: [Color(0xFF0E2232), Color(0xFF141926)],
    ),
    LinearGradient(
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
      colors: [Color(0xFF241B12), Color(0xFF141926)],
    ),
    LinearGradient(
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
      colors: [Color(0xFF1D1B36), Color(0xFF141926)],
    ),
  ];

  static const List<Color> _borderGlowColors = [
    Color(0xFF34D399),
    Color(0xFF38BDF8),
    Color(0xFFFBBF24),
    Color(0xFFA78BFA),
  ];

  @override
  Widget build(BuildContext context) {
    final exam = widget.exam;
    final gradient = _cardGradients[widget.cardIndex % _cardGradients.length];
    final glowColor =
        _borderGlowColors[widget.cardIndex % _borderGlowColors.length];
    final attempts = ref.watch(examAttemptsByExamProvider(exam.id)).valueOrNull ?? [];
    final latestAttempt = attempts.isNotEmpty ? attempts.first : null;
    final isTaken = latestAttempt != null && latestAttempt.isGraded;
    final isInProgress = latestAttempt != null && !latestAttempt.isGraded;

    return GestureDetector(
      onTapDown: (_) => setState(() => _isPressed = true),
      onTapUp: (_) => setState(() => _isPressed = false),
      onTapCancel: () => setState(() => _isPressed = false),
      child: AnimatedScale(
        scale: _isPressed ? 0.96 : 1.0,
        duration: const Duration(milliseconds: 120),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(16),
          child: Container(
            padding: const EdgeInsets.all(10),
            decoration: BoxDecoration(
              gradient: gradient,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(
                color: _isPressed ? AppColors.borderFocused : glowColor.withAlpha(60),
                width: 1.2,
              ),
              boxShadow: [
                BoxShadow(
                  color: glowColor.withAlpha(20),
                  blurRadius: 8,
                  offset: const Offset(0, 3),
                ),
              ],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                // Top Meta Pill Row: Status / History + Marks
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    if (isTaken)
                      Flexible(
                        child: FittedBox(
                          fit: BoxFit.scaleDown,
                          alignment: Alignment.centerLeft,
                          child: GestureDetector(
                            onTap: () => showExamAttemptHistorySheet(
                              context: context,
                              examId: exam.id,
                              examTitle: exam.title,
                              onRetake: () => _startExam(context),
                            ),
                            child: Container(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 6,
                                vertical: 3,
                              ),
                              decoration: BoxDecoration(
                                color: (latestAttempt.isPassed
                                        ? const Color(0xFF10B981)
                                        : const Color(0xFFF59E0B))
                                    .withAlpha(90),
                                borderRadius: BorderRadius.circular(6),
                                border: Border.all(
                                  color: Colors.white.withAlpha(60),
                                  width: 0.8,
                                ),
                              ),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Icon(
                                    latestAttempt.isPassed
                                        ? Icons.check_circle_rounded
                                        : Icons.refresh_rounded,
                                    size: 11,
                                    color: Colors.white,
                                  ),
                                  const SizedBox(width: 3),
                                  Text(
                                    '${latestAttempt.percentage?.toStringAsFixed(0)}% • ${latestAttempt.isPassed ? "Pass" : "Retry"}',
                                    style: const TextStyle(
                                      fontSize: 10,
                                      fontWeight: FontWeight.w800,
                                      color: Colors.white,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      )
                    else if (isInProgress)
                      Flexible(
                        child: FittedBox(
                          fit: BoxFit.scaleDown,
                          alignment: Alignment.centerLeft,
                          child: Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 3,
                            ),
                            decoration: BoxDecoration(
                              color: const Color(0xFF0284C7).withAlpha(100),
                              borderRadius: BorderRadius.circular(6),
                            ),
                            child: const Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Icon(
                                  Icons.play_arrow_rounded,
                                  size: 11,
                                  color: Colors.white,
                                ),
                                SizedBox(width: 3),
                                Text(
                                  'In Progress',
                                  style: TextStyle(
                                    fontSize: 10,
                                    fontWeight: FontWeight.bold,
                                    color: Colors.white,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      )
                    else
                      Flexible(
                        child: FittedBox(
                          fit: BoxFit.scaleDown,
                          alignment: Alignment.centerLeft,
                          child: Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 3,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.black.withAlpha(75),
                              borderRadius: BorderRadius.circular(6),
                            ),
                            child: const Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Icon(
                                  Icons.star_rounded,
                                  size: 11,
                                  color: AppColors.brandAmber,
                                ),
                                SizedBox(width: 3),
                                Text(
                                  'New Exam',
                                  style: TextStyle(
                                    fontSize: 10,
                                    fontWeight: FontWeight.bold,
                                    color: Colors.white,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    const SizedBox(width: 4),
                    Flexible(
                      child: FittedBox(
                        fit: BoxFit.scaleDown,
                        alignment: Alignment.centerRight,
                        child: Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 6,
                            vertical: 3,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.black.withAlpha(75),
                            borderRadius: BorderRadius.circular(6),
                          ),
                          child: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              const Icon(
                                Icons.emoji_events_rounded,
                                size: 11,
                                color: AppColors.brandAmber,
                              ),
                              const SizedBox(width: 3),
                              Text(
                                '${exam.totalMarks} Marks',
                                style: const TextStyle(
                                  fontSize: 10,
                                  fontWeight: FontWeight.bold,
                                  color: Colors.white,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ],
                ),

                // Center AI Robot 3D Mascot Graphic
                Expanded(
                  child: Center(
                    child: Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4),
                      child: FittedBox(
                        fit: BoxFit.scaleDown,
                        child: SizedBox(
                          width: 82,
                          height: 72,
                          child: AiRobotMascot(
                            size: 64,
                            backgroundColor: glowColor,
                          ),
                        ),
                      ),
                    ),
                  ),
                ),

                // Title & Info
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      exam.title,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.w800,
                        color: Colors.white,
                        height: 1.2,
                      ),
                    ),
                    const SizedBox(height: 3),
                    Row(
                      children: [
                        const Icon(
                          Icons.timer_outlined,
                          size: 11,
                          color: Colors.white70,
                        ),
                        const SizedBox(width: 3),
                        Expanded(
                          child: Text(
                            '${exam.durationMinutes} min',
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                              fontSize: 10,
                              color: Colors.white70,
                            ),
                          ),
                        ),
                        if (attempts.isNotEmpty)
                          GestureDetector(
                            onTap: () => showExamAttemptHistorySheet(
                              context: context,
                              examId: exam.id,
                              examTitle: exam.title,
                              onRetake: () => _startExam(context),
                            ),
                            child: Container(
                              padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 1),
                              decoration: BoxDecoration(
                                color: Colors.black.withAlpha(60),
                                borderRadius: BorderRadius.circular(4),
                              ),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  const Icon(Icons.history_rounded, size: 10, color: Colors.white70),
                                  const SizedBox(width: 2),
                                  Text(
                                    '${attempts.length}',
                                    style: const TextStyle(
                                      fontSize: 9,
                                      fontWeight: FontWeight.w700,
                                      color: Colors.white,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                      ],
                    ),
                  ],
                ),

                const SizedBox(height: 6),

                // Bottom CTA Pill Button ("Join" / "Retake" / "Resume")
                SizedBox(
                  width: double.infinity,
                  height: 32,
                  child: ElevatedButton(
                    onPressed: _isStarting ? null : () => _confirmStart(context),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppColors.bgTertiary,
                      foregroundColor: glowColor,
                      elevation: 0,
                      padding: EdgeInsets.zero,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                        side: BorderSide(color: glowColor.withAlpha(80)),
                      ),
                    ),
                    child: _isStarting
                        ? const SizedBox(
                            width: 12,
                            height: 12,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              color: Colors.white,
                            ),
                          )
                        : FittedBox(
                            fit: BoxFit.scaleDown,
                            child: Text(
                              isInProgress
                                  ? 'Resume'
                                  : isTaken
                                      ? 'Retake'
                                      : 'Join',
                              style: const TextStyle(
                                fontSize: 12.5,
                                fontWeight: FontWeight.bold,
                                letterSpacing: 0.2,
                              ),
                            ),
                          ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _confirmStart(BuildContext context) async {
    final exam = widget.exam;
    final start = await showModalBottomSheet<bool>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (sheetContext) => SafeArea(
        child: SingleChildScrollView(
          physics: const BouncingScrollPhysics(),
          padding: const EdgeInsets.all(22),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Container(
                width: 44,
                height: 4,
                decoration: BoxDecoration(
                  color: AppColors.borderStrong,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 16),

              // Animated Mascot
              const AiRobotMascot(
                size: 72,
                backgroundColor: AppColors.brandEmerald,
              ),
              const SizedBox(height: 14),

              Text(
                exam.title,
                textAlign: TextAlign.center,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: AppColors.textPrimary,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                '${exam.subject} • Grade ${exam.grade}',
                style: const TextStyle(
                  fontSize: 13,
                  color: AppColors.brandEmerald,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 18),

              _ConfirmRow(label: 'Exam Duration', value: exam.durationLabel),
              _ConfirmRow(
                label: 'Total Questions / Marks',
                value: '${exam.totalMarks} Marks',
              ),
              _ConfirmRow(
                label: 'Passing Standard',
                value: 'Pass ${exam.passMarks} Marks',
              ),
              const SizedBox(height: 14),

              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: AppColors.bgTertiary,
                  borderRadius: BorderRadius.circular(14),
                  border: Border.all(color: AppColors.borderStrong),
                ),
                child: const Row(
                  children: [
                    AiRobotMascot(size: 26),
                    SizedBox(width: 10),
                    Expanded(
                      child: Text(
                        'The countdown starts as soon as you tap Join Exam. Good luck!',
                        style: TextStyle(
                          fontSize: 12,
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 14),
              ExamDownloadButton(
                examId: widget.exam.id,
                title: widget.exam.title,
              ),
              const SizedBox(height: 22),

              Row(
                children: [
                  Expanded(
                    child: SizedBox(
                      height: 46,
                      child: OutlinedButton(
                        onPressed: () => Navigator.of(sheetContext).pop(false),
                        style: OutlinedButton.styleFrom(
                          foregroundColor: AppColors.textSecondary,
                          side: const BorderSide(color: AppColors.borderStrong),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(14),
                          ),
                        ),
                        child: const Text('Cancel'),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: SizedBox(
                      height: 46,
                      child: ElevatedButton.icon(
                        onPressed: () => Navigator.of(sheetContext).pop(true),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.brandEmerald,
                          foregroundColor: Colors.white,
                          elevation: 0,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(14),
                          ),
                        ),
                        icon: const Icon(Icons.play_arrow_rounded, size: 20),
                        label: const Text(
                          'Join Exam',
                          style: TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
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

    final examId = widget.exam.id;
    final downloaded =
        ref.read(downloadStoreProvider).isDownloaded(DownloadType.exam, examId);
    final online = await ref.read(connectivityServiceProvider).isOnline();
    // Guests are always graded on-device; so is anyone currently offline. Only
    // an online, signed-in user takes the authoritative server path.
    final wantLocal = isGuest(ref) || !online;

    if (wantLocal) {
      if (!context.mounted) return;
      setState(() => _isStarting = false);
      if (!downloaded) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              online
                  ? 'Download this exam first to take it as a guest.'
                  : "You're offline. Download this exam first to take it.",
            ),
          ),
        );
        return;
      }
      _goToLocalExam(context);
      return;
    }

    final result = await ref.read(startExamUseCaseProvider)(examId);
    if (!context.mounted) return;
    setState(() => _isStarting = false);

    result.fold(
      (failure) {
        // Started online but the network dropped — grade the downloaded copy.
        if (isOfflineFailure(failure) && downloaded) {
          _goToLocalExam(context);
          return;
        }
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(failure.message)),
        );
      },
      (attempt) {
        ref.read(pendingExamAttemptProvider.notifier).state = attempt;
        ref.invalidate(myAllExamAttemptsProvider);
        ref.invalidate(examAttemptsByExamProvider(examId));
        context.push(
          AppRoutes.examAttemptPath(
            attemptId: attempt.attemptId,
            examId: attempt.examId,
          ),
        );
      },
    );
  }

  /// Offline / guest branch: synthesize the attempt from the downloaded copy and
  /// hand it to the notifier's adopt-path via [pendingExamAttemptProvider] — the
  /// same seam the server flow uses — then navigate with the `local-…` id.
  void _goToLocalExam(BuildContext context) {
    final exam = ref.read(downloadStoreProvider).getOfflineExam(widget.exam.id);
    if (exam == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('This exam is not downloaded.')),
      );
      return;
    }
    final localId = newLocalAttemptId();
    ref.read(pendingExamAttemptProvider.notifier).state =
        buildLocalExamAttempt(exam, attemptId: localId);
    context.push(
      AppRoutes.examAttemptPath(attemptId: localId, examId: widget.exam.id),
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
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: const TextStyle(fontSize: 13, color: AppColors.textMuted),
          ),
          Text(
            value,
            style: const TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
        ],
      ),
    );
  }
}
