import 'package:flutter/material.dart';

/// Memere Professional Color System — Duolingo-Inspired (Strict 4-Color Palette).
///
/// 1. Brand Emerald Green (`#10B981`): Primary actions, level progress, completion indicators.
/// 2. Brand Warm Amber (`#F59E0B`): Streaks 🔥, XP points ⚡, daily targets.
/// 3. Background Canvas (`#0B0F17`): Deep slate dark canvas.
/// 4. Card Surface (`#1E293B`): Tactile cards with 3D depth borders (`#334155`).
abstract class AppColors {
  // ── 4 Master Theme Colors ────────────────────────────────────────────────
  static const Color brandEmerald = Color(0xFF10B981);
  static const Color brandEmeraldDark = Color(0xFF059669);
  static const Color brandAmber = Color(0xFFF59E0B);
  static const Color brandAmberDark = Color(0xFFD97706);

  // ── Backgrounds ──────────────────────────────────────────────────────────
  static const Color bgPrimary = Color(0xFF0B0F17);
  static const Color bgSecondary = Color(0xFF1E293B);
  static const Color bgTertiary = Color(0xFF131C2E);
  static const Color bgQuaternary = Color(0xFF334155);
  static const Color bgOverlay = Color(0xB8000000);

  // ── Accent ───────────────────────────────────────────────────────────────
  static const Color accentPrimary = Color(0xFF10B981);
  static const Color accentPrimaryDeep = Color(0xFF059669);
  static const Color accentSecondary = Color(0xFFF59E0B);
  static const Color accentTertiary = Color(0xFF94A3B8);
  static const Color accentGlow = Color(0x2610B981);

  // ── Text ─────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFFF8FAFC);
  static const Color textSecondary = Color(0xFFCBD5E1);
  static const Color textMuted = Color(0xFF94A3B8);
  static const Color textDisabled = Color(0xFF64748B);
  static const Color textInverse = Color(0xFF0B0F17);

  // ── Status ───────────────────────────────────────────────────────────────
  static const Color success = Color(0xFF10B981);
  static const Color successSurface = Color(0x2210B981);
  static const Color error = Color(0xFFEF4444);
  static const Color errorSurface = Color(0x22EF4444);
  static const Color warning = Color(0xFFF59E0B);
  static const Color warningSurface = Color(0x22F59E0B);
  static const Color info = Color(0xFF38BDF8);
  static const Color infoSurface = Color(0x2238BDF8);
  static const Color lavender = Color(0xFF10B981);
  static const Color lavenderSurface = Color(0x2210B981);
  static const Color coral = Color(0xFFF59E0B);
  static const Color coralSurface = Color(0x22F59E0B);

  // ── Subject Tag Colors (Unified under Brand Palette) ────────────────────
  static const Color subjectNeutral = Color(0xFF94A3B8);
  static const Color subjectMath = Color(0xFF10B981);
  static const Color subjectPhysics = Color(0xFF10B981);
  static const Color subjectChem = Color(0xFF10B981);
  static const Color subjectBio = Color(0xFF10B981);
  static const Color subjectEng = Color(0xFFF59E0B);
  static const Color subjectHist = Color(0xFFF59E0B);
  static const Color subjectGeo = Color(0xFF10B981);
  static const Color subjectEcon = Color(0xFFF59E0B);
  static const Color subjectCivics = Color(0xFF10B981);

  // ── Borders & Dividers (Tactile 3D Depth) ─────────────────────────────────
  static const Color border = Color(0xFF1E293B);
  static const Color borderStrong = Color(0xFF334155);
  static const Color borderFocused = Color(0xFF10B981);
  static const Color divider = Color(0xFF1E293B);

  // ── Gradients ────────────────────────────────────────────────────────────
  static const LinearGradient backgroundGradient = LinearGradient(
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
    colors: [Color(0xFF0F172A), Color(0xFF0B0F17)],
  );

  static const LinearGradient primaryGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF10B981), Color(0xFF059669)],
  );

  static const LinearGradient headerGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF1E293B), Color(0xFF0F172A)],
  );

  static const LinearGradient cardGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF1E293B), Color(0xFF131C2E)],
  );

  static const LinearGradient warmGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFFF59E0B), Color(0xFFD97706)],
  );

  static const LinearGradient examGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF1E293B), Color(0xFF0B0F17)],
  );

  static const RadialGradient glowGradient = RadialGradient(
    center: Alignment.center,
    radius: 0.8,
    colors: [Color(0x2410B981), Color(0x0010B981)],
  );

  // ── Refinements ─────────────────────────────────────────────────────────
  static const Color scrim = Color(0xCC000000);
  static const Color hairline = Color(0x33334155);
  static const Color pressedOverlay = Color(0x1F10B981);
  static const Color skeletonBase = Color(0xFF131C2E);
  static const Color skeletonHighlight = Color(0xFF1E293B);
}
