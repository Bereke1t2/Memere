import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../domain/entities/offline_video_entity.dart';
import '../providers/offline_video_provider.dart';

class DownloadButton extends ConsumerWidget {
  const DownloadButton({
    super.key,
    required this.videoId,
    required this.lessonId,
    required this.courseId,
    required this.title,
    this.isFullButton = true,
  });

  final String videoId;
  final String lessonId;
  final String courseId;
  final String title;
  final bool isFullButton;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final downloadsAsync = ref.watch(offlineDownloadsProvider);
    final download = downloadsAsync.valueOrNull
        ?.where((item) => item.videoId == videoId)
        .firstOrNull;

    // Live download progress for THIS video (0.0–1.0), or null when not running.
    final activeProgress = ref.watch(offlineDownloadProgressProvider)[videoId];
    final progressPct = activeProgress != null
        ? (activeProgress * 100).clamp(0, 100).round()
        : null;

    final isLoading = downloadsAsync.isLoading;
    final disabled = videoId.trim().isEmpty || isLoading;
    final status = download?.isExpired == true
        ? OfflineVideoStatus.expired
        : download?.status;

    if (!isFullButton) {
      return IconButton(
        tooltip: _tooltip(status),
        onPressed: disabled
            ? null
            : () => _handlePressed(context, ref, download, status),
        icon: _iconFor(status, isLoading, activeProgress),
      );
    }

    // Full Action Button
    final isDownloaded = status == OfflineVideoStatus.downloaded;
    final isDownloading = activeProgress != null ||
        status == OfflineVideoStatus.downloading ||
        status == OfflineVideoStatus.queued ||
        isLoading;

    return ElevatedButton.icon(
      onPressed: disabled ? null : () => _handlePressed(context, ref, download, status),
      style: ElevatedButton.styleFrom(
        backgroundColor: isDownloaded
            ? const Color(0x1810B981)
            : const Color(0xFF1E293B),
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
      icon: isDownloading
          ? SizedBox(
              width: 14,
              height: 14,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                value: activeProgress, // null → indeterminate spinner
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
        isDownloading
            ? (progressPct != null
                ? 'Downloading… $progressPct%'
                : 'Downloading Video...')
            : (isDownloaded
                ? 'Downloaded • Offline Ready'
                : 'Download Video (Watch Offline)'),
        style: TextStyle(
          fontWeight: FontWeight.w600,
          fontSize: 13,
          color: isDownloaded ? AppColors.brandEmerald : Colors.white,
        ),
      ),
    );
  }

  Future<void> _handlePressed(
    BuildContext context,
    WidgetRef ref,
    OfflineVideoEntity? download,
    OfflineVideoStatus? status,
  ) async {
    if (download != null && status == OfflineVideoStatus.downloaded) {
      final shouldDelete = await showDialog<bool>(
        context: context,
        builder: (ctx) => AlertDialog(
          backgroundColor: const Color(0xFF141824),
          title: const Text('Remove Offline Video?', style: TextStyle(color: Colors.white, fontSize: 16)),
          content: const Text(
            'This will remove the downloaded video file from your device storage.',
            style: TextStyle(color: Colors.white70, fontSize: 13),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(ctx, false),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () => Navigator.pop(ctx, true),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFFEF4444)),
              child: const Text('Remove', style: TextStyle(color: Colors.white)),
            ),
          ],
        ),
      );

      if (shouldDelete == true) {
        await ref.read(offlineDownloadsProvider.notifier).removeDownload(videoId);
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Offline video removed.')),
          );
        }
      }
      return;
    }

    if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Starting video download for offline viewing...')),
      );
    }

    await ref.read(offlineDownloadsProvider.notifier).startDownload(
          videoId: videoId,
          lessonId: lessonId,
          courseId: courseId,
          title: title,
        );

    final state = ref.read(offlineDownloadsProvider);
    if (state.hasError && context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(offlineFailureMessage(state.error!))),
      );
    } else if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Video successfully downloaded for offline access! ✓')),
      );
    }
  }

  String _tooltip(OfflineVideoStatus? status) {
    switch (status) {
      case OfflineVideoStatus.downloaded:
        return 'Downloaded (Offline Ready)';
      case OfflineVideoStatus.failed:
        return 'Retry download';
      case OfflineVideoStatus.expired:
        return 'Refresh download';
      case OfflineVideoStatus.queued:
      case OfflineVideoStatus.downloading:
        return 'Downloading';
      case null:
        return 'Download for Offline';
    }
  }

  Widget _iconFor(OfflineVideoStatus? status, bool isLoading,
      [double? activeProgress]) {
    if (activeProgress != null ||
        isLoading ||
        status == OfflineVideoStatus.queued ||
        status == OfflineVideoStatus.downloading) {
      return SizedBox(
        width: 18,
        height: 18,
        child: CircularProgressIndicator(
          strokeWidth: 2,
          value: activeProgress, // null → indeterminate spinner
          color: AppColors.brandEmerald,
        ),
      );
    }

    switch (status) {
      case OfflineVideoStatus.downloaded:
        return const Icon(Icons.check_circle_rounded, color: AppColors.brandEmerald);
      case OfflineVideoStatus.failed:
        return const Icon(Icons.refresh_rounded);
      case OfflineVideoStatus.expired:
        return const Icon(Icons.update_rounded);
      case OfflineVideoStatus.queued:
      case OfflineVideoStatus.downloading:
        return const Icon(Icons.downloading_rounded);
      case null:
        return const Icon(Icons.file_download_outlined);
    }
  }
}
