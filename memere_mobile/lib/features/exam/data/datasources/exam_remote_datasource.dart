import '../../../../core/network/dio_client.dart';
import '../../../../core/storage/hive/models/offline_exam.dart';
import '../../domain/entities/exam_answer_payload.dart';
import '../models/exam_attempt_analytics_model.dart';
import '../models/exam_attempt_history_model.dart';
import '../models/exam_attempt_model.dart';
import '../models/exam_result_model.dart';
import '../models/mock_exam_model.dart';
import '../models/paginated_mock_exams_model.dart';

abstract class ExamRemoteDataSource {
  Future<List<MockExamModel>> listExamsByCourse(String courseId);

  Future<PaginatedMockExamsModel> listMockExams({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  });

  Future<ExamAttemptModel> startExam(String examId);

  Future<void> saveProgress({
    required String attemptId,
    required ExamAnswerPayload answers,
  });

  Future<ExamResultModel> submitExam({
    required String attemptId,
    required ExamAnswerPayload answers,
  });

  Future<ExamResultModel> getResult(String attemptId);

  Future<ExamAttemptAnalyticsModel> getAnalytics(String attemptId);

  Future<List<ExamAttemptHistoryModel>> listMyAttempts({String? examId});

  /// Fetches the full exam WITH answer keys for offline storage & on-device
  /// grading (`GET /mock-exams/:id/download`; guest-accessible for
  /// published/free content via optionalAuth).
  Future<OfflineExam> getForDownload(String examId);
}

class ExamRemoteDataSourceImpl implements ExamRemoteDataSource {
  const ExamRemoteDataSourceImpl(this._client);
  final DioClient _client;

  @override
  Future<List<MockExamModel>> listExamsByCourse(String courseId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/courses/$courseId/exams',
    );
    final data = response.data;
    if (data == null) return const [];
    final items = data['data'];
    if (items is! List) return const [];
    return items
        .whereType<Map<String, dynamic>>()
        .map(MockExamModel.fromJson)
        .toList();
  }

  @override
  Future<PaginatedMockExamsModel> listMockExams({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  }) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/mock-exams',
      queryParameters: {
        'limit': limit,
        if (after != null && after.isNotEmpty) 'after': after,
        if (subject != null && subject.isNotEmpty) 'subject': subject,
        if (grade != null) 'grade': grade,
      },
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing mock exams response body');
    }
    return PaginatedMockExamsModel.fromJson(data);
  }

  @override
  Future<ExamAttemptModel> startExam(String examId) async {
    final response = await _client.post<Map<String, dynamic>>(
      '/mock-exams/$examId/start',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing exam attempt response body');
    }
    return ExamAttemptModel.fromJson(_unwrapMap(data));
  }

  @override
  Future<void> saveProgress({
    required String attemptId,
    required ExamAnswerPayload answers,
  }) async {
    await _client.patch(
      '/exam-attempts/$attemptId',
      data: {'answers': answers},
    );
  }

  @override
  Future<ExamResultModel> submitExam({
    required String attemptId,
    required ExamAnswerPayload answers,
  }) async {
    final response = await _client.post<Map<String, dynamic>>(
      '/exam-attempts/$attemptId/submit',
      data: {'answers': answers},
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing exam submit response body');
    }
    return ExamResultModel.fromJson(_unwrapMap(data));
  }

  @override
  Future<ExamResultModel> getResult(String attemptId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/exam-attempts/$attemptId/results',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing exam result response body');
    }
    return ExamResultModel.fromJson(_unwrapMap(data));
  }

  @override
  Future<ExamAttemptAnalyticsModel> getAnalytics(String attemptId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/exam-attempts/$attemptId/analytics',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing exam analytics response body');
    }
    return ExamAttemptAnalyticsModel.fromJson(_unwrapMap(data));
  }

  @override
  Future<List<ExamAttemptHistoryModel>> listMyAttempts({String? examId}) async {
    final path = examId != null && examId.isNotEmpty
        ? '/mock-exams/$examId/attempts'
        : '/exam-attempts/my';
    final response = await _client.get<dynamic>(path);
    final data = response.data;
    if (data is List) {
      return data
          .whereType<Map<String, dynamic>>()
          .map(ExamAttemptHistoryModel.fromJson)
          .toList();
    }
    return [];
  }

  @override
  Future<OfflineExam> getForDownload(String examId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/mock-exams/$examId/download',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing exam download response body');
    }
    return OfflineExam.fromJson(_unwrapMap(data));
  }
}

Map<String, dynamic> _unwrapMap(Map<String, dynamic> raw) {
  if (raw.containsKey('data') && raw['data'] is Map<String, dynamic>) {
    return raw['data'] as Map<String, dynamic>;
  }
  return raw;
}
