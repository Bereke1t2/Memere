import '../../domain/entities/offline_video_entity.dart';
import 'video_model_helpers.dart';

class OfflineVideoModel extends OfflineVideoEntity {
  const OfflineVideoModel({
    required super.videoId,
    required super.lessonId,
    required super.courseId,
    required super.title,
    required super.localPath,
    required super.downloadedAt,
    required super.expiresAt,
    required super.fileSizeBytes,
    required super.status,
  });

  factory OfflineVideoModel.fromJson(Map<String, dynamic> json) {
    return OfflineVideoModel(
      videoId: videoStringValue(json['video_id']),
      lessonId: videoStringValue(json['lesson_id']),
      courseId: videoStringValue(json['course_id']),
      title: videoStringValue(json['title'], fallback: 'Downloaded lesson'),
      localPath: videoStringValue(json['local_path']),
      downloadedAt: videoDateValue(json['downloaded_at']) ?? DateTime.now(),
      expiresAt: videoDateValue(json['expires_at']) ?? DateTime.now(),
      fileSizeBytes: videoIntValue(json['file_size_bytes']),
      status: _parseStatus(videoStringValue(json['status'])),
    );
  }

  factory OfflineVideoModel.fromEntity(OfflineVideoEntity entity) {
    return OfflineVideoModel(
      videoId: entity.videoId,
      lessonId: entity.lessonId,
      courseId: entity.courseId,
      title: entity.title,
      localPath: entity.localPath,
      downloadedAt: entity.downloadedAt,
      expiresAt: entity.expiresAt,
      fileSizeBytes: entity.fileSizeBytes,
      status: entity.status,
    );
  }

  Map<String, dynamic> toJson() => {
        'video_id': videoId,
        'lesson_id': lessonId,
        'course_id': courseId,
        'title': title,
        'local_path': localPath,
        'downloaded_at': downloadedAt.toIso8601String(),
        'expires_at': expiresAt.toIso8601String(),
        'file_size_bytes': fileSizeBytes,
        'status': status.name,
      };

  static OfflineVideoStatus _parseStatus(String value) {
    switch (value.toLowerCase()) {
      case 'queued':
        return OfflineVideoStatus.queued;
      case 'downloading':
        return OfflineVideoStatus.downloading;
      case 'downloaded':
        return OfflineVideoStatus.downloaded;
      case 'expired':
        return OfflineVideoStatus.expired;
      default:
        return OfflineVideoStatus.failed;
    }
  }
}
