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
