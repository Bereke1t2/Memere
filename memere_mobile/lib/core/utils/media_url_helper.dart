import '../constants/env.dart';

/// Sanitizes media and API URLs from backend (e.g. MinIO/S3 localhost or relative file paths)
/// to point to the reachable API host address on mobile devices/emulators.
String fixMediaUrl(String rawUrl) {
  final trimmed = rawUrl.trim();
  if (trimmed.isEmpty) return trimmed;
  try {
    final apiUri = Uri.parse(Env.baseUrl);
    final scheme = apiUri.scheme.isNotEmpty ? apiUri.scheme : 'http';
    var host = apiUri.host;
    if (host.isEmpty) {
      host = '10.0.2.2';
    }

    final portPart = (apiUri.hasPort && apiUri.port != 0)
        ? ':${apiUri.port}'
        : (scheme == 'https' ? '' : ':8080');
    final origin = '$scheme://$host$portPart';

    var url = trimmed;
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      if (url.startsWith('/uploads/') ||
          url.startsWith('uploads/') ||
          url.startsWith('/storage/') ||
          url.startsWith('storage/')) {
        final cleanPath = url.startsWith('/') ? url : '/$url';
        url = '$origin$cleanPath';
      } else if (url.startsWith('/')) {
        url = '$origin$url';
      } else {
        url = '$origin/api/v1/$url';
      }
    }

    return url
        .replaceAll('http://localhost:8080', origin)
        .replaceAll('http://127.0.0.1:8080', origin)
        .replaceAll('http://localhost:9000', origin)
        .replaceAll('http://minio:9000', origin)
        .replaceAll('http://127.0.0.1:9000', origin);
  } catch (_) {
    return trimmed;
  }
}

