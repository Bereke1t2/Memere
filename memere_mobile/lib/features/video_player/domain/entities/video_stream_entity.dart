class VideoStreamEntity {
  const VideoStreamEntity({
    required this.masterUrl,
    required this.expiresIn,
    this.thumbnailUrl,
    required this.durationSeconds,
  });

  final String masterUrl;
  final int expiresIn;
  final String? thumbnailUrl;
  final int durationSeconds;

  DateTime get expiresAt => DateTime.now().add(Duration(seconds: expiresIn));
}
