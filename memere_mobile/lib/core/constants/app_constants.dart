abstract class AppConstants {
  // ── Storage Keys ─────────────────────────────────────────────────────────
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String onboardingSeenKey = 'onboarding_seen';
  static const String userBoxKey = 'user_box';
  static const String courseBoxKey = 'course_box';
  static const String prefsBoxKey = 'prefs_box';
  static const String syncQueueKey = 'sync_queue';

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

  // ── Connectivity ─────────────────────────────────────────────────────────
  static const Duration connectionTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);
}
