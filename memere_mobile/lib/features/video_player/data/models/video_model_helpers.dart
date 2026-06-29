String videoStringValue(Object? value, {String fallback = ''}) {
  if (value is String) return value;
  if (value == null) return fallback;
  return value.toString();
}

String? videoNullableString(Object? value) {
  if (value is String && value.trim().isNotEmpty) return value;
  return null;
}

int videoIntValue(Object? value, {int fallback = 0}) {
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value) ?? fallback;
  return fallback;
}

double videoDoubleValue(Object? value, {double fallback = 0}) {
  if (value is int) return value.toDouble();
  if (value is double) return value;
  if (value is num) return value.toDouble();
  if (value is String) return double.tryParse(value) ?? fallback;
  return fallback;
}

bool videoBoolValue(Object? value, {bool fallback = false}) {
  if (value is bool) return value;
  if (value is String) return value.toLowerCase() == 'true';
  if (value is num) return value != 0;
  return fallback;
}

DateTime? videoDateValue(Object? value) {
  if (value is! String || value.isEmpty) return null;
  return DateTime.tryParse(value);
}
