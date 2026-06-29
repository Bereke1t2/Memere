import '../../domain/entities/video_download_entity.dart';
import 'video_model_helpers.dart';

class VideoDownloadModel extends VideoDownloadEntity {
  const VideoDownloadModel({
    required super.downloadUrl,
    required super.token,
    required super.expiresIn,
  });

  factory VideoDownloadModel.fromJson(Map<String, dynamic> json) {
    return VideoDownloadModel(
      downloadUrl: videoStringValue(json['download_url']),
      token: videoStringValue(json['token']),
      expiresIn: videoIntValue(json['expires_in']),
    );
  }
}
