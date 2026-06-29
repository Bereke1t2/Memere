import 'offline_video_entity.dart';

class DownloadTaskEntity {
  const DownloadTaskEntity({
    required this.videoId,
    required this.lessonId,
    required this.courseId,
    required this.title,
    required this.progress,
    required this.status,
    this.errorMessage,
  });

  final String videoId;
  final String lessonId;
  final String courseId;
  final String title;
  final double progress;
  final OfflineVideoStatus status;
  final String? errorMessage;
}
