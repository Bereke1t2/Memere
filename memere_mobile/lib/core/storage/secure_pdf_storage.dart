import 'dart:io';
import 'dart:typed_data';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:path_provider/path_provider.dart';
import 'package:syncfusion_flutter_pdf/pdf.dart';

import '../storage/secure_storage_service.dart';
import '../utils/media_url_helper.dart';

/// Secure Local PDF Storage Service.
/// Downloads, validates, and stores PDF files in app-private sandbox storage
/// (`/data/user/0/.../app_flutter/pdfs/`), preventing raw files from being shared outside Memere.
class SecurePdfStorage {
  /// Generates a safe filename key based on the pdfUrl or lesson title
  static String getFileKey(String url) {
    final clean = url.trim().isEmpty ? 'sample.pdf' : url.trim();
    return 'pdf_${clean.replaceAll(RegExp(r'[^a-zA-Z0-9]'), '_')}';
  }

  /// Gets local File reference for a given fileKey
  static Future<File> getPdfFile(String fileKey) async {
    final dir = await getApplicationSupportDirectory();
    final pdfDir = Directory('${dir.path}/pdfs');
    if (!await pdfDir.exists()) {
      await pdfDir.create(recursive: true);
    }
    return File('${pdfDir.path}/$fileKey.pdf');
  }

  /// Checks if a valid PDF file exists in app local private storage
  static Future<bool> isDownloaded(String fileKey) async {
    try {
      final file = await getPdfFile(fileKey);
      return await file.exists() && (await file.length()) > 200;
    } catch (_) {
      return false;
    }
  }

  /// Deletes a local PDF file from app private storage
  static Future<void> deletePdf(String fileKey) async {
    try {
      final file = await getPdfFile(fileKey);
      if (await file.exists()) {
        await file.delete();
      }
    } catch (_) {}
  }

  /// Downloads a PDF from backend with live progress callback, validates PDF header,
  /// saves to app sandbox storage, and returns the local File.
  /// If no remote PDF file is uploaded but lesson text content exists, compiles
  /// the actual lesson text content into a clean PDF document.
  static Future<File> downloadPdf({
    required String pdfUrl,
    required String fileKey,
    String? title,
    String? content,
    void Function(double progress)? onProgress,
  }) async {
    final file = await getPdfFile(fileKey);

    final rawUrl = pdfUrl.trim();
    final isPlaceholder = rawUrl.isEmpty || rawUrl == 'sample.pdf';

    // 1. If remote PDF URL is available, attempt network download
    if (!isPlaceholder) {
      final resolvedUrl = fixMediaUrl(rawUrl);
      final uri = Uri.tryParse(resolvedUrl);
      final safeUrl = uri != null ? uri.toString() : Uri.encodeFull(resolvedUrl);

      final token = await SecureStorageService().getAccessToken();
      final options = Options(
        headers: token != null && token.isNotEmpty ? {'Authorization': 'Bearer $token'} : null,
        responseType: ResponseType.bytes,
        receiveTimeout: const Duration(seconds: 30),
        sendTimeout: const Duration(seconds: 15),
      );

      final dio = Dio();
      try {
        final response = await dio.get<List<int>>(
          safeUrl,
          options: options,
          onReceiveProgress: (received, total) {
            if (total > 0 && onProgress != null) {
              onProgress(received / total);
            }
          },
        );

        if (response.statusCode == 200 && response.data != null && response.data!.length > 100) {
          final bytes = Uint8List.fromList(response.data!);
          if (_isValidPdfBytes(bytes)) {
            await file.writeAsBytes(bytes, flush: true);
            return file;
          } else {
            // Check if server returned a JSON error message instead of PDF
            final responseStr = String.fromCharCodes(bytes.take(200));
            throw Exception('Server returned non-PDF content: $responseStr');
          }
        } else {
          throw Exception('Backend returned status code ${response.statusCode}');
        }
      } catch (e) {
        // If real PDF download failed, but lesson has actual text content, fallback to compiling actual text content
        if (content != null && content.trim().isNotEmpty) {
          final textPdfBytes = _generatePdfFromLessonContent(
            title ?? 'Lesson Study Guide',
            content,
          );
          await file.writeAsBytes(textPdfBytes, flush: true);
          return file;
        }
        // Otherwise rethrow explicit error so UI displays exact issue to student
        throw Exception('Could not download PDF from server: ${e.toString()}');
      }
    }

    // 2. If lesson has text content but no uploaded PDF file, compile text content to PDF
    final textPdfBytes = _generatePdfFromLessonContent(
      title ?? 'Lesson Study Guide',
      content ?? 'Comprehensive Grade 12 National Exam preparation material.',
    );

    await file.writeAsBytes(textPdfBytes, flush: true);
    return file;
  }

  /// Scans the binary bytes for `%PDF-` signature within the first 1024 bytes
  static bool _isValidPdfBytes(Uint8List bytes) {
    if (bytes.length < 5) return false;
    final limit = bytes.length < 1024 ? bytes.length - 4 : 1020;
    for (int i = 0; i < limit; i++) {
      if (bytes[i] == 0x25 && // %
          bytes[i + 1] == 0x50 && // P
          bytes[i + 2] == 0x44 && // D
          bytes[i + 3] == 0x46 && // F
          bytes[i + 4] == 0x2D) { // -
        return true;
      }
    }
    return false;
  }

  /// Sanitizes text to ensure all characters are supported by Syncfusion PdfStandardFont (WinAnsi / ASCII)
  static String _sanitizeTextForPdfStandardFont(String text) {
    if (text.isEmpty) return text;

    // Replace common Unicode punctuation & formatting symbols with safe ASCII equivalents
    var cleaned = text
        .replaceAll('“', '"')
        .replaceAll('”', '"')
        .replaceAll('‘', "'")
        .replaceAll('’', "'")
        .replaceAll('—', '-')
        .replaceAll('–', '-')
        .replaceAll('•', '*')
        .replaceAll('…', '...')
        .replaceAll('\u00A0', ' ')
        .replaceAll('≤', '<=')
        .replaceAll('≥', '>=')
        .replaceAll('≠', '!=');

    // Strip/replace any code point outside 32..255 (WinAnsi range) to prevent PdfStandardFont ArgumentError
    final buffer = StringBuffer();
    for (final char in cleaned.runes) {
      if ((char >= 32 && char <= 255) || char == 10 || char == 13 || char == 9) {
        buffer.writeCharCode(char);
      } else {
        buffer.write(' ');
      }
    }

    final result = buffer.toString().trim();
    return result.isEmpty ? 'Study Notes Content' : result;
  }

  /// Generates a valid Syncfusion PDF document containing the actual lesson title and content
  static Uint8List _generatePdfFromLessonContent(String title, String content) {
    final safeTitle = _sanitizeTextForPdfStandardFont(title.isEmpty ? 'Lesson Study Guide' : title);
    final safeContent = _sanitizeTextForPdfStandardFont(
      content.trim().isNotEmpty ? content : 'No detailed notes provided for this lesson.',
    );

    try {
      final document = PdfDocument();
      final page = document.pages.add();
      final graphics = page.graphics;

      final titleFont = PdfStandardFont(PdfFontFamily.helvetica, 16, style: PdfFontStyle.bold);
      final subtitleFont = PdfStandardFont(PdfFontFamily.helvetica, 11, style: PdfFontStyle.italic);
      final bodyFont = PdfStandardFont(PdfFontFamily.helvetica, 11);

      // Draw Title
      graphics.drawString(
        safeTitle,
        titleFont,
        bounds: const Rect.fromLTWH(0, 0, 500, 26),
      );

      // Draw Header Subtitle
      graphics.drawString(
        'Memere Grade 12 Exam Prep Study Guide',
        subtitleFont,
        bounds: const Rect.fromLTWH(0, 28, 500, 18),
      );

      // Draw Line Separator
      graphics.drawLine(
        PdfPen(PdfColor(200, 200, 200), width: 1),
        const Offset(0, 50),
        const Offset(500, 50),
      );

      // Draw Content Body Text
      final layoutElement = PdfTextElement(
        text: safeContent,
        font: bodyFont,
        brush: PdfSolidBrush(PdfColor(30, 30, 30)),
      );

      final layoutFormat = PdfLayoutFormat(
        layoutType: PdfLayoutType.paginate,
      );

      layoutElement.draw(
        page: page,
        bounds: const Rect.fromLTWH(0, 60, 500, 700),
        format: layoutFormat,
      );

      final List<int> bytes = document.saveSync();
      document.dispose();
      return Uint8List.fromList(bytes);
    } catch (_) {
      // Emergency fallback if PDF drawing fails for any reason
      final document = PdfDocument();
      final page = document.pages.add();
      final font = PdfStandardFont(PdfFontFamily.helvetica, 12);
      page.graphics.drawString('Memere Lesson Notes', font, bounds: const Rect.fromLTWH(0, 0, 500, 30));
      final List<int> bytes = document.saveSync();
      document.dispose();
      return Uint8List.fromList(bytes);
    }
  }
}
