import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:memere_mobile/features/courses/presentation/helpers/html_presentation_transformer.dart';

void main() {
  group('HtmlPresentationTransformer', () {
    test('injects responsive viewport tag when no viewport exists', () {
      const rawHtml = '<html><head><title>Test</title></head><body><h1>Hello</h1></body></html>';
      final result = HtmlPresentationTransformer.wrapHtml(
        rawHtml: rawHtml,
        isDark: false,
      );

      expect(
        result,
        contains('<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0, minimum-scale=0.25, user-scalable=yes">'),
      );
      expect(result, contains('<h1>Hello</h1>'));
    });

    test('normalizes restrictive fixed-width viewport tag', () {
      const rawHtml = '''
<html>
  <head>
    <meta name="viewport" content="width=1280, initial-scale=1.0, user-scalable=no">
    <title>Presentation</title>
  </head>
  <body>
    <div class="slide" style="width: 1280px;">Slide 1</div>
  </body>
</html>
''';
      final result = HtmlPresentationTransformer.wrapHtml(
        rawHtml: rawHtml,
        isDark: true,
      );

      expect(
        result,
        contains('<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0, minimum-scale=0.25, user-scalable=yes">'),
      );
      expect(result, isNot(contains('width=1280')));
      expect(result, contains('<div class="slide" style="width: 1280px;">Slide 1</div>'));
    });

    test('handles HTML with no <head> tag by creating one or prepending', () {
      const rawHtml = '<div class="lesson-content">Simple body content</div>';
      final result = HtmlPresentationTransformer.wrapHtml(
        rawHtml: rawHtml,
        isDark: true,
      );

      expect(result, contains('<meta name="viewport"'));
      expect(result, contains('<style id="memere-responsive-presentation-styles">'));
      expect(result, contains('<script id="memere-responsive-presentation-script">'));
      expect(result, contains('Simple body content'));
    });

    test('injects CSS containment and responsive rules for tables, code, and media', () {
      const rawHtml = '''
<html>
<head><title>Lesson</title></head>
<body>
  <table><tr><td>Row 1</td><td>Row 2</td></tr></table>
  <pre><code>int x = 42;</code></pre>
  <img src="diagram.png" width="800" height="600" />
</body>
</html>
''';
      final result = HtmlPresentationTransformer.wrapHtml(
        rawHtml: rawHtml,
        isDark: true,
      );

      // Box-sizing and root containment
      expect(result, contains('box-sizing: border-box;'));
      expect(result, contains('overflow-x: hidden;'));

      // Table scroll wrapper
      expect(result, contains('.memere-table-wrapper'));
      expect(result, contains('-webkit-overflow-scrolling: touch;'));

      // Media containment
      expect(result, contains('img {'));
      expect(result, contains('max-width: 100% !important;'));
      expect(result, contains('object-fit: contain;'));

      // Code blocks
      expect(result, contains('pre {'));
      expect(result, contains('overflow-x: auto;'));

      // Presentation decks
      expect(result, contains('.reveal-viewport'));
      expect(result, contains('.marp-container'));
      expect(result, contains('.slide'));
    });

    test('applies appropriate dark and light mode color variables', () {
      const rawHtml = '<html><head></head><body><p>Theme Test</p></body></html>';

      final darkResult = HtmlPresentationTransformer.wrapHtml(
        rawHtml: rawHtml,
        isDark: true,
      );
      expect(darkResult, contains('color-scheme: dark;'));
      expect(darkResult, contains('--bg-primary: #0B0E14;'));
      expect(darkResult, contains('--text-primary: #E2E8F0;'));

      final lightResult = HtmlPresentationTransformer.wrapHtml(
        rawHtml: rawHtml,
        isDark: false,
        linkColor: const Color(0xFF10B981),
      );
      expect(lightResult, contains('color-scheme: light;'));
      expect(lightResult, contains('--bg-primary: #FFFFFF;'));
      expect(lightResult, contains('--text-primary: #0F172A;'));
      expect(lightResult, contains('--link-color: #10b981;'));
    });

    test('injects client-side JS auto-fitter script for dynamic resizing & table wrapping', () {
      const rawHtml = '<html><head></head><body><div class="slides">Slide</div></body></html>';
      final result = HtmlPresentationTransformer.wrapHtml(
        rawHtml: rawHtml,
        isDark: true,
      );

      expect(result, contains('initResponsivePresentation'));
      expect(result, contains('memere-table-wrapper'));
      expect(result, contains('autoFitSlides'));
      expect(result, contains('orientationchange'));
    });
  });
}
