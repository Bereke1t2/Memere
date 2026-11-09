@Tags(['golden'])
library;

import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/core/constants/app_colors.dart';
import 'package:memere_mobile/shared/widgets/assessment/attempt_header.dart';
import 'package:memere_mobile/shared/widgets/assessment/question_palette_sheet.dart';
import 'package:memere_mobile/shared/widgets/assessment/question_status.dart';
import 'package:memere_mobile/shared/widgets/assessment/submit_summary_sheet.dart';

/// Image previews of the attempt chrome, sized to a 100-question exam because
/// that is the case these were built for.
void main() {
  setUpAll(() async {
    TestWidgetsFlutterBinding.ensureInitialized();
    Future<void> load(String family, List<String> paths) async {
      final found = paths.where((p) => File(p).existsSync()).toList();
      if (found.isEmpty) return;
      final loader = FontLoader(family);
      for (final path in found) {
        loader.addFont(
          File(path).readAsBytes().then((b) => ByteData.view(Uint8List.fromList(b).buffer)),
        );
      }
      await loader.load();
    }

    const substitute = [
      '/usr/share/fonts/noto/NotoSans-Regular.ttf',
      '/usr/share/fonts/noto/NotoSans-Bold.ttf',
    ];
    for (final f in const ['Roboto', 'Sora', 'DM_Sans', 'Noto Sans', 'Noto Sans Ethiopic']) {
      await load(f, substitute);
    }

    // Icons render as boxes without this. The font ships with the SDK but its
    // location varies by install, so every candidate is tried and a miss is
    // simply skipped.
    await load('MaterialIcons', [
      '${Platform.environment['HOME']}/fvm/versions/3.44.4/bin/cache/artifacts/material_fonts/MaterialIcons-Regular.otf',
      '/usr/lib/flutter/bin/cache/artifacts/material_fonts/MaterialIcons-Regular.otf',
      '${Platform.environment['HOME']}/.pub-cache/hosted/pub.dev/provider-6.1.2/extension/devtools/build/assets/fonts/MaterialIcons-Regular.otf',
    ]);
  });

  Future<void> shoot(
    WidgetTester tester,
    String name,
    Widget child, {
    double height = 820,
  }) async {
    await tester.binding.setSurfaceSize(Size(420, height));
    await tester.pumpWidget(
      MaterialApp(
        debugShowCheckedModeBanner: false,
        theme: ThemeData.dark(useMaterial3: true),
        home: Scaffold(backgroundColor: AppColors.bgPrimary, body: child),
      ),
    );
    await tester.pumpAndSettle();
    await expectLater(
      find.byType(MaterialApp),
      matchesGoldenFile('previews/$name.png'),
    );
  }

  /// A paper part-way through, with the blanks scattered rather than bunched
  /// at the end — that is what a real attempt looks like, and it is the case
  /// the grid has to make obvious at a glance.
  List<QuestionStatus> partlyDone(int total) => List.generate(total, (i) {
        final answered = i % 9 != 5;
        final flagged = i % 17 == 3 || i % 23 == 5;
        return QuestionStatus.of(answered: answered, flagged: flagged);
      });

  testWidgets('attempt header', (tester) async {
    await shoot(
      tester,
      'attempt_header',
      Column(
        children: [
          AttemptHeader(
            currentIndex: 13,
            total: 100,
            answered: 84,
            timer: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: AppColors.bgTertiary,
                borderRadius: BorderRadius.circular(999),
                border: Border.all(color: AppColors.borderStrong),
              ),
              child: const Text('1:12:40'),
            ),
            textScaleLabel: '110%',
            onOpenPalette: () {},
            onDecreaseText: () {},
            onIncreaseText: () {},
          ),
        ],
      ),
      height: 220,
    );
  });

  testWidgets('question palette — 100 questions', (tester) async {
    await shoot(
      tester,
      'question_palette',
      QuestionPaletteSheet(
        statuses: partlyDone(100),
        currentIndex: 13,
        onSelected: (_) {},
      ),
      height: 900,
    );
  });

  testWidgets('submit summary — with blanks', (tester) async {
    await shoot(
      tester,
      'submit_summary_blanks',
      SubmitSummarySheet(
        title: 'Submit exam?',
        total: 100,
        answered: 90,
        flagged: 6,
        isTimeUp: false,
        onReviewBlank: () {},
      ),
      height: 560,
    );
  });

  testWidgets('submit summary — complete', (tester) async {
    await shoot(
      tester,
      'submit_summary_complete',
      const SubmitSummarySheet(
        title: 'Submit exam?',
        total: 100,
        answered: 100,
        flagged: 0,
        isTimeUp: false,
      ),
      height: 440,
    );
  });
}
