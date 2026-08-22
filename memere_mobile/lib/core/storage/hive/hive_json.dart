/// Small, defensive JSON coercion helpers shared by the offline Hive models.
///
/// The offline stores ingest two shapes of JSON: the backend `/download` wire
/// payload (on first save) and the model's own `toJson` (on reload from Hive).
/// These helpers tolerate both — and any partially-corrupt entry — without
/// throwing, mirroring the `video_model_helpers` idiom used by offline videos.
library;

String hstr(dynamic v, {String fallback = ''}) {
  if (v is String) return v;
  return v?.toString() ?? fallback;
}

String? hstrOrNull(dynamic v) {
  if (v == null) return null;
  if (v is String) return v.isEmpty ? null : v;
  return v.toString();
}

int hint(dynamic v, {int fallback = 0}) {
  if (v is int) return v;
  if (v is num) return v.toInt();
  return int.tryParse('$v') ?? fallback;
}

int? hintOrNull(dynamic v) {
  if (v == null) return null;
  if (v is int) return v;
  if (v is num) return v.toInt();
  return int.tryParse('$v');
}

double hdouble(dynamic v, {double fallback = 0}) {
  if (v is double) return v;
  if (v is num) return v.toDouble();
  return double.tryParse('$v') ?? fallback;
}

bool hbool(dynamic v) {
  if (v is bool) return v;
  if (v is num) return v != 0;
  final s = '$v'.toLowerCase();
  return s == 'true' || s == '1';
}

DateTime hdate(dynamic v) {
  if (v is String) return DateTime.tryParse(v) ?? DateTime.now();
  return DateTime.now();
}

DateTime? hdateOrNull(dynamic v) {
  if (v is String) return DateTime.tryParse(v);
  return null;
}

/// Decodes a JSON list of maps into a typed list, skipping malformed entries.
List<T> hlist<T>(dynamic v, T Function(Map<String, dynamic>) fromJson) {
  if (v is! List) return <T>[];
  final out = <T>[];
  for (final e in v) {
    if (e is Map) {
      out.add(fromJson(Map<String, dynamic>.from(e)));
    }
  }
  return out;
}
