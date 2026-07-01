import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../core/constants/app_motion.dart';
import '../../core/constants/app_shadows.dart';
import '../../core/constants/app_sizes.dart';
import '../../core/constants/app_text_styles.dart';

class AppPageBackground extends StatelessWidget {
  const AppPageBackground({
    super.key,
    required this.child,
  });

  final Widget child;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: const BoxDecoration(gradient: AppColors.backgroundGradient),
      child: child,
    );
  }
}

class AppPressable extends StatefulWidget {
  const AppPressable({
    super.key,
    required this.child,
    required this.onTap,
    this.borderRadius = AppSizes.radiusMd,
    this.enabled = true,
  });

  final Widget child;
  final VoidCallback? onTap;
  final double borderRadius;
  final bool enabled;

  @override
  State<AppPressable> createState() => _AppPressableState();
}

class _AppPressableState extends State<AppPressable> {
  bool _pressed = false;

  bool get _enabled => widget.enabled && widget.onTap != null;

  void _setPressed(bool value) {
    if (!_enabled || _pressed == value) return;
    setState(() => _pressed = value);
  }

  @override
  Widget build(BuildContext context) {
    return Semantics(
      button: _enabled,
      child: GestureDetector(
        behavior: HitTestBehavior.opaque,
        onTap: _enabled ? widget.onTap : null,
        onTapDown: (_) => _setPressed(true),
        onTapCancel: () => _setPressed(false),
        onTapUp: (_) => _setPressed(false),
        child: AnimatedScale(
          scale: _pressed ? 0.985 : 1,
          duration: AppMotion.fast,
          curve: AppMotion.standard,
          child: AnimatedOpacity(
            opacity: widget.enabled ? 1 : 0.58,
            duration: AppMotion.fast,
            curve: AppMotion.standard,
            child: widget.child,
          ),
        ),
      ),
    );
  }
}

class AppSurface extends StatelessWidget {
  const AppSurface({
    super.key,
    required this.child,
    this.onTap,
    this.padding = const EdgeInsets.all(AppSizes.md),
    this.margin,
    this.radius = AppSizes.radiusLg,
    this.color = AppColors.bgSecondary,
    this.gradient,
    this.borderColor = AppColors.hairline,
    this.shadows = AppShadows.sm,
    this.clipBehavior = Clip.antiAlias,
  });

  final Widget child;
  final VoidCallback? onTap;
  final EdgeInsetsGeometry padding;
  final EdgeInsetsGeometry? margin;
  final double radius;
  final Color color;
  final Gradient? gradient;
  final Color borderColor;
  final List<BoxShadow> shadows;
  final Clip clipBehavior;

  @override
  Widget build(BuildContext context) {
    final surface = Container(
      margin: margin,
      padding: padding,
      clipBehavior: clipBehavior,
      decoration: BoxDecoration(
        color: gradient == null ? color : null,
        gradient: gradient,
        borderRadius: BorderRadius.circular(radius),
        border: Border.all(color: borderColor),
        boxShadow: shadows,
      ),
      child: child,
    );

    if (onTap == null) return surface;
    return AppPressable(
      onTap: onTap,
      borderRadius: radius,
      child: surface,
    );
  }
}

class AppIconTile extends StatelessWidget {
  const AppIconTile({
    super.key,
    required this.icon,
    this.color = AppColors.accentPrimary,
    this.size = 48,
    this.iconSize = AppSizes.iconMd,
    this.gradient,
  });

  final IconData icon;
  final Color color;
  final double size;
  final double iconSize;
  final Gradient? gradient;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: gradient == null ? color.withAlpha(28) : null,
        gradient: gradient,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        border: Border.all(
          color: gradient == null ? color.withAlpha(70) : AppColors.hairline,
        ),
        boxShadow: gradient == null ? null : AppShadows.accentGlow,
      ),
      child: Icon(
        icon,
        size: iconSize,
        color: gradient == null ? color : Colors.white,
      ),
    );
  }
}

class AppBadge extends StatelessWidget {
  const AppBadge({
    super.key,
    required this.label,
    this.color = AppColors.accentPrimary,
    this.icon,
  });

  final String label;
  final Color color;
  final IconData? icon;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSizes.sm,
        vertical: AppSizes.xs,
      ),
      decoration: BoxDecoration(
        color: color.withAlpha(30),
        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
        border: Border.all(color: color.withAlpha(96)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (icon != null) ...[
            Icon(icon, size: AppSizes.iconXs, color: color),
            const SizedBox(width: AppSizes.xs),
          ],
          Text(
            label,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: AppTextStyles.labelSmall.copyWith(color: color),
          ),
        ],
      ),
    );
  }
}

class AppSectionHeader extends StatelessWidget {
  const AppSectionHeader({
    super.key,
    required this.title,
    this.subtitle,
    this.trailing,
  });

  final String title;
  final String? subtitle;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(title, style: AppTextStyles.headlineSmall),
              if (subtitle != null) ...[
                const SizedBox(height: AppSizes.xs),
                Text(
                  subtitle!,
                  style: AppTextStyles.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ],
          ),
        ),
        if (trailing != null) ...[
          const SizedBox(width: AppSizes.md),
          trailing!,
        ],
      ],
    );
  }
}

class AppStaggeredReveal extends StatefulWidget {
  const AppStaggeredReveal({
    super.key,
    required this.child,
    this.index = 0,
    this.offset = const Offset(0, 0.04),
  });

  final Widget child;
  final int index;
  final Offset offset;

  @override
  State<AppStaggeredReveal> createState() => _AppStaggeredRevealState();
}

class _AppStaggeredRevealState extends State<AppStaggeredReveal>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  late final Animation<double> _opacity;
  late final Animation<Offset> _slide;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: AppMotion.slow,
    );
    _opacity = CurvedAnimation(parent: _controller, curve: AppMotion.standard);
    _slide = Tween<Offset>(
      begin: widget.offset,
      end: Offset.zero,
    ).animate(_opacity);

    Future<void>.delayed(AppMotion.stagger * widget.index, () {
      if (mounted) _controller.forward();
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return FadeTransition(
      opacity: _opacity,
      child: SlideTransition(
        position: _slide,
        child: widget.child,
      ),
    );
  }
}
