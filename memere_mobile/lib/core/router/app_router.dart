import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/screens/splash_screen.dart';
import '../../features/auth/presentation/screens/onboarding_screen.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/register_screen.dart';
import '../../features/auth/presentation/providers/auth_state_provider.dart';
import '../../features/courses/presentation/screens/course_detail_screen.dart';
import '../../features/courses/presentation/screens/course_list_screen.dart';
import '../../features/exam/presentation/screens/exam_analytics_screen.dart';
import '../../features/learning/presentation/screens/my_learning_screen.dart';
import '../../features/profile/presentation/screens/profile_screen.dart';
import '../../features/exam/presentation/screens/exam_attempt_screen.dart';
import '../../features/exam/presentation/screens/exam_result_screen.dart';
import '../../features/exam/presentation/screens/mock_exam_catalog_screen.dart';
import '../../features/payment/presentation/screens/payment_result_screen.dart';
import '../../features/payment/presentation/screens/payment_webview_screen.dart';
import '../../features/payment/presentation/screens/purchase_history_screen.dart';
import '../../features/payment/presentation/screens/subscription_plans_screen.dart';
import '../../features/quiz/presentation/screens/quiz_attempt_screen.dart';
import '../../features/quiz/presentation/screens/quiz_detail_screen.dart';
import '../../features/quiz/presentation/screens/quiz_result_screen.dart';
import '../../features/video_player/presentation/screens/video_player_screen.dart';
import 'app_shell.dart';

abstract class AppRoutes {
  static const splash = '/';
  static const onboarding = '/onboarding';
  static const login = '/login';
  static const register = '/register';
  static const home = '/home';
  static const learn = '/learn';
  static const profile = '/profile';
  static const courseDetail = '/courses/:courseId';
  static const videoPlayer = '/videos/:videoId';
  static const quizDetail = '/quizzes/:quizId';
  static const quizAttempt = '/quiz-attempts/:attemptId';
  static const quizResult = '/quiz-attempts/:attemptId/result';
  static const mockExams = '/mock-exams';
  static const examAttempt = '/exam-attempts/:attemptId';
  static const examResult = '/exam-attempts/:attemptId/results';
  static const examAnalytics = '/exam-attempts/:attemptId/analytics';
  static const paymentWebView = '/payments/:paymentId/webview';
  static const paymentResult = '/payments/:paymentId/result';
  static const purchaseHistory = '/payments';
  static const subscriptionPlans = '/subscription-plans';

  static String courseDetailPath(String courseId) => '/courses/$courseId';
  static String quizDetailPath(String quizId) => '/quizzes/$quizId';
  static String quizResultPath(String attemptId) =>
      '/quiz-attempts/$attemptId/result';

  static String examAttemptPath({
    required String attemptId,
    required String examId,
  }) {
    return Uri(
      path: '/exam-attempts/$attemptId',
      queryParameters: {'examId': examId},
    ).toString();
  }

  static String examResultPath(String attemptId) =>
      '/exam-attempts/$attemptId/results';

  static String examAnalyticsPath(String attemptId) =>
      '/exam-attempts/$attemptId/analytics';

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

  static String paymentWebViewPath({
    required String paymentId,
    required String courseId,
    required String redirectUrl,
  }) {
    return Uri(
      path: '/payments/$paymentId/webview',
      queryParameters: {
        'courseId': courseId,
        'redirectUrl': redirectUrl,
      },
    ).toString();
  }

  static String paymentResultPath({
    required String paymentId,
    String courseId = '',
  }) {
    return Uri(
      path: '/payments/$paymentId/result',
      queryParameters: {'courseId': courseId},
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
      StatefulShellRoute.indexedStack(
        builder: (_, __, navigationShell) =>
            AppShell(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.home,
                builder: (_, __) => const CourseListScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.learn,
                builder: (_, __) => const MyLearningScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.mockExams,
                builder: (_, __) => const MockExamCatalogScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.profile,
                builder: (_, __) => const ProfileScreen(),
              ),
            ],
          ),
        ],
      ),
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
      GoRoute(
        path: AppRoutes.examAttempt,
        builder: (_, state) {
          final attemptId = state.pathParameters['attemptId']!;
          final examId = state.uri.queryParameters['examId'] ?? '';
          return ExamAttemptScreen(attemptId: attemptId, examId: examId);
        },
      ),
      GoRoute(
        path: AppRoutes.examResult,
        builder: (_, state) {
          final attemptId = state.pathParameters['attemptId']!;
          return ExamResultScreen(attemptId: attemptId);
        },
      ),
      GoRoute(
        path: AppRoutes.examAnalytics,
        builder: (_, state) {
          final attemptId = state.pathParameters['attemptId']!;
          return ExamAnalyticsScreen(attemptId: attemptId);
        },
      ),
      GoRoute(
        path: AppRoutes.paymentWebView,
        builder: (_, state) {
          final paymentId = state.pathParameters['paymentId']!;
          final query = state.uri.queryParameters;
          return PaymentWebViewScreen(
            paymentId: paymentId,
            redirectUrl: query['redirectUrl'] ?? '',
            courseId: query['courseId'] ?? '',
          );
        },
      ),
      GoRoute(
        path: AppRoutes.paymentResult,
        builder: (_, state) {
          final paymentId = state.pathParameters['paymentId']!;
          final courseId = state.uri.queryParameters['courseId'] ?? '';
          return PaymentResultScreen(
            paymentId: paymentId,
            courseId: courseId,
          );
        },
      ),
      GoRoute(
        path: AppRoutes.purchaseHistory,
        builder: (_, __) => const PurchaseHistoryScreen(),
      ),
      GoRoute(
        path: AppRoutes.subscriptionPlans,
        builder: (_, __) => const SubscriptionPlansScreen(),
      ),
    ],
  );
});
