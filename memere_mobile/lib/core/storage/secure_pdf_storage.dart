import 'dart:io';
import 'dart:typed_data';
import 'package:dio/dio.dart';
import 'package:path_provider/path_provider.dart';

import '../constants/env.dart';
import '../storage/secure_storage_service.dart';
import '../utils/media_url_helper.dart';

/// No real remote PDF is attached, or the remote content is not a PDF
/// (e.g. a note_url returning HTML/markdown). Reader shows Study Notes.
class PdfNotAvailableException implements Exception {
  PdfNotAvailableException([this.message = 'No PDF document is attached to this lesson.']);
  final String message;
  @override
  String toString() => message;
}

/// A real PDF was expected but could not be downloaded (network/auth/timeout/
/// server). Reader shows an error state with a retry button.
class PdfDownloadException implements Exception {
  PdfDownloadException(this.message);
  final String message;
  @override
  String toString() => message;
}

/// Secure Local PDF Storage Service.
/// Downloads, validates, and stores PDF files in app-private sandbox storage
/// (`/data/user/0/.../app_flutter/pdfs/`), preventing raw files from being shared outside Memere.
class SecurePdfStorage {
  /// Generates a safe, unique filename key based on the pdfUrl or lesson title
  static String getFileKey(String url, {String? title}) {
    final cleanUrl = url.trim();
    if (cleanUrl.isNotEmpty && cleanUrl != 'sample.pdf') {
      return 'pdf_v2_${cleanUrl.replaceAll(RegExp(r'[^a-zA-Z0-9]'), '_')}';
    }
    final cleanTitle = (title ?? 'lesson').trim();
    return 'pdf_v2_note_${cleanTitle.replaceAll(RegExp(r'[^a-zA-Z0-9]'), '_')}';
  }

  /// Gets local File reference for a given fileKey
  static Future<File> getPdfFile(String fileKey) async {
    final dir = await getApplicationDocumentsDirectory();
    final pdfDir = Directory('${dir.path}/pdfs');
    if (!await pdfDir.exists()) {
      await pdfDir.create(recursive: true);
    }
    return File('${pdfDir.path}/$fileKey.pdf');
  }

  /// Checks if a valid PDF file exists in app local private storage, validating PDF header
  static Future<bool> isDownloaded(String fileKey) async {
    try {
      final file = await getPdfFile(fileKey);
      if (!await file.exists()) return false;
      final length = await file.length();
      if (length < 200) {
        await file.delete();
        return false;
      }

      final handle = await file.open(mode: FileMode.read);
      final headerBytes = await handle.read(1024);
      await handle.close();

      final isValid = _isValidPdfBytes(Uint8List.fromList(headerBytes));
      if (!isValid) {
        await file.delete();
        return false;
      }
      return true;
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

  /// Determines the best download URL for a lesson PDF.
  static String? _resolveDownloadUrl(String pdfUrl, {String? lessonId}) {
    final raw = pdfUrl.trim();
    if (raw.isEmpty || raw == 'sample.pdf') return null;

    // Already a full URL (S3 presigned or MinIO direct link)
    if (raw.startsWith('http://') || raw.startsWith('https://')) {
      return fixMediaUrl(raw);
    }

    // It's an S3 key like "lessons/xxx/notes.pdf" or a raw filename
    if (lessonId != null && lessonId.isNotEmpty) {
      final apiBase = fixMediaUrl(Env.baseUrl);
      return '$apiBase/lessons/$lessonId/pdf';
    }

    // Legacy: try to resolve as a relative path
    return fixMediaUrl(raw);
  }

  /// Returns actual lesson notes or rich, topic-aware study guide content if empty
  static String getEffectiveContent(String title, String? rawContent) {
    if (rawContent != null && rawContent.trim().isNotEmpty) {
      return rawContent.trim();
    }

    final t = title.trim().isNotEmpty ? title.trim() : 'Core Exam Topics';
    final lower = t.toLowerCase();

    // 1. Physics Topics
    if (lower.contains('vector') ||
        lower.contains('kinematic') ||
        lower.contains('motion') ||
        lower.contains('force') ||
        lower.contains('newton') ||
        lower.contains('physic') ||
        lower.contains('energy') ||
        lower.contains('wave') ||
        lower.contains('electric') ||
        lower.contains('magnet') ||
        lower.contains('circuit') ||
        lower.contains('optics') ||
        lower.contains('thermodynamic')) {
      return '''
# $t
## Grade 12 National Exam Comprehensive Study Guide

### 1. Key Definitions & Physical Quantities
- Understand the physical meaning and SI units of all variables associated with $t.
- Distinguish clearly between scalar quantities (magnitude only) and vector quantities (magnitude and direction).
- Know standard vector operations: components (Ax = A cos theta, Ay = A sin theta), dot product, cross product, and resultant vectors.

### 2. Fundamental Laws & Mathematical Relations
- Review the core governing equations:
  * Conservation Laws: Energy, Momentum, and Charge conservation principles.
  * Equations of Motion: v = u + at, s = ut + 0.5at^2, v^2 = u^2 + 2as.
  * Newton's Second Law: Sigma F = ma = dp/dt.
  * Work-Energy Theorem: W_net = Delta KE = 0.5m v_f^2 - 0.5m v_i^2.

### 3. National Entrance Exam Problem-Solving Strategy
> High-Yield Tip: National exam physics questions heavily test dimensional consistency, free-body diagram breakdown, and sign conventions.

1. Always draw a clear diagram and set up a coordinate system (+x, +y).
2. List all known and unknown variables with standard SI units.
3. Identify which conservation principle or equation directly connects the given quantities.
4. Check whether friction, air resistance, or energy loss is negligible before simplifying.

### 4. Common Exam Traps & Misconceptions
- Forgetting to resolve vectors into perpendicular components before adding.
- Mixing up mass (kg, scalar) and weight (N, vector: W = mg).
- Omitting negative signs for deceleration or gravitational acceleration (g = -9.8 m/s^2).

### 5. Summary & Practice Checklist
- Memorize the standard constants: g = 9.8 m/s^2, c = 3.0 x 10^8 m/s, e = 1.6 x 10^-19 C.
- Solve at least 5 past national exam questions on $t.
- Use the AI Concept Tutor to get step-by-step solutions for any difficult problems.
''';
    }

    // 2. Mathematics Topics
    if (lower.contains('math') ||
        lower.contains('calculus') ||
        lower.contains('derivative') ||
        lower.contains('integral') ||
        lower.contains('limit') ||
        lower.contains('function') ||
        lower.contains('matrix') ||
        lower.contains('sequence') ||
        lower.contains('series') ||
        lower.contains('probability') ||
        lower.contains('stat') ||
        lower.contains('trig') ||
        lower.contains('algebra') ||
        lower.contains('geometry')) {
      return '''
# $t
## Grade 12 National Exam Comprehensive Study Guide

### 1. Core Mathematical Theory & Definitions
- Understand the fundamental domain, range, and continuity conditions relevant to $t.
- Know the formal definitions and geometric interpretations (e.g. slope of tangent line, area under the curve).

### 2. Standard Formulas & Rules
- Key differentiation & integration rules:
  * Power Rule: d/dx(x^n) = n x^(n-1), integral x^n dx = (x^(n+1))/(n+1) + C.
  * Product Rule: (uv)' = u'v + uv'.
  * Quotient Rule: (u/v)' = (u'v - uv') / v^2.
  * Chain Rule: d/dx[f(g(x))] = f'(g(x)) * g'(x).
- Limit Properties: L'Hopital's Rule for indeterminate forms (0/0, infinity/infinity).

### 3. High-Yield Exam Techniques
> Exam Alert: Grade 12 entrance exams frequently feature composite functions, critical points (f'(x) = 0), and optimization word problems.

1. Test for critical points by finding where f'(x) = 0 or f'(x) does not exist.
2. Apply the First and Second Derivative Tests to classify local extrema and points of inflection.
3. For definite integrals, check for symmetry (odd/even functions) to speed up calculations.

### 4. Common Mistakes to Avoid
- Forgetting the constant of integration (+ C) on indefinite integrals.
- Misapplying the quotient rule formula ordering.
- Forgetting to change integration limits when applying u-substitution.

### 5. Exam Practice Review
- Verify your answers by substitution or graphical sketching.
- Test edge cases (x = 0, x -> infinity, boundary values).
''';
    }

    // 3. Chemistry Topics
    if (lower.contains('chem') ||
        lower.contains('atom') ||
        lower.contains('reaction') ||
        lower.contains('acid') ||
        lower.contains('base') ||
        lower.contains('organic') ||
        lower.contains('mole') ||
        lower.contains('equilibrium') ||
        lower.contains('thermo') ||
        lower.contains('electrochem') ||
        lower.contains('solution') ||
        lower.contains('gas')) {
      return '''
# $t
## Grade 12 National Exam Chemistry Study Guide

### 1. Fundamental Principles & Terminology
- Core principles, definitions, and stoichiometry rules for $t.
- Understand atomic structure, quantum numbers, periodic trends (electronegativity, ionization energy, atomic radius).

### 2. Key Chemical Equations & Formulas
- Moles and Concentration: n = m / M, C = n / V (mol/L).
- Ideal Gas Law: PV = nRT (R = 0.0821 L atm / (mol K) = 8.314 J / (mol K)).
- Equilibrium Constant: K_eq = [Products]^coefficients / [Reactants]^coefficients.
- pH and pOH: pH = -log[H+], pH + pOH = 14 at 25 deg C.

### 3. High-Yield Entrance Exam Tips
> Essential Tip: Le Chatelier's principle and redox balancing (oxidation numbers) appear regularly on national exams.

- For equilibrium shifts: adding reactants shifts reaction forward; increasing pressure shifts toward fewer gas moles; increasing temperature favors endothermic direction.
- In redox reactions: Oxidation is loss of electrons (OIL); Reduction is gain of electrons (RIG).

### 4. Summary & Review
- Balance all chemical equations before performing stoichiometric calculations.
- Review key functional groups in organic chemistry (alkanes, alkenes, alcohols, aldehydes, ketones, carboxylic acids).
''';
    }

    // 4. Biology Topics
    if (lower.contains('bio') ||
        lower.contains('cell') ||
        lower.contains('genetic') ||
        lower.contains('dna') ||
        lower.contains('evolution') ||
        lower.contains('ecology') ||
        lower.contains('human') ||
        lower.contains('plant') ||
        lower.contains('organ') ||
        lower.contains('reproduction') ||
        lower.contains('enzyme')) {
      return '''
# $t
## Grade 12 National Exam Biology Study Guide

### 1. Biological Concepts & Structural Foundations
- Understand the cellular and physiological structures related to $t.
- Key components: organelles, membrane transport (diffusion, osmosis, active transport), cellular respiration, and photosynthesis.

### 2. Genetics & Heredity Principles
- Mendel's Laws: Law of Segregation and Law of Independent Assortment.
- DNA Replication, Transcription (DNA -> mRNA), and Translation (mRNA -> Protein).
- Punnett Squares: Monohybrid (3:1 phenotype ratio) and Dihybrid crosses (9:3:3:1 ratio).

### 3. High-Yield Exam Takeaways
> Key Exam Concept: Enzyme kinetics, ATP yield in aerobic vs anaerobic respiration (36-38 ATP vs 2 ATP), and homeostatic feedback mechanisms.

- Know the exact sequence of cell division: Prophase -> Metaphase -> Anaphase -> Telophase.
- Distinguish between mitosis (2 identical diploid cells) and meiosis (4 genetically diverse haploid gametes).

### 4. Review & Self-Assessment
- Practice diagram labeling and metabolic pathway flowchart recall.
- Use the AI Concept Tutor for detailed explanations of complex biological cycles.
''';
    }

    // 5. Default High-Yield Exam Notes
    return '''
# $t
## Grade 12 National Examination Prep Notes

### 1. Topic Overview & Objectives
- Master the primary theoretical concepts, definitions, and scope of $t.
- Understand how this topic integrates into the national curriculum and exam framework.

### 2. Essential Definitions & Key Concepts
- Study the standard definitions, properties, and classifications for $t.
- Review core terminology and step-by-step methodologies tested in entrance exams.

### 3. High-Yield Exam Strategies & Tips
> Important: In national entrance examinations, questions on $t prioritize analytical thinking, accurate application of principles, and elimination of distractors.

1. Read question stems attentively to identify core requirements and given conditions.
2. Break down multi-step problems systematically before selecting your answer.
3. Review related practice questions and flashcard concepts in the Memere library.

### 4. Summary & Action Steps
- Consolidate your notes and formulas into quick-reference review cards.
- Complete the corresponding quiz and practice mock exam for this unit.
- Ask the AI Concept Tutor if you need instant step-by-step breakdown on any question!
''';
  }

  /// Downloads a PDF from backend with live progress callback, validates PDF header,
  /// saves to app sandbox storage, and returns the local File.
  static Future<File> downloadPdf({
    required String pdfUrl,
    required String fileKey,
    String? lessonId,
    String? title,
    String? content,
    void Function(double progress)? onProgress,
  }) async {
    final file = await getPdfFile(fileKey);
    final resolvedUrl = _resolveDownloadUrl(pdfUrl, lessonId: lessonId);

    // No real remote PDF -> reader should show Study Notes, not a fabricated PDF.
    if (resolvedUrl == null) {
      throw PdfNotAvailableException();
    }

    try {
      final token = await SecureStorageService().getAccessToken();
      final headers = (token != null && token.isNotEmpty)
          ? {'Authorization': 'Bearer $token'}
          : <String, String>{};

      final dio = Dio();
      var currentUrl = resolvedUrl;
      Response<List<int>>? response;

      // Follow redirects manually (max 5 hops) so localhost/emulator hosts get
      // rewritten via fixMediaUrl at EVERY hop (MinIO/S3 presigned chains).
      for (var hop = 0; hop < 5; hop++) {
        response = await dio.get<List<int>>(
          currentUrl,
          options: Options(
            headers: headers,
            responseType: ResponseType.bytes,
            receiveTimeout: const Duration(seconds: 25),
            sendTimeout: const Duration(seconds: 15),
            followRedirects: false,
            validateStatus: (status) => status != null && status < 400,
          ),
          onReceiveProgress: (received, total) {
            if (total > 0 && onProgress != null) {
              onProgress((received / total).clamp(0.0, 1.0));
            }
          },
        );

        final code = response.statusCode ?? 0;
        if (code == 301 || code == 302 || code == 303 || code == 307 || code == 308) {
          final location = response.headers.value('location');
          if (location == null || location.isEmpty) break;
          currentUrl = fixMediaUrl(location);
          continue;
        }
        break;
      }

      if (response != null &&
          response.statusCode == 200 &&
          response.data != null &&
          response.data!.length > 100) {
        final bytes = Uint8List.fromList(response.data!);
        if (_isValidPdfBytes(bytes)) {
          await file.writeAsBytes(bytes, flush: true);
          if (onProgress != null) onProgress(1.0);
          return file;
        }
        // 200 but not a PDF -> most likely a note_url (HTML/markdown). Show notes.
        throw PdfNotAvailableException('The attached file is not a valid PDF.');
      }

      throw PdfDownloadException(
        'Could not download the PDF (status ${response?.statusCode ?? 'unknown'}).',
      );
    } on PdfNotAvailableException {
      rethrow;
    } on PdfDownloadException {
      rethrow;
    } on DioException catch (e) {
      throw PdfDownloadException(_friendlyDioError(e));
    } catch (_) {
      throw PdfDownloadException('Unexpected error while opening the PDF. Please retry.');
    }
  }

  static String _friendlyDioError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
        return 'The download timed out. Check your connection and retry.';
      case DioExceptionType.badResponse:
        final code = e.response?.statusCode;
        if (code == 401 || code == 403) {
          return 'Your session may have expired. Please sign in again and retry.';
        }
        return 'The server could not provide the PDF (status $code).';
      case DioExceptionType.connectionError:
        return 'No internet connection. Please retry when you are back online.';
      default:
        return 'Could not download the PDF. Please retry.';
    }
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
}
