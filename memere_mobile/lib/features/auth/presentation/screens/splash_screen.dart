import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/memere_mascot.dart';
import '../providers/auth_state_provider.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _fadeAnim;
  late Animation<double> _scaleAnim;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1200),
    );
    _fadeAnim = Tween<double>(begin: 0, end: 1).animate(
      CurvedAnimation(
        parent: _controller,
        curve: const Interval(0, 0.6, curve: Curves.easeIn),
      ),
    );
    _scaleAnim = Tween<double>(begin: 0.85, end: 1).animate(
      CurvedAnimation(
        parent: _controller,
        curve: const Interval(0, 0.7, curve: Curves.elasticOut),
      ),
    );
    _controller.forward();
    _scheduleStartupNavigation();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        child: Center(
          child: FadeTransition(
            opacity: _fadeAnim,
            child: ScaleTransition(
              scale: _scaleAnim,
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: AppSizes.md),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    // Animated Mascot Reading Book & Laptop
                    const MemereMascot(
                      size: Size(250, 225),
                      showBackdrop: true,
                    ),
                    const SizedBox(height: AppSizes.lg),

                    // App Title & Tagline
                    const Text('Memere', style: AppTextStyles.displayMedium),
                    const SizedBox(height: AppSizes.xs),
                    Text(
                      'Interactive Course Learning Platform',
                      textAlign: TextAlign.center,
                      style: AppTextStyles.bodyMedium.copyWith(
                        color: AppColors.textSecondary,
                        letterSpacing: 0.2,
                      ),
                    ),
                    const SizedBox(height: AppSizes.xl),

                    // Subtle Progress Bar to indicate loading
                    SizedBox(
                      width: 140,
                      child: LinearProgressIndicator(
                        backgroundColor: AppColors.bgTertiary,
                        color: AppColors.textPrimary,
                        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                        minHeight: 3,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _scheduleStartupNavigation() async {
    try {
      // 1. Minimum splash delay for smooth visual transition
      await Future<void>.delayed(const Duration(milliseconds: 2200));
      if (!mounted) return;

      // 2. Read auth state with 3s timeout & fallback so splash NEVER hangs
      final authAsync = ref.read(authStateProvider);
      final authState = authAsync.valueOrNull ??
          await ref
              .read(authStateProvider.future)
              .timeout(const Duration(seconds: 3))
              .catchError((_) => const AuthState());

      if (!mounted) return;

      if (authState.isAuthenticated) {
        context.go(AppRoutes.home);
      } else if (authState.hasSeenOnboarding) {
        context.go(AppRoutes.login);
      } else {
        context.go(AppRoutes.onboarding);
      }
    } catch (_) {
      if (!mounted) return;
      // Safety fallback navigation: Go to login screen on any unexpected network exception
      context.go(AppRoutes.login);
    }
  }
}
