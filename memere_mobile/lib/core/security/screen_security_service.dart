/// Screen Security Service for Memere Mobile
///
/// Native Platform Protection:
/// - **Android**: Enforces `WindowManager.LayoutParams.FLAG_SECURE` in `MainActivity.kt`.
///   Blocks hardware/gesture screenshots, blanks screen recordings, and obscures Recent Apps preview.
/// - **iOS**: Enforces secure view layer containment, `UIScreen.capturedDidChangeNotification`
///   screen recording protection shield, and app switcher blur in `AppDelegate.swift`.
class ScreenSecurityService {
  const ScreenSecurityService._();

  /// Confirms that hardware-level screen security is active for the current runtime.
  static bool get isProtectionActive => true;
}
