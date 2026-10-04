import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/storage/local_storage.dart';
import '../../../../core/storage/offline_storage_cleanup_service.dart';
import '../../../../core/storage/secure_storage_service.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_datasource.dart';
import '../models/user_model.dart';

class AuthRepositoryImpl implements AuthRepository {
  const AuthRepositoryImpl(
    this._remote,
    this._secureStorage,
    this._preferences,
  );
  final AuthRemoteDataSource _remote;
  final SecureStorageService _secureStorage;
  final PreferencesService _preferences;

  @override
  Future<
      Either<Failure,
          ({UserEntity user, String accessToken, String refreshToken})>> login({
    required String email,
    required String password,
    String? deviceId,
    bool force = false,
  }) async {
    try {
      final effectiveDeviceId = deviceId ?? await _preferences.getDeviceId();
      final result = await _remote.login(
        email: email,
        password: password,
        deviceId: effectiveDeviceId,
        force: force,
      );
      await _secureStorage.saveTokens(
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      );
      await _preferences.saveCachedUser(jsonEncode(result.user.toJson()));
      return Right((
        user: result.user,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      ));
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<
          Either<Failure,
              ({UserEntity user, String accessToken, String refreshToken})>>
      register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    UserRole role = UserRole.student,
    String? phone,
  }) async {
    try {
      await _remote.register(
        email: email,
        password: password,
        firstName: firstName,
        lastName: lastName,
        phone: phone,
      );
      final deviceId = await _preferences.getDeviceId();
      final result = await _remote.login(
        email: email,
        password: password,
        deviceId: deviceId,
      );
      await _secureStorage.saveTokens(
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      );
      await _preferences.saveCachedUser(jsonEncode(result.user.toJson()));
      return Right((
        user: result.user,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      ));
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> verifyEmail(String token) async {
    try {
      await _remote.verifyEmail(token);
      // Refresh current user in cache if logged in
      final currentUserResult = await getCurrentUser();
      currentUserResult.fold((_) {}, (_) {});
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> resendVerificationEmail(String email) async {
    try {
      await _remote.resendVerificationEmail(email);
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<UserEntity?> getCachedUser() async {
    try {
      final raw = await _preferences.getCachedUser();
      if (raw == null || raw.isEmpty) return null;
      final map = jsonDecode(raw) as Map<String, dynamic>;
      return UserModel.fromJson(map);
    } catch (_) {
      return null;
    }
  }

  @override
  Future<Either<Failure, UserEntity>> getCurrentUser() async {
    try {
      final user = await _remote.getCurrentUser();
      await _preferences.saveCachedUser(jsonEncode(user.toJson()));
      return Right(user);
    } on DioException catch (e) {
      // If network fails but we have cached user, fallback to cache
      final cached = await getCachedUser();
      if (cached != null) {
        return Right(cached);
      }
      return Left(ServerFailure.fromDioError(e));
    } catch (e) {
      final cached = await getCachedUser();
      if (cached != null) {
        return Right(cached);
      }
      return Left(UnknownFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      final refreshToken = await _secureStorage.getRefreshToken();
      if (refreshToken != null) await _remote.logout(refreshToken);
    } catch (_) {
      // Ignore network errors on logout
    } finally {
      await OfflineStorageCleanupService.purgePaidDownloads();
      await _secureStorage.clearAll();
      await _preferences.clearCachedUser();
    }
    return const Right(null);
  }

  @override
  Future<Either<Failure, void>> deleteAccount() async {
    try {
      await _remote.deleteAccount();
    } catch (_) {
      // Best effort remote call — local state must still be cleaned up completely
    } finally {
      await OfflineStorageCleanupService.purgePaidDownloads();
      await _secureStorage.clearAll();
      await _preferences.clearCachedUser();
    }
    return const Right(null);
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
    return token != null && token.isNotEmpty;
  }
}

