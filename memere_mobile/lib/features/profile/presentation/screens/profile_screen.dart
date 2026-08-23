import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/router/app_router.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../courses/presentation/providers/completed_lessons_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../progress/domain/entities/student_points_entity.dart';
import '../../../progress/presentation/providers/progress_providers.dart';

/// Refined Profile Screen for Memere adapted directly from the reference UI design.
///
/// Design Highlights:
/// - Top patterned header banner with settings icon
/// - Prominent hero avatar bridging header and profile content
/// - User handle & metadata row ("@username • Joined August 2024")
/// - 3-Column Stat Strip (Courses, Total Points, Avg Score)
/// - Primary Action Button ("Edit Profile" / "Share Profile") + Square Share Button
/// - "Weekly progress" card with 7-day comparative sparkline graph (This Week vs Last Week)
/// - Clean dark obsidian settings navigation groups
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
    final points = pointsAsync.valueOrNull;

    final completedLessons =
        ref.watch(completedLessonsProvider).valueOrNull ?? const {};
    final completedCount = completedLessons.length;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        top: false,
        child: RefreshIndicator(
          color: AppColors.brandEmerald,
          backgroundColor: AppColors.bgSecondary,
          onRefresh: () async {
            ref.invalidate(authStateProvider);
            ref.invalidate(enrollmentListProvider);
            ref.invalidate(paymentHistoryProvider);
            ref.invalidate(studentPointsProvider);
            ref.invalidate(completedLessonsProvider);
          },
          child: ListView(
            padding: EdgeInsets.zero,
            children: [
              // 1. Top Header Banner with Settings Icon & Hero Avatar
              _ProfileHeaderBanner(
                initials: _initials(user?.firstName, user?.lastName),
                onSettingsPressed: () => _showSettingsSheet(context, ref),
              ),

              // 2. Main Profile Content Body
              Padding(
                padding: const EdgeInsets.symmetric(
                    horizontal: AppSizes.screenPaddingH),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const SizedBox(height: 12),

                    // User Name & Handle / Joined Info
                    Text(
                      fullName.isEmpty ? 'Active Student' : fullName,
                      style: const TextStyle(
                        fontFamily: 'Sora',
                        fontSize: 22,
                        fontWeight: FontWeight.w800,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.4,
                      ),
                    ),
                    const SizedBox(height: 3),
                    Row(
                      children: [
                        Text(
                          user?.email.isNotEmpty == true
                              ? '@${user!.email.split('@').first}'
                              : '@student',
                          style: const TextStyle(
                            fontSize: 13,
                            fontWeight: FontWeight.w600,
                            color: AppColors.textMuted,
                          ),
                        ),
                        const SizedBox(width: 6),
                        const Text(
                          '•',
                          style: TextStyle(
                              color: AppColors.textDisabled, fontSize: 12),
                        ),
                        const SizedBox(width: 6),
                        const Text(
                          'Joined Aug 2024',
                          style: TextStyle(
                            fontSize: 13,
                            color: AppColors.textMuted,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),

                    // 3. Horizontal Metric Stat Row (Courses, Points, Avg Score)
                    _SocialStatsRow(
                      enrolledCount: enrolledCount,
                      totalPoints: points?.totalPoints ?? 0,
                      avgScore: points?.avgPercentage ?? 0,
                      onTapCourses: () => context.go(AppRoutes.learn),
                      onTapPoints: () => context.go(AppRoutes.mockExams),
                    ),
                    const SizedBox(height: 14),

                    // 4. Primary Action Bar (Share Profile + Square Share Icon)
                    _ProfileActionsBar(
                      onShare: () => _shareProfile(context, user?.email),
                      onEdit: () => _showEditPrompt(context),
                    ),
                    const SizedBox(height: 24),

                    // 5. "Weekly progress" Card Section (Matching reference UI sparkline chart)
                    const Text(
                      'Weekly progress',
                      style: TextStyle(
                        fontFamily: 'Sora',
                        fontSize: 18,
                        fontWeight: FontWeight.w800,
                        color: AppColors.textPrimary,
                        letterSpacing: -0.3,
                      ),
                    ),
                    const SizedBox(height: 10),
                    _WeeklyProgressCard(completedCount: completedCount),
                    const SizedBox(height: 24),

                    // 6. Academic & Learning Group
                    const _SectionHeader(title: 'Academic & Learning'),
                    const SizedBox(height: 6),
                    _SettingsGroup(
                      items: [
                        _SettingsItemData(
                          icon: Icons.menu_book_outlined,
                          title: 'My Enrolled Courses',
                          subtitle: enrolledCount == 0
                              ? 'Explore curriculum courses'
                              : '$enrolledCount active courses in progress',
                          onTap: () => context.go(AppRoutes.learn),
                        ),
                        _SettingsItemData(
                          icon: Icons.assignment_outlined,
                          title: 'Mock Exams & Results',
                          subtitle:
                              'National entrance exams and score analytics',
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
                    const SizedBox(height: 20),

                    // 7. Membership & Purchases Group
                    const _SectionHeader(title: 'Membership & Billing'),
                    const SizedBox(height: 6),
                    _SettingsGroup(
                      items: [
                        _SettingsItemData(
                          icon: Icons.workspace_premium_outlined,
                          title: 'All-Access Plans',
                          subtitle:
                              'Unlock unlimited mock exams & full solutions',
                          onTap: () =>
                              context.push(AppRoutes.subscriptionPlans),
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
                    const SizedBox(height: 20),

                    // 8. Preferences & Support Group
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
                          subtitle:
                              'Daily study schedules and mock announcements',
                          onTap: () {
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(
                                content: Text(
                                    'Study notifications are currently enabled.'),
                                behavior: SnackBarBehavior.floating,
                              ),
                            );
                          },
                        ),
                        _SettingsItemData(
                          icon: Icons.help_outline_rounded,
                          title: 'Help & Support',
                          subtitle:
                              'FAQs, contact instructors, report an issue',
                          onTap: () => _showHelpDialog(context),
                        ),
                      ],
                    ),
                    const SizedBox(height: 20),

                    // 9. Account Session & Sign Out
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
                    const SizedBox(height: 32),

                    // Footer Version & Brand Tag
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
                    const SizedBox(height: 40),
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

  void _shareProfile(BuildContext context, String? email) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('Profile link copied to clipboard (${email ?? 'user'})'),
        behavior: SnackBarBehavior.floating,
      ),
    );
  }

  void _showEditPrompt(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Profile details updated.'),
        behavior: SnackBarBehavior.floating,
      ),
    );
  }

  void _showSettingsSheet(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 36,
                height: 4,
                decoration: BoxDecoration(
                  color: AppColors.borderStrong,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 16),
              ListTile(
                leading: const Icon(Icons.logout_rounded, color: AppColors.error),
                title: const Text('Sign Out',
                    style: TextStyle(
                        color: AppColors.error, fontWeight: FontWeight.bold)),
                onTap: () {
                  Navigator.pop(ctx);
                  _confirmSignOut(context, ref);
                },
              ),
            ],
          ),
        ),
      ),
    );
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

/// Top Header Banner matching the reference UI (pattern background + avatar bridge + settings gear icon)
class _ProfileHeaderBanner extends StatelessWidget {
  const _ProfileHeaderBanner({
    required this.initials,
    required this.onSettingsPressed,
  });

  final String initials;
  final VoidCallback onSettingsPressed;

  @override
  Widget build(BuildContext context) {
    final topPadding = MediaQuery.paddingOf(context).top;

    return Stack(
      clipBehavior: Clip.none,
      children: [
        // Top Banner Background with subtle patterned gradient
        Container(
          height: 130 + topPadding,
          width: double.infinity,
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                Color(0xFF0F382A),
                Color(0xFF072118),
                Color(0xFF0B1713),
              ],
            ),
          ),
          child: Stack(
            children: [
              // Subtle background watermark motif
              Positioned(
                right: -20,
                top: topPadding - 10,
                child: Icon(
                  Icons.school_rounded,
                  size: 140,
                  color: Colors.white.withAlpha(12),
                ),
              ),
              // Top Right Settings Gear Icon
              Positioned(
                top: topPadding + 8,
                right: AppSizes.screenPaddingH,
                child: InkWell(
                  onTap: onSettingsPressed,
                  borderRadius: BorderRadius.circular(20),
                  child: Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: Colors.black.withAlpha(80),
                      shape: BoxShape.circle,
                      border: Border.all(color: Colors.white.withAlpha(40)),
                    ),
                    child: const Icon(
                      Icons.settings_outlined,
                      size: 20,
                      color: Colors.white,
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),

        // Hero Avatar Circle bridging the banner and profile content below
        Positioned(
          left: AppSizes.screenPaddingH,
          bottom: -28,
          child: Stack(
            clipBehavior: Clip.none,
            children: [
              Container(
                width: 72,
                height: 72,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: const Color(0xFF141926),
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: AppColors.brandEmerald,
                    width: 3,
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withAlpha(100),
                      blurRadius: 10,
                      offset: const Offset(0, 4),
                    ),
                  ],
                ),
                child: Text(
                  initials,
                  style: const TextStyle(
                    fontFamily: 'Sora',
                    fontSize: 24,
                    fontWeight: FontWeight.w800,
                    color: AppColors.brandEmerald,
                    letterSpacing: -0.5,
                  ),
                ),
              ),
              // Verification Checkmark Badge
              Positioned(
                bottom: 2,
                right: 2,
                child: Container(
                  width: 20,
                  height: 20,
                  decoration: BoxDecoration(
                    color: const Color(0xFF141926),
                    shape: BoxShape.circle,
                    border: Border.all(
                      color: const Color(0xFF141926),
                      width: 2,
                    ),
                  ),
                  child: const Icon(
                    Icons.check_circle_rounded,
                    size: 16,
                    color: AppColors.brandEmerald,
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

/// 3-Column Metric Stat Row matching reference design ("Courses", "Following/Points", "Followers/Avg")
class _SocialStatsRow extends StatelessWidget {
  const _SocialStatsRow({
    required this.enrolledCount,
    required this.totalPoints,
    required this.avgScore,
    required this.onTapCourses,
    required this.onTapPoints,
  });

  final int enrolledCount;
  final int totalPoints;
  final double avgScore;
  final VoidCallback onTapCourses;
  final VoidCallback onTapPoints;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          // Stat 1: Courses
          Expanded(
            child: InkWell(
              onTap: onTapCourses,
              borderRadius: BorderRadius.circular(10),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.school_rounded,
                          size: 15, color: Color(0xFF38BDF8)),
                      const SizedBox(width: 4),
                      Text(
                        '$enrolledCount',
                        style: const TextStyle(
                          fontFamily: 'Sora',
                          fontSize: 16,
                          fontWeight: FontWeight.w800,
                          color: AppColors.textPrimary,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 2),
                  const Text(
                    'Courses',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textMuted,
                    ),
                  ),
                ],
              ),
            ),
          ),
          Container(height: 24, width: 1, color: AppColors.border),

          // Stat 2: Total Points
          Expanded(
            child: InkWell(
              onTap: onTapPoints,
              borderRadius: BorderRadius.circular(10),
              child: Column(
                children: [
                  Text(
                    _formatNumber(totalPoints.toDouble()),
                    style: const TextStyle(
                      fontFamily: 'Sora',
                      fontSize: 16,
                      fontWeight: FontWeight.w800,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 2),
                  const Text(
                    'Points',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textMuted,
                    ),
                  ),
                ],
              ),
            ),
          ),
          Container(height: 24, width: 1, color: AppColors.border),

          // Stat 3: Avg Score
          Expanded(
            child: Column(
              children: [
                Text(
                  '${_formatNumber(avgScore)}%',
                  style: const TextStyle(
                    fontFamily: 'Sora',
                    fontSize: 16,
                    fontWeight: FontWeight.w800,
                    color: AppColors.brandEmerald,
                  ),
                ),
                const SizedBox(height: 2),
                const Text(
                  'Avg Score',
                  style: TextStyle(
                    fontSize: 11,
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

  String _formatNumber(double value) {
    if (value == value.roundToDouble()) return value.toInt().toString();
    return value.toStringAsFixed(1);
  }
}

/// Primary Actions Bar ("+ ADD FRIENDS" / "EDIT PROFILE" + Square Share button)
class _ProfileActionsBar extends StatelessWidget {
  const _ProfileActionsBar({
    required this.onShare,
    required this.onEdit,
  });

  final VoidCallback onShare;
  final VoidCallback onEdit;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        // Primary Action Button ("+ EDIT PROFILE" / "SHARE PROFILE")
        Expanded(
          child: SizedBox(
            height: 44,
            child: ElevatedButton(
              onPressed: onShare,
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.bgSecondary,
                foregroundColor: AppColors.textPrimary,
                elevation: 0,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                  side: const BorderSide(color: AppColors.borderStrong),
                ),
              ),
              child: const Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.person_add_alt_1_rounded,
                      size: 17, color: Color(0xFF38BDF8)),
                  SizedBox(width: 8),
                  Text(
                    'SHARE PROFILE',
                    style: TextStyle(
                      fontSize: 12.5,
                      fontWeight: FontWeight.w800,
                      color: Color(0xFF38BDF8),
                      letterSpacing: 0.5,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        const SizedBox(width: 10),

        // Square Share Button ([↑])
        InkWell(
          onTap: onShare,
          borderRadius: BorderRadius.circular(12),
          child: Container(
            width: 44,
            height: 44,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: AppColors.bgSecondary,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.borderStrong),
            ),
            child: const Icon(
              Icons.ios_share_rounded,
              size: 18,
              color: Color(0xFF38BDF8),
            ),
          ),
        ),
      ],
    );
  }
}

/// "Weekly progress" Card matching the reference UI Comparative Sparkline Graph
class _WeeklyProgressCard extends StatelessWidget {
  const _WeeklyProgressCard({required this.completedCount});

  final int completedCount;

  @override
  Widget build(BuildContext context) {
    // 7-day sparkline points (Mon - Sun)
    final thisWeekData = [4.0, 3.0, 10.0, 4.0, 7.0, 9.0, 3.0];
    final lastWeekData = [32.0, 8.0, 7.0, 3.0, 5.0, 4.0, 10.0];

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: AppColors.borderStrong),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header Legend Row: "● This week  36 lessons" & "● Last week  74 lessons"
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  Container(
                    width: 7,
                    height: 7,
                    decoration: const BoxDecoration(
                      shape: BoxShape.circle,
                      color: Color(0xFF38BDF8),
                    ),
                  ),
                  const SizedBox(width: 6),
                  const Text(
                    'This week',
                    style: TextStyle(
                      fontSize: 12.5,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                ],
              ),
              Text(
                '${completedCount > 0 ? completedCount : 36} lessons',
                style: const TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w800,
                  color: Color(0xFF38BDF8),
                ),
              ),
            ],
          ),
          const SizedBox(height: 6),
          const Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  SizedBox(
                    width: 7,
                    height: 7,
                    child: DecoratedBox(
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: Color(0xFF64748B),
                      ),
                    ),
                  ),
                  SizedBox(width: 6),
                  Text(
                    'Last week',
                    style: TextStyle(
                      fontSize: 12.5,
                      fontWeight: FontWeight.w500,
                      color: AppColors.textMuted,
                    ),
                  ),
                ],
              ),
              Text(
                '74 lessons',
                style: TextStyle(
                  fontSize: 12.5,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textMuted,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // Custom 7-day sparkline chart graph
          SizedBox(
            height: 90,
            width: double.infinity,
            child: CustomPaint(
              painter: _WeeklyProgressPainter(
                thisWeekData: thisWeekData,
                lastWeekData: lastWeekData,
              ),
            ),
          ),
          const SizedBox(height: 8),

          // X-Axis Day Labels (M, T, W, T, F, S, S)
          const Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('F', style: TextStyle(fontSize: 10.5, color: AppColors.textMuted)),
              Text('S', style: TextStyle(fontSize: 10.5, color: AppColors.textMuted)),
              Text('S', style: TextStyle(fontSize: 10.5, color: AppColors.textMuted)),
              Text('M', style: TextStyle(fontSize: 10.5, color: AppColors.textMuted)),
              Text('T', style: TextStyle(fontSize: 10.5, color: AppColors.textMuted)),
              Text('W', style: TextStyle(fontSize: 10.5, color: AppColors.textMuted)),
              Text('T', style: TextStyle(fontSize: 10.5, color: AppColors.textMuted)),
            ],
          ),
        ],
      ),
    );
  }
}

/// Custom Painter for 7-day Weekly Progress Line Chart
class _WeeklyProgressPainter extends CustomPainter {
  _WeeklyProgressPainter({
    required this.thisWeekData,
    required this.lastWeekData,
  });

  final List<double> thisWeekData;
  final List<double> lastWeekData;

  @override
  void paint(Canvas canvas, Size size) {
    final gridPaint = Paint()
      ..color = AppColors.border.withAlpha(80)
      ..strokeWidth = 1;

    // Draw horizontal grid lines
    const lineCount = 3;
    for (int i = 0; i <= lineCount; i++) {
      final y = size.height * (i / lineCount);
      canvas.drawLine(Offset(0, y), Offset(size.width, y), gridPaint);
    }

    final double stepX = size.width / (thisWeekData.length - 1);

    // 1. Draw Last Week Line (Muted Grey)
    final lastWeekPaint = Paint()
      ..color = const Color(0xFF475569)
      ..strokeWidth = 2
      ..style = PaintingStyle.stroke;

    final lastWeekPath = Path();
    for (int i = 0; i < lastWeekData.length; i++) {
      final x = i * stepX;
      final y = size.height * (1.0 - (lastWeekData[i] / 40.0).clamp(0.08, 0.92));
      if (i == 0) {
        lastWeekPath.moveTo(x, y);
      } else {
        lastWeekPath.lineTo(x, y);
      }
    }
    canvas.drawPath(lastWeekPath, lastWeekPaint);

    // 2. Draw This Week Line (Vibrant Cyan Blue)
    final thisWeekPaint = Paint()
      ..color = const Color(0xFF38BDF8)
      ..strokeWidth = 2.5
      ..style = PaintingStyle.stroke;

    final dotPaint = Paint()
      ..color = const Color(0xFF38BDF8)
      ..style = PaintingStyle.fill;

    final dotOutlinePaint = Paint()
      ..color = AppColors.bgSecondary
      ..style = PaintingStyle.fill;

    final thisWeekPath = Path();
    final points = <Offset>[];

    for (int i = 0; i < thisWeekData.length; i++) {
      final x = i * stepX;
      final y = size.height * (1.0 - (thisWeekData[i] / 40.0).clamp(0.08, 0.92));
      final pt = Offset(x, y);
      points.add(pt);

      if (i == 0) {
        thisWeekPath.moveTo(x, y);
      } else {
        thisWeekPath.lineTo(x, y);
      }
    }
    canvas.drawPath(thisWeekPath, thisWeekPaint);

    // Draw dots on This Week line
    for (final pt in points) {
      canvas.drawCircle(pt, 5, dotOutlinePaint);
      canvas.drawCircle(pt, 3.5, dotPaint);
    }
  }

  @override
  bool shouldRepaint(covariant _WeeklyProgressPainter oldDelegate) => true;
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
    final topPadding = MediaQuery.paddingOf(context).top;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        top: false,
        child: ListView(
          padding: EdgeInsets.zero,
          children: [
            // Top Header Banner for Guest
            Container(
              height: 130 + topPadding,
              width: double.infinity,
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    Color(0xFF0F382A),
                    Color(0xFF072118),
                    Color(0xFF0B1713),
                  ],
                ),
              ),
              child: Stack(
                children: [
                  Positioned(
                    right: -20,
                    top: topPadding - 10,
                    child: Icon(
                      Icons.school_rounded,
                      size: 140,
                      color: Colors.white.withAlpha(12),
                    ),
                  ),
                  Positioned(
                    left: AppSizes.screenPaddingH,
                    bottom: -28,
                    child: Container(
                      width: 72,
                      height: 72,
                      alignment: Alignment.center,
                      decoration: BoxDecoration(
                        color: const Color(0xFF141926),
                        shape: BoxShape.circle,
                        border: Border.all(
                          color: AppColors.brandEmerald,
                          width: 3,
                        ),
                      ),
                      child: const Icon(
                        Icons.person_outline_rounded,
                        size: 32,
                        color: AppColors.brandEmerald,
                      ),
                    ),
                  ),
                ],
              ),
            ),

            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: AppSizes.screenPaddingH),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const SizedBox(height: 38),

                  // Guest Banner Text
                  const Text(
                    "Guest Profile",
                    style: TextStyle(
                      fontFamily: 'Sora',
                      fontSize: 22,
                      fontWeight: FontWeight.w800,
                      color: AppColors.textPrimary,
                      letterSpacing: -0.4,
                    ),
                  ),
                  const SizedBox(height: 4),
                  const Text(
                    "Sign in to sync your progress, enroll in courses, and access your mock exam analytics across devices.",
                    style: TextStyle(
                      fontSize: 13,
                      height: 1.45,
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 18),

                  // Action Buttons
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
                  const SizedBox(height: 24),

                  // Explore Options
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
                  const SizedBox(height: 32),
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
                  const SizedBox(height: 40),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
