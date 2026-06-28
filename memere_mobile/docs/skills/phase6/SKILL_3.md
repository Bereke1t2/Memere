# phase6/SKILL_3.md — Purchase History, Subscriptions & Phase 6 Checklist
# Memere Mobile (memere_mobile) — Phase 6, Part 3
# READ SKILL.md → phase6/SKILL_1.md → phase6/SKILL_2.md → then this file.

---

## OBJECTIVE

Finish Phase 6 by adding:
purchase history → enrollment list → subscription plan catalog →
subscription checkout foundation → access-state polish → final validation.

By the end of this skill, students can review purchases/enrollments and start
subscription checkout using the same idempotent payment flow.

---

## PHASE 6 FINAL BOUNDARY

Included:

- purchase history
- enrollment list
- active subscription view
- subscription plan catalog
- subscription checkout WebView flow
- cancel subscription action

Deferred:

- subscription-only content personalization
- detailed invoices/receipts PDF
- refunds
- coupon browsing/management
- teacher earnings
- admin revenue

---

## PART A — PURCHASE HISTORY SCREEN

### FILE A1 — `lib/features/payment/presentation/screens/purchase_history_screen.dart`

Route:

```dart
static const purchaseHistory = '/payments';
```

Required UI:

- top app bar: `Purchases`
- payment history tab/list
- enrollments tab/list
- loading skeleton
- empty state
- error retry

Payment row shows:

- amount/currency
- status
- course ID or course title if locally resolvable
- status badge

Enrollment row shows:

- course ID/title
- source
- enrolled date
- expiry if present

---

## PART B — HISTORY PROVIDERS

### FILE B1 — `lib/features/payment/presentation/providers/purchase_history_provider.dart`

Create:

- `paymentHistoryProvider`
- `enrollmentListProvider`

Methods:

- refresh payments
- refresh enrollments
- retry failed load

Use:

- `ListPaymentsUseCase`
- `ListEnrollmentsUseCase`

---

## PART C — HISTORY WIDGETS

### FILE C1 — `lib/features/payment/presentation/widgets/payment_history_tile.dart`

Show:

- payment status badge
- amount
- course label
- provider if available later

Status colors:

- pending: warning
- completed: success
- failed: error
- refunded: info/secondary

---

### FILE C2 — `lib/features/payment/presentation/widgets/enrollment_tile.dart`

Show:

- course label
- source badge
- enrolled date
- expiry

Source colors:

- purchase: accent
- free: success
- subscription: info
- coupon: warning

---

### FILE C3 — `lib/features/payment/presentation/widgets/payment_empty_state.dart`

Reusable empty/error state.

---

## PART D — SUBSCRIPTION PLAN CATALOG

### FILE D1 — `lib/features/payment/presentation/screens/subscription_plans_screen.dart`

Route:

```dart
static const subscriptionPlans = '/subscription-plans';
```

Required UI:

- title: `All-access plans`
- monthly/annual plan cards
- price
- period
- provider selection
- subscribe button
- active subscription state

Do not over-market. Keep it clear and transactional.

---

### FILE D2 — `lib/features/payment/presentation/providers/subscription_provider.dart`

Create:

- `subscriptionPlansProvider`
- `mySubscriptionProvider`
- subscription checkout flow provider or reuse `checkout_flow_provider`

Actions:

- load plans
- load active subscription
- initiate subscription payment
- cancel subscription

---

### FILE D3 — `lib/features/payment/presentation/widgets/subscription_plan_card.dart`

Show:

- plan
- price
- period days
- selected state
- subscribe button

Annual plan may show simple value text only if computed safely.

---

### FILE D4 — `lib/features/payment/presentation/widgets/active_subscription_card.dart`

Show:

- plan
- status
- current period start/end
- canceled date if present
- cancel button when active

Cancel action:

- show confirmation
- call `CancelSubscriptionUseCase`
- refresh active subscription

---

## PART E — SUBSCRIPTION CHECKOUT FLOW

Use the same payment WebView and status polling patterns from `SKILL_2.md`.

Differences:

- no `courseId`
- body uses `plan`
- endpoint is `POST /subscriptions`
- success refreshes active subscription

Idempotency:

- generate one key per subscription checkout attempt
- reuse key for retry of same attempt

After completed status:

- refresh `mySubscriptionProvider`
- show success result
- optionally navigate to course catalog

---

## PART F — COURSE ACCESS POLISH

Update access logic:

- access if active enrollment exists
- access if active subscription exists
- free courses can enroll directly

Refresh access after:

- free enroll
- completed course payment
- completed subscription payment
- subscription cancel

Canceling a subscription should keep access until `current_period_end` if backend
still reports active access. Do not remove access locally just because cancel was tapped.

---

## PART G — COURSE DETAIL PAYMENT STATES

Course detail CTA states:

- loading access state
- free and not enrolled: `Start learning`
- paid and no access: `Enroll for <price>`
- pending payment: `Payment pending`
- access granted: `Continue learning`
- subscription active: `Continue learning`

Pending payment state can be read from the latest payment history for that course.

---

## PART H — ERROR HANDLING

Common errors:

- `ALREADY_ENROLLED`: refresh access
- `COURSE_IS_FREE`: use free enrollment
- `IDEMPOTENCY_KEY_REQUIRED`: mobile bug, fix code
- `UNKNOWN_PROVIDER`: hide unsupported provider
- payment pending timeout: show check-again action
- WebView load failed: show retry/open externally option only if safe

Do not swallow payment failures silently.

---

## PART I — LOCAL DEV PROVIDER

If backend has mock provider enabled, allow `PaymentProvider.mock` in debug mode.

Rules:

- hide mock in release builds
- label it clearly as `Mock`
- use only for local testing

---

## PART J — TESTING NOTES

Widget tests:

- free CTA calls free enroll
- paid CTA opens provider sheet
- completed payment updates CTA
- history empty state renders
- subscription plan cards render

Unit tests:

- payment status parsing
- enrollment active helper
- idempotency key reuse in checkout flow
- provider visibility for debug/release

---

## VALIDATION

Run:

```bash
dart format lib/core/router lib/features/payment lib/features/courses
flutter analyze
flutter build apk --debug
```

Manual backend run:

```bash
cd ../Memere-backend
make up
make migrate-up
make seed
make run
```

Physical Android:

```bash
adb reverse tcp:8080 tcp:8080
flutter run
```

Manual checks:

- free course enrolls
- paid course initiates checkout
- WebView opens
- status polling completes or shows pending
- enrollments refresh
- purchase history loads
- subscription plans load
- subscription checkout uses same WebView/status flow

---

## PHASE 6 FINAL CHECKLIST

### SKILL_1 complete when:

- [ ] Payment/enrollment/subscription entities/models compile
- [ ] Repository/use cases cover free enroll, payment, history, subscriptions
- [ ] Remote datasource calls all payment/enrollment endpoints
- [ ] Initiate calls send `Idempotency-Key`
- [ ] Failure mapping is consistent

### SKILL_2 complete when:

- [ ] Course detail CTA reflects access state
- [ ] Free enrollment works
- [ ] Provider selection works
- [ ] Paid checkout opens WebView
- [ ] Payment status polling works
- [ ] Completed payment refreshes access
- [ ] Payment result screen handles completed/pending/failed

### SKILL_3 complete when:

- [ ] Purchase history screen loads payments
- [ ] Enrollment list loads
- [ ] Subscription plans load
- [ ] Subscription checkout starts
- [ ] Active subscription view loads
- [ ] Cancel subscription action works
- [ ] Access state accounts for enrollments and subscriptions
- [ ] `flutter analyze` has 0 errors
- [ ] Android debug APK builds

---

## PHASE 6 → PHASE 7 HANDOFF

**Phase 6 is complete when all 3 SKILL files are done and the checklist above passes.**

To start Phase 7, tell Antigravity:

```text
Phase 6 is complete. Read SKILL.md and all phase7 skill files.
We are starting Phase 7: Progress Dashboard.
Reference: memere_mobile/docs/memere_Design_Specification.md
Build in order as specified in the phase7 skill files.
```

**What Phase 7 will build:**

- student dashboard
- enrolled course progress
- streaks
- completion percentages
- recent activity
- weak areas summary
