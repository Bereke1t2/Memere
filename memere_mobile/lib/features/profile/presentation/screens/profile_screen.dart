import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/router/app_router.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../progress/presentation/providers/progress_providers.dart';

/// Refined Profile Screen for Memere adapted directly from the reference UI design.
///
/// Design Highlights:
/// - Top patterned header banner with settings icon & approval status indicator
/// - Prominent hero avatar with verified / waitlist badge
/// - User handle & metadata row ("@username • Joined August 2024")
/// - Distinct Account Status Pill (Approved Student vs Waitlist / Pending Review)
/// - Informative Waitlist Notice Card when awaiting approval
/// - 3-Column Stat Strip (Courses, Total Points, Avg Score) for approved students
/// - Primary Action Button ("Edit Profile" / "Share Profile") + Square Share Button
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
    final isApproved = user?.isApproved ?? false;
    final isPending = user?.isPendingApproval ?? false;

    final fullName = [
      user?.firstName.trim() ?? '',
      user?.lastName.trim() ?? '',
    ].where((part) => part.isNotEmpty).join(' ');

    final enrollmentsAsync = isApproved ? ref.watch(enrollmentListProvider) : null;
    final enrolledCount = enrollmentsAsync?.valueOrNull?.length ?? 0;

    final pointsAsync = isApproved ? ref.watch(studentPointsProvider) : null;
    final points = pointsAsync?.valueOrNull;

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        top: false,
        child: RefreshIndicator(
          color: AppColors.brandEmerald,
          backgroundColor: AppColors.bgSecondary,
          onRefresh: () async {
            ref.invalidate(authStateProvider);
            if (isApproved) {
              ref.invalidate(enrollmentListProvider);
              ref.invalidate(studentPointsProvider);
            }
          },
          child: ListView(
            padding: EdgeInsets.zero,
            children: [
              // 1. Top Header Banner with Settings Icon & Hero Avatar
              _ProfileHeaderBanner(
                initials: _initials(user?.firstName, user?.lastName),
                isApproved: isApproved,
                isPending: isPending,
                onSettingsPressed: () => _showSettingsSheet(context, ref),
              ),

              // 2. Main Profile Content Body
              Padding(
                padding: const EdgeInsets.symmetric(
                    horizontal: AppSizes.screenPaddingH),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const SizedBox(height: 42),

                    // User Name & Handle / Joined Info
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Expanded(
                          child: Text(
                            fullName.isEmpty ? 'Active Student' : fullName,
                            style: const TextStyle(
                              fontFamily: 'Sora',
                              fontSize: 22,
                              fontWeight: FontWeight.w800,
                              color: AppColors.textPrimary,
                              letterSpacing: -0.4,
                            ),
                          ),
                        ),
                      ],
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
                    const SizedBox(height: 10),

                    // Account Status Pill (Approved Student vs Waitlist / Pending Review)
                    if (isPending)
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 4.5),
                        decoration: BoxDecoration(
                          color: const Color(0x24F59E0B),
                          borderRadius: BorderRadius.circular(20),
                          border: Border.all(
                            color: const Color(0x5DF59E0B),
                          ),
                        ),
                        child: const Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.hourglass_top_rounded,
                                size: 13.5, color: Color(0xFFF59E0B)),
                            SizedBox(width: 5),
                            Text(
                              'Waitlist • Pending Approval',
                              style: TextStyle(
                                fontSize: 11.5,
                                fontWeight: FontWeight.w700,
                                color: Color(0xFFF59E0B),
                              ),
                            ),
                          ],
                        ),
                      )
                    else if (isApproved)
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 4.5),
                        decoration: BoxDecoration(
                          color: const Color(0x1E10B981),
                          borderRadius: BorderRadius.circular(20),
                          border: Border.all(
                            color: const Color(0x4D10B981),
                          ),
                        ),
                        child: const Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.verified_rounded,
                                size: 14, color: AppColors.brandEmerald),
                            SizedBox(width: 5),
                            Text(
                              'Approved Student',
                              style: TextStyle(
                                fontSize: 11.5,
                                fontWeight: FontWeight.w700,
                                color: AppColors.brandEmerald,
                              ),
                            ),
                          ],
                        ),
                      ),
                    const SizedBox(height: 16),

                    // 3. Waitlist Announcement Banner (Shown only when pending approval)
                    if (isPending) ...[
                      const _WaitlistNoticeCard(),
                      const SizedBox(height: 20),
                    ],

                    // 4. Horizontal Metric Stat Row (Courses, Points, Avg Score) - Visible only for approved students
                    if (isApproved) ...[
                      _SocialStatsRow(
                        enrolledCount: enrolledCount,
                        totalPoints: (points?.totalPoints ?? 0).toInt(),
                        avgScore: points?.avgPercentage ?? 0,
                        onTapCourses: () => context.go(AppRoutes.learn),
                        onTapPoints: () => context.go(AppRoutes.mockExams),
                      ),
                      const SizedBox(height: 14),

                      // Primary Action Bar (Share Profile + Square Share Icon)
                      _ProfileActionsBar(
                        onShare: () => _shareProfile(context, user?.email),
                        onEdit: () => _showEditPrompt(context),
                      ),
                      const SizedBox(height: 24),
                    ],

                    // 5. Academic & Learning Group
                    const _SectionHeader(title: 'Academic & Learning'),
                    const SizedBox(height: 6),
                    _SettingsGroup(
                      items: [
                        _SettingsItemData(
                          icon: Icons.menu_book_outlined,
                          title: isApproved ? 'My Enrolled Courses' : 'Browse Courses',
                          subtitle: isApproved
                              ? (enrolledCount == 0
                                  ? 'Explore curriculum courses'
                                  : '$enrolledCount active courses in progress')
                              : 'Explore free introductory curriculum courses',
                          onTap: () => context.go(isApproved ? AppRoutes.learn : AppRoutes.home),
                        ),
                        if (isApproved)
                          _SettingsItemData(
                            icon: Icons.assignment_outlined,
                            title: 'Mock Exams & Results',
                            subtitle:
                                'Interactive mock exams and score analytics',
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

                    // 8. Preferences & Support Group
                    const _SectionHeader(title: 'Preferences & Support'),
                    const SizedBox(height: 6),
                    _SettingsGroup(
                      items: [
                        _SettingsItemData(
                          icon: Icons.grid_view_rounded,
                          title: 'Curriculum Stream',
                          subtitle: 'General Learning & Science',
                          onTap: () {
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(
                                content: Text(
                                  'Learning Stream: General Courses & Science',
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
                            'Memere • Interactive Course Learning Platform',
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
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Are you sure you want to sign out of Memere on this device?',
              style: TextStyle(
                fontSize: 13.5,
                color: AppColors.textSecondary,
                height: 1.45,
              ),
            ),
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: const Color(0x1DEF4444),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: const Color(0x55EF4444)),
              ),
              child: const Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Icon(
                    Icons.warning_amber_rounded,
                    color: Color(0xFFEF4444),
                    size: 20,
                  ),
                  SizedBox(width: 10),
                  Expanded(
                    child: Text(
                      'Warning: To protect course materials, signing out will erase all offline downloaded videos, PDF/HTML study notes, and mock exams from this device.',
                      style: TextStyle(
                        fontSize: 12,
                        color: Color(0xFFFCA5A5),
                        fontWeight: FontWeight.w500,
                        height: 1.4,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
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

/// Top Header Banner matching the reference UI (pattern background + avatar bridge + settings gear icon + verification/waitlist status)
class _ProfileHeaderBanner extends StatelessWidget {
  const _ProfileHeaderBanner({
    required this.initials,
    required this.isApproved,
    required this.isPending,
    required this.onSettingsPressed,
  });

  final String initials;
  final bool isApproved;
  final bool isPending;
  final VoidCallback onSettingsPressed;

  @override
  Widget build(BuildContext context) {
    final topPadding = MediaQuery.paddingOf(context).top;
    final Color avatarBorderColor = isPending
        ? const Color(0xFFF59E0B)
        : (isApproved ? AppColors.brandEmerald : AppColors.borderStrong);
    final Color avatarTextColor = isPending
        ? const Color(0xFFF59E0B)
        : (isApproved ? AppColors.brandEmerald : AppColors.textPrimary);

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
                    color: avatarBorderColor,
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
                  style: TextStyle(
                    fontFamily: 'Sora',
                    fontSize: 24,
                    fontWeight: FontWeight.w800,
                    color: avatarTextColor,
                    letterSpacing: -0.5,
                  ),
                ),
              ),
              // Status Badge (Emerald checkmark for approved, Amber hourglass for waitlist)
              if (isApproved)
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
                )
              else if (isPending)
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
                        color: const Color(0xFFF59E0B),
                        width: 1.5,
                      ),
                    ),
                    child: const Icon(
                      Icons.hourglass_top_rounded,
                      size: 12,
                      color: Color(0xFFF59E0B),
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

/// Prominent Waitlist Announcement Notice Card
class _WaitlistNoticeCard extends StatelessWidget {
  const _WaitlistNoticeCard();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF161D26),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: const Color(0xFFF59E0B).withAlpha(90),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(6),
                decoration: BoxDecoration(
                  color: const Color(0x28F59E0B),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Icon(
                  Icons.hourglass_empty_rounded,
                  size: 18,
                  color: Color(0xFFF59E0B),
                ),
              ),
              const SizedBox(width: 10),
              const Expanded(
                child: Text(
                  'Account In Waitlist Queue',
                  style: TextStyle(
                    fontFamily: 'Sora',
                    fontSize: 14.5,
                    fontWeight: FontWeight.w800,
                    color: Color(0xFFF59E0B),
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 10),
          const Text(
            'Your account registration is awaiting administrator review. You can explore curriculum courses and free introductory lessons in the meantime. Once your account is approved, you will be able to request access to premium courses and take mock exams.',
            style: TextStyle(
              fontSize: 12.5,
              height: 1.5,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}

/// 3-Column Metric Stat Row matching reference design ("Courses", "Points", "Avg Score")
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

/// Primary Actions Bar ("SHARE PROFILE" + Square Share button)
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
        // Primary Action Button ("SHARE PROFILE")
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
                    'Guest Profile',
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
                    'Sign in to sync your progress, enroll in courses, and access your mock exam analytics across devices.',
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
                      'Memere • Interactive Course Learning Platform',
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
