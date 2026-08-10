import '../constants/env.dart';

/// Sanitizes media and API URLs from backend (e.g. MinIO/S3 localhost or relative file paths)
/// to point to the reachable API host address on mobile devices/emulators.
String fixMediaUrl(String rawUrl) {
  final trimmed = rawUrl.trim();
  if (trimmed.isEmpty) return trimmed;
  try {
    final apiUri = Uri.parse(Env.baseUrl);
    var host = apiUri.host;
    if (host.isEmpty || host == 'localhost' || host == '127.0.0.1') {
      host = '192.168.0.201';
    }

    var url = trimmed;
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      if (url.startsWith('/')) {
        url = 'http://$host:8080$url';
      } else {
        url = 'http://$host:8080/api/v1/$url';
      }
    }

    return url
        .replaceAll('http://localhost:8080', 'http://$host:8080')
        .replaceAll('http://127.0.0.1:8080', 'http://$host:8080')
        .replaceAll('http://localhost:9000', 'http://$host:9000')
        .replaceAll('http://minio:9000', 'http://$host:9000')
        .replaceAll('http://127.0.0.1:9000', 'http://$host:9000');
  } catch (_) {
    return trimmed;
  }
}
