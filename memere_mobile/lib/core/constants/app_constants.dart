abstract class AppConstants {
  // ── Storage Keys ─────────────────────────────────────────────────────────
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String onboardingSeenKey = 'onboarding_seen';
  static const String userBoxKey = 'user_box';
  static const String courseBoxKey = 'course_box';
  static const String prefsBoxKey = 'prefs_box';
  static const String syncQueueKey = 'sync_queue';

  // ── Offline stores (Hive) ────────────────────────────────────────────────
  // Structured, queryable offline data. All hold JSON strings (Box<String>).
  static const String downloadsIndexBoxKey = 'downloads_index';
  static const String savedItemsBoxKey = 'saved_items';
  static const String offlineAttemptResultsBoxKey = 'offline_attempt_results';
  // Encrypted boxes: these carry correct-answer keys for on-device grading
  // (the narrow, sanctioned relaxation of "answers never sent to client").
  static const String downloadedQuizzesBoxKey = 'downloaded_quizzes';
  static const String downloadedExamsBoxKey = 'downloaded_exams';
  // Secret name under flutter_secure_storage holding the base64 AES key that
  // encrypts the two boxes above (HiveAesCipher).
  static const String hiveEncryptionKeyName = 'hive_encryption_key';

  // ── Cache TTL ────────────────────────────────────────────────────────────
  static const Duration cacheTtl = Duration(hours: 1);
  static const Duration downloadedVideoTtl = Duration(days: 30);

  // ── Pagination ───────────────────────────────────────────────────────────
  static const int defaultPageLimit = 20;

  // ── Video ────────────────────────────────────────────────────────────────
  static const int videoProgressSaveIntervalSeconds = 30;

  // ── Auth ─────────────────────────────────────────────────────────────────
  static const int accessTokenTtlMinutes = 15;
  static const int refreshTokenTtlDays = 30;
  static const int maxLoginAttempts = 5;

  // ── Exam ─────────────────────────────────────────────────────────────────
  static const int examAutoSaveIntervalSeconds = 30;

  // ── Payments ─────────────────────────────────────────────────────────────
  static const int paymentPollIntervalSeconds = 3;
  static const int paymentPollTimeoutSeconds = 120;

  // ── Connectivity ─────────────────────────────────────────────────────────
  static const Duration connectionTimeout = Duration(seconds: 5);
  static const Duration receiveTimeout = Duration(seconds: 15);
}
