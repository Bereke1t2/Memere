import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../domain/entities/quiz_result_entity.dart';

class QuizScoreSummary extends StatelessWidget {
  const QuizScoreSummary({
    super.key,
    required this.result,
  });

  final QuizResultEntity result;

  @override
  Widget build(BuildContext context) {
    const emeraldColor = Color(0xFF10B981);
    const amberColor = Color(0xFFF59E0B);
    final isPassed = result.passed;
    final primaryColor = isPassed ? emeraldColor : amberColor;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(22),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(
          color: primaryColor.withValues(alpha: 0.35),
          width: 1.5,
        ),
        boxShadow: [
          BoxShadow(
            color: primaryColor.withValues(alpha: 0.08),
            blurRadius: 16,
            offset: const Offset(0, 6),
          ),
        ],
      ),
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
            decoration: BoxDecoration(
              color: primaryColor.withValues(alpha: 0.15),
              borderRadius: BorderRadius.circular(20),
              border: Border.all(color: primaryColor.withValues(alpha: 0.4)),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  isPassed
                      ? Icons.emoji_events_rounded
                      : Icons.psychology_rounded,
                  color: primaryColor,
                  size: 18,
                ),
                const SizedBox(width: 8),
                Text(
                  isPassed ? 'QUIZ PASSED 🎉' : 'KEEP PRACTICING 💪',
                  style: TextStyle(
                    fontSize: 13,
                    fontWeight: FontWeight.w800,
                    color: primaryColor,
                    letterSpacing: 0.5,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 18),
          Text(
            '${result.percentage.toStringAsFixed(0)}%',
            style: TextStyle(
              fontSize: 48,
              fontWeight: FontWeight.w900,
              color: primaryColor,
              height: 1.0,
            ),
          ),
          const SizedBox(height: 6),
          Text(
            '${_trimNumber(result.score)} of ${result.totalPoints} Total Points',
            style: const TextStyle(
              fontSize: 14.5,
              fontWeight: FontWeight.w600,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 14),
          Text(
            isPassed
                ? 'Great effort! Review the detailed question feedback below to reinforce your understanding.'
                : 'Don\'t give up! Review the correct answers and explanations below, then try again.',
            textAlign: TextAlign.center,
            style: const TextStyle(
              fontSize: 13,
              color: Color(0xFF94A3B8),
              height: 1.4,
            ),
          ),
        ],
      ),
    );
  }
}

String _trimNumber(double value) {
  return value % 1 == 0 ? value.toInt().toString() : value.toStringAsFixed(1);
}
