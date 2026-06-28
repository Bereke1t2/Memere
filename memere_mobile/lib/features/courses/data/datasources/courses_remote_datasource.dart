import '../../../../core/network/dio_client.dart';
import '../models/course_detail_model.dart';
import '../models/paginated_courses_model.dart';

abstract class CoursesRemoteDataSource {
  Future<PaginatedCoursesModel> listCourses({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  });

  Future<CourseDetailModel> getCourseDetail(String courseId);
}

class CoursesRemoteDataSourceImpl implements CoursesRemoteDataSource {
  const CoursesRemoteDataSourceImpl(this._client);
  final DioClient _client;

  @override
  Future<PaginatedCoursesModel> listCourses({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  }) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/courses',
      queryParameters: {
        'limit': limit,
        if (after != null && after.isNotEmpty) 'after': after,
        if (subject != null && subject.isNotEmpty) 'subject': subject,
        if (grade != null) 'grade': grade,
      },
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing courses response body');
    }
    return PaginatedCoursesModel.fromJson(data);
  }

  @override
  Future<CourseDetailModel> getCourseDetail(String courseId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/courses/$courseId',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing course detail response body');
    }
    return CourseDetailModel.fromJson(data);
  }
}
