import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/user_entity.dart';
import '../../data/datasources/auth_remote_datasource.dart';
import '../../data/repositories/auth_repository_impl.dart';
import '../../../../core/network/dio_client.dart';
import '../../../../core/storage/local_storage.dart';
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
    final hasSeenOnboarding = await prefs.hasSeenOnboarding();
    final isLoggedIn = await repo.isLoggedIn();
    if (!isLoggedIn) {
      return AuthState(hasSeenOnboarding: hasSeenOnboarding);
    }

    final result = await repo.getCurrentUser();
    return result.fold(
      (_) => AuthState(hasSeenOnboarding: hasSeenOnboarding),
      (user) => AuthState(
        user: user,
        isAuthenticated: true,
        hasSeenOnboarding: hasSeenOnboarding,
      ),
    );
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

  Future<void> login(String email, String password) async {
    final previous = state.valueOrNull;
    state = const AsyncLoading();
    final useCase = ref.read(loginUseCaseProvider);
    final result = await useCase(LoginParams(email: email, password: password));
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
      },
    );
  }

  Future<void> register(RegisterParams params) async {
    state = const AsyncLoading();
    final useCase = ref.read(registerUseCaseProvider);
    final result = await useCase(params);
    await result.fold(
      (failure) async {
        state = AsyncError(failure, StackTrace.current);
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
      },
    );
  }

  Future<void> logout() async {
    final repo = ref.read(authRepositoryProvider);
    await repo.logout();
    final hasSeenOnboarding =
        await ref.read(preferencesServiceProvider).hasSeenOnboarding();
    state = AsyncData(AuthState(hasSeenOnboarding: hasSeenOnboarding));
  }
}
