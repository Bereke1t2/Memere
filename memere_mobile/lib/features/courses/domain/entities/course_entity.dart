class CourseEntity {
  const CourseEntity({
    required this.id,
    required this.teacherId,
    required this.title,
    required this.slug,
    required this.description,
    required this.shortDescription,
    required this.subject,
    required this.grade,
    required this.thumbnailUrl,
    required this.price,
    required this.currency,
    required this.isFree,
    required this.isPublished,
    required this.language,
    required this.level,
    required this.totalDurationSeconds,
    required this.totalLessons,
    required this.ratingAvg,
    required this.enrollmentCount,
    required this.createdAt,
    required this.updatedAt,
  });

  final String id;
  final String teacherId;
  final String title;
  final String slug;
  final String description;
  final String shortDescription;
  final String subject;
  final int grade;
  final String? thumbnailUrl;
  final double price;
  final String currency;
  final bool isFree;
  final bool isPublished;
  final String language;
  final CourseLevel level;
  final int totalDurationSeconds;
  final int totalLessons;
  final double ratingAvg;
  final int enrollmentCount;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  bool get isPaid => !isFree && price > 0;

  String get priceLabel {
    if (isFree || price <= 0) return 'Free';
    final amount =
        price % 1 == 0 ? price.toInt().toString() : price.toStringAsFixed(2);
    return '$currency $amount';
  }

  String get durationLabel {
    if (totalDurationSeconds <= 0) return '0m';
    final hours = totalDurationSeconds ~/ 3600;
    final minutes = (totalDurationSeconds % 3600) ~/ 60;
    if (hours > 0 && minutes > 0) return '${hours}h ${minutes}m';
    if (hours > 0) return '${hours}h';
    return '${minutes == 0 ? 1 : minutes}m';
  }

  String get lessonCountLabel {
    if (totalLessons == 1) return '1 lesson';
    return '$totalLessons lessons';
  }

  String get levelLabel {
    switch (level) {
      case CourseLevel.beginner:
        return 'Beginner';
      case CourseLevel.intermediate:
        return 'Intermediate';
      case CourseLevel.advanced:
        return 'Advanced';
    }
  }
}

enum CourseLevel { beginner, intermediate, advanced }
