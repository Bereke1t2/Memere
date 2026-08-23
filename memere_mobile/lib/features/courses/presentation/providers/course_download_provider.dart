import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/storage/secure_pdf_storage.dart';
import '../../../exam/presentation/providers/exam_download_provider.dart';
import '../../../exam/presentation/providers/exam_providers.dart';
import '../../../quiz/presentation/providers/quiz_download_provider.dart';
import '../../../quiz/presentation/providers/quiz_providers.dart';
import '../../../video_player/presentation/providers/offline_video_provider.dart';
import '../../domain/entities/course_detail_entity.dart';

/// Orchestrates a "download the whole course" — every downloadable asset across
/// the course: each lesson's video (MP4) and PDF, plus every quiz and exam
/// (with answer keys, for offline grading). It composes the existing per-asset
/// download entry points rather than reimplementing them, running them
/// sequentially and reporting one aggregate progress value.
final courseDownloadProvider = NotifierProvider.family<CourseDownloadController,
    CourseDownloadProgress, String>(
  CourseDownloadController.new,
);

enum CourseDownloadStatus { idle, running, done, failed }

/// Aggregate progress across every asset queued for a whole-course download.
class CourseDownloadProgress {
  const CourseDownloadProgress({
    this.status = CourseDownloadStatus.idle,
    this.total = 0,
    this.completed = 0,
    this.failed = 0,
    this.currentLabel,
  });

  final CourseDownloadStatus status;

  /// Total number of assets queued.
  final int total;

  /// Assets that downloaded successfully (or had nothing to fetch, e.g. a
  /// notes-only lesson with no real PDF).
  final int completed;

  /// Assets that hit a hard error (network/server) and were skipped.
  final int failed;

  /// The asset currently downloading — drives the "Downloading X…" caption.
  final String? currentLabel;

  int get processed => completed + failed;

  double get fraction => total == 0 ? 0.0 : (processed / total).clamp(0.0, 1.0);

  int get percent => (fraction * 100).round();

  bool get isRunning => status == CourseDownloadStatus.running;

  CourseDownloadProgress copyWith({
    CourseDownloadStatus? status,
    int? total,
    int? completed,
    int? failed,
    String? currentLabel,
  }) {
    return CourseDownloadProgress(
      status: status ?? this.status,
      total: total ?? this.total,
      completed: completed ?? this.completed,
      failed: failed ?? this.failed,
      currentLabel: currentLabel ?? this.currentLabel,
    );
  }
}

enum _TaskResult { success, failed }

/// A single unit of offline work with a human label for the progress caption.
class _DownloadTask {
  _DownloadTask({required this.label, required this.run});

  final String label;
  final Future<_TaskResult> Function() run;
}

class CourseDownloadController
    extends FamilyNotifier<CourseDownloadProgress, String> {
  @override
  CourseDownloadProgress build(String courseId) =>
      const CourseDownloadProgress();

  /// Downloads everything in the course for offline use. Idempotent while
  /// running (a second call is ignored). Individual failures are tolerated and
  /// counted so the whole batch still completes.
  Future<void> downloadCourse(CourseDetailEntity detail) async {
    if (state.isRunning) return;

    final tasks = await _buildTasks(detail);

    if (tasks.isEmpty) {
      state = const CourseDownloadProgress(status: CourseDownloadStatus.done);
      return;
    }

    state = CourseDownloadProgress(
      status: CourseDownloadStatus.running,
      total: tasks.length,
      currentLabel: tasks.first.label,
    );

    var completed = 0;
    var failed = 0;
    for (final task in tasks) {
      state = state.copyWith(currentLabel: task.label);
      _TaskResult result;
      try {
        result = await task.run();
      } catch (_) {
        result = _TaskResult.failed;
      }
      if (result == _TaskResult.failed) {
        failed++;
      } else {
        completed++;
      }
      state = state.copyWith(completed: completed, failed: failed);
    }

    // Refresh the per-video download list so each lesson's own download button
    // reflects what the batch just fetched.
    ref.invalidate(offlineDownloadsProvider);

    state = CourseDownloadProgress(
      status: completed == 0 && failed > 0
          ? CourseDownloadStatus.failed
          : CourseDownloadStatus.done,
      total: tasks.length,
      completed: completed,
      failed: failed,
    );
  }

  Future<List<_DownloadTask>> _buildTasks(CourseDetailEntity detail) async {
    final videoRepo = ref.read(offlineVideoRepositoryProvider);
    final quizNotifier = ref.read(downloadedQuizzesProvider.notifier);
    final examNotifier = ref.read(downloadedExamsProvider.notifier);

    final tasks = <_DownloadTask>[];
    final seenQuizIds = <String>{};

    for (final section in detail.sections) {
      for (final lesson in section.lessons) {
        if (lesson.hasVideo) {
          final videoId = lesson.videoId!;
          tasks.add(_DownloadTask(
            label: lesson.title,
            run: () async {
              final result = await videoRepo.startDownload(
                videoId: videoId,
                lessonId: lesson.id,
                courseId: lesson.courseId,
                title: lesson.title,
              );
              return result.fold(
                (_) => _TaskResult.failed,
                (_) => _TaskResult.success,
              );
            },
          ));
        }

        // Only fetch a PDF when the lesson carries a real pdfUrl; notes-only
        // lessons have nothing to download (the reader renders study notes).
        if (lesson.hasPdf) {
          final url = lesson.pdfUrl!;
          tasks.add(_DownloadTask(
            label: lesson.title,
            run: () async {
              try {
                await SecurePdfStorage.downloadPdf(
                  pdfUrl: url,
                  fileKey: SecurePdfStorage.getFileKey(url, title: lesson.title),
                  lessonId: lesson.id,
                  title: lesson.title,
                  content: lesson.content,
                );
                return _TaskResult.success;
              } on PdfNotAvailableException {
                // Nothing downloadable here — not a failure.
                return _TaskResult.success;
              } catch (_) {
                return _TaskResult.failed;
              }
            },
          ));
        }

        if (lesson.hasQuiz && seenQuizIds.add(lesson.quizId!)) {
          final quizId = lesson.quizId!;
          tasks.add(_DownloadTask(
            label: lesson.title,
            run: () async {
              final failure =
                  await quizNotifier.download(quizId: quizId, title: lesson.title);
              return failure == null ? _TaskResult.success : _TaskResult.failed;
            },
          ));
        }
      }
    }

    // Standalone course quizzes (best-effort; empty on error).
    final quizzes = await ref.read(courseQuizzesProvider(arg).future);
    for (final quiz in quizzes) {
      if (seenQuizIds.add(quiz.id)) {
        tasks.add(_DownloadTask(
          label: quiz.title,
          run: () async {
            final failure =
                await quizNotifier.download(quizId: quiz.id, title: quiz.title);
            return failure == null ? _TaskResult.success : _TaskResult.failed;
          },
        ));
      }
    }

    // Course exams.
    final exams = await ref.read(courseExamsProvider(arg).future);
    for (final exam in exams) {
      tasks.add(_DownloadTask(
        label: exam.title,
        run: () async {
          final failure =
              await examNotifier.download(examId: exam.id, title: exam.title);
          return failure == null ? _TaskResult.success : _TaskResult.failed;
        },
      ));
    }

    return tasks;
  }
}
