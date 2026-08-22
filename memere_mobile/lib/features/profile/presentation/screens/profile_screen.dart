import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';

/// Refined, mature, and professional Profile & Account Screen for Memere.
///
/// Design Highlights:
/// - Clean obsidian surfaces with subtle borders and shadows (no childish candy colors or mascots)
/// - Refined Hero Student Card with initials avatar, verification badge, and curriculum stream metadata
/// - Quick Academic Metric Strip with press feedback
/// - Cohesive, grouped navigation settings with unified neutral icon badges
/// - Elegant sign-out confirmation and support modals
class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authValue = ref.watch(authStateProvider).valueOrNull;
    // Guests browse without an account — show a sign-in card instead of
    // account-only sections (enrollments/purchases would 401).
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

    final topPadding = MediaQuery.paddingOf(context).top;

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
          },
          child: ListView(
            padding: EdgeInsets.fromLTRB(
              AppSizes.screenPaddingH,
              topPadding + AppSizes.sm,
              AppSizes.screenPaddingH,
              AppSizes.xxl,
            ),
            children: [
              // 1. Sleek Top Bar
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
              const SizedBox(height: AppSizes.md),

              // 2. Refined Hero Student Profile Card
              _StudentHeroCard(
                name: fullName.isEmpty ? 'Student' : fullName,
                email: user?.email ?? '',
                phoneNumber: user?.phone ?? '',
                role: user?.role.name ?? 'student',
                initials: _initials(user?.firstName, user?.lastName),
              ),
              const SizedBox(height: AppSizes.md),

              // 3. Academic Metric Strip
              _LearningStatsRow(
                enrolledCount: enrolledCount,
                purchasesCount: purchasesCount,
              ),
              const SizedBox(height: AppSizes.xl),

              // 4. Academic & Learning Section
              const _SectionHeader(title: 'Academic & Learning'),
              const SizedBox(height: AppSizes.xs),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.menu_book_outlined,
                    title: 'Enrolled Courses',
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
                    subtitle: 'Offline study guides and bookmarked materials',
                    onTap: () => context.go(AppRoutes.saved),
                  ),
                ],
              ),
              const SizedBox(height: AppSizes.lg),

              // 5. Memberships & Purchases Section
              const _SectionHeader(title: 'Membership & Billing'),
              const SizedBox(height: AppSizes.xs),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.workspace_premium_outlined,
                    iconAccent: AppColors.brandAmber,
                    title: 'All-Access Plans',
                    subtitle: 'Unlock unlimited mock exams & full solutions',
                    onTap: () => context.push(AppRoutes.subscriptionPlans),
                  ),
                  _SettingsItemData(
                    icon: Icons.receipt_outlined,
                    title: 'Purchase History & Receipts',
                    subtitle: purchasesCount == 0
                        ? 'View transaction history and invoices'
                        : '$purchasesCount recorded transactions',
                    onTap: () => context.push(AppRoutes.purchaseHistory),
                  ),
                ],
              ),
              const SizedBox(height: AppSizes.lg),

              // 6. Preferences & Support Section
              const _SectionHeader(title: 'Preferences & Support'),
              const SizedBox(height: AppSizes.xs),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.school_outlined,
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
                    title: 'Exam Reminders & Notifications',
                    subtitle: 'Daily study schedules and mock announcements',
                    onTap: () {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Study notifications are currently enabled.'),
                          behavior: SnackBarBehavior.floating,
                        ),
                      );
                    },
                  ),
                  _SettingsItemData(
                    icon: Icons.help_outline_rounded,
                    title: 'Help Center & Support',
                    subtitle: 'FAQs, contact instructors, report an issue',
                    onTap: () => _showHelpDialog(context),
                  ),
                ],
              ),
              const SizedBox(height: AppSizes.lg),

              // 7. Account Session & Sign Out
              const _SectionHeader(title: 'Account Session'),
              const SizedBox(height: AppSizes.xs),
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
              const SizedBox(height: AppSizes.xl),

              // 8. Footer Brand & Version
              const Center(
                child: Column(
                  children: [
                    Text(
                      'Memere • Ethiopian University Entrance Exam Prep',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                        color: AppColors.textMuted,
                      ),
                    ),
                    SizedBox(height: 4),
                    Text(
                      'Version 1.0.0 (Build 42)',
                      style: TextStyle(
                        fontSize: 11,
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
              style: AppTextStyles.titleLarge,
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
        actionsPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
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
              padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 10),
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
        title: const Text('Help & Support', style: AppTextStyles.titleLarge),
        content: const Text(
          'Need help with course materials, mock exams, or payment confirmations? Contact the Memere academic support team at support@memere.et.',
          style: TextStyle(
            fontSize: 13.5,
            color: AppColors.textSecondary,
            height: 1.45,
          ),
        ),
        actionsPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        actions: [
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.brandEmerald,
              foregroundColor: Colors.white,
              elevation: 0,
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
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

/// Sleek Top Bar for Profile Screen
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
        if (canPop)
          InkWell(
            onTap: onPop,
            borderRadius: BorderRadius.circular(20),
            child: Container(
              width: 38,
              height: 38,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: AppColors.bgSecondary,
                shape: BoxShape.circle,
                border: Border.all(color: AppColors.borderStrong),
              ),
              child: const Icon(
                Icons.arrow_back_ios_new_rounded,
                color: AppColors.textPrimary,
                size: 15,
              ),
            ),
          )
        else
          const Expanded(
            child: Text(
              'Profile & Account',
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                fontFamily: 'Sora',
                fontSize: 22,
                fontWeight: FontWeight.w800,
                color: AppColors.textPrimary,
                letterSpacing: -0.5,
              ),
            ),
          ),
        const SizedBox(width: 8),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
          decoration: BoxDecoration(
            color: AppColors.bgSecondary,
            borderRadius: BorderRadius.circular(20),
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
                  fontSize: 11.5,
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

/// Refined Hero Student Profile Card with Minimalist Obsidian Styling
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
      padding: const EdgeInsets.all(18),
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: AppColors.borderStrong),
        boxShadow: AppShadows.sm,
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          // Sleek Avatar with Active Ring
          Stack(
            clipBehavior: Clip.none,
            children: [
              Container(
                width: 58,
                height: 58,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: const Color(0xFF181820),
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: const Color(0xFF2C2C38),
                    width: 1.5,
                  ),
                ),
                child: Text(
                  initials,
                  style: const TextStyle(
                    fontFamily: 'Sora',
                    fontSize: 20,
                    fontWeight: FontWeight.w800,
                    color: AppColors.textPrimary,
                    letterSpacing: -0.5,
                  ),
                ),
              ),
              Positioned(
                bottom: 0,
                right: 0,
                child: Container(
                  width: 16,
                  height: 16,
                  decoration: BoxDecoration(
                    color: AppColors.bgSecondary,
                    shape: BoxShape.circle,
                    border: Border.all(
                      color: AppColors.bgSecondary,
                      width: 2,
                    ),
                  ),
                  child: const Icon(
                    Icons.check_circle_rounded,
                    size: 14,
                    color: AppColors.brandEmerald,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(width: 14),

          // Name, Email, & Metadata Badges
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  name,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    fontFamily: 'Sora',
                    fontSize: 17,
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
                      fontSize: 12.5,
                      color: AppColors.textMuted,
                    ),
                  ),
                if (phoneNumber.isNotEmpty) ...[
                  const SizedBox(height: 1),
                  Text(
                    phoneNumber,
                    style: const TextStyle(
                      fontSize: 11.5,
                      color: AppColors.textDisabled,
                    ),
                  ),
                ],
                const SizedBox(height: 8),
                Wrap(
                  spacing: 6,
                  runSpacing: 4,
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 3,
                      ),
                      decoration: BoxDecoration(
                        color: const Color(0xFF181820),
                        borderRadius: BorderRadius.circular(6),
                        border: Border.all(color: const Color(0xFF282834)),
                      ),
                      child: const Text(
                        'GRADE 12 • NATURAL SCIENCE',
                        style: TextStyle(
                          fontSize: 9.5,
                          fontWeight: FontWeight.w700,
                          color: AppColors.textSecondary,
                          letterSpacing: 0.4,
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
    );
  }
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
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 12),
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(16),
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
                fontSize: 13,
                fontWeight: FontWeight.w700,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: 2),
            Text(
              label,
              style: const TextStyle(
                fontSize: 11,
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
      padding: const EdgeInsets.only(left: 4, bottom: 6),
      child: Text(
        title.toUpperCase(),
        style: const TextStyle(
          fontSize: 11,
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
    this.iconAccent,
    this.isDestructive = false,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;
  final Color? iconAccent;
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
        borderRadius: BorderRadius.circular(18),
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
          final Color iconColor = item.isDestructive
              ? AppColors.error
              : (item.iconAccent ?? AppColors.textSecondary);
          final Color iconBg = item.isDestructive
              ? const Color(0x1DEF4444)
              : (item.iconAccent != null
                  ? item.iconAccent!.withAlpha(20)
                  : const Color(0xFF181820));
          final Color iconBorder = item.isDestructive
              ? const Color(0x35EF4444)
              : (item.iconAccent != null
                  ? item.iconAccent!.withAlpha(45)
                  : const Color(0xFF282834));

          return InkWell(
            onTap: item.onTap,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 13),
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
                  const SizedBox(width: 14),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          item.title,
                          style: TextStyle(
                            fontSize: 14,
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
                            fontSize: 11.5,
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

/// Profile tab as seen by a guest (no account). Browsing, downloads, saved
/// items, and taking quizzes/exams all work signed-out; this card explains what
/// an account adds and offers sign-in / sign-up without blocking the app.
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
          padding: EdgeInsets.fromLTRB(
            AppSizes.screenPaddingH,
            topPadding + AppSizes.sm,
            AppSizes.screenPaddingH,
            AppSizes.xxl,
          ),
          children: [
            const SizedBox(height: AppSizes.sm),
            const Text(
              'Profile & Account',
              style: TextStyle(
                fontFamily: 'Sora',
                fontSize: 22,
                fontWeight: FontWeight.w800,
                color: AppColors.textPrimary,
                letterSpacing: -0.5,
              ),
            ),
            const SizedBox(height: AppSizes.lg),

            // Guest hero card
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: AppColors.bgSecondary,
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: AppColors.borderStrong),
                boxShadow: AppShadows.sm,
              ),
              child: Column(
                children: [
                  Container(
                    width: 64,
                    height: 64,
                    alignment: Alignment.center,
                    decoration: const BoxDecoration(
                      color: Color(0xFF181820),
                      shape: BoxShape.circle,
                    ),
                    child: const Icon(
                      Icons.person_outline_rounded,
                      size: 30,
                      color: AppColors.textMuted,
                    ),
                  ),
                  const SizedBox(height: AppSizes.md),
                  const Text(
                    "You're browsing as a guest",
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontFamily: 'Sora',
                      fontSize: 17,
                      fontWeight: FontWeight.w700,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: AppSizes.xs),
                  const Text(
                    'Sign in to enroll in courses, sync your progress across '
                    'devices, and keep your purchases. Your downloads and saved '
                    'items stay on this device either way.',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 13,
                      height: 1.45,
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: AppSizes.lg),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppColors.brandEmerald,
                        foregroundColor: Colors.white,
                        elevation: 0,
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      onPressed: () => context.push(AppRoutes.login),
                      child: const Text(
                        'Sign in',
                        style:
                            TextStyle(fontWeight: FontWeight.w700, fontSize: 15),
                      ),
                    ),
                  ),
                  const SizedBox(height: AppSizes.sm),
                  SizedBox(
                    width: double.infinity,
                    child: OutlinedButton(
                      style: OutlinedButton.styleFrom(
                        foregroundColor: AppColors.textPrimary,
                        side: const BorderSide(color: AppColors.borderStrong),
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      onPressed: () => context.push(AppRoutes.register),
                      child: const Text(
                        'Create account',
                        style:
                            TextStyle(fontWeight: FontWeight.w700, fontSize: 15),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: AppSizes.xl),

            // Browse actions available without an account
            const _SectionHeader(title: 'Explore without an account'),
            const SizedBox(height: AppSizes.xs),
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
            const SizedBox(height: AppSizes.xxl),
            const Center(
              child: Text(
                'Memere • Ethiopian University Entrance Exam Prep',
                style: TextStyle(
                  fontSize: 12,
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

