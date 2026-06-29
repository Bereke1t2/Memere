import '../../../../core/network/dio_client.dart';
import '../models/lesson_progress_model.dart';
import '../models/video_download_model.dart';
import '../models/video_status_model.dart';
import '../models/video_stream_model.dart';

abstract class VideoRemoteDataSource {
  Future<VideoStatusModel> getVideoStatus(String videoId);
  Future<VideoStreamModel> getStream(String videoId);
  Future<VideoDownloadModel> getDownloadUrl(String videoId);
  Future<LessonProgressModel> saveVideoProgress({
    required String lessonId,
    required int positionSeconds,
  });
  Future<void> markLessonComplete(String lessonId);
}

class VideoRemoteDataSourceImpl implements VideoRemoteDataSource {
  const VideoRemoteDataSourceImpl(this._client);
  final DioClient _client;

  @override
  Future<VideoStatusModel> getVideoStatus(String videoId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/videos/$videoId/status',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing video status response body');
    }
    return VideoStatusModel.fromJson(data);
  }

  @override
  Future<VideoStreamModel> getStream(String videoId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/videos/$videoId/stream',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing video stream response body');
    }
    return VideoStreamModel.fromJson(data);
  }

  @override
  Future<VideoDownloadModel> getDownloadUrl(String videoId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/videos/$videoId/download-url',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing video download response body');
    }
    return VideoDownloadModel.fromJson(data);
  }

  @override
  Future<LessonProgressModel> saveVideoProgress({
    required String lessonId,
    required int positionSeconds,
  }) async {
    final response = await _client.put<Map<String, dynamic>>(
      '/lessons/$lessonId/video-progress',
      data: {'position_seconds': positionSeconds},
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing lesson progress response body');
    }
    return LessonProgressModel.fromJson(data);
  }

  @override
  Future<void> markLessonComplete(String lessonId) async {
    await _client.post('/lessons/$lessonId/complete');
  }
}
