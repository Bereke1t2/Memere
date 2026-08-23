import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/router/app_router.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../progress/domain/entities/student_points_entity.dart';
import '../../../progress/presentation/providers/progress_providers.dart';

/// Clean, simple, mature, and professional Profile & Account Screen for Memere.
///
/// Design Language:
/// - Dark Obsidian surfaces with subtle borders matching the Home, Exam, and Learning Hubs
/// - Clean hero student card with initials avatar and active stream pill
/// - Simple 3-stat metric strip
/// - Unified dark neutral settings groups with subtle chevrons
/// - Zero colorful distractions or candy mascots
class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authValue = ref.watch(authStateProvider).valueOrNull;
    if (!(authValue?.isAuthenticated ?? false)) {
      return const _GuestProfileView();
    }
    final user = authValue?.user;
    final fullName = [
      user?.firstName.trim() ?? '',
      user?.lastName.trim() ?? '',
    ].where((part) => part.isNotEmpty).join(' ');

    final enrollmentsAsync = ref.watch(enrollmentListProvider);
    final enrolledCount = enrollmentsAsync.valueOrNull?.length ?? 0;

    final purchasesAsync = ref.watch(paymentHistoryProvider);
    final purchasesCount = purchasesAsync.valueOrNull?.length ?? 0;

    final pointsAsync = ref.watch(studentPointsProvider);

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        child: RefreshIndicator(
          color: AppColors.brandEmerald,
          backgroundColor: AppColors.bgSecondary,
          onRefresh: () async {
            ref.invalidate(authStateProvider);
            ref.invalidate(enrollmentListProvider);
            ref.invalidate(paymentHistoryProvider);
            ref.invalidate(studentPointsProvider);
          },
          child: ListView(
            padding: const EdgeInsets.fromLTRB(
              AppSizes.screenPaddingH,
              AppSizes.sm,
              AppSizes.screenPaddingH,
              AppSizes.xxl,
            ),
            children: [
              // 1. Sleek Top Bar matching Home & Learning Hubs
              _ProfileTopBar(
                canPop: context.canPop(),
                onPop: () {
                  if (context.canPop()) {
                    context.pop();
                  } else {
                    context.go(AppRoutes.home);
                  }
                },
              ),
              const SizedBox(height: 16),

              // 2. Refined Hero Student Profile Card
              _StudentHeroCard(
                name: fullName.isEmpty ? 'Student' : fullName,
                email: user?.email ?? '',
                phoneNumber: user?.phone ?? '',
                role: user?.role.name ?? 'student',
                initials: _initials(user?.firstName, user?.lastName),
              ),
              const SizedBox(height: 14),

              // 3. Points Card (Auth-only, synced points)
              ...pointsAsync.when(
                data: (points) => [
                  _PointsCard(points: points),
                  const SizedBox(height: 14),
                ],
                loading: () => const [
                  _PointsCardSkeleton(),
                  SizedBox(height: 14),
                ],
                error: (_, __) => const <Widget>[],
              ),

              // 4. Clean Academic Metric Strip
              _LearningStatsRow(
                enrolledCount: enrolledCount,
                purchasesCount: purchasesCount,
              ),
              const SizedBox(height: 20),

              // 5. Academic & Learning Section
              const _SectionHeader(title: 'Academic & Learning'),
              const SizedBox(height: 6),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.school_outlined,
                    title: 'My Enrolled Courses',
                    subtitle: enrolledCount == 0
                        ? 'Explore curriculum courses'
                        : '$enrolledCount active courses in progress',
                    onTap: () => context.go(AppRoutes.learn),
                  ),
                  _SettingsItemData(
                    icon: Icons.assignment_outlined,
                    title: 'Mock Exams & Results',
                    subtitle: 'National entrance exams and score analytics',
                    onTap: () => context.go(AppRoutes.mockExams),
                  ),
                  _SettingsItemData(
                    icon: Icons.bookmark_outline_rounded,
                    title: 'Saved Notes & Library',
                    subtitle: 'Offline study guides and saved materials',
                    onTap: () => context.go(AppRoutes.saved),
                  ),
                ],
              ),
              const SizedBox(height: 18),

              // 6. Membership & Billing Section
              const _SectionHeader(title: 'Membership & Billing'),
              const SizedBox(height: 6),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.workspace_premium_outlined,
                    title: 'All-Access Plans',
                    subtitle: 'Unlock unlimited mock exams & full solutions',
                    onTap: () => context.push(AppRoutes.subscriptionPlans),
                  ),
                  _SettingsItemData(
                    icon: Icons.receipt_outlined,
                    title: 'Purchase History',
                    subtitle: purchasesCount == 0
                        ? 'View transaction history and invoices'
                        : '$purchasesCount recorded transactions',
                    onTap: () => context.push(AppRoutes.purchaseHistory),
                  ),
                ],
              ),
              const SizedBox(height: 18),

              // 7. Preferences & Support Section
              const _SectionHeader(title: 'Preferences & Support'),
              const SizedBox(height: 6),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.grid_view_rounded,
                    title: 'Curriculum Stream',
                    subtitle: 'Natural Science (Grade 12 EUEE)',
                    onTap: () {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text(
                            'Curriculum Stream: Natural Science (Grade 12)',
                          ),
                          behavior: SnackBarBehavior.floating,
                        ),
                      );
                    },
                  ),
                  _SettingsItemData(
                    icon: Icons.notifications_none_rounded,
                    title: 'Study Notifications',
                    subtitle: 'Daily study schedules and mock announcements',
                    onTap: () {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content:
                              Text('Study notifications are currently enabled.'),
                          behavior: SnackBarBehavior.floating,
                        ),
                      );
                    },
                  ),
                  _SettingsItemData(
                    icon: Icons.help_outline_rounded,
                    title: 'Help & Support',
                    subtitle: 'FAQs, contact instructors, report an issue',
                    onTap: () => _showHelpDialog(context),
                  ),
                ],
              ),
              const SizedBox(height: 18),

              // 8. Account Session Section
              const _SectionHeader(title: 'Account Session'),
              const SizedBox(height: 6),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.logout_rounded,
                    title: 'Sign Out',
                    subtitle: 'Log out of your account on this device',
                    isDestructive: true,
                    onTap: () => _confirmSignOut(context, ref),
                  ),
                ],
              ),
              const SizedBox(height: 28),

              // 9. Footer Brand & Version
              const Center(
                child: Column(
                  children: [
                    Text(
                      'Memere • Ethiopian University Entrance Exam Prep',
                      style: TextStyle(
                        fontSize: 11.5,
                        fontWeight: FontWeight.w600,
                        color: AppColors.textMuted,
                      ),
                    ),
                    SizedBox(height: 4),
                    Text(
                      'Version 1.0.0 (Build 42)',
                      style: TextStyle(
                        fontSize: 10.5,
                        color: AppColors.textDisabled,
                      ),
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

  String _initials(String? first, String? last) {
    final f = (first ?? '').trim();
    final l = (last ?? '').trim();
    final buffer = StringBuffer();
    if (f.isNotEmpty) buffer.write(f[0]);
    if (l.isNotEmpty) buffer.write(l[0]);
    final result = buffer.toString().toUpperCase();
    return result.isEmpty ? 'S' : result;
  }

  Future<void> _confirmSignOut(BuildContext context, WidgetRef ref) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        backgroundColor: AppColors.bgSecondary,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
          side: const BorderSide(color: AppColors.borderStrong),
        ),
        title: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(8),
              decoration: const BoxDecoration(
                color: Color(0x1DEF4444),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.logout_rounded,
                color: AppColors.error,
                size: 20,
              ),
            ),
            const SizedBox(width: 12),
            const Text(
              'Sign Out',
              style: TextStyle(
                fontSize: 17,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
          ],
        ),
        content: const Text(
          'Are you sure you want to sign out? You can sign back in anytime to access your courses, mock exams, and saved study materials.',
          style: TextStyle(
            fontSize: 13.5,
            color: AppColors.textSecondary,
            height: 1.45,
          ),
        ),
        actionsPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text(
              'Cancel',
              style: TextStyle(
                color: AppColors.textMuted,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
              foregroundColor: Colors.white,
              elevation: 0,
              padding:
                  const EdgeInsets.symmetric(horizontal: 18, vertical: 10),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(10),
              ),
            ),
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: const Text(
              'Sign Out',
              style: TextStyle(fontWeight: FontWeight.w700),
            ),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    await ref.read(authStateProvider.notifier).logout();
    if (context.mounted) context.go(AppRoutes.login);
  }

  void _showHelpDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (dialogContext) => AlertDialog(
        backgroundColor: AppColors.bgSecondary,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
          side: const BorderSide(color: AppColors.borderStrong),
        ),
        title: const Text(
          'Help & Support',
          style: TextStyle(
            fontSize: 17,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        content: const Text(
          'Need help with course materials, mock exams, or payment confirmations? Contact the Memere academic support team at support@memere.et.',
          style: TextStyle(
            fontSize: 13.5,
            color: AppColors.textSecondary,
            height: 1.45,
          ),
        ),
        actionsPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        actions: [
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.brandEmerald,
              foregroundColor: Colors.white,
              elevation: 0,
              padding:
                  const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(10),
              ),
            ),
            onPressed: () => Navigator.of(dialogContext).pop(),
            child: const Text(
              'Done',
              style: TextStyle(fontWeight: FontWeight.w700),
            ),
          ),
        ],
      ),
    );
  }
}

/// Sleek Top Bar matching Home & Learning Hubs
class _ProfileTopBar extends StatelessWidget {
  const _ProfileTopBar({
    required this.canPop,
    required this.onPop,
  });

  final bool canPop;
  final VoidCallback onPop;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Row(
          children: [
            if (canPop) ...[
              InkWell(
                onTap: onPop,
                borderRadius: BorderRadius.circular(20),
                child: Container(
                  width: 36,
                  height: 36,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    shape: BoxShape.circle,
                    border: Border.all(color: AppColors.borderStrong),
                  ),
                  child: const Icon(
                    Icons.arrow_back_ios_new_rounded,
                    color: AppColors.textPrimary,
                    size: 14,
                  ),
                ),
              ),
              const SizedBox(width: 10),
            ] else ...[
              const Icon(
                Icons.person_rounded,
                color: AppColors.brandEmerald,
                size: 22,
              ),
              const SizedBox(width: 8),
            ],
            const Text(
              'Profile & Account',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w900,
                color: AppColors.textPrimary,
                letterSpacing: -0.4,
              ),
            ),
          ],
        ),

        // Active Student Status Pill
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
          decoration: BoxDecoration(
            color: AppColors.bgSecondary,
            borderRadius: BorderRadius.circular(16),
            border: Border.all(color: AppColors.borderStrong),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 6,
                height: 6,
                decoration: const BoxDecoration(
                  shape: BoxShape.circle,
                  color: AppColors.brandEmerald,
                ),
              ),
              const SizedBox(width: 6),
              const Text(
                'Active Student',
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

/// Refined Hero Student Profile Card
class _StudentHeroCard extends StatelessWidget {
  const _StudentHeroCard({
    required this.name,
    required this.email,
    required this.phoneNumber,
    required this.role,
    required this.initials,
  });

  final String name;
  final String email;
  final String phoneNumber;
  final String role;
  final String initials;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(18),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          // Sleek Emerald-Tinted Avatar Circle
          Stack(
            clipBehavior: Clip.none,
            children: [
              Container(
                width: 54,
                height: 54,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: AppColors.brandEmerald.withAlpha(25),
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: AppColors.brandEmerald.withAlpha(80),
                    width: 1.5,
                  ),
                ),
                child: Text(
                  initials,
                  style: const TextStyle(
                    fontSize: 19,
                    fontWeight: FontWeight.w800,
                    color: AppColors.brandEmerald,
                    letterSpacing: -0.5,
                  ),
                ),
              ),
              Positioned(
                bottom: 0,
                right: 0,
                child: Container(
                  width: 15,
                  height: 15,
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    shape: BoxShape.circle,
                    border: Border.all(
                      color: AppColors.bgSecondary,
                      width: 1.5,
                    ),
                  ),
                  child: const Icon(
                    Icons.check_circle_rounded,
                    size: 13,
                    color: AppColors.brandEmerald,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(width: 14),

          // Name, Email, & Stream Pill
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  name,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w700,
                    color: AppColors.textPrimary,
                    letterSpacing: -0.3,
                  ),
                ),
                const SizedBox(height: 2),
                if (email.isNotEmpty)
                  Text(
                    email,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textMuted,
                    ),
                  ),
                if (phoneNumber.isNotEmpty) ...[
                  const SizedBox(height: 1),
                  Text(
                    phoneNumber,
                    style: const TextStyle(
                      fontSize: 11,
                      color: AppColors.textDisabled,
                    ),
                  ),
                ],
                const SizedBox(height: 6),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(
                    color: AppColors.bgTertiary,
                    borderRadius: BorderRadius.circular(6),
                    border: Border.all(color: AppColors.borderStrong),
                  ),
                  child: const Text(
                    'GRADE 12 • NATURAL SCIENCE',
                    style: TextStyle(
                      fontSize: 9.5,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF38BDF8),
                      letterSpacing: 0.4,
                    ),
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

/// Cumulative Points Card (Auth Only)
class _PointsCard extends StatelessWidget {
  const _PointsCard({required this.points});

  final StudentPointsEntity points;

  @override
  Widget build(BuildContext context) {
    if (points.isEmpty) return const _PointsEmptyCard();

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(18),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: AppColors.brandEmerald.withAlpha(25),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.brandEmerald.withAlpha(60)),
            ),
            child: const Icon(
              Icons.military_tech_rounded,
              size: 22,
              color: AppColors.brandEmerald,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'TOTAL POINTS',
                  style: TextStyle(
                    fontSize: 10,
                    fontWeight: FontWeight.w700,
                    color: AppColors.textMuted,
                    letterSpacing: 0.8,
                  ),
                ),
                const SizedBox(height: 2),
                Row(
                  crossAxisAlignment: CrossAxisAlignment.baseline,
                  textBaseline: TextBaseline.alphabetic,
                  children: [
                    Text(
                      _formatNumber(points.totalPoints),
                      style: const TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.w800,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.5,
                      ),
                    ),
                    const SizedBox(width: 5),
                    const Text(
                      'pts',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                        color: AppColors.textMuted,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  _breakdownLabel(points),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    fontSize: 11,
                    color: AppColors.textMuted,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 10),

          // Average score pill badge
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
            decoration: BoxDecoration(
              color: AppColors.bgTertiary,
              borderRadius: BorderRadius.circular(10),
              border: Border.all(color: AppColors.borderStrong),
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  '${_formatNumber(points.avgPercentage)}%',
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w800,
                    color: AppColors.brandEmerald,
                  ),
                ),
                const SizedBox(height: 1),
                const Text(
                  'avg score',
                  style: TextStyle(
                    fontSize: 9,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textMuted,
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

class _PointsEmptyCard extends StatelessWidget {
  const _PointsEmptyCard();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(18),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: AppColors.bgTertiary,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.borderStrong),
            ),
            child: const Icon(
              Icons.military_tech_outlined,
              size: 22,
              color: AppColors.textMuted,
            ),
          ),
          const SizedBox(width: 12),
          const Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Earn your first points',
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w700,
                    color: AppColors.textPrimary,
                  ),
                ),
                SizedBox(height: 2),
                Text(
                  'Complete quizzes and mock exams to start building your score.',
                  style: TextStyle(
                    fontSize: 11.5,
                    height: 1.35,
                    color: AppColors.textMuted,
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

class _PointsCardSkeleton extends StatelessWidget {
  const _PointsCardSkeleton();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(18),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: AppColors.bgTertiary,
              borderRadius: BorderRadius.circular(12),
            ),
          ),
          const SizedBox(width: 12),
          const Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _SkeletonBar(width: 80, height: 10),
                SizedBox(height: 8),
                _SkeletonBar(width: 60, height: 18),
                SizedBox(height: 8),
                _SkeletonBar(width: 120, height: 10),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _SkeletonBar extends StatelessWidget {
  const _SkeletonBar({required this.width, required this.height});

  final double width;
  final double height;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        color: AppColors.bgTertiary,
        borderRadius: BorderRadius.circular(6),
      ),
    );
  }
}

String _formatNumber(double value) {
  if (value == value.roundToDouble()) return value.toInt().toString();
  return value.toStringAsFixed(1);
}

String _breakdownLabel(StudentPointsEntity p) {
  final quizzes = '${p.quizCount} ${p.quizCount == 1 ? 'quiz' : 'quizzes'}';
  final exams = '${p.examCount} ${p.examCount == 1 ? 'exam' : 'exams'}';
  return '$quizzes • $exams';
}

/// Academic Metric Strip (Courses, Mock Exams, Saved Library)
class _LearningStatsRow extends StatelessWidget {
  const _LearningStatsRow({
    required this.enrolledCount,
    required this.purchasesCount,
  });

  final int enrolledCount;
  final int purchasesCount;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: _StatCard(
            icon: Icons.menu_book_outlined,
            label: 'Courses',
            value: '$enrolledCount Enrolled',
            onTap: () => context.go(AppRoutes.learn),
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: _StatCard(
            icon: Icons.assignment_outlined,
            label: 'Mock Prep',
            value: 'National Exams',
            onTap: () => context.go(AppRoutes.mockExams),
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: _StatCard(
            icon: Icons.bookmark_outline_rounded,
            label: 'Library',
            value: 'Saved Notes',
            onTap: () => context.go(AppRoutes.saved),
          ),
        ),
      ],
    );
  }
}

class _StatCard extends StatelessWidget {
  const _StatCard({
    required this.icon,
    required this.label,
    required this.value,
    required this.onTap,
  });

  final IconData icon;
  final String label;
  final String value;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(14),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 12),
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(14),
          border: Border.all(color: AppColors.borderStrong),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(icon, size: 18, color: AppColors.textMuted),
            const SizedBox(height: 8),
            Text(
              value,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(
                fontSize: 12.5,
                fontWeight: FontWeight.w700,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: 2),
            Text(
              label,
              style: const TextStyle(
                fontSize: 10.5,
                color: AppColors.textMuted,
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _SectionHeader extends StatelessWidget {
  const _SectionHeader({required this.title});

  final String title;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 2, bottom: 4),
      child: Text(
        title.toUpperCase(),
        style: const TextStyle(
          fontSize: 10.5,
          fontWeight: FontWeight.w700,
          color: AppColors.textMuted,
          letterSpacing: 0.8,
        ),
      ),
    );
  }
}

class _SettingsItemData {
  const _SettingsItemData({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
    this.isDestructive = false,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;
  final bool isDestructive;
}

class _SettingsGroup extends StatelessWidget {
  const _SettingsGroup({required this.items});

  final List<_SettingsItemData> items;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.borderStrong),
      ),
      clipBehavior: Clip.antiAlias,
      child: ListView.separated(
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        padding: EdgeInsets.zero,
        itemCount: items.length,
        separatorBuilder: (_, __) =>
            const Divider(height: 1, color: AppColors.border),
        itemBuilder: (context, index) {
          final item = items[index];
          final Color iconColor =
              item.isDestructive ? AppColors.error : AppColors.textSecondary;
          final Color iconBg = item.isDestructive
              ? const Color(0x1DEF4444)
              : AppColors.bgTertiary;
          final Color iconBorder = item.isDestructive
              ? const Color(0x35EF4444)
              : AppColors.borderStrong;

          return InkWell(
            onTap: item.onTap,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
              child: Row(
                children: [
                  Container(
                    width: 36,
                    height: 36,
                    alignment: Alignment.center,
                    decoration: BoxDecoration(
                      color: iconBg,
                      borderRadius: BorderRadius.circular(10),
                      border: Border.all(color: iconBorder),
                    ),
                    child: Icon(item.icon, size: 18, color: iconColor),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          item.title,
                          style: TextStyle(
                            fontSize: 13.5,
                            fontWeight: FontWeight.w600,
                            color: item.isDestructive
                                ? AppColors.error
                                : AppColors.textPrimary,
                            letterSpacing: -0.2,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          item.subtitle,
                          style: const TextStyle(
                            fontSize: 11,
                            color: AppColors.textMuted,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Icon(
                    Icons.chevron_right_rounded,
                    size: 18,
                    color: item.isDestructive
                        ? AppColors.error.withAlpha(150)
                        : AppColors.textDisabled,
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

/// Guest Profile View (Unauthenticated)
class _GuestProfileView extends StatelessWidget {
  const _GuestProfileView();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(
            AppSizes.screenPaddingH,
            AppSizes.sm,
            AppSizes.screenPaddingH,
            AppSizes.xxl,
          ),
          children: [
            const Row(
              children: [
                Icon(
                  Icons.person_rounded,
                  color: AppColors.brandEmerald,
                  size: 22,
                ),
                SizedBox(width: 8),
                Text(
                  'Profile & Account',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w900,
                    color: AppColors.textPrimary,
                    letterSpacing: -0.4,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // Guest Hero Card
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(18),
              decoration: BoxDecoration(
                color: AppColors.bgSecondary,
                borderRadius: BorderRadius.circular(18),
                border: Border.all(color: AppColors.borderStrong),
              ),
              child: Column(
                children: [
                  Container(
                    width: 56,
                    height: 56,
                    alignment: Alignment.center,
                    decoration: BoxDecoration(
                      color: AppColors.bgTertiary,
                      shape: BoxShape.circle,
                      border: Border.all(color: AppColors.borderStrong),
                    ),
                    child: const Icon(
                      Icons.person_outline_rounded,
                      size: 26,
                      color: AppColors.textMuted,
                    ),
                  ),
                  const SizedBox(height: 14),
                  const Text(
                    "You're browsing as a guest",
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 6),
                  const Text(
                    'Sign in to enroll in courses, sync your progress across '
                    'devices, and keep your purchases.',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 12.5,
                      height: 1.4,
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppColors.brandEmerald,
                        foregroundColor: Colors.white,
                        elevation: 0,
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      onPressed: () => context.push(AppRoutes.login),
                      child: const Text(
                        'Sign in',
                        style: TextStyle(
                            fontWeight: FontWeight.w700, fontSize: 14),
                      ),
                    ),
                  ),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: OutlinedButton(
                      style: OutlinedButton.styleFrom(
                        foregroundColor: AppColors.textPrimary,
                        side: const BorderSide(color: AppColors.borderStrong),
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      onPressed: () => context.push(AppRoutes.register),
                      child: const Text(
                        'Create account',
                        style: TextStyle(
                            fontWeight: FontWeight.w700, fontSize: 14),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Browse options available without an account
            const _SectionHeader(title: 'Explore without an account'),
            const SizedBox(height: 6),
            _SettingsGroup(
              items: [
                _SettingsItemData(
                  icon: Icons.menu_book_outlined,
                  title: 'Browse Courses',
                  subtitle: 'Explore the full curriculum catalog',
                  onTap: () => context.go(AppRoutes.home),
                ),
                _SettingsItemData(
                  icon: Icons.assignment_outlined,
                  title: 'Mock Exams',
                  subtitle: 'Take national entrance mock exams',
                  onTap: () => context.go(AppRoutes.mockExams),
                ),
                _SettingsItemData(
                  icon: Icons.bookmark_outline_rounded,
                  title: 'Saved & Downloaded',
                  subtitle: 'Study offline — no account needed',
                  onTap: () => context.go(AppRoutes.saved),
                ),
              ],
            ),
            const SizedBox(height: 28),
            const Center(
              child: Text(
                'Memere • Ethiopian University Entrance Exam Prep',
                style: TextStyle(
                  fontSize: 11.5,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textMuted,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
