// ignore_for_file: unused_import
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../../constants/app_constants.dart';
import '../../constants/env.dart';

class AuthInterceptor extends Interceptor {
  AuthInterceptor(this._dio);

  final Dio _dio;
  bool _isRefreshing = false;

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    const storage = FlutterSecureStorage();
    final token = await storage.read(key: AppConstants.accessTokenKey);
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401 && !_isRefreshing) {
      _isRefreshing = true;
      try {
        final newToken = await _refreshToken();
        if (newToken != null) {
          err.requestOptions.headers['Authorization'] = 'Bearer $newToken';
          final response = await _dio.fetch(err.requestOptions);
          handler.resolve(response);
          return;
        }
      } catch (_) {
        await _clearTokensAndRedirect();
      } finally {
        _isRefreshing = false;
      }
    }
    handler.next(err);
  }

  Future<String?> _refreshToken() async {
    const storage = FlutterSecureStorage();
    final refreshToken = await storage.read(key: AppConstants.refreshTokenKey);
    if (refreshToken == null) return null;

    final response = await Dio().post(
      '${Env.baseUrl}/auth/refresh',
      data: {'refresh_token': refreshToken},
    );
    final data = response.data as Map<String, dynamic>;
    final newAccessToken = data['access_token'] as String;
    await storage.write(key: AppConstants.accessTokenKey, value: newAccessToken);
    return newAccessToken;
  }

  Future<void> _clearTokensAndRedirect() async {
    const storage = FlutterSecureStorage();
    await storage.deleteAll();
    // Router redirect will handle navigation via authStateProvider
  }
}
