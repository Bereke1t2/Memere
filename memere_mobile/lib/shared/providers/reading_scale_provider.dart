import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// How much to scale assessment body text, chosen by the student.
///
/// Imported exam questions run long and carry dense notation, and students sit
/// them on every size of phone. A fixed type size that suits one reader leaves
/// another squinting through a 100-question paper, so the size is theirs to set
/// and is remembered between sessions.
///
/// The scale deliberately does not follow the OS text-size setting: a student
/// who enlarges system text for everyday use often still wants a different size
/// for a dense maths paper, and the reverse.
class ReadingScale extends StateNotifier<double> {
  ReadingScale(this._prefs) : super(defaultScale) {
    _restore();
  }

  final Future<SharedPreferences> _prefs;

  static const double defaultScale = 1.0;
  static const double minScale = 0.9;
  static const double maxScale = 1.45;
  static const double _step = 0.1;
  static const String _key = 'assessment_reading_scale';

  bool get canDecrease => state > minScale + 0.001;
  bool get canIncrease => state < maxScale - 0.001;

  /// Human-readable size for the control, e.g. "110%".
  String get label => '${(state * 100).round()}%';

  Future<void> _restore() async {
    final prefs = await _prefs;
    final stored = prefs.getDouble(_key);
    if (stored != null) state = _clamp(stored);
  }

  Future<void> increase() => _set(state + _step);

  Future<void> decrease() => _set(state - _step);

  Future<void> reset() => _set(defaultScale);

  Future<void> _set(double value) async {
    final next = _clamp(value);
    if (next == state) return;
    state = next;
    final prefs = await _prefs;
    await prefs.setDouble(_key, next);
  }

  static double _clamp(double v) =>
      v < minScale ? minScale : (v > maxScale ? maxScale : v);
}

final _preferencesProvider = FutureProvider<SharedPreferences>(
  (ref) => SharedPreferences.getInstance(),
);

final readingScaleProvider = StateNotifierProvider<ReadingScale, double>((ref) {
  return ReadingScale(ref.watch(_preferencesProvider.future));
});
