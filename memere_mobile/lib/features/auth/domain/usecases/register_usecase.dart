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

  Future<
      Either<Failure,
          ({UserEntity user, String accessToken, String refreshToken})>> call(
      RegisterParams params) {
    if (params.email.trim().isEmpty) {
      return Future.value(
          const Left(ValidationFailure('Email is required', field: 'email')));
    }
    if (!params.email.contains('@')) {
      return Future.value(const Left(
          ValidationFailure('Invalid email address', field: 'email')));
    }
    if (params.password.length < 8) {
      return Future.value(const Left(ValidationFailure(
          'Password must be at least 8 characters',
          field: 'password')));
    }
    if (params.firstName.trim().isEmpty) {
      return Future.value(const Left(
          ValidationFailure('First name is required', field: 'firstName')));
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
