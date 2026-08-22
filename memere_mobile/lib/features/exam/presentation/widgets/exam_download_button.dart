import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../shared/widgets/asset_download_button.dart';
import '../providers/exam_download_provider.dart';

/// Download / remove control for a single mock exam. Downloading an exam is a
/// LOCAL feature (it stores the exam — with answer keys — for on-device
/// grading), so it is intentionally NOT wrapped in the account/subscription
/// gate: guests can download published/free exams and take them offline.
class ExamDownloadButton extends ConsumerWidget {
  const ExamDownloadButton({
    super.key,
    required this.examId,
    required this.title,
  });

  final String examId;
  final String title;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final downloadsAsync = ref.watch(downloadedExamsProvider);
    final isDownloaded =
        downloadsAsync.valueOrNull?.any((item) => item.id == examId) ?? false;
    final isBusy = ref.watch(examDownloadBusyProvider).contains(examId) ||
        downloadsAsync.isLoading;

    return AssetDownloadButton(
      isDownloaded: isDownloaded,
      isBusy: isBusy,
      enabled: examId.trim().isNotEmpty,
      idleLabel: 'Download exam (take offline)',
      removeMessage:
          'This removes the offline copy of this exam from your device. You can download it again later.',
      onDownload: () async {
        final failure = await ref
            .read(downloadedExamsProvider.notifier)
            .download(examId: examId, title: title);
        if (!context.mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              failure == null
                  ? 'Exam downloaded — you can take it offline. ✓'
                  : failure.message,
            ),
          ),
        );
      },
      onRemove: () async {
        await ref.read(downloadedExamsProvider.notifier).remove(examId);
        if (!context.mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Exam download removed.')),
        );
      },
    );
  }
}
