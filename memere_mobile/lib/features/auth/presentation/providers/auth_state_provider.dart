import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import '../../domain/entities/user_entity.dart';
import '../../data/datasources/auth_remote_datasource.dart';
import '../../data/repositories/auth_repository_impl.dart';
import '../../domain/repositories/auth_repository.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../../../../core/storage/local_storage.dart';
import '../../../../core/storage/secure_storage_service.dart';
import '../../domain/usecases/login_usecase.dart';
import '../../domain/usecases/register_usecase.dart';
import '../../../payment/presentation/providers/purchase_history_provider.dart';
import '../../../payment/presentation/providers/subscription_provider.dart';

// ── Dependency Providers ──────────────────────────────────────────────────────

final authRemoteDataSourceProvider = Provider<AuthRemoteDataSource>((ref) {
  return AuthRemoteDataSourceImpl(ref.watch(dioClientProvider));
});

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepositoryImpl(
    ref.watch(authRemoteDataSourceProvider),
    ref.watch(secureStorageServiceProvider),
    ref.watch(preferencesServiceProvider),
  );
});

final loginUseCaseProvider =
    Provider((ref) => LoginUseCase(ref.watch(authRepositoryProvider)));
final registerUseCaseProvider =
    Provider((ref) => RegisterUseCase(ref.watch(authRepositoryProvider)));

// ── Auth State ────────────────────────────────────────────────────────────────

class AuthState {
  const AuthState({
    this.user,
    this.isAuthenticated = false,
    this.hasSeenOnboarding = false,
  });

  final UserEntity? user;
  final bool isAuthenticated;
  final bool hasSeenOnboarding;
}

final authStateProvider = AsyncNotifierProvider<AuthStateNotifier, AuthState>(
  AuthStateNotifier.new,
);

class AuthStateNotifier extends AsyncNotifier<AuthState> {
  @override
  Future<AuthState> build() async {
    final repo = ref.watch(authRepositoryProvider);
    final prefs = ref.watch(preferencesServiceProvider);
    final hasSeenOnboarding =
        await prefs.hasSeenOnboarding().catchError((_) => false);
    final isLoggedIn = await repo.isLoggedIn().catchError((_) => false);

    if (!isLoggedIn) {
      return AuthState(hasSeenOnboarding: hasSeenOnboarding);
    }

    // 1. Immediately hydrate from cached user for instantaneous startup
    final cachedUser = await repo.getCachedUser().catchError((_) => null);

    // 2. Fetch fresh user profile in background without blocking startup
    _syncUserProfileInBackground(repo);

    if (cachedUser != null) {
      return AuthState(
        user: cachedUser,
        isAuthenticated: true,
        hasSeenOnboarding: hasSeenOnboarding,
      );
    }

    // If no local cache exists, fetch from backend
    try {
      final result = await repo.getCurrentUser();
      return result.fold(
        (_) => AuthState(
          isAuthenticated: true,
          hasSeenOnboarding: hasSeenOnboarding,
        ),
        (user) => AuthState(
          user: user,
          isAuthenticated: true,
          hasSeenOnboarding: hasSeenOnboarding,
        ),
      );
    } catch (_) {
      return AuthState(
        isAuthenticated: true,
        hasSeenOnboarding: hasSeenOnboarding,
      );
    }
  }

  void _syncUserProfileInBackground(AuthRepository repo) {
    Future.microtask(() async {
      try {
        final result = await repo.getCurrentUser();
        result.fold(
          (_) {},
          (freshUser) {
            final curr = state.valueOrNull;
            if (curr != null && curr.isAuthenticated) {
              state = AsyncData(
                AuthState(
                  user: freshUser,
                  isAuthenticated: true,
                  hasSeenOnboarding: curr.hasSeenOnboarding,
                ),
              );
            }
          },
        );
      } catch (_) {}
    });
  }

  Future<void> markOnboardingSeen() async {
    await ref.read(preferencesServiceProvider).markOnboardingSeen();
    final previous = state.valueOrNull;
    state = AsyncData(
      AuthState(
        user: previous?.user,
        isAuthenticated: previous?.isAuthenticated ?? false,
        hasSeenOnboarding: true,
      ),
    );
  }

  Future<void> login(
    String email,
    String password, {
    String? deviceId,
    bool force = false,
  }) async {
    final previous = state.valueOrNull;
    state = const AsyncLoading();
    final useCase = ref.read(loginUseCaseProvider);
    final result = await useCase(
      LoginParams(
        email: email,
        password: password,
        deviceId: deviceId,
        force: force,
      ),
    );
    await result.fold(
      (failure) async {
        final hasSeenOnboarding = previous?.hasSeenOnboarding ??
            await ref.read(preferencesServiceProvider).hasSeenOnboarding();
        state = AsyncError<AuthState>(failure, StackTrace.current).copyWithPrevious(
          AsyncData(AuthState(hasSeenOnboarding: hasSeenOnboarding)),
        );
      },
      (data) async {
        await ref.read(preferencesServiceProvider).markOnboardingSeen();
        state = AsyncData(
          AuthState(
            user: data.user,
            isAuthenticated: true,
            hasSeenOnboarding: true,
          ),
        );
        _invalidateUserDataProviders();
      },
    );
  }

  Future<void> register(RegisterParams params) async {
    final previous = state.valueOrNull;
    state = const AsyncLoading();
    final useCase = ref.read(registerUseCaseProvider);
    final result = await useCase(params);
    await result.fold(
      (failure) async {
        final hasSeenOnboarding = previous?.hasSeenOnboarding ??
            await ref.read(preferencesServiceProvider).hasSeenOnboarding();
        state = AsyncError<AuthState>(failure, StackTrace.current).copyWithPrevious(
          AsyncData(AuthState(hasSeenOnboarding: hasSeenOnboarding)),
        );
      },
      (data) async {
        await ref.read(preferencesServiceProvider).markOnboardingSeen();
        state = AsyncData(
          AuthState(
            user: data.user,
            isAuthenticated: true,
            hasSeenOnboarding: true,
          ),
        );
        _invalidateUserDataProviders();
      },
    );
  }

  Future<Either<Failure, void>> verifyEmail(String token) async {
    final repo = ref.read(authRepositoryProvider);
    final result = await repo.verifyEmail(token);
    return result.fold(
      (failure) => Left(failure),
      (_) async {
        final userResult = await repo.getCurrentUser();
        userResult.fold((_) {}, (freshUser) {
          final curr = state.valueOrNull;
          state = AsyncData(
            AuthState(
              user: freshUser,
              isAuthenticated: true,
              hasSeenOnboarding: curr?.hasSeenOnboarding ?? true,
            ),
          );
        });
        return const Right(null);
      },
    );
  }

  Future<Either<Failure, void>> resendVerificationEmail(String email) async {
    final repo = ref.read(authRepositoryProvider);
    return repo.resendVerificationEmail(email);
  }

  Future<void> refreshUser() async {
    final repo = ref.read(authRepositoryProvider);
    final result = await repo.getCurrentUser();
    result.fold((_) {}, (user) {
      final curr = state.valueOrNull;
      state = AsyncData(
        AuthState(
          user: user,
          isAuthenticated: curr?.isAuthenticated ?? true,
          hasSeenOnboarding: curr?.hasSeenOnboarding ?? true,
        ),
      );
    });
  }

  Future<void> logout() async {
    final repo = ref.read(authRepositoryProvider);
    await repo.logout();
    final hasSeenOnboarding =
        await ref.read(preferencesServiceProvider).hasSeenOnboarding();
    state = AsyncData(AuthState(hasSeenOnboarding: hasSeenOnboarding));
    _invalidateUserDataProviders();
  }

  Future<void> deleteAccount() async {
    final repo = ref.read(authRepositoryProvider);
    await repo.deleteAccount();
    final hasSeenOnboarding =
        await ref.read(preferencesServiceProvider).hasSeenOnboarding();
    state = AsyncData(AuthState(hasSeenOnboarding: hasSeenOnboarding));
    _invalidateUserDataProviders();
  }

  void _invalidateUserDataProviders() {
    ref.invalidate(enrollmentListProvider);
    ref.invalidate(paymentHistoryProvider);
    ref.invalidate(mySubscriptionProvider);
  }
}

