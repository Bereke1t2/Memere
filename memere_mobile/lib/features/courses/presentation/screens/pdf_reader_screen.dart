import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_pdfview/flutter_pdfview.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/storage/secure_pdf_storage.dart';
import '../../../../shared/widgets/ai_robot_mascot.dart';

/// In-App PDF & Study Notes Reader for Memere.
/// Uses the high-performance native PDFView background engine (from btluBook-Store)
/// while maintaining the full Obsidian header, dual-mode tabs, font controls,
/// and note-taking UI.
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
  double _fontSize = 15.0;
  int _activeTab = 0; // 0: Study Notes, 1: PDF Document
  bool _isSwipeHorizontal = false;

  PDFViewController? _pdfViewController;
  final List<String> _userNotes = [];

  String get _effectiveContent =>
      SecurePdfStorage.getEffectiveContent(widget.title, widget.content);

  String get _fileKey =>
      SecurePdfStorage.getFileKey(widget.pdfUrl, title: widget.title);

  int get _progressPercent =>
      _totalPages > 0 ? ((_currentPage / _totalPages) * 100).clamp(0, 100).toInt() : 0;

  @override
  void initState() {
    super.initState();

    // Default to PDF Document if a specific PDF is attached, otherwise Study Notes
    if (widget.pdfUrl.trim().isNotEmpty && widget.pdfUrl.trim() != 'sample.pdf') {
      _activeTab = 1;
    } else {
      _activeTab = 0;
    }

    _initializePdf();
  }

  Future<void> _initializePdf() async {
    setState(() {
      _isDownloading = true;
      _downloadProgress = 0.1;
      _pdfRenderError = null;
    });

    try {
      final isSaved = await SecurePdfStorage.isDownloaded(_fileKey);
      if (isSaved) {
        final pdfFile = await SecurePdfStorage.getPdfFile(_fileKey);
        if (mounted && await pdfFile.exists()) {
          setState(() {
            _localPdfPath = pdfFile.path;
            _isSavedToDownloads = true;
            _showPdfCanvas = true;
            _isDownloading = false;
          });
          return;
        }
      }
    } catch (_) {}

    if (mounted) {
      _downloadPdf();
    }
  }

  Future<void> _downloadPdf() async {
    if (_isDownloading && _downloadProgress > 0.3) return;
    setState(() {
      _isDownloading = true;
      _downloadProgress = 0.15;
      _pdfRenderError = null;
    });

    try {
      final pdfFile = await SecurePdfStorage.downloadPdf(
        pdfUrl: widget.pdfUrl,
        fileKey: _fileKey,
        lessonId: widget.lessonId,
        title: widget.title,
        content: _effectiveContent,
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
          content: Text('Document removed from local cache.'),
          duration: Duration(seconds: 2),
        ),
      );
    } catch (_) {}
  }

  void _showJumpToPageDialog() {
    final textController = TextEditingController(text: '$_currentPage');
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          backgroundColor: _isNightMode ? AppColors.bgSecondary : Colors.white,
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
          title: Text(
            'Jump to Page',
            style: TextStyle(
              fontFamily: 'Sora',
              fontWeight: FontWeight.bold,
              fontSize: 16,
              color: _isNightMode ? Colors.white : Colors.black87,
            ),
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                'Enter page number (1 to $_totalPages):',
                style: TextStyle(
                  fontSize: 13,
                  color: _isNightMode ? AppColors.textMuted : Colors.black54,
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: textController,
                keyboardType: TextInputType.number,
                autofocus: true,
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: _isNightMode ? Colors.white : Colors.black,
                ),
                decoration: InputDecoration(
                  filled: true,
                  fillColor: _isNightMode ? const Color(0xFF141824) : const Color(0xFFF1F5F9),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                    borderSide: BorderSide.none,
                  ),
                ),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                final page = int.tryParse(textController.text.trim());
                if (page != null && page >= 1 && page <= _totalPages) {
                  _pdfViewController?.setPage(page - 1);
                  setState(() => _currentPage = page);
                  Navigator.pop(context);
                }
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.brandEmerald,
                foregroundColor: Colors.white,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
              ),
              child: const Text('Go'),
            ),
          ],
        );
      },
    );
  }

  void _showAddNoteBottomSheet() {
    final noteCtrl = TextEditingController();
    final cardColor = _isNightMode ? AppColors.bgSecondary : Colors.white;
    final textColor = _isNightMode ? Colors.white : const Color(0xFF0F172A);

    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      isScrollControlled: true,
      builder: (context) {
        return Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(context).viewInsets.bottom,
          ),
          child: Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: cardColor,
              borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
              border: Border.all(
                color: _isNightMode ? AppColors.borderStrong : const Color(0xFFE2E8F0),
              ),
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: Container(
                    width: 36,
                    height: 4,
                    decoration: BoxDecoration(
                      color: Colors.white24,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                ),
                const SizedBox(height: 14),
                Row(
                  children: [
                    const Icon(Icons.edit_note_rounded, color: AppColors.brandEmerald, size: 22),
                    const SizedBox(width: 8),
                    Text(
                      'Study Notes for ${widget.title}',
                      style: TextStyle(
                        fontFamily: 'Sora',
                        fontSize: 15,
                        fontWeight: FontWeight.bold,
                        color: textColor,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: noteCtrl,
                  maxLines: 4,
                  autofocus: true,
                  style: TextStyle(color: textColor, fontSize: 14),
                  decoration: InputDecoration(
                    hintText: 'Jot down key formulas, insights, or questions...',
                    hintStyle: const TextStyle(color: AppColors.textMuted, fontSize: 13),
                    filled: true,
                    fillColor: _isNightMode ? const Color(0xFF141824) : const Color(0xFFF8FAFC),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(
                        color: _isNightMode ? AppColors.borderStrong : const Color(0xFFE2E8F0),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                ElevatedButton.icon(
                  onPressed: () {
                    final text = noteCtrl.text.trim();
                    if (text.isNotEmpty) {
                      setState(() {
                        _userNotes.add(text);
                      });
                      Navigator.pop(context);
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Note saved to this lesson!'),
                          behavior: SnackBarBehavior.floating,
                        ),
                      );
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.brandEmerald,
                    foregroundColor: Colors.white,
                    minimumSize: const Size(double.infinity, 44),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                  icon: const Icon(Icons.check_rounded, size: 18),
                  label: const Text('Save Note'),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  void _copyNotesToClipboard() {
    Clipboard.setData(ClipboardData(text: _effectiveContent));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Lesson notes copied to clipboard.'),
        behavior: SnackBarBehavior.floating,
      ),
    );
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
                'Studying "${widget.title}"...\nAsk questions, request step-by-step formula derivations, or solve entrance exam questions with AI!',
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
    final bgColor = _isNightMode ? AppColors.bgPrimary : const Color(0xFFF8FAFC);
    final cardColor = _isNightMode ? AppColors.bgSecondary : Colors.white;
    final textColor = _isNightMode ? AppColors.textPrimary : const Color(0xFF0F172A);
    final mutedColor = _isNightMode ? AppColors.textMuted : const Color(0xFF64748B);
    final borderColor = _isNightMode ? AppColors.borderStrong : const Color(0xFFE2E8F0);

    return Scaffold(
      backgroundColor: bgColor,
      floatingActionButton: FloatingActionButton(
        backgroundColor: AppColors.brandEmerald,
        foregroundColor: Colors.white,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        elevation: 8,
        tooltip: 'Add Note / Ask AI Tutor',
        onPressed: _showAddNoteBottomSheet,
        child: const Icon(Icons.edit_note_rounded, size: 26),
      ),
      body: SafeArea(
        child: Column(
          children: [
            // 1. Premium Obsidian Header
            _buildObsidianHeader(cardColor, textColor, mutedColor, borderColor),

            // 2. Mode Switcher Bar ([ Study Notes ] [ PDF Document ])
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 6, 16, 6),
              child: Container(
                height: 38,
                padding: const EdgeInsets.all(3),
                decoration: BoxDecoration(
                  color: cardColor,
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: borderColor),
                ),
                child: Row(
                  children: [
                    Expanded(
                      child: _ReaderTabItem(
                        label: 'Study Notes',
                        icon: Icons.article_outlined,
                        isSelected: _activeTab == 0,
                        isNightMode: _isNightMode,
                        onTap: () => setState(() => _activeTab = 0),
                      ),
                    ),
                    const SizedBox(width: 4),
                    Expanded(
                      child: _ReaderTabItem(
                        label: 'PDF Document',
                        icon: Icons.picture_as_pdf_outlined,
                        isSelected: _activeTab == 1,
                        isNightMode: _isNightMode,
                        onTap: () {
                          setState(() => _activeTab = 1);
                          if (_localPdfPath == null && !_isDownloading) {
                            _downloadPdf();
                          }
                        },
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // 3. Sub-Toolbar for Study Notes (Font sizing & copy action)
            if (_activeTab == 0)
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 2, 16, 6),
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 5),
                  decoration: BoxDecoration(
                    color: cardColor,
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(color: borderColor),
                  ),
                  child: Row(
                    children: [
                      const Text(
                        'Font Size',
                        style: TextStyle(
                          fontSize: 11.5,
                          fontWeight: FontWeight.w600,
                          color: AppColors.textMuted,
                        ),
                      ),
                      const SizedBox(width: 8),
                      InkWell(
                        onTap: _fontSize > 12.0
                            ? () => setState(() => _fontSize -= 1.5)
                            : null,
                        borderRadius: BorderRadius.circular(6),
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                          decoration: BoxDecoration(
                            color: _isNightMode ? const Color(0xFF181820) : const Color(0xFFE2E8F0),
                            borderRadius: BorderRadius.circular(6),
                          ),
                          child: Text(
                            'A-',
                            style: TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.bold,
                              color: _fontSize > 12.0 ? textColor : mutedColor,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 6),
                      Text(
                        '${_fontSize.toInt()}pt',
                        style: TextStyle(
                          fontSize: 11.5,
                          fontWeight: FontWeight.w600,
                          color: textColor,
                        ),
                      ),
                      const SizedBox(width: 6),
                      InkWell(
                        onTap: _fontSize < 22.0
                            ? () => setState(() => _fontSize += 1.5)
                            : null,
                        borderRadius: BorderRadius.circular(6),
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                          decoration: BoxDecoration(
                            color: _isNightMode ? const Color(0xFF181820) : const Color(0xFFE2E8F0),
                            borderRadius: BorderRadius.circular(6),
                          ),
                          child: Text(
                            'A+',
                            style: TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.bold,
                              color: _fontSize < 22.0 ? textColor : mutedColor,
                            ),
                          ),
                        ),
                      ),
                      const Spacer(),
                      InkWell(
                        onTap: _copyNotesToClipboard,
                        borderRadius: BorderRadius.circular(6),
                        child: Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 4),
                          child: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(Icons.copy_rounded, size: 14, color: mutedColor),
                              const SizedBox(width: 4),
                              Text(
                                'Copy Notes',
                                style: TextStyle(
                                  fontSize: 11.5,
                                  fontWeight: FontWeight.w600,
                                  color: textColor,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),

            // 4. Main Body Content
            Expanded(
              child: _activeTab == 0
                  ? _buildStudyNotesView(textColor, mutedColor, borderColor, cardColor)
                  : (_isDownloading
                      ? _buildLoadingView(cardColor, textColor, mutedColor, borderColor)
                      : (_showPdfCanvas && _localPdfPath != null
                          ? _buildPdfCanvasView(cardColor, textColor, mutedColor, borderColor)
                          : _buildErrorOrEmptyView(cardColor, textColor, mutedColor, borderColor))),
            ),
          ],
        ),
      ),
    );
  }

  /// Floating Header with Back, Title, Progress Pill, and Action Buttons
  Widget _buildObsidianHeader(
    Color cardColor,
    Color textColor,
    Color mutedColor,
    Color borderColor,
  ) {
    return Container(
      padding: const EdgeInsets.fromLTRB(12, 8, 12, 8),
      decoration: BoxDecoration(
        color: cardColor,
        border: Border(bottom: BorderSide(color: borderColor)),
      ),
      child: Row(
        children: [
          // Back icon
          GestureDetector(
            onTap: () => context.pop(),
            child: Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: _isNightMode ? const Color(0xFF161B26) : const Color(0xFFF1F5F9),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                Icons.arrow_back_ios_new_rounded,
                color: textColor,
                size: 16,
              ),
            ),
          ),
          const SizedBox(width: 12),

          // Title & Progress Subtitle
          Expanded(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.title,
                  style: TextStyle(
                    fontFamily: 'Sora',
                    fontSize: 14,
                    fontWeight: FontWeight.w700,
                    color: textColor,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                if (_activeTab == 1 && _showPdfCanvas)
                  Text(
                    'Page $_currentPage of $_totalPages',
                    style: TextStyle(fontSize: 11, color: mutedColor),
                  )
                else
                  Text(
                    'Grade 12 National Exam Prep',
                    style: TextStyle(fontSize: 11, color: mutedColor),
                  ),
              ],
            ),
          ),
          const SizedBox(width: 8),

          // Progress % Badge
          if (_activeTab == 1 && _showPdfCanvas && _totalPages > 0)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 4),
              decoration: BoxDecoration(
                color: const Color(0x1810B981),
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: const Color(0x3510B981)),
              ),
              child: Text(
                '$_progressPercent%',
                style: const TextStyle(
                  fontSize: 11.5,
                  fontWeight: FontWeight.w800,
                  color: AppColors.brandEmerald,
                ),
              ),
            ),

          // Night / Light Mode Toggle
          IconButton(
            icon: Icon(
              _isNightMode ? Icons.wb_sunny_outlined : Icons.nightlight_round,
              color: textColor,
              size: 19,
            ),
            tooltip: _isNightMode ? 'Light Mode' : 'Dark Mode',
            onPressed: () => setState(() => _isNightMode = !_isNightMode),
          ),

          // Options Menu (Reload / Clear Cache)
          if (_activeTab == 1 && _isSavedToDownloads)
            PopupMenuButton<String>(
              icon: Icon(Icons.more_vert_rounded, color: mutedColor, size: 20),
              color: cardColor,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              onSelected: (val) {
                if (val == 'reload') {
                  _downloadPdf();
                } else if (val == 'delete') {
                  _deletePdf();
                }
              },
              itemBuilder: (context) => [
                PopupMenuItem(
                  value: 'reload',
                  child: Row(
                    children: [
                      Icon(Icons.refresh_rounded, size: 16, color: textColor),
                      const SizedBox(width: 8),
                      Text('Reload Document', style: TextStyle(color: textColor, fontSize: 13)),
                    ],
                  ),
                ),
                const PopupMenuItem(
                  value: 'delete',
                  child: Row(
                    children: [
                      Icon(Icons.delete_outline_rounded, size: 16, color: Color(0xFFEF4444)),
                      SizedBox(width: 8),
                      Text('Clear Cache', style: TextStyle(color: Color(0xFFEF4444), fontSize: 13)),
                    ],
                  ),
                ),
              ],
            ),

          // AI Tutor Quick Button
          IconButton(
            icon: const Icon(Icons.smart_toy_outlined, color: AppColors.brandEmerald, size: 20),
            tooltip: 'AI Concept Tutor',
            onPressed: _openAiTutor,
          ),
        ],
      ),
    );
  }

  /// Rich Study Notes View for reading lecture notes and summaries
  Widget _buildStudyNotesView(
    Color textColor,
    Color mutedColor,
    Color borderColor,
    Color cardColor,
  ) {
    final noteContent = _effectiveContent;
    final wordCount = noteContent.split(RegExp(r'\s+')).length;
    final estimatedReadTime = ((wordCount / 180).ceil()).clamp(1, 45);

    return ListView(
      padding: const EdgeInsets.fromLTRB(20, 16, 20, 90),
      children: [
        // Header Meta Pill
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              decoration: BoxDecoration(
                color: const Color(0x1810B981),
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: const Color(0x3510B981)),
              ),
              child: const Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.school_rounded,
                    size: 13,
                    color: AppColors.brandEmerald,
                  ),
                  SizedBox(width: 5),
                  Text(
                    'STUDY NOTES',
                    style: TextStyle(
                      fontSize: 10.5,
                      fontWeight: FontWeight.w800,
                      color: AppColors.brandEmerald,
                      letterSpacing: 0.5,
                    ),
                  ),
                ],
              ),
            ),
            Text(
              '~$estimatedReadTime min read',
              style: TextStyle(
                fontSize: 12,
                color: mutedColor,
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),

        // Title
        Text(
          widget.title,
          style: TextStyle(
            fontFamily: 'Sora',
            fontSize: 21,
            fontWeight: FontWeight.w800,
            color: textColor,
            letterSpacing: -0.4,
            height: 1.25,
          ),
        ),
        const SizedBox(height: 14),

        Divider(color: borderColor, height: 1),
        const SizedBox(height: 16),

        // Formatted Note Paragraphs & Sections
        ..._parseAndRenderNotes(noteContent, textColor, mutedColor, borderColor, cardColor),

        // Personal User Notes Section (if any added)
        if (_userNotes.isNotEmpty) ...[
          const SizedBox(height: 24),
          Text(
            'My Revision Notes (${_userNotes.length})',
            style: TextStyle(
              fontFamily: 'Sora',
              fontSize: 15,
              fontWeight: FontWeight.bold,
              color: textColor,
            ),
          ),
          const SizedBox(height: 10),
          ..._userNotes.map(
            (n) => Container(
              margin: const EdgeInsets.only(bottom: 8),
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: _isNightMode ? const Color(0xFF141824) : const Color(0xFFF1F5F9),
                borderRadius: BorderRadius.circular(10),
                border: Border.all(color: const Color(0x3510B981)),
              ),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Icon(Icons.bookmark_outline_rounded, size: 16, color: AppColors.brandEmerald),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      n,
                      style: TextStyle(fontSize: 13, color: textColor),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],

        const SizedBox(height: 32),

        // End of Lesson Note Card
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: cardColor,
            borderRadius: BorderRadius.circular(16),
            border: Border.all(color: borderColor),
          ),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                alignment: Alignment.center,
                decoration: const BoxDecoration(
                  color: Color(0x1810B981),
                  shape: BoxShape.circle,
                ),
                child: const Icon(
                  Icons.check_circle_rounded,
                  color: AppColors.brandEmerald,
                  size: 22,
                ),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'End of Study Notes',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                        color: textColor,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      'Ready to test your knowledge? Try the mock practice or ask the AI Tutor.',
                      style: TextStyle(
                        fontSize: 11.5,
                        color: mutedColor,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  /// Parses text lines into clean, formatted typographic blocks
  List<Widget> _parseAndRenderNotes(
    String content,
    Color textColor,
    Color mutedColor,
    Color borderColor,
    Color cardColor,
  ) {
    final widgets = <Widget>[];
    final lines = content.split('\n');

    for (int i = 0; i < lines.length; i++) {
      final line = lines[i].trim();
      if (line.isEmpty) {
        widgets.add(const SizedBox(height: 10));
        continue;
      }

      // 1. Heading 1 (# Heading)
      if (line.startsWith('# ')) {
        widgets.add(
          Padding(
            padding: const EdgeInsets.only(top: 14, bottom: 8),
            child: Text(
              line.substring(2).trim(),
              style: TextStyle(
                fontFamily: 'Sora',
                fontSize: _fontSize + 4,
                fontWeight: FontWeight.w800,
                color: textColor,
                letterSpacing: -0.3,
              ),
            ),
          ),
        );
      }
      // 2. Heading 2 (## Heading)
      else if (line.startsWith('## ')) {
        widgets.add(
          Padding(
            padding: const EdgeInsets.only(top: 12, bottom: 6),
            child: Text(
              line.substring(3).trim(),
              style: TextStyle(
                fontFamily: 'Sora',
                fontSize: _fontSize + 2,
                fontWeight: FontWeight.w700,
                color: textColor,
                letterSpacing: -0.2,
              ),
            ),
          ),
        );
      }
      // 3. Heading 3 (### Heading)
      else if (line.startsWith('### ')) {
        widgets.add(
          Padding(
            padding: const EdgeInsets.only(top: 10, bottom: 4),
            child: Text(
              line.substring(4).trim(),
              style: TextStyle(
                fontSize: _fontSize + 1,
                fontWeight: FontWeight.w600,
                color: textColor,
              ),
            ),
          ),
        );
      }
      // 4. Quote / Callout Block (> Note...)
      else if (line.startsWith('> ') || line.startsWith('Important:') || line.startsWith('Note:')) {
        final text = line.startsWith('> ') ? line.substring(2) : line;
        widgets.add(
          Container(
            margin: const EdgeInsets.symmetric(vertical: 8),
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
            decoration: BoxDecoration(
              color: _isNightMode
                  ? const Color(0xFF141824)
                  : const Color(0xFFF1F5F9),
              borderRadius: BorderRadius.circular(10),
              border: Border(
                left: const BorderSide(color: AppColors.brandEmerald, width: 3.5),
                top: BorderSide(color: borderColor),
                right: BorderSide(color: borderColor),
                bottom: BorderSide(color: borderColor),
              ),
            ),
            child: Text(
              text,
              style: TextStyle(
                fontSize: _fontSize - 0.5,
                fontWeight: FontWeight.w500,
                color: textColor,
                height: 1.45,
              ),
            ),
          ),
        );
      }
      // 5. Bullet List Items (- or * or •)
      else if (line.startsWith('- ') || line.startsWith('* ') || line.startsWith('• ')) {
        final bulletText = line.substring(2).trim();
        widgets.add(
          Padding(
            padding: const EdgeInsets.only(bottom: 6, left: 4),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Padding(
                  padding: const EdgeInsets.only(top: 6, right: 8),
                  child: Container(
                    width: 5,
                    height: 5,
                    decoration: const BoxDecoration(
                      color: AppColors.brandEmerald,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
                Expanded(
                  child: Text(
                    _stripMarkdownMarkers(bulletText),
                    style: TextStyle(
                      fontSize: _fontSize,
                      color: textColor,
                      height: 1.5,
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      }
      // 6. Numbered items (1. 2. etc.)
      else if (RegExp(r'^\d+\.\s+').hasMatch(line)) {
        final match = RegExp(r'^(\d+\.)\s+(.*)$').firstMatch(line);
        final numPrefix = match?.group(1) ?? '';
        final bodyText = match?.group(2) ?? line;
        widgets.add(
          Padding(
            padding: const EdgeInsets.only(bottom: 6, left: 4),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(
                  width: 24,
                  child: Text(
                    numPrefix,
                    style: const TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.bold,
                      color: AppColors.brandEmerald,
                    ),
                  ),
                ),
                Expanded(
                  child: Text(
                    _stripMarkdownMarkers(bodyText),
                    style: TextStyle(
                      fontSize: _fontSize,
                      color: textColor,
                      height: 1.5,
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      }
      // 7. Standard Paragraph Text
      else {
        widgets.add(
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: SelectableText(
              _stripMarkdownMarkers(line),
              style: TextStyle(
                fontSize: _fontSize,
                color: textColor.withAlpha(240),
                height: 1.6,
                letterSpacing: 0.1,
              ),
            ),
          ),
        );
      }
    }

    return widgets;
  }

  String _stripMarkdownMarkers(String text) {
    return text
        .replaceAll('**', '')
        .replaceAll('__', '')
        .replaceAll('`', '');
  }

  /// Loading View for PDF
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
              'Opening Study Document...',
              style: TextStyle(
                fontFamily: 'Sora',
                fontSize: 16,
                fontWeight: FontWeight.bold,
                color: textColor,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              'Rendering high-yield exam guide & diagrams.',
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
            const SizedBox(height: 24),
            TextButton.icon(
              onPressed: () => setState(() => _activeTab = 0),
              icon: const Icon(Icons.article_outlined, size: 16),
              label: const Text('Read Study Notes instead'),
              style: TextButton.styleFrom(foregroundColor: AppColors.brandEmerald),
            ),
          ],
        ),
      ),
    );
  }

  /// Full PDF Reader Canvas using native PDFView from btluBook-Store
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
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(
                          Icons.error_outline_rounded,
                          color: Color(0xFFEF4444),
                          size: 48,
                        ),
                        const SizedBox(height: 12),
                        Text(
                          'Could not render PDF document',
                          style: TextStyle(
                            color: textColor,
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          _pdfRenderError!,
                          style: TextStyle(color: mutedColor, fontSize: 12),
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 20),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            ElevatedButton.icon(
                              onPressed: _downloadPdf,
                              style: ElevatedButton.styleFrom(
                                backgroundColor: AppColors.brandEmerald,
                              ),
                              icon: const Icon(Icons.refresh_rounded, size: 16),
                              label: const Text('Re-render Document'),
                            ),
                            const SizedBox(width: 12),
                            OutlinedButton.icon(
                              onPressed: () => setState(() => _activeTab = 0),
                              icon: const Icon(Icons.article_outlined, size: 16),
                              label: const Text('Read Notes'),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                )
              : PDFView(
                  key: ValueKey(
                    '${_localPdfPath}_${_isSwipeHorizontal}_${File(_localPdfPath!).existsSync() ? File(_localPdfPath!).lastModifiedSync().millisecondsSinceEpoch : 0}',
                  ),
                  filePath: _localPdfPath!,
                  enableSwipe: true,
                  swipeHorizontal: _isSwipeHorizontal,
                  autoSpacing: false,
                  pageFling: true,
                  pageSnap: true,
                  fitPolicy: FitPolicy.BOTH,
                  preventLinkNavigation: false,
                  defaultPage: _currentPage > 0 ? _currentPage - 1 : 0,
                  onRender: (pages) {
                    if (mounted) {
                      setState(() {
                        _totalPages = pages ?? 1;
                        _pdfRenderError = null;
                      });
                    }
                  },
                  onError: (error) {
                    if (mounted) {
                      setState(() {
                        _pdfRenderError = error.toString();
                      });
                    }
                  },
                  onPageError: (page, error) {
                    if (mounted) {
                      setState(() {
                        _pdfRenderError = 'Page $page: ${error.toString()}';
                      });
                    }
                  },
                  onViewCreated: (PDFViewController pdfViewController) {
                    _pdfViewController = pdfViewController;
                  },
                  onPageChanged: (int? page, int? total) {
                    if (mounted) {
                      setState(() {
                        _currentPage = (page ?? 0) + 1;
                        if (total != null && total > 0) {
                          _totalPages = total;
                        }
                      });
                    }
                  },
                ),
        ),

        // Bottom Reader Navigation Bar
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          decoration: BoxDecoration(
            color: cardColor,
            border: Border(top: BorderSide(color: borderColor)),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              // Previous Page Button
              IconButton(
                onPressed: _currentPage > 1
                    ? () {
                        _pdfViewController?.setPage((_currentPage - 2).clamp(0, _totalPages - 1));
                      }
                    : null,
                icon: const Icon(Icons.chevron_left_rounded),
                color: textColor,
                tooltip: 'Previous Page',
              ),

              // Page Indicator Pill (Tap to Jump)
              InkWell(
                onTap: _showJumpToPageDialog,
                borderRadius: BorderRadius.circular(20),
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                  decoration: BoxDecoration(
                    color: _isNightMode ? const Color(0xFF0F172A) : const Color(0xFFF1F5F9),
                    borderRadius: BorderRadius.circular(20),
                    border: Border.all(color: borderColor),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        'Page $_currentPage of $_totalPages',
                        style: TextStyle(
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                          color: textColor,
                        ),
                      ),
                      const SizedBox(width: 4),
                      Icon(Icons.unfold_more_rounded, size: 14, color: mutedColor),
                    ],
                  ),
                ),
              ),

              // Next Page Button
              IconButton(
                onPressed: _currentPage < _totalPages
                    ? () {
                        _pdfViewController?.setPage(_currentPage.clamp(0, _totalPages - 1));
                      }
                    : null,
                icon: const Icon(Icons.chevron_right_rounded),
                color: textColor,
                tooltip: 'Next Page',
              ),

              // Scroll Direction Toggle (Horizontal / Vertical)
              IconButton(
                onPressed: () {
                  setState(() {
                    _isSwipeHorizontal = !_isSwipeHorizontal;
                  });
                },
                icon: Icon(
                  !_isSwipeHorizontal
                      ? Icons.swap_horiz_rounded
                      : Icons.swap_vert_rounded,
                  size: 20,
                ),
                color: mutedColor,
                tooltip: !_isSwipeHorizontal
                    ? 'Switch to Horizontal Flip'
                    : 'Switch to Vertical Scroll',
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
            const Icon(
              Icons.menu_book_rounded,
              color: AppColors.brandEmerald,
              size: 54,
            ),
            const SizedBox(height: 14),
            Text(
              widget.title,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontFamily: 'Sora',
                fontSize: 17,
                fontWeight: FontWeight.bold,
                color: textColor,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              _pdfRenderError ?? 'Ready to open lesson study document.',
              style: TextStyle(fontSize: 13, color: mutedColor),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                ElevatedButton.icon(
                  onPressed: _downloadPdf,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.brandEmerald,
                    foregroundColor: Colors.white,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  icon: const Icon(Icons.download_rounded, size: 18),
                  label: const Text('Load Document'),
                ),
                const SizedBox(width: 12),
                OutlinedButton.icon(
                  onPressed: () => setState(() => _activeTab = 0),
                  icon: const Icon(Icons.article_outlined, size: 18),
                  label: const Text('Read Notes'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _ReaderTabItem extends StatelessWidget {
  const _ReaderTabItem({
    required this.label,
    required this.icon,
    required this.isSelected,
    required this.isNightMode,
    required this.onTap,
  });

  final String label;
  final IconData icon;
  final bool isSelected;
  final bool isNightMode;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
        decoration: BoxDecoration(
          color: isSelected
              ? AppColors.brandEmerald
              : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              size: 14,
              color: isSelected
                  ? Colors.white
                  : (isNightMode ? AppColors.textMuted : const Color(0xFF64748B)),
            ),
            const SizedBox(width: 6),
            Text(
              label,
              style: TextStyle(
                fontSize: 12,
                fontWeight: isSelected ? FontWeight.bold : FontWeight.w600,
                color: isSelected
                    ? Colors.white
                    : (isNightMode ? AppColors.textSecondary : const Color(0xFF334155)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}


