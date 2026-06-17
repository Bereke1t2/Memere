import { z } from "zod";

// ---- Primitives ----------------------------------------------------------------

export const ErrorEnvelopeSchema = z.object({
  code: z.string(),
  message: z.string(),
  details: z.unknown().optional(),
});

// ---- User ----------------------------------------------------------------------

export const UserSchema = z.object({
  id: z.string(),
  email: z.string(),
  phone: z.string().optional().nullable(),
  role: z.string(),
  first_name: z.string(),
  last_name: z.string(),
  avatar_url: z.string().optional().nullable(),
  is_active: z.boolean(),
  is_email_verified: z.boolean(),
  last_login_at: z.string().optional().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type User = z.infer<typeof UserSchema>;

// ---- Auth ----------------------------------------------------------------------

export const AuthResponseSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string(),
  expires_in: z.number(),
  user: UserSchema.optional().nullable(),
});

export type AuthResponse = z.infer<typeof AuthResponseSchema>;

// ---- Course --------------------------------------------------------------------

export const CourseSchema = z.object({
  id: z.string(),
  teacher_id: z.string(),
  title: z.string(),
  slug: z.string().optional(),
  description: z.string(),
  short_description: z.string().optional().nullable(),
  subject: z.string(),
  grade: z.number(),
  thumbnail_url: z.string().optional().nullable(),
  price: z.number(),
  currency: z.string(),
  is_free: z.boolean(),
  is_published: z.boolean(),
  language: z.string().optional(),
  level: z.string().optional(),
  total_duration_seconds: z.number().optional(),
  total_lessons: z.number().optional(),
  rating_avg: z.number().optional(),
  enrollment_count: z.number().optional(),
  metadata: z.record(z.string(), z.unknown()).optional().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Course = z.infer<typeof CourseSchema>;

// ---- Section -------------------------------------------------------------------

export const SectionSchema = z.object({
  id: z.string(),
  course_id: z.string().optional(),
  title: z.string(),
  order: z.number().optional(),
  is_published: z.boolean().optional(),
  total_lessons: z.number().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export type Section = z.infer<typeof SectionSchema>;

// ---- Payment -------------------------------------------------------------------

export const AdminPaymentSchema = z.object({
  id: z.string(),
  student_id: z.string(),
  amount: z.string(),
  currency: z.string(),
  status: z.string(),
  provider: z.string(),
  provider_txn_id: z.string().optional().nullable(),
  created_at: z.string(),
});

export type AdminPayment = z.infer<typeof AdminPaymentSchema>;

// ---- Analytics -----------------------------------------------------------------

// Money fields (gross_revenue, refunded_amount, mrr) are shopspring decimal →
// serialise as JSON strings.
export const OverviewSchema = z.object({
  total_students: z.number(),
  total_teachers: z.number(),
  total_admins: z.number(),
  gross_revenue: z.string(),
  refunded_amount: z.string(),
  completed_payments: z.number(),
  mrr: z.string(),
  from: z.string(),
  to: z.string(),
});

export type Overview = z.infer<typeof OverviewSchema>;

export const RevenueBreakdownItemSchema = z.object({
  provider: z.string(),
  gross: z.string(),
  units: z.number(),
});

export type RevenueBreakdownItem = z.infer<typeof RevenueBreakdownItemSchema>;

export const EngagementStatsSchema = z.object({
  avg_quiz_pass_rate: z.number(),
  avg_exam_pass_rate: z.number(),
  avg_completion_pct: z.number(),
});

export type EngagementStats = z.infer<typeof EngagementStatsSchema>;

// ---- Reconcile -----------------------------------------------------------------

export const ReconcileResponseSchema = z.object({
  reconciled: z.number(),
});

export type ReconcileResponse = z.infer<typeof ReconcileResponseSchema>;

// ---- Paginated list helpers ----------------------------------------------------

export function PaginatedSchema<T extends z.ZodTypeAny>(itemSchema: T) {
  return z.object({
    items: itemSchema.array(),
    next: z.string(),
  });
}

export const PaginatedUsersSchema = z.object({
  users: UserSchema.array(),
  next: z.string(),
});

export const PaginatedCoursesSchema = z.object({
  courses: CourseSchema.array(),
  next: z.string(),
});

export const PaginatedPaymentsSchema = z.object({
  payments: AdminPaymentSchema.array(),
  next: z.string(),
});

export const RevenueBreakdownResponseSchema = z.object({
  breakdown: RevenueBreakdownItemSchema.array(),
});

export type PaginatedUsers = z.infer<typeof PaginatedUsersSchema>;
export type PaginatedCourses = z.infer<typeof PaginatedCoursesSchema>;
export type PaginatedPayments = z.infer<typeof PaginatedPaymentsSchema>;
