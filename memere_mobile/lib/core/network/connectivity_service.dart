import 'dart:async';

import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Wraps `connectivity_plus` (v6, which reports a `List<ConnectivityResult>`)
/// into a simple online/offline signal for the offline-grading + sync features.
///
/// Connectivity is a *reachability* hint, not a guarantee — a device can be on
/// wifi yet unable to reach the API. The attempt/submit paths therefore still
/// treat a `ServerFailure` with code `NO_INTERNET`/`TIMEOUT` (see
/// [isOfflineFailure]) as authoritative; this service just avoids pointless
/// network attempts and drives the sync trigger on reconnect.
class ConnectivityService {
  ConnectivityService([Connectivity? connectivity])
      : _connectivity = connectivity ?? Connectivity();

  final Connectivity _connectivity;

  /// A device can report several transports at once (e.g. wifi + vpn); only an
  /// all-`none` list means offline.
  static bool _isOnline(List<ConnectivityResult> results) =>
      results.any((result) => result != ConnectivityResult.none);

  /// One-shot current status. `true` = at least one active transport.
  Future<bool> isOnline() async {
    final results = await _connectivity.checkConnectivity();
    return _isOnline(results);
  }

  /// Emits the current status immediately, then only on real offline↔online
  /// transitions (transport flaps that don't change the boolean are collapsed).
  Stream<bool> onStatusChange() async* {
    var last = await isOnline();
    yield last;
    await for (final results in _connectivity.onConnectivityChanged) {
      final online = _isOnline(results);
      if (online != last) {
        last = online;
        yield online;
      }
    }
  }
}

final connectivityServiceProvider = Provider<ConnectivityService>((ref) {
  return ConnectivityService();
});

/// Live online/offline status (`true` = online), seeded with the current value.
/// For a one-shot check prefer `ref.read(connectivityServiceProvider).isOnline()`.
final connectivityStreamProvider = StreamProvider<bool>((ref) {
  return ref.watch(connectivityServiceProvider).onStatusChange();
});
