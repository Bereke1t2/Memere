import 'dart:async';
import 'dart:io';

import 'package:chewie/chewie.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:video_player/video_player.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/utils/media_url_helper.dart';
import '../../domain/entities/video_status_entity.dart';
import '../../domain/entities/video_stream_entity.dart';
import '../../domain/usecases/save_video_progress_usecase.dart';
import '../../../auth/presentation/providers/auth_state_provider.dart';
import 'offline_video_provider.dart';
import 'video_providers.dart';

class VideoPlayerParams {
  const VideoPlayerParams({
    required this.videoId,
    required this.lessonId,
    required this.courseId,
  });

  final String videoId;
  final String lessonId;
  final String courseId;

  @override
  bool operator ==(Object other) {
    return other is VideoPlayerParams &&
        other.videoId == videoId &&
        other.lessonId == lessonId &&
        other.courseId == courseId;
  }

  @override
  int get hashCode => Object.hash(videoId, lessonId, courseId);
}

class VideoPlaybackState {
  const VideoPlaybackState({
    this.status,
    this.stream,
    this.videoController,
    this.chewieController,
    this.isInitializing = false,
    this.isSavingProgress = false,
    this.lastSavedPositionSeconds = 0,
    this.isCompleted = false,
    this.isOffline = false,
    this.errorMessage,
  });

  final VideoStatusEntity? status;
  final VideoStreamEntity? stream;
  final VideoPlayerController? videoController;
  final ChewieController? chewieController;
  final bool isInitializing;
  final bool isSavingProgress;
  final int lastSavedPositionSeconds;
  final bool isCompleted;
  final bool isOffline;
  final String? errorMessage;

  bool get isReady =>
      videoController != null &&
      chewieController != null &&
      (videoController?.value.isInitialized ?? false);

  VideoPlaybackState copyWith({
    VideoStatusEntity? status,
    VideoStreamEntity? stream,
    VideoPlayerController? videoController,
    ChewieController? chewieController,
    bool? isInitializing,
    bool? isSavingProgress,
    int? lastSavedPositionSeconds,
    bool? isCompleted,
    bool? isOffline,
    String? errorMessage,
    bool clearError = false,
  }) {
    return VideoPlaybackState(
      status: status ?? this.status,
      stream: stream ?? this.stream,
      videoController: videoController ?? this.videoController,
      chewieController: chewieController ?? this.chewieController,
      isInitializing: isInitializing ?? this.isInitializing,
      isSavingProgress: isSavingProgress ?? this.isSavingProgress,
      lastSavedPositionSeconds:
          lastSavedPositionSeconds ?? this.lastSavedPositionSeconds,
      isCompleted: isCompleted ?? this.isCompleted,
      isOffline: isOffline ?? this.isOffline,
      errorMessage: clearError ? null : errorMessage ?? this.errorMessage,
    );
  }
}

final videoPlayerControllerProvider = AsyncNotifierProviderFamily<
    VideoPlayerControllerNotifier,
    VideoPlaybackState,
    VideoPlayerParams>(VideoPlayerControllerNotifier.new);

class VideoPlayerControllerNotifier
    extends FamilyAsyncNotifier<VideoPlaybackState, VideoPlayerParams> {
  VideoPlayerController? _videoController;
  ChewieController? _chewieController;
  bool _wasPlaying = false;

  @override
  Future<VideoPlaybackState> build(VideoPlayerParams arg) async {
    ref.onDispose(() {
      unawaited(_saveProgress(force: true));
      _disposeControllers();
    });

    if (arg.videoId.trim().isEmpty && arg.lessonId.trim().isEmpty) {
      return const VideoPlaybackState(
        errorMessage: 'Lesson video information is missing.',
      );
    }

    // 1. Check if video has been downloaded for OFFLINE PLAYBACK
    try {
      final offlineResult =
          await ref.read(offlineVideoRepositoryProvider).getDownload(arg.videoId);
      final offlineVideo = offlineResult.fold((_) => null, (v) => v);
      if (offlineVideo != null && offlineVideo.localPath.isNotEmpty) {
        final file = File(offlineVideo.localPath);
        if (await file.exists() && (await file.length()) > 1000) {
          _videoController = VideoPlayerController.file(file);
          await _videoController!.initialize();
          _videoController!.addListener(_handlePlaybackTick);
          _chewieController = _buildChewieController(_videoController!, isOffline: true);

          return VideoPlaybackState(
            videoController: _videoController,
            chewieController: _chewieController,
            isOffline: true,
          );
        }
      }
    } catch (_) {}

    // 2. LIVE STREAMING: Fetch status & stream from server
    VideoStatusEntity? status;
    try {
      final statusResult =
          await ref.read(getVideoStatusUseCaseProvider)(arg.videoId);
      status = statusResult.fold((_) => null, (value) => value);
    } catch (_) {}

    VideoStreamEntity? stream;
    try {
      final streamResult =
          await ref.read(getVideoStreamUseCaseProvider)(arg.videoId);
      stream = streamResult.fold((_) => null, (value) => value);
    } catch (_) {}

    // 3. Resolve the signed stream URL. The backend signs a single self-contained
    //    MP4 that plays progressively over HTTP Range. We deliberately do NOT fall
    //    back to any sample clip: a video with no playable URL must surface as a
    //    real error (matching the offline-download behaviour), never silently play
    //    unrelated footage — that fake fallback is exactly what produced the
    //    misleading "check your internet" message when its host was unreachable.
    final streamUrl = (stream != null && stream.masterUrl.trim().isNotEmpty)
        ? fixMediaUrl(stream.masterUrl)
        : '';

    if (streamUrl.isEmpty) {
      return VideoPlaybackState(
        status: status,
        stream: stream,
        errorMessage: _streamUnavailableMessage(status),
      );
    }

    // 4. Initialize the network player against the signed MP4.
    try {
      _videoController = VideoPlayerController.networkUrl(Uri.parse(streamUrl));
      await _videoController!.initialize();
      _videoController!.addListener(_handlePlaybackTick);

      _chewieController = _buildChewieController(_videoController!, isOffline: false);

      return VideoPlaybackState(
        status: status,
        stream: stream,
        videoController: _videoController,
        chewieController: _chewieController,
        isOffline: false,
      );
    } catch (_) {
      // Report the real failure (no sample-clip fallback) so the user sees that
      // THIS video failed and can retry, and so logs point at the real problem.
      return VideoPlaybackState(
        status: status,
        stream: stream,
        errorMessage: _streamUnavailableMessage(status),
      );
    }
  }

  /// Human message for a stream that won't load. If the server says the video is
  /// still processing (or failed), say so; otherwise it's a load/network error.
  String _streamUnavailableMessage(VideoStatusEntity? status) {
    if (status != null && status.isProcessing) {
      return 'This video is still being processed. Please check back shortly.';
    }
    if (status != null && status.isFailed) {
      return 'This video could not be processed. Please contact support.';
    }
    return 'Unable to load this video. Please check your connection and retry.';
  }

  ChewieController _buildChewieController(
    VideoPlayerController controller, {
    required bool isOffline,
  }) {
    return ChewieController(
      videoPlayerController: controller,
      autoPlay: true,
      allowFullScreen: true,
      allowMuting: true,
      allowPlaybackSpeedChanging: true,
      playbackSpeeds: const [0.75, 1.0, 1.25, 1.5, 2.0],
      aspectRatio: 16 / 9,
      showControlsOnInitialize: true,
      cupertinoProgressColors: ChewieProgressColors(
        playedColor: AppColors.brandEmerald,
        handleColor: AppColors.brandEmerald,
        bufferedColor: Colors.white24,
        backgroundColor: Colors.white10,
      ),
      materialProgressColors: ChewieProgressColors(
        playedColor: AppColors.brandEmerald,
        handleColor: AppColors.brandEmerald,
        bufferedColor: Colors.white24,
        backgroundColor: Colors.white10,
      ),
    );
  }

  Future<void> retry() async {
    _disposeControllers();
    state = const AsyncLoading();
    state = await AsyncValue.guard(() => build(arg));
  }

  Future<void> refreshStream() => retry();

  Future<void> saveProgressNow() => _saveProgress(force: true);

  /// Guests have no server-side progress; writes would 401. Skip the network
  /// call and keep completion state on-device instead.
  bool get _isSignedOut =>
      !(ref.read(authStateProvider).valueOrNull?.isAuthenticated ?? false);

  Future<void> markComplete() async {
    if (arg.lessonId.trim().isEmpty) return;
    if (_isSignedOut) {
      state = AsyncData(
        (state.valueOrNull ?? const VideoPlaybackState()).copyWith(
          isCompleted: true,
          clearError: true,
        ),
      );
      return;
    }
    final result = await ref.read(markLessonCompleteUseCaseProvider)(
      arg.lessonId,
    );
    result.fold(
      (failure) => state = AsyncData(
        (state.valueOrNull ?? const VideoPlaybackState()).copyWith(
          errorMessage: failure.message,
        ),
      ),
      (_) => state = AsyncData(
        (state.valueOrNull ?? const VideoPlaybackState()).copyWith(
          isCompleted: true,
          clearError: true,
        ),
      ),
    );
  }

  void _handlePlaybackTick() {
    final controller = _videoController;
    if (controller == null || !controller.value.isInitialized) return;

    final isPlaying = controller.value.isPlaying;
    final currentSeconds = controller.value.position.inSeconds;
    final currentState = state.valueOrNull;
    final lastSavedSeconds = currentState?.lastSavedPositionSeconds ?? 0;

    if (_wasPlaying && !isPlaying && currentSeconds != lastSavedSeconds) {
      unawaited(_saveProgress(force: true));
    }

    _wasPlaying = isPlaying;
    if (!isPlaying) return;

    if (currentSeconds - lastSavedSeconds >=
        AppConstants.videoProgressSaveIntervalSeconds) {
      unawaited(_saveProgress());
    }
  }

  Future<void> _saveProgress({bool force = false}) async {
    if (_isSignedOut) return;
    final controller = _videoController;
    final currentState = state.valueOrNull;
    if (controller == null ||
        !controller.value.isInitialized ||
        currentState == null ||
        currentState.isSavingProgress ||
        arg.lessonId.trim().isEmpty) {
      return;
    }

    final positionSeconds = controller.value.position.inSeconds;
    if (!force &&
        positionSeconds - currentState.lastSavedPositionSeconds <
            AppConstants.videoProgressSaveIntervalSeconds) {
      return;
    }
    if (positionSeconds == currentState.lastSavedPositionSeconds) return;

    state = AsyncData(currentState.copyWith(isSavingProgress: true));
    final result = await ref.read(saveVideoProgressUseCaseProvider)(
      SaveVideoProgressParams(
        lessonId: arg.lessonId,
        positionSeconds: positionSeconds,
      ),
    );

    final latestState = state.valueOrNull ?? currentState;
    result.fold(
      (failure) {
        state = AsyncData(
          latestState.copyWith(
            isSavingProgress: false,
            errorMessage: failure.message,
          ),
        );
      },
      (progress) {
        state = AsyncData(
          latestState.copyWith(
            isSavingProgress: false,
            lastSavedPositionSeconds: progress.videoProgressSeconds,
            isCompleted: progress.isCompleted,
            clearError: true,
          ),
        );
      },
    );
  }

  void _disposeControllers() {
    _videoController?.removeListener(_handlePlaybackTick);
    _chewieController?.dispose();
    _videoController?.dispose();
    _chewieController = null;
    _videoController = null;
  }
}

String videoFailureMessage(Object error) {
  if (error is Failure) return error.message;
  return 'Could not load video. Please try again.';
}
