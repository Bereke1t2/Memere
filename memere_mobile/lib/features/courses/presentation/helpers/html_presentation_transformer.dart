import 'package:flutter/material.dart';

/// Helper utility to transform and wrap HTML presentations and documents
/// for responsive, overflow-safe rendering in mobile WebViews.
class HtmlPresentationTransformer {
  /// Normalizes viewport meta tags, injects responsive CSS containment rules,
  /// theme variables, and an adaptive JavaScript auto-fitter for slide decks and tables.
  static String wrapHtml({
    required String rawHtml,
    required bool isDark,
    Color? linkColor,
  }) {
    var html = rawHtml;
    final textColor = isDark ? '#E2E8F0' : '#0F172A';
    final bgColor = isDark ? '#0B0E14' : '#FFFFFF';
    final borderColor = isDark ? '#334155' : '#E2E8F0';
    final codeBg = isDark ? '#1E293B' : '#F1F5F9';
    final codeColor = isDark ? '#38BDF8' : '#0369A1';
    final linkHex = linkColor != null
        ? '#${linkColor.toARGB32().toRadixString(16).padLeft(8, '0').substring(2)}'
        : '#10B981';

    const responsiveViewportTag =
        '<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0, minimum-scale=0.25, user-scalable=yes">';

    // 1. Normalize viewport meta tag (replace existing restrictive viewport or inject new)
    final viewportRegex = RegExp(
      r'<meta\s+[^>]*name=["\x27]viewport["\x27][^>]*>',
      caseSensitive: false,
    );

    if (viewportRegex.hasMatch(html)) {
      html = html.replaceFirst(viewportRegex, responsiveViewportTag);
    } else if (html.toLowerCase().contains('<head>')) {
      html = html.replaceFirst(
          RegExp(r'<head>', caseSensitive: false), '<head>$responsiveViewportTag');
    } else if (html.toLowerCase().contains('<head ')) {
      html = html.replaceFirst(
          RegExp(r'(<head[^>]*>)', caseSensitive: false), '\$1$responsiveViewportTag');
    } else {
      html = '$responsiveViewportTag$html';
    }

    // 2. Responsive CSS & Theme Stylesheet
    final responsiveStyles = '''
<style id="memere-responsive-presentation-styles">
  :root {
    color-scheme: ${isDark ? 'dark' : 'light'};
    --bg-primary: $bgColor;
    --text-primary: $textColor;
    --border-color: $borderColor;
    --code-bg: $codeBg;
    --code-color: $codeColor;
    --link-color: $linkHex;
  }
  html {
    box-sizing: border-box;
    -webkit-text-size-adjust: 100%;
    width: 100%;
    max-width: 100vw;
  }
  *, *::before, *::after {
    box-sizing: inherit;
  }
  html, body {
    margin: 0;
    padding: 0;
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
    -webkit-font-smoothing: antialiased;
  }
  body {
    background-color: var(--bg-primary);
    color: var(--text-primary);
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    line-height: 1.6;
    word-wrap: break-word;
    overflow-wrap: break-word;
    padding: 14px;
  }
  h1, h2, h3, h4, h5, h6 {
    max-width: 100%;
    overflow-wrap: break-word;
    line-height: 1.3;
  }
  p, ul, ol, dl, blockquote {
    max-width: 100%;
  }
  img {
    max-width: 100% !important;
    height: auto;
    object-fit: contain;
  }
  svg {
    max-width: 100% !important;
    height: auto;
  }
  video, audio, embed, object {
    max-width: 100% !important;
  }
  iframe {
    max-width: 100% !important;
    border: 0;
  }
  canvas {
    max-width: 100% !important;
    height: auto;
  }
  a {
    color: var(--link-color);
  }
  .memere-table-wrapper {
    width: 100%;
    max-width: 100%;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    margin: 14px 0;
    border-radius: 8px;
  }
  table {
    border-collapse: collapse;
    width: 100%;
    max-width: 100%;
    margin: 0;
  }
  th, td {
    border: 1px solid var(--border-color);
    padding: 8px 12px;
    text-align: left;
  }
  th {
    background-color: var(--code-bg);
    font-weight: 600;
  }
  pre {
    background: var(--code-bg) !important;
    color: var(--code-color) !important;
    border-radius: 8px;
    padding: 12px;
    max-width: 100%;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    white-space: pre;
    word-wrap: normal;
  }
  code {
    background: var(--code-bg);
    color: var(--code-color);
    border-radius: 4px;
    padding: 2px 5px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    font-size: 0.9em;
  }
  pre code {
    background: transparent !important;
    padding: 0;
  }
  .katex-display, .MathJax_Display, .math-display {
    max-width: 100%;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    padding: 4px 0;
  }
  .reveal-viewport, .reveal {
    width: 100% !important;
    height: 100% !important;
    max-width: 100vw !important;
    max-height: 100vh !important;
    overflow: hidden !important;
  }
  .marp-container, .deck-container, .presentation-container {
    max-width: 100% !important;
  }
  .slide, [data-slide], section.slide {
    max-width: 100% !important;
    box-sizing: border-box;
  }
</style>
''';

    // 3. Client-Side JavaScript Auto-Fit & Responsive Containment Script
    const responsiveScript = '''
<script id="memere-responsive-presentation-script">
(function() {
  function initResponsivePresentation() {
    try {
      // 1. Wrap uncontained tables in a horizontal scroll container
      var tables = document.querySelectorAll('table');
      tables.forEach(function(table) {
        if (!table.parentElement || !table.parentElement.classList.contains('memere-table-wrapper')) {
          var wrapper = document.createElement('div');
          wrapper.className = 'memere-table-wrapper';
          table.parentNode.insertBefore(wrapper, table);
          wrapper.appendChild(table);
        }
      });

      // 2. Slide presentation proportional auto-fit
      function autoFitSlides() {
        var winWidth = window.innerWidth || document.documentElement.clientWidth;
        if (!winWidth || winWidth <= 0) return;

        // Reveal.js layout sync
        if (window.Reveal && typeof window.Reveal.layout === 'function') {
          window.Reveal.layout();
          return;
        }

        // Custom slide containers or fixed-width presentation roots
        var slideNodes = document.querySelectorAll('.slides, .slide, .deck-container, .presentation-deck, [data-slide-container]');
        if (slideNodes.length > 0) {
          slideNodes.forEach(function(node) {
            var origWidth = node.getAttribute('data-original-width');
            if (!origWidth) {
              origWidth = node.offsetWidth || 960;
              node.setAttribute('data-original-width', origWidth);
            } else {
              origWidth = parseFloat(origWidth);
            }

            if (origWidth > winWidth) {
              var targetWidth = winWidth - 16;
              var scaleRatio = targetWidth / origWidth;
              node.style.transformOrigin = 'center top';
              node.style.transform = 'scale(' + scaleRatio + ')';
              node.style.margin = '0 auto';
            } else {
              node.style.transform = 'none';
            }
          });
        }
      }

      autoFitSlides();
      window.addEventListener('resize', autoFitSlides);
      window.addEventListener('orientationchange', function() {
        setTimeout(autoFitSlides, 100);
      });
    } catch (e) {
      console.warn('Responsive presentation init warning:', e);
    }
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initResponsivePresentation);
  } else {
    initResponsivePresentation();
  }
})();
</script>
''';

    final injection = '$responsiveStyles\n$responsiveScript';

    if (html.toLowerCase().contains('</head>')) {
      html = html.replaceFirst(
          RegExp(r'</head>', caseSensitive: false), '$injection</head>');
    } else {
      html = '<head>$injection</head>$html';
    }

    return html;
  }
}
