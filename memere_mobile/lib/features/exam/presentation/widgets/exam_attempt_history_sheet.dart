import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/router/app_router.dart';
import '../../domain/entities/exam_attempt_history_entity.dart';
import '../providers/exam_providers.dart';

/// Shows a bottom sheet with the student's attempt history for an exam.
void showExamAttemptHistorySheet({
  required BuildContext context,
  required String examId,
  required String examTitle,
  required VoidCallback onRetake,
}) {
  showModalBottomSheet(
    context: context,
    isScrollControlled: true,
    backgroundColor: const Color(0xFF0F0F14),
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      side: BorderSide(color: Color(0xFF22222C), width: 1),
    ),
    builder: (sheetContext) => DraggableScrollableSheet(
      initialChildSize: 0.65,
      minChildSize: 0.40,
      maxChildSize: 0.90,
      expand: false,
      builder: (_, scrollController) => _HistorySheetContent(
        scrollController: scrollController,
        examId: examId,
        examTitle: examTitle,
        onRetake: () {
          Navigator.of(sheetContext).pop();
          onRetake();
        },
      ),
    ),
  );
}

class _HistorySheetContent extends ConsumerWidget {
  const _HistorySheetContent({
    required this.scrollController,
    required this.examId,
    required this.examTitle,
    required this.onRetake,
  });

  final ScrollController scrollController;
  final String examId;
  final String examTitle;
  final VoidCallback onRetake;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final attemptsAsync = ref.watch(examAttemptsByExamProvider(examId));

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 18),
      child: Column(
        children: [
          const SizedBox(height: 12),
          // Drag handle
          Center(
            child: Container(
              width: 38,
              height: 4,
              decoration: BoxDecoration(
                color: const Color(0xFF3F3F46),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 16),
          // Title
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Exam Attempt History',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      examTitle,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                        color: Color(0xFF71717A),
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ),
              IconButton(
                onPressed: () => Navigator.of(context).pop(),
                icon: const Icon(Icons.close_rounded, color: Color(0xFF71717A)),
              ),
            ],
          ),
          const SizedBox(height: 14),

          // Attempts List
          Expanded(
            child: attemptsAsync.when(
              loading: () => const Center(
                child: CircularProgressIndicator(color: Color(0xFF10B981), strokeWidth: 2),
              ),
              error: (err, _) => Center(
                child: Text('Could not load history: $err', style: const TextStyle(color: Colors.white)),
              ),
              data: (attempts) {
                if (attempts.isEmpty) {
                  return const Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.history_toggle_off_rounded, size: 40, color: Color(0xFF52525B)),
                        SizedBox(height: 10),
                        Text(
                          'No previous attempts found.',
                          style: TextStyle(color: Color(0xFFA1A1AA), fontSize: 14),
                        ),
                      ],
                    ),
                  );
                }

                return ListView.builder(
                  controller: scrollController,
                  itemCount: attempts.length,
                  itemBuilder: (context, index) {
                    final attempt = attempts[index];
                    final attemptNumber = attempts.length - index;
                    return _AttemptHistoryCard(
                      attempt: attempt,
                      attemptNumber: attemptNumber,
                    );
                  },
                );
              },
            ),
          ),

          // Retake CTA at bottom
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 14),
            child: SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton.icon(
                onPressed: onRetake,
                icon: const Icon(Icons.replay_rounded, size: 18),
                label: const Text(
                  'Retake This Exam',
                  style: TextStyle(fontWeight: FontWeight.w700, fontSize: 15),
                ),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF10B981),
                  foregroundColor: Colors.black,
                  elevation: 0,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _AttemptHistoryCard extends StatelessWidget {
  const _AttemptHistoryCard({
    required this.attempt,
    required this.attemptNumber,
  });

  final ExamAttemptHistoryEntity attempt;
  final int attemptNumber;

  @override
  Widget build(BuildContext context) {
    final isGraded = attempt.isGraded;
    final scorePct = attempt.percentage ?? 0.0;
    final isPassed = attempt.isPassed;

    Color badgeColor = isPassed ? const Color(0xFF10B981) : const Color(0xFFF59E0B);
    String statusText = isPassed ? 'Passed' : 'Needs Practice';

    if (!isGraded) {
      badgeColor = const Color(0xFF38BDF8);
      statusText = 'In Progress';
    }

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: const Color(0xFF14141C),
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: const Color(0xFF22222E)),
      ),
      child: Row(
        children: [
          // Score Circle or In-Progress Icon
          Container(
            width: 44,
            height: 44,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: badgeColor.withAlpha(25),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: badgeColor.withAlpha(80)),
            ),
            child: isGraded
                ? Text(
                    '${scorePct.toStringAsFixed(0)}%',
                    style: TextStyle(
                      color: badgeColor,
                      fontWeight: FontWeight.w800,
                      fontSize: 13,
                    ),
                  )
                : Icon(Icons.timer_outlined, color: badgeColor, size: 20),
          ),
          const SizedBox(width: 12),

          // Details
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Text(
                      'Attempt #$attemptNumber',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.w700,
                        fontSize: 14,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: badgeColor.withAlpha(25),
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: Text(
                        statusText,
                        style: TextStyle(
                          color: badgeColor,
                          fontSize: 10,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 3),
                Text(
                  _formatDate(attempt.submittedAt ?? attempt.startedAt),
                  style: const TextStyle(color: Color(0xFF71717A), fontSize: 11),
                ),
              ],
            ),
          ),

          // View Results Button
          if (isGraded)
            OutlinedButton(
              onPressed: () {
                Navigator.of(context).pop();
                context.push(AppRoutes.examResultPath(attempt.attemptId));
              },
              style: OutlinedButton.styleFrom(
                foregroundColor: const Color(0xFF38BDF8),
                side: const BorderSide(color: Color(0xFF263342)),
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                minimumSize: Size.zero,
                tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
              ),
              child: const Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text('Review', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
                  SizedBox(width: 2),
                  Icon(Icons.arrow_forward_ios_rounded, size: 10),
                ],
              ),
            ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')} ${date.hour.toString().padLeft(2, '0')}:${date.minute.toString().padLeft(2, '0')}';
  }
}
