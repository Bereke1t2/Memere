import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

import '../../core/constants/app_colors.dart';

/// What a run of assessment copy turned out to be.
enum ContentFragmentKind {
  /// Ordinary prose.
  text,

  /// `$…$` or `\(…\)` — set on the text baseline, inside the sentence.
  inlineMath,

  /// `$$…$$` or `\[…\]` — set on its own line, scrollable if wide.
  displayMath,

  /// `**…**` or `__…__`.
  bold,
}

/// One run of content, as produced by [MathText.parseFragments].
///
/// Splitting the lexing out of the rendering keeps the part that can silently
/// go wrong — which delimiter won, where a formula started and stopped — in a
/// pure function that is cheap to test directly.
@immutable
class ContentFragment {
  const ContentFragment(this.kind, this.value, {this.source});

  final ContentFragmentKind kind;

  /// The content with its delimiters removed, ready to render.
  final String value;

  /// The original text including delimiters, shown verbatim when a formula
  /// cannot be parsed.
  final String? source;

  @override
  bool operator ==(Object other) =>
      other is ContentFragment &&
      other.kind == kind &&
      other.value == value &&
      other.source == source;

  @override
  int get hashCode => Object.hash(kind, value, source);

  @override
  String toString() => '${kind.name}(${value.length > 40 ? '${value.substring(0, 40)}…' : value})';
}

/// Renders assessment copy that mixes prose with LaTeX.
///
/// Question banks are authored for a MathJax-rendering page, so their text,
/// options and explanations contain TeX delimiters: `$x^2$` inline, `$$…$$` for a
/// displayed equation, and the `\(…\)` / `\[…\]` forms. Rendering those through a
/// plain [Text] shows a student the raw source, which is why this exists.
///
/// Four content concerns are handled here so every assessment surface treats
/// them identically:
///
///  * **Maths** is laid out by flutter_math_fork as a [WidgetSpan] inside the
///    surrounding sentence, so an inline formula sits on the text baseline
///    rather than breaking the paragraph apart.
///  * **A formula that fails to parse** falls back to its own source in a
///    monospace accent rather than throwing or rendering blank — a student is
///    better served by seeing `\frac{1}{2}` than by an empty gap where the
///    answer should be.
///  * **Wide equations** scroll horizontally on their own instead of overflowing
///    the card, which matters on a phone.
///  * **Stray HTML** (`<br>`, `<b>`, `&amp;`) is tolerated. Imported content is
///    cleaned before it is stored, but content predating the importer can still
///    carry tags, and a student should never be shown markup.
class MathText extends StatelessWidget {
  const MathText(
    this.text, {
    super.key,
    required this.style,
    this.mathColor,
    this.textAlign = TextAlign.start,
    this.selectable = false,
  });

  final String text;

  /// Base style for the prose. Maths inherits its size and colour from this so
  /// a formula never looks pasted into the sentence.
  final TextStyle style;

  /// Overrides the colour of formulas only — used where maths should stay
  /// legible against a tinted surface that dims the prose.
  final Color? mathColor;

  final TextAlign textAlign;

  /// Lets a student copy a question or a worked solution. Off by default: it is
  /// deliberately disabled on live attempt screens.
  final bool selectable;

  @override
  Widget build(BuildContext context) {
    final cleaned = normalizeContent(text);

    // Nothing to interpret: skip the span machinery entirely so the common case
    // stays as cheap as a plain Text.
    if (!looksRich(cleaned)) {
      return selectable
          ? SelectableText(cleaned, style: style, textAlign: textAlign)
          : Text(cleaned, style: style, textAlign: textAlign);
    }

    final span = TextSpan(children: _buildSpans(cleaned, style, mathColor));
    return selectable
        ? SelectableText.rich(span, textAlign: textAlign)
        : Text.rich(span, textAlign: textAlign);
  }

  /// True when [content] contains anything this widget would treat specially.
  /// Used to skip rich layout for ordinary prose.
  static bool looksRich(String content) =>
      content.contains(r'$') ||
      content.contains(r'\(') ||
      content.contains(r'\[') ||
      content.contains('**') ||
      content.contains('__');

  /// Strips leftover HTML and decodes the handful of entities that appear in
  /// authored content, turning `<br>` and closing block tags into line breaks.
  ///
  /// LaTeX is left completely alone: `$`, backslashes and braces are content
  /// here, not markup.
  static String normalizeContent(String raw) {
    if (raw.isEmpty) return raw;

    var out = raw;
    if (out.contains('<')) {
      out = out
          .replaceAll(RegExp(r'<\s*br\s*/?\s*>', caseSensitive: false), '\n')
          .replaceAll(
            RegExp(r'<\s*/\s*(p|div|li|tr)\s*>', caseSensitive: false),
            '\n',
          )
          .replaceAll(RegExp(r'<[^>]+>'), '');
    }
    if (out.contains('&')) {
      out = out
          .replaceAll('&nbsp;', ' ')
          .replaceAll('&lt;', '<')
          .replaceAll('&gt;', '>')
          .replaceAll('&quot;', '"')
          .replaceAll('&#39;', "'")
          // Ampersand last, so a double-encoded entity is not mangled.
          .replaceAll('&amp;', '&');
    }

    return out
        .split('\n')
        .map((line) => line.trimRight())
        .join('\n')
        .replaceAll(RegExp(r'\n{3,}'), '\n\n')
        .trim();
  }

  /// Splits [content] into prose, emphasis and maths runs.
  ///
  /// The alternatives are ordered longest-delimiter-first so `$$…$$` is matched
  /// before the `$…$` rule can claim its opening pair. Inline maths also stops
  /// at a newline, which keeps a lone `$` (a price, a stray character) from
  /// swallowing the rest of a paragraph looking for its partner.
  ///
  /// Exposed for testing: this is where the delimiter precedence lives, and
  /// getting it wrong silently mangles a formula rather than throwing.
  @visibleForTesting
  static List<ContentFragment> parseFragments(String content) {
    final out = <ContentFragment>[];
    if (content.isEmpty) return out;

    final pattern = RegExp(
      r'(\$\$[\s\S]+?\$\$'
      r'|\\\[[\s\S]+?\\\]'
      r'|\$[^$\n]+?\$'
      r'|\\\([\s\S]+?\\\)'
      r'|\*\*[^*]+?\*\*'
      r'|__[^_\n]+?__)',
    );

    var cursor = 0;
    for (final match in pattern.allMatches(content)) {
      if (match.start > cursor) {
        out.add(ContentFragment(
          ContentFragmentKind.text,
          content.substring(cursor, match.start),
        ));
      }

      final token = match.group(0)!;
      final isDisplay = token.startsWith(r'$$') || token.startsWith(r'\[');
      final isInline = token.startsWith(r'\(') ||
          (token.startsWith(r'$') && !isDisplay && token.length > 2);

      if (isDisplay || (isInline && !token.startsWith(r'$'))) {
        // Both display forms and \(…\) strip two characters per side.
        out.add(ContentFragment(
          isDisplay
              ? ContentFragmentKind.displayMath
              : ContentFragmentKind.inlineMath,
          token.substring(2, token.length - 2).trim(),
          source: token,
        ));
      } else if (isInline) {
        out.add(ContentFragment(
          ContentFragmentKind.inlineMath,
          token.substring(1, token.length - 1).trim(),
          source: token,
        ));
      } else if (token.startsWith('**') || token.startsWith('__')) {
        out.add(ContentFragment(
          ContentFragmentKind.bold,
          token.substring(2, token.length - 2),
          source: token,
        ));
      } else {
        out.add(ContentFragment(ContentFragmentKind.text, token));
      }

      cursor = match.end;
    }

    if (cursor < content.length) {
      out.add(
        ContentFragment(ContentFragmentKind.text, content.substring(cursor)),
      );
    }
    return out;
  }

  static List<InlineSpan> _buildSpans(
    String content,
    TextStyle baseStyle,
    Color? mathColor,
  ) {
    final spans = <InlineSpan>[];
    for (final fragment in parseFragments(content)) {
      switch (fragment.kind) {
        case ContentFragmentKind.text:
          spans.add(TextSpan(text: fragment.value, style: baseStyle));
        case ContentFragmentKind.inlineMath:
        case ContentFragmentKind.displayMath:
          spans.add(_mathSpan(
            fragment.value,
            fragment.source ?? fragment.value,
            baseStyle,
            mathColor,
            display: fragment.kind == ContentFragmentKind.displayMath,
          ));
        case ContentFragmentKind.bold:
          final bold = baseStyle.copyWith(fontWeight: FontWeight.w700);
          // Emphasis can itself wrap maths ("**Correct: $x=2$**"), so recurse
          // rather than emitting the inner source as literal text.
          if (looksRich(fragment.value)) {
            spans.addAll(_buildSpans(fragment.value, bold, mathColor));
          } else {
            spans.add(TextSpan(text: fragment.value, style: bold));
          }
      }
    }
    return spans;
  }

  static InlineSpan _mathSpan(
    String body,
    String source,
    TextStyle baseStyle,
    Color? mathColor, {
    required bool display,
  }) {
    final mathStyle = baseStyle.copyWith(
      color: mathColor ?? baseStyle.color,
      // flutter_math lays out its own glyphs; an inherited height multiplier
      // would offset the formula from the baseline it sits on.
      height: null,
      fontWeight: display ? baseStyle.fontWeight : null,
    );

    final rendered = Math.tex(
      body,
      mathStyle: display ? MathStyle.display : MathStyle.text,
      textStyle: mathStyle,
      onErrorFallback: (_) => Text(
        source,
        style: baseStyle.copyWith(
          fontFamily: 'monospace',
          color: AppColors.brandEmerald,
        ),
      ),
    );

    return WidgetSpan(
      alignment: PlaceholderAlignment.middle,
      child: display
          // A displayed equation gets its own scrollable line so a long
          // derivation cannot push the card sideways.
          ? Padding(
              padding: const EdgeInsets.symmetric(horizontal: 2, vertical: 4),
              child: SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                physics: const BouncingScrollPhysics(),
                child: rendered,
              ),
            )
          : Padding(
              padding: const EdgeInsets.symmetric(horizontal: 1.5),
              child: rendered,
            ),
    );
  }
}
