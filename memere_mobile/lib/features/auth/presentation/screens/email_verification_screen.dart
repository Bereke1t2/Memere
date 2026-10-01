import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../../shared/widgets/app_text_field.dart';
import '../../../../shared/widgets/memere_mascot.dart';
import '../providers/auth_state_provider.dart';

class EmailVerificationScreen extends ConsumerStatefulWidget {
  const EmailVerificationScreen({
    super.key,
    this.email,
    this.initialToken,
  });

  final String? email;
  final String? initialToken;

  @override
  ConsumerState<EmailVerificationScreen> createState() =>
      _EmailVerificationScreenState();
}

class _EmailVerificationScreenState
    extends ConsumerState<EmailVerificationScreen> {
  final _tokenCtrl = TextEditingController();
  final _emailCtrl = TextEditingController();
  bool _isLoading = false;
  bool _isResending = false;
  bool _isVerifiedSuccess = false;
  String? _errorMessage;
  int _resendCountdown = 0;
  Timer? _countdownTimer;

  @override
  void initState() {
    super.initState();
    if (widget.initialToken != null && widget.initialToken!.isNotEmpty) {
      _tokenCtrl.text = widget.initialToken!;
      // Auto verify if token was provided via link
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _verifyToken(widget.initialToken!);
      });
    }

    final currentUser = ref.read(authStateProvider).valueOrNull?.user;
    final initialEmail = widget.email ?? currentUser?.email ?? '';
    _emailCtrl.text = initialEmail;

    if (currentUser?.isEmailVerified == true) {
      _isVerifiedSuccess = true;
    }
  }

  @override
  void dispose() {
    _tokenCtrl.dispose();
    _emailCtrl.dispose();
    _countdownTimer?.cancel();
    super.dispose();
  }

  void _startResendCooldown() {
    setState(() => _resendCountdown = 60);
    _countdownTimer?.cancel();
    _countdownTimer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (!mounted) {
        timer.cancel();
        return;
      }
      if (_resendCountdown <= 1) {
        timer.cancel();
        setState(() => _resendCountdown = 0);
      } else {
        setState(() => _resendCountdown--);
      }
    });
  }

  Future<void> _verifyToken(String token) async {
    final trimmed = token.trim();
    if (trimmed.isEmpty) {
      setState(() => _errorMessage = 'Please enter or paste your verification token');
      return;
    }

    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    final result =
        await ref.read(authStateProvider.notifier).verifyEmail(trimmed);

    if (!mounted) return;
    setState(() => _isLoading = false);

    result.fold(
      (failure) {
        setState(() {
          if (failure is ServerFailure &&
              (failure.code == 'INVALID_VERIFICATION_TOKEN' ||
                  failure.message.toLowerCase().contains('invalid') ||
                  failure.message.toLowerCase().contains('expired'))) {
            _errorMessage =
                'Verification link or token has expired or is invalid. Please request a new verification email below.';
          } else {
            _errorMessage = failure.message;
          }
        });
      },
      (_) {
        setState(() {
          _isVerifiedSuccess = true;
          _errorMessage = null;
        });
      },
    );
  }

  Future<void> _resendEmail() async {
    final email = _emailCtrl.text.trim();
    if (email.isEmpty || !email.contains('@')) {
      setState(() => _errorMessage = 'Please enter a valid email address');
      return;
    }

    setState(() {
      _isResending = true;
      _errorMessage = null;
    });

    final result = await ref
        .read(authStateProvider.notifier)
        .resendVerificationEmail(email);

    if (!mounted) return;
    setState(() => _isResending = false);

    result.fold(
      (failure) {
        if (failure is ServerFailure &&
            (failure.code == 'EMAIL_ALREADY_VERIFIED' ||
                failure.message.toLowerCase().contains('already verified'))) {
          setState(() {
            _isVerifiedSuccess = true;
            _errorMessage = null;
          });
        } else {
          setState(() => _errorMessage = failure.message);
        }
      },
      (_) {
        _startResendCooldown();
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Row(
              children: [
                const Icon(Icons.mark_email_read_rounded,
                    color: Colors.white, size: 20),
                const SizedBox(width: 10),
                Expanded(
                  child: Text('Verification email sent to $email. Please check your inbox.'),
                ),
              ],
            ),
            backgroundColor: AppColors.brandEmerald,
            behavior: SnackBarBehavior.floating,
            margin: const EdgeInsets.all(AppSizes.md),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppSizes.radiusMd),
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(
        title: const Text('Email Verification'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_new_rounded),
          onPressed: () {
            if (context.canPop()) {
              context.pop();
            } else {
              context.go(AppRoutes.home);
            }
          },
        ),
      ),
      body: AppPageBackground(
        child: SafeArea(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSizes.screenPaddingH,
              vertical: AppSizes.screenPaddingV,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const SizedBox(height: AppSizes.md),

                // Icon / Mascot Illustration
                Center(
                  child: _isVerifiedSuccess
                      ? Container(
                          width: 88,
                          height: 88,
                          alignment: Alignment.center,
                          decoration: BoxDecoration(
                            color: AppColors.brandEmerald.withAlpha(30),
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(
                            Icons.verified_rounded,
                            color: AppColors.brandEmerald,
                            size: 48,
                          ),
                        )
                      : Container(
                          width: 88,
                          height: 88,
                          alignment: Alignment.center,
                          decoration: BoxDecoration(
                            color: AppColors.accentPrimary.withAlpha(25),
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(
                            Icons.mark_email_unread_rounded,
                            color: AppColors.accentPrimary,
                            size: 44,
                          ),
                        ),
                ),

                const SizedBox(height: AppSizes.lg),

                Text(
                  _isVerifiedSuccess
                      ? 'Email Verified!'
                      : 'Verify Your Email',
                  textAlign: TextAlign.center,
                  style: AppTextStyles.headlineLarge.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: AppSizes.xs),
                Text(
                  _isVerifiedSuccess
                      ? 'Your email is verified. You now have full access to study materials, lessons, and mock exams.'
                      : 'We sent a verification token to your email address. Please enter the token or tap the link in the email to activate your student account.',
                  textAlign: TextAlign.center,
                  style: AppTextStyles.bodyMedium.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),

                const SizedBox(height: AppSizes.xl),

                if (_errorMessage != null) ...[
                  Container(
                    padding: const EdgeInsets.all(AppSizes.md),
                    decoration: BoxDecoration(
                      color: AppColors.errorSurface,
                      borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                      border: Border.all(
                        color: AppColors.error.withAlpha(70),
                      ),
                    ),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Icon(
                          Icons.error_outline_rounded,
                          color: AppColors.error,
                          size: 20,
                        ),
                        const SizedBox(width: 10),
                        Expanded(
                          child: Text(
                            _errorMessage!,
                            style: AppTextStyles.bodySmall.copyWith(
                              color: AppColors.textPrimary,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSizes.lg),
                ],

                if (_isVerifiedSuccess) ...[
                  AppSurface(
                    padding: const EdgeInsets.all(AppSizes.lg),
                    color: AppColors.bgSecondary,
                    shadows: AppShadows.md,
                    radius: AppSizes.radiusXl,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        AppButton(
                          label: 'Start Learning',
                          onPressed: () => context.go(AppRoutes.home),
                          suffixIcon: Icons.arrow_forward_rounded,
                        ),
                      ],
                    ),
                  ),
                ] else ...[
                  AppSurface(
                    padding: const EdgeInsets.all(AppSizes.lg),
                    color: AppColors.bgSecondary,
                    shadows: AppShadows.md,
                    radius: AppSizes.radiusXl,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        AppTextField(
                          controller: _emailCtrl,
                          labelText: 'Registered Email',
                          hintText: 'student@example.com',
                          prefixIcon: Icons.email_outlined,
                          keyboardType: TextInputType.emailAddress,
                        ),
                        const SizedBox(height: AppSizes.md),
                        AppTextField(
                          controller: _tokenCtrl,
                          labelText: 'Verification Token',
                          hintText: 'Paste token from verification email',
                          prefixIcon: Icons.vpn_key_outlined,
                          textInputAction: TextInputAction.done,
                          onFieldSubmitted: (v) => _verifyToken(v),
                        ),
                        const SizedBox(height: AppSizes.lg),
                        AppButton(
                          label: 'Verify Email',
                          isLoading: _isLoading,
                          onPressed: _isLoading
                              ? null
                              : () => _verifyToken(_tokenCtrl.text),
                          suffixIcon: Icons.check_circle_outline_rounded,
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: AppSizes.lg),

                  Center(
                    child: _resendCountdown > 0
                        ? Text(
                            'Resend available in ${_resendCountdown}s',
                            style: AppTextStyles.bodySmall.copyWith(
                              color: AppColors.textMuted,
                            ),
                          )
                        : TextButton.icon(
                            icon: _isResending
                                ? const SizedBox(
                                    width: 16,
                                    height: 16,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                      color: AppColors.accentPrimary,
                                    ),
                                  )
                                : const Icon(Icons.refresh_rounded,
                                    size: 18, color: AppColors.accentPrimary),
                            label: Text(
                              _isResending
                                  ? 'Sending...'
                                  : 'Resend Verification Email',
                              style: AppTextStyles.labelMedium.copyWith(
                                color: AppColors.accentPrimary,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            onPressed: _isResending ? null : _resendEmail,
                          ),
                  ),
                ],

                const SizedBox(height: AppSizes.xl),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
