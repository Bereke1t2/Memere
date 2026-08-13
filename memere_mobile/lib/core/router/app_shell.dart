import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../constants/app_colors.dart';
import '../constants/app_motion.dart';
import '../constants/app_shadows.dart';
import '../constants/app_sizes.dart';
import '../constants/app_text_styles.dart';
import '../../shared/widgets/app_surface.dart';

/// Hosts the four primary tabs in a persistent bottom navigation. Each tab
/// keeps its own navigation stack via the [StatefulNavigationShell].
class AppShell extends StatelessWidget {
  const AppShell({super.key, required this.navigationShell});

  final StatefulNavigationShell navigationShell;

  void _onDestinationSelected(int index) {
    navigationShell.goBranch(
      index,
      initialLocation: index == navigationShell.currentIndex,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      extendBody: true,
      body: AnimatedSwitcher(
        duration: AppMotion.base,
        switchInCurve: AppMotion.standard,
        switchOutCurve: AppMotion.exit,
        transitionBuilder: (child, animation) {
          return FadeTransition(
            opacity: animation,
            child: SlideTransition(
              position: Tween<Offset>(
                begin: const Offset(0.02, 0),
                end: Offset.zero,
              ).animate(animation),
              child: child,
            ),
          );
        },
        child: KeyedSubtree(
          key: ValueKey(navigationShell.currentIndex),
          child: navigationShell,
        ),
      ),
      bottomNavigationBar: SafeArea(
        top: false,
        minimum: const EdgeInsets.fromLTRB(
          AppSizes.md,
          0,
          AppSizes.md,
          AppSizes.md,
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(AppSizes.radiusXl),
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 18, sigmaY: 18),
            child: Container(
              height: AppSizes.bottomNavHeight,
              padding: const EdgeInsets.symmetric(horizontal: AppSizes.xs),
              decoration: BoxDecoration(
                color: AppColors.bgSecondary.withAlpha(214),
                borderRadius: BorderRadius.circular(AppSizes.radiusXl),
                border: Border.all(color: AppColors.borderStrong),
                boxShadow: AppShadows.lg,
              ),
              child: Row(
                children: [
                  _NavItem(
                    label: 'Home',
                    icon: Icons.home_outlined,
                    selectedIcon: Icons.home_rounded,
                    selected: navigationShell.currentIndex == 0,
                    onTap: () => _onDestinationSelected(0),
                  ),
                  _NavItem(
                    label: 'Lessons',
                    icon: Icons.play_circle_outline_rounded,
                    selectedIcon: Icons.play_circle_rounded,
                    selected: navigationShell.currentIndex == 1,
                    onTap: () => _onDestinationSelected(1),
                  ),
                  _NavItem(
                    label: 'Exam',
                    icon: Icons.assignment_outlined,
                    selectedIcon: Icons.assignment_rounded,
                    selected: navigationShell.currentIndex == 2,
                    onTap: () => _onDestinationSelected(2),
                  ),
                  _NavItem(
                    label: 'Saved',
                    icon: Icons.bookmark_border_rounded,
                    selectedIcon: Icons.bookmark_rounded,
                    selected: navigationShell.currentIndex == 3,
                    onTap: () => _onDestinationSelected(3),
                  ),
                  _NavItem(
                    label: 'Profile',
                    icon: Icons.person_outline_rounded,
                    selectedIcon: Icons.person_rounded,
                    selected: navigationShell.currentIndex == 4,
                    onTap: () => _onDestinationSelected(4),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _NavItem extends StatelessWidget {
  const _NavItem({
    required this.label,
    required this.icon,
    required this.selectedIcon,
    required this.selected,
    required this.onTap,
  });

  final String label;
  final IconData icon;
  final IconData selectedIcon;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: AppPressable(
        onTap: onTap,
        borderRadius: AppSizes.radiusLg,
        child: SizedBox(
          height: AppSizes.bottomNavHeight,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              AnimatedContainer(
                duration: AppMotion.base,
                curve: AppMotion.standard,
                padding: EdgeInsets.symmetric(
                  horizontal: selected ? AppSizes.md : AppSizes.sm,
                  vertical: AppSizes.xs,
                ),
                decoration: BoxDecoration(
                  color: selected
                      ? AppColors.brandEmerald.withOpacity(0.18)
                      : Colors.transparent,
                  borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                ),
                child: Icon(
                  selected ? selectedIcon : icon,
                  color: selected ? AppColors.brandEmerald : AppColors.textMuted,
                  size: selected ? AppSizes.iconLg : AppSizes.iconMd,
                ),
              ),
              if (selected)
                Padding(
                  padding: const EdgeInsets.only(top: 2),
                  child: Text(
                    label,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      color: AppColors.brandEmerald,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
