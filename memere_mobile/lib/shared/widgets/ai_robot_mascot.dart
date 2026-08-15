import 'package:flutter/material.dart';
import '../../core/constants/app_colors.dart';

/// Reusable AI Robot Persona Mascot widget for Memere.
///
/// Can be rendered as:
/// - Avatar / Icon Tile
/// - Card Mascot
/// - Floating Action Button (FAB) for instant AI Concept Tutor guidance.
class AiRobotMascot extends StatelessWidget {
  const AiRobotMascot({
    super.key,
    this.size = 40,
    this.backgroundColor = AppColors.brandEmerald,
    this.iconColor = Colors.white,
  });

  final double size;
  final Color backgroundColor;
  final Color iconColor;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: backgroundColor.withAlpha(40),
        shape: BoxShape.circle,
        border: Border.all(color: backgroundColor.withAlpha(120), width: 1.5),
        boxShadow: [
          BoxShadow(
            color: backgroundColor.withAlpha(50),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Icon(
        Icons.smart_toy_rounded,
        size: size * 0.55,
        color: backgroundColor,
      ),
    );
  }
}

/// Floating Action Button (FAB) featuring the AI Robot Persona Mascot
class AiTutorFab extends StatelessWidget {
  const AiTutorFab({
    super.key,
    required this.onPressed,
  });

  final VoidCallback onPressed;

  @override
  Widget build(BuildContext context) {
    return FloatingActionButton(
      onPressed: onPressed,
      backgroundColor: AppColors.brandEmerald,
      foregroundColor: Colors.white,
      elevation: 6,
      shape: const CircleBorder(),
      child: const Icon(
        Icons.smart_toy_rounded,
        size: 26,
      ),
    );
  }
}
