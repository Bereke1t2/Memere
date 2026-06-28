import '../../../../core/network/dio_client.dart';
import '../models/user_model.dart';

abstract class AuthRemoteDataSource {
  Future<({UserModel user, String accessToken, String refreshToken})> login({
    required String email,
    required String password,
  });
  Future<UserModel> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    String? phone,
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
    required String email,
    required String password,
  }) async {
    final response =
        await _client.post<Map<String, dynamic>>('/auth/login', data: {
      'email': email,
      'password': password,
    });
    return _authPayloadFromResponse(response.data);
  }

  @override
  Future<UserModel> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    String? phone,
  }) async {
    final response =
        await _client.post<Map<String, dynamic>>('/auth/register', data: {
      'email': email,
      'password': password,
      'first_name': firstName,
      'last_name': lastName,
      'role': 'student',
      if (phone != null) 'phone': phone,
    });
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing registration response body');
    }
    return UserModel.fromJson(data);
  }

  @override
  Future<UserModel> getCurrentUser() async {
    final response = await _client.get<Map<String, dynamic>>('/auth/me');
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing user response body');
    }
    return UserModel.fromJson(data);
  }

  @override
  Future<void> logout(String refreshToken) async {
    await _client.post('/auth/logout', data: {'refresh_token': refreshToken});
  }

  @override
  Future<void> forgotPassword(String email) async {
    await _client.post('/auth/forgot-password', data: {'email': email});
  }

  ({UserModel user, String accessToken, String refreshToken})
      _authPayloadFromResponse(Map<String, dynamic>? data) {
    if (data == null) {
      throw const FormatException('Missing auth response body');
    }
    final userJson = data['user'];
    if (userJson is! Map<String, dynamic>) {
      throw const FormatException('Missing auth user payload');
    }
    return (
      user: UserModel.fromJson(userJson),
      accessToken: data['access_token'] as String,
      refreshToken: data['refresh_token'] as String,
    );
  }
}
