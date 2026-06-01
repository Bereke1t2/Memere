import 'package:flutter/material.dart';

/// ExamPrep Color System — Dark First
/// All colors reference the dark theme. Light theme values noted where different.
abstract class AppColors {
  // ── Backgrounds ──────────────────────────────────────────────────────────
  static const Color bgPrimary   = Color(0xFF0D0D0D); // main background
  static const Color bgSecondary = Color(0xFF1A1A1A); // card, bottom sheet
  static const Color bgTertiary  = Color(0xFF252525); // elevated surface, input
  static const Color bgOverlay   = Color(0x80000000); // modal overlay 50%

  // ── Accent ───────────────────────────────────────────────────────────────
  static const Color accentPrimary   = Color(0xFF6C63FF); // primary CTA
  static const Color accentSecondary = Color(0xFF03DAC6); // teal — secondary CTA
  static const Color accentGlow      = Color(0x336C63FF); // purple glow overlay

  // ── Text ─────────────────────────────────────────────────────────────────
  static const Color textPrimary   = Color(0xFFF5F5F5);
  static const Color textSecondary = Color(0xFF9E9E9E);
  static const Color textDisabled  = Color(0xFF4A4A4A);
  static const Color textInverse   = Color(0xFF0D0D0D); // text on accent bg

  // ── Status ───────────────────────────────────────────────────────────────
  static const Color success        = Color(0xFF4CAF50);
  static const Color successSurface = Color(0x1A4CAF50);
  static const Color error          = Color(0xFFCF6679);
  static const Color errorSurface   = Color(0x1ACF6679);
  static const Color warning        = Color(0xFFFFB347);
  static const Color warningSurface = Color(0x1AFFB347);
  static const Color info           = Color(0xFF64B5F6);
  static const Color infoSurface    = Color(0x1A64B5F6);

  // ── Subject Tag Colors ────────────────────────────────────────────────────
  /// Each subject gets a consistent color for tags/pills
  static const Color subjectMath    = Color(0xFF6C63FF);
  static const Color subjectPhysics = Color(0xFF4FC3F7);
  static const Color subjectChem    = Color(0xFF81C784);
  static const Color subjectBio     = Color(0xFF4DB6AC);
  static const Color subjectEng     = Color(0xFFFFB74D);
  static const Color subjectHist    = Color(0xFFBA68C8);
  static const Color subjectGeo     = Color(0xFF4DD0E1);
  static const Color subjectEcon    = Color(0xFF9CCC65);

  // ── Borders & Dividers ───────────────────────────────────────────────────
  static const Color border         = Color(0xFF2C2C2C);
  static const Color borderFocused  = Color(0xFF6C63FF);
  static const Color divider        = Color(0xFF1E1E1E);

  // ── Gradients ────────────────────────────────────────────────────────────
  static const LinearGradient primaryGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF6C63FF), Color(0xFF9C5CFF)],
  );

  static const LinearGradient cardGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF1E1E2E), Color(0xFF252535)],
  );

  static const RadialGradient glowGradient = RadialGradient(
    center: Alignment.center,
    radius: 0.8,
    colors: [Color(0x406C63FF), Color(0x006C63FF)],
  );
}
