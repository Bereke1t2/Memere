import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/screens/splash_screen.dart';
import '../../features/auth/presentation/screens/onboarding_screen.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/register_screen.dart';
import '../../features/auth/presentation/providers/auth_state_provider.dart';
import '../../features/courses/presentation/screens/course_detail_screen.dart';
import '../../features/courses/presentation/screens/course_list_screen.dart';

abstract class AppRoutes {
  static const splash = '/';
  static const onboarding = '/onboarding';
  static const login = '/login';
  static const register = '/register';
  static const home = '/home';
  static const courseDetail = '/courses/:courseId';

  static String courseDetailPath(String courseId) => '/courses/$courseId';
}

final appRouterProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: AppRoutes.splash,
    redirect: (context, state) {
      final location = state.matchedLocation;
      final isLoggedIn = authState.valueOrNull?.isAuthenticated ?? false;
      final hasSeenOnboarding =
          authState.valueOrNull?.hasSeenOnboarding ?? false;
      final onSplash = location == AppRoutes.splash;
      final onOnboarding = location == AppRoutes.onboarding;
      final onLogin = location == AppRoutes.login;
      final onRegister = location == AppRoutes.register;
      final onAuthPage = onLogin || onRegister || onOnboarding || onSplash;

      if (authState.isLoading) {
        return onSplash ? null : AppRoutes.splash;
      }

      if (!isLoggedIn) {
        if (!hasSeenOnboarding && !onAuthPage) return AppRoutes.onboarding;
        if (hasSeenOnboarding && (onOnboarding || !onAuthPage)) {
          return AppRoutes.login;
        }
      }

      if (isLoggedIn && onAuthPage && !onSplash) {
        return AppRoutes.home;
      }

      return null;
    },
    routes: [
      GoRoute(path: AppRoutes.splash, builder: (_, __) => const SplashScreen()),
      GoRoute(
          path: AppRoutes.onboarding,
          builder: (_, __) => const OnboardingScreen()),
      GoRoute(path: AppRoutes.login, builder: (_, __) => const LoginScreen()),
      GoRoute(
          path: AppRoutes.register, builder: (_, __) => const RegisterScreen()),
      GoRoute(
          path: AppRoutes.home, builder: (_, __) => const CourseListScreen()),
      GoRoute(
        path: AppRoutes.courseDetail,
        builder: (_, state) {
          final courseId = state.pathParameters['courseId']!;
          return CourseDetailScreen(courseId: courseId);
        },
      ),
    ],
  );
});
