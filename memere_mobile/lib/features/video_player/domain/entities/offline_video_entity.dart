class OfflineVideoEntity {
  const OfflineVideoEntity({
    required this.videoId,
    required this.lessonId,
    required this.courseId,
    required this.title,
    required this.localPath,
    required this.downloadedAt,
    required this.expiresAt,
    required this.fileSizeBytes,
    required this.status,
  });

  final String videoId;
  final String lessonId;
  final String courseId;
  final String title;
  final String localPath;
  final DateTime downloadedAt;
  final DateTime expiresAt;
  final int fileSizeBytes;
  final OfflineVideoStatus status;

  bool get isExpired => DateTime.now().isAfter(expiresAt);
  bool get isPlayable => status == OfflineVideoStatus.downloaded && !isExpired;
}

enum OfflineVideoStatus { queued, downloading, downloaded, failed, expired }
