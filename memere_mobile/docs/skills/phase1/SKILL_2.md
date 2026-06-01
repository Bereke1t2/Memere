# phase1/SKILL_2.md — Design System (Colors, Typography, Theme, Shared Widgets)
# ExamPrep Mobile (memere_mobile) — Phase 1, Part 2
# READ SKILL.md → SKILL_1.md → then this file.

---

## OBJECTIVE

Build the complete design system: colors, typography, spacing, theme, and reusable shared
widgets. Inspired by the ChatGPT/Claude dark mobile UI from the Figma reference. Every
screen in every phase will import from this system — it must be complete and stable before
any feature work begins.

Reference: `memere_mobile/docs/memere_Design_Specification.md`

---

## FILE 1 — `lib/core/constants/app_colors.dart`

```dart
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
```

---

## FILE 2 — `lib/core/constants/app_text_styles.dart`

```dart
import 'package:flutter/material.dart';
import 'app_colors.dart';

/// ExamPrep Typography System
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
    letterSpacing: -0.5,
  );

  static const TextStyle displayMedium = TextStyle(
    fontFamily: 'Sora',
    fontSize: 28,
    fontWeight: FontWeight.w700,
    color: AppColors.textPrimary,
    height: 1.25,
    letterSpacing: -0.3,
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
    letterSpacing: 0.2,
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
    letterSpacing: 0.3,
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
```

---

## FILE 3 — `lib/core/constants/app_sizes.dart`

```dart
/// Spacing, border radius, icon sizes, elevation constants
abstract class AppSizes {
  // ── Spacing (8px grid) ───────────────────────────────────────────────────
  static const double xs   = 4.0;
  static const double sm   = 8.0;
  static const double md   = 16.0;
  static const double lg   = 24.0;
  static const double xl   = 32.0;
  static const double xxl  = 48.0;
  static const double xxxl = 64.0;

  // ── Border Radius ────────────────────────────────────────────────────────
  static const double radiusSm   = 8.0;
  static const double radiusMd   = 12.0;
  static const double radiusLg   = 16.0;
  static const double radiusXl   = 24.0;
  static const double radiusFull = 999.0;

  // ── Icon Sizes ───────────────────────────────────────────────────────────
  static const double iconXs = 14.0;
  static const double iconSm = 18.0;
  static const double iconMd = 24.0;
  static const double iconLg = 32.0;
  static const double iconXl = 48.0;

  // ── Component Heights ────────────────────────────────────────────────────
  static const double buttonHeight      = 52.0;
  static const double buttonHeightSm    = 40.0;
  static const double inputHeight       = 52.0;
  static const double appBarHeight      = 56.0;
  static const double bottomNavHeight   = 64.0;
  static const double cardMinHeight     = 120.0;
  static const double courseCardHeight  = 200.0;
  static const double avatarSm          = 32.0;
  static const double avatarMd          = 44.0;
  static const double avatarLg          = 64.0;

  // ── Elevation ────────────────────────────────────────────────────────────
  static const double elevationNone  = 0.0;
  static const double elevationSm    = 2.0;
  static const double elevationMd    = 6.0;
  static const double elevationLg    = 12.0;
  static const double elevationXl    = 24.0;

  // ── Screen Padding ───────────────────────────────────────────────────────
  static const double screenPaddingH = 20.0; // horizontal screen edge padding
  static const double screenPaddingV = 24.0; // vertical screen edge padding
}
```

---

## FILE 4 — `lib/core/constants/app_constants.dart`

```dart
abstract class AppConstants {
  // ── Storage Keys ─────────────────────────────────────────────────────────
  static const String accessTokenKey  = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userBoxKey      = 'user_box';
  static const String courseBoxKey    = 'course_box';
  static const String prefsBoxKey     = 'prefs_box';
  static const String syncQueueKey    = 'sync_queue';

  // ── Cache TTL ────────────────────────────────────────────────────────────
  static const Duration cacheTtl             = Duration(hours: 1);
  static const Duration downloadedVideoTtl   = Duration(days: 30);

  // ── Pagination ───────────────────────────────────────────────────────────
  static const int defaultPageLimit = 20;

  // ── Video ────────────────────────────────────────────────────────────────
  static const int videoProgressSaveIntervalSeconds = 30;

  // ── Auth ─────────────────────────────────────────────────────────────────
  static const int accessTokenTtlMinutes  = 15;
  static const int refreshTokenTtlDays    = 30;
  static const int maxLoginAttempts       = 5;

  // ── Exam ─────────────────────────────────────────────────────────────────
  static const int examAutoSaveIntervalSeconds = 30;

  // ── Connectivity ─────────────────────────────────────────────────────────
  static const Duration connectionTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout    = Duration(seconds: 30);
}
```

---

## FILE 5 — `lib/core/theme/app_theme.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../constants/app_colors.dart';
import '../constants/app_text_styles.dart';
import '../constants/app_sizes.dart';

abstract class AppTheme {
  static ThemeData dark() {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      scaffoldBackgroundColor: AppColors.bgPrimary,
      primaryColor: AppColors.accentPrimary,
      colorScheme: const ColorScheme.dark(
        primary: AppColors.accentPrimary,
        secondary: AppColors.accentSecondary,
        surface: AppColors.bgSecondary,
        error: AppColors.error,
        onPrimary: AppColors.textInverse,
        onSecondary: AppColors.textInverse,
        onSurface: AppColors.textPrimary,
        onError: AppColors.textPrimary,
      ),

      // ── AppBar ─────────────────────────────────────────────────────────
      appBarTheme: const AppBarTheme(
        backgroundColor: AppColors.bgPrimary,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        titleTextStyle: AppTextStyles.titleLarge,
        iconTheme: IconThemeData(color: AppColors.textPrimary, size: AppSizes.iconMd),
        systemOverlayStyle: SystemUiOverlayStyle(
          statusBarColor: Colors.transparent,
          statusBarIconBrightness: Brightness.light,
        ),
      ),

      // ── Cards ──────────────────────────────────────────────────────────
      cardTheme: CardThemeData(
        color: AppColors.bgSecondary,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusLg),
          side: const BorderSide(color: AppColors.border, width: 1),
        ),
      ),

      // ── Input Fields ───────────────────────────────────────────────────
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: AppColors.bgTertiary,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: AppSizes.md,
          vertical: AppSizes.md,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          borderSide: const BorderSide(color: AppColors.border),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          borderSide: const BorderSide(color: AppColors.border),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          borderSide: const BorderSide(color: AppColors.borderFocused, width: 1.5),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          borderSide: const BorderSide(color: AppColors.error),
        ),
        focusedErrorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          borderSide: const BorderSide(color: AppColors.error, width: 1.5),
        ),
        hintStyle: AppTextStyles.bodyMedium.copyWith(color: AppColors.textDisabled),
        labelStyle: AppTextStyles.bodyMedium.copyWith(color: AppColors.textSecondary),
        errorStyle: AppTextStyles.bodySmall.copyWith(color: AppColors.error),
        prefixIconColor: AppColors.textSecondary,
        suffixIconColor: AppColors.textSecondary,
      ),

      // ── Elevated Button ────────────────────────────────────────────────
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.accentPrimary,
          foregroundColor: Colors.white,
          minimumSize: const Size(double.infinity, AppSizes.buttonHeight),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          ),
          elevation: 0,
          textStyle: AppTextStyles.labelLarge,
        ),
      ),

      // ── Outlined Button ────────────────────────────────────────────────
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: AppColors.accentPrimary,
          minimumSize: const Size(double.infinity, AppSizes.buttonHeight),
          side: const BorderSide(color: AppColors.accentPrimary),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          ),
          textStyle: AppTextStyles.labelLarge,
        ),
      ),

      // ── Text Button ────────────────────────────────────────────────────
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: AppColors.accentPrimary,
          textStyle: AppTextStyles.labelMedium,
        ),
      ),

      // ── Bottom Navigation ──────────────────────────────────────────────
      bottomNavigationBarTheme: const BottomNavigationBarThemeData(
        backgroundColor: AppColors.bgSecondary,
        selectedItemColor: AppColors.accentPrimary,
        unselectedItemColor: AppColors.textDisabled,
        elevation: 0,
        type: BottomNavigationBarType.fixed,
      ),

      // ── Chip ───────────────────────────────────────────────────────────
      chipTheme: ChipThemeData(
        backgroundColor: AppColors.bgTertiary,
        selectedColor: AppColors.accentGlow,
        labelStyle: AppTextStyles.labelSmall,
        side: const BorderSide(color: AppColors.border),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        ),
        padding: const EdgeInsets.symmetric(horizontal: AppSizes.sm, vertical: AppSizes.xs),
      ),

      // ── Divider ────────────────────────────────────────────────────────
      dividerTheme: const DividerThemeData(
        color: AppColors.divider,
        thickness: 1,
        space: 1,
      ),

      // ── Text ───────────────────────────────────────────────────────────
      textTheme: const TextTheme(
        displayLarge:   AppTextStyles.displayLarge,
        displayMedium:  AppTextStyles.displayMedium,
        headlineLarge:  AppTextStyles.headlineLarge,
        headlineMedium: AppTextStyles.headlineMedium,
        headlineSmall:  AppTextStyles.headlineSmall,
        titleLarge:     AppTextStyles.titleLarge,
        titleMedium:    AppTextStyles.titleMedium,
        bodyLarge:      AppTextStyles.bodyLarge,
        bodyMedium:     AppTextStyles.bodyMedium,
        bodySmall:      AppTextStyles.bodySmall,
        labelLarge:     AppTextStyles.labelLarge,
        labelMedium:    AppTextStyles.labelMedium,
        labelSmall:     AppTextStyles.labelSmall,
      ),
    );
  }
}
```

---

## FILE 6 — `lib/shared/widgets/app_button.dart`

```dart
import 'package:flutter/material.dart';
import '../../core/constants/app_colors.dart';
import '../../core/constants/app_sizes.dart';
import '../../core/constants/app_text_styles.dart';

enum AppButtonVariant { primary, secondary, outline, ghost, danger }

class AppButton extends StatelessWidget {
  const AppButton({
    super.key,
    required this.label,
    required this.onPressed,
    this.variant = AppButtonVariant.primary,
    this.isLoading = false,
    this.isDisabled = false,
    this.prefixIcon,
    this.suffixIcon,
    this.height = AppSizes.buttonHeight,
    this.width,
  });

  final String label;
  final VoidCallback? onPressed;
  final AppButtonVariant variant;
  final bool isLoading;
  final bool isDisabled;
  final IconData? prefixIcon;
  final IconData? suffixIcon;
  final double height;
  final double? width;

  @override
  Widget build(BuildContext context) {
    final bool disabled = isDisabled || isLoading || onPressed == null;

    return SizedBox(
      height: height,
      width: width ?? double.infinity,
      child: _buildButton(disabled),
    );
  }

  Widget _buildButton(bool disabled) {
    switch (variant) {
      case AppButtonVariant.primary:
        return _PrimaryButton(
          label: label, onPressed: disabled ? null : onPressed,
          isLoading: isLoading, prefixIcon: prefixIcon, suffixIcon: suffixIcon,
        );
      case AppButtonVariant.secondary:
        return _SecondaryButton(
          label: label, onPressed: disabled ? null : onPressed,
          isLoading: isLoading,
        );
      case AppButtonVariant.outline:
        return _OutlineButton(
          label: label, onPressed: disabled ? null : onPressed,
          isLoading: isLoading,
        );
      case AppButtonVariant.ghost:
        return _GhostButton(label: label, onPressed: disabled ? null : onPressed);
      case AppButtonVariant.danger:
        return _DangerButton(
          label: label, onPressed: disabled ? null : onPressed,
          isLoading: isLoading,
        );
    }
  }
}

class _ButtonContent extends StatelessWidget {
  const _ButtonContent({
    required this.label,
    required this.labelColor,
    this.isLoading = false,
    this.prefixIcon,
    this.suffixIcon,
    this.loadingColor = Colors.white,
  });

  final String label;
  final Color labelColor;
  final bool isLoading;
  final IconData? prefixIcon;
  final IconData? suffixIcon;
  final Color loadingColor;

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return SizedBox(
        width: 20, height: 20,
        child: CircularProgressIndicator(
          strokeWidth: 2, valueColor: AlwaysStoppedAnimation(loadingColor),
        ),
      );
    }
    return Row(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        if (prefixIcon != null) ...[
          Icon(prefixIcon, size: AppSizes.iconSm, color: labelColor),
          const SizedBox(width: AppSizes.sm),
        ],
        Text(label, style: AppTextStyles.labelLarge.copyWith(color: labelColor)),
        if (suffixIcon != null) ...[
          const SizedBox(width: AppSizes.sm),
          Icon(suffixIcon, size: AppSizes.iconSm, color: labelColor),
        ],
      ],
    );
  }
}

class _PrimaryButton extends StatelessWidget {
  const _PrimaryButton({
    required this.label, required this.onPressed,
    required this.isLoading, this.prefixIcon, this.suffixIcon,
  });
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;
  final IconData? prefixIcon;
  final IconData? suffixIcon;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: onPressed != null ? AppColors.primaryGradient : null,
        color: onPressed == null ? AppColors.textDisabled : null,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
      ),
      child: ElevatedButton(
        onPressed: onPressed,
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.transparent,
          shadowColor: Colors.transparent,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          ),
        ),
        child: _ButtonContent(
          label: label, labelColor: Colors.white,
          isLoading: isLoading, prefixIcon: prefixIcon, suffixIcon: suffixIcon,
        ),
      ),
    );
  }
}

class _SecondaryButton extends StatelessWidget {
  const _SecondaryButton({required this.label, required this.onPressed, required this.isLoading});
  final String label; final VoidCallback? onPressed; final bool isLoading;

  @override
  Widget build(BuildContext context) => ElevatedButton(
    onPressed: onPressed,
    style: ElevatedButton.styleFrom(
      backgroundColor: AppColors.bgTertiary,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        side: const BorderSide(color: AppColors.border),
      ),
    ),
    child: _ButtonContent(label: label, labelColor: AppColors.textPrimary, isLoading: isLoading),
  );
}

class _OutlineButton extends StatelessWidget {
  const _OutlineButton({required this.label, required this.onPressed, required this.isLoading});
  final String label; final VoidCallback? onPressed; final bool isLoading;

  @override
  Widget build(BuildContext context) => OutlinedButton(
    onPressed: onPressed,
    child: _ButtonContent(label: label, labelColor: AppColors.accentPrimary, isLoading: isLoading, loadingColor: AppColors.accentPrimary),
  );
}

class _GhostButton extends StatelessWidget {
  const _GhostButton({required this.label, required this.onPressed});
  final String label; final VoidCallback? onPressed;

  @override
  Widget build(BuildContext context) => TextButton(
    onPressed: onPressed,
    child: Text(label, style: AppTextStyles.labelMedium.copyWith(color: AppColors.accentPrimary)),
  );
}

class _DangerButton extends StatelessWidget {
  const _DangerButton({required this.label, required this.onPressed, required this.isLoading});
  final String label; final VoidCallback? onPressed; final bool isLoading;

  @override
  Widget build(BuildContext context) => ElevatedButton(
    onPressed: onPressed,
    style: ElevatedButton.styleFrom(backgroundColor: AppColors.errorSurface),
    child: _ButtonContent(label: label, labelColor: AppColors.error, isLoading: isLoading, loadingColor: AppColors.error),
  );
}
```

---

## FILE 7 — `lib/shared/widgets/app_text_field.dart`

```dart
import 'package:flutter/material.dart';
import '../../core/constants/app_colors.dart';
import '../../core/constants/app_sizes.dart';
import '../../core/constants/app_text_styles.dart';

class AppTextField extends StatefulWidget {
  const AppTextField({
    super.key,
    required this.controller,
    required this.hintText,
    this.labelText,
    this.prefixIcon,
    this.suffixIcon,
    this.isPassword = false,
    this.keyboardType = TextInputType.text,
    this.textInputAction = TextInputAction.next,
    this.validator,
    this.onChanged,
    this.onFieldSubmitted,
    this.autofillHints,
    this.enabled = true,
    this.maxLines = 1,
    this.focusNode,
  });

  final TextEditingController controller;
  final String hintText;
  final String? labelText;
  final IconData? prefixIcon;
  final Widget? suffixIcon;
  final bool isPassword;
  final TextInputType keyboardType;
  final TextInputAction textInputAction;
  final String? Function(String?)? validator;
  final void Function(String)? onChanged;
  final void Function(String)? onFieldSubmitted;
  final List<String>? autofillHints;
  final bool enabled;
  final int maxLines;
  final FocusNode? focusNode;

  @override
  State<AppTextField> createState() => _AppTextFieldState();
}

class _AppTextFieldState extends State<AppTextField> {
  bool _obscureText = true;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (widget.labelText != null) ...[
          Text(widget.labelText!, style: AppTextStyles.labelMedium.copyWith(color: AppColors.textSecondary)),
          const SizedBox(height: AppSizes.xs),
        ],
        TextFormField(
          controller: widget.controller,
          focusNode: widget.focusNode,
          obscureText: widget.isPassword && _obscureText,
          keyboardType: widget.keyboardType,
          textInputAction: widget.textInputAction,
          validator: widget.validator,
          onChanged: widget.onChanged,
          onFieldSubmitted: widget.onFieldSubmitted,
          autofillHints: widget.autofillHints,
          enabled: widget.enabled,
          maxLines: widget.isPassword ? 1 : widget.maxLines,
          style: AppTextStyles.bodyMedium,
          decoration: InputDecoration(
            hintText: widget.hintText,
            prefixIcon: widget.prefixIcon != null
                ? Icon(widget.prefixIcon, size: AppSizes.iconSm, color: AppColors.textSecondary)
                : null,
            suffixIcon: widget.isPassword
                ? IconButton(
                    icon: Icon(
                      _obscureText ? Icons.visibility_off_outlined : Icons.visibility_outlined,
                      size: AppSizes.iconSm,
                      color: AppColors.textSecondary,
                    ),
                    onPressed: () => setState(() => _obscureText = !_obscureText),
                  )
                : widget.suffixIcon,
          ),
        ),
      ],
    );
  }
}
```

---

## FILE 8 — `lib/shared/widgets/loading_overlay.dart`

```dart
import 'package:flutter/material.dart';
import '../../core/constants/app_colors.dart';

class LoadingOverlay extends StatelessWidget {
  const LoadingOverlay({super.key, required this.child, required this.isLoading});
  final Widget child;
  final bool isLoading;

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        child,
        if (isLoading)
          const ColoredBox(
            color: AppColors.bgOverlay,
            child: Center(
              child: CircularProgressIndicator(
                color: AppColors.accentPrimary,
                strokeWidth: 2.5,
              ),
            ),
          ),
      ],
    );
  }
}
```

---

## FILE 9 — `lib/shared/extensions/context_extensions.dart`

```dart
import 'package:flutter/material.dart';
import '../../core/constants/app_sizes.dart';

extension BuildContextX on BuildContext {
  ThemeData get theme => Theme.of(this);
  TextTheme get textTheme => Theme.of(this).textTheme;
  ColorScheme get colorScheme => Theme.of(this).colorScheme;
  MediaQueryData get mediaQuery => MediaQuery.of(this);
  Size get screenSize => MediaQuery.of(this).size;
  double get screenWidth => MediaQuery.of(this).size.width;
  double get screenHeight => MediaQuery.of(this).size.height;
  EdgeInsets get viewPadding => MediaQuery.of(this).viewPadding;
  bool get isKeyboardOpen => MediaQuery.of(this).viewInsets.bottom > 0;

  void showSnackBar(String message, {bool isError = false}) {
    ScaffoldMessenger.of(this).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: isError ? const Color(0xFFCF6679) : const Color(0xFF4CAF50),
        behavior: SnackBarBehavior.floating,
        margin: const EdgeInsets.all(AppSizes.md),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppSizes.radiusSm)),
      ),
    );
  }
}
```

---

## PHASE 1 CHECKLIST — SKILL_2 COMPLETE WHEN:

- [ ] All 9 files above created at correct paths
- [ ] `AppTheme.dark()` referenced in `app.dart` compiles
- [ ] `AppButton` renders in 5 variants without errors
- [ ] `AppTextField` renders with password toggle
- [ ] `flutter analyze` passes with 0 errors

---

## NEXT: Go to `phase1/SKILL_3.md` — Core Infrastructure + Auth Feature + Screens
