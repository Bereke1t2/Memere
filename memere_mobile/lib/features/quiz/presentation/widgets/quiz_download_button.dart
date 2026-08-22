import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../shared/widgets/asset_download_button.dart';
import '../providers/quiz_download_provider.dart';

/// Download / remove control for a single quiz. Downloading a quiz is a LOCAL
/// feature (it stores the quiz — with answer keys — for on-device grading), so
/// it is intentionally NOT wrapped in the account/subscription gate: guests can
/// download published/free quizzes and take them offline.
class QuizDownloadButton extends ConsumerWidget {
  const QuizDownloadButton({
    super.key,
    required this.quizId,
    required this.title,
  });

  final String quizId;
  final String title;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final downloadsAsync = ref.watch(downloadedQuizzesProvider);
    final isDownloaded =
        downloadsAsync.valueOrNull?.any((item) => item.id == quizId) ?? false;
    final isBusy = ref.watch(quizDownloadBusyProvider).contains(quizId) ||
        downloadsAsync.isLoading;

    return AssetDownloadButton(
      isDownloaded: isDownloaded,
      isBusy: isBusy,
      enabled: quizId.trim().isNotEmpty,
      idleLabel: 'Download quiz (take offline)',
      removeMessage:
          'This removes the offline copy of this quiz from your device. You can download it again later.',
      onDownload: () async {
        final failure = await ref
            .read(downloadedQuizzesProvider.notifier)
            .download(quizId: quizId, title: title);
        if (!context.mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              failure == null
                  ? 'Quiz downloaded — you can take it offline. ✓'
                  : failure.message,
            ),
          ),
        );
      },
      onRemove: () async {
        await ref.read(downloadedQuizzesProvider.notifier).remove(quizId);
        if (!context.mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Quiz download removed.')),
        );
      },
    );
  }
}
