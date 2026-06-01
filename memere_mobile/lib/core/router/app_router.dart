import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/screens/splash_screen.dart';
import '../../features/auth/presentation/screens/onboarding_screen.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/register_screen.dart';
import '../../features/auth/presentation/providers/auth_state_provider.dart';

abstract class AppRoutes {
  static const splash     = '/';
  static const onboarding = '/onboarding';
  static const login      = '/login';
  static const register   = '/register';
  static const home       = '/home';   // Phase 2
  static const courses    = '/courses'; // Phase 2
}

final appRouterProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: AppRoutes.splash,
    redirect: (context, state) {
      final isLoggedIn = authState.valueOrNull?.isAuthenticated ?? false;
      final onAuthPage = state.matchedLocation == AppRoutes.login ||
          state.matchedLocation == AppRoutes.register ||
          state.matchedLocation == AppRoutes.onboarding ||
          state.matchedLocation == AppRoutes.splash;

      // Still loading auth state → stay on splash
      if (authState.isLoading) return AppRoutes.splash;

      // Not logged in and not on auth page → go to login
      if (!isLoggedIn && !onAuthPage) return AppRoutes.login;

      // Logged in but on auth page → go to home
      if (isLoggedIn && onAuthPage && state.matchedLocation != AppRoutes.splash) {
        return AppRoutes.home;
      }

      return null; // no redirect
    },
    routes: [
      GoRoute(path: AppRoutes.splash,     builder: (_, __) => const SplashScreen()),
      GoRoute(path: AppRoutes.onboarding, builder: (_, __) => const OnboardingScreen()),
      GoRoute(path: AppRoutes.login,      builder: (_, __) => const LoginScreen()),
      GoRoute(path: AppRoutes.register,   builder: (_, __) => const RegisterScreen()),
      // home route added in Phase 2
    ],
  );
});
