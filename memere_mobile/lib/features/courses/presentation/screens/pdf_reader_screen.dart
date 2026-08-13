import 'dart:io';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:syncfusion_flutter_pdfviewer/pdfviewer.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../../../core/constants/app_text_styles.dart';
import '../../../../core/storage/secure_pdf_storage.dart';
import '../../../../core/utils/media_url_helper.dart';

/// Professional PDF Document Reader Screen.
/// Powered by Syncfusion PDF Engine for guaranteed, high-performance in-app viewing on all Android devices.
class PdfReaderScreen extends StatefulWidget {
  const PdfReaderScreen({
    super.key,
    required this.title,
    required this.pdfUrl,
    this.content,
  });

  final String title;
  final String pdfUrl;
  final String? content;

  @override
  State<PdfReaderScreen> createState() => _PdfReaderScreenState();
}

class _PdfReaderScreenState extends State<PdfReaderScreen> {
  String? _localPdfPath;
  bool _isDownloading = false;
  double _downloadProgress = 0.0;
  bool _isSavedToDownloads = false;
  bool _showPdfCanvas = false;
  String? _pdfRenderError;
  int _currentPage = 1;
  int _totalPages = 1;
  bool _isNightMode = true;
  late PdfViewerController _pdfViewerController;

  String get _fileKey => SecurePdfStorage.getFileKey(widget.pdfUrl, title: widget.title);

  @override
  void initState() {
    super.initState();
    _pdfViewerController = PdfViewerController();
    _checkDownloadStatus();
  }

  @override
  void dispose() {
    _pdfViewerController.dispose();
    super.dispose();
  }

  Future<void> _checkDownloadStatus() async {
    try {
      final isSaved = await SecurePdfStorage.isDownloaded(_fileKey);
      if (isSaved) {
        final pdfFile = await SecurePdfStorage.getPdfFile(_fileKey);
        if (mounted && await pdfFile.exists()) {
          setState(() {
            _localPdfPath = pdfFile.path;
            _isSavedToDownloads = true;
          });
          return;
        }
      }
    } catch (_) {}
  }

  Future<void> _downloadPdf() async {
    if (_isDownloading) return;
    setState(() {
      _isDownloading = true;
      _downloadProgress = 0.05;
      _pdfRenderError = null;
    });

    try {
      final pdfFile = await SecurePdfStorage.downloadPdf(
        pdfUrl: widget.pdfUrl,
        fileKey: _fileKey,
        title: widget.title,
        content: widget.content,
        onProgress: (progress) {
          if (mounted) {
            setState(() {
              _downloadProgress = progress;
            });
          }
        },
      );

      if (!mounted) return;
      setState(() {
        _isDownloading = false;
        _isSavedToDownloads = true;
        _localPdfPath = pdfFile.path;
        _showPdfCanvas = true;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('PDF downloaded successfully! Opening reader...'),
          duration: Duration(seconds: 2),
        ),
      );
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _isDownloading = false;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Download failed: ${e.toString()}')),
      );
    }
  }

  Future<void> _deletePdf() async {
    try {
      await SecurePdfStorage.deletePdf(_fileKey);
      if (!mounted) return;
      setState(() {
        _localPdfPath = null;
        _isSavedToDownloads = false;
        _showPdfCanvas = false;
        _pdfRenderError = null;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('PDF deleted from local storage.'),
          duration: Duration(seconds: 2),
        ),
      );
    } catch (_) {}
  }

  Future<void> _openInBrowser() async {
    try {
      final rawUrl = widget.pdfUrl.trim();
      final targetUrl = rawUrl.isEmpty ? 'sample.pdf' : rawUrl;
      final resolvedUrl = fixMediaUrl(targetUrl);
      final uri = Uri.parse(resolvedUrl);
      if (await canLaunchUrl(uri)) {
        await launchUrl(uri, mode: LaunchMode.externalApplication);
      }
    } catch (_) {}
  }

  @override
  Widget build(BuildContext context) {
    final bgColor = _isNightMode ? const Color(0xFF0B0F17) : const Color(0xFFF8FAFC);
    final cardColor = _isNightMode ? const Color(0xFF131C2E) : Colors.white;
    final textColor = _isNightMode ? const Color(0xFFF8FAFC) : const Color(0xFF0F172A);
    final mutedColor = _isNightMode ? const Color(0xFF94A3B8) : const Color(0xFF64748B);
    final borderColor = _isNightMode ? const Color(0xFF1E293B) : const Color(0xFFE2E8F0);

    final progressPct = (_downloadProgress * 100).clamp(0, 100).toInt();

    return Scaffold(
      backgroundColor: bgColor,
      appBar: AppBar(
        backgroundColor: cardColor,
        elevation: 0,
        scrolledUnderElevation: 0,
        leading: IconButton(
          icon: Icon(
            _showPdfCanvas ? Icons.arrow_back_rounded : Icons.close_rounded,
            color: textColor,
          ),
          onPressed: () {
            if (_showPdfCanvas) {
              setState(() => _showPdfCanvas = false);
            } else {
              context.pop();
            }
          },
          tooltip: _showPdfCanvas ? 'Back to Notes' : 'Close',
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              widget.title,
              style: AppTextStyles.titleMedium.copyWith(
                color: textColor,
                fontWeight: FontWeight.bold,
                fontSize: 15,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            Text(
              _showPdfCanvas && _localPdfPath != null
                  ? 'Page $_currentPage of $_totalPages'
                  : (_isSavedToDownloads ? 'PDF Ready' : 'PDF Study Material'),
              style: AppTextStyles.caption.copyWith(color: mutedColor, fontSize: 11),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: Icon(Icons.open_in_browser_rounded, color: textColor, size: 20),
            tooltip: 'Open in Browser',
            onPressed: _openInBrowser,
          ),
          if (_isSavedToDownloads) ...[
            IconButton(
              icon: const Icon(Icons.delete_outline_rounded, color: Color(0xFFEF4444), size: 20),
              tooltip: 'Delete Local PDF',
              onPressed: _deletePdf,
            ),
          ],
          IconButton(
            icon: Icon(
              _isNightMode ? Icons.wb_sunny_outlined : Icons.nightlight_round,
              color: textColor,
              size: 20,
            ),
            tooltip: _isNightMode ? 'Light Mode' : 'Dark Mode',
            onPressed: () => setState(() => _isNightMode = !_isNightMode),
          ),
        ],
      ),
      body: SafeArea(
        child: _showPdfCanvas && _localPdfPath != null
            ? _buildPdfCanvasView(cardColor, textColor, mutedColor, borderColor)
            : _buildOverviewView(cardColor, textColor, mutedColor, borderColor, progressPct),
      ),
    );
  }

  /// Full Syncfusion PDF Reader View Canvas
  Widget _buildPdfCanvasView(
    Color cardColor,
    Color textColor,
    Color mutedColor,
    Color borderColor,
  ) {
    return Column(
      children: [
        // Top Return Bar (No Overflow)
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          color: _isNightMode ? const Color(0xFF0F172A) : const Color(0xFFF1F5F9),
          child: Row(
            children: [
              Expanded(
                child: Text(
                  'PDF Document Reader',
                  style: TextStyle(fontSize: 12, color: mutedColor, fontWeight: FontWeight.w500),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              TextButton.icon(
                onPressed: () => setState(() => _showPdfCanvas = false),
                icon: const Icon(Icons.article_outlined, size: 14),
                label: const Text('Notes', style: TextStyle(fontSize: 12)),
              ),
              IconButton(
                icon: const Icon(Icons.delete_outline_rounded, color: Color(0xFFEF4444), size: 18),
                onPressed: _deletePdf,
                tooltip: 'Delete PDF',
              ),
            ],
          ),
        ),

        // Syncfusion Native Canvas PDF Viewer
        Expanded(
          child: _pdfRenderError != null
              ? Center(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Icon(Icons.error_outline_rounded, color: Color(0xFFEF4444), size: 48),
                        const SizedBox(height: 12),
                        Text(
                          'Could not render PDF file',
                          style: TextStyle(color: textColor, fontWeight: FontWeight.bold, fontSize: 16),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          _pdfRenderError!,
                          style: TextStyle(color: mutedColor, fontSize: 12),
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 20),
                        ElevatedButton.icon(
                          onPressed: _downloadPdf,
                          style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF1E293B)),
                          icon: const Icon(Icons.refresh_rounded, size: 16),
                          label: const Text('Re-download PDF'),
                        ),
                      ],
                    ),
                  ),
                )
              : SfPdfViewer.file(
                  File(_localPdfPath!),
                  key: ValueKey('${_localPdfPath}_${File(_localPdfPath!).existsSync() ? File(_localPdfPath!).lastModifiedSync().millisecondsSinceEpoch : 0}'),
                  controller: _pdfViewerController,
                  canShowScrollHead: true,
                  canShowScrollStatus: true,
                  enableDoubleTapZooming: true,
                  onDocumentLoaded: (PdfDocumentLoadedDetails details) {
                    if (mounted) {
                      setState(() {
                        _totalPages = details.document.pages.count;
                        _pdfRenderError = null;
                      });
                    }
                  },
                  onPageChanged: (PdfPageChangedDetails details) {
                    if (mounted) {
                      setState(() {
                        _currentPage = details.newPageNumber;
                      });
                    }
                  },
                  onDocumentLoadFailed: (PdfDocumentLoadFailedDetails details) {
                    if (mounted) {
                      setState(() {
                        _pdfRenderError = '${details.error}: ${details.description}';
                      });
                    }
                  },
                ),
        ),

        // Bottom Page Controls
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
          decoration: BoxDecoration(
            color: cardColor,
            border: Border(top: BorderSide(color: borderColor)),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              IconButton(
                onPressed: _currentPage > 1
                    ? () => _pdfViewerController.previousPage()
                    : null,
                icon: const Icon(Icons.chevron_left_rounded),
                color: textColor,
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: _isNightMode ? const Color(0xFF0F172A) : const Color(0xFFF1F5F9),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(color: borderColor),
                ),
                child: Text(
                  'Page $_currentPage of $_totalPages',
                  style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: textColor),
                ),
              ),
              IconButton(
                onPressed: _currentPage < _totalPages
                    ? () => _pdfViewerController.nextPage()
                    : () => setState(() => _showPdfCanvas = false),
                icon: Icon(
                  _currentPage < _totalPages ? Icons.chevron_right_rounded : Icons.check_rounded,
                  color: textColor,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  /// Main Overview View: Download/Read + Delete Buttons + ALWAYS Preserves Lesson Text Notes
  Widget _buildOverviewView(
    Color cardColor,
    Color textColor,
    Color mutedColor,
    Color borderColor,
    int progressPct,
  ) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Lesson Title Box
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: cardColor,
              borderRadius: BorderRadius.circular(14),
              border: Border.all(color: borderColor),
            ),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: _isNightMode ? const Color(0xFF0F172A) : const Color(0xFFF1F5F9),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(Icons.picture_as_pdf_outlined, color: textColor, size: 22),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        widget.title,
                        style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: textColor),
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 2),
                      Text(
                        _isSavedToDownloads ? 'PDF Ready in App Storage' : 'PDF Document Available',
                        style: TextStyle(fontSize: 12, color: mutedColor),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // Action Buttons: Download/Read + Delete
          if (_isDownloading) ...[
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: cardColor,
                borderRadius: BorderRadius.circular(14),
                border: Border.all(color: borderColor),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text('Downloading PDF...', style: TextStyle(fontSize: 13, color: textColor, fontWeight: FontWeight.w600)),
                      Text('$progressPct%', style: const TextStyle(fontSize: 13, fontWeight: FontWeight.bold, color: Color(0xFF38BDF8))),
                    ],
                  ),
                  const SizedBox(height: 10),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(4),
                    child: LinearProgressIndicator(
                      value: _downloadProgress > 0 ? _downloadProgress : null,
                      color: const Color(0xFF38BDF8),
                      backgroundColor: const Color(0xFF0F172A),
                      minHeight: 8,
                    ),
                  ),
                ],
              ),
            ),
          ] else ...[
            Row(
              children: [
                Expanded(
                  child: SizedBox(
                    height: 48,
                    child: ElevatedButton.icon(
                      onPressed: _isSavedToDownloads
                          ? () => setState(() => _showPdfCanvas = true)
                          : _downloadPdf,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: _isSavedToDownloads
                            ? const Color(0xFF10B981)
                            : const Color(0xFF1E293B),
                        foregroundColor: Colors.white,
                        elevation: 0,
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(10),
                          side: BorderSide(
                            color: _isSavedToDownloads ? const Color(0xFF10B981) : const Color(0xFF334155),
                          ),
                        ),
                      ),
                      icon: Icon(
                        _isSavedToDownloads ? Icons.menu_book_rounded : Icons.arrow_downward_rounded,
                        size: 20,
                      ),
                      label: Text(
                        _isSavedToDownloads ? 'Read PDF Document' : 'Download PDF',
                        style: const TextStyle(fontSize: 15, fontWeight: FontWeight.bold),
                      ),
                    ),
                  ),
                ),
                if (_isSavedToDownloads) ...[
                  const SizedBox(width: 10),
                  SizedBox(
                    height: 48,
                    width: 48,
                    child: OutlinedButton(
                      onPressed: _deletePdf,
                      style: OutlinedButton.styleFrom(
                        padding: EdgeInsets.zero,
                        side: const BorderSide(color: Color(0xFFEF4444)),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(10),
                        ),
                      ),
                      child: const Icon(Icons.delete_outline_rounded, color: Color(0xFFEF4444), size: 20),
                    ),
                  ),
                ],
              ],
            ),
          ],
          const SizedBox(height: 20),

          // Lesson Text Notes & Content (ALWAYS PRESERVED, NEVER DELETED OR HIDDEN)
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(18),
            decoration: BoxDecoration(
              color: cardColor,
              borderRadius: BorderRadius.circular(14),
              border: Border.all(color: borderColor),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.article_outlined, size: 18, color: mutedColor),
                    const SizedBox(width: 8),
                    Text(
                      'Lesson Notes & Study Guide',
                      style: TextStyle(fontSize: 14, fontWeight: FontWeight.bold, color: textColor),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Divider(color: borderColor),
                const SizedBox(height: 12),
                Text(
                  widget.content != null && widget.content!.trim().isNotEmpty
                      ? widget.content!
                      : 'Comprehensive Grade 12 National Exam preparation material:\n\n'
                          '• Core definitions, formulas, and key concepts for "${widget.title}".\n'
                          '• Step-by-step exam guidelines and practice questions.\n'
                          '• Download the PDF above for offline viewing.',
                  style: TextStyle(fontSize: 14, height: 1.6, color: textColor),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
