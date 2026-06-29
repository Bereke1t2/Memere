import '../../domain/entities/video_stream_entity.dart';
import 'video_model_helpers.dart';

class VideoStreamModel extends VideoStreamEntity {
  const VideoStreamModel({
    required super.masterUrl,
    required super.expiresIn,
    super.thumbnailUrl,
    required super.durationSeconds,
  });

  factory VideoStreamModel.fromJson(Map<String, dynamic> json) {
    return VideoStreamModel(
      masterUrl: videoStringValue(json['master_url']),
      expiresIn: videoIntValue(json['expires_in']),
      thumbnailUrl: videoNullableString(json['thumbnail_url']),
      durationSeconds: videoIntValue(json['duration_seconds']),
    );
  }
}
