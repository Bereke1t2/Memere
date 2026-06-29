String quizStringValue(Object? value, {String fallback = ''}) {
  if (value is String) return value;
  if (value == null) return fallback;
  return value.toString();
}

String? quizNullableString(Object? value) {
  if (value is String && value.trim().isNotEmpty) return value;
  return null;
}

int quizIntValue(Object? value, {int fallback = 0}) {
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value) ?? fallback;
  return fallback;
}

int? quizNullableInt(Object? value) {
  if (value == null) return null;
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value);
  return null;
}

double quizDoubleValue(Object? value, {double fallback = 0}) {
  if (value is int) return value.toDouble();
  if (value is double) return value;
  if (value is num) return value.toDouble();
  if (value is String) return double.tryParse(value) ?? fallback;
  return fallback;
}

bool quizBoolValue(Object? value, {bool fallback = false}) {
  if (value is bool) return value;
  if (value is String) return value.toLowerCase() == 'true';
  if (value is num) return value != 0;
  return fallback;
}

DateTime? quizDateValue(Object? value) {
  if (value is! String || value.isEmpty) return null;
  return DateTime.tryParse(value);
}

List<String> quizStringList(Object? value) {
  if (value is List) {
    return value
        .map((item) => quizStringValue(item))
        .where((item) => item.trim().isNotEmpty)
        .toList();
  }
  return const [];
}
