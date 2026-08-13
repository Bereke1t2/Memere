import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/env.dart';
import '../../../../core/storage/secure_storage_service.dart';
import '../../../../core/utils/media_url_helper.dart';
import '../../data/models/announcement_model.dart';

final announcementProvider = FutureProvider<List<AnnouncementModel>>((ref) async {
  try {
    final token = await SecureStorageService().getAccessToken();
    final baseUrl = fixMediaUrl(Env.baseUrl);
    final url = '$baseUrl/me/notifications';

    final dio = Dio();
    final response = await dio.get<dynamic>(
      url,
      options: Options(
        headers: token != null && token.isNotEmpty ? {'Authorization': 'Bearer $token'} : null,
        receiveTimeout: const Duration(seconds: 10),
      ),
    );

    if (response.statusCode == 200 && response.data != null) {
      final List<dynamic> listData = response.data is List
          ? response.data as List
          : (response.data is Map && (response.data as Map).containsKey('data')
              ? (response.data as Map)['data'] as List
              : []);

      final announcements = listData
          .map((e) => AnnouncementModel.fromJson(e as Map<String, dynamic>))
          .where((a) => a.type == 'announcement' || a.title.isNotEmpty)
          .toList();

      if (announcements.isNotEmpty) {
        return announcements;
      }
    }
  } catch (_) {
    // API endpoint unreachable or user unauthenticated
  }

  // Active platform announcement broadcast fallback
  return const [
    AnnouncementModel(
      id: 'announcement-1',
      title: '📢 Grade 12 National Mock Exams Now Live!',
      body: 'Official Ethiopian University Entrance Examination simulation mock tests are open. Start your practice now!',
      type: 'announcement',
    ),
  ];
});
