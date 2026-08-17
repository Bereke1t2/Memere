import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/network/dio_client.dart';
import '../../data/models/announcement_model.dart';

/// Provider that fetches real announcements published by the admin on the backend.
/// Filters `/me/notifications` for type == 'announcement' or general system broadcasts.
final announcementProvider = FutureProvider<List<AnnouncementModel>>((ref) async {
  try {
    final dioClient = ref.watch(dioClientProvider);
    final response = await dioClient.get<dynamic>('/me/notifications?limit=50');

    if (response.statusCode == 200 && response.data != null) {
      final List<dynamic> listData = response.data is List
          ? response.data as List
          : (response.data is Map && (response.data as Map).containsKey('notifications')
              ? (response.data as Map)['notifications'] as List
              : (response.data is Map && (response.data as Map).containsKey('data')
                  ? (response.data as Map)['data'] as List
                  : []));

      final announcements = listData
          .map((e) => AnnouncementModel.fromJson(e as Map<String, dynamic>))
          .where((a) =>
              a.type.toLowerCase() == 'announcement' ||
              (a.type.toLowerCase() == 'system' && a.title.isNotEmpty))
          .toList();

      // Sort by creation date (newest first)
      announcements.sort((a, b) {
        if (a.createdAt == null && b.createdAt == null) return 0;
        if (a.createdAt == null) return 1;
        if (b.createdAt == null) return -1;
        return b.createdAt!.compareTo(a.createdAt!);
      });

      return announcements;
    }
  } catch (_) {
    // API endpoint unreachable or user unauthenticated
  }

  return const [];
});

