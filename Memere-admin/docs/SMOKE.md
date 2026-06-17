# Phase 1 Smoke Test

## Prerequisites

- Backend running at `API_BASE_URL` with seeded admin user
- `NEXT_PUBLIC_API_BASE_URL` (if used) or `API_BASE_URL` set in `.env.local`

## Steps

```bash
pnpm build && pnpm start
```

1. **Login** — open `http://localhost:3000/login`, enter admin credentials → redirected to `/`
2. **Dashboard KPIs** — total students, teachers, gross revenue, MRR, completed payments, and refunded amount all render with real values from the backend
3. **Date range** — click "7d", "90d", "This year" → KPIs and revenue chart update; no NaN or crash on empty ranges
4. **Engagement** — quiz pass rate, exam pass rate, course completion tiles render with progress bars
5. **Revenue chart** — bar chart shows gross revenue per provider; empty state shown when no data in range
6. **Sidebar nav** — click each link (Users, Courses, Payments, Revenue, Announcements) → placeholder pages load with active highlight
7. **Sidebar collapse** — click toggle on desktop → icon-only mode with tooltips
8. **Mobile** — resize to ≤768px → sidebar hidden; hamburger in header opens drawer
9. **Theme toggle** — header sun/moon button switches light/dark; persists on reload
10. **User menu** — avatar in header shows name/email/role; click **Logout** → `POST /api/auth/logout` → redirected to `/login`
11. **Auth guard** — log out, then navigate directly to `/` → redirected to `/login`

---

# Phase 2 Smoke Test

## Users

12. **Users list** — navigate to `/users` → table loads with name, email, role badge, status badge, joined date
13. **Role filter** — select "Teacher" → table refetches showing only teachers; select "All roles" → resets
14. **Pagination** — click Next → second page loads; click Prev → returns to first page; change limit to 50
15. **Search** — type a name/email → filters current page; note in UI confirms search is page-scoped
16. **User detail** — click a row → `/users/[id]` loads; breadcrumb shows user's name; full profile renders
17. **Suspend** — click "Suspend user" → dialog opens; "Confirm" disabled until reason typed; enter reason → confirm → toast success, status badge flips to Suspended
18. **Reactivate** — click "Reactivate user" → confirm dialog → toast success, status flips to Active
19. **Change role** — click "Change role" → select new role → save → toast success, role badge updates
20. **Backend error** — force a 403/error → toast shows backend's mapped message; dialog stays open
21. **Not found** — navigate to `/users/nonexistent-id` → friendly not-found page with back link

## Courses

22. **Courses list** — navigate to `/courses` → table loads with title, subject, grade, price (Free vs amount), status badge, date
23. **Pagination** — Next/Prev and limit selector work correctly
24. **Course detail** — click a row → `/courses/[id]` loads; breadcrumb shows course title; metadata + sections list renders
25. **Unpublish** — click "Unpublish" (only on published courses) → dialog opens; confirm disabled until reason typed → enter reason → confirm → toast, status flips to Unpublished; button disappears
26. **Not found** — navigate to `/courses/nonexistent-id` → friendly not-found page with back link

---

# Phase 3 Smoke Test

## Payments

27. **Payments list** — navigate to `/payments` → table loads with short IDs (monospace), amounts, status badges, provider, date
28. **Status filter** — select "Completed" → only completed payments shown; select "All statuses" → resets; pagination works
29. **Copy ID** — click a payment ID chip → full ID copied to clipboard (no row navigation)
30. **Student link** — click student ID in row → navigates to `/users/[id]` without triggering row click
31. **Payment detail** — click a row → `/payments/[id]` loads; breadcrumb shows short payment ID; full record renders
32. **Refund** — open a completed payment → "Refund" button visible; dialog shows formatted amount + "cannot be undone"; confirm → toast success, status flips to Refunded; button disappears
33. **Refund on non-completed** — open a pending/failed/refunded payment → no Refund button shown
34. **Reconcile** — click "Reconcile pending" on payments list → toast shows real `{reconciled}` count
35. **Backend error** — force an error on refund/reconcile → toast shows backend's mapped message; no false success
36. **Not found** — navigate to `/payments/nonexistent-id` → friendly not-found page

## Revenue

37. **Revenue page** — navigate to `/revenue` → 4 KPI tiles (Gross, Refunded, MRR, Completed Payments) load from backend
38. **Date range** — switch between 7d / 30d / 90d / This year → KPIs + charts + breakdown table update; all numbers match backend
39. **Zero range** — select a future range with no data → charts show empty state (no NaN, no broken layout)
40. **Revenue bar chart** — shows gross per provider; tooltip displays formatted money string from backend
41. **Units pie chart** — shows payment count per provider; tooltip shows count
42. **Provider breakdown table** — provider, gross (formatted), payments count; currency always reads from `PLATFORM_CURRENCY` (analytics endpoints have no currency field — documented gap)
43. **Financial overview + trend** — section renders with 4-bucket trend; gap note clearly states "backend has no subscription-list endpoint in v1"
44. **Money formatting** — all amounts display via `formatMoney`; no client-side arithmetic on displayed figures; `parseFloat` used only for chart axis values

---

# Phase 4 Smoke Test

## Announcements

45. **Composer** — navigate to `/announcements` → form renders with title, message, audience select, optional data fields
46. **Validation** — submit with empty title/body → inline errors; form does not submit
47. **Send to students** — fill in title + body, select "Students" → click "Preview & send" → confirm dialog shows segment; click "Send announcement" → toast success; entry appears in session history below
48. **Send to all — double confirm** — select "All users" → confirm dialog shows red warning + checkbox; "Send announcement" disabled until checkbox checked → check it → send → toast success
49. **Session history** — history list shows title, segment badge, time; clearly labelled "Session-only"
50. **Backend error** — simulate error → toast shows backend's mapped message; dialog stays open

## Security & hardening

51. **Security headers** — `curl -sI http://localhost:3000/login | grep -iE 'content-security|x-frame|x-content|referrer'` → all 5 headers present
52. **No token in client** — DevTools → Application → Local Storage / Session Storage → empty; Cookies shows `mm_access` and `mm_refresh` with HttpOnly flag
53. **Token refresh** — expire the access token (shorten backend TTL or wait); perform an action → action succeeds transparently via silent refresh
54. **Expired session** — revoke the refresh token → next action redirects to `/login` with "session expired" toast
55. **CSRF** — send a cross-origin POST to `/api/admin/users/:id/suspend` with a different `Origin` header → 403 response
56. **Auth guard** — log out, then navigate directly to `/users`, `/courses`, `/payments`, `/revenue`, `/announcements` → all redirect to `/login`

## Final end-to-end checklist (against deployed environment)

57. Login as admin → dashboard KPIs and charts load with live data
58. Users: list → role filter → paginate → detail → suspend (reason) → reactivate → change role
59. Courses: list → paginate → detail (sections visible) → unpublish (reason)
60. Payments: list → status filter → detail → refund (strong confirm) → reconcile
61. Revenue: KPI tiles + bar + pie + trend for 30d range; switch to 7d → all update; zero range → graceful empty states
62. Announcements: send to teachers → double-confirm send to all
63. Logout → redirect to `/login`; `/` while logged out → redirect to `/login`
