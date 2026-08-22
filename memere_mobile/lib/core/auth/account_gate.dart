import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/presentation/providers/auth_state_provider.dart';
import '../constants/app_colors.dart';
import '../constants/app_sizes.dart';
import '../router/app_router.dart';

/// True when no account is signed in — the app is being used as a guest.
bool isGuest(WidgetRef ref) =>
    !(ref.read(authStateProvider).valueOrNull?.isAuthenticated ?? false);

/// Gates an account-only ACTION (enroll, checkout, subscribe, server-side
/// progress writes, notifications, certificates).
///
/// Returns `true` when an account exists and the caller may proceed. For a guest
/// it presents a non-blocking sign-in prompt and returns `false`. Local-only
/// features — browsing, downloading, bookmarking, and taking guest attempts —
/// must NOT call this; they work without an account by design.
Future<bool> requireAccount(
  BuildContext context,
  WidgetRef ref, {
  String title = 'Sign in to continue',
  String message =
      'Create a free account or sign in to use this feature. Your downloads and '
      'saved items stay on this device either way.',
}) async {
  if (!isGuest(ref)) return true;
  if (!context.mounted) return false;

  final choice = await showModalBottomSheet<_GateChoice>(
    context: context,
    backgroundColor: AppColors.bgSecondary,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
    ),
    builder: (_) => _AccountGateSheet(title: title, message: message),
  );

  if (choice == null || !context.mounted) return false;
  switch (choice) {
    case _GateChoice.signIn:
      context.push(AppRoutes.login);
    case _GateChoice.register:
      context.push(AppRoutes.register);
  }
  return false;
}

enum _GateChoice { signIn, register }

class _AccountGateSheet extends StatelessWidget {
  const _AccountGateSheet({required this.title, required this.message});

  final String title;
  final String message;

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(
          AppSizes.screenPaddingH,
          AppSizes.lg,
          AppSizes.screenPaddingH,
          AppSizes.md,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Container(
              width: 40,
              height: 4,
              margin: const EdgeInsets.only(bottom: AppSizes.lg),
              decoration: BoxDecoration(
                color: AppColors.borderStrong,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            Align(
              alignment: Alignment.center,
              child: Container(
                width: 52,
                height: 52,
                alignment: Alignment.center,
                decoration: const BoxDecoration(
                  color: Color(0x1D22C55E),
                  shape: BoxShape.circle,
                ),
                child: const Icon(
                  Icons.lock_outline_rounded,
                  color: AppColors.brandEmerald,
                  size: 24,
                ),
              ),
            ),
            const SizedBox(height: AppSizes.md),
            Text(
              title,
              textAlign: TextAlign.center,
              style: const TextStyle(
                fontFamily: 'Sora',
                fontSize: 18,
                fontWeight: FontWeight.w800,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: AppSizes.xs),
            Text(
              message,
              textAlign: TextAlign.center,
              style: const TextStyle(
                fontSize: 13,
                height: 1.45,
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: AppSizes.lg),
            ElevatedButton(
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.brandEmerald,
                foregroundColor: Colors.white,
                elevation: 0,
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              onPressed: () => Navigator.of(context).pop(_GateChoice.signIn),
              child: const Text(
                'Sign in',
                style: TextStyle(fontWeight: FontWeight.w700, fontSize: 15),
              ),
            ),
            const SizedBox(height: AppSizes.sm),
            OutlinedButton(
              style: OutlinedButton.styleFrom(
                foregroundColor: AppColors.textPrimary,
                side: const BorderSide(color: AppColors.borderStrong),
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              onPressed: () => Navigator.of(context).pop(_GateChoice.register),
              child: const Text(
                'Create account',
                style: TextStyle(fontWeight: FontWeight.w700, fontSize: 15),
              ),
            ),
            const SizedBox(height: AppSizes.xs),
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text(
                'Not now',
                style: TextStyle(
                  color: AppColors.textMuted,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
