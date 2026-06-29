import '../../../../core/network/dio_client.dart';
import '../../domain/entities/quiz_answer_payload.dart';
import '../models/quiz_attempt_model.dart';
import '../models/quiz_model.dart';
import '../models/quiz_result_model.dart';

abstract class QuizRemoteDataSource {
  Future<List<QuizModel>> listQuizzesByCourse(String courseId);
  Future<QuizModel> getQuiz(String quizId);
  Future<QuizAttemptModel> startAttempt(String quizId);
  Future<void> saveProgress({
    required String attemptId,
    required QuizAnswerPayload answers,
  });
  Future<QuizResultModel> submitAttempt({
    required String attemptId,
    required QuizAnswerPayload answers,
  });
  Future<QuizResultModel> getResult(String attemptId);
}

class QuizRemoteDataSourceImpl implements QuizRemoteDataSource {
  const QuizRemoteDataSourceImpl(this._client);
  final DioClient _client;

  @override
  Future<List<QuizModel>> listQuizzesByCourse(String courseId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/courses/$courseId/quizzes',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing course quizzes response body');
    }
    final items = data['data'];
    if (items is! List) return const [];
    return items
        .whereType<Map<String, dynamic>>()
        .map(QuizModel.fromJson)
        .toList();
  }

  @override
  Future<QuizModel> getQuiz(String quizId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/quizzes/$quizId',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing quiz response body');
    }
    return QuizModel.fromJson(data);
  }

  @override
  Future<QuizAttemptModel> startAttempt(String quizId) async {
    final response = await _client.post<Map<String, dynamic>>(
      '/quizzes/$quizId/attempts',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing quiz attempt response body');
    }
    return QuizAttemptModel.fromJson(data);
  }

  @override
  Future<void> saveProgress({
    required String attemptId,
    required QuizAnswerPayload answers,
  }) async {
    await _client.patch(
      '/quiz-attempts/$attemptId',
      data: {'answers': answers},
    );
  }

  @override
  Future<QuizResultModel> submitAttempt({
    required String attemptId,
    required QuizAnswerPayload answers,
  }) async {
    final response = await _client.post<Map<String, dynamic>>(
      '/quiz-attempts/$attemptId/submit',
      data: {'answers': answers},
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing quiz submit response body');
    }
    return QuizResultModel.fromJson(data);
  }

  @override
  Future<QuizResultModel> getResult(String attemptId) async {
    final response = await _client.get<Map<String, dynamic>>(
      '/quiz-attempts/$attemptId/result',
    );
    final data = response.data;
    if (data == null) {
      throw const FormatException('Missing quiz result response body');
    }
    return QuizResultModel.fromJson(data);
  }
}
