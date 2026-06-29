class VideoDownloadEntity {
  const VideoDownloadEntity({
    required this.downloadUrl,
    required this.token,
    required this.expiresIn,
  });

  final String downloadUrl;
  final String token;
  final int expiresIn;

  DateTime get expiresAt => DateTime.now().add(Duration(seconds: expiresIn));
}
