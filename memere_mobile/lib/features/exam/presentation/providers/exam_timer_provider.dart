import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

class ExamTimerSeed {
  const ExamTimerSeed({
    required this.attemptId,
    this.expiresAt,
    this.remainingSeconds,
  });

  final String attemptId;
  final DateTime? expiresAt;
  final int? remainingSeconds;

  @override
  bool operator ==(Object other) {
    return other is ExamTimerSeed &&
        other.attemptId == attemptId &&
        other.expiresAt == expiresAt &&
        other.remainingSeconds == remainingSeconds;
  }

  @override
  int get hashCode => Object.hash(attemptId, expiresAt, remainingSeconds);
}

/// Seconds remaining for display only. Server `expiresAt`/`remainingSeconds`
/// remain authoritative; this never extends or resets time locally.
final examTimerProvider =
    StreamProvider.autoDispose.family<int?, ExamTimerSeed>((ref, seed) {
  final initial = _secondsRemaining(seed);
  if (initial == null) return Stream<int?>.value(null);

  late final StreamController<int?> controller;
  Timer? timer;

  controller = StreamController<int?>(
    onListen: () {
      controller.add(_secondsRemaining(seed));
      timer = Timer.periodic(const Duration(seconds: 1), (_) {
        controller.add(_secondsRemaining(seed));
      });
    },
    onCancel: () {
      timer?.cancel();
    },
  );

  ref.onDispose(() {
    timer?.cancel();
    controller.close();
  });

  return controller.stream;
});

int? _secondsRemaining(ExamTimerSeed seed) {
  if (seed.expiresAt != null) {
    final seconds = seed.expiresAt!.difference(DateTime.now()).inSeconds;
    return seconds < 0 ? 0 : seconds;
  }
  if (seed.remainingSeconds == null) return null;
  return seed.remainingSeconds! < 0 ? 0 : seed.remainingSeconds;
}
