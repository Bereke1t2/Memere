# phase6/SKILL_2.md — Course Checkout, WebView Flow & Enrollment Confirmation
# Memere Mobile (memere_mobile) — Phase 6, Part 2
# READ SKILL.md → phase6/SKILL_1.md → then this file.

---

## OBJECTIVE

Build the student checkout experience:
course detail CTA → free enrollment or paid checkout → provider selection →
idempotent payment initiate → provider WebView → payment status polling →
enrollment confirmation.

By the end of this skill, a student can gain course access through free
enrollment or a backend-verified paid checkout.

---

## PAYMENT FLOW RULES

- For free courses, call `POST /courses/:id/enroll-free`.
- For paid courses, call `POST /payments/initiate`.
- Generate one idempotency key per checkout attempt and reuse it for retries.
- Open only backend-provided `redirect_url` in WebView.
- Poll backend payment status after WebView redirect/close.
- Treat `completed` status as payment success.
- Confirm course access by refreshing enrollments or course access state.
- Do not treat WebView URL alone as success.

---

## PART A — COURSE DETAIL CTA UPDATE

### FILE A1 — `lib/features/courses/presentation/screens/course_detail_screen.dart`

Replace Phase 2 placeholder CTA.

Rules:

- if course is free: button `Start learning`
- if course is paid and not enrolled: button `Enroll for <price>`
- if enrolled: button `Continue learning`
- if subscription active later grants access: button `Continue learning`

For now, enrollment state comes from `listEnrollments`.

---

## PART B — COURSE ACCESS PROVIDER

### FILE B1 — `lib/features/payment/presentation/providers/course_access_provider.dart`

Create provider family:

```dart
final courseAccessProvider =
    FutureProvider.family<CourseAccessState, String>((ref, courseId) async {
  // load enrollments and decide if courseId exists and is active
});
```

State:

```dart
class CourseAccessState {
  const CourseAccessState({
    required this.hasAccess,
    this.enrollment,
  });
}
```

Use:

- `ListEnrollmentsUseCase`
- optional subscription provider when implemented

Refresh this provider after free enroll or completed payment.

---

## PART C — CHECKOUT FLOW PROVIDER

### FILE C1 — `lib/features/payment/presentation/providers/checkout_flow_provider.dart`

Create an `AsyncNotifierProvider.family` keyed by course ID.

State:

```dart
class CheckoutFlowState {
  const CheckoutFlowState({
    this.selectedProvider,
    this.idempotencyKey,
    this.initiation,
    this.payment,
    this.isPolling = false,
  });
}
```

Methods:

- `selectProvider(PaymentProvider provider)`
- `startFreeEnrollment(courseId)`
- `startPaidCheckout(courseId, provider, couponCode)`
- `pollStatus(paymentId)`
- `reset()`

Idempotency behavior:

- generate key before first initiate
- keep same key for retry of same checkout
- reset key only when user intentionally starts a new checkout

Use `uuid` package already in `pubspec.yaml`.

---

## PART D — PROVIDER SELECTION

### FILE D1 — `lib/features/payment/presentation/widgets/payment_provider_sheet.dart`

Show as bottom sheet.

Providers:

- Chapa
- Telebirr
- Stripe
- Mock only in debug/dev builds

Each option shows:

- icon or initials
- provider name
- short payment method description

Actions:

- cancel
- continue

Do not show card fields.

---

## PART E — CHECKOUT WEBVIEW SCREEN

### FILE E1 — `lib/features/payment/presentation/screens/payment_webview_screen.dart`

Constructor:

```dart
class PaymentWebViewScreen extends ConsumerWidget {
  const PaymentWebViewScreen({
    super.key,
    required this.paymentId,
    required this.redirectUrl,
    required this.courseId,
  });

  final String paymentId;
  final String redirectUrl;
  final String courseId;
}
```

Route:

```dart
static const paymentWebView = '/payments/:paymentId/webview';

static String paymentWebViewPath({
  required String paymentId,
  required String courseId,
  required String redirectUrl,
}) {
  return '/payments/$paymentId/webview?courseId=$courseId&redirectUrl=${Uri.encodeComponent(redirectUrl)}';
}
```

Use `webview_flutter`.

Required UI:

- app bar title `Complete payment`
- secure checkout WebView
- loading progress indicator
- close button with confirmation

Navigation delegate:

- watch URL changes for backend return URL patterns if configured
- when return/cancel-like URL is detected, stop loading and poll status
- if user closes manually, poll status before deciding final state

Do not parse provider secrets from URL.

---

## PART F — PAYMENT STATUS POLLING

### FILE F1 — `lib/features/payment/presentation/providers/payment_status_polling_provider.dart`

Polling rules:

- poll `GET /payments/:id/status`
- interval: 2-3 seconds
- timeout: 2 minutes
- stop when status is `completed`, `failed`, or `refunded`
- on completed: refresh enrollments and course access
- on pending timeout: show pending state with retry/check again

Do not hammer backend.

---

## PART G — RESULT/CONFIRMATION SCREEN

### FILE G1 — `lib/features/payment/presentation/screens/payment_result_screen.dart`

Route:

```dart
static const paymentResult = '/payments/:paymentId/result';
```

Show based on status:

- completed: success, `Course unlocked`
- failed: failure, retry
- pending: pending, check again

Actions:

- go to course
- view purchases

For completed status, navigate to course detail or first lesson only after access
provider refresh confirms enrollment.

---

## PART H — FREE ENROLLMENT FLOW

For free course CTA:

1. call `EnrollFreeUseCase`
2. refresh `courseAccessProvider(courseId)`
3. show success snackbar
4. update CTA to `Continue learning`

If backend rejects because course is paid:

- show paid provider sheet

If already enrolled:

- refresh access state and show `Continue learning`

---

## PART I — PAID CHECKOUT FLOW

For paid course CTA:

1. open provider sheet
2. generate idempotency key
3. call `InitiateCoursePaymentUseCase`
4. navigate to WebView with `redirectUrl`
5. poll status when WebView returns/closes
6. show result screen
7. refresh access state

If initiate returns:

- `COURSE_IS_FREE`: run free enrollment flow
- `ALREADY_ENROLLED`: refresh access and show continue
- `IDEMPOTENCY_KEY_REQUIRED`: bug in mobile, fix immediately

---

## PART J — ROUTER UPDATE

### FILE J1 — `lib/core/router/app_router.dart`

Add:

- payment WebView route
- payment result route
- purchase history route from `SKILL_3.md`

Keep all payment routes behind auth.

---

## PART K — WEBVIEW SAFETY

Rules:

- only open `https://` provider URLs, except local/mock dev URLs when explicitly enabled
- show error for invalid URL
- never inject JavaScript that reads payment form data
- clear loading/error state when page fails
- user can close WebView, but app must poll backend after close

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/payment lib/features/courses
flutter analyze
flutter build apk --debug
```

Manual test:

- free course enrolls without WebView
- paid course opens provider sheet
- initiate sends `Idempotency-Key`
- WebView opens `redirect_url`
- closing/returning polls backend status
- completed payment refreshes enrollment/access
- failed/pending states are clear

---

## SKILL_2 CHECKLIST

- [ ] Course detail CTA uses real free/paid/enrolled state
- [ ] Course access provider reads enrollments
- [ ] Provider selection sheet works
- [ ] Paid initiate sends stable idempotency key
- [ ] WebView opens backend redirect URL
- [ ] WebView close/return polls backend status
- [ ] Payment result screen handles completed/pending/failed
- [ ] Completed payment refreshes access
- [ ] Free enrollment refreshes access
- [ ] No card/payment details are collected in app
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds
