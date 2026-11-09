import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/shared/widgets/math_text.dart';

void main() {
  Widget host(String content, {bool selectable = false}) => MaterialApp(
        home: Scaffold(
          body: MathText(
            content,
            style: const TextStyle(fontSize: 14, color: Colors.white),
            selectable: selectable,
          ),
        ),
      );

  group('normalizeContent', () {
    test('turns <br> and closing block tags into line breaks', () {
      expect(
        MathText.normalizeContent('one<br>two<br/>three</p>four'),
        'one\ntwo\nthree\nfour',
      );
    });

    test('strips leftover tags without touching LaTeX', () {
      // Backslashes, braces and $ are content here, not markup.
      expect(
        MathText.normalizeContent(r'<b>Area</b> is $\frac{1}{2}bh$'),
        r'Area is $\frac{1}{2}bh$',
      );
    });

    test('decodes entities, ampersand last so double-encoding survives', () {
      expect(MathText.normalizeContent('a &lt; b &amp;&amp; c'), 'a < b && c');
      expect(MathText.normalizeContent('&amp;lt;'), '&lt;');
    });

    test('collapses runs of blank lines and trims', () {
      expect(MathText.normalizeContent('  a\n\n\n\nb  '), 'a\n\nb');
    });

    test('leaves ordinary prose untouched', () {
      const plain = 'Which organelle makes most ATP?';
      expect(MathText.normalizeContent(plain), plain);
    });
  });

  group('looksRich', () {
    test('detects every delimiter the widget handles', () {
      expect(MathText.looksRich(r'$x$'), isTrue);
      expect(MathText.looksRich(r'\(x\)'), isTrue);
      expect(MathText.looksRich(r'\[x\]'), isTrue);
      expect(MathText.looksRich('**bold**'), isTrue);
    });

    test('is false for prose, so plain questions skip rich layout', () {
      expect(MathText.looksRich('Bacterial cell wall is unique due to its'), isFalse);
    });
  });

  group('parseFragments — delimiter precedence and boundaries', () {
    List<ContentFragmentKind> kinds(String s) =>
        MathText.parseFragments(s).map((f) => f.kind).toList();

    test('inline maths keeps the surrounding sentence intact', () {
      final f = MathText.parseFragments(r'Area of $\triangle ABC$ in units');
      expect(kinds(r'Area of $\triangle ABC$ in units'), [
        ContentFragmentKind.text,
        ContentFragmentKind.inlineMath,
        ContentFragmentKind.text,
      ]);
      expect(f[0].value, 'Area of ');
      expect(f[1].value, r'\triangle ABC');
      expect(f[2].value, ' in units');
    });

    test(r'display $$..$$ is matched before the inline rule can split it', () {
      // A naive $..$ rule would claim the opening $$ and mangle the formula.
      final f = MathText.parseFragments(r'$$x+y$$');
      expect(f, hasLength(1));
      expect(f.single.kind, ContentFragmentKind.displayMath);
      expect(f.single.value, 'x+y');
    });

    test(r'\(..\) is inline and \[..\] is display', () {
      final f = MathText.parseFragments(r'a \(x^2\) b \[y^2\] c');
      expect(
        f.map((e) => e.kind),
        [
          ContentFragmentKind.text,
          ContentFragmentKind.inlineMath,
          ContentFragmentKind.text,
          ContentFragmentKind.displayMath,
          ContentFragmentKind.text,
        ],
      );
      expect(f[1].value, 'x^2');
      expect(f[3].value, 'y^2');
    });

    test('several formulas in one sentence are all found', () {
      expect(
        MathText.parseFragments(r'If $a=1$ and $b=2$ then $a+b=3$')
            .where((f) => f.kind == ContentFragmentKind.inlineMath)
            .map((f) => f.value),
        [r'a=1', r'b=2', r'a+b=3'],
      );
    });

    test('a lone unmatched \$ stays prose and cannot swallow a paragraph', () {
      // Inline maths stops at a newline, which bounds the damage of a stray $.
      expect(kinds('The price is \$5\nNext line'), [ContentFragmentKind.text]);
    });

    test('a \$ pair spanning a newline is not treated as maths', () {
      expect(kinds('costs \$5\nor \$9 total'), [ContentFragmentKind.text]);
    });

    test('empty delimiters are left as prose rather than empty maths', () {
      expect(kinds(r'$$'), [ContentFragmentKind.text]);
    });

    test('the source including delimiters is retained for the fallback', () {
      final f = MathText.parseFragments(r'$\frac{1}{2}$').single;
      expect(f.value, r'\frac{1}{2}');
      expect(f.source, r'$\frac{1}{2}$');
    });

    test('bold is captured and its markers consumed', () {
      final f = MathText.parseFragments('**Correct Choice: B**').single;
      expect(f.kind, ContentFragmentKind.bold);
      expect(f.value, 'Correct Choice: B');
    });

    test('a real imported explanation lexes into prose and two formulas', () {
      // Shape taken from the ESSLCE question banks.
      const raw = r'The period $T$ is:'
          '\n'
          r'$$T = \frac{2\pi}{B} = \frac{4\pi}{3}$$';
      final f = MathText.parseFragments(raw);
      expect(f.map((e) => e.kind), [
        ContentFragmentKind.text,
        ContentFragmentKind.inlineMath,
        ContentFragmentKind.text,
        ContentFragmentKind.displayMath,
      ]);
      expect(f[1].value, 'T');
      expect(f[3].value, r'T = \frac{2\pi}{B} = \frac{4\pi}{3}');
    });

    test('prose round-trips with nothing added or dropped', () {
      const raw = r'Given $a=1$ and $b=2$, find $a+b$ exactly.';
      final rebuilt = MathText.parseFragments(raw)
          .map((f) => f.kind == ContentFragmentKind.text ? f.value : f.source!)
          .join();
      expect(rebuilt, raw);
    });
  });

  group('rendering', () {
    testWidgets('plain prose renders as a single Text and no Math', (tester) async {
      await tester.pumpWidget(host('Which genus oxidises hydrogen sulphide?'));
      expect(find.text('Which genus oxidises hydrogen sulphide?'), findsOneWidget);
      expect(find.byType(Math), findsNothing);
    });

    testWidgets('inline \$..\$ renders one Math in text style', (tester) async {
      await tester.pumpWidget(host(r'Area of $\triangle ABC$ in units'));
      final maths = tester.widgetList<Math>(find.byType(Math)).toList();
      expect(maths, hasLength(1));
      expect(maths.single.mathStyle, MathStyle.text);
      // The prose around the formula survives.
      expect(find.textContaining('Area of'), findsOneWidget);
    });

    testWidgets(r'display $$..$$ renders in display style and scrolls',
        (tester) async {
      await tester.pumpWidget(host(r'Then: $$\frac{2\pi}{9}$$'));
      final maths = tester.widgetList<Math>(find.byType(Math)).toList();
      expect(maths, hasLength(1));
      expect(maths.single.mathStyle, MathStyle.display);
      // A wide equation must not push the card sideways.
      expect(find.byType(SingleChildScrollView), findsOneWidget);
    });

    testWidgets(r'\(..\) and \[..\] both reach the renderer', (tester) async {
      await tester.pumpWidget(host(r'inline \(a^2\) and block \[b^2\]'));
      final maths = tester.widgetList<Math>(find.byType(Math)).toList();
      expect(maths, hasLength(2));
      expect(maths.map((m) => m.mathStyle),
          containsAll([MathStyle.text, MathStyle.display]));
    });

    testWidgets('several formulas in one sentence all render', (tester) async {
      await tester.pumpWidget(host(r'If $a=1$ and $b=2$ then $a+b=3$'));
      expect(find.byType(Math), findsNWidgets(3));
    });

    testWidgets(r'a lone unmatched $ is left as text, not maths',
        (tester) async {
      // Inline maths stops at a newline, so a stray $ cannot swallow a paragraph.
      await tester.pumpWidget(host('The price is \$5\nNext line'));
      expect(find.byType(Math), findsNothing);
    });

    testWidgets('invalid TeX falls back to its own source rather than blank',
        (tester) async {
      await tester.pumpWidget(host(r'broken $\frac{$'));
      await tester.pumpAndSettle();
      // The fallback shows the delimiters too, so a student can see the source.
      expect(find.textContaining(r'\frac{'), findsOneWidget);
    });

    testWidgets('bold wrapping maths renders both the emphasis and the formula',
        (tester) async {
      await tester.pumpWidget(host(r'**Answer: $x=2$**'));
      expect(find.byType(Math), findsOneWidget);
      // The ** markers are consumed, not shown.
      expect(find.textContaining('**'), findsNothing);
    });

    testWidgets('HTML in an explanation is cleaned before display',
        (tester) async {
      await tester.pumpWidget(host('<b>Antonyms:</b> opposite of <i>calm</i>'));
      expect(find.text('Antonyms: opposite of calm'), findsOneWidget);
    });

    testWidgets('selectable mode allows copying a worked solution',
        (tester) async {
      await tester.pumpWidget(host(r'Given $x=2$', selectable: true));
      expect(find.byType(SelectableText), findsOneWidget);
    });

    testWidgets('a real imported explanation renders maths, not markup',
        (tester) async {
      // Shape taken from the ESSLCE question banks: a label, a <br>, and TeX.
      await tester.pumpWidget(host(
        r'<b>Step-by-step derivation:</b><br>'
        r'The period $T$ is:<br>$$T = \frac{2\pi}{B} = \frac{4\pi}{3}$$',
      ));
      expect(find.byType(Math), findsNWidgets(2));
      expect(find.textContaining('<br>'), findsNothing);
      expect(find.textContaining('<b>'), findsNothing);
    });
  });
}
