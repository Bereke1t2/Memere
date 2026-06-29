import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

import '../models/offline_video_model.dart';

abstract class OfflineVideoLocalDataSource {
  Future<List<OfflineVideoModel>> getDownloads();
  Future<OfflineVideoModel?> getDownload(String videoId);
  Future<void> saveDownload(OfflineVideoModel video);
  Future<void> removeDownload(String videoId);
  Future<void> clearExpiredDownloads();
}

class OfflineVideoLocalDataSourceImpl implements OfflineVideoLocalDataSource {
  const OfflineVideoLocalDataSourceImpl();

  static const _downloadsKey = 'offline_video_downloads';

  @override
  Future<List<OfflineVideoModel>> getDownloads() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getStringList(_downloadsKey) ?? [];
    return raw
        .map((item) => jsonDecode(item))
        .whereType<Map<String, dynamic>>()
        .map(OfflineVideoModel.fromJson)
        .toList();
  }

  @override
  Future<OfflineVideoModel?> getDownload(String videoId) async {
    final downloads = await getDownloads();
    for (final download in downloads) {
      if (download.videoId == videoId) return download;
    }
    return null;
  }

  @override
  Future<void> saveDownload(OfflineVideoModel video) async {
    final downloads = await getDownloads();
    final next = [
      for (final download in downloads)
        if (download.videoId != video.videoId) download,
      video,
    ];
    await _writeDownloads(next);
  }

  @override
  Future<void> removeDownload(String videoId) async {
    final downloads = await getDownloads();
    await _writeDownloads(
      downloads.where((download) => download.videoId != videoId).toList(),
    );
  }

  @override
  Future<void> clearExpiredDownloads() async {
    final downloads = await getDownloads();
    await _writeDownloads(
      downloads.where((download) => !download.isExpired).toList(),
    );
  }

  Future<void> _writeDownloads(List<OfflineVideoModel> downloads) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setStringList(
      _downloadsKey,
      downloads.map((download) => jsonEncode(download.toJson())).toList(),
    );
  }
}
