import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../domain/entities/video_status_entity.dart';
import '../providers/video_player_controller_provider.dart';
import '../widgets/hls_video_player.dart';
import '../widgets/video_action_bar.dart';
import '../widgets/video_player_error_state.dart';

class VideoPlayerScreen extends ConsumerWidget {
  const VideoPlayerScreen({
    super.key,
    required this.videoId,
    required this.lessonId,
    required this.courseId,
    this.title,
  });

  final String videoId;
  final String lessonId;
  final String courseId;
  final String? title;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final params = VideoPlayerParams(
      videoId: videoId,
      lessonId: lessonId,
      courseId: courseId,
    );
    final playbackAsync = ref.watch(videoPlayerControllerProvider(params));

    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        leading: IconButton(
          tooltip: 'Back',
          onPressed: context.pop,
          icon: const Icon(Icons.arrow_back_rounded),
        ),
        title: Text(
          title == null || title!.trim().isEmpty ? 'Lesson video' : title!,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
      ),
      body: SafeArea(
        top: false,
        child: playbackAsync.when(
          loading: () => const _VideoLoadingState(),
          error: (error, _) => VideoPlayerErrorState(
            title: 'Could not load video',
            message: videoFailureMessage(error),
            onRetry: () {
              ref.read(videoPlayerControllerProvider(params).notifier).retry();
            },
            onBack: context.pop,
          ),
          data: (playback) {
            if (playback.errorMessage != null) {
              return VideoPlayerErrorState(
                title: 'Video unavailable',
                message: playback.errorMessage!,
                onRetry: () {
                  ref
                      .read(videoPlayerControllerProvider(params).notifier)
                      .retry();
                },
                onBack: context.pop,
              );
            }

            final status = playback.status;
            if (status != null && !status.isReady) {
              return _VideoStatusState(
                status: status,
                onRetry: () {
                  ref
                      .read(videoPlayerControllerProvider(params).notifier)
                      .retry();
                },
              );
            }

            final controller = playback.videoController;
            final chewieController = playback.chewieController;
            if (controller == null || chewieController == null) {
              return VideoPlayerErrorState(
                title: 'Video unavailable',
                message: 'The stream could not be initialized.',
                onRetry: () {
                  ref
                      .read(videoPlayerControllerProvider(params).notifier)
                      .retry();
                },
                onBack: context.pop,
              );
            }

            return ListView(
              children: [
                HlsVideoPlayer(
                  controller: controller,
                  chewieController: chewieController,
                ),
                Padding(
                  padding: const EdgeInsets.all(AppSizes.screenPaddingH),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        title == null || title!.trim().isEmpty
                            ? 'Lesson video'
                            : title!,
                        style: AppTextStyles.headlineSmall,
                      ),
                      const SizedBox(height: AppSizes.xs),
                      Text(
                        playback.isCompleted
                            ? 'Completed'
                            : 'Progress saves while you watch.',
                        style: AppTextStyles.bodySmall.copyWith(
                          color: playback.isCompleted
                              ? AppColors.success
                              : AppColors.textSecondary,
                        ),
                      ),
                      const SizedBox(height: AppSizes.md),
                      VideoActionBar(
                        videoId: videoId,
                        lessonId: lessonId,
                        courseId: courseId,
                        title: title ?? 'Lesson video',
                        isSavingProgress: playback.isSavingProgress,
                        onRefresh: () {
                          ref
                              .read(
                                videoPlayerControllerProvider(params).notifier,
                              )
                              .refreshStream();
                        },
                        onMarkComplete: () {
                          ref
                              .read(
                                videoPlayerControllerProvider(params).notifier,
                              )
                              .markComplete();
                        },
                      ),
                    ],
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }
}

class _VideoLoadingState extends StatelessWidget {
  const _VideoLoadingState();

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: CircularProgressIndicator(
        color: AppColors.accentPrimary,
        strokeWidth: 2,
      ),
    );
  }
}

class _VideoStatusState extends StatelessWidget {
  const _VideoStatusState({
    required this.status,
    required this.onRetry,
  });

  final VideoStatusEntity status;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    final failed = status.isFailed;
    return VideoPlayerErrorState(
      title: failed ? 'Video processing failed' : 'Video is not ready yet',
      message: failed
          ? status.error ?? 'This video could not be prepared for playback.'
          : 'This lesson video is still processing. Please check again later.',
      onRetry: onRetry,
      onBack: context.pop,
    );
  }
}
