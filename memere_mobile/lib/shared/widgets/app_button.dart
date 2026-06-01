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
