import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:webview_flutter/webview_flutter.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/router/app_router.dart';
import '../widgets/payment_empty_state.dart';

/// Hosts the provider's secure checkout page. The app never reads payment form
/// data from the page — completion is decided by backend status polling on the
/// result screen, not by the WebView URL alone.
class PaymentWebViewScreen extends ConsumerStatefulWidget {
  const PaymentWebViewScreen({
    super.key,
    required this.paymentId,
    required this.redirectUrl,
    required this.courseId,
  });

  final String paymentId;
  final String redirectUrl;
  final String courseId;

  @override
  ConsumerState<PaymentWebViewScreen> createState() =>
      _PaymentWebViewScreenState();
}

class _PaymentWebViewScreenState extends ConsumerState<PaymentWebViewScreen> {
  WebViewController? _controller;
  bool _loading = true;
  bool _navigated = false;
  String? _urlError;

  // URL fragments that indicate the provider has handed control back to us.
  static const _returnMarkers = [
    'payment/callback',
    'payment/return',
    'payments/return',
    'payment-success',
    'payment-failed',
    'payment-cancel',
    'checkout/complete',
    'return',
    'cancel',
  ];

  @override
  void initState() {
    super.initState();
    _initController();
  }

  bool _isAllowedUrl(Uri uri) {
    if (uri.scheme == 'https') return true;
    // Allow plain http only for local/mock dev hosts in debug builds.
    if (!kReleaseMode &&
        uri.scheme == 'http' &&
        (uri.host == 'localhost' ||
            uri.host == '10.0.2.2' ||
            uri.host == '127.0.0.1')) {
      return true;
    }
    return false;
  }

  void _initController() {
    final uri = Uri.tryParse(widget.redirectUrl);
    if (uri == null || !_isAllowedUrl(uri)) {
      setState(() {
        _urlError = 'This checkout link is invalid or insecure.';
        _loading = false;
      });
      return;
    }

    _controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setBackgroundColor(AppColors.bgPrimary)
      ..setNavigationDelegate(
        NavigationDelegate(
          onProgress: (_) {},
          onPageStarted: (_) {
            if (mounted) setState(() => _loading = true);
          },
          onPageFinished: (_) {
            if (mounted) setState(() => _loading = false);
          },
          onWebResourceError: (error) {
            if (mounted) {
              setState(() {
                _loading = false;
                _urlError = 'Could not load the checkout page.';
              });
            }
          },
          onNavigationRequest: (request) {
            if (_looksLikeReturn(request.url)) {
              _finishToResult();
              return NavigationDecision.prevent;
            }
            return NavigationDecision.navigate;
          },
        ),
      )
      ..loadRequest(uri);
  }

  bool _looksLikeReturn(String url) {
    final lower = url.toLowerCase();
    return _returnMarkers.any(lower.contains);
  }

  /// Leaves the WebView and lets the result screen poll the backend for the
  /// authoritative status. Guarded so we only navigate once.
  void _finishToResult() {
    if (_navigated || !mounted) return;
    _navigated = true;
    context.pushReplacement(
      AppRoutes.paymentResultPath(
        paymentId: widget.paymentId,
        courseId: widget.courseId,
      ),
    );
  }

  Future<void> _confirmClose() async {
    final shouldClose = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        backgroundColor: AppColors.bgSecondary,
        title:
            const Text('Leave checkout?', style: AppTextStyles.headlineSmall),
        content: Text(
          'If you already paid, we will confirm it on the next screen.',
          style:
              AppTextStyles.bodyMedium.copyWith(color: AppColors.textSecondary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text('Stay'),
          ),
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: const Text('Leave'),
          ),
        ],
      ),
    );
    if (shouldClose == true) _finishToResult();
  }

  @override
  Widget build(BuildContext context) {
    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) {
        if (!didPop) _confirmClose();
      },
      child: Scaffold(
        backgroundColor: AppColors.bgPrimary,
        appBar: AppBar(
          title: const Text('Complete payment'),
          leading: IconButton(
            icon: const Icon(Icons.close_rounded),
            onPressed: _confirmClose,
          ),
        ),
        body: SafeArea(
          child: _urlError != null
              ? PaymentEmptyState(
                  icon: Icons.lock_outline_rounded,
                  title: 'Checkout unavailable',
                  body: _urlError!,
                  buttonLabel: 'Check payment status',
                  onPressed: _finishToResult,
                )
              : Stack(
                  children: [
                    if (_controller != null)
                      WebViewWidget(controller: _controller!),
                    if (_loading)
                      const LinearProgressIndicator(
                        color: AppColors.accentPrimary,
                        backgroundColor: AppColors.bgTertiary,
                      ),
                  ],
                ),
        ),
      ),
    );
  }
}
