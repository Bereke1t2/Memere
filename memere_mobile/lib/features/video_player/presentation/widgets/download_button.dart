import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../domain/entities/offline_video_entity.dart';
import '../providers/offline_video_provider.dart';

class DownloadButton extends ConsumerWidget {
  const DownloadButton({
    super.key,
    required this.videoId,
    required this.lessonId,
    required this.courseId,
    required this.title,
  });

  final String videoId;
  final String lessonId;
  final String courseId;
  final String title;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final downloadsAsync = ref.watch(offlineDownloadsProvider);
    final download = downloadsAsync.valueOrNull
        ?.where((item) => item.videoId == videoId)
        .firstOrNull;

    final isLoading = downloadsAsync.isLoading;
    final disabled = videoId.trim().isEmpty || isLoading;
    final status = download?.isExpired == true
        ? OfflineVideoStatus.expired
        : download?.status;

    return IconButton(
      tooltip: _tooltip(status),
      onPressed: disabled
          ? null
          : () => _handlePressed(context, ref, download, status),
      icon: _iconFor(status, isLoading),
    );
  }

  Future<void> _handlePressed(
    BuildContext context,
    WidgetRef ref,
    OfflineVideoEntity? download,
    OfflineVideoStatus? status,
  ) async {
    if (download != null && status == OfflineVideoStatus.downloaded) {
      await ref.read(offlineDownloadsProvider.notifier).removeDownload(videoId);
      return;
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
    }
  }

  String _tooltip(OfflineVideoStatus? status) {
    switch (status) {
      case OfflineVideoStatus.downloaded:
        return 'Remove download';
      case OfflineVideoStatus.failed:
        return 'Retry download';
      case OfflineVideoStatus.expired:
        return 'Refresh download';
      case OfflineVideoStatus.queued:
      case OfflineVideoStatus.downloading:
        return 'Downloading';
      case null:
        return 'Download';
    }
  }

  Widget _iconFor(OfflineVideoStatus? status, bool isLoading) {
    if (isLoading ||
        status == OfflineVideoStatus.queued ||
        status == OfflineVideoStatus.downloading) {
      return const SizedBox(
        width: AppSizes.iconSm,
        height: AppSizes.iconSm,
        child: CircularProgressIndicator(
          strokeWidth: 2,
          color: AppColors.accentPrimary,
        ),
      );
    }

    switch (status) {
      case OfflineVideoStatus.downloaded:
        return const Icon(Icons.download_done_rounded);
      case OfflineVideoStatus.failed:
        return const Icon(Icons.refresh_rounded);
      case OfflineVideoStatus.expired:
        return const Icon(Icons.update_rounded);
      case OfflineVideoStatus.queued:
      case OfflineVideoStatus.downloading:
        return const Icon(Icons.downloading_rounded);
      case null:
        return const Icon(Icons.download_rounded);
    }
  }
}
