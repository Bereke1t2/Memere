import 'package:flutter/material.dart';

import '../../../../core/constants/app_text_styles.dart';

class ExamTimerBadge extends StatelessWidget {
  const ExamTimerBadge({
    super.key,
    required this.seconds,
  });

  final int? seconds;

  @override
  Widget build(BuildContext context) {
    if (seconds == null) {
      return const SizedBox.shrink();
    }
    final value = seconds!;
    final critical = value <= 60;
    final warning = value <= 300;

    final color = critical
        ? const Color(0xFFDD5F55)
        : warning
            ? const Color(0xFFD99A2B)
            : const Color(0xFF35A862);

    final bgTint = critical
        ? const Color(0x28DD5F55)
        : warning
            ? const Color(0x28D99A2B)
            : const Color(0x2035A862);

    final borderColor = critical
        ? const Color(0xFFDD5F55).withAlpha(120)
        : warning
            ? const Color(0xFFD99A2B).withAlpha(120)
            : const Color(0xFF35A862).withAlpha(100);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: bgTint,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: borderColor, width: 1.2),
        boxShadow: critical || warning
            ? [
                BoxShadow(
                  color: color.withAlpha(50),
                  blurRadius: 8,
                  offset: const Offset(0, 2),
                ),
              ]
            : null,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            critical
                ? Icons.alarm_rounded
                : Icons.timer_outlined,
            size: 15,
            color: color,
          ),
          const SizedBox(width: 5),
          Text(
            formatExamTimer(value),
            style: AppTextStyles.labelMedium.copyWith(
              color: color,
              fontWeight: FontWeight.w700,
              letterSpacing: 0.5,
            ),
          ),
        ],
      ),
    );
  }
}

/// Formats seconds as `h:mm:ss` (or `m:ss` under an hour).
String formatExamTimer(int seconds) {
  final safe = seconds < 0 ? 0 : seconds;
  final hours = safe ~/ 3600;
  final minutes = (safe % 3600) ~/ 60;
  final secs = safe % 60;
  if (hours > 0) {
    return '$hours:${minutes.toString().padLeft(2, '0')}:'
        '${secs.toString().padLeft(2, '0')}';
  }
  return '$minutes:${secs.toString().padLeft(2, '0')}';
}
