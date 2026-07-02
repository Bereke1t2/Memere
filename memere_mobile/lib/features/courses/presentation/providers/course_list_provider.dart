import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../domain/entities/course_entity.dart';
import '../../domain/usecases/list_courses_usecase.dart';
import 'courses_providers.dart';

const phase2Subjects = [
  'Mathematics',
  'Physics',
  'Chemistry',
  'Biology',
  'English',
  'Civics',
  'History',
  'Geography',
  'Economics',
];

class CourseListState {
  const CourseListState({
    this.courses = const [],
    this.filteredCourses = const [],
    this.nextCursor,
    this.selectedSubject,
    this.selectedGrade = 12,
    this.searchQuery = '',
    this.isLoadingMore = false,
    this.loadMoreError,
  });

  final List<CourseEntity> courses;
  final List<CourseEntity> filteredCourses;
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

  CourseListState copyWith({
    List<CourseEntity>? courses,
    List<CourseEntity>? filteredCourses,
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
    return CourseListState(
      courses: courses ?? this.courses,
      filteredCourses: filteredCourses ?? this.filteredCourses,
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

final courseListProvider =
    AsyncNotifierProvider<CourseListNotifier, CourseListState>(
  CourseListNotifier.new,
);

class CourseListNotifier extends AsyncNotifier<CourseListState> {
  @override
  Future<CourseListState> build() async {
    return _fetchCourses(
      subject: null,
      grade: 12,
      searchQuery: '',
    );
  }

  Future<void> refresh() async {
    final previous = state.valueOrNull ?? const CourseListState();
    state = const AsyncLoading<CourseListState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchCourses(
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

    final useCase = ref.read(listCoursesUseCaseProvider);
    final result = await useCase(
      ListCoursesParams(
        limit: AppConstants.defaultPageLimit,
        after: current.nextCursor,
        subject: current.selectedSubject,
        grade: current.selectedGrade,
      ),
    );

    result.fold(
      (failure) {
        state = AsyncData(
          current.copyWith(
            isLoadingMore: false,
            loadMoreError: failure,
          ),
        );
      },
      (page) {
        final merged = _mergeCourses(current.courses, page.courses);
        state = AsyncData(
          current.copyWith(
            courses: merged,
            filteredCourses: _applySearch(merged, current.searchQuery),
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
        filteredCourses: _applySearch(current.courses, value),
      ),
    );
  }

  Future<void> setSubject(String? subject) async {
    final current = state.valueOrNull ?? const CourseListState();
    state = const AsyncLoading<CourseListState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchCourses(
        subject: subject,
        grade: current.selectedGrade,
        searchQuery: current.searchQuery,
      ),
    );
  }

  Future<void> setGrade(int? grade) async {
    final current = state.valueOrNull ?? const CourseListState();
    state = const AsyncLoading<CourseListState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchCourses(
        subject: current.selectedSubject,
        grade: grade,
        searchQuery: current.searchQuery,
      ),
    );
  }

  Future<void> clearFilters() async {
    state = const AsyncLoading<CourseListState>().copyWithPrevious(state);
    state = await AsyncValue.guard(
      () => _fetchCourses(
        subject: null,
        grade: 12,
        searchQuery: '',
      ),
    );
  }

  Future<CourseListState> _fetchCourses({
    required String? subject,
    required int? grade,
    required String searchQuery,
  }) async {
    final useCase = ref.read(listCoursesUseCaseProvider);
    final result = await useCase(
      ListCoursesParams(
        limit: AppConstants.defaultPageLimit,
        subject: subject,
        grade: grade,
        searchQuery: searchQuery,
      ),
    );

    return result.fold(
      (failure) => throw failure,
      (page) {
        return CourseListState(
          courses: page.courses,
          filteredCourses: _applySearch(page.courses, searchQuery),
          nextCursor: page.nextCursor,
          selectedSubject: subject,
          selectedGrade: grade,
          searchQuery: searchQuery,
        );
      },
    );
  }

  List<CourseEntity> _mergeCourses(
    List<CourseEntity> existing,
    List<CourseEntity> incoming,
  ) {
    final seen = existing.map((course) => course.id).toSet();
    return [
      ...existing,
      ...incoming.where((course) => seen.add(course.id)),
    ];
  }

  List<CourseEntity> _applySearch(
    List<CourseEntity> courses,
    String query,
  ) {
    final normalized = query.trim().toLowerCase();
    if (normalized.isEmpty) return courses;

    return courses.where((course) {
      return course.title.toLowerCase().contains(normalized) ||
          course.subject.toLowerCase().contains(normalized) ||
          course.shortDescription.toLowerCase().contains(normalized) ||
          course.description.toLowerCase().contains(normalized);
    }).toList();
  }
}
