import '../../../../core/network/dio_client.dart';
import '../models/student_points_model.dart';

abstract class ProgressRemoteDataSource {
  Future<StudentPointsModel> getMyPoints();
}

class ProgressRemoteDataSourceImpl implements ProgressRemoteDataSource {
  const ProgressRemoteDataSourceImpl(this._client);
  final DioClient _client;

  @override
  Future<StudentPointsModel> getMyPoints() async {
    final response = await _client.get<Map<String, dynamic>>('/me/points');
    return StudentPointsModel.fromJson(_requireBody(response.data));
  }

  Map<String, dynamic> _requireBody(Map<String, dynamic>? data) {
    if (data == null) {
      throw const FormatException('Missing points response body');
    }
    return data;
  }
}
