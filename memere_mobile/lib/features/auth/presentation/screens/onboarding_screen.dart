import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_motion.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../../shared/widgets/memere_mascot.dart';
import '../providers/auth_state_provider.dart';

class _OnboardingPage {
  const _OnboardingPage({
    required this.title,
    required this.body,
  });

  final String title;
  final String body;
}

final _pages = [
  const _OnboardingPage(
    title: 'Meet Memere',
    body:
        'Selam. Learn with interactive AI-guided courses and comprehensive study tools.',
  ),
  const _OnboardingPage(
    title: 'Study Every Subject',
    body:
        'Math, Physics, Chemistry, Biology, English, Civics, and more in one focused plan.',
  ),
  const _OnboardingPage(
    title: 'Practice With Confidence',
    body:
        'Move from lessons to quizzes and mock exams with progress that keeps you on track.',
  ),
];

class OnboardingScreen extends ConsumerStatefulWidget {
  const OnboardingScreen({super.key});

  @override
  ConsumerState<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends ConsumerState<OnboardingScreen> {
  final _pageController = PageController();
  int _currentPage = 0;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  Future<void> _finish() async {
    await ref.read(authStateProvider.notifier).markOnboardingSeen();
    if (!mounted) return;
    context.go(AppRoutes.login);
  }

  void _next() {
    if (_currentPage < _pages.length - 1) {
      _pageController.nextPage(
        duration: AppMotion.slow,
        curve: AppMotion.standard,
      );
      return;
    }
    _finish();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: DecoratedBox(
        decoration: const BoxDecoration(gradient: AppColors.backgroundGradient),
        child: SafeArea(
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(
                  AppSizes.screenPaddingH,
                  AppSizes.sm,
                  AppSizes.screenPaddingH,
                  0,
                ),
                child: Row(
                  children: [
                    const AppIconTile(
                      icon: Icons.school_rounded,
                      size: 44,
                      iconSize: AppSizes.iconSm,
                      color: AppColors.accentPrimary,
                    ),
                    const SizedBox(width: AppSizes.sm),
                    Text(
                      'Memere',
                      style: AppTextStyles.titleLarge.copyWith(
                        color: AppColors.textPrimary,
                      ),
                    ),
                    const Spacer(),
                    TextButton(
                      onPressed: _finish,
                      child: Text(
                        'Skip',
                        style: AppTextStyles.labelMedium.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              Expanded(
                child: PageView.builder(
                  controller: _pageController,
                  itemCount: _pages.length,
                  onPageChanged: (i) => setState(() => _currentPage = i),
                  itemBuilder: (_, i) => _OnboardingPageWidget(
                    controller: _pageController,
                    index: i,
                    page: _pages[i],
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(
                  AppSizes.screenPaddingH,
                  0,
                  AppSizes.screenPaddingH,
                  AppSizes.lg,
                ),
                child: Column(
                  children: [
                    _PageDots(
                      pageCount: _pages.length,
                      currentPage: _currentPage,
                    ),
                    const SizedBox(height: AppSizes.xl),
                    AppButton(
                      label: _currentPage == _pages.length - 1
                          ? 'Start Learning'
                          : 'Next',
                      variant: AppButtonVariant.secondary,
                      onPressed: _next,
                      suffixIcon: Icons.arrow_forward_rounded,
                    ),
                    const SizedBox(height: AppSizes.md),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'Already have an account? ',
                          style: AppTextStyles.bodySmall.copyWith(
                            color: AppColors.textSecondary,
                          ),
                        ),
                        GestureDetector(
                          onTap: _finish,
                          child: Text(
                            'Sign In',
                            style: AppTextStyles.labelMedium.copyWith(
                              color: AppColors.accentPrimary,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _OnboardingPageWidget extends StatelessWidget {
  const _OnboardingPageWidget({
    required this.controller,
    required this.index,
    required this.page,
  });

  final PageController controller;
  final int index;
  final _OnboardingPage page;

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: controller,
      builder: (context, child) {
        final current = controller.hasClients && controller.page != null
            ? controller.page!
            : 0.0;
        final delta = index - current;
        return Transform.translate(
          offset: Offset(delta * 22, 0),
          child: child,
        );
      },
      child: Padding(
        padding:
            const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const MemereMascot(size: Size(296, 260)),
            const SizedBox(height: AppSizes.xl),
            Text(
              page.title,
              style: AppTextStyles.displayMedium.copyWith(
                color: AppColors.textPrimary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: AppSizes.md),
            Text(
              page.body,
              style: AppTextStyles.bodyLarge.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
}

class _PageDots extends StatelessWidget {
  const _PageDots({
    required this.pageCount,
    required this.currentPage,
  });

  final int pageCount;
  final int currentPage;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(
        pageCount,
        (i) => AnimatedContainer(
          duration: AppMotion.base,
          curve: AppMotion.standard,
          margin: const EdgeInsets.symmetric(horizontal: 4),
          width: currentPage == i ? 30 : 8,
          height: 8,
          decoration: BoxDecoration(
            color: currentPage == i
                ? AppColors.accentPrimary
                : AppColors.borderStrong,
            borderRadius: BorderRadius.circular(AppSizes.radiusFull),
          ),
        ),
      ),
    );
  }
}
