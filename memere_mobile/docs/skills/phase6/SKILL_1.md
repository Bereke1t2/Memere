# phase6/SKILL_1.md — Payment, Enrollment & Subscription Data Layer
# Memere Mobile (memere_mobile) — Phase 6, Part 1
# READ SKILL.md → all phase1 skill files → all phase2 skill files → all phase3 skill files → all phase4 skill files → all phase5 skill files → then this file.

---

## OBJECTIVE

Build the Payment & Enrollment foundation:
payment entities → enrollment entities → subscription plan entities →
repository contract → use cases → backend DTO models → remote datasource →
repository implementation → Riverpod dependency providers.

By the end of this skill, the app can enroll in free courses, initiate paid
course checkout with an idempotency key, poll payment status, list enrollments,
and list payment history. WebView checkout is built in `SKILL_2.md`.

---

## PHASE 6 BOUNDARY

Phase 6 builds **course access and payment flows**:

- free course enrollment
- paid course checkout
- Chapa/Telebirr/Stripe/mock provider selection
- WebView checkout
- payment status polling
- enrollment confirmation
- purchase history
- subscription plan catalog and subscription checkout foundation

Do **not** build these yet:

- refunds UI
- coupon management UI
- teacher earnings dashboard
- admin revenue dashboard
- certificates
- notification center

---

## SECURITY AND PAYMENT RULES

These are non-negotiable:

- Every payment initiation must send an `Idempotency-Key` header.
- Never retry a payment initiation with a new key for the same user tap.
- Never log payment redirect URLs, provider refs, tokens, or card data.
- Never collect card details in-app; provider checkout handles that.
- Payment completion must be verified by backend status polling, not WebView URL alone.
- Enrollment is confirmed only after backend says payment is `completed` or free enroll succeeds.
- Webhook routes are backend/provider only; mobile must not call them.

---

## BACKEND CONTRACT

Student routes:

```text
POST /api/v1/courses/:id/enroll-free
GET /api/v1/me/enrollments?limit=50
POST /api/v1/payments/initiate
GET /api/v1/payments?limit=50
GET /api/v1/payments/:id/status
GET /api/v1/subscription-plans
POST /api/v1/subscriptions
GET /api/v1/me/subscription
POST /api/v1/subscriptions/:id/cancel
```

Do not call:

```text
POST /api/v1/webhooks/payments/:provider
POST /api/v1/payments/:id/refund
```

Those are provider/admin paths.

---

## REQUEST/RESPONSE SHAPES

### `POST /courses/:id/enroll-free`

Response: `204 No Content`

### `GET /me/enrollments`

```json
{
  "data": [
    {
      "id": "uuid",
      "course_id": "uuid",
      "source": "free",
      "enrolled_at": "2026-01-01T00:00:00Z",
      "expires_at": "2026-02-01T00:00:00Z"
    }
  ]
}
```

### `POST /payments/initiate`

Headers:

```text
Idempotency-Key: <uuid-v4>
```

Body:

```json
{
  "course_id": "uuid",
  "provider": "chapa",
  "coupon_code": "OPTIONAL"
}
```

Response:

```json
{
  "payment_id": "uuid",
  "redirect_url": "https://provider-checkout",
  "amount": "250.00",
  "currency": "ETB"
}
```

### `GET /payments/:id/status`

```json
{
  "payment_id": "uuid",
  "status": "pending",
  "amount": "250.00",
  "currency": "ETB",
  "course_id": "uuid"
}
```

### `GET /payments`

```json
{
  "data": [
    {
      "payment_id": "uuid",
      "status": "completed",
      "amount": "250.00",
      "currency": "ETB",
      "course_id": "uuid"
    }
  ]
}
```

### `GET /subscription-plans`

```json
{
  "data": [
    {
      "plan": "monthly",
      "price": "300.00",
      "currency": "ETB",
      "period_days": 30
    }
  ]
}
```

### `POST /subscriptions`

Headers:

```text
Idempotency-Key: <uuid-v4>
```

Body:

```json
{
  "plan": "monthly",
  "provider": "chapa",
  "coupon_code": "OPTIONAL"
}
```

Response: same as payment initiate response.

---

## REQUIRED FOLDER STRUCTURE

Create or verify:

```bash
mkdir -p lib/features/payment/data/datasources
mkdir -p lib/features/payment/data/models
mkdir -p lib/features/payment/data/repositories
mkdir -p lib/features/payment/domain/entities
mkdir -p lib/features/payment/domain/repositories
mkdir -p lib/features/payment/domain/usecases
mkdir -p lib/features/payment/presentation/providers
mkdir -p lib/features/payment/presentation/screens
mkdir -p lib/features/payment/presentation/widgets
```

---

## PART A — DOMAIN ENTITIES

### FILE A1 — `lib/features/payment/domain/entities/payment_entity.dart`

Fields:

- `paymentId`
- `status`
- `amount`
- `currency`
- `courseId`

Use enum:

```dart
enum PaymentStatus { pending, completed, failed, refunded }
```

Helpers:

- `bool get isPending`
- `bool get isCompleted`
- `bool get isFailed`
- `String get amountLabel`

---

### FILE A2 — `lib/features/payment/domain/entities/payment_initiation_entity.dart`

Fields:

- `paymentId`
- `redirectUrl`
- `amount`
- `currency`

---

### FILE A3 — `lib/features/payment/domain/entities/enrollment_entity.dart`

Fields:

- `id`
- `courseId`
- `source`
- `enrolledAt`
- `expiresAt`

Use enum:

```dart
enum EnrollmentSource { purchase, subscription, free, coupon }
```

Helper:

- `bool get isActive`

---

### FILE A4 — `lib/features/payment/domain/entities/payment_provider_entity.dart`

Use enum:

```dart
enum PaymentProvider { chapa, telebirr, stripe, mock }
```

Rules:

- show `mock` only in debug/dev mode
- default Ethiopia provider order: Chapa, Telebirr, Stripe

---

### FILE A5 — `lib/features/payment/domain/entities/subscription_plan_entity.dart`

Fields:

- `plan`
- `price`
- `currency`
- `periodDays`

Use enum:

```dart
enum SubscriptionPlanType { monthly, annual }
```

---

### FILE A6 — `lib/features/payment/domain/entities/subscription_entity.dart`

Fields:

- `id`
- `plan`
- `status`
- `currentPeriodStart`
- `currentPeriodEnd`
- `canceledAt`

Use enum:

```dart
enum SubscriptionStatus { active, pastDue, canceled, expired }
```

---

## PART B — REPOSITORY CONTRACT

### FILE B1 — `lib/features/payment/domain/repositories/payment_repository.dart`

Use `Either<Failure, T>`.

Methods:

```dart
Future<Either<Failure, void>> enrollFree(String courseId);

Future<Either<Failure, List<EnrollmentEntity>>> listEnrollments({
  int limit = 50,
});

Future<Either<Failure, PaymentInitiationEntity>> initiateCoursePayment({
  required String courseId,
  required PaymentProvider provider,
  required String idempotencyKey,
  String? couponCode,
});

Future<Either<Failure, PaymentEntity>> getPaymentStatus(String paymentId);

Future<Either<Failure, List<PaymentEntity>>> listPayments({
  int limit = 50,
});

Future<Either<Failure, List<SubscriptionPlanEntity>>> listSubscriptionPlans();

Future<Either<Failure, PaymentInitiationEntity>> initiateSubscriptionPayment({
  required SubscriptionPlanType plan,
  required PaymentProvider provider,
  required String idempotencyKey,
  String? couponCode,
});

Future<Either<Failure, SubscriptionEntity>> getMySubscription();

Future<Either<Failure, void>> cancelSubscription(String subscriptionId);
```

---

## PART C — USE CASES

Create:

```text
enroll_free_usecase.dart
list_enrollments_usecase.dart
initiate_course_payment_usecase.dart
get_payment_status_usecase.dart
list_payments_usecase.dart
list_subscription_plans_usecase.dart
initiate_subscription_payment_usecase.dart
get_my_subscription_usecase.dart
cancel_subscription_usecase.dart
```

Rules:

- validate IDs are not empty
- validate idempotency key is not empty for initiate calls
- validate provider is supported
- validate plan for subscription initiate
- return `ValidationFailure` for invalid input

---

## PART D — DATA MODELS

Create:

```text
payment_model.dart
payment_initiation_model.dart
enrollment_model.dart
subscription_plan_model.dart
subscription_model.dart
```

Parsing rules:

- money values come from backend as strings; keep as `String` or parse with care
- parse date strings with `DateTime.parse`
- fallback unknown payment status to `pending`
- fallback unknown enrollment source to `purchase`
- fallback unknown subscription status to `expired`
- fallback unknown provider only in request validation, not response parsing

---

## PART E — REMOTE DATASOURCE

### FILE E0 — `lib/core/network/dio_client.dart`

If not already supported, extend request helpers to accept Dio `Options`:

```dart
Future<Response<T>> post<T>(
  String path, {
  dynamic data,
  Map<String, dynamic>? queryParameters,
  Options? options,
}) =>
    _dio.post<T>(
      path,
      data: data,
      queryParameters: queryParameters,
      options: options,
    );
```

Do the same for `get`, `put`, `patch`, and `delete` only if needed. Keep the
existing call sites working.

### FILE E1 — `lib/features/payment/data/datasources/payment_remote_datasource.dart`

Use `DioClient` only.

Methods:

```dart
Future<void> enrollFree(String courseId);
Future<List<EnrollmentModel>> listEnrollments({int limit = 50});
Future<PaymentInitiationModel> initiateCoursePayment({
  required String courseId,
  required String provider,
  required String idempotencyKey,
  String? couponCode,
});
Future<PaymentModel> getPaymentStatus(String paymentId);
Future<List<PaymentModel>> listPayments({int limit = 50});
Future<List<SubscriptionPlanModel>> listSubscriptionPlans();
Future<PaymentInitiationModel> initiateSubscriptionPayment({
  required String plan,
  required String provider,
  required String idempotencyKey,
  String? couponCode,
});
Future<SubscriptionModel> getMySubscription();
Future<void> cancelSubscription(String subscriptionId);
```

Endpoint examples:

```dart
await _client.post('/courses/$courseId/enroll-free');
await _client.get<Map<String, dynamic>>('/me/enrollments', queryParameters: {'limit': limit});
await _client.post<Map<String, dynamic>>(
  '/payments/initiate',
  data: {
    'course_id': courseId,
    'provider': provider,
    if (couponCode != null && couponCode.isNotEmpty) 'coupon_code': couponCode,
  },
  options: Options(headers: {'Idempotency-Key': idempotencyKey}),
);
await _client.get<Map<String, dynamic>>('/payments/$paymentId/status');
await _client.get<Map<String, dynamic>>('/payments', queryParameters: {'limit': limit});
await _client.get<Map<String, dynamic>>('/subscription-plans');
await _client.post<Map<String, dynamic>>(
  '/subscriptions',
  data: {
    'plan': plan,
    'provider': provider,
    if (couponCode != null && couponCode.isNotEmpty) 'coupon_code': couponCode,
  },
  options: Options(headers: {'Idempotency-Key': idempotencyKey}),
);
await _client.get<Map<String, dynamic>>('/me/subscription');
await _client.post('/subscriptions/$subscriptionId/cancel');
```

Do not create raw Dio in the datasource.

---

## PART F — REPOSITORY IMPLEMENTATION

### FILE F1 — `lib/features/payment/data/repositories/payment_repository_impl.dart`

Map errors:

- `DioException` → `ServerFailure.fromDioError`
- any other exception → `UnknownFailure`

Do not create idempotency keys in the repository. Create them in the payment
flow provider so retries for the same user action reuse the same key.

---

## PART G — PROVIDERS

### FILE G1 — `lib/features/payment/presentation/providers/payment_providers.dart`

Create:

- `paymentRemoteDataSourceProvider`
- `paymentRepositoryProvider`
- all use case providers

Mirror auth/courses/video/quiz/exam provider style.

---

## VALIDATION

After this skill:

```bash
dart format lib/features/payment
flutter analyze
```

Do not continue to checkout UI until this passes.

---

## SKILL_1 CHECKLIST

- [ ] Payment/enrollment/subscription entities compile
- [ ] Repository contract covers free enroll, payment, history, subscription methods
- [ ] Use cases validate IDs, provider, plan, and idempotency key
- [ ] Models parse backend JSON exactly
- [ ] Remote datasource sends `Idempotency-Key` header on initiate calls
- [ ] DioClient supports custom request options if needed
- [ ] Repository maps failures consistently
- [ ] Dependency providers are wired
- [ ] `flutter analyze` has 0 errors
