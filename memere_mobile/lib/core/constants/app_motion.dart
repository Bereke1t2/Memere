import 'package:flutter/animation.dart';

/// Shared motion tokens for transitions, feedback, and reveals.
abstract class AppMotion {
  static const Duration fast = Duration(milliseconds: 120);
  static const Duration base = Duration(milliseconds: 220);
  static const Duration slow = Duration(milliseconds: 320);
  static const Duration stagger = Duration(milliseconds: 40);

  static const Curve standard = Cubic(0.2, 0.0, 0.0, 1.0);
  static const Curve exit = Cubic(0.3, 0.0, 0.8, 0.15);
  static const Curve emphasized = Curves.easeOutBack;
}
