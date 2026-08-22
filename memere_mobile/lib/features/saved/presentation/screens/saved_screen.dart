import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/storage/hive/models/downloaded_item.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../../shared/widgets/memere_mascot.dart';
import '../../../exam/presentation/providers/exam_download_provider.dart';
import '../../../quiz/presentation/providers/quiz_download_provider.dart';

/// The Saved tab surfaces content the user has explicitly downloaded for offline
/// use — quizzes and exams carrying their answer keys, gradeable on-device with
/// no network. Each row can be removed (which deletes the encrypted payload and
/// its manifest entry); the content can always be downloaded again later.
class SavedScreen extends ConsumerWidget {
  const SavedScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final quizzesAsync = ref.watch(downloadedQuizzesProvider);
    final examsAsync = ref.watch(downloadedExamsProvider);

    final downloads = <DownloadedItem>[
      ...quizzesAsync.valueOrNull ?? const [],
      ...examsAsync.valueOrNull ?? const [],
    ]..sort((a, b) => b.downloadedAt.compareTo(a.downloadedAt));

    // Only a genuine first-load (no cached rows yet) warrants a spinner; once
    // either notifier has resolved we render whatever it holds.
    final isLoading = downloads.isEmpty &&
        (quizzesAsync.isLoading || examsAsync.isLoading);

    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: AppPageBackground(
        child: SafeArea(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(
              AppSizes.screenPaddingH,
              AppSizes.lg,
              AppSizes.screenPaddingH,
              AppSizes.xl,
            ),
            children: [
              const Text('Saved', style: AppTextStyles.displayMedium),
              const SizedBox(height: AppSizes.xs),
              Text(
                'Quizzes and exams you downloaded, ready to take offline.',
                style: AppTextStyles.bodyMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(height: AppSizes.xl),
              if (isLoading)
                const Padding(
                  padding: EdgeInsets.only(top: AppSizes.xxl),
                  child: Center(child: CircularProgressIndicator()),
                )
              else if (downloads.isEmpty)
                const _EmptyDownloads()
              else ...[
                AppSectionHeader(
                  title: 'Available offline',
                  subtitle:
                      '${downloads.length} ${downloads.length == 1 ? 'item' : 'items'}',
                ),
                const SizedBox(height: AppSizes.md),
                for (final item in downloads) ...[
                  _DownloadTile(item: item),
                  const SizedBox(height: AppSizes.sm),
                ],
              ],
            ],
          ),
        ),
      ),
    );
  }
}

class _EmptyDownloads extends StatelessWidget {
  const _EmptyDownloads();

  @override
  Widget build(BuildContext context) {
    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.lg),
      radius: AppSizes.radiusXl,
      color: AppColors.bgSecondary,
      shadows: AppShadows.md,
      child: Column(
        children: [
          const MemereMascot(
            size: Size(220, 198),
            showBackdrop: false,
          ),
          const SizedBox(height: AppSizes.lg),
          const Text(
            'Nothing downloaded yet',
            textAlign: TextAlign.center,
            style: AppTextStyles.headlineMedium,
          ),
          const SizedBox(height: AppSizes.sm),
          Text(
            'Download a quiz or exam to take it — and get graded — even without a connection.',
            textAlign: TextAlign.center,
            style: AppTextStyles.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}

class _DownloadTile extends ConsumerWidget {
  const _DownloadTile({required this.item});

  final DownloadedItem item;

  bool get _isQuiz => item.type == DownloadType.quiz;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final accent = _isQuiz ? AppColors.accentPrimary : AppColors.accentSecondary;
    final subtitle = item.subtitle;

    return AppSurface(
      padding: const EdgeInsets.all(AppSizes.md),
      child: Row(
        children: [
          AppIconTile(
            icon: _isQuiz ? Icons.quiz_rounded : Icons.assignment_rounded,
            color: accent,
          ),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    AppBadge(label: _isQuiz ? 'Quiz' : 'Exam', color: accent),
                  ],
                ),
                const SizedBox(height: AppSizes.xs),
                Text(
                  item.title,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: AppTextStyles.bodyMedium.copyWith(
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                if (subtitle != null && subtitle.isNotEmpty) ...[
                  const SizedBox(height: AppSizes.xs),
                  Text(
                    subtitle,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: AppTextStyles.bodySmall.copyWith(
                      color: AppColors.textMuted,
                    ),
                  ),
                ],
              ],
            ),
          ),
          IconButton(
            tooltip: 'Remove download',
            onPressed: () => _confirmRemove(context, ref),
            icon: const Icon(
              Icons.delete_outline_rounded,
              color: AppColors.textMuted,
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _confirmRemove(BuildContext context, WidgetRef ref) async {
    final kind = _isQuiz ? 'quiz' : 'exam';
    final shouldDelete = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: const Color(0xFF141824),
        title: Text(
          'Remove downloaded $kind?',
          style: const TextStyle(color: Colors.white, fontSize: 16),
        ),
        content: Text(
          'This removes the offline copy of "${item.title}" from your device. '
          'You can download it again anytime you\'re online.',
          style: const TextStyle(color: Colors.white70, fontSize: 13),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: ElevatedButton.styleFrom(backgroundColor: AppColors.error),
            child: const Text('Remove', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
    );

    if (shouldDelete != true) return;

    if (_isQuiz) {
      await ref.read(downloadedQuizzesProvider.notifier).remove(item.id);
    } else {
      await ref.read(downloadedExamsProvider.notifier).remove(item.id);
    }

    if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Removed "${item.title}" from downloads.')),
      );
    }
  }
}
