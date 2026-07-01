import 'package:flutter/material.dart';

/// Memere Color System — dark-first, natural and focused.
///
/// Design rules:
/// - Green stays the brand accent because it feels connected to learning.
/// - Amber and blue are supporting colors for progress, exams, and highlights.
/// - Surfaces use warm neutral charcoal with a small green cast, never flat black.
/// - Status colors are muted and only used for meaning.
abstract class AppColors {
  // ── Backgrounds ──────────────────────────────────────────────────────────
  static const Color bgPrimary = Color(0xFF08110E); // main background
  static const Color bgSecondary = Color(0xFF101A16); // card, bottom sheet
  static const Color bgTertiary = Color(0xFF17241F); // elevated surface, input
  static const Color bgQuaternary = Color(0xFF20312A); // selected controls
  static const Color bgOverlay = Color(0x99000000); // modal overlay 60%

  // ── Accent ───────────────────────────────────────────────────────────────
  static const Color accentPrimary = Color(0xFF35B87E); // primary CTA
  static const Color accentPrimaryDeep = Color(0xFF147852);
  static const Color accentSecondary = Color(0xFFF2B85B); // warm emphasis
  static const Color accentTertiary = Color(0xFF6E9FE8); // exam / analytics
  static const Color accentGlow = Color(0x3335B87E); // green glow overlay

  // ── Text ─────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFFF3F7F4);
  static const Color textSecondary = Color(0xFFA9B9B0);
  static const Color textMuted = Color(0xFF83948B);
  static const Color textDisabled = Color(0xFF64746C);
  static const Color textInverse =
      Color(0xFF06100C); // text on bright accent bg

  // ── Status (muted) ───────────────────────────────────────────────────────
  static const Color success = Color(0xFF55D494);
  static const Color successSurface = Color(0x1A4FC78F);
  static const Color error = Color(0xFFEF6A65);
  static const Color errorSurface = Color(0x1AEF6A65);
  static const Color warning = Color(0xFFF2B85B);
  static const Color warningSurface = Color(0x1AF0B45B);
  static const Color info = Color(0xFF74A8E8);
  static const Color infoSurface = Color(0x1A70A7D8);
  static const Color lavender = Color(0xFFB69BF2);
  static const Color lavenderSurface = Color(0x1AB69BF2);
  static const Color coral = Color(0xFFFF8A7A);
  static const Color coralSurface = Color(0x1AFF8A7A);

  // ── Subject Tag Color ─────────────────────────────────────────────────────
  /// Subject colors are muted enough for repeated cards but distinct enough
  /// to help recognition while browsing.
  static const Color subjectNeutral = Color(0xFF9BAB9F);
  static const Color subjectMath = Color(0xFF74A8E8);
  static const Color subjectPhysics = Color(0xFFB69BF2);
  static const Color subjectChem = Color(0xFF55D494);
  static const Color subjectBio = Color(0xFF7ACF73);
  static const Color subjectEng = Color(0xFFFF8A7A);
  static const Color subjectHist = Color(0xFFF2B85B);
  static const Color subjectGeo = Color(0xFF68C7C0);
  static const Color subjectEcon = Color(0xFFE89CC8);

  // ── Borders & Dividers ───────────────────────────────────────────────────
  static const Color border = Color(0xFF263B33);
  static const Color borderStrong = Color(0xFF365B4A);
  static const Color borderFocused = Color(0xFF55D494);
  static const Color divider = Color(0xFF17261F);

  // ── Gradients ────────────────────────────────────────────────────────────
  static const LinearGradient backgroundGradient = LinearGradient(
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
    colors: [Color(0xFF0E2119), Color(0xFF08110E)],
  );

  static const LinearGradient primaryGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF55D494), Color(0xFF1E8A60)],
  );

  static const LinearGradient cardGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF152820), Color(0xFF101A16)],
  );

  static const LinearGradient warmGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFFF2B85B), Color(0xFFD87548)],
  );

  static const LinearGradient examGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF74A8E8), Color(0xFF6F7CE4)],
  );

  static const RadialGradient glowGradient = RadialGradient(
    center: Alignment.center,
    radius: 0.8,
    colors: [Color(0x4035B87E), Color(0x0035B87E)],
  );

  // ── Refinements ─────────────────────────────────────────────────────────
  static const Color scrim = Color(0xCC06100C);
  static const Color hairline = Color(0x1FFFFFFF);
  static const Color pressedOverlay = Color(0x17FFFFFF);
  static const Color skeletonBase = Color(0xFF14231D);
  static const Color skeletonHighlight = Color(0xFF263C33);
}
