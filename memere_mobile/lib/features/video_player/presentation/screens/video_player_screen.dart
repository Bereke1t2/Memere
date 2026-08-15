import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/ai_robot_mascot.dart';
import '../../domain/entities/video_status_entity.dart';
import '../providers/video_player_controller_provider.dart';
import '../widgets/hls_video_player.dart';
import '../widgets/video_action_bar.dart';
import '../widgets/video_player_error_state.dart';

class VideoPlayerScreen extends ConsumerStatefulWidget {
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
  ConsumerState<VideoPlayerScreen> createState() => _VideoPlayerScreenState();
}

class _VideoPlayerScreenState extends ConsumerState<VideoPlayerScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  void _openAiTutor() {
    final lessonTitle = (widget.title == null || widget.title!.trim().isEmpty)
        ? 'Lesson Video'
        : widget.title!;

    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Row(
                children: [
                  AiRobotMascot(size: 32),
                  SizedBox(width: 10),
                  Text(
                    'AI Concept Tutor',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textPrimary,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 14),
              Text(
                'Asking AI Tutor about "$lessonTitle"...\nGet instant step-by-step explanations of video formulas and entrance exam topics!',
                style: const TextStyle(fontSize: 13, color: AppColors.textMuted),
              ),
              const SizedBox(height: 20),
              ElevatedButton.icon(
                onPressed: () => Navigator.pop(context),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.brandEmerald,
                  foregroundColor: Colors.white,
                  minimumSize: const Size(double.infinity, 44),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                icon: const Icon(Icons.smart_toy_rounded, size: 18),
                label: const Text('Start AI Discussion'),
              ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final params = VideoPlayerParams(
      videoId: widget.videoId,
      lessonId: widget.lessonId,
      courseId: widget.courseId,
    );
    final playbackAsync = ref.watch(videoPlayerControllerProvider(params));
    final lessonTitle = (widget.title == null || widget.title!.trim().isEmpty)
        ? 'Lesson Video'
        : widget.title!;

    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        elevation: 0,
        leading: IconButton(
          tooltip: 'Back to Course',
          onPressed: context.pop,
          icon: const Icon(Icons.arrow_back_rounded, color: Colors.white),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              lessonTitle,
              style: AppTextStyles.titleMedium.copyWith(
                color: Colors.white,
                fontWeight: FontWeight.bold,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const Text(
              'Grade 12 National Exam Prep • 1080p HD',
              style: TextStyle(fontSize: 11, color: AppColors.textMuted),
            ),
          ],
        ),
      ),
      floatingActionButton: AiTutorFab(onPressed: _openAiTutor),
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
                  ref.read(videoPlayerControllerProvider(params).notifier).retry();
                },
                onBack: context.pop,
              );
            }

            final status = playback.status;
            if (status != null && !status.isReady) {
              return _VideoStatusState(
                status: status,
                onRetry: () {
                  ref.read(videoPlayerControllerProvider(params).notifier).retry();
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
                  ref.read(videoPlayerControllerProvider(params).notifier).retry();
                },
                onBack: context.pop,
              );
            }

            return Column(
              children: [
                // Standard 16:9 Widescreen HD Video Player
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: AppSizes.sm),
                  child: HlsVideoPlayer(
                    controller: controller,
                    chewieController: chewieController,
                  ),
                ),
                const SizedBox(height: AppSizes.md),

                // Scrollable Content Area below Video
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: AppSizes.screenPaddingH),
                    decoration: const BoxDecoration(
                      color: Color(0xFF0B0F17),
                      borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const SizedBox(height: AppSizes.md),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    lessonTitle,
                                    style: AppTextStyles.titleMedium.copyWith(
                                      fontWeight: FontWeight.bold,
                                      color: Colors.white,
                                    ),
                                    maxLines: 2,
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                  const SizedBox(height: 4),
                                  Row(
                                    children: [
                                      Container(
                                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                                        decoration: BoxDecoration(
                                          color: const Color(0xFF1E293B),
                                          borderRadius: BorderRadius.circular(4),
                                          border: Border.all(color: const Color(0xFF334155)),
                                        ),
                                        child: const Text(
                                          '1080P HD',
                                          style: TextStyle(
                                            fontSize: 9,
                                            fontWeight: FontWeight.w600,
                                            color: Color(0xFF94A3B8),
                                          ),
                                        ),
                                      ),
                                      const SizedBox(width: 8),
                                      Text(
                                        playback.isCompleted
                                            ? 'Status: Completed ✓'
                                            : 'Auto-saves progress',
                                        style: AppTextStyles.caption.copyWith(
                                          color: playback.isCompleted
                                              ? const Color(0xFF10B981)
                                              : const Color(0xFF64748B),
                                          fontSize: 11,
                                        ),
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: AppSizes.md),

                        // Action buttons bar
                        VideoActionBar(
                          videoId: widget.videoId,
                          lessonId: widget.lessonId,
                          courseId: widget.courseId,
                          title: lessonTitle,
                          isSavingProgress: playback.isSavingProgress,
                          isCompleted: playback.isCompleted,
                          onRefresh: () {
                            ref
                                .read(videoPlayerControllerProvider(params).notifier)
                                .refreshStream();
                          },
                          onMarkComplete: () {
                            ref
                                .read(videoPlayerControllerProvider(params).notifier)
                                .markComplete();
                          },
                        ),
                        const SizedBox(height: AppSizes.md),
                        const Divider(color: AppColors.border),

                        // Tab Bar
                        TabBar(
                          controller: _tabController,
                          indicatorColor: AppColors.brandEmerald,
                          labelColor: AppColors.textPrimary,
                          unselectedLabelColor: AppColors.textMuted,
                          labelStyle: AppTextStyles.labelMedium.copyWith(fontWeight: FontWeight.bold),
                          tabs: const [
                            Tab(text: 'Overview & Notes'),
                            Tab(text: 'Discussion & Q&A'),
                          ],
                        ),
                        const SizedBox(height: AppSizes.md),

                        // Tab Views
                        Expanded(
                          child: TabBarView(
                            controller: _tabController,
                            children: [
                              SingleChildScrollView(
                                child: Container(
                                  padding: const EdgeInsets.all(AppSizes.md),
                                  decoration: BoxDecoration(
                                    color: AppColors.bgSecondary,
                                    borderRadius: BorderRadius.circular(AppSizes.radiusLg),
                                    border: Border.all(color: AppColors.border),
                                  ),
                                  child: Column(
                                    crossAxisAlignment: CrossAxisAlignment.start,
                                    children: [
                                      const Text(
                                        'Lesson Summary',
                                        style: TextStyle(
                                          fontWeight: FontWeight.bold,
                                          fontSize: 14,
                                          color: AppColors.textPrimary,
                                        ),
                                      ),
                                      const SizedBox(height: 6),
                                      Text(
                                        'Video lesson for "$lessonTitle". Watch carefully to understand core exam concepts, formulas, and step-by-step solutions.',
                                        style: AppTextStyles.bodySmall.copyWith(
                                          color: AppColors.textSecondary,
                                          height: 1.6,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                              SingleChildScrollView(
                                child: Container(
                                  padding: const EdgeInsets.all(AppSizes.md),
                                  decoration: BoxDecoration(
                                    color: AppColors.bgSecondary,
                                    borderRadius: BorderRadius.circular(AppSizes.radiusLg),
                                    border: Border.all(color: AppColors.border),
                                  ),
                                  child: Column(
                                    children: [
                                      const AiRobotMascot(size: 40),
                                      const SizedBox(height: 10),
                                      const Text(
                                        'Student Community & AI Tutor',
                                        style: TextStyle(fontWeight: FontWeight.bold, color: AppColors.textPrimary),
                                      ),
                                      const SizedBox(height: 4),
                                      Text(
                                        'Ask questions and discuss solutions with fellow Grade 12 students or tap the AI Robot FAB.',
                                        style: AppTextStyles.caption.copyWith(color: AppColors.textMuted),
                                        textAlign: TextAlign.center,
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(height: AppSizes.md),
                      ],
                    ),
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
        color: AppColors.brandEmerald,
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
