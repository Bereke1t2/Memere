import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/screens/splash_screen.dart';
import '../../features/auth/presentation/screens/onboarding_screen.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/register_screen.dart';
import '../../features/auth/presentation/providers/auth_state_provider.dart';
import '../../features/courses/presentation/screens/course_detail_screen.dart';
import '../../features/courses/presentation/screens/course_list_screen.dart';
import '../../features/quiz/presentation/screens/quiz_attempt_screen.dart';
import '../../features/quiz/presentation/screens/quiz_detail_screen.dart';
import '../../features/quiz/presentation/screens/quiz_result_screen.dart';
import '../../features/video_player/presentation/screens/video_player_screen.dart';

abstract class AppRoutes {
  static const splash = '/';
  static const onboarding = '/onboarding';
  static const login = '/login';
  static const register = '/register';
  static const home = '/home';
  static const courseDetail = '/courses/:courseId';
  static const videoPlayer = '/videos/:videoId';
  static const quizDetail = '/quizzes/:quizId';
  static const quizAttempt = '/quiz-attempts/:attemptId';
  static const quizResult = '/quiz-attempts/:attemptId/result';

  static String courseDetailPath(String courseId) => '/courses/$courseId';
  static String quizDetailPath(String quizId) => '/quizzes/$quizId';
  static String quizResultPath(String attemptId) =>
      '/quiz-attempts/$attemptId/result';

  static String videoPlayerPath({
    required String videoId,
    required String lessonId,
    required String courseId,
    String? title,
  }) {
    final query = {
      'lessonId': lessonId,
      'courseId': courseId,
      if (title != null && title.isNotEmpty) 'title': title,
    };
    return Uri(path: '/videos/$videoId', queryParameters: query).toString();
  }

  static String quizAttemptPath({
    required String attemptId,
    required String quizId,
  }) {
    return Uri(
      path: '/quiz-attempts/$attemptId',
      queryParameters: {'quizId': quizId},
    ).toString();
  }
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
      GoRoute(
        path: AppRoutes.videoPlayer,
        builder: (_, state) {
          final videoId = state.pathParameters['videoId']!;
          final query = state.uri.queryParameters;
          return VideoPlayerScreen(
            videoId: videoId,
            lessonId: query['lessonId'] ?? '',
            courseId: query['courseId'] ?? '',
            title: query['title'],
          );
        },
      ),
      GoRoute(
        path: AppRoutes.quizDetail,
        builder: (_, state) {
          final quizId = state.pathParameters['quizId']!;
          return QuizDetailScreen(quizId: quizId);
        },
      ),
      GoRoute(
        path: AppRoutes.quizAttempt,
        builder: (_, state) {
          final attemptId = state.pathParameters['attemptId']!;
          final quizId = state.uri.queryParameters['quizId'] ?? '';
          return QuizAttemptScreen(attemptId: attemptId, quizId: quizId);
        },
      ),
      GoRoute(
        path: AppRoutes.quizResult,
        builder: (_, state) {
          final attemptId = state.pathParameters['attemptId']!;
          return QuizResultScreen(attemptId: attemptId);
        },
      ),
    ],
  );
});
