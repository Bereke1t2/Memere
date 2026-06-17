import { apiFetch } from "./server";
import {
  AuthResponseSchema,
  UserSchema,
  OverviewSchema,
  EngagementStatsSchema,
  RevenueBreakdownResponseSchema,
  PaginatedUsersSchema,
  PaginatedCoursesSchema,
  PaginatedPaymentsSchema,
  AdminPaymentSchema,
  ReconcileResponseSchema,
  CourseSchema,
  SectionSchema,
  type AuthResponse,
  type User,
  type Overview,
  type EngagementStats,
  type RevenueBreakdownItem,
  type PaginatedUsers,
  type PaginatedCourses,
  type PaginatedPayments,
  type AdminPayment,
  type ReconcileResponse,
  type Course,
  type Section,
} from "./schemas";

// ---- Auth ----------------------------------------------------------------------

export async function login(
  email: string,
  password: string
): Promise<AuthResponse> {
  const data = await apiFetch("/auth/login", {
    method: "POST",
    body: { email, password },
    schema: AuthResponseSchema,
  });
  return data!;
}

export async function me(): Promise<User> {
  const data = await apiFetch("/auth/me", { schema: UserSchema });
  return data!;
}

export async function logout(refreshToken: string): Promise<void> {
  await apiFetch("/auth/logout", {
    method: "POST",
    body: { refresh_token: refreshToken },
  });
}

export async function refresh(refreshToken: string): Promise<AuthResponse> {
  const data = await apiFetch("/auth/refresh", {
    method: "POST",
    body: { refresh_token: refreshToken },
    schema: AuthResponseSchema,
  });
  return data!;
}

// ---- Analytics (Phase 1) -------------------------------------------------------

export async function getOverview(from: string, to: string): Promise<Overview> {
  const data = await apiFetch(
    `/admin/analytics/overview?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`,
    { schema: OverviewSchema }
  );
  return data!;
}

export async function getEngagement(): Promise<EngagementStats> {
  const data = await apiFetch("/admin/analytics/engagement", {
    schema: EngagementStatsSchema,
  });
  return data!;
}

export async function getRevenueBreakdown(
  from: string,
  to: string
): Promise<RevenueBreakdownItem[]> {
  const data = await apiFetch(
    `/admin/analytics/revenue?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`,
    { schema: RevenueBreakdownResponseSchema }
  );
  return data!.breakdown;
}

// ---- Users (Phase 2) -----------------------------------------------------------

export async function listUsers(params: {
  limit?: number;
  after?: string;
  role?: string;
}): Promise<PaginatedUsers> {
  const qs = new URLSearchParams();
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.after) qs.set("after", params.after);
  if (params.role) qs.set("role", params.role);
  const data = await apiFetch(`/admin/users?${qs}`, {
    schema: PaginatedUsersSchema,
  });
  return data!;
}

export async function getUser(id: string): Promise<User> {
  const data = await apiFetch(`/admin/users/${id}`, { schema: UserSchema });
  return data!;
}

export async function suspendUser(id: string, reason: string): Promise<void> {
  await apiFetch(`/admin/users/${id}/suspend`, {
    method: "POST",
    body: { reason },
  });
}

export async function reactivateUser(id: string): Promise<void> {
  await apiFetch(`/admin/users/${id}/reactivate`, { method: "POST" });
}

export async function changeRole(id: string, role: string): Promise<void> {
  await apiFetch(`/admin/users/${id}/role`, {
    method: "POST",
    body: { role },
  });
}

// ---- Courses (Phase 2) ---------------------------------------------------------

export async function listCourses(params: {
  limit?: number;
  after?: string;
}): Promise<PaginatedCourses> {
  const qs = new URLSearchParams();
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.after) qs.set("after", params.after);
  const data = await apiFetch(`/admin/courses?${qs}`, {
    schema: PaginatedCoursesSchema,
  });
  return data!;
}

export async function getCourse(id: string): Promise<Course> {
  const data = await apiFetch(`/courses/${id}`, { schema: CourseSchema });
  return data!;
}

export async function getCourseSections(id: string): Promise<Section[]> {
  const data = await apiFetch(`/courses/${id}/sections`, {
    schema: SectionSchema.array(),
  });
  return data ?? [];
}

export async function unpublishCourse(id: string, reason: string): Promise<void> {
  await apiFetch(`/admin/courses/${id}/unpublish`, {
    method: "POST",
    body: { reason },
  });
}

// ---- Payments (Phase 3) --------------------------------------------------------

export async function listPayments(params: {
  limit?: number;
  after?: string;
  status?: string;
}): Promise<PaginatedPayments> {
  const qs = new URLSearchParams();
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.after) qs.set("after", params.after);
  if (params.status) qs.set("status", params.status);
  const data = await apiFetch(`/admin/payments?${qs}`, {
    schema: PaginatedPaymentsSchema,
  });
  return data!;
}

export async function getPayment(id: string): Promise<AdminPayment> {
  const data = await apiFetch(`/admin/payments/${id}`, {
    schema: AdminPaymentSchema,
  });
  return data!;
}

export async function refundPayment(id: string): Promise<void> {
  await apiFetch(`/payments/${id}/refund`, { method: "POST" });
}

export async function reconcilePayments(): Promise<ReconcileResponse> {
  const data = await apiFetch("/admin/payments/reconcile", {
    method: "POST",
    schema: ReconcileResponseSchema,
  });
  return data!;
}

// ---- Announcements (Phase 4) ---------------------------------------------------

export async function broadcast(input: {
  title: string;
  body: string;
  segment: "all" | "students" | "teachers" | "subscribers";
  data?: Record<string, string>;
}): Promise<void> {
  await apiFetch("/admin/announcements", {
    method: "POST",
    body: input,
  });
}
