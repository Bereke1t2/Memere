import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../providers/auth_state_provider.dart';

class _OnboardingPage {
  const _OnboardingPage(
      {required this.title,
      required this.body,
      required this.icon,
      required this.color});
  final String title, body;
  final IconData icon;
  final Color color;
}

final _pages = [
  const _OnboardingPage(
    title: 'Master Every Subject',
    body:
        'Video lessons from expert teachers covering all Grade 12 subjects for the national exam.',
    icon: Icons.play_circle_fill_rounded,
    color: AppColors.accentPrimary,
  ),
  const _OnboardingPage(
    title: 'Practice Like It\'s Real',
    body:
        'Timed mock exams that simulate the exact national exam format with detailed analytics.',
    icon: Icons.timer_rounded,
    color: AppColors.accentSecondary,
  ),
  const _OnboardingPage(
    title: 'Study Offline, Anywhere',
    body:
        'Download lessons and study without internet — critical for students across Ethiopia.',
    icon: Icons.download_done_rounded,
    color: AppColors.warning,
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
          duration: const Duration(milliseconds: 350), curve: Curves.easeInOut);
    } else {
      _finish();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        child: Column(
          children: [
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                onPressed: _finish,
                child: Text('Skip',
                    style: AppTextStyles.labelMedium
                        .copyWith(color: AppColors.textSecondary)),
              ),
            ),
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                itemCount: _pages.length,
                onPageChanged: (i) => setState(() => _currentPage = i),
                itemBuilder: (_, i) => _OnboardingPageWidget(page: _pages[i]),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: AppSizes.screenPaddingH),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: List.generate(
                      _pages.length,
                      (i) => AnimatedContainer(
                        duration: const Duration(milliseconds: 250),
                        margin: const EdgeInsets.symmetric(horizontal: 4),
                        width: _currentPage == i ? 24 : 8,
                        height: 8,
                        decoration: BoxDecoration(
                          color: _currentPage == i
                              ? AppColors.accentPrimary
                              : AppColors.border,
                          borderRadius: BorderRadius.circular(
                            AppSizes.radiusFull,
                          ),
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(height: AppSizes.xl),
                  AppButton(
                    label: _currentPage == _pages.length - 1
                        ? 'Get Started'
                        : 'Next',
                    onPressed: _next,
                    suffixIcon: Icons.arrow_forward_rounded,
                  ),
                  const SizedBox(height: AppSizes.md),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Text(
                        'Already have an account? ',
                        style: AppTextStyles.bodySmall,
                      ),
                      GestureDetector(
                        onTap: _finish,
                        child: Text('Sign In',
                            style: AppTextStyles.labelMedium
                                .copyWith(color: AppColors.accentPrimary)),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSizes.lg),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _OnboardingPageWidget extends StatelessWidget {
  const _OnboardingPageWidget({required this.page});
  final _OnboardingPage page;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 120,
            height: 120,
            decoration: BoxDecoration(
              color: page.color.withAlpha(30),
              shape: BoxShape.circle,
              border: Border.all(color: page.color.withAlpha(76), width: 1.5),
            ),
            child: Icon(page.icon, size: 54, color: page.color),
          ),
          const SizedBox(height: AppSizes.xl),
          Text(page.title,
              style: AppTextStyles.headlineLarge, textAlign: TextAlign.center),
          const SizedBox(height: AppSizes.md),
          Text(
            page.body,
            style: AppTextStyles.bodyLarge
                .copyWith(color: AppColors.textSecondary),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}
