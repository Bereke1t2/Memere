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

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailCtrl = TextEditingController();
  final _passwordCtrl = TextEditingController();

  @override
  void dispose() {
    _emailCtrl.dispose();
    _passwordCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!(_formKey.currentState?.validate() ?? false)) return;
    await ref.read(authStateProvider.notifier).login(
          _emailCtrl.text.trim(),
          _passwordCtrl.text,
        );
  }

  @override
  Widget build(BuildContext context) {
    final authAsync = ref.watch(authStateProvider);

    ref.listen(authStateProvider, (_, next) {
      final auth = next.valueOrNull;
      if (auth?.isAuthenticated ?? false) {
        context.go(AppRoutes.home);
        return;
      }

      if (next.hasError) {
        final error = next.error;
        var message = 'Invalid email or password. Please check your credentials or tap Create Account below.';
        IconData icon = Icons.lock_reset_rounded;

        if (error is ServerFailure && error.statusCode != 401) {
          message = error.message;
          icon = Icons.error_outline_rounded;
        } else if (error.toString().contains('Connection refused')) {
          message = 'Cannot connect to backend server. Make sure backend is running!';
          icon = Icons.wifi_off_rounded;
        }

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Row(
              children: [
                Icon(icon, color: Colors.white, size: 20),
                const SizedBox(width: 10),
                Expanded(child: Text(message)),
              ],
            ),
            backgroundColor: AppColors.bgQuaternary,
            behavior: SnackBarBehavior.floating,
            margin: const EdgeInsets.all(AppSizes.md),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppSizes.radiusMd),
              side: const BorderSide(color: AppColors.borderStrong),
            ),
          ),
        );
      }
    });

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: AppPageBackground(
        child: SafeArea(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSizes.screenPaddingH,
              vertical: AppSizes.screenPaddingV,
            ),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const SizedBox(height: AppSizes.sm),

                  // Mascot Header Illustration
                  Center(
                    child: Column(
                      children: [
                        const MemereMascot(
                          size: Size(210, 185),
                          showBackdrop: false,
                        ),
                        const SizedBox(height: AppSizes.sm),
                        Text(
                          'Memere',
                          style: AppTextStyles.displayMedium.copyWith(
                            fontWeight: FontWeight.bold,
                            letterSpacing: -0.5,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          'Interactive Course Learning Platform',
                          style: AppTextStyles.bodyMedium.copyWith(
                            color: AppColors.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: AppSizes.lg),
                  Text('Welcome Back', style: AppTextStyles.headlineMedium.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 4),
                  Text(
                    'Sign in to access your course materials and mock exams',
                    style: AppTextStyles.bodyMedium.copyWith(color: AppColors.textMuted),
                  ),
                  const SizedBox(height: AppSizes.md),

                  // Interactive Form Card
                  AppSurface(
                    padding: const EdgeInsets.all(AppSizes.lg),
                    color: AppColors.bgSecondary,
                    shadows: AppShadows.md,
                    radius: AppSizes.radiusXl,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        AppTextField(
                          controller: _emailCtrl,
                          hintText: 'student@memere.edu.et',
                          labelText: 'Email Address',
                          prefixIcon: Icons.email_outlined,
                          keyboardType: TextInputType.emailAddress,
                          textInputAction: TextInputAction.next,
                          autofillHints: const [AutofillHints.email],
                          validator: (v) {
                            if (v == null || v.trim().isEmpty) {
                              return 'Email is required';
                            }
                            if (!v.contains('@')) return 'Enter a valid email address';
                            return null;
                          },
                        ),
                        const SizedBox(height: AppSizes.md),
                        AppTextField(
                          controller: _passwordCtrl,
                          hintText: 'Enter your password',
                          labelText: 'Password',
                          prefixIcon: Icons.lock_outline_rounded,
                          isPassword: true,
                          textInputAction: TextInputAction.done,
                          autofillHints: const [AutofillHints.password],
                          onFieldSubmitted: (_) => _submit(),
                          validator: (v) {
                            if (v == null || v.isEmpty) {
                              return 'Password is required';
                            }
                            return null;
                          },
                        ),
                        Align(
                          alignment: Alignment.centerRight,
                          child: TextButton(
                            onPressed: () {},
                            child: Text(
                              'Forgot password?',
                              style: AppTextStyles.labelMedium.copyWith(
                                color: AppColors.accentPrimary,
                              ),
                            ),
                          ),
                        ),
                        const SizedBox(height: AppSizes.sm),
                        AppButton(
                          label: 'Sign In to Learn',
                          onPressed: authAsync.isLoading ? null : _submit,
                          isLoading: authAsync.isLoading,
                          suffixIcon: Icons.arrow_forward_rounded,
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: AppSizes.xl),

                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        "Don't have an account? ",
                        style: AppTextStyles.bodyMedium.copyWith(color: AppColors.textSecondary),
                      ),
                      GestureDetector(
                        onTap: () => context.go(AppRoutes.register),
                        child: Text(
                          'Create Account',
                          style: AppTextStyles.labelMedium.copyWith(
                            color: AppColors.accentPrimary,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSizes.md),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
