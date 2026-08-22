import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/network/connectivity_service.dart';
import 'core/offline/sync_service.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';

class MemereApp extends ConsumerStatefulWidget {
  const MemereApp({super.key});

  @override
  ConsumerState<MemereApp> createState() => _MemereAppState();
}

class _MemereAppState extends ConsumerState<MemereApp>
    with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    // Coming back to the foreground is a natural moment to flush any offline
    // submissions queued while the app was backgrounded or the network was down.
    if (state == AppLifecycleState.resumed) {
      unawaited(ref.read(syncServiceProvider).drain());
    }
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(appRouterProvider);

    // Drain the offline-submission queue on every offline→online transition.
    // The stream seeds its current value on first listen, so a cold start while
    // online also triggers one drain to clear any leftover queue. The service
    // itself no-ops for guests, when offline, or while a drain is in flight.
    ref.listen<AsyncValue<bool>>(connectivityStreamProvider, (previous, next) {
      final wasOnline = previous?.valueOrNull ?? false;
      final isOnline = next.valueOrNull ?? false;
      if (!wasOnline && isOnline) {
        unawaited(ref.read(syncServiceProvider).drain());
      }
    });

    return MaterialApp.router(
      title: 'Memere',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.dark(),
      darkTheme: AppTheme.dark(),
      themeMode: ThemeMode.dark,
      routerConfig: router,
    );
  }
}
