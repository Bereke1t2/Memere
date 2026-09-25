import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/features/courses/presentation/screens/pdf_reader_screen.dart';

void main() {
  Widget buildTestApp({
    required Size screenSize,
    required String title,
    required String pdfUrl,
    String? content,
  }) {
    return ProviderScope(
      child: MaterialApp(
        home: MediaQuery(
          data: MediaQueryData(size: screenSize),
          child: SizedBox(
            width: screenSize.width,
            height: screenSize.height,
            child: PdfReaderScreen(
              title: title,
              pdfUrl: pdfUrl,
              content: content,
            ),
          ),
        ),
      ),
    );
  }

  group('HTML Presentation View Screen Sizes Test', () {
    const testSizes = [
      Size(320, 568),  // Compact Mobile (iPhone SE)
      Size(390, 844),  // Standard Mobile (iPhone 13/14)
      Size(430, 932),  // Large Mobile (iPhone 15 Pro Max)
      Size(768, 1024), // Tablet Portrait
      Size(1024, 768), // Tablet Landscape
    ];

    for (final size in testSizes) {
      testWidgets('renders cleanly on screen size ${size.width}x${size.height}', (tester) async {
        tester.view.physicalSize = size;
        tester.view.devicePixelRatio = 1.0;
        addTearDown(tester.view.resetPhysicalSize);

        await tester.pumpWidget(
          buildTestApp(
            screenSize: size,
            title: 'Cell Biology Presentation',
            pdfUrl: 'cell_biology.html',
            content: '<h1>Cell Structure</h1><p>Notes on organelles and membranes.</p>',
          ),
        );

        await tester.pump();

        // Check header title
        expect(find.text('Cell Biology Presentation'), findsOneWidget);

        // Check Mode Switcher has HTML Document tab
        expect(find.text('HTML Document'), findsOneWidget);
        expect(find.text('Study Notes'), findsOneWidget);

        // Verify no render overflow exceptions occurred
        expect(tester.takeException(), isNull);
      });
    }

    testWidgets('switching between Study Notes and HTML Document works smoothly', (tester) async {
      tester.view.physicalSize = const Size(390, 844);
      tester.view.devicePixelRatio = 1.0;
      addTearDown(tester.view.resetPhysicalSize);

      await tester.pumpWidget(
        buildTestApp(
          screenSize: const Size(390, 844),
          title: 'Physics Wave Optics',
          pdfUrl: 'optics_presentation.html',
          content: '## Wave Optics Summary\nKey formulas and diffraction notes.',
        ),
      );

      await tester.pump();

      // Initially active tab is HTML Document (since pdfUrl is provided and not empty/sample)
      expect(find.text('HTML Document'), findsOneWidget);

      // Tap on Study Notes tab
      await tester.tap(find.text('Study Notes'));
      await tester.pumpAndSettle();

      // Study notes sub-toolbar and title should be visible
      expect(find.text('Font Size'), findsOneWidget);
      expect(find.text('Copy Notes'), findsOneWidget);
      expect(find.text('Physics Wave Optics'), findsWidgets);

      // Tap back on HTML Document tab
      await tester.tap(find.text('HTML Document'));
      await tester.pump();

      expect(tester.takeException(), isNull);
    });
  });
}
