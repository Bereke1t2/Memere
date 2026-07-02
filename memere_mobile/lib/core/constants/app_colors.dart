import 'package:flutter/material.dart';

/// Memere Color System — CPD Hub-inspired monochrome.
///
/// Black canvas, white primary accent, graphite surfaces, and quiet gray text.
/// Subject/status colors are deliberately desaturated so the UI stays black and
/// white first.
abstract class AppColors {
  // ── Backgrounds ──────────────────────────────────────────────────────────
  static const Color bgPrimary = Color(0xFF050505);
  static const Color bgSecondary = Color(0xFF101010);
  static const Color bgTertiary = Color(0xFF181818);
  static const Color bgQuaternary = Color(0xFF242424);
  static const Color bgOverlay = Color(0xB8000000);

  // ── Accent ───────────────────────────────────────────────────────────────
  static const Color accentPrimary = Color(0xFFFFFFFF);
  static const Color accentPrimaryDeep = Color(0xFFBDBDBD);
  static const Color accentSecondary = Color(0xFFE0E0E0);
  static const Color accentTertiary = Color(0xFF9E9E9E);
  static const Color accentGlow = Color(0x26FFFFFF);

  // ── Text ─────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFFFFFFFF);
  static const Color textSecondary = Color(0xFFB8B8B8);
  static const Color textMuted = Color(0xFF8A8A8A);
  static const Color textDisabled = Color(0xFF5E5E5E);
  static const Color textInverse = Color(0xFF050505);

  // ── Status ───────────────────────────────────────────────────────────────
  static const Color success = Color(0xFFD7D7D7);
  static const Color successSurface = Color(0x22D7D7D7);
  static const Color error = Color(0xFFBDBDBD);
  static const Color errorSurface = Color(0x22BDBDBD);
  static const Color warning = Color(0xFFCFCFCF);
  static const Color warningSurface = Color(0x22CFCFCF);
  static const Color info = Color(0xFFB5B5B5);
  static const Color infoSurface = Color(0x22B5B5B5);
  static const Color lavender = Color(0xFFA7A7A7);
  static const Color lavenderSurface = Color(0x22A7A7A7);
  static const Color coral = Color(0xFF969696);
  static const Color coralSurface = Color(0x22969696);

  // ── Subject Tag Color ─────────────────────────────────────────────────────
  static const Color subjectNeutral = Color(0xFF9A9A9A);
  static const Color subjectMath = Color(0xFFE6E6E6);
  static const Color subjectPhysics = Color(0xFFD6D6D6);
  static const Color subjectChem = Color(0xFFC8C8C8);
  static const Color subjectBio = Color(0xFFBABABA);
  static const Color subjectEng = Color(0xFFACACAC);
  static const Color subjectHist = Color(0xFF9E9E9E);
  static const Color subjectGeo = Color(0xFF909090);
  static const Color subjectEcon = Color(0xFF828282);
  static const Color subjectCivics = Color(0xFF747474);

  // ── Borders & Dividers ───────────────────────────────────────────────────
  static const Color border = Color(0xFF2A2A2A);
  static const Color borderStrong = Color(0xFF3A3A3A);
  static const Color borderFocused = Color(0xFFFFFFFF);
  static const Color divider = Color(0xFF202020);

  // ── Gradients ────────────────────────────────────────────────────────────
  static const LinearGradient backgroundGradient = LinearGradient(
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
    colors: [Color(0xFF080808), Color(0xFF050505)],
  );

  static const LinearGradient primaryGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFFFFFFFF), Color(0xFFBEBEBE)],
  );

  static const LinearGradient headerGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF1E1E1E), Color(0xFF101010)],
  );

  static const LinearGradient cardGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF191919), Color(0xFF101010)],
  );

  static const LinearGradient warmGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF2E2E2E), Color(0xFF171717)],
  );

  static const LinearGradient examGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF2A2A2A), Color(0xFF111111)],
  );

  static const RadialGradient glowGradient = RadialGradient(
    center: Alignment.center,
    radius: 0.8,
    colors: [Color(0x24FFFFFF), Color(0x00FFFFFF)],
  );

  // ── Refinements ─────────────────────────────────────────────────────────
  static const Color scrim = Color(0xCC000000);
  static const Color hairline = Color(0x333A3A3A);
  static const Color pressedOverlay = Color(0x1FFFFFFF);
  static const Color skeletonBase = Color(0xFF141414);
  static const Color skeletonHighlight = Color(0xFF252525);
}
