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
        String message = 'Server error occurred';
        String? code;
        if (data is Map) {
          if (data['message'] != null && data['message'].toString().trim().isNotEmpty) {
            message = data['message'].toString();
          } else if (data['error'] != null && data['error'].toString().trim().isNotEmpty) {
            message = data['error'].toString();
          }
          if (data['code'] != null) {
            code = data['code'].toString();
          }
        }
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

/// True when a [Failure] means the device could not reach the server (no
/// connectivity or a timeout) — the signal to fall back to on-device grading of
/// downloaded content. Both cases surface as a [ServerFailure] with a null
/// `statusCode` and one of these codes (see [ServerFailure.fromDioError]).
bool isOfflineFailure(Failure failure) =>
    failure is ServerFailure &&
    (failure.code == 'NO_INTERNET' || failure.code == 'TIMEOUT');
