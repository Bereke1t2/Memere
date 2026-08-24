import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:hive_flutter/hive_flutter.dart';

import '../../constants/app_constants.dart';

/// Central opener + registry for the app's structured offline stores.
///
/// Every box holds JSON strings (`Box<String>`); models hand-write
/// `toJson`/`fromJson` (no `.g.dart` codegen), matching the existing
/// offline-video model idiom. Two boxes carry correct-answer keys for offline
/// grading and are therefore AES-encrypted at rest via [HiveAesCipher], keyed
/// from a secret held in `flutter_secure_storage`.
///
/// Call [openAppHiveBoxes] once after `Hive.initFlutter()` in `main.dart`.
class AppHiveBoxes {
  AppHiveBoxes._();

  static const FlutterSecureStorage _keyStorage = FlutterSecureStorage(
    aOptions: AndroidOptions(),
    iOptions: IOSOptions(accessibility: KeychainAccessibility.first_unlock),
  );

  static bool _opened = false;

  /// Opens every offline box exactly once. Idempotent: a second call is a no-op.
  ///
  /// Each open is fault-tolerant — a box whose on-disk file is unreadable
  /// (e.g. a rotated encryption key or a corrupt file) is deleted and recreated
  /// empty rather than crashing startup. The cost is losing that box's cached
  /// data, which the user can re-download.
  static Future<void> openAll() async {
    if (_opened) return;

    final cipher = HiveAesCipher(await _encryptionKey());

    await Future.wait<Box<String>>([
      _openBox(AppConstants.downloadsIndexBoxKey),
      _openBox(AppConstants.savedItemsBoxKey),
      _openBox(AppConstants.syncQueueKey),
      _openBox(AppConstants.offlineAttemptResultsBoxKey),
      // Course structure cache — no answer keys, so plain (unencrypted).
      _openBox(AppConstants.courseDetailsBoxKey),
      // Encrypted: these carry the answer key (Non-Negotiable #1 relaxation).
      _openBox(AppConstants.downloadedQuizzesBoxKey, cipher: cipher),
      _openBox(AppConstants.downloadedExamsBoxKey, cipher: cipher),
    ]);

    _opened = true;
  }

  // ── Box accessors ──────────────────────────────────────────────────────────
  // All return an already-open box; call [openAll] during startup first.

  static Box<String> get downloadsIndex =>
      Hive.box<String>(AppConstants.downloadsIndexBoxKey);

  static Box<String> get savedItems =>
      Hive.box<String>(AppConstants.savedItemsBoxKey);

  static Box<String> get syncQueue =>
      Hive.box<String>(AppConstants.syncQueueKey);

  static Box<String> get offlineAttemptResults =>
      Hive.box<String>(AppConstants.offlineAttemptResultsBoxKey);

  static Box<String> get courseDetails =>
      Hive.box<String>(AppConstants.courseDetailsBoxKey);

  static Box<String> get downloadedQuizzes =>
      Hive.box<String>(AppConstants.downloadedQuizzesBoxKey);

  static Box<String> get downloadedExams =>
      Hive.box<String>(AppConstants.downloadedExamsBoxKey);

  // ── internals ────────────────────────────────────────────────────────────

  static Future<Box<String>> _openBox(String name,
      {HiveAesCipher? cipher}) async {
    try {
      return await Hive.openBox<String>(name, encryptionCipher: cipher);
    } catch (_) {
      // Corrupt/unreadable box (or key mismatch): reset it rather than block boot.
      await Hive.deleteBoxFromDisk(name);
      return Hive.openBox<String>(name, encryptionCipher: cipher);
    }
  }

  /// Returns the 32-byte AES key for the encrypted boxes, generating and
  /// persisting one on first run. Stored base64url-encoded in secure storage.
  static Future<List<int>> _encryptionKey() async {
    final existing = await _keyStorage.read(key: AppConstants.hiveEncryptionKeyName);
    if (existing != null && existing.isNotEmpty) {
      try {
        final decoded = base64Url.decode(existing);
        if (decoded.length == 32) return decoded;
      } catch (_) {
        // fall through and regenerate
      }
    }
    final key = Hive.generateSecureKey(); // 32 cryptographically-random bytes
    await _keyStorage.write(
      key: AppConstants.hiveEncryptionKeyName,
      value: base64Url.encode(key),
    );
    return key;
  }
}

/// Top-level convenience for `main.dart`.
Future<void> openAppHiveBoxes() => AppHiveBoxes.openAll();
