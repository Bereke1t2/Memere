import '../../domain/entities/video_status_entity.dart';
import 'video_model_helpers.dart';

class VideoStatusModel extends VideoStatusEntity {
  const VideoStatusModel({
    required super.videoId,
    required super.status,
    required super.durationSeconds,
    required super.hasThumbnail,
    super.error,
  });

  factory VideoStatusModel.fromJson(Map<String, dynamic> json) {
    return VideoStatusModel(
      videoId: videoStringValue(json['video_id']),
      status: _parseStatus(videoStringValue(json['status'])),
      durationSeconds: videoIntValue(json['duration_seconds']),
      hasThumbnail: videoBoolValue(json['has_thumbnail']),
      error: videoNullableString(json['error']),
    );
  }

  static VideoProcessingStatus _parseStatus(String value) {
    switch (value.toLowerCase()) {
      case 'processing':
        return VideoProcessingStatus.processing;
      case 'ready':
        return VideoProcessingStatus.ready;
      case 'failed':
        return VideoProcessingStatus.failed;
      default:
        return VideoProcessingStatus.pending;
    }
  }
}
