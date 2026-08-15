import 'dart:io';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:syncfusion_flutter_pdfviewer/pdfviewer.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/storage/secure_pdf_storage.dart';
import '../../../../shared/widgets/ai_robot_mascot.dart';

/// Secure In-App PDF Document & Lesson Notes Reader.
/// Stored in app-private sandbox storage (`/data/user/0/.../app_flutter/pdfs/`)
/// with DRM-protected playback preventing external sharing or export.
class PdfReaderScreen extends StatefulWidget {
  const PdfReaderScreen({
    super.key,
    required this.title,
    required this.pdfUrl,
    this.lessonId,
    this.content,
  });

  final String title;
  final String pdfUrl;
  final String? lessonId;
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
    _initializeAndOpen();
  }

  @override
  void dispose() {
    _pdfViewerController.dispose();
    super.dispose();
  }

  Future<void> _initializeAndOpen() async {
    try {
      final isSaved = await SecurePdfStorage.isDownloaded(_fileKey);
      if (isSaved) {
        final pdfFile = await SecurePdfStorage.getPdfFile(_fileKey);
        if (mounted && await pdfFile.exists()) {
          setState(() {
            _localPdfPath = pdfFile.path;
            _isSavedToDownloads = true;
            _showPdfCanvas = true;
          });
          return;
        }
      }
    } catch (_) {}

    // If not already cached, start automatic download & compilation
    if (mounted) {
      _downloadPdf();
    }
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
        lessonId: widget.lessonId,
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
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _isDownloading = false;
        _pdfRenderError = e.toString();
      });
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
          content: Text('PDF removed from offline storage cache.'),
          duration: Duration(seconds: 2),
        ),
      );
    } catch (_) {}
  }

  void _openAiTutor() {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Row(
                children: [
                  AiRobotMascot(size: 32),
                  SizedBox(width: 10),
                  Text(
                    'AI Concept Tutor',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textPrimary,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 14),
              Text(
                'Studying "${widget.title}"...\nAsk questions or get instant step-by-step explanations for formulas and exam questions!',
                style: const TextStyle(fontSize: 13, color: AppColors.textMuted),
              ),
              const SizedBox(height: 20),
              ElevatedButton.icon(
                onPressed: () => Navigator.pop(context),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.brandEmerald,
                  foregroundColor: Colors.white,
                  minimumSize: const Size(double.infinity, 44),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                icon: const Icon(Icons.smart_toy_rounded, size: 18),
                label: const Text('Start Q&A Session'),
              ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final bgColor = _isNightMode ? const Color(0xFF0B0F17) : const Color(0xFFF8FAFC);
    final cardColor = _isNightMode ? const Color(0xFF131C2E) : Colors.white;
    final textColor = _isNightMode ? const Color(0xFFF8FAFC) : const Color(0xFF0F172A);
    final mutedColor = _isNightMode ? const Color(0xFF94A3B8) : const Color(0xFF64748B);
    final borderColor = _isNightMode ? const Color(0xFF1E293B) : const Color(0xFFE2E8F0);

    return Scaffold(
      backgroundColor: bgColor,
      appBar: AppBar(
        backgroundColor: cardColor,
        elevation: 0,
        scrolledUnderElevation: 0,
        leading: IconButton(
          icon: Icon(
            Icons.arrow_back_rounded,
            color: textColor,
            size: 20,
          ),
          onPressed: () => context.pop(),
          tooltip: 'Back',
        ),
        title: Text(
          widget.title,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.bold,
            color: textColor,
          ),
        ),
        actions: [
          if (_isSavedToDownloads) ...[
            IconButton(
              icon: const Icon(Icons.cached_rounded, size: 20),
              color: mutedColor,
              tooltip: 'Reload PDF Document',
              onPressed: _downloadPdf,
            ),
            IconButton(
              icon: const Icon(Icons.delete_outline_rounded, size: 20),
              color: mutedColor,
              tooltip: 'Clear Offline Cache',
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
      floatingActionButton: AiTutorFab(onPressed: _openAiTutor),
      body: SafeArea(
        child: _isDownloading
            ? _buildLoadingView(cardColor, textColor, mutedColor, borderColor)
            : (_showPdfCanvas && _localPdfPath != null
                ? _buildPdfCanvasView(cardColor, textColor, mutedColor, borderColor)
                : _buildErrorOrEmptyView(cardColor, textColor, mutedColor, borderColor)),
      ),
    );
  }

  /// Downloading & Preparation View
  Widget _buildLoadingView(
    Color cardColor,
    Color textColor,
    Color mutedColor,
    Color borderColor,
  ) {
    final progressPct = (_downloadProgress * 100).clamp(0, 100).toInt();

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(28),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const AiRobotMascot(size: 64, backgroundColor: AppColors.brandEmerald),
            const SizedBox(height: 20),
            Text(
              'Loading Study Material...',
              style: TextStyle(
                fontSize: 17,
                fontWeight: FontWeight.bold,
                color: textColor,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              'Decrypting and caching in secure local storage.',
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 12, color: mutedColor),
            ),
            const SizedBox(height: 24),
            SizedBox(
              width: 220,
              child: Column(
                children: [
                  ClipRRect(
                    borderRadius: BorderRadius.circular(6),
                    child: LinearProgressIndicator(
                      value: _downloadProgress > 0 ? _downloadProgress : null,
                      color: AppColors.brandEmerald,
                      backgroundColor: const Color(0xFF1E293B),
                      minHeight: 8,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '$progressPct%',
                    style: const TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                      color: AppColors.brandEmerald,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// Full Syncfusion PDF Reader Canvas
  Widget _buildPdfCanvasView(
    Color cardColor,
    Color textColor,
    Color mutedColor,
    Color borderColor,
  ) {
    return Column(
      children: [
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
                          style: ElevatedButton.styleFrom(backgroundColor: AppColors.brandEmerald),
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
        // Bottom Navigation & Page Numbering Bar
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
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
                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
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
                    : null,
                icon: const Icon(Icons.chevron_right_rounded),
                color: textColor,
              ),
            ],
          ),
        ),
      ],
    );
  }

  /// Error or Empty Fallback
  Widget _buildErrorOrEmptyView(
    Color cardColor,
    Color textColor,
    Color mutedColor,
    Color borderColor,
  ) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.menu_book_rounded, color: AppColors.brandEmerald, size: 54),
            const SizedBox(height: 14),
            Text(
              widget.title,
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: textColor),
            ),
            const SizedBox(height: 8),
            Text(
              _pdfRenderError ?? 'Ready to open lesson study guide.',
              style: TextStyle(fontSize: 13, color: mutedColor),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 20),
            ElevatedButton.icon(
              onPressed: _downloadPdf,
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.brandEmerald,
                foregroundColor: Colors.white,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              ),
              icon: const Icon(Icons.download_rounded, size: 18),
              label: const Text('Open Study Document'),
            ),
          ],
        ),
      ),
    );
  }
}
