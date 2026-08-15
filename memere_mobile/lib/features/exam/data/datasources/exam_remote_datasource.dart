import '../../../../core/network/dio_client.dart';
import '../../domain/entities/exam_answer_payload.dart';
import '../models/exam_attempt_analytics_model.dart';
import '../models/exam_attempt_history_model.dart';
import '../models/exam_attempt_model.dart';
import '../models/exam_result_model.dart';
import '../models/paginated_mock_exams_model.dart';

abstract class ExamRemoteDataSource {
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
}

class ExamRemoteDataSourceImpl implements ExamRemoteDataSource {
  const ExamRemoteDataSourceImpl(this._client);
  final DioClient _client;

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
    return ExamAttemptModel.fromJson(data);
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
    return ExamResultModel.fromJson(data);
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
    return ExamResultModel.fromJson(data);
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
    return ExamAttemptAnalyticsModel.fromJson(data);
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
}
