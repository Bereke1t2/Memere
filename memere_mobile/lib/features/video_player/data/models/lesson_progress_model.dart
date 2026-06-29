import '../../domain/entities/lesson_progress_entity.dart';
import 'video_model_helpers.dart';

class LessonProgressModel extends LessonProgressEntity {
  const LessonProgressModel({
    required super.lessonId,
    required super.isCompleted,
    super.completedAt,
    required super.videoProgressSeconds,
  });

  factory LessonProgressModel.fromJson(Map<String, dynamic> json) {
    return LessonProgressModel(
      lessonId: videoStringValue(json['lesson_id']),
      isCompleted: videoBoolValue(json['is_completed']),
      completedAt: videoDateValue(json['completed_at']),
      videoProgressSeconds: videoIntValue(json['video_progress_seconds']),
    );
  }
}
