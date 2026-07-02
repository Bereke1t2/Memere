import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../../core/constants/app_colors.dart';
import '../../core/constants/app_motion.dart';
import '../../core/constants/app_shadows.dart';
import '../../core/constants/app_sizes.dart';
import '../../core/constants/app_text_styles.dart';

enum AppButtonVariant { primary, secondary, outline, ghost, danger }

class AppButton extends StatefulWidget {
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
  State<AppButton> createState() => _AppButtonState();
}

class _AppButtonState extends State<AppButton> {
  bool _pressed = false;

  bool get _disabled =>
      widget.isDisabled || widget.isLoading || widget.onPressed == null;

  void _setPressed(bool value) {
    if (_disabled || _pressed == value) return;
    setState(() => _pressed = value);
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTapDown: _disabled ? null : (_) => _setPressed(true),
      onTapCancel: _disabled ? null : () => _setPressed(false),
      onTapUp: _disabled ? null : (_) => _setPressed(false),
      child: AnimatedScale(
        scale: _pressed ? 0.95 : 1,
        duration: AppMotion.fast,
        curve: AppMotion.standard,
        child: SizedBox(
          height: widget.height,
          width: widget.width ?? double.infinity,
          child: _buildButton(_disabled),
        ),
      ),
    );
  }

  Widget _buildButton(bool disabled) {
    final onPressed = disabled ? null : _withHaptic(widget.onPressed);
    switch (widget.variant) {
      case AppButtonVariant.primary:
        return _PrimaryButton(
          label: widget.label,
          onPressed: onPressed,
          isLoading: widget.isLoading,
          prefixIcon: widget.prefixIcon,
          suffixIcon: widget.suffixIcon,
          height: widget.height,
        );
      case AppButtonVariant.secondary:
        return _SecondaryButton(
          label: widget.label,
          onPressed: onPressed,
          isLoading: widget.isLoading,
          prefixIcon: widget.prefixIcon,
          suffixIcon: widget.suffixIcon,
          height: widget.height,
        );
      case AppButtonVariant.outline:
        return _OutlineButton(
          label: widget.label,
          onPressed: onPressed,
          isLoading: widget.isLoading,
          prefixIcon: widget.prefixIcon,
          suffixIcon: widget.suffixIcon,
          height: widget.height,
        );
      case AppButtonVariant.ghost:
        return _GhostButton(
          label: widget.label,
          onPressed: onPressed,
          prefixIcon: widget.prefixIcon,
          suffixIcon: widget.suffixIcon,
        );
      case AppButtonVariant.danger:
        return _DangerButton(
          label: widget.label,
          onPressed: onPressed,
          isLoading: widget.isLoading,
          prefixIcon: widget.prefixIcon,
          suffixIcon: widget.suffixIcon,
          height: widget.height,
        );
    }
  }

  VoidCallback? _withHaptic(VoidCallback? onPressed) {
    if (onPressed == null) return null;
    return () {
      HapticFeedback.selectionClick();
      onPressed();
    };
  }
}

class _ButtonContent extends StatelessWidget {
  const _ButtonContent({
    required this.label,
    required this.labelColor,
    this.isLoading = false,
    this.prefixIcon,
    this.suffixIcon,
    this.loadingColor = AppColors.textPrimary,
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
        width: 20,
        height: 20,
        child: CircularProgressIndicator(
          strokeWidth: 2,
          valueColor: AlwaysStoppedAnimation(loadingColor),
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
        Flexible(
          child: Text(
            label,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: AppTextStyles.labelLarge.copyWith(color: labelColor),
          ),
        ),
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
    required this.label,
    required this.onPressed,
    required this.isLoading,
    required this.height,
    this.prefixIcon,
    this.suffixIcon,
  });
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;
  final double height;
  final IconData? prefixIcon;
  final IconData? suffixIcon;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        color: onPressed != null
            ? AppColors.accentPrimary
            : AppColors.bgQuaternary,
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        boxShadow: onPressed != null ? AppShadows.sm : null,
      ),
      child: ElevatedButton(
        onPressed: onPressed,
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.transparent,
          shadowColor: Colors.transparent,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppSizes.radiusFull),
            side: BorderSide.none,
          ),
          minimumSize: Size(0, height),
          padding: const EdgeInsets.symmetric(horizontal: AppSizes.md),
        ),
        child: _ButtonContent(
          label: label,
          labelColor: AppColors.textInverse,
          isLoading: isLoading,
          prefixIcon: prefixIcon,
          suffixIcon: suffixIcon,
          loadingColor: AppColors.textInverse,
        ),
      ),
    );
  }
}

class _SecondaryButton extends StatelessWidget {
  const _SecondaryButton({
    required this.label,
    required this.onPressed,
    required this.isLoading,
    required this.height,
    this.prefixIcon,
    this.suffixIcon,
  });
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;
  final double height;
  final IconData? prefixIcon;
  final IconData? suffixIcon;

  @override
  Widget build(BuildContext context) => ElevatedButton(
        onPressed: onPressed,
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.bgTertiary,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppSizes.radiusFull),
            side: const BorderSide(color: AppColors.borderStrong),
          ),
          minimumSize: Size(0, height),
          padding: const EdgeInsets.symmetric(horizontal: AppSizes.md),
        ),
        child: _ButtonContent(
          label: label,
          labelColor: AppColors.textPrimary,
          isLoading: isLoading,
          prefixIcon: prefixIcon,
          suffixIcon: suffixIcon,
        ),
      );
}

class _OutlineButton extends StatelessWidget {
  const _OutlineButton({
    required this.label,
    required this.onPressed,
    required this.isLoading,
    required this.height,
    this.prefixIcon,
    this.suffixIcon,
  });
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;
  final double height;
  final IconData? prefixIcon;
  final IconData? suffixIcon;

  @override
  Widget build(BuildContext context) => OutlinedButton(
        onPressed: onPressed,
        style: OutlinedButton.styleFrom(
          minimumSize: Size(0, height),
          padding: const EdgeInsets.symmetric(horizontal: AppSizes.md),
        ),
        child: _ButtonContent(
          label: label,
          labelColor: AppColors.accentPrimary,
          isLoading: isLoading,
          loadingColor: AppColors.accentPrimary,
          prefixIcon: prefixIcon,
          suffixIcon: suffixIcon,
        ),
      );
}

class _GhostButton extends StatelessWidget {
  const _GhostButton({
    required this.label,
    required this.onPressed,
    this.prefixIcon,
    this.suffixIcon,
  });
  final String label;
  final VoidCallback? onPressed;
  final IconData? prefixIcon;
  final IconData? suffixIcon;

  @override
  Widget build(BuildContext context) => TextButton(
        onPressed: onPressed,
        child: _ButtonContent(
          label: label,
          labelColor: AppColors.accentPrimary,
          prefixIcon: prefixIcon,
          suffixIcon: suffixIcon,
        ),
      );
}

class _DangerButton extends StatelessWidget {
  const _DangerButton({
    required this.label,
    required this.onPressed,
    required this.isLoading,
    required this.height,
    this.prefixIcon,
    this.suffixIcon,
  });
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;
  final double height;
  final IconData? prefixIcon;
  final IconData? suffixIcon;

  @override
  Widget build(BuildContext context) => ElevatedButton(
        onPressed: onPressed,
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.errorSurface,
          minimumSize: Size(0, height),
          padding: const EdgeInsets.symmetric(horizontal: AppSizes.md),
        ),
        child: _ButtonContent(
          label: label,
          labelColor: AppColors.error,
          isLoading: isLoading,
          loadingColor: AppColors.error,
          prefixIcon: prefixIcon,
          suffixIcon: suffixIcon,
        ),
      );
}
