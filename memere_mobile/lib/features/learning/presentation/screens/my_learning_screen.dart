import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../payment/presentation/widgets/enrollment_tile.dart';
import '../../../payment/presentation/widgets/payment_empty_state.dart';

/// "Learn" tab — the student's enrolled courses with a tap-through to continue.
class MyLearningScreen extends ConsumerWidget {
  const MyLearningScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isAuthenticated =
        ref.watch(authStateProvider).valueOrNull?.isAuthenticated ?? false;
    // Enrollments are an account feature (the endpoint requires auth). Guests
    // get a sign-in prompt; their downloads live under Saved ▸ Downloaded.
    if (!isAuthenticated) {
      return Scaffold(
        backgroundColor: AppColors.bgPrimary,
        body: AppPageBackground(
          child: SafeArea(
            child: PaymentEmptyState(
              icon: Icons.school_outlined,
              title: 'Track your courses',
              body:
                  'Sign in to enroll and pick up where you left off across '
                  'devices. You can still browse and download courses to study '
                  'offline without an account.',
              buttonLabel: 'Browse courses',
              onPressed: () => context.go(AppRoutes.home),
            ),
          ),
        ),
      );
    }

    final async = ref.watch(enrollmentListProvider);

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: AppPageBackground(
        child: SafeArea(
          child: async.when(
            loading: () => const _LearningSkeleton(),
            error: (error, _) => PaymentEmptyState(
              icon: Icons.error_outline_rounded,
              title: 'Could not load your courses',
              body: error is Failure ? error.message : 'Please try again.',
              buttonLabel: 'Retry',
              onPressed: () => ref.invalidate(enrollmentListProvider),
              iconColor: AppColors.error,
            ),
            data: (enrollments) {
              if (enrollments.isEmpty) {
                return PaymentEmptyState(
                  icon: Icons.school_outlined,
                  title: 'No courses yet',
                  body:
                      'Enroll in a course to start learning. Your courses will '
                      'show up here.',
                  buttonLabel: 'Browse courses',
                  onPressed: () => context.go(AppRoutes.home),
                );
              }
              return RefreshIndicator(
                color: AppColors.accentPrimary,
                backgroundColor: AppColors.bgSecondary,
                onRefresh: () async {
                  ref.invalidate(enrollmentListProvider);
                  await ref.read(enrollmentListProvider.future);
                },
                child: ListView.separated(
                  physics: const AlwaysScrollableScrollPhysics(),
                  padding: const EdgeInsets.fromLTRB(
                    AppSizes.screenPaddingH,
                    AppSizes.md,
                    AppSizes.screenPaddingH,
                    AppSizes.xl,
                  ),
                  itemCount: enrollments.length + 1,
                  separatorBuilder: (_, __) =>
                      const SizedBox(height: AppSizes.md),
                  itemBuilder: (_, index) {
                    if (index == 0) {
                      return _LearningHeader(courseCount: enrollments.length);
                    }
                    final enrollment = enrollments[index - 1];
                    return AppStaggeredReveal(
                      index: index - 1,
                      child: EnrollmentTile(
                        enrollment: enrollment,
                        onTap: () => context.push(
                          AppRoutes.courseDetailPath(enrollment.courseId),
                        ),
                      ),
                    );
                  },
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}

class _LearningHeader extends StatelessWidget {
  const _LearningHeader({required this.courseCount});

  final int courseCount;

  @override
  Widget build(BuildContext context) {
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.md),
      gradient: AppColors.cardGradient,
      shadows: AppShadows.md,
      child: Row(
        children: [
          const AppIconTile(
            icon: Icons.school_rounded,
            gradient: AppColors.primaryGradient,
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('My learning', style: AppTextStyles.headlineMedium),
                const SizedBox(height: AppSizes.xs),
                Text(
                  courseCount == 1
                      ? '1 active course'
                      : '$courseCount active courses',
                  style: AppTextStyles.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          const AppBadge(
            label: 'Continue',
            color: AppColors.success,
            icon: Icons.play_arrow_rounded,
          ),
        ],
      ),
    );
  }
}

class _LearningSkeleton extends StatelessWidget {
  const _LearningSkeleton();

  @override
  Widget build(BuildContext context) {
    return ListView.separated(
      padding: const EdgeInsets.all(AppSizes.screenPaddingH),
      itemCount: 6,
      separatorBuilder: (_, __) => const SizedBox(height: AppSizes.sm),
      itemBuilder: (_, __) => Container(
        height: 76,
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          border: Border.all(color: AppColors.border),
        ),
      ),
    );
  }
}
