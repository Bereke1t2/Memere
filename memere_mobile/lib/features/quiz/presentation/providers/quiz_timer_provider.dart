import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

class QuizTimerSeed {
  const QuizTimerSeed({
    required this.attemptId,
    this.expiresAt,
    this.remainingSeconds,
  });

  final String attemptId;
  final DateTime? expiresAt;
  final int? remainingSeconds;

  @override
  bool operator ==(Object other) {
    return other is QuizTimerSeed &&
        other.attemptId == attemptId &&
        other.expiresAt == expiresAt &&
        other.remainingSeconds == remainingSeconds;
  }

  @override
  int get hashCode => Object.hash(attemptId, expiresAt, remainingSeconds);
}

final quizTimerProvider =
    StreamProvider.autoDispose.family<int?, QuizTimerSeed>((ref, seed) {
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

int? _secondsRemaining(QuizTimerSeed seed) {
  if (seed.expiresAt != null) {
    final seconds = seed.expiresAt!.difference(DateTime.now()).inSeconds;
    return seconds < 0 ? 0 : seconds;
  }
  if (seed.remainingSeconds == null) return null;
  return seed.remainingSeconds! < 0 ? 0 : seed.remainingSeconds;
}
