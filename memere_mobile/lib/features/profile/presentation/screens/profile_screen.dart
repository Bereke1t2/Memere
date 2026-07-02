import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';

/// "Profile" tab — account summary plus access to purchases, plans and sign out.
class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authStateProvider).valueOrNull?.user;
    final fullName = [
      user?.firstName.trim() ?? '',
      user?.lastName.trim() ?? '',
    ].where((part) => part.isNotEmpty).join(' ');

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: AppPageBackground(
        child: SafeArea(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(
              AppSizes.screenPaddingH,
              AppSizes.md,
              AppSizes.screenPaddingH,
              AppSizes.xl,
            ),
            children: [
              const AppSectionHeader(
                title: 'Profile',
                subtitle: 'Account, payments, and learning access',
              ),
              const SizedBox(height: AppSizes.md),
              _ProfileHeader(
                name: fullName.isEmpty ? 'Student' : fullName,
                email: user?.email ?? '',
                initials: _initials(user?.firstName, user?.lastName),
              ),
              const SizedBox(height: AppSizes.lg),
              _ProfileRow(
                icon: Icons.receipt_long_outlined,
                label: 'Purchases',
                subtitle: 'Payment history and enrollments',
                onTap: () => context.push(AppRoutes.purchaseHistory),
              ),
              const SizedBox(height: AppSizes.sm),
              _ProfileRow(
                icon: Icons.workspace_premium_outlined,
                label: 'All-access plans',
                subtitle: 'Subscribe or manage your plan',
                onTap: () => context.push(AppRoutes.subscriptionPlans),
              ),
              const SizedBox(height: AppSizes.sm),
              _ProfileRow(
                icon: Icons.logout_rounded,
                label: 'Sign out',
                isDestructive: true,
                onTap: () => _confirmSignOut(context, ref),
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
        title: const Text('Sign out?', style: AppTextStyles.headlineSmall),
        content: Text(
          'You will need to log in again to access your courses.',
          style:
              AppTextStyles.bodyMedium.copyWith(color: AppColors.textSecondary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: const Text('Sign out'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    await ref.read(authStateProvider.notifier).logout();
    if (context.mounted) context.go(AppRoutes.login);
  }
}

class _ProfileHeader extends StatelessWidget {
  const _ProfileHeader({
    required this.name,
    required this.email,
    required this.initials,
  });

  final String name;
  final String email;
  final String initials;

  @override
  Widget build(BuildContext context) {
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.md),
      gradient: AppColors.cardGradient,
      shadows: AppShadows.md,
      child: Row(
        children: [
          Container(
            width: AppSizes.avatarLg,
            height: AppSizes.avatarLg,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: AppColors.bgTertiary,
              shape: BoxShape.circle,
              border: Border.all(color: AppColors.borderStrong),
            ),
            child: Text(
              initials,
              style: AppTextStyles.headlineSmall.copyWith(
                color: AppColors.textPrimary,
              ),
            ),
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  name,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: AppTextStyles.headlineSmall,
                ),
                if (email.isNotEmpty) ...[
                  const SizedBox(height: AppSizes.xs),
                  Text(
                    email,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodySmall,
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _ProfileRow extends StatelessWidget {
  const _ProfileRow({
    required this.icon,
    required this.label,
    required this.onTap,
    this.subtitle,
    this.isDestructive = false,
  });

  final IconData icon;
  final String label;
  final String? subtitle;
  final VoidCallback onTap;
  final bool isDestructive;

  @override
  Widget build(BuildContext context) {
    final color = isDestructive ? AppColors.error : AppColors.textPrimary;
    final iconColor = isDestructive ? AppColors.error : AppColors.accentPrimary;

    return AppSurface(
      onTap: onTap,
      padding: const EdgeInsets.all(AppSizes.md),
      shadows: AppShadows.sm,
      child: Row(
        children: [
          AppIconTile(
            icon: icon,
            color: iconColor,
            size: AppSizes.avatarMd,
            iconSize: AppSizes.iconSm,
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: AppTextStyles.titleMedium.copyWith(color: color),
                ),
                if (subtitle != null) ...[
                  const SizedBox(height: AppSizes.xs),
                  Text(subtitle!, style: AppTextStyles.bodySmall),
                ],
              ],
            ),
          ),
          if (!isDestructive)
            const Icon(
              Icons.chevron_right_rounded,
              color: AppColors.textSecondary,
            ),
        ],
      ),
    );
  }
}
