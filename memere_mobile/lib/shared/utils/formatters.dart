String formatDurationSeconds(int seconds) {
  if (seconds <= 0) return '0m';
  final hours = seconds ~/ 3600;
  final minutes = (seconds % 3600) ~/ 60;

  if (hours > 0 && minutes > 0) return '${hours}h ${minutes}m';
  if (hours > 0) return '${hours}h';
  return '${minutes == 0 ? 1 : minutes}m';
}

String formatCompactCount(int count) {
  if (count < 1000) return count.toString();
  if (count < 1000000) {
    final value = count / 1000;
    final text =
        value >= 10 ? value.toStringAsFixed(0) : value.toStringAsFixed(1);
    return '${_trimTrailingZero(text)}k';
  }
  final value = count / 1000000;
  final text =
      value >= 10 ? value.toStringAsFixed(0) : value.toStringAsFixed(1);
  return '${_trimTrailingZero(text)}m';
}

String formatPrice(double price, String currency) {
  if (price <= 0) return 'Free';
  final amount =
      price % 1 == 0 ? price.toInt().toString() : price.toStringAsFixed(2);
  return '$currency $amount';
}

String _trimTrailingZero(String value) {
  return value.endsWith('.0') ? value.substring(0, value.length - 2) : value;
}
