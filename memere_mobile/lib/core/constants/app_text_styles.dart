import 'package:flutter/material.dart';
import 'app_colors.dart';

/// Memere Typography System
/// Display font: Sora (headings, titles)
/// Body font: DM Sans (body, labels, buttons)
abstract class AppTextStyles {
  // ── Display (Sora) ───────────────────────────────────────────────────────
  static const TextStyle displayLarge = TextStyle(
    fontFamily: 'Sora',
    fontSize: 36,
    fontWeight: FontWeight.w700,
    color: AppColors.textPrimary,
    height: 1.2,
    letterSpacing: 0,
  );

  static const TextStyle displayMedium = TextStyle(
    fontFamily: 'Sora',
    fontSize: 28,
    fontWeight: FontWeight.w700,
    color: AppColors.textPrimary,
    height: 1.25,
    letterSpacing: 0,
  );

  // ── Headlines (Sora) ─────────────────────────────────────────────────────
  static const TextStyle headlineLarge = TextStyle(
    fontFamily: 'Sora',
    fontSize: 22,
    fontWeight: FontWeight.w700,
    color: AppColors.textPrimary,
    height: 1.3,
  );

  static const TextStyle headlineMedium = TextStyle(
    fontFamily: 'Sora',
    fontSize: 18,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.35,
  );

  static const TextStyle headlineSmall = TextStyle(
    fontFamily: 'Sora',
    fontSize: 16,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.4,
  );

  // ── Title (Sora) ─────────────────────────────────────────────────────────
  static const TextStyle titleLarge = TextStyle(
    fontFamily: 'Sora',
    fontSize: 18,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.4,
  );

  static const TextStyle titleMedium = TextStyle(
    fontFamily: 'Sora',
    fontSize: 15,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.4,
  );

  // ── Body (DM Sans) ───────────────────────────────────────────────────────
  static const TextStyle bodyLarge = TextStyle(
    fontFamily: 'DM_Sans',
    fontSize: 16,
    fontWeight: FontWeight.w400,
    color: AppColors.textPrimary,
    height: 1.5,
  );

  static const TextStyle bodyMedium = TextStyle(
    fontFamily: 'DM_Sans',
    fontSize: 14,
    fontWeight: FontWeight.w400,
    color: AppColors.textPrimary,
    height: 1.5,
  );

  static const TextStyle bodySmall = TextStyle(
    fontFamily: 'DM_Sans',
    fontSize: 12,
    fontWeight: FontWeight.w400,
    color: AppColors.textSecondary,
    height: 1.5,
  );

  // ── Labels / Buttons (DM Sans) ───────────────────────────────────────────
  static const TextStyle labelLarge = TextStyle(
    fontFamily: 'DM_Sans',
    fontSize: 15,
    fontWeight: FontWeight.w600,
    color: AppColors.textPrimary,
    height: 1.2,
    letterSpacing: 0,
  );

  static const TextStyle labelMedium = TextStyle(
    fontFamily: 'DM_Sans',
    fontSize: 13,
    fontWeight: FontWeight.w500,
    color: AppColors.textPrimary,
    height: 1.2,
  );

  static const TextStyle labelSmall = TextStyle(
    fontFamily: 'DM_Sans',
    fontSize: 11,
    fontWeight: FontWeight.w500,
    color: AppColors.textSecondary,
    height: 1.2,
    letterSpacing: 0,
  );

  // ── Caption ──────────────────────────────────────────────────────────────
  static const TextStyle caption = TextStyle(
    fontFamily: 'DM_Sans',
    fontSize: 11,
    fontWeight: FontWeight.w400,
    color: AppColors.textDisabled,
    height: 1.4,
  );

  // ── Helpers ──────────────────────────────────────────────────────────────
  static TextStyle withColor(TextStyle style, Color color) =>
      style.copyWith(color: color);

  static TextStyle withSize(TextStyle style, double size) =>
      style.copyWith(fontSize: size);
}
