import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../core/storage/hive/models/saved_item.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../courses/domain/entities/course_entity.dart';
import '../../../courses/presentation/providers/course_download_provider.dart';
import '../../../payment/domain/entities/enrollment_entity.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../payment/presentation/widgets/enrollment_tile.dart';
import '../../../saved/presentation/providers/saved_courses_provider.dart';

/// "My Courses" tab — the courses the student favorited (local, guest-friendly)
/// plus the courses they're enrolled in (an account feature). Favorites let a
/// guest keep a personal shortlist without signing in; enrolled courses sync
/// across devices once they do.
class MyLearningScreen extends ConsumerWidget {
  const MyLearningScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isAuthenticated =
        ref.watch(authStateProvider).valueOrNull?.isAuthenticated ?? false;
    final favorites = ref.watch(savedCoursesProvider).valueOrNull ?? const [];
    final downloadedCourses = ref.watch(downloadedCoursesProvider);
    // Only touch the enrollment endpoint (which requires auth) for signed-in
    // users; guests still see their local favorites above.
    final enrollmentAsync =
        isAuthenticated ? ref.watch(enrollmentListProvider) : null;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: AppPageBackground(
        child: SafeArea(
          child: RefreshIndicator(
            color: AppColors.accentPrimary,
            backgroundColor: AppColors.bgSecondary,
            onRefresh: () async {
              ref.invalidate(savedCoursesProvider);
              if (isAuthenticated) {
                ref.invalidate(enrollmentListProvider);
                await ref.read(enrollmentListProvider.future);
              }
            },
            child: ListView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.fromLTRB(
                AppSizes.screenPaddingH,
                AppSizes.lg,
                AppSizes.screenPaddingH,
                AppSizes.xl,
              ),
              children: [
                const Text('My Courses', style: AppTextStyles.displayMedium),
                const SizedBox(height: AppSizes.xs),
                Text(
                  'Your favorites and the courses you’re enrolled in.',
                  style: AppTextStyles.bodyMedium.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
                const SizedBox(height: AppSizes.xl),

                // ── Favorites (local, works signed-out) ──
                AppSectionHeader(
                  title: 'Favorites',
                  subtitle: favorites.isEmpty
                      ? null
                      : '${favorites.length} ${favorites.length == 1 ? 'course' : 'courses'}',
                ),
                const SizedBox(height: AppSizes.md),
                if (favorites.isEmpty)
                  const _EmptyFavorites()
                else
                  for (final item in favorites) ...[
                    _FavoriteCourseTile(item: item),
                    const SizedBox(height: AppSizes.sm),
                  ],

                const SizedBox(height: AppSizes.xl),

                // ── Downloaded (available offline, works signed-out) ──
                if (downloadedCourses.isNotEmpty) ...[
                  AppSectionHeader(
                    title: 'Downloaded',
                    subtitle:
                        '${downloadedCourses.length} available offline',
                  ),
                  const SizedBox(height: AppSizes.md),
                  for (final course in downloadedCourses) ...[
                    _DownloadedCourseTile(course: course),
                    const SizedBox(height: AppSizes.sm),
                  ],
                  const SizedBox(height: AppSizes.xl),
                ],

                // ── Enrolled (account feature) ──
                const AppSectionHeader(title: 'Enrolled'),
                const SizedBox(height: AppSizes.md),
                if (!isAuthenticated)
                  const _EnrolledSignInPrompt()
                else
                  ..._buildEnrolled(context, ref, enrollmentAsync!),
              ],
            ),
          ),
        ),
      ),
    );
  }

  List<Widget> _buildEnrolled(
    BuildContext context,
    WidgetRef ref,
    AsyncValue<List<EnrollmentEntity>> async,
  ) {
    return [
      async.when(
        loading: () => const _EnrolledSkeleton(),
        error: (error, _) => _EnrolledMessageCard(
          icon: Icons.error_outline_rounded,
          iconColor: AppColors.error,
          title: 'Could not load your courses',
          body: error is Failure ? error.message : 'Please try again.',
          actionLabel: 'Retry',
          onAction: () => ref.invalidate(enrollmentListProvider),
        ),
        data: (enrollments) {
          if (enrollments.isEmpty) {
            return _EnrolledMessageCard(
              icon: Icons.school_outlined,
              title: 'No courses yet',
              body:
                  'Enroll in a course to start learning — it will show up here.',
              actionLabel: 'Browse courses',
              onAction: () => context.go(AppRoutes.home),
            );
          }
          return Column(
            children: [
              for (var i = 0; i < enrollments.length; i++) ...[
                AppStaggeredReveal(
                  index: i,
                  child: EnrollmentTile(
                    enrollment: enrollments[i],
                    onTap: () => context.push(
                      AppRoutes.courseDetailPath(enrollments[i].courseId),
                    ),
                  ),
                ),
                if (i != enrollments.length - 1)
                  const SizedBox(height: AppSizes.md),
              ],
            ],
          );
        },
      ),
    ];
  }
}

/// A favorited course row — taps through to the course, trailing bookmark
/// removes it from the local favorites store.
class _FavoriteCourseTile extends ConsumerWidget {
  const _FavoriteCourseTile({required this.item});

  final SavedItem item;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final subtitle = item.subtitle;
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.md),
      child: InkWell(
        onTap: () => context.push(AppRoutes.courseDetailPath(item.courseId)),
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        child: Row(
          children: [
            const AppIconTile(
              icon: Icons.menu_book_rounded,
              color: AppColors.accentPrimary,
            ),
            const SizedBox(width: AppSizes.md),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    item.title,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodyMedium.copyWith(
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  if (subtitle != null && subtitle.isNotEmpty) ...[
                    const SizedBox(height: AppSizes.xs),
                    Text(
                      subtitle,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: AppTextStyles.bodySmall.copyWith(
                        color: AppColors.textMuted,
                      ),
                    ),
                  ],
                ],
              ),
            ),
            IconButton(
              tooltip: 'Remove from favorites',
              onPressed: () async {
                await ref
                    .read(savedCoursesProvider.notifier)
                    .remove(item.courseId);
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content:
                          Text('Removed "${item.title}" from My Courses.'),
                    ),
                  );
                }
              },
              icon: const Icon(
                Icons.bookmark_rounded,
                color: AppColors.brandEmerald,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// A downloaded course row — available offline; taps through to the course,
/// which opens from the cached structure when there's no connection.
class _DownloadedCourseTile extends StatelessWidget {
  const _DownloadedCourseTile({required this.course});

  final CourseEntity course;

  @override
  Widget build(BuildContext context) {
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.md),
      child: InkWell(
        onTap: () => context.push(AppRoutes.courseDetailPath(course.id)),
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        child: Row(
          children: [
            const AppIconTile(
              icon: Icons.download_done_rounded,
              color: AppColors.brandEmerald,
            ),
            const SizedBox(width: AppSizes.md),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    course.title,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodyMedium.copyWith(
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  if (course.subject.isNotEmpty) ...[
                    const SizedBox(height: AppSizes.xs),
                    Text(
                      course.subject,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: AppTextStyles.bodySmall.copyWith(
                        color: AppColors.textMuted,
                      ),
                    ),
                  ],
                ],
              ),
            ),
            const Icon(
              Icons.chevron_right_rounded,
              color: AppColors.textMuted,
            ),
          ],
        ),
      ),
    );
  }
}

class _EmptyFavorites extends StatelessWidget {
  const _EmptyFavorites();

  @override
  Widget build(BuildContext context) {
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.lg),
      color: AppColors.bgSecondary,
      child: Row(
        children: [
          const AppIconTile(
            icon: Icons.bookmark_border_rounded,
            color: AppColors.accentPrimary,
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'No favorites yet',
                  style: AppTextStyles.headlineSmall,
                ),
                const SizedBox(height: AppSizes.xs),
                Text(
                  'Tap the bookmark on any course to keep it here for quick access.',
                  style: AppTextStyles.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _EnrolledSignInPrompt extends StatelessWidget {
  const _EnrolledSignInPrompt();

  @override
  Widget build(BuildContext context) {
    return _EnrolledMessageCard(
      icon: Icons.lock_outline_rounded,
      title: 'Sign in to track enrolled courses',
      body:
          'Your enrolled courses sync across devices once you sign in. Favorites '
          'and downloads stay on this device either way.',
      actionLabel: 'Browse courses',
      onAction: () => context.go(AppRoutes.home),
    );
  }
}

/// Compact inline card used for the enrolled section's guest / empty / error
/// states — sized for a list child (unlike the full-bleed PaymentEmptyState).
class _EnrolledMessageCard extends StatelessWidget {
  const _EnrolledMessageCard({
    required this.icon,
    required this.title,
    required this.body,
    required this.actionLabel,
    required this.onAction,
    this.iconColor = AppColors.accentPrimary,
  });

  final IconData icon;
  final String title;
  final String body;
  final String actionLabel;
  final VoidCallback onAction;
  final Color iconColor;

  @override
  Widget build(BuildContext context) {
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.lg),
      color: AppColors.bgSecondary,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              AppIconTile(icon: icon, color: iconColor),
              const SizedBox(width: AppSizes.md),
              Expanded(
                child: Text(title, style: AppTextStyles.headlineSmall),
              ),
            ],
          ),
          const SizedBox(height: AppSizes.sm),
          Text(
            body,
            style: AppTextStyles.bodySmall.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: AppSizes.md),
          Align(
            alignment: Alignment.centerRight,
            child: TextButton(
              onPressed: onAction,
              child: Text(actionLabel),
            ),
          ),
        ],
      ),
    );
  }
}

class _EnrolledSkeleton extends StatelessWidget {
  const _EnrolledSkeleton();

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        for (var i = 0; i < 3; i++) ...[
          Container(
            height: 76,
            decoration: BoxDecoration(
              color: AppColors.bgSecondary,
              borderRadius: BorderRadius.circular(AppSizes.radiusMd),
              border: Border.all(color: AppColors.border),
            ),
          ),
          if (i != 2) const SizedBox(height: AppSizes.sm),
        ],
      ],
    );
  }
}
