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
