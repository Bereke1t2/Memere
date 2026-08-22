import 'package:fpdart/fpdart.dart';

import '../../../../core/errors/failures.dart';
import '../entities/exam_answer_payload.dart';
import '../entities/exam_attempt_analytics_entity.dart';
import '../entities/exam_attempt_entity.dart';
import '../entities/exam_attempt_history_entity.dart';
import '../entities/exam_result_entity.dart';
import '../entities/mock_exam_entity.dart';
import '../entities/paginated_mock_exams_entity.dart';

abstract class ExamRepository {
  Future<Either<Failure, List<MockExamEntity>>> listExamsByCourse(String courseId);

  Future<Either<Failure, PaginatedMockExamsEntity>> listMockExams({
    int limit = 20,
    String? after,
    String? subject,
    int? grade,
  });

  Future<Either<Failure, ExamAttemptEntity>> startExam(String examId);

  Future<Either<Failure, void>> saveProgress({
    required String attemptId,
    required ExamAnswerPayload answers,
  });

  Future<Either<Failure, ExamResultEntity>> submitExam({
    required String attemptId,
    required ExamAnswerPayload answers,
  });

  Future<Either<Failure, ExamResultEntity>> getResult(String attemptId);

  Future<Either<Failure, ExamAttemptAnalyticsEntity>> getAnalytics(
    String attemptId,
  );

  Future<Either<Failure, List<ExamAttemptHistoryEntity>>> listMyAttempts({
    String? examId,
  });
}
