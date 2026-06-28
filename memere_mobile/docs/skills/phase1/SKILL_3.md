# phase1/SKILL_3.md — Core Infrastructure, Auth Feature & All Phase 1 Screens
# ExamPrep Mobile (memere_mobile) — Phase 1, Part 3
# READ SKILL.md → SKILL_1.md → SKILL_2.md → then this file.

---

## OBJECTIVE

Build the complete auth pipeline end-to-end:
Network layer (Dio + interceptors) → Error types → Storage services →
GoRouter → Auth domain (entity + repo interface + use cases) →
Auth data (model + remote/local datasources + repo impl) →
Auth presentation (providers + screens) →
Splash screen + Onboarding screens + Login screen + Register screen

All screens must match the dark ChatGPT-inspired design system defined in SKILL_2.md.

---

## PART A — CORE INFRASTRUCTURE

---

### FILE A1 — `lib/core/errors/failures.dart`

```dart
import 'package:dio/dio.dart';

/// Base failure type
abstract class Failure {
  const Failure(this.message);
  final String message;
}

class ServerFailure extends Failure {
  const ServerFailure(super.message, {this.code, this.statusCode});
  final String? code;
  final int? statusCode;

  factory ServerFailure.fromDioError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return const ServerFailure('Connection timed out. Please check your internet.', code: 'TIMEOUT');
      case DioExceptionType.badResponse:
        final statusCode = e.response?.statusCode;
        final data = e.response?.data;
        final message = (data is Map && data['message'] != null)
            ? data['message'] as String
            : 'Server error occurred';
        final code = (data is Map && data['code'] != null) ? data['code'] as String : null;
        return ServerFailure(message, code: code, statusCode: statusCode);
      case DioExceptionType.connectionError:
        return const ServerFailure('No internet connection.', code: 'NO_INTERNET');
      default:
        return ServerFailure(e.message ?? 'An unexpected error occurred.');
    }
  }
}

class CacheFailure extends Failure {
  const CacheFailure(super.message);
}

class NetworkFailure extends Failure {
  const NetworkFailure(super.message);
}

class AuthFailure extends Failure {
  const AuthFailure(super.message, {this.code});
  final String? code;
}

class UnknownFailure extends Failure {
  const UnknownFailure(super.message);
}

class ValidationFailure extends Failure {
  const ValidationFailure(super.message, {this.field});
  final String? field;
}
```

---

### FILE A2 — `lib/core/network/dio_client.dart`

```dart
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../constants/app_constants.dart';
import '../constants/env.dart';
import 'interceptors/auth_interceptor.dart';
import 'interceptors/logging_interceptor.dart';

final dioClientProvider = Provider<DioClient>((ref) {
  return DioClient(ref);
});

class DioClient {
  DioClient(this._ref) {
    _dio = Dio(
      BaseOptions(
        baseUrl: Env.baseUrl,
        connectTimeout: AppConstants.connectionTimeout,
        receiveTimeout: AppConstants.receiveTimeout,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );
    _dio.interceptors.addAll([
      AuthInterceptor(_ref, _dio),
      LoggingInterceptor(),
    ]);
  }

  late final Dio _dio;
  final Ref _ref;

  Dio get dio => _dio;

  Future<Response<T>> get<T>(String path, {Map<String, dynamic>? queryParameters}) =>
      _dio.get<T>(path, queryParameters: queryParameters);

  Future<Response<T>> post<T>(String path, {dynamic data, Map<String, dynamic>? queryParameters}) =>
      _dio.post<T>(path, data: data, queryParameters: queryParameters);

  Future<Response<T>> put<T>(String path, {dynamic data}) =>
      _dio.put<T>(path, data: data);

  Future<Response<T>> patch<T>(String path, {dynamic data}) =>
      _dio.patch<T>(path, data: data);

  Future<Response<T>> delete<T>(String path) =>
      _dio.delete<T>(path);
}
```

---

### FILE A3 — `lib/core/network/interceptors/auth_interceptor.dart`

```dart
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../../constants/app_constants.dart';
import '../../constants/env.dart';

class AuthInterceptor extends Interceptor {
  AuthInterceptor(this._ref, this._dio);

  final Ref _ref;
  final Dio _dio;
  bool _isRefreshing = false;

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    const storage = FlutterSecureStorage();
    final token = await storage.read(key: AppConstants.accessTokenKey);
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401 && !_isRefreshing) {
      _isRefreshing = true;
      try {
        final newToken = await _refreshToken();
        if (newToken != null) {
          err.requestOptions.headers['Authorization'] = 'Bearer $newToken';
          final response = await _dio.fetch(err.requestOptions);
          handler.resolve(response);
          return;
        }
      } catch (_) {
        await _clearTokensAndRedirect();
      } finally {
        _isRefreshing = false;
      }
    }
    handler.next(err);
  }

  Future<String?> _refreshToken() async {
    const storage = FlutterSecureStorage();
    final refreshToken = await storage.read(key: AppConstants.refreshTokenKey);
    if (refreshToken == null) return null;

    final response = await Dio().post(
      '${Env.baseUrl}/auth/refresh',
      data: {'refresh_token': refreshToken},
    );
    final newAccessToken = response.data['access_token'] as String;
    await storage.write(key: AppConstants.accessTokenKey, value: newAccessToken);
    return newAccessToken;
  }

  Future<void> _clearTokensAndRedirect() async {
    const storage = FlutterSecureStorage();
    await storage.deleteAll();
    // Router redirect will handle navigation via authStateProvider
  }
}
```

---

### FILE A4 — `lib/core/network/interceptors/logging_interceptor.dart`

```dart
import 'package:dio/dio.dart';
import 'package:logger/logger.dart';

class LoggingInterceptor extends Interceptor {
  final _logger = Logger(
    printer: PrettyPrinter(
      methodCount: 0,
      errorMethodCount: 5,
      lineLength: 80,
      colors: true,
      printEmojis: true,
    ),
  );

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    _logger.d('[REQ] ${options.method} ${options.path}');
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    _logger.d('[RES] ${response.statusCode} ${response.requestOptions.path}');
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    _logger.e('[ERR] ${err.response?.statusCode} ${err.requestOptions.path}: ${err.message}');
    handler.next(err);
  }
}
```

---

### FILE A5 — `lib/core/storage/secure_storage_service.dart`

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../constants/app_constants.dart';

final secureStorageServiceProvider = Provider<SecureStorageService>((ref) {
  return SecureStorageService();
});

class SecureStorageService {
  static const _storage = FlutterSecureStorage(
    aOptions: AndroidOptions(encryptedSharedPreferences: true),
    iOptions: IOSOptions(accessibility: KeychainAccessibility.first_unlock),
  );

  Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
  }) async {
    await Future.wait([
      _storage.write(key: AppConstants.accessTokenKey, value: accessToken),
      _storage.write(key: AppConstants.refreshTokenKey, value: refreshToken),
    ]);
  }

  Future<String?> getAccessToken() =>
      _storage.read(key: AppConstants.accessTokenKey);

  Future<String?> getRefreshToken() =>
      _storage.read(key: AppConstants.refreshTokenKey);

  Future<void> clearTokens() async {
    await Future.wait([
      _storage.delete(key: AppConstants.accessTokenKey),
      _storage.delete(key: AppConstants.refreshTokenKey),
    ]);
  }

  Future<void> write(String key, String value) =>
      _storage.write(key: key, value: value);

  Future<String?> read(String key) => _storage.read(key: key);

  Future<void> delete(String key) => _storage.delete(key: key);

  Future<void> clearAll() => _storage.deleteAll();
}
```

---

### FILE A6 — `lib/core/router/app_router.dart`

```dart
import 'package:flutter/material.dart';
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
```

---

## PART B — AUTH DOMAIN LAYER

---

### FILE B1 — `lib/features/auth/domain/entities/user_entity.dart`

```dart
/// Pure Dart — zero external dependencies
class UserEntity {
  const UserEntity({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.role,
    this.phone,
    this.avatarUrl,
    this.isEmailVerified = false,
  });

  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final UserRole role;
  final String? phone;
  final String? avatarUrl;
  final bool isEmailVerified;

  String get fullName => '$firstName $lastName';
  bool get isStudent => role == UserRole.student;
  bool get isTeacher => role == UserRole.teacher;
  bool get isAdmin   => role == UserRole.admin;
}

enum UserRole { student, teacher, admin }
```

---

### FILE B2 — `lib/features/auth/domain/repositories/auth_repository.dart`

```dart
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user_entity.dart';

abstract class AuthRepository {
  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>> login({
    required String email,
    required String password,
  });

  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    UserRole role = UserRole.student,
    String? phone,
  });

  Future<Either<Failure, UserEntity>> getCurrentUser();

  Future<Either<Failure, void>> logout();

  Future<Either<Failure, void>> forgotPassword(String email);

  Future<bool> isLoggedIn();
}
```

---

### FILE B3 — `lib/features/auth/domain/usecases/login_usecase.dart`

```dart
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user_entity.dart';
import '../repositories/auth_repository.dart';

class LoginParams {
  const LoginParams({required this.email, required this.password});
  final String email;
  final String password;
}

class LoginUseCase {
  const LoginUseCase(this._repository);
  final AuthRepository _repository;

  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>>
      call(LoginParams params) {
    if (params.email.trim().isEmpty) {
      return Future.value(const Left(ValidationFailure('Email is required', field: 'email')));
    }
    if (params.password.isEmpty) {
      return Future.value(const Left(ValidationFailure('Password is required', field: 'password')));
    }
    return _repository.login(email: params.email.trim(), password: params.password);
  }
}
```

---

### FILE B4 — `lib/features/auth/domain/usecases/register_usecase.dart`

```dart
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user_entity.dart';
import '../repositories/auth_repository.dart';

class RegisterParams {
  const RegisterParams({
    required this.email,
    required this.password,
    required this.firstName,
    required this.lastName,
    this.phone,
  });
  final String email;
  final String password;
  final String firstName;
  final String lastName;
  final String? phone;
}

class RegisterUseCase {
  const RegisterUseCase(this._repository);
  final AuthRepository _repository;

  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>>
      call(RegisterParams params) {
    if (params.email.trim().isEmpty) {
      return Future.value(const Left(ValidationFailure('Email is required', field: 'email')));
    }
    if (!params.email.contains('@')) {
      return Future.value(const Left(ValidationFailure('Invalid email address', field: 'email')));
    }
    if (params.password.length < 8) {
      return Future.value(const Left(ValidationFailure('Password must be at least 8 characters', field: 'password')));
    }
    if (params.firstName.trim().isEmpty) {
      return Future.value(const Left(ValidationFailure('First name is required', field: 'firstName')));
    }
    return _repository.register(
      email: params.email.trim(),
      password: params.password,
      firstName: params.firstName.trim(),
      lastName: params.lastName.trim(),
      phone: params.phone,
    );
  }
}
```

---

## PART C — AUTH DATA LAYER

---

### FILE C1 — `lib/features/auth/data/models/user_model.dart`

```dart
import '../../domain/entities/user_entity.dart';

class UserModel extends UserEntity {
  const UserModel({
    required super.id,
    required super.email,
    required super.firstName,
    required super.lastName,
    required super.role,
    super.phone,
    super.avatarUrl,
    super.isEmailVerified,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id:              json['id'] as String,
      email:           json['email'] as String,
      firstName:       json['first_name'] as String,
      lastName:        json['last_name'] as String,
      role:            _parseRole(json['role'] as String),
      phone:           json['phone'] as String?,
      avatarUrl:       json['avatar_url'] as String?,
      isEmailVerified: json['is_email_verified'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() => {
    'id':                 id,
    'email':              email,
    'first_name':         firstName,
    'last_name':          lastName,
    'role':               role.name,
    'phone':              phone,
    'avatar_url':         avatarUrl,
    'is_email_verified':  isEmailVerified,
  };

  static UserRole _parseRole(String role) {
    switch (role) {
      case 'teacher': return UserRole.teacher;
      case 'admin':   return UserRole.admin;
      default:        return UserRole.student;
    }
  }
}
```

---

### FILE C2 — `lib/features/auth/data/datasources/auth_remote_datasource.dart`

```dart
import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../models/user_model.dart';

abstract class AuthRemoteDataSource {
  Future<({UserModel user, String accessToken, String refreshToken})> login({
    required String email, required String password,
  });
  Future<({UserModel user, String accessToken, String refreshToken})> register({
    required String email, required String password,
    required String firstName, required String lastName, String? phone,
  });
  Future<UserModel> getCurrentUser();
  Future<void> logout(String refreshToken);
  Future<void> forgotPassword(String email);
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  const AuthRemoteDataSourceImpl(this._client);
  final DioClient _client;

  @override
  Future<({UserModel user, String accessToken, String refreshToken})> login({
    required String email, required String password,
  }) async {
    final response = await _client.post('/auth/login', data: {
      'email': email, 'password': password,
    });
    final data = response.data as Map<String, dynamic>;
    return (
      user: UserModel.fromJson(data['user'] as Map<String, dynamic>),
      accessToken: data['access_token'] as String,
      refreshToken: data['refresh_token'] as String,
    );
  }

  @override
  Future<({UserModel user, String accessToken, String refreshToken})> register({
    required String email, required String password,
    required String firstName, required String lastName, String? phone,
  }) async {
    final response = await _client.post('/auth/register', data: {
      'email': email, 'password': password,
      'first_name': firstName, 'last_name': lastName,
      if (phone != null) 'phone': phone,
    });
    final data = response.data as Map<String, dynamic>;
    return (
      user: UserModel.fromJson(data['user'] as Map<String, dynamic>),
      accessToken: data['access_token'] as String,
      refreshToken: data['refresh_token'] as String,
    );
  }

  @override
  Future<UserModel> getCurrentUser() async {
    final response = await _client.get('/auth/me');
    return UserModel.fromJson(response.data as Map<String, dynamic>);
  }

  @override
  Future<void> logout(String refreshToken) async {
    await _client.post('/auth/logout', data: {'refresh_token': refreshToken});
  }

  @override
  Future<void> forgotPassword(String email) async {
    await _client.post('/auth/forgot-password', data: {'email': email});
  }
}
```

---

### FILE C3 — `lib/features/auth/data/repositories/auth_repository_impl.dart`

```dart
import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/storage/secure_storage_service.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_datasource.dart';

class AuthRepositoryImpl implements AuthRepository {
  const AuthRepositoryImpl(this._remote, this._secureStorage);
  final AuthRemoteDataSource _remote;
  final SecureStorageService _secureStorage;

  @override
  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>> login({
    required String email, required String password,
  }) async {
    try {
      final result = await _remote.login(email: email, password: password);
      await _secureStorage.saveTokens(
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      );
      return Right((user: result.user, accessToken: result.accessToken, refreshToken: result.refreshToken));
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>> register({
    required String email, required String password,
    required String firstName, required String lastName,
    UserRole role = UserRole.student, String? phone,
  }) async {
    try {
      final result = await _remote.register(
        email: email, password: password,
        firstName: firstName, lastName: lastName, phone: phone,
      );
      await _secureStorage.saveTokens(
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      );
      return Right((user: result.user, accessToken: result.accessToken, refreshToken: result.refreshToken));
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, UserEntity>> getCurrentUser() async {
    try {
      final user = await _remote.getCurrentUser();
      return Right(user);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      final refreshToken = await _secureStorage.getRefreshToken();
      if (refreshToken != null) await _remote.logout(refreshToken);
      await _secureStorage.clearAll();
      return const Right(null);
    } catch (e) {
      await _secureStorage.clearAll(); // always clear locally
      return const Right(null);
    }
  }

  @override
  Future<Either<Failure, void>> forgotPassword(String email) async {
    try {
      await _remote.forgotPassword(email);
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<bool> isLoggedIn() async {
    final token = await _secureStorage.getAccessToken();
    return token != null;
  }
}
```

---

## PART D — AUTH PRESENTATION LAYER

---

### FILE D1 — `lib/features/auth/presentation/providers/auth_state_provider.dart`

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/user_entity.dart';
import '../../data/datasources/auth_remote_datasource.dart';
import '../../data/repositories/auth_repository_impl.dart';
import '../../../../core/network/dio_client.dart';
import '../../../../core/storage/secure_storage_service.dart';
import '../../domain/usecases/login_usecase.dart';
import '../../domain/usecases/register_usecase.dart';

// ── Dependency Providers ──────────────────────────────────────────────────────

final authRemoteDataSourceProvider = Provider<AuthRemoteDataSource>((ref) {
  return AuthRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final authRepositoryProvider = Provider((ref) {
  return AuthRepositoryImpl(
    ref.watch(authRemoteDataSourceProvider),
    ref.watch(secureStorageServiceProvider),
  );
});

final loginUseCaseProvider = Provider((ref) => LoginUseCase(ref.watch(authRepositoryProvider)));
final registerUseCaseProvider = Provider((ref) => RegisterUseCase(ref.watch(authRepositoryProvider)));

// ── Auth State ────────────────────────────────────────────────────────────────

class AuthState {
  const AuthState({this.user, this.isAuthenticated = false});
  final UserEntity? user;
  final bool isAuthenticated;
}

final authStateProvider = AsyncNotifierProvider<AuthStateNotifier, AuthState>(
  AuthStateNotifier.new,
);

class AuthStateNotifier extends AsyncNotifier<AuthState> {
  @override
  Future<AuthState> build() async {
    final repo = ref.watch(authRepositoryProvider);
    final isLoggedIn = await repo.isLoggedIn();
    if (!isLoggedIn) return const AuthState();

    final result = await repo.getCurrentUser();
    return result.fold(
      (_) => const AuthState(),
      (user) => AuthState(user: user, isAuthenticated: true),
    );
  }

  Future<void> login(String email, String password) async {
    state = const AsyncLoading();
    final useCase = ref.read(loginUseCaseProvider);
    final result = await useCase(LoginParams(email: email, password: password));
    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (data) => AsyncData(AuthState(user: data.user, isAuthenticated: true)),
    );
  }

  Future<void> register(RegisterParams params) async {
    state = const AsyncLoading();
    final useCase = ref.read(registerUseCaseProvider);
    final result = await useCase(params);
    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (data) => AsyncData(AuthState(user: data.user, isAuthenticated: true)),
    );
  }

  Future<void> logout() async {
    final repo = ref.read(authRepositoryProvider);
    await repo.logout();
    state = const AsyncData(AuthState());
  }
}
```

---

### FILE D2 — `lib/features/auth/presentation/screens/splash_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../providers/auth_state_provider.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _fadeAnim;
  late Animation<double> _scaleAnim;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(vsync: this, duration: const Duration(milliseconds: 1200));
    _fadeAnim  = Tween<double>(begin: 0, end: 1).animate(
        CurvedAnimation(parent: _controller, curve: const Interval(0, 0.6, curve: Curves.easeIn)));
    _scaleAnim = Tween<double>(begin: 0.8, end: 1).animate(
        CurvedAnimation(parent: _controller, curve: const Interval(0, 0.7, curve: Curves.elasticOut)));
    _controller.forward();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    ref.listen(authStateProvider, (_, next) {
      if (next.isLoading) return;
      if (next.hasError) { context.go(AppRoutes.login); return; }
      final auth = next.value!;
      if (auth.isAuthenticated) {
        context.go(AppRoutes.home);
      } else {
        Future.delayed(const Duration(milliseconds: 2000), () {
          if (mounted) context.go(AppRoutes.onboarding);
        });
      }
    });

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: Center(
        child: FadeTransition(
          opacity: _fadeAnim,
          child: ScaleTransition(
            scale: _scaleAnim,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // App logo / icon
                Container(
                  width: 88,
                  height: 88,
                  decoration: BoxDecoration(
                    gradient: AppColors.primaryGradient,
                    borderRadius: BorderRadius.circular(24),
                    boxShadow: [
                      BoxShadow(
                        color: AppColors.accentPrimary.withOpacity(0.35),
                        blurRadius: 32,
                        spreadRadius: 4,
                      ),
                    ],
                  ),
                  child: const Icon(Icons.school_rounded, size: 44, color: Colors.white),
                ),
                const SizedBox(height: 20),
                Text('ExamPrep', style: AppTextStyles.displayMedium),
                const SizedBox(height: 6),
                Text(
                  'Grade 12 University Entrance',
                  style: AppTextStyles.bodyMedium.copyWith(color: AppColors.textSecondary),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
```

---

### FILE D3 — `lib/features/auth/presentation/screens/onboarding_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';

class _OnboardingPage {
  const _OnboardingPage({required this.title, required this.body, required this.icon, required this.color});
  final String title, body;
  final IconData icon;
  final Color color;
}

final _pages = [
  const _OnboardingPage(
    title: 'Master Every Subject',
    body: 'Video lessons from expert teachers covering all Grade 12 subjects for the national exam.',
    icon: Icons.play_circle_fill_rounded,
    color: AppColors.accentPrimary,
  ),
  const _OnboardingPage(
    title: 'Practice Like It\'s Real',
    body: 'Timed mock exams that simulate the exact national exam format with detailed analytics.',
    icon: Icons.timer_rounded,
    color: AppColors.accentSecondary,
  ),
  const _OnboardingPage(
    title: 'Study Offline, Anywhere',
    body: 'Download lessons and study without internet — critical for students across Ethiopia.',
    icon: Icons.download_done_rounded,
    color: AppColors.warning,
  ),
];

class OnboardingScreen extends StatefulWidget {
  const OnboardingScreen({super.key});

  @override
  State<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends State<OnboardingScreen> {
  final _pageController = PageController();
  int _currentPage = 0;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  void _next() {
    if (_currentPage < _pages.length - 1) {
      _pageController.nextPage(duration: const Duration(milliseconds: 350), curve: Curves.easeInOut);
    } else {
      context.go(AppRoutes.login);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
        child: Column(
          children: [
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                onPressed: () => context.go(AppRoutes.login),
                child: Text('Skip', style: AppTextStyles.labelMedium.copyWith(color: AppColors.textSecondary)),
              ),
            ),
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                itemCount: _pages.length,
                onPageChanged: (i) => setState(() => _currentPage = i),
                itemBuilder: (_, i) => _OnboardingPageWidget(page: _pages[i]),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
              child: Column(
                children: [
                  // Page indicators
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: List.generate(_pages.length, (i) => AnimatedContainer(
                      duration: const Duration(milliseconds: 250),
                      margin: const EdgeInsets.symmetric(horizontal: 4),
                      width: _currentPage == i ? 24 : 8,
                      height: 8,
                      decoration: BoxDecoration(
                        color: _currentPage == i ? AppColors.accentPrimary : AppColors.border,
                        borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                      ),
                    )),
                  ),
                  const SizedBox(height: AppSizes.xl),
                  AppButton(
                    label: _currentPage == _pages.length - 1 ? 'Get Started' : 'Next',
                    onPressed: _next,
                    suffixIcon: Icons.arrow_forward_rounded,
                  ),
                  const SizedBox(height: AppSizes.md),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text('Already have an account? ', style: AppTextStyles.bodySmall),
                      GestureDetector(
                        onTap: () => context.go(AppRoutes.login),
                        child: Text('Sign In', style: AppTextStyles.labelMedium.copyWith(color: AppColors.accentPrimary)),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSizes.lg),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _OnboardingPageWidget extends StatelessWidget {
  const _OnboardingPageWidget({required this.page});
  final _OnboardingPage page;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 120, height: 120,
            decoration: BoxDecoration(
              color: page.color.withOpacity(0.12),
              shape: BoxShape.circle,
              border: Border.all(color: page.color.withOpacity(0.3), width: 1.5),
            ),
            child: Icon(page.icon, size: 54, color: page.color),
          ),
          const SizedBox(height: AppSizes.xl),
          Text(page.title, style: AppTextStyles.headlineLarge, textAlign: TextAlign.center),
          const SizedBox(height: AppSizes.md),
          Text(
            page.body,
            style: AppTextStyles.bodyLarge.copyWith(color: AppColors.textSecondary),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}
```

---

### FILE D4 — `lib/features/auth/presentation/screens/login_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_text_field.dart';
import '../providers/auth_state_provider.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _formKey    = GlobalKey<FormState>();
  final _emailCtrl  = TextEditingController();
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
      if (next.hasError) {
        final failure = next.error as Failure?;
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(failure?.message ?? 'Login failed. Please try again.'),
          backgroundColor: AppColors.error,
          behavior: SnackBarBehavior.floating,
          margin: const EdgeInsets.all(AppSizes.md),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppSizes.radiusSm)),
        ));
      }
    });

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: SafeArea(
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
                const SizedBox(height: AppSizes.xl),

                // ── Logo ────────────────────────────────────────────────
                Container(
                  width: 56, height: 56,
                  decoration: BoxDecoration(
                    gradient: AppColors.primaryGradient,
                    borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                  ),
                  child: const Icon(Icons.school_rounded, size: 28, color: Colors.white),
                ),
                const SizedBox(height: AppSizes.xl),

                // ── Headline ─────────────────────────────────────────────
                Text('Welcome back', style: AppTextStyles.displayMedium),
                const SizedBox(height: AppSizes.sm),
                Text(
                  'Sign in to continue your exam prep',
                  style: AppTextStyles.bodyLarge.copyWith(color: AppColors.textSecondary),
                ),
                const SizedBox(height: AppSizes.xxxl),

                // ── Email ────────────────────────────────────────────────
                AppTextField(
                  controller: _emailCtrl,
                  hintText: 'your@email.com',
                  labelText: 'Email',
                  prefixIcon: Icons.email_outlined,
                  keyboardType: TextInputType.emailAddress,
                  textInputAction: TextInputAction.next,
                  autofillHints: const [AutofillHints.email],
                  validator: (v) {
                    if (v == null || v.trim().isEmpty) return 'Email is required';
                    if (!v.contains('@')) return 'Enter a valid email';
                    return null;
                  },
                ),
                const SizedBox(height: AppSizes.md),

                // ── Password ─────────────────────────────────────────────
                AppTextField(
                  controller: _passwordCtrl,
                  hintText: '••••••••',
                  labelText: 'Password',
                  prefixIcon: Icons.lock_outline_rounded,
                  isPassword: true,
                  textInputAction: TextInputAction.done,
                  autofillHints: const [AutofillHints.password],
                  onFieldSubmitted: (_) => _submit(),
                  validator: (v) {
                    if (v == null || v.isEmpty) return 'Password is required';
                    return null;
                  },
                ),

                // ── Forgot Password ──────────────────────────────────────
                Align(
                  alignment: Alignment.centerRight,
                  child: TextButton(
                    onPressed: () {/* Phase 2 */},
                    child: Text('Forgot password?',
                        style: AppTextStyles.labelMedium.copyWith(color: AppColors.accentPrimary)),
                  ),
                ),
                const SizedBox(height: AppSizes.lg),

                // ── Sign In Button ───────────────────────────────────────
                AppButton(
                  label: 'Sign In',
                  onPressed: authAsync.isLoading ? null : _submit,
                  isLoading: authAsync.isLoading,
                ),
                const SizedBox(height: AppSizes.xl),

                // ── Divider ──────────────────────────────────────────────
                Row(children: [
                  const Expanded(child: Divider(color: AppColors.border)),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: AppSizes.md),
                    child: Text('or', style: AppTextStyles.bodySmall.copyWith(color: AppColors.textDisabled)),
                  ),
                  const Expanded(child: Divider(color: AppColors.border)),
                ]),
                const SizedBox(height: AppSizes.xl),

                // ── Register CTA ─────────────────────────────────────────
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text("Don't have an account? ", style: AppTextStyles.bodyMedium.copyWith(color: AppColors.textSecondary)),
                    GestureDetector(
                      onTap: () => context.go(AppRoutes.register),
                      child: Text('Create one', style: AppTextStyles.labelMedium.copyWith(color: AppColors.accentPrimary)),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
```

---

### FILE D5 — `lib/features/auth/presentation/screens/register_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_text_field.dart';
import '../providers/auth_state_provider.dart';
import '../../domain/usecases/register_usecase.dart';

class RegisterScreen extends ConsumerStatefulWidget {
  const RegisterScreen({super.key});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _formKey     = GlobalKey<FormState>();
  final _firstCtrl   = TextEditingController();
  final _lastCtrl    = TextEditingController();
  final _emailCtrl   = TextEditingController();
  final _passwordCtrl = TextEditingController();
  final _phoneCtrl   = TextEditingController();

  @override
  void dispose() {
    _firstCtrl.dispose(); _lastCtrl.dispose();
    _emailCtrl.dispose(); _passwordCtrl.dispose(); _phoneCtrl.dispose();
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
        phone: _phoneCtrl.text.trim().isEmpty ? null : _phoneCtrl.text.trim(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final authAsync = ref.watch(authStateProvider);

    ref.listen(authStateProvider, (_, next) {
      if (next.hasError) {
        final failure = next.error as Failure?;
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(failure?.message ?? 'Registration failed. Please try again.'),
          backgroundColor: AppColors.error,
          behavior: SnackBarBehavior.floating,
          margin: const EdgeInsets.all(AppSizes.md),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppSizes.radiusSm)),
        ));
      }
    });

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      appBar: AppBar(
        backgroundColor: AppColors.bgPrimary,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_new_rounded),
          onPressed: () => context.go(AppRoutes.login),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(
            horizontal: AppSizes.screenPaddingH, vertical: AppSizes.md,
          ),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Create account', style: AppTextStyles.displayMedium),
                const SizedBox(height: AppSizes.sm),
                Text('Start your exam prep journey today',
                    style: AppTextStyles.bodyLarge.copyWith(color: AppColors.textSecondary)),
                const SizedBox(height: AppSizes.xxl),

                Row(children: [
                  Expanded(child: AppTextField(
                    controller: _firstCtrl, hintText: 'First name', labelText: 'First Name',
                    textInputAction: TextInputAction.next,
                    validator: (v) => (v == null || v.trim().isEmpty) ? 'Required' : null,
                  )),
                  const SizedBox(width: AppSizes.md),
                  Expanded(child: AppTextField(
                    controller: _lastCtrl, hintText: 'Last name', labelText: 'Last Name',
                    textInputAction: TextInputAction.next,
                    validator: (v) => (v == null || v.trim().isEmpty) ? 'Required' : null,
                  )),
                ]),
                const SizedBox(height: AppSizes.md),

                AppTextField(
                  controller: _emailCtrl, hintText: 'your@email.com', labelText: 'Email',
                  prefixIcon: Icons.email_outlined, keyboardType: TextInputType.emailAddress,
                  textInputAction: TextInputAction.next,
                  autofillHints: const [AutofillHints.email],
                  validator: (v) {
                    if (v == null || v.trim().isEmpty) return 'Email is required';
                    if (!v.contains('@')) return 'Enter a valid email';
                    return null;
                  },
                ),
                const SizedBox(height: AppSizes.md),

                AppTextField(
                  controller: _phoneCtrl, hintText: '+251 9XX XXX XXX', labelText: 'Phone (optional)',
                  prefixIcon: Icons.phone_outlined, keyboardType: TextInputType.phone,
                  textInputAction: TextInputAction.next,
                ),
                const SizedBox(height: AppSizes.md),

                AppTextField(
                  controller: _passwordCtrl, hintText: 'Min. 8 characters', labelText: 'Password',
                  prefixIcon: Icons.lock_outline_rounded, isPassword: true,
                  textInputAction: TextInputAction.done,
                  autofillHints: const [AutofillHints.newPassword],
                  onFieldSubmitted: (_) => _submit(),
                  validator: (v) {
                    if (v == null || v.isEmpty) return 'Password is required';
                    if (v.length < 8) return 'Min. 8 characters';
                    return null;
                  },
                ),
                const SizedBox(height: AppSizes.xl),

                AppButton(
                  label: 'Create Account',
                  onPressed: authAsync.isLoading ? null : _submit,
                  isLoading: authAsync.isLoading,
                ),
                const SizedBox(height: AppSizes.lg),

                Row(mainAxisAlignment: MainAxisAlignment.center, children: [
                  Text('Already have an account? ',
                      style: AppTextStyles.bodyMedium.copyWith(color: AppColors.textSecondary)),
                  GestureDetector(
                    onTap: () => context.go(AppRoutes.login),
                    child: Text('Sign In',
                        style: AppTextStyles.labelMedium.copyWith(color: AppColors.accentPrimary)),
                  ),
                ]),
                const SizedBox(height: AppSizes.lg),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
```

---

## PHASE 1 FINAL CHECKLIST

### SKILL_3 complete when:

- [ ] All Part A files (networking, storage, router) compile
- [ ] All Part B files (domain entities, interfaces, use cases) compile
- [ ] All Part C files (models, datasources, repo impl) compile
- [ ] All Part D files (providers, 4 screens) compile
- [ ] App shows: Splash → Onboarding → Login → Register flow with correct dark UI
- [ ] Login screen validates form and shows snackbar on error
- [ ] `flutter analyze` 0 errors
- [ ] Auth state persists on hot restart (token read from Secure Storage)

---

## PHASE 1 → PHASE 2 HANDOFF

**Phase 1 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 2, tell Antigravity:
```
Phase 1 is complete. Read SKILL.md and all phase2 skill files.
We are starting Phase 2: Course Browsing.
Reference: memere_mobile/docs/memere_Design_Specification.md
```

**What Phase 2 will build:**
- Course entity, model, repository
- Course listing screen (grid/list toggle, search, subject filter chips)
- Course detail screen (sections accordion, lesson list, enroll button)
- Home screen with bottom navigation shell
- Subject filter system (Math, Physics, Chemistry, etc.)
- Shimmer loading states for all course screens
- Hive caching for course metadata (1-hour TTL)
