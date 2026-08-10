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
import '../../domain/usecases/register_usecase.dart';
import '../providers/auth_state_provider.dart';

class RegisterScreen extends ConsumerStatefulWidget {
  const RegisterScreen({super.key});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _formKey = GlobalKey<FormState>();
  final _firstCtrl = TextEditingController();
  final _lastCtrl = TextEditingController();
  final _emailCtrl = TextEditingController();
  final _passwordCtrl = TextEditingController();
  final _phoneCtrl = TextEditingController();

  @override
  void dispose() {
    _firstCtrl.dispose();
    _lastCtrl.dispose();
    _emailCtrl.dispose();
    _passwordCtrl.dispose();
    _phoneCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!(_formKey.currentState?.validate() ?? false)) return;
    await ref.read(authStateProvider.notifier).register(
          RegisterParams(
            email: _emailCtrl.text.trim(),
            password: _passwordCtrl.text,
            firstName: _firstCtrl.text.trim(),
            lastName: _lastCtrl.text.trim(),
            phone:
                _phoneCtrl.text.trim().isEmpty ? null : _phoneCtrl.text.trim(),
          ),
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
        var message = 'Registration failed. Please try again.';
        if (error is Failure) {
          message = error.message;
        } else if (error.toString().contains('Connection refused')) {
          message = 'Cannot connect to backend server. Make sure backend is running!';
        }

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Row(
              children: [
                const Icon(Icons.wifi_off_rounded, color: Colors.white, size: 20),
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
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_new_rounded),
          onPressed: () => context.go(AppRoutes.login),
        ),
      ),
      body: AppPageBackground(
        child: SafeArea(
          top: false,
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSizes.screenPaddingH,
              vertical: AppSizes.md,
            ),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Top Mascot Illustration Header
                  Center(
                    child: Column(
                      children: [
                        const MemereMascot(
                          size: Size(190, 168),
                          showBackdrop: false,
                        ),
                        const SizedBox(height: AppSizes.xs),
                        Text(
                          'Memere',
                          style: AppTextStyles.titleLarge.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSizes.md),
                  Text('Create Student Account', style: AppTextStyles.headlineMedium.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 4),
                  Text(
                    'Start preparing for Grade 12 National Exams today',
                    style: AppTextStyles.bodyMedium.copyWith(color: AppColors.textMuted),
                  ),
                  const SizedBox(height: AppSizes.lg),
                  AppSurface(
                    padding: const EdgeInsets.all(AppSizes.lg),
                    color: AppColors.bgSecondary,
                    shadows: AppShadows.md,
                    radius: AppSizes.radiusXl,
                    child: Column(
                      children: [
                        Row(
                          children: [
                            Expanded(
                              child: AppTextField(
                                controller: _firstCtrl,
                                hintText: 'Abebe',
                                labelText: 'First Name',
                                textInputAction: TextInputAction.next,
                                validator: (v) =>
                                    (v == null || v.trim().isEmpty)
                                        ? 'Required'
                                        : null,
                              ),
                            ),
                            const SizedBox(width: AppSizes.md),
                            Expanded(
                              child: AppTextField(
                                controller: _lastCtrl,
                                hintText: 'Kebede',
                                labelText: 'Last Name',
                                textInputAction: TextInputAction.next,
                                validator: (v) =>
                                    (v == null || v.trim().isEmpty)
                                        ? 'Required'
                                        : null,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: AppSizes.md),
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
                          controller: _phoneCtrl,
                          hintText: '+251 9XX XXX XXX',
                          labelText: 'Phone (optional)',
                          prefixIcon: Icons.phone_outlined,
                          keyboardType: TextInputType.phone,
                          textInputAction: TextInputAction.next,
                        ),
                        const SizedBox(height: AppSizes.md),
                        AppTextField(
                          controller: _passwordCtrl,
                          hintText: 'Min. 8 characters',
                          labelText: 'Password',
                          prefixIcon: Icons.lock_outline_rounded,
                          isPassword: true,
                          textInputAction: TextInputAction.done,
                          autofillHints: const [AutofillHints.newPassword],
                          onFieldSubmitted: (_) => _submit(),
                          validator: (v) {
                            if (v == null || v.isEmpty) {
                              return 'Password is required';
                            }
                            if (v.length < 8) return 'Min. 8 characters';
                            return null;
                          },
                        ),
                        const SizedBox(height: AppSizes.xl),
                        AppButton(
                          label: 'Create Free Account',
                          onPressed: authAsync.isLoading ? null : _submit,
                          isLoading: authAsync.isLoading,
                          suffixIcon: Icons.arrow_forward_rounded,
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSizes.lg),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        'Already have an account? ',
                        style: AppTextStyles.bodyMedium.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                      GestureDetector(
                        onTap: () => context.go(AppRoutes.login),
                        child: Text(
                          'Sign In',
                          style: AppTextStyles.labelMedium.copyWith(
                            color: AppColors.accentPrimary,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSizes.lg),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
