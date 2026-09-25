import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/core/constants/app_sizes.dart';
import 'package:memere_mobile/features/courses/domain/entities/course_entity.dart';
import 'package:memere_mobile/features/courses/domain/entities/paginated_courses_entity.dart';
import 'package:memere_mobile/features/courses/presentation/providers/course_list_provider.dart';
import 'package:memere_mobile/features/courses/presentation/screens/course_list_screen.dart';
import 'package:memere_mobile/features/courses/presentation/widgets/course_card.dart';

class _FakeCourseListNotifier extends CourseListNotifier {
  _FakeCourseListNotifier(this._initialState);
  final CourseListState _initialState;

  @override
  Future<CourseListState> build() async => _initialState;
}

PaginatedCoursesEntity _createMockCourses(int count) {
  return PaginatedCoursesEntity(
    courses: List.generate(
      count,
      (i) => CourseEntity(
        id: 'course-$i',
        teacherId: 'teacher-1',
        title: 'Grade 12 Physics Course #$i',
        slug: 'grade-12-physics-$i',
        description: 'Comprehensive entrance exam preparation for unit $i.',
        shortDescription: 'Unit $i exam preparation',
        subject: 'Physics',
        grade: 12,
        thumbnailUrl: null,
        price: 0,
        currency: 'ETB',
        isFree: true,
        isPublished: true,
        language: 'English',
        level: CourseLevel.beginner,
        totalDurationSeconds: 3600,
        totalLessons: 10,
        ratingAvg: 4.8,
        enrollmentCount: 150,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      ),
    ),
    nextCursor: null,
    limit: count,
  );
}

void main() {
  Widget buildHomeApp({
    required Size screenSize,
    required EdgeInsets padding,
    required PaginatedCoursesEntity mockData,
  }) {
    return ProviderScope(
      overrides: [
        courseListProvider.overrideWith(
          () => _FakeCourseListNotifier(
            CourseListState(
              courses: mockData.courses,
              filteredCourses: mockData.courses,
            ),
          ),
        ),
      ],
      child: MaterialApp(
        home: MediaQuery(
          data: MediaQueryData(
            size: screenSize,
            padding: padding,
          ),
          child: Scaffold(
            extendBody: true,
            body: const CourseListScreen(),
            bottomNavigationBar: SafeArea(
              top: false,
              minimum: const EdgeInsets.fromLTRB(
                AppSizes.md,
                0,
                AppSizes.md,
                AppSizes.md,
              ),
              child: Container(
                height: AppSizes.bottomNavHeight,
                color: Colors.blue,
              ),
            ),
          ),
        ),
      ),
    );
  }

  group('Home Screen Bottom Navigation Clearance Tests', () {
    const testConfigs = [
      (
        'Android 3-Button Navigation (360x800, padding.bottom = 0)',
        Size(360, 800),
        EdgeInsets.only(top: 24, bottom: 0),
      ),
      (
        'Gesture Navigation / iPhone (390x844, padding.bottom = 34)',
        Size(390, 844),
        EdgeInsets.only(top: 48, bottom: 34),
      ),
      (
        'Tablet / Large Screen (768x1024, padding.bottom = 20)',
        Size(768, 1024),
        EdgeInsets.only(top: 24, bottom: 20),
      ),
    ];

    for (final (description, size, padding) in testConfigs) {
      testWidgets(description, (tester) async {
        tester.view.physicalSize = size;
        tester.view.devicePixelRatio = 1.0;
        addTearDown(tester.view.resetPhysicalSize);

        final mockData = _createMockCourses(6);

        await tester.pumpWidget(
          buildHomeApp(
            screenSize: size,
            padding: padding,
            mockData: mockData,
          ),
        );

        await tester.pump();
        await tester.pump(const Duration(milliseconds: 500));

        // Scroll progressively to ensure all lazy sliver items are built and bottom reached
        final scrollable = find.byType(Scrollable).first;
        for (int i = 0; i < 8; i++) {
          await tester.drag(scrollable, const Offset(0, -500));
          await tester.pump(const Duration(milliseconds: 100));
        }
        await tester.pump(const Duration(milliseconds: 300));

        // Find the last course card
        final lastCardFinder = find.byType(CourseRowCard).last;
        expect(lastCardFinder, findsOneWidget);

        // Get the bottom coordinate of the last course card
        final cardBottomY = tester.getBottomRight(lastCardFinder).dy;

        // The floating bottom navigation bar top edge
        final navBarFinder = find.byType(SafeArea).last;
        final navBarTopY = tester.getTopLeft(navBarFinder).dy;

        // Ensure the last course card bottom is ABOVE the top of the bottom navigation bar
        expect(
          cardBottomY,
          lessThanOrEqualTo(navBarTopY),
          reason: 'The last course card should be fully visible above the bottom navigation bar.',
        );

        // Verify no rendering overflow exceptions occurred
        expect(tester.takeException(), isNull);
      });
    }
  });
}
