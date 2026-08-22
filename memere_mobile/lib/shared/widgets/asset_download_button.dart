import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';

/// Presentational download / remove toggle shared by the quiz and exam download
/// controls. It owns only the tap → confirm-remove interaction; the downloaded/
/// busy state and the actual actions are injected by the feature wrapper, so it
/// carries no network or Riverpod dependency.
class AssetDownloadButton extends StatelessWidget {
  const AssetDownloadButton({
    super.key,
    required this.isDownloaded,
    required this.isBusy,
    required this.onDownload,
    required this.onRemove,
    this.enabled = true,
    this.idleLabel = 'Download to take offline',
    this.downloadedLabel = 'Downloaded • Offline ready',
    this.busyLabel = 'Downloading…',
    this.removeTitle = 'Remove download?',
    this.removeMessage =
        'This removes the offline copy from this device. You can download it again later.',
  });

  final bool isDownloaded;
  final bool isBusy;
  final Future<void> Function() onDownload;
  final Future<void> Function() onRemove;
  final bool enabled;
  final String idleLabel;
  final String downloadedLabel;
  final String busyLabel;
  final String removeTitle;
  final String removeMessage;

  @override
  Widget build(BuildContext context) {
    final disabled = !enabled || isBusy;
    return ElevatedButton.icon(
      onPressed: disabled ? null : () => _handlePressed(context),
      style: ElevatedButton.styleFrom(
        backgroundColor:
            isDownloaded ? const Color(0x1810B981) : const Color(0xFF1E293B),
        foregroundColor: isDownloaded ? AppColors.brandEmerald : Colors.white,
        minimumSize: const Size.fromHeight(44),
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(10),
          side: BorderSide(
            color: isDownloaded
                ? const Color(0x4510B981)
                : const Color(0xFF334155),
          ),
        ),
      ),
      icon: isBusy
          ? const SizedBox(
              width: 14,
              height: 14,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: AppColors.brandEmerald,
              ),
            )
          : Icon(
              isDownloaded
                  ? Icons.check_circle_rounded
                  : Icons.file_download_outlined,
              size: 18,
              color: isDownloaded ? AppColors.brandEmerald : Colors.white,
            ),
      label: Text(
        isBusy
            ? busyLabel
            : isDownloaded
                ? downloadedLabel
                : idleLabel,
        style: TextStyle(
          fontWeight: FontWeight.w600,
          fontSize: 13,
          color: isDownloaded ? AppColors.brandEmerald : Colors.white,
        ),
      ),
    );
  }

  Future<void> _handlePressed(BuildContext context) async {
    if (!isDownloaded) {
      await onDownload();
      return;
    }
    final shouldRemove = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: const Color(0xFF141824),
        title: Text(
          removeTitle,
          style: const TextStyle(color: Colors.white, fontSize: 16),
        ),
        content: Text(
          removeMessage,
          style: const TextStyle(color: Colors.white70, fontSize: 13),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFEF4444),
            ),
            child: const Text('Remove', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
    );
    if (shouldRemove == true) await onRemove();
  }
}
