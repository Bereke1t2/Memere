# fix1.md — Fix: downloaded PDF shows lesson notes instead of the real PDF

> **Execution notes for the implementing agent (Gemini).** This fix now spans **two
> repos**. Do the parts in order — Part 1 is what actually makes real PDFs appear.
> **No visual/UI/layout changes anywhere.** No encryption work.
> - **Part 1 — Mobile** (`memere_mobile/`): `secure_pdf_storage.dart` + `pdf_reader_screen.dart`.
> - **Part 2 — Backend** (`Memere-backend/`): `internal/delivery/http/course_handler.go`.
>
> After editing mobile, run `flutter analyze`. After editing backend, run `go build ./...`.

---

## ⚠️ Prerequisite (config, not code) — MinIO must be reachable from the backend

Symptom that points here: the backend logs the PDF request as **`status=302`** with
**`latency≈25s`**. That means `store.Get()` (S3 `GetObject`) is **timing out** trying to
reach MinIO, then falling back to a presigned redirect the device also can't reach.

`Memere-backend/.env` has `S3_ENDPOINT=http://192.168.0.201:9000` — a LAN IP that goes
stale when the machine's address changes. Everything runs locally (the app reaches the
API via `127.0.0.1`), so set it to a reachable address and **restart the backend**:
```
S3_ENDPOINT=http://127.0.0.1:9000
```
Verify first with `curl -s -o /dev/null -w "%{http_code}\n" --connect-timeout 3 http://127.0.0.1:9000/minio/health/live`
(expect `200`; start MinIO with `docker compose up -d` if nothing responds). A healthy
`GetObject` returns in milliseconds — never 25 s. If, after fixing this, `GetObject` is
fast but the object 404s, the earlier upload never landed → **re-upload the PDF from the
admin page**.

---

## Problem

Opening a lesson's **PDF Document** tab shows *study-notes text rendered as a PDF*
instead of the actual uploaded PDF — even though the PDF was uploaded correctly from
the admin panel.

## What actually happens (verified end-to-end)

The upload side is **correct**: the admin panel sends a real `application/pdf`, the
backend validates the `%PDF-` header (`course_handler.go:380`), stores it in MinIO/S3
at `lessons/{id}/notes.pdf` (`:385-391`), and saves that **storage key** into
`lessons.pdf_url` (`:393`). So the stored object is a genuine PDF.

The **download side is broken by two independent bugs**, and a third masking behavior
that made the failure look like "notes in a PDF":

### 🔴 Bug X (PRIMARY, mobile) — the request 404s because `/api/v1` is doubled
The whole app uses a Dio client whose `baseUrl` **already contains `/api/v1`**
(`dio_client.dart:24` → `http://<host>:8080/api/v1`), so every datasource uses
relative paths like `'/courses'` (`courses_remote_datasource.dart:28,46`).

But `_resolveDownloadUrl` (`secure_pdf_storage.dart:83-86`) builds the PDF URL with a
**second** `/api/v1`:
```dart
final apiBase = fixMediaUrl(Env.baseUrl);          // = http://<host>:8080/api/v1
return '$apiBase/api/v1/lessons/$lessonId/pdf';     // => .../api/v1/api/v1/lessons/{id}/pdf
```
The real route is registered at `/api/v1/lessons/:id/pdf`, so the doubled path
**404s for every lesson**. The download never succeeds → (old code) falls back to a
generated notes-PDF. **This is the root reason real PDFs never show.**

### 🟠 Bug Y (SECONDARY, backend) — the download response never signals a clean end-of-body
The backend must return the PDF in a way the mobile client can finish reading. Two ways
it currently gets this wrong:
- **302 redirect to a MinIO presigned URL** (`http://192.168.0.201:9000/...` or
  `http://127.0.0.1:9000/...`): on a phone/emulator that host is often unreachable, is
  cleartext `http` (Android blocks it by default), and the SigV4 signature covers the
  Host header — so if the app rewrites the host the signature breaks (`403`).
- **Streaming with `io.Copy` and no `Content-Length`** (chunked/unbounded): the Dio
  client (`ResponseType.bytes`) has no definite body length to wait for, so after the
  bytes arrive the download future **never resolves** — the reader hangs on the loading
  overlay at ~100%. ← **this is the "stuck at 100%" symptom.**

Fix: the backend should return a **bounded response with an explicit `Content-Length`**
(read the object and use `c.Data`), served from the API host the app already reaches.

### 🟡 Masking behavior (mobile) — silent fake-PDF fallback
`downloadPdf` currently **always returns a PDF**: on any failure it *generates* a PDF
from the lesson's notes text (`_generatePdfFromLessonContent`) and caches it. That
fake file passes `isDownloaded()` (which only checks for a `%PDF-` header), so it's
never re-fetched. This is why the failure looked like "the notes rendered as a PDF"
rather than an error. We remove this and signal the true outcome instead.

> Note: there is **no encryption** anywhere (despite the class name `SecurePdfStorage`);
> files are stored as plain bytes in the app-private sandbox. Out of scope here.

## Intended behavior after the fix
- Lesson has a real PDF → **the real PDF displays** (Bug X + Bug Y fixed).
- No real PDF attached, or the URL returns non-PDF content → **route to the Study Notes tab** (no fabricated PDF).
- Real PDF expected but download fails (network/auth/timeout/server) → show the **existing** error + retry view.

The reader UI already supports all three outcomes (`_buildErrorOrEmptyView`, the
Study Notes tab, the "Read Notes"/"Load Document" buttons), so **no widgets change**.

---

# PART 1 — MOBILE (`memere_mobile/`)

## File 1 — `lib/core/storage/secure_pdf_storage.dart`

### 1a. 🔴 Fix the doubled `/api/v1` in `_resolveDownloadUrl` (currently lines 83-86) — THE PRIMARY FIX
`Env.baseUrl` already ends in `/api/v1`, so **remove the extra `/api/v1`**:
```dart
// It's an S3 key like "lessons/xxx/notes.pdf" or a raw filename
if (lessonId != null && lessonId.isNotEmpty) {
  final apiBase = fixMediaUrl(Env.baseUrl); // already ends in /api/v1
  return '$apiBase/lessons/$lessonId/pdf';
}
```
(Only the URL string changes — do not touch the rest of the method.)

### 1b. Add two exception types (after the imports, before `class SecurePdfStorage`)
```dart
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
```

### 1c. Bump the cache-key version in `getFileKey` (currently lines 17-24)
This flushes stale fake note-PDFs exactly once, forcing a clean re-download.
- `'pdf_${...}'` → `'pdf_v2_${...}'`
- `'pdf_note_${...}'` → `'pdf_v2_note_${...}'`

The reader derives its key from the same method, so this stays consistent.

### 1d. Rewrite `downloadPdf` (currently lines 307-397)
Remove the notes-PDF fallback entirely; throw typed outcomes; follow redirects in a loop.
Replace the whole method with:
```dart
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
```

Add this helper anywhere in the class:
```dart
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
```

### 1e. Delete now-dead code
- Delete `_generatePdfFromLessonContent` (currently lines 447-525).
- Delete `_sanitizeTextForPdfStandardFont` (currently lines 415-444).
- Remove the import `import 'package:syncfusion_flutter_pdf/pdf.dart';` (currently line 6).
- **Keep** `getEffectiveContent` — it is still used by the Study Notes tab.
- Optional: if `grep -rn "syncfusion_flutter_pdf" memere_mobile/lib` shows no other
  use, remove `syncfusion_flutter_pdf: ^33.2.15` from `pubspec.yaml`. Do **not** touch
  `syncfusion_flutter_pdfviewer` or `flutter_pdfview`.

## File 2 — `lib/features/courses/presentation/screens/pdf_reader_screen.dart`

`secure_pdf_storage.dart` is already imported (line 9), so the new exceptions are in scope.

### 2a. Add a state field (near line 40, with the other `bool` fields)
```dart
bool _noPdfAvailable = false;
```

### 2b. Reset the flag when a (re)download starts
Inside `_downloadPdf`'s opening `setState(...)` (currently lines 106-110), add:
```dart
_noPdfAvailable = false;
```
so the "Load Document" / "Re-render" retry buttons truly retry.

### 2c. Replace the single `catch` in `_downloadPdf` (currently lines 135-141) with typed handling
```dart
} on PdfNotAvailableException {
  if (!mounted) return;
  setState(() {
    _isDownloading = false;
    _noPdfAvailable = true;
    _activeTab = 0; // route to Study Notes
  });
} on PdfDownloadException catch (e) {
  if (!mounted) return;
  setState(() {
    _isDownloading = false;
    _pdfRenderError = e.message; // shown by _buildErrorOrEmptyView with retry
  });
} catch (e) {
  if (!mounted) return;
  setState(() {
    _isDownloading = false;
    _pdfRenderError = e.toString();
  });
}
```

### 2d. Don't auto-retry for known no-PDF lessons — PDF tab `onTap` (currently lines 482-487)
Add `&& !_noPdfAvailable` to the guard:
```dart
onTap: () {
  setState(() => _activeTab = 1);
  if (_localPdfPath == null && !_isDownloading && !_noPdfAvailable) {
    _downloadPdf();
  }
},
```
The empty view's "Load Document" button still calls `_downloadPdf` (which resets the
flag), so a manual retry always works.

---

# PART 2 — BACKEND (`Memere-backend/`)

## File 3 — `internal/delivery/http/course_handler.go`

### 3a. 🟠 Stream the PDF bytes instead of 302-redirecting to a presigned URL
In `DownloadLessonPDF`, the object-store branch currently tries `PresignGet` **first**
and only streams as a fallback (currently lines 430-445). Reverse it: **stream first**
(reliable from any device because it reuses the API connection), and keep the presigned
redirect only as a fallback (useful in prod where the CDN/S3 host is publicly reachable).

Replace the `if h.store != nil { ... }` block with:
```go
if h.store != nil {
	// Read the object fully and return it with an explicit Content-Length via
	// c.Data. This is the critical detail: io.Copy(c.Writer, rc) sends the body
	// with chunked/unbounded transfer encoding (no Content-Length), and the
	// mobile Dio client (ResponseType.bytes) then has no definite end-of-body
	// to wait for -> its download future never completes -> the reader hangs at
	// ~100% on the loading overlay. c.Data sets Content-Length so the client
	// knows exactly when the body ends and completes cleanly.
	rc, err := h.store.Get(c.Request.Context(), pdfPath)
	if err == nil {
		defer rc.Close()
		data, rerr := io.ReadAll(rc)
		if rerr == nil && len(data) > 0 {
			c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filepath.Base(pdfPath)))
			c.Data(http.StatusOK, "application/pdf", data)
			return
		}
	}

	// Fallback (mainly prod): hand back a presigned URL when the object can't be
	// read here. Reached only on a read error, so it won't mask the dev path.
	if presigned, perr := h.store.PresignGet(c.Request.Context(), pdfPath, 1*time.Hour); perr == nil && presigned != "" {
		c.Redirect(http.StatusFound, presigned)
		return
	}
}
```
Leave everything else in the handler unchanged (the empty-`PdfURL` 404 at the top, and
the absolute-URL redirect for `http(s)://` values, both stay as-is). No imports change —
`io`, `filepath`, `time`, `net/http` are already imported.

> **Why not `io.Copy`?** `c.Data` buffers the whole PDF in memory (fine — uploads are
> capped at 50 MB) and, crucially, sets `Content-Length`. Streaming with `io.Copy` omits
> `Content-Length` → chunked response → the mobile download future never resolves →
> **the "stuck at 100%" hang you are seeing.**

---

## Do NOT change
- No visual/layout/widget/color/copy changes (beyond the exception message strings in Part 1).
- No encryption; keep app-private sandbox storage as-is.
- Do not touch the success path, `_isValidPdfBytes`, `getEffectiveContent`, the Study
  Notes rendering, the upload handler, or storage config.

## Edge case (document, don't over-engineer)
If the backend ever returns an error page as **HTTP 200 HTML**, the mobile side
classifies it as "not available" and routes to notes rather than showing a retry. This
is rare (real errors are non-200 → retry path) and is recoverable via "Load Document".

## Verification
1. **Backend:** `cd Memere-backend && go build ./...` — must pass. Then
   `curl -i -H "Authorization: Bearer <token>" http://<host>:8080/api/v1/lessons/<id>/pdf`
   must return **`HTTP/1.1 200`**, header **`Content-Length: <n>`** (a real number),
   `Content-Type: application/pdf`, and a body starting with `%PDF-`. It must **NOT**
   show `Transfer-Encoding: chunked` and must **NOT** be a `302` redirect — either of
   those reproduces the "stuck at 100%" hang.
2. **Mobile:** `cd memere_mobile && flutter analyze` — must pass (confirms deletions left no references).
3. `flutter run` and check:
   - Lesson with a real uploaded PDF → PDF tab shows the **actual PDF** (not notes). ← the core fix
   - Lesson with no PDF (`pdfUrl` empty / `sample.pdf`) → lands on Study Notes; PDF tab shows the empty view with "Read Notes"; no fabricated PDF.
   - Lesson whose `pdfUrl` is really a `note_url` (non-PDF) → routes to Study Notes.
   - Real PDF with network off / bad token → PDF tab shows error + "Load Document" retry; retry works when back online.
   - Previously cached fake note-PDF → after the `getFileKey` version bump, the real PDF re-downloads on next open.
