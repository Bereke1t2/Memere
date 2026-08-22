import '../hive_json.dart';

/// What a bookmark points at. `.name` round-trips as the discriminator.
enum SavedType { course, lesson, quiz, exam }

SavedType savedTypeFromString(String v) {
  switch (v.toLowerCase()) {
    case 'course':
      return SavedType.course;
    case 'lesson':
      return SavedType.lesson;
    case 'quiz':
      return SavedType.quiz;
    case 'exam':
      return SavedType.exam;
    default:
      return SavedType.course;
  }
}

/// A user bookmark held in the `saved_items` box. Purely local (guest bookmarks
/// are never synced to the server — the confirmed "keep local" decision).
class SavedItem {
  const SavedItem({
    required this.id,
    required this.type,
    required this.courseId,
    required this.title,
    required this.savedAt,
    this.subtitle,
  });

  /// The bookmarked content id (courseId / lessonId / quizId / examId).
  final String id;
  final SavedType type;
  final String courseId;
  final String title;
  final String? subtitle;
  final DateTime savedAt;

  /// Unique key within the box; a course and a lesson can share no key.
  String get storageKey => '${type.name}:$id';

  factory SavedItem.fromJson(Map<String, dynamic> json) => SavedItem(
        id: hstr(json['id']),
        type: savedTypeFromString(hstr(json['type'])),
        courseId: hstr(json['course_id']),
        title: hstr(json['title'], fallback: 'Saved item'),
        subtitle: hstrOrNull(json['subtitle']),
        savedAt: hdate(json['saved_at']),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'type': type.name,
        'course_id': courseId,
        'title': title,
        'subtitle': subtitle,
        'saved_at': savedAt.toIso8601String(),
      };
}
