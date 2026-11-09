import 'package:flutter/material.dart';

/// Memere colour system — an obsidian-dark canvas carrying a single leaf-green
/// accent, scoped tightly so the accent means something when it appears.
///
/// ## Why these values
///
/// The palette was previously built from stock framework defaults — emerald
/// `#10B981`, amber `#F59E0B`, violet `#8B5CF6`, red `#EF4444`, sky `#38BDF8`,
/// zinc greys — on a near-pure `#050505` black. Those are the swatches every
/// generated dark theme reaches for, which is exactly why the app read as
/// generic rather than as itself.
///
/// Each hue was re-chosen for an Ethiopian ESSLCE exam-prep app that students
/// read for hours, often late at night, on a phone:
///
///  * **Accent** moved from a cyan-leaning mint (hue 160°, 84% saturation) to a
///    leaf green (hue 144°, 52%). It reads as study and growth rather than as a
///    dashboard status light, and desaturating it stops it vibrating against the
///    dark ground during long sessions.
///  * **Warning** moved from neon amber to an ochre gold — warmer, closer to
///    manuscript and chalk than to a hazard label.
///  * **Error** moved from fire-engine red to a clay red, so a wrong answer
///    reads as a correction and not an alarm.
///  * **Info** moved from an electric sky blue to a muted steel blue that sits
///    behind the accent instead of competing with it.
///  * **Neutrals** carry a barely-perceptible green bias rather than the old
///    blue-grey zinc, so greys read as chosen and belong to the accent's family.
///  * **The canvas** lifted off pure black. `#050505` is harsh against white
///    text and smears on OLED panels while scrolling.
///
/// ## Contrast
///
/// Every colour used for text or an icon clears WCAG AA (4.5:1) against all
/// three dark surfaces. Values intended as fills, borders or gradient stops are
/// exempt by role and annotated where that matters. Re-verify with a contrast
/// checker before changing any value here — several of these clear AA by a
/// small margin and a casual "brighten it a bit" can quietly break legibility.
abstract class AppColors {
  // ── 3 Level Badge Colors (Strict Scoped Usage) ───────────────────────────
  static const Color levelBeginner = Color(0xFF35A862);
  static const Color levelBeginnerDark = Color(0xFF24814A);
  static const Color levelIntermediate = Color(0xFFD99A2B);
  static const Color levelIntermediateDark = Color(0xFFB07A1E);
  static const Color levelAdvanced = Color(0xFF8A79C9);
  static const Color levelAdvancedDark = Color(0xFF6F5FB0);

  // ── Master Brand Tokens ──────────────────────────────────────────────────
  // Named "emerald"/"amber" for continuity with existing call sites; the values
  // are now leaf green and ochre gold.
  static const Color brandEmerald = Color(0xFF35A862);
  static const Color brandEmeraldDark = Color(0xFF24814A);
  static const Color brandAmber = Color(0xFFD99A2B);
  static const Color brandAmberDark = Color(0xFFB07A1E);

  // ── Backgrounds (Obsidian, lifted off pure black) ────────────────────────
  static const Color bgPrimary = Color(0xFF080A09);
  static const Color bgSecondary = Color(0xFF101311);
  static const Color bgTertiary = Color(0xFF171B18);
  static const Color bgQuaternary = Color(0xFF232825);
  static const Color bgOverlay = Color(0xCC000000);

  // ── Accent & Neutral Surfaces ────────────────────────────────────────────
  static const Color accentPrimary = Color(0xFF35A862);
  static const Color accentPrimaryDeep = Color(0xFF24814A);
  static const Color accentSecondary = Color(0xFFD99A2B);
  static const Color accentTertiary = Color(0xFF9AA09A);
  static const Color accentGlow = Color(0x2635A862);

  // ── Text ─────────────────────────────────────────────────────────────────
  static const Color textPrimary = Color(0xFFFFFFFF);
  static const Color textSecondary = Color(0xFFD7D9D4);
  static const Color textMuted = Color(0xFF9AA09A);
  /// Decorative only — below AA by design, never use for content a student
  /// needs to read.
  static const Color textDisabled = Color(0xFF575D58);
  static const Color textInverse = Color(0xFF080A09);

  // ── Status ───────────────────────────────────────────────────────────────
  static const Color success = Color(0xFF35A862);
  static const Color successSurface = Color(0x2235A862);
  static const Color error = Color(0xFFDD5F55);
  static const Color errorSurface = Color(0x22DD5F55);
  static const Color warning = Color(0xFFD99A2B);
  static const Color warningSurface = Color(0x22D99A2B);
  static const Color info = Color(0xFF5FA8D3);
  static const Color infoSurface = Color(0x225FA8D3);
  static const Color lavender = Color(0xFF8A79C9);
  static const Color lavenderSurface = Color(0x228A79C9);
  static const Color coral = Color(0xFFD99A2B);
  static const Color coralSurface = Color(0x22D99A2B);

  // ── Subject Tag Colors (neutral by design — subjects are not colour-coded,
  //    so the accent stays reserved for state and actions) ─────────────────
  static const Color subjectNeutral = Color(0xFF7E847C);
  static const Color subjectMath = Color(0xFF7E847C);
  static const Color subjectPhysics = Color(0xFF7E847C);
  static const Color subjectChem = Color(0xFF7E847C);
  static const Color subjectBio = Color(0xFF7E847C);
  static const Color subjectEng = Color(0xFF7E847C);
  static const Color subjectHist = Color(0xFF7E847C);
  static const Color subjectGeo = Color(0xFF7E847C);
  static const Color subjectEcon = Color(0xFF7E847C);
  static const Color subjectCivics = Color(0xFF7E847C);

  // ── Borders & Dividers (fills, not text — AA does not apply) ─────────────
  static const Color border = Color(0xFF191D1A);
  static const Color borderStrong = Color(0xFF2A302C);
  static const Color borderFocused = Color(0xFF515A53);
  static const Color divider = Color(0xFF191D1A);

  // ── Gradients ────────────────────────────────────────────────────────────
  static const LinearGradient backgroundGradient = LinearGradient(
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
    colors: [Color(0xFF0C0F0D), Color(0xFF080A09)],
  );

  static const LinearGradient primaryGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF35A862), Color(0xFF24814A)],
  );

  static const LinearGradient headerGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF101311), Color(0xFF080A09)],
  );

  static const LinearGradient cardGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF101311), Color(0xFF0B0E0C)],
  );

  static const LinearGradient warmGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFFD99A2B), Color(0xFFB07A1E)],
  );

  static const LinearGradient examGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF101311), Color(0xFF080A09)],
  );

  static const RadialGradient glowGradient = RadialGradient(
    center: Alignment.center,
    radius: 0.8,
    colors: [Color(0x2435A862), Color(0x0035A862)],
  );

  // ── Refinements ─────────────────────────────────────────────────────────
  static const Color scrim = Color(0xCC000000);
  static const Color hairline = Color(0x332A302C);
  static const Color pressedOverlay = Color(0x1F35A862);
  static const Color skeletonBase = Color(0xFF0D1110);
  static const Color skeletonHighlight = Color(0xFF1B201C);
}
