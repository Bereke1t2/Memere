// ignore_for_file: unused_import
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
