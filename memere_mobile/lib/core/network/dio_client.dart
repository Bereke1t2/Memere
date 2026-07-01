import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../constants/app_constants.dart';
import '../constants/env.dart';
import 'interceptors/auth_interceptor.dart';
import 'interceptors/logging_interceptor.dart';

final dioClientProvider = Provider<DioClient>((ref) {
  return DioClient();
});

class DioClient {
  DioClient() {
    _dio = Dio(
      BaseOptions(
        baseUrl: Env.baseUrl,
        connectTimeout: AppConstants.connectionTimeout,
        receiveTimeout: AppConstants.receiveTimeout,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );
    _dio.interceptors.addAll([
      AuthInterceptor(_dio),
      LoggingInterceptor(),
    ]);
  }

  late final Dio _dio;

  Dio get dio => _dio;

  Future<Response<T>> get<T>(String path,
          {Map<String, dynamic>? queryParameters, Options? options}) =>
      _dio.get<T>(path, queryParameters: queryParameters, options: options);

  Future<Response<T>> post<T>(String path,
          {dynamic data,
          Map<String, dynamic>? queryParameters,
          Options? options}) =>
      _dio.post<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );

  Future<Response<T>> put<T>(String path, {dynamic data, Options? options}) =>
      _dio.put<T>(path, data: data, options: options);

  Future<Response<T>> patch<T>(String path, {dynamic data, Options? options}) =>
      _dio.patch<T>(path, data: data, options: options);

  Future<Response<T>> delete<T>(String path, {Options? options}) =>
      _dio.delete<T>(path, options: options);
}
