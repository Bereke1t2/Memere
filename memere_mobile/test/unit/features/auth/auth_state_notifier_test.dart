import 'package:flutter_test/flutter_test.dart';
import 'package:fpdart/fpdart.dart';
import 'package:memere_mobile/core/errors/failures.dart';
import 'package:memere_mobile/features/auth/domain/entities/user_entity.dart';
import 'package:memere_mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:memere_mobile/features/auth/domain/usecases/login_usecase.dart';

class FakeAuthRepository implements AuthRepository {
  bool loggedIn = true;
  UserEntity? cachedUser;
  UserEntity? remoteUser;
  bool shouldFailRemote = false;
  String? lastDeviceId;
  bool lastForce = false;
  int verifyEmailCalls = 0;
  int resendEmailCalls = 0;

  @override
  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>> login({
    required String email,
    required String password,
    String? deviceId,
    bool force = false,
  }) async {
    lastDeviceId = deviceId;
    lastForce = force;
    if (email == 'conflict@mirkuz.edu.et' && !force) {
      return const Left(ServerFailure('Account active on another device', code: 'ACTIVE_SESSION_EXISTS', statusCode: 409));
    }
    final user = remoteUser ??
        UserEntity(
          id: 'user-123',
          email: email,
          firstName: 'Abebe',
          lastName: 'Kebede',
          role: UserRole.student,
          isEmailVerified: false,
        );
    cachedUser = user;
    loggedIn = true;
    return Right((user: user, accessToken: 'access-123', refreshToken: 'refresh-123'));
  }

  @override
  Future<Either<Failure, ({UserEntity user, String accessToken, String refreshToken})>> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    UserRole role = UserRole.student,
    String? phone,
  }) async {
    final user = UserEntity(
      id: 'new-user',
      email: email,
      firstName: firstName,
      lastName: lastName,
      role: role,
      isEmailVerified: false,
    );
    cachedUser = user;
    loggedIn = true;
    return Right((user: user, accessToken: 'access-123', refreshToken: 'refresh-123'));
  }

  @override
  Future<Either<Failure, void>> verifyEmail(String token) async {
    verifyEmailCalls++;
    if (token == 'invalid-token') {
      return const Left(ServerFailure('Invalid or expired token', code: 'INVALID_VERIFICATION_TOKEN', statusCode: 400));
    }
    if (cachedUser != null) {
      cachedUser = UserEntity(
        id: cachedUser!.id,
        email: cachedUser!.email,
        firstName: cachedUser!.firstName,
        lastName: cachedUser!.lastName,
        role: cachedUser!.role,
        isEmailVerified: true,
      );
    }
    return const Right(null);
  }

  @override
  Future<Either<Failure, void>> resendVerificationEmail(String email) async {
    resendEmailCalls++;
    return const Right(null);
  }

  @override
  Future<UserEntity?> getCachedUser() async => cachedUser;

  @override
  Future<Either<Failure, UserEntity>> getCurrentUser() async {
    if (shouldFailRemote) {
      return const Left(ServerFailure('Network timeout', statusCode: 500));
    }
    return Right(cachedUser ??
        const UserEntity(
          id: 'user-123',
          email: 'student@mirkuz.edu.et',
          firstName: 'Abebe',
          lastName: 'Kebede',
          role: UserRole.student,
          isEmailVerified: false,
        ));
  }

  @override
  Future<Either<Failure, void>> logout() async {
    loggedIn = false;
    cachedUser = null;
    return const Right(null);
  }

  @override
  Future<Either<Failure, void>> deleteAccount() async {
    loggedIn = false;
    cachedUser = null;
    return const Right(null);
  }

  @override
  Future<Either<Failure, void>> forgotPassword(String email) async => const Right(null);

  @override
  Future<bool> isLoggedIn() async => loggedIn;
}

void main() {
  group('LoginUseCase', () {
    test('passes deviceId and force to repository', () async {
      final fakeRepo = FakeAuthRepository();
      final useCase = LoginUseCase(fakeRepo);

      final result = await useCase(
        const LoginParams(
          email: 'student@mirkuz.edu.et',
          password: 'password123',
          deviceId: 'device-uuid-xyz',
          force: true,
        ),
      );

      expect(result.isRight(), isTrue);
      expect(fakeRepo.lastDeviceId, 'device-uuid-xyz');
      expect(fakeRepo.lastForce, isTrue);
    });

    test('returns ACTIVE_SESSION_EXISTS error on conflict without force', () async {
      final fakeRepo = FakeAuthRepository();
      final useCase = LoginUseCase(fakeRepo);

      final result = await useCase(
        const LoginParams(
          email: 'conflict@mirkuz.edu.et',
          password: 'password123',
          force: false,
        ),
      );

      expect(result.isLeft(), isTrue);
      result.fold(
        (failure) {
          expect(failure, isA<ServerFailure>());
          expect((failure as ServerFailure).code, 'ACTIVE_SESSION_EXISTS');
        },
        (_) => fail('expected failure'),
      );
    });
  });

  group('Email Verification', () {
    test('verifyEmail succeeds with valid token and updates verification state', () async {
      final fakeRepo = FakeAuthRepository();
      fakeRepo.cachedUser = const UserEntity(
        id: 'u-1',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        role: UserRole.student,
        isEmailVerified: false,
      );

      final res = await fakeRepo.verifyEmail('valid-token-123');
      expect(res.isRight(), isTrue);
      expect(fakeRepo.verifyEmailCalls, 1);
      expect(fakeRepo.cachedUser?.isEmailVerified, isTrue);
    });

    test('verifyEmail returns INVALID_VERIFICATION_TOKEN for expired token', () async {
      final fakeRepo = FakeAuthRepository();
      final res = await fakeRepo.verifyEmail('invalid-token');
      expect(res.isLeft(), isTrue);
      res.fold(
        (failure) {
          expect((failure as ServerFailure).code, 'INVALID_VERIFICATION_TOKEN');
        },
        (_) => fail('expected failure'),
      );
    });

    test('resendVerificationEmail triggers repository call', () async {
      final fakeRepo = FakeAuthRepository();
      final res = await fakeRepo.resendVerificationEmail('student@example.com');
      expect(res.isRight(), isTrue);
      expect(fakeRepo.resendEmailCalls, 1);
    });
  });
}
