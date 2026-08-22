import '../hive_json.dart';

/// Whether a queued offline submission is a quiz or an exam attempt — selects
/// the endpoints the sync service replays it through.
enum SubmissionKind { quiz, exam }

SubmissionKind submissionKindFromString(String v) =>
    v.toLowerCase() == 'exam' ? SubmissionKind.exam : SubmissionKind.quiz;

/// A quiz/exam attempt taken offline and awaiting authoritative server grading.
///
/// Holds the RAW answers (never a score the server would trust) plus the
/// provisional on-device score shown to the user immediately. On reconnect the
/// sync service replays start→submit through the normal endpoints so the server
/// grades authoritatively; the provisional result is then replaced.
class PendingSubmission {
  const PendingSubmission({
    required this.localId,
    required this.kind,
    required this.contentId,
    required this.rawAnswers,
    required this.takenAt,
    this.provisionalScore,
    this.provisionalPct,
    this.serverAttemptId,
    this.attempts = 0,
  });

  /// Client-generated id; also the key into the `offline_attempt_results` box.
  final String localId;
  final SubmissionKind kind;

  /// The quiz or exam id this attempt belongs to.
  final String contentId;

  /// questionId → selected answer id(s), exactly as the submit endpoint expects.
  final Map<String, dynamic> rawAnswers;

  /// The provisional on-device score/percentage (informational until synced).
  final int? provisionalScore;
  final double? provisionalPct;
  final DateTime takenAt;

  /// Set once the server-side attempt has been created during replay.
  final String? serverAttemptId;

  /// Replay attempt counter, for backoff / give-up policy in the sync service.
  final int attempts;

  PendingSubmission copyWith({
    String? serverAttemptId,
    int? attempts,
  }) =>
      PendingSubmission(
        localId: localId,
        kind: kind,
        contentId: contentId,
        rawAnswers: rawAnswers,
        takenAt: takenAt,
        provisionalScore: provisionalScore,
        provisionalPct: provisionalPct,
        serverAttemptId: serverAttemptId ?? this.serverAttemptId,
        attempts: attempts ?? this.attempts,
      );

  factory PendingSubmission.fromJson(Map<String, dynamic> json) =>
      PendingSubmission(
        localId: hstr(json['local_id']),
        kind: submissionKindFromString(hstr(json['kind'])),
        contentId: hstr(json['content_id']),
        rawAnswers: json['raw_answers'] is Map
            ? Map<String, dynamic>.from(json['raw_answers'] as Map)
            : <String, dynamic>{},
        provisionalScore: hintOrNull(json['provisional_score']),
        provisionalPct: json['provisional_pct'] == null
            ? null
            : hdouble(json['provisional_pct']),
        takenAt: hdate(json['taken_at']),
        serverAttemptId: hstrOrNull(json['server_attempt_id']),
        attempts: hint(json['attempts']),
      );

  Map<String, dynamic> toJson() => {
        'local_id': localId,
        'kind': kind.name,
        'content_id': contentId,
        'raw_answers': rawAnswers,
        'provisional_score': provisionalScore,
        'provisional_pct': provisionalPct,
        'taken_at': takenAt.toIso8601String(),
        'server_attempt_id': serverAttemptId,
        'attempts': attempts,
      };
}
