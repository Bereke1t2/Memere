class VideoStatusEntity {
  const VideoStatusEntity({
    required this.videoId,
    required this.status,
    required this.durationSeconds,
    required this.hasThumbnail,
    this.error,
  });

  final String videoId;
  final VideoProcessingStatus status;
  final int durationSeconds;
  final bool hasThumbnail;
  final String? error;

  bool get isReady => status == VideoProcessingStatus.ready;
  bool get isProcessing =>
      status == VideoProcessingStatus.pending ||
      status == VideoProcessingStatus.processing;
  bool get isFailed => status == VideoProcessingStatus.failed;
}

enum VideoProcessingStatus { pending, processing, ready, failed }
