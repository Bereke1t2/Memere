import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// Questions a student has marked to come back to, for one attempt.
///
/// Flagging is a study aid, not graded state, so it never leaves the device and
/// never reaches the grading API. It is persisted per attempt rather than held
/// in memory because an attempt survives the app being closed — a student who
/// flags six questions, takes a call and comes back should still have their six
/// flags.
///
/// Keyed by attempt id so two attempts in flight never share flags.
class FlaggedQuestions extends StateNotifier<Set<String>> {
  FlaggedQuestions(this._attemptId, this._prefs) : super(const {}) {
    _restore();
  }

  final String _attemptId;
  final Future<SharedPreferences> _prefs;

  String get _key => 'flagged_questions_$_attemptId';

  Future<void> _restore() async {
    final prefs = await _prefs;
    final stored = prefs.getStringList(_key);
    if (stored != null && stored.isNotEmpty) state = stored.toSet();
  }

  bool isFlagged(String questionId) => state.contains(questionId);

  Future<void> toggle(String questionId) async {
    final next = Set<String>.from(state);
    if (!next.remove(questionId)) next.add(questionId);
    state = next;
    final prefs = await _prefs;
    await prefs.setStringList(_key, next.toList());
  }

  /// Drops every flag for this attempt. Called once the attempt is submitted,
  /// so flags do not accumulate in storage for attempts that are over.
  Future<void> clear() async {
    state = const {};
    final prefs = await _prefs;
    await prefs.remove(_key);
  }
}

final _preferencesProvider = FutureProvider<SharedPreferences>(
  (ref) => SharedPreferences.getInstance(),
);

final flaggedQuestionsProvider =
    StateNotifierProvider.family<FlaggedQuestions, Set<String>, String>(
  (ref, attemptId) =>
      FlaggedQuestions(attemptId, ref.watch(_preferencesProvider.future)),
);
