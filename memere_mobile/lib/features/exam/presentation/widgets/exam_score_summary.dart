import 'package:flutter/material.dart';

import '../../domain/entities/exam_result_entity.dart';

class ExamScoreSummary extends StatelessWidget {
  const ExamScoreSummary({
    super.key,
    required this.result,
  });

  final ExamResultEntity result;

  @override
  Widget build(BuildContext context) {
    final passed = result.passed;
    final color = passed ? const Color(0xFF10B981) : const Color(0xFFF59E0B);
    final gradientColors = passed
        ? [const Color(0xFF064E3B), const Color(0xFF0F1F18)]
        : [const Color(0xFF451A03), const Color(0xFF1F150F)];

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: gradientColors,
        ),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: color.withAlpha(80),
          width: 1.5,
        ),
        boxShadow: [
          BoxShadow(
            color: color.withAlpha(30),
            blurRadius: 16,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: color.withAlpha(40),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  passed ? Icons.emoji_events_rounded : Icons.psychology_rounded,
                  color: color,
                  size: 24,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      passed ? 'Congratulations!' : 'Keep Practicing!',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      passed
                          ? 'You successfully passed the examination.'
                          : 'Review weak topics to improve your score next time.',
                      style: const TextStyle(
                        color: Color(0xFFD4D4D8),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '${result.percentage.toStringAsFixed(0)}%',
                style: TextStyle(
                  color: color,
                  fontSize: 48,
                  fontWeight: FontWeight.w900,
                  letterSpacing: -1.5,
                  height: 1.0,
                ),
              ),
              const SizedBox(width: 16),
              Padding(
                padding: const EdgeInsets.only(bottom: 6),
                child: Text(
                  '${_trimNumber(result.score)} / ${result.totalMarks} Marks',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w700,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 10),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
            decoration: BoxDecoration(
              color: Colors.black.withAlpha(60),
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: Colors.white.withAlpha(20)),
            ),
            child: Text(
              'Required to pass: ${result.passMarks} marks (50%)',
              style: const TextStyle(
                color: Color(0xFFA1A1AA),
                fontSize: 12,
                fontWeight: FontWeight.w500,
              ),
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
