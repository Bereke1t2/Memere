import '../../domain/entities/course_detail_entity.dart';
import 'course_model.dart';
import 'course_section_model.dart';

class CourseDetailModel extends CourseDetailEntity {
  const CourseDetailModel({
    required super.course,
    required super.sections,
  });

  factory CourseDetailModel.fromJson(Map<String, dynamic> json) {
    final sectionsJson = json['sections'];
    final sections = sectionsJson is List
        ? sectionsJson
            .whereType<Map<String, dynamic>>()
            .map(CourseSectionModel.fromJson)
            .toList()
        : <CourseSectionModel>[];

    return CourseDetailModel(
      course: CourseModel.fromJson(json),
      sections: sections,
    );
  }
}
