import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../constants/app_colors.dart';
import '../constants/app_motion.dart';
import '../constants/app_shadows.dart';
import '../constants/app_sizes.dart';
import '../constants/app_text_styles.dart';
import '../../shared/widgets/app_surface.dart';

/// Hosts the four primary tabs (Browse · Learn · Exams · Profile) in a persistent
/// bottom navigation. Each tab keeps its own navigation stack via the
/// [StatefulNavigationShell]; detail/checkout screens push over this shell.
class AppShell extends StatelessWidget {
  const AppShell({super.key, required this.navigationShell});

  final StatefulNavigationShell navigationShell;

  void _onDestinationSelected(int index) {
    // Re-tapping the active tab returns it to its root.
    navigationShell.goBranch(
      index,
      initialLocation: index == navigationShell.currentIndex,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
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
          AppSizes.sm,
        ),
        child: AppSurface(
          padding: EdgeInsets.zero,
          radius: AppSizes.radiusXl,
          color: AppColors.bgSecondary.withAlpha(238),
          borderColor: AppColors.hairline,
          shadows: AppShadows.lg,
          child: ClipRRect(
            borderRadius: BorderRadius.circular(AppSizes.radiusXl),
            child: NavigationBarTheme(
              data: NavigationBarThemeData(
                height: AppSizes.bottomNavHeight,
                backgroundColor: Colors.transparent,
                indicatorColor: AppColors.bgQuaternary,
                indicatorShape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(AppSizes.radiusLg),
                ),
                surfaceTintColor: Colors.transparent,
                labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
                iconTheme: WidgetStateProperty.resolveWith(
                  (states) => IconThemeData(
                    size: AppSizes.iconMd,
                    color: states.contains(WidgetState.selected)
                        ? AppColors.accentPrimary
                        : AppColors.textSecondary,
                  ),
                ),
                labelTextStyle: WidgetStateProperty.resolveWith(
                  (states) => AppTextStyles.labelSmall.copyWith(
                    color: states.contains(WidgetState.selected)
                        ? AppColors.accentPrimary
                        : AppColors.textSecondary,
                  ),
                ),
              ),
              child: NavigationBar(
                selectedIndex: navigationShell.currentIndex,
                onDestinationSelected: _onDestinationSelected,
                destinations: const [
                  NavigationDestination(
                    icon: Icon(Icons.explore_outlined),
                    selectedIcon: Icon(Icons.explore_rounded),
                    label: 'Browse',
                  ),
                  NavigationDestination(
                    icon: Icon(Icons.school_outlined),
                    selectedIcon: Icon(Icons.school_rounded),
                    label: 'Learn',
                  ),
                  NavigationDestination(
                    icon: Icon(Icons.assignment_outlined),
                    selectedIcon: Icon(Icons.assignment_rounded),
                    label: 'Exams',
                  ),
                  NavigationDestination(
                    icon: Icon(Icons.person_outline_rounded),
                    selectedIcon: Icon(Icons.person_rounded),
                    label: 'Profile',
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
