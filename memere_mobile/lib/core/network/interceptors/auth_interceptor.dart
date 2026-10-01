import 'dart:async';
import 'package:dio/dio.dart';
import '../../constants/env.dart';
import '../../storage/offline_storage_cleanup_service.dart';
import '../../storage/secure_storage_service.dart';
import '../../utils/media_url_helper.dart';

String _resolveUrl(String url) {
  return fixMediaUrl(url);
}

class AuthInterceptor extends Interceptor {
  AuthInterceptor(this._dio, this._secureStorage);

  final Dio _dio;
  final SecureStorageService _secureStorage;
  Completer<String?>? _refreshCompleter;

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

    if (err.response?.statusCode == 401 && !isAuthEndpoint) {
      if (_refreshCompleter != null) {
        try {
          final newToken = await _refreshCompleter!.future;
          if (newToken != null && newToken.isNotEmpty) {
            err.requestOptions.headers['Authorization'] = 'Bearer $newToken';
            final response = await _dio.fetch(err.requestOptions);
            handler.resolve(response);
            return;
          }
        } catch (_) {
          // fall through to handler.next(err)
        }
        handler.next(err);
        return;
      }

      _refreshCompleter = Completer<String?>();
      try {
        final newToken = await _refreshToken();
        _refreshCompleter!.complete(newToken);
        if (newToken != null && newToken.isNotEmpty) {
          err.requestOptions.headers['Authorization'] = 'Bearer $newToken';
          final response = await _dio.fetch(err.requestOptions);
          handler.resolve(response);
          return;
        }
      } on DioException catch (e) {
        _refreshCompleter?.completeError(e);
        // Only wipe tokens if the refresh endpoint itself rejected the refresh token (401/403)
        if (e.response?.statusCode == 401 || e.response?.statusCode == 403) {
          await _clearTokensAndRedirect();
        }
      } catch (e) {
        _refreshCompleter?.completeError(e);
      } finally {
        _refreshCompleter = null;
      }
    }
    handler.next(err);
  }

  Future<String?> _refreshToken() async {
    final refreshToken = await _secureStorage.getRefreshToken();
    if (refreshToken == null || refreshToken.isEmpty) {
      await _clearTokensAndRedirect();
      return null;
    }

    final refreshUrl = _resolveUrl('${Env.baseUrl}/auth/refresh');
    final response = await Dio(
      BaseOptions(
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 15),
      ),
    ).post<Map<String, dynamic>>(
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
    await OfflineStorageCleanupService.purgePaidDownloads();
    await _secureStorage.clearTokens();
  }
}

