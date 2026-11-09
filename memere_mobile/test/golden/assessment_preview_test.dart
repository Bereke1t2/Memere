import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/core/constants/app_colors.dart';
import 'package:memere_mobile/features/quiz/domain/entities/question_feedback_entity.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_answer_entity.dart';
import 'package:memere_mobile/features/quiz/domain/entities/quiz_question_entity.dart';
import 'package:memere_mobile/features/quiz/presentation/widgets/question_feedback_tile.dart';
import 'package:memere_mobile/features/quiz/presentation/widgets/quiz_question_card.dart';

/// Renders the assessment widgets to PNG so the result can be looked at, not
/// just compiled. These are previews rather than regression goldens: run with
/// `flutter test --update-goldens test/golden` and open the files under
/// `test/golden/previews/`.
///
/// The fixtures are real ESSLCE question-bank content — LaTeX, a long
/// multi-line prompt, a worked derivation — because those are the cases that
/// actually break layout.
void main() {
  // Widget tests ship a dummy font that draws every glyph as a filled box,
  // which makes a preview useless for judging type or spacing. Loading a real
  // face and the KaTeX faces flutter_math needs renders actual text and actual
  // formulas. Both are read from disk, so a machine without them still runs the
  // test — it just falls back to boxes.
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

    // The app's real faces (Sora, DM Sans) are fetched at runtime and are not
    // on disk here, so a stand-in is registered under each family the styles
    // ask for. These previews are therefore accurate for layout, spacing and
    // colour, but not for the exact typeface.
    const substitute = [
      '/usr/share/fonts/noto/NotoSans-Regular.ttf',
      '/usr/share/fonts/noto/NotoSans-Bold.ttf',
      '/usr/share/fonts/liberation/LiberationSans-Regular.ttf',
    ];
    for (final family in const [
      'Roboto',
      'Sora',
      'DM Sans',
      'Noto Sans',
      'Noto Sans Ethiopic',
    ]) {
      await load(family, substitute);
    }

    // flutter_math asks for its glyphs under the package-prefixed family name
    // (see make_symbol.dart), which is what a dependency's declared fonts are
    // registered as. Loading them under the bare name silently does nothing.
    const katexDir = String.fromEnvironment(
      'KATEX_DIR',
      defaultValue: '',
    );
    final dir = katexDir.isNotEmpty
        ? katexDir
        : '${Platform.environment['HOME']}/.pub-cache/hosted/pub.dev/'
            'flutter_math_fork-0.7.4/lib/katex_fonts/fonts';
    for (final family in const [
      'KaTeX_Main',
      'KaTeX_Math',
      'KaTeX_Size1',
      'KaTeX_Size2',
      'KaTeX_Size3',
      'KaTeX_Size4',
      'KaTeX_AMS',
      'KaTeX_Caligraphic',
      'KaTeX_Fraktur',
      'KaTeX_SansSerif',
      'KaTeX_Script',
      'KaTeX_Typewriter',
    ]) {
      await load('packages/flutter_math_fork/$family', [
        '$dir/$family-Regular.ttf',
        '$dir/$family-Bold.ttf',
        '$dir/$family-Italic.ttf',
        '$dir/$family-BoldItalic.ttf',
      ]);
    }
  });

  Future<void> shoot(
    WidgetTester tester,
    String name,
    Widget child, {
    double width = 420,
    double height = 900,
  }) async {
    await tester.binding.setSurfaceSize(Size(width, height));
    await tester.pumpWidget(
      MaterialApp(
        debugShowCheckedModeBanner: false,
        theme: ThemeData.dark(useMaterial3: true),
        home: Scaffold(
          backgroundColor: AppColors.bgPrimary,
          body: SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: child,
          ),
        ),
      ),
    );
    await tester.pumpAndSettle();
    await expectLater(
      find.byType(MaterialApp),
      matchesGoldenFile('previews/$name.png'),
    );
  }

  QuizAnswerEntity option(String id, String text, int order) =>
      QuizAnswerEntity(id: id, text: text, orderIndex: order);

  QuizAnswerFeedbackEntity reviewed(String id, String text, {bool isCorrect = false}) =>
      QuizAnswerFeedbackEntity(id: id, text: text, isCorrect: isCorrect);

  testWidgets('question card — maths prompt with LaTeX options', (tester) async {
    final question = QuizQuestionEntity(
      id: 'q1',
      text: r'Let $f(x) = -\frac{1}{3}\cos\left(\frac{3}{2}x - \frac{\pi}{3}\right) + 1$. '
          r'What are the period and phase shift of the graph of $f$ respectively?',
      type: QuizQuestionType.multipleChoice,
      points: 5,
      subject: 'Mathematics',
      answers: [
        option('a', r'$\frac{4\pi}{3}$ and $\frac{2\pi}{9}$', 0),
        option('b', r'$\frac{2\pi}{3}$ and $\frac{2\pi}{9}$', 1),
        option('c', r'$\frac{2\pi}{3}$ and $\frac{4\pi}{9}$', 2),
        option('d', r'$\frac{\pi}{9}$ and $\frac{4\pi}{9}$', 3),
      ],
    );

    await shoot(
      tester,
      'question_card_maths',
      QuizQuestionCard(
        question: question,
        selectedAnswer: const ['a'],
        onSingleAnswerSelected: (_, __) {},
        onMultiAnswerToggled: (_, __) {},
        onShortAnswerChanged: (_, __) {},
      ),
    );
  });

  testWidgets('question card — long multi-line prompt', (tester) async {
    final question = QuizQuestionEntity(
      id: 'q2',
      text: 'Three students designed the following experimental settings to '
          'understand the actual taste of a food.\n'
          'Student 1: Blind folded two students, one of them tightly holding the '
          'nose while the other without holding the nose and tasting the food.\n'
          'Student 2: Blind folded one student and the other left with open eyes, '
          'both holding their nose tightly while tasting the food.\n'
          'Student 3: None of the participants blind folded and only one of them '
          'holding her nose tightly while tasting the food.\n'
          'Which student designed the most appropriate experiment?',
      type: QuizQuestionType.multipleChoice,
      points: 1,
      subject: 'Biology',
      answers: [
        option('a', 'Student 1', 0),
        option('b', 'Student 2', 1),
        option('c', 'Student 3', 2),
        option('d', 'None of them', 3),
      ],
    );

    await shoot(
      tester,
      'question_card_long_prompt',
      QuizQuestionCard(
        question: question,
        selectedAnswer: null,
        onSingleAnswerSelected: (_, __) {},
        onMultiAnswerToggled: (_, __) {},
        onShortAnswerChanged: (_, __) {},
      ),
      height: 1000,
    );
  });

  testWidgets('feedback tile — correct, with a worked derivation', (tester) async {
    final feedback = QuestionFeedbackEntity(
      questionId: 'q1',
      correct: true,
      pointsAwarded: 5,
      pointsPossible: 5,
      questionText: r'What is the period of $f(x) = \cos\left(\frac{3}{2}x\right)$?',
      explanation: 'Step-by-step derivation:\n'
          r'The general form is $f(x) = A \cos(Bx - C) + D$, with $B = \frac{3}{2}$.'
          '\n'
          r'$$T = \frac{2\pi}{B} = \frac{2\pi}{3/2} = \frac{4\pi}{3}$$'
          '\n\nSource: Q38 [2015 E.C. | Booklet 12]',
      selectedAnswers: const ['a'],
      correctAnswerIds: const ['a'],
      answers: [
        reviewed('a', r'$\frac{4\pi}{3}$', isCorrect: true),
        reviewed('b', r'$\frac{2\pi}{3}$'),
        reviewed('c', r'$\frac{\pi}{3}$'),
      ],
    );

    await shoot(tester, 'feedback_correct', QuestionFeedbackTile(feedback: feedback, index: 0));
  });

  testWidgets('feedback tile — incorrect', (tester) async {
    final feedback = QuestionFeedbackEntity(
      questionId: 'q2',
      correct: false,
      pointsAwarded: 0,
      pointsPossible: 1,
      questionText: 'Which organelle is the main site of ATP production?',
      explanation: 'Antonyms: Mitochondria carry out oxidative phosphorylation, '
          "which produces most of the cell's ATP.",
      selectedAnswers: const ['c'],
      correctAnswerIds: const ['b'],
      answers: [
        reviewed('a', 'Nucleus'),
        reviewed('b', 'Mitochondrion', isCorrect: true),
        reviewed('c', 'Ribosome'),
        reviewed('d', 'Golgi apparatus'),
      ],
    );

    await shoot(tester, 'feedback_incorrect', QuestionFeedbackTile(feedback: feedback, index: 1));
  });
}
