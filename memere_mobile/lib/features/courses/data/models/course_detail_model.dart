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

  /// Round-trips through [CourseDetailModel.fromJson]: the course fields live at
  /// the top level (as the API returns them) alongside a `sections` array. Only
  /// meaningful for a model built from `fromJson` (remote or the offline cache),
  /// where `course` is a [CourseModel] and each section a [CourseSectionModel].
  Map<String, dynamic> toJson() => {
        ...(course as CourseModel).toJson(),
        'sections': sections
            .whereType<CourseSectionModel>()
            .map((section) => section.toJson())
            .toList(),
      };
}
