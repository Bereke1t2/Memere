import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/memere_mascot.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';

/// Clean, professional Profile & Account Screen for Memere.
///
/// Features:
/// - Hero Student Card with avatar, student credentials, and animated mascot
/// - Quick Learning Stats Strip (Enrolled Courses, Mock Exams, Saved Items)
/// - Organized Settings & Navigation Sections (Learning, Billing, Preferences, Security)
/// - Confirmation Dialog for Sign Out
class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authStateProvider).valueOrNull?.user;
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
              // 1. Top App Bar with back button (if can pop) & Screen Title
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

              // 2. Hero Student Profile Card with Avatar & Character Mascot
              _StudentHeroCard(
                name: fullName.isEmpty ? 'Student' : fullName,
                email: user?.email ?? '',
                phoneNumber: user?.phone ?? '',
                role: user?.role.name ?? 'student',
                initials: _initials(user?.firstName, user?.lastName),
              ),
              const SizedBox(height: AppSizes.lg),

              // 3. Quick Learning Stats Row
              _LearningStatsRow(
                enrolledCount: enrolledCount,
                purchasesCount: purchasesCount,
              ),
              const SizedBox(height: AppSizes.xl),

              // 4. Learning & Academic Section
              const _SectionHeader(title: 'Learning & Academics'),
              const SizedBox(height: AppSizes.xs),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.menu_book_rounded,
                    iconColor: const Color(0xFF38BDF8),
                    title: 'My Enrolled Courses',
                    subtitle: '$enrolledCount active courses in progress',
                    onTap: () => context.go(AppRoutes.learn),
                  ),
                  _SettingsItemData(
                    icon: Icons.assignment_turned_in_outlined,
                    iconColor: const Color(0xFFF59E0B),
                    title: 'Mock Exams & Results',
                    subtitle: 'Take national entrance exams & view history',
                    onTap: () => context.go(AppRoutes.mockExams),
                  ),
                  _SettingsItemData(
                    icon: Icons.bookmark_border_rounded,
                    iconColor: AppColors.brandEmerald,
                    title: 'Saved Notes & Downloads',
                    subtitle: 'Offline study guides and PDF materials',
                    onTap: () => context.go(AppRoutes.saved),
                  ),
                ],
              ),
              const SizedBox(height: AppSizes.lg),

              // 5. Billing & Memberships Section
              const _SectionHeader(title: 'Billing & Memberships'),
              const SizedBox(height: AppSizes.xs),
              _SettingsGroup(
                items: [
                  _SettingsItemData(
                    icon: Icons.workspace_premium_rounded,
                    iconColor: const Color(0xFFA855F7),
                    title: 'All-Access Plans',
                    subtitle: 'Upgrade for unlimited mock exams & courses',
                    onTap: () => context.push(AppRoutes.subscriptionPlans),
                  ),
                  _SettingsItemData(
                    icon: Icons.receipt_long_outlined,
                    iconColor: const Color(0xFF38BDF8),
                    title: 'Payment & Purchase History',
                    subtitle: 'View receipts and active transactions',
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
                    icon: Icons.translate_rounded,
                    iconColor: const Color(0xFF14B8A6),
                    title: 'Curriculum Stream',
                    subtitle: 'Natural Science (Grade 12)',
                    onTap: () {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Natural Science curriculum stream selected.'),
                        ),
                      );
                    },
                  ),
                  _SettingsItemData(
                    icon: Icons.notifications_none_rounded,
                    iconColor: const Color(0xFFF43F5E),
                    title: 'Exam Reminders & Notifications',
                    subtitle: 'Daily study goals and announcement alerts',
                    onTap: () {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Notifications are currently enabled.'),
                        ),
                      );
                    },
                  ),
                  _SettingsItemData(
                    icon: Icons.help_outline_rounded,
                    iconColor: const Color(0xFF94A3B8),
                    title: 'Help Center & Feedback',
                    subtitle: 'FAQs, contact support, report issues',
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
                    iconColor: AppColors.error,
                    title: 'Sign Out',
                    subtitle: 'Log out of your account on this device',
                    isDestructive: true,
                    onTap: () => _confirmSignOut(context, ref),
                  ),
                ],
              ),
              const SizedBox(height: AppSizes.xl),

              // 8. Footer Brand & Version
              Center(
                child: Column(
                  children: [
                    const Text(
                      'Memere • Ethiopian Entrance Exam Prep',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                        color: AppColors.textMuted,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'Version 1.0.0 (Build 42)',
                      style: TextStyle(
                        fontSize: 11,
                        color: AppColors.textMuted.withAlpha(150),
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
    return result.isEmpty ? 'M' : result;
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
        title: const Row(
          children: [
            Icon(Icons.logout_rounded, color: AppColors.error, size: 22),
            SizedBox(width: 10),
            Text('Sign Out', style: AppTextStyles.titleLarge),
          ],
        ),
        content: const Text(
          'Are you sure you want to sign out? You can sign back in anytime to access your courses and exams.',
          style: TextStyle(
            fontSize: 13.5,
            color: AppColors.textSecondary,
            height: 1.4,
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text('Cancel', style: TextStyle(color: AppColors.textMuted)),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
            ),
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: const Text('Sign Out'),
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
          'For questions about course enrollments, mock exams, or payment receipts, reach out to our team at support@memere.et.',
          style: TextStyle(
            fontSize: 13.5,
            color: AppColors.textSecondary,
            height: 1.4,
          ),
        ),
        actions: [
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.brandEmerald,
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
            ),
            onPressed: () => Navigator.of(dialogContext).pop(),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }
}

/// Top Bar for Profile Screen
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
                size: 16,
              ),
            ),
          )
        else
          const Text(
            'Account & Profile',
            style: TextStyle(
              fontSize: 22,
              fontWeight: FontWeight.w900,
              color: AppColors.textPrimary,
              letterSpacing: -0.5,
            ),
          ),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
          decoration: BoxDecoration(
            color: const Color(0x1810B981),
            borderRadius: BorderRadius.circular(20),
            border: Border.all(color: const Color(0x3510B981)),
          ),
          child: const Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.verified_rounded, size: 13, color: AppColors.brandEmerald),
              SizedBox(width: 4),
              Text(
                'Student Account',
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.bold,
                  color: AppColors.brandEmerald,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

/// Hero Student Profile Card with Character Mascot
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
        gradient: const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFF1D4ED8), // Royal Blue
            Color(0xFF1E3A8A), // Deep Blue
            Color(0xFF0F172A), // Dark Slate
          ],
        ),
        borderRadius: BorderRadius.circular(24),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF1D4ED8).withAlpha(80),
            blurRadius: 16,
            offset: const Offset(0, 6),
          ),
        ],
      ),
      child: Stack(
        children: [
          // Ambient depth circle
          Positioned(
            right: -20,
            bottom: -20,
            child: Container(
              width: 120,
              height: 120,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: Colors.white.withAlpha(12),
              ),
            ),
          ),

          // Content Row
          Row(
            children: [
              // Avatar Circle
              Container(
                width: 60,
                height: 60,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: AppColors.brandEmerald,
                  shape: BoxShape.circle,
                  border: Border.all(color: Colors.white, width: 2),
                  boxShadow: [
                    BoxShadow(
                      color: AppColors.brandEmerald.withAlpha(100),
                      blurRadius: 10,
                      offset: const Offset(0, 3),
                    ),
                  ],
                ),
                child: Text(
                  initials,
                  style: const TextStyle(
                    fontSize: 22,
                    fontWeight: FontWeight.w900,
                    color: Colors.white,
                  ),
                ),
              ),
              const SizedBox(width: 14),

              // Name, Email, & Badge
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Flexible(
                          child: Text(
                            name,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                              color: Colors.white,
                              letterSpacing: -0.3,
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 2),
                    if (email.isNotEmpty)
                      Text(
                        email,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(
                          fontSize: 12.5,
                          color: Colors.white.withAlpha(210),
                        ),
                      ),
                    if (phoneNumber.isNotEmpty) ...[
                      const SizedBox(height: 1),
                      Text(
                        phoneNumber,
                        style: TextStyle(
                          fontSize: 11.5,
                          color: Colors.white.withAlpha(170),
                        ),
                      ),
                    ],
                    const SizedBox(height: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                      decoration: BoxDecoration(
                        color: Colors.black.withAlpha(60),
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(color: Colors.white.withAlpha(35)),
                      ),
                      child: Text(
                        'GRADE 12 • ${role.toUpperCase()}',
                        style: const TextStyle(
                          fontSize: 10,
                          fontWeight: FontWeight.w800,
                          color: Color(0xFF67E8F9),
                          letterSpacing: 0.5,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 6),

              // Floating Character Mascot (Micro animated)
              const SizedBox(
                width: 64,
                height: 70,
                child: MemereMascot(
                  size: Size(64, 70),
                  showBackdrop: false,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

/// Learning Stats Row (Enrolled Courses, Mock Exams, Saved Library)
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
            icon: Icons.school_rounded,
            iconColor: const Color(0xFF38BDF8),
            label: 'Courses',
            value: '$enrolledCount',
            onTap: () => context.go(AppRoutes.learn),
          ),
        ),
        const SizedBox(width: 10),
        Expanded(
          child: _StatCard(
            icon: Icons.quiz_rounded,
            iconColor: const Color(0xFFF59E0B),
            label: 'Exams',
            value: 'Mock prep',
            onTap: () => context.go(AppRoutes.mockExams),
          ),
        ),
        const SizedBox(width: 10),
        Expanded(
          child: _StatCard(
            icon: Icons.bookmark_rounded,
            iconColor: AppColors.brandEmerald,
            label: 'Saved',
            value: 'Library',
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
    required this.iconColor,
    required this.label,
    required this.value,
    required this.onTap,
  });

  final IconData icon;
  final Color iconColor;
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
          children: [
            Icon(icon, size: 20, color: iconColor),
            const SizedBox(height: 6),
            Text(
              value,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w800,
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
          fontWeight: FontWeight.w800,
          color: Color(0xFF64748B),
          letterSpacing: 0.8,
        ),
      ),
    );
  }
}

class _SettingsItemData {
  const _SettingsItemData({
    required this.icon,
    required this.iconColor,
    required this.title,
    required this.subtitle,
    required this.onTap,
    this.isDestructive = false,
  });

  final IconData icon;
  final Color iconColor;
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
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: AppColors.borderStrong),
      ),
      clipBehavior: Clip.antiAlias,
      child: ListView.separated(
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        padding: EdgeInsets.zero,
        itemCount: items.length,
        separatorBuilder: (_, __) =>
            const Divider(height: 1, color: AppColors.borderStrong),
        itemBuilder: (context, index) {
          final item = items[index];
          return InkWell(
            onTap: item.onTap,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 13),
              child: Row(
                children: [
                  Container(
                    width: 38,
                    height: 38,
                    alignment: Alignment.center,
                    decoration: BoxDecoration(
                      color: item.iconColor.withAlpha(25),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: item.iconColor.withAlpha(50)),
                    ),
                    child: Icon(item.icon, size: 18, color: item.iconColor),
                  ),
                  const SizedBox(width: 14),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          item.title,
                          style: TextStyle(
                            fontSize: 14.5,
                            fontWeight: FontWeight.bold,
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
                    size: 20,
                    color: item.isDestructive
                        ? AppColors.error.withAlpha(160)
                        : const Color(0xFF475569),
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
