import 'dart:io';

import 'package:path_provider/path_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'hive/hive_boxes.dart';

/// Anti-Fraud Offline Storage Cleanup Service.
///
/// When a user logs out or is kicked due to a login from another phone, this service
/// purges all local offline-downloaded paid files (videos, PDF/HTML study notes,
/// offline quizzes & mock exams) from private app sandbox storage and Hive boxes.
/// This prevents users from downloading paid content and accessing it offline after logout.
class OfflineStorageCleanupService {
  OfflineStorageCleanupService._();

  static const String _offlineVideosKey = 'offline_video_downloads';

  /// Purges all offline downloaded paid content across file storage, SharedPreferences, and Hive.
  static Future<void> purgePaidDownloads() async {
    try {
      final appDir = await getApplicationDocumentsDirectory();

      // 1. Delete all offline video files (.mp4) in app sandbox
      final videosDir = Directory('${appDir.path}/offline_videos');
      if (await videosDir.exists()) {
        try {
          await videosDir.delete(recursive: true);
        } catch (_) {}
      }

      // 2. Clear SharedPreferences offline video metadata
      try {
        final prefs = await SharedPreferences.getInstance();
        await prefs.remove(_offlineVideosKey);
      } catch (_) {}

      // 3. Delete all cached PDF & HTML study documents in app sandbox
      final pdfsDir = Directory('${appDir.path}/pdfs');
      if (await pdfsDir.exists()) {
        try {
          await pdfsDir.delete(recursive: true);
        } catch (_) {}
      }

      // 4. Clear Hive boxes holding downloaded quizzes, exams, and course details
      try {
        if (AppHiveBoxes.downloadedQuizzes.isOpen) {
          await AppHiveBoxes.downloadedQuizzes.clear();
        }
      } catch (_) {}

      try {
        if (AppHiveBoxes.downloadedExams.isOpen) {
          await AppHiveBoxes.downloadedExams.clear();
        }
      } catch (_) {}

      try {
        if (AppHiveBoxes.downloadsIndex.isOpen) {
          await AppHiveBoxes.downloadsIndex.clear();
        }
      } catch (_) {}

      try {
        if (AppHiveBoxes.courseDetails.isOpen) {
          await AppHiveBoxes.courseDetails.clear();
        }
      } catch (_) {}
    } catch (_) {}
  }
}
