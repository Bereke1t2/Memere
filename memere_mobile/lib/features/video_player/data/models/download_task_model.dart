import '../../domain/entities/download_task_entity.dart';
import '../../domain/entities/offline_video_entity.dart';

class DownloadTaskModel extends DownloadTaskEntity {
  const DownloadTaskModel({
    required super.videoId,
    required super.lessonId,
    required super.courseId,
    required super.title,
    required super.progress,
    required super.status,
    super.errorMessage,
  });

  DownloadTaskModel copyWith({
    double? progress,
    OfflineVideoStatus? status,
    String? errorMessage,
    bool clearError = false,
  }) {
    return DownloadTaskModel(
      videoId: videoId,
      lessonId: lessonId,
      courseId: courseId,
      title: title,
      progress: progress ?? this.progress,
      status: status ?? this.status,
      errorMessage: clearError ? null : errorMessage ?? this.errorMessage,
    );
  }
}
