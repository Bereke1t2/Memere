import 'course_entity.dart';
import 'course_section_entity.dart';

class CourseDetailEntity {
  const CourseDetailEntity({
    required this.course,
    required this.sections,
  });

  final CourseEntity course;
  final List<CourseSectionEntity> sections;

  int get lessonCount => sections.fold<int>(
        0,
        (total, section) => total + section.lessons.length,
      );
}
