import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../domain/entities/mock_exam_entity.dart';
import '../../domain/usecases/list_mock_exams_usecase.dart';
import 'exam_providers.dart';

const mockExamSubjects = [
  'Mathematics',
  'Physics',
  'Chemistry',
  'Biology',
  'English',
  'History',
  'Geography',
  'Economics',
];

class MockExamCatalogState {
  const MockExamCatalogState({
    this.exams = const [],
    this.filteredExams = const [],
    this.nextCursor,
    this.selectedSubject,
    this.selectedGrade = 12,
    this.searchQuery = '',
    this.isLoadingMore = false,
    this.loadMoreError,
  });

  final List<MockExamEntity> exams;
  final List<MockExamEntity> filteredExams;
  final String? nextCursor;
  final String? selectedSubject;
  final int? selectedGrade;
  final String searchQuery;
  final bool isLoadingMore;
  final Failure? loadMoreError;

  bool get hasMore => nextCursor != null && nextCursor!.isNotEmpty;
  bool get hasActiveFilters =>
      searchQuery.trim().isNotEmpty ||
      selectedSubject != null ||
      selectedGrade != 12;

  MockExamCatalogState copyWith({
    List<MockExamEntity>? exams,
    List<MockExamEntity>? filteredExams,
    String? nextCursor,
    bool clearNextCursor = false,
    String? selectedSubject,
    bool clearSelectedSubject = false,
    int? selectedGrade,
    bool clearSelectedGrade = false,
    String? searchQuery,
    bool? isLoadingMore,
    Failure? loadMoreError,
    bool clearLoadMoreError = false,
  }) {
    return MockExamCatalogState(
      exams: exams ?? this.exams,
      filteredExams: filteredExams ?? this.filteredExams,
      nextCursor: clearNextCursor ? null : nextCursor ?? this.nextCursor,
      selectedSubject:
          clearSelectedSubject ? null : selectedSubject ?? this.selectedSubject,
      selectedGrade:
          clearSelectedGrade ? null : selectedGrade ?? this.selectedGrade,
      searchQuery: searchQuery ?? this.searchQuery,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      loadMoreError: clearLoadMoreError ? null : loadMoreError,
    );
  }
}

final mockExamCatalogProvider =
    AsyncNotifierProvider<MockExamCatalogNotifier, MockExamCatalogState>(
  MockExamCatalogNotifier.new,
);

class MockExamCatalogNotifier extends AsyncNotifier<MockExamCatalogState> {
  @override
  Future<MockExamCatalogState> build() async {
    return _fetchExams(
      subject: null,
      grade: 12,
      searchQuery: '',
    );
  }

  Future<void> refresh() async {
    final previous = state.valueOrNull ?? const MockExamCatalogState();
    state = const AsyncLoading<MockExamCatalogState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchExams(
        subject: previous.selectedSubject,
        grade: previous.selectedGrade,
        searchQuery: previous.searchQuery,
      ),
    );
  }

  Future<void> loadMore() async {
    final current = state.valueOrNull;
    if (current == null || current.isLoadingMore || !current.hasMore) return;

    state = AsyncData(
      current.copyWith(isLoadingMore: true, clearLoadMoreError: true),
    );

    final useCase = ref.read(listMockExamsUseCaseProvider);
    final result = await useCase(
      ListMockExamsParams(
        limit: AppConstants.defaultPageLimit,
        after: current.nextCursor,
        subject: current.selectedSubject,
        grade: current.selectedGrade,
      ),
    );

    result.fold(
      (failure) {
        state = AsyncData(
          current.copyWith(isLoadingMore: false, loadMoreError: failure),
        );
      },
      (page) {
        final merged = _mergeExams(current.exams, page.exams);
        state = AsyncData(
          current.copyWith(
            exams: merged,
            filteredExams: _applySearch(merged, current.searchQuery),
            nextCursor: page.nextCursor,
            clearNextCursor: page.nextCursor == null,
            isLoadingMore: false,
            clearLoadMoreError: true,
          ),
        );
      },
    );
  }

  void setSearchQuery(String value) {
    final current = state.valueOrNull;
    if (current == null) return;
    state = AsyncData(
      current.copyWith(
        searchQuery: value,
        filteredExams: _applySearch(current.exams, value),
      ),
    );
  }

  Future<void> setSubject(String? subject) async {
    final current = state.valueOrNull ?? const MockExamCatalogState();
    state = const AsyncLoading<MockExamCatalogState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchExams(
        subject: subject,
        grade: current.selectedGrade,
        searchQuery: current.searchQuery,
      ),
    );
  }

  Future<void> setGrade(int? grade) async {
    final current = state.valueOrNull ?? const MockExamCatalogState();
    state = const AsyncLoading<MockExamCatalogState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchExams(
        subject: current.selectedSubject,
        grade: grade,
        searchQuery: current.searchQuery,
      ),
    );
  }

  Future<void> clearFilters() async {
    state = const AsyncLoading<MockExamCatalogState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchExams(
        subject: null,
        grade: 12,
        searchQuery: '',
      ),
    );
  }

  Future<MockExamCatalogState> _fetchExams({
    required String? subject,
    required int? grade,
    required String searchQuery,
  }) async {
    final useCase = ref.read(listMockExamsUseCaseProvider);
    final result = await useCase(
      ListMockExamsParams(
        limit: AppConstants.defaultPageLimit,
        subject: subject,
        grade: grade,
      ),
    );

    return result.fold(
      (failure) => throw failure,
      (page) {
        return MockExamCatalogState(
          exams: page.exams,
          filteredExams: _applySearch(page.exams, searchQuery),
          nextCursor: page.nextCursor,
          selectedSubject: subject,
          selectedGrade: grade,
          searchQuery: searchQuery,
        );
      },
    );
  }

  List<MockExamEntity> _mergeExams(
    List<MockExamEntity> existing,
    List<MockExamEntity> incoming,
  ) {
    final seen = existing.map((exam) => exam.id).toSet();
    return [
      ...existing,
      ...incoming.where((exam) => seen.add(exam.id)),
    ];
  }

  List<MockExamEntity> _applySearch(
    List<MockExamEntity> exams,
    String query,
  ) {
    final normalized = query.trim().toLowerCase();
    if (normalized.isEmpty) return exams;

    return exams.where((exam) {
      return exam.title.toLowerCase().contains(normalized) ||
          exam.subject.toLowerCase().contains(normalized) ||
          (exam.instructions ?? '').toLowerCase().contains(normalized);
    }).toList();
  }
}
