import 'package:dio/dio.dart';
import '../../constants/env.dart';
import '../../storage/secure_storage_service.dart';
import '../../utils/media_url_helper.dart';

String _resolveUrl(String url) {
  return fixMediaUrl(url);
}

class AuthInterceptor extends Interceptor {
  AuthInterceptor(this._dio, this._secureStorage);

  final Dio _dio;
  final SecureStorageService _secureStorage;
  bool _isRefreshing = false;

  @override
  void onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _secureStorage.getAccessToken();
    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    final path = err.requestOptions.path;
    final isAuthEndpoint = path.contains('/auth/login') ||
        path.contains('/auth/register') ||
        path.contains('/auth/refresh') ||
        path.contains('/auth/forgot-password');

    if (err.response?.statusCode == 401 && !_isRefreshing && !isAuthEndpoint) {
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
    final refreshToken = await _secureStorage.getRefreshToken();
    if (refreshToken == null || refreshToken.isEmpty) return null;

    final refreshUrl = _resolveUrl('${Env.baseUrl}/auth/refresh');
    final response = await Dio().post<Map<String, dynamic>>(
      refreshUrl,
      data: {'refresh_token': refreshToken},
    );
    final data = response.data;
    if (data == null) {
      return null;
    }
    final newAccessToken = data['access_token'] as String;
    final newRefreshToken = data['refresh_token'] as String?;
    await _secureStorage.saveTokens(
      accessToken: newAccessToken,
      refreshToken: newRefreshToken ?? refreshToken,
    );
    return newAccessToken;
  }

  Future<void> _clearTokensAndRedirect() async {
    await _secureStorage.clearTokens();
  }
}
