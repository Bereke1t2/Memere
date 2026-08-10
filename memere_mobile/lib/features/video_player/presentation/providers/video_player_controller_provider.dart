import 'dart:async';

import 'package:chewie/chewie.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:video_player/video_player.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/utils/media_url_helper.dart';
import '../../domain/entities/video_status_entity.dart';
import '../../domain/entities/video_stream_entity.dart';
import '../../domain/usecases/save_video_progress_usecase.dart';
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

    if (arg.videoId.trim().isEmpty) {
      return const VideoPlaybackState(
        errorMessage: 'This lesson does not have a video attached yet.',
      );
    }
    if (arg.lessonId.trim().isEmpty) {
      return const VideoPlaybackState(
        errorMessage: 'Lesson information is missing.',
      );
    }

    final statusResult =
        await ref.read(getVideoStatusUseCaseProvider)(arg.videoId);
    final status =
        statusResult.fold((failure) => throw failure, (value) => value);

    if (!status.isReady) {
      return VideoPlaybackState(status: status);
    }

    final streamResult =
        await ref.read(getVideoStreamUseCaseProvider)(arg.videoId);
    final stream =
        streamResult.fold((failure) => throw failure, (value) => value);

    try {
      final resolvedUrl = fixMediaUrl(stream.masterUrl);
      _videoController =
          VideoPlayerController.networkUrl(Uri.parse(resolvedUrl));
      await _videoController!.initialize();
      _videoController!.addListener(_handlePlaybackTick);

      _chewieController = ChewieController(
        videoPlayerController: _videoController!,
        autoPlay: true,
        allowFullScreen: true,
        allowMuting: true,
        aspectRatio: 16 / 9,
        showControlsOnInitialize: true,
        cupertinoProgressColors: ChewieProgressColors(
          playedColor: const Color(0xFFFF5252),
          handleColor: const Color(0xFFFF5252),
          bufferedColor: Colors.white24,
          backgroundColor: Colors.white10,
        ),
        materialProgressColors: ChewieProgressColors(
          playedColor: const Color(0xFFFF5252),
          handleColor: const Color(0xFFFF5252),
          bufferedColor: Colors.white24,
          backgroundColor: Colors.white10,
        ),
      );

      return VideoPlaybackState(
        status: status,
        stream: stream,
        videoController: _videoController,
        chewieController: _chewieController,
      );
    } catch (_) {
      return VideoPlaybackState(
        status: status,
        stream: stream,
        errorMessage:
            'Video material is being updated by your instructor. Please try again shortly.',
      );
    }
  }

  Future<void> retry() async {
    _disposeControllers();
    state = const AsyncLoading();
    state = await AsyncValue.guard(() => build(arg));
  }

  Future<void> refreshStream() => retry();

  Future<void> saveProgressNow() => _saveProgress(force: true);

  Future<void> markComplete() async {
    if (arg.lessonId.trim().isEmpty) return;
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
