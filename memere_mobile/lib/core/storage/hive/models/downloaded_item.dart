import '../hive_json.dart';

/// The kind of asset a downloaded item represents. `.name` round-trips as the
/// stored/queryable discriminator.
enum DownloadType { video, pdf, quiz, exam }

/// Lifecycle of a single downloaded asset.
enum DownloadItemStatus { queued, downloading, downloaded, failed, expired }

DownloadType downloadTypeFromString(String v) {
  switch (v.toLowerCase()) {
    case 'video':
      return DownloadType.video;
    case 'pdf':
      return DownloadType.pdf;
    case 'quiz':
      return DownloadType.quiz;
    case 'exam':
      return DownloadType.exam;
    default:
      return DownloadType.pdf;
  }
}

DownloadItemStatus downloadItemStatusFromString(String v) {
  switch (v.toLowerCase()) {
    case 'queued':
      return DownloadItemStatus.queued;
    case 'downloading':
      return DownloadItemStatus.downloading;
    case 'downloaded':
      return DownloadItemStatus.downloaded;
    case 'expired':
      return DownloadItemStatus.expired;
    default:
      return DownloadItemStatus.failed;
  }
}

/// One entry in the unified `downloads_index` box — the single source that
/// enumerates every downloaded asset (video, PDF, quiz, exam) across courses.
/// The heavy payloads live elsewhere (video files, `SecurePdfStorage`, the
/// encrypted quiz/exam boxes); this row is the lightweight, queryable manifest.
class DownloadedItem {
  const DownloadedItem({
    required this.id,
    required this.type,
    required this.courseId,
    required this.title,
    required this.sizeBytes,
    required this.status,
    required this.downloadedAt,
    this.sectionId,
    this.lessonId,
    this.subtitle,
    this.localPath,
    this.expiresAt,
  });

  /// The content id (videoId / lessonId for a PDF / quizId / examId).
  final String id;
  final DownloadType type;
  final String courseId;
  final String? sectionId;
  final String? lessonId;
  final String title;
  final String? subtitle;
  final int sizeBytes;

  /// File path for file-backed assets (video/PDF); null for quiz/exam JSON.
  final String? localPath;
  final DownloadItemStatus status;
  final DateTime downloadedAt;
  final DateTime? expiresAt;

  /// Unique key within the index box; disambiguates same-id assets of different
  /// kinds (e.g. a lesson's PDF vs. its quiz).
  String get storageKey => '${type.name}:$id';

  bool get isDownloaded => status == DownloadItemStatus.downloaded;

  DownloadedItem copyWith({
    int? sizeBytes,
    String? localPath,
    DownloadItemStatus? status,
    DateTime? expiresAt,
  }) =>
      DownloadedItem(
        id: id,
        type: type,
        courseId: courseId,
        sectionId: sectionId,
        lessonId: lessonId,
        title: title,
        subtitle: subtitle,
        sizeBytes: sizeBytes ?? this.sizeBytes,
        localPath: localPath ?? this.localPath,
        status: status ?? this.status,
        downloadedAt: downloadedAt,
        expiresAt: expiresAt ?? this.expiresAt,
      );

  factory DownloadedItem.fromJson(Map<String, dynamic> json) => DownloadedItem(
        id: hstr(json['id']),
        type: downloadTypeFromString(hstr(json['type'])),
        courseId: hstr(json['course_id']),
        sectionId: hstrOrNull(json['section_id']),
        lessonId: hstrOrNull(json['lesson_id']),
        title: hstr(json['title'], fallback: 'Downloaded item'),
        subtitle: hstrOrNull(json['subtitle']),
        sizeBytes: hint(json['size_bytes']),
        localPath: hstrOrNull(json['local_path']),
        status: downloadItemStatusFromString(hstr(json['status'])),
        downloadedAt: hdate(json['downloaded_at']),
        expiresAt: hdateOrNull(json['expires_at']),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'type': type.name,
        'course_id': courseId,
        'section_id': sectionId,
        'lesson_id': lessonId,
        'title': title,
        'subtitle': subtitle,
        'size_bytes': sizeBytes,
        'local_path': localPath,
        'status': status.name,
        'downloaded_at': downloadedAt.toIso8601String(),
        'expires_at': expiresAt?.toIso8601String(),
      };
}
