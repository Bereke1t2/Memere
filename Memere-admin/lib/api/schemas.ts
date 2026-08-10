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
  slug: z.string().optional().nullable(),
  description: z.string().optional().nullable().default(""),
  short_description: z.string().optional().nullable(),
  subject: z.string().optional().default("General"),
  grade: z.number().optional().default(1),
  thumbnail_url: z.string().optional().nullable(),
  price: z.number().optional().default(0),
  currency: z.string().optional().default("ETB"),
  is_free: z.boolean().optional().default(false),
  is_published: z.boolean().optional().default(false),
  language: z.string().optional().nullable(),
  level: z.string().optional().nullable(),
  total_duration_seconds: z.number().optional(),
  total_lessons: z.number().optional(),
  rating_avg: z.number().optional(),
  enrollment_count: z.number().optional(),
  metadata: z.record(z.string(), z.unknown()).optional().nullable(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
}).passthrough();

export type Course = z.infer<typeof CourseSchema>;

// ---- Section -------------------------------------------------------------------

export const SectionSchema = z.object({
  id: z.string(),
  course_id: z.string().optional(),
  title: z.string(),
  description: z.string().optional().nullable(),
  order: z.number().optional(),
  order_index: z.number().optional(),
  is_published: z.boolean().optional(),
  total_lessons: z.number().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
}).passthrough();

export const SectionListResponseSchema = z.object({
  data: SectionSchema.array(),
});

// ---- Lesson -------------------------------------------------------------------

export const LessonSchema = z.object({
  id: z.string(),
  section_id: z.string().optional(),
  course_id: z.string().optional(),
  title: z.string(),
  type: z.string(),
  order_index: z.number().optional(),
  is_free_preview: z.boolean().optional(),
  duration_seconds: z.number().optional(),
  is_published: z.boolean().optional(),
  video_id: z.string().optional().nullable(),
  content: z.string().optional().nullable(),
  pdf_url: z.string().optional().nullable(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
}).passthrough();
export type Lesson = z.infer<typeof LessonSchema>;

export const LessonListResponseSchema = z.object({ data: LessonSchema.array() });

export const AddLessonInputSchema = z.object({
  title: z.string().min(1, "Title is required"),
  type: z.enum(["video", "note", "quiz", "mixed"]),
  is_free_preview: z.boolean(),
  duration_seconds: z.number().min(0).optional(),
  is_published: z.boolean().optional(),
  content: z.string().optional(),
  pdf_url: z.string().optional(),
});
export type AddLessonInput = z.infer<typeof AddLessonInputSchema>;

// ---- Quiz ---------------------------------------------------------------------

export const QuizSchema = z.object({
  id: z.string(),
  course_id: z.string().optional(),
  title: z.string(),
  time_limit_seconds: z.number().optional(),
  pass_percentage: z.number().optional(),
  randomize_questions: z.boolean().optional(),
  max_attempts: z.number().optional(),
  question_count: z.number().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});
export type Quiz = z.infer<typeof QuizSchema>;

export const QuizListResponseSchema = z.object({ data: QuizSchema.array() });

export const CreateQuizInputSchema = z.object({
  title: z.string().min(1, "Title is required"),
  time_limit_seconds: z.number().min(0),
  pass_percentage: z.number().min(0).max(100),
  randomize_questions: z.boolean(),
  max_attempts: z.number().min(1),
});
export type CreateQuizInput = z.infer<typeof CreateQuizInputSchema>;

export const AnswerSchema = z.object({
  text: z.string().min(1),
  is_correct: z.boolean(),
});
export type Answer = z.infer<typeof AnswerSchema>;

export const QuizQuestionSchema = z.object({
  id: z.string(),
  quiz_id: z.string().optional(),
  text: z.string(),
  type: z.string(),
  points: z.number().optional(),
  explanation: z.string().optional().nullable(),
  order_index: z.number().optional(),
  subject: z.string().optional().nullable(),
  topic: z.string().optional().nullable(),
  answers: AnswerSchema.array().optional(),
});
export type QuizQuestion = z.infer<typeof QuizQuestionSchema>;

export const AddQuizQuestionInputSchema = z.object({
  text: z.string().min(1, "Question text is required"),
  type: z.enum(["multiple_choice", "true_false"]),
  points: z.number().min(1),
  explanation: z.string().optional(),
  order_index: z.number().min(0),
  subject: z.string().optional(),
  topic: z.string().optional(),
  answers: AnswerSchema.array().min(2, "At least 2 answers required"),
});
export type AddQuizQuestionInput = z.infer<typeof AddQuizQuestionInputSchema>;

// ---- Exam ---------------------------------------------------------------------

export const ExamSchema = z.object({
  id: z.string(),
  course_id: z.string().optional(),
  title: z.string(),
  subject: z.string().optional(),
  grade: z.number().optional(),
  duration_minutes: z.number().optional(),
  pass_marks: z.number().optional(),
  instructions: z.string().optional().nullable(),
  is_published: z.boolean().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});
export type Exam = z.infer<typeof ExamSchema>;

export const ExamListResponseSchema = z.object({ data: ExamSchema.array() });

export const CreateExamInputSchema = z.object({
  title: z.string().min(1, "Title is required"),
  subject: z.string().min(1, "Subject is required"),
  grade: z.number().int().min(1).max(12),
  duration_minutes: z.number().min(1),
  pass_marks: z.number().min(0),
  instructions: z.string().optional(),
});
export type CreateExamInput = z.infer<typeof CreateExamInputSchema>;

export const ExamStatsSchema = z.object({
  exam_id: z.string(),
  total_attempts: z.number(),
  avg_score: z.number().optional(),
  pass_rate: z.number().optional(),
  highest_score: z.number().optional(),
  lowest_score: z.number().optional(),
}).passthrough();
export type ExamStats = z.infer<typeof ExamStatsSchema>;

export const ExamQuestionInputSchema = z.object({
  text: z.string().min(1, "Question text is required"),
  type: z.enum(["multiple_choice", "true_false"]),
  marks: z.number().min(1),
  order_index: z.number().min(0),
  answers: AnswerSchema.array().min(2, "At least 2 answers required"),
  explanation: z.string().optional(),
  question_id: z.string().optional(),
});
export type ExamQuestionInput = z.infer<typeof ExamQuestionInputSchema>;

// ---- Video pipeline -----------------------------------------------------------

export const VideoUploadResponseSchema = z.object({
  video_id: z.string(),
  upload_url: z.string(),
  expires_at: z.string().optional(),
});
export type VideoUploadResponse = z.infer<typeof VideoUploadResponseSchema>;

export const VideoStatusSchema = z.object({
  id: z.string(),
  lesson_id: z.string().optional().nullable(),
  status: z.string(),
  duration_seconds: z.number().optional().nullable(),
  thumbnail_url: z.string().optional().nullable(),
  hls_url: z.string().optional().nullable(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
}).passthrough();
export type VideoStatus = z.infer<typeof VideoStatusSchema>;

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
  next: z.string().nullable(),
});

export const PaginatedCoursesSchema = z.object({
  courses: CourseSchema.array(),
  next: z.string().nullable(),
});

export const PaginatedPaymentsSchema = z.object({
  payments: AdminPaymentSchema.array(),
  next: z.string().nullable(),
});

export const RevenueBreakdownResponseSchema = z.object({
  breakdown: RevenueBreakdownItemSchema.array(),
});

export type PaginatedUsers = z.infer<typeof PaginatedUsersSchema>;
export type PaginatedCourses = z.infer<typeof PaginatedCoursesSchema>;
export type PaginatedPayments = z.infer<typeof PaginatedPaymentsSchema>;

// ---- Teacher: course list (public /courses endpoint) ---------------------------

export const TeacherCourseListSchema = z.object({
  data: CourseSchema.array(),
  next_cursor: z.string().nullable(),
  limit: z.number(),
});
export type TeacherCourseList = z.infer<typeof TeacherCourseListSchema>;

// ---- Teacher: create / update course ------------------------------------------

export const CreateCourseInputSchema = z.object({
  title: z.string().min(3, "Title must be at least 3 characters").max(200),
  description: z.string().min(10, "Description must be at least 10 characters"),
  subject: z.string().min(1, "Subject is required"),
  grade: z.number().int().min(1, "Grade is required").max(12),
  level: z.enum(["beginner", "intermediate", "advanced"]),
  price: z.number().min(0),
  language: z.string(),
});
export type CreateCourseInput = z.infer<typeof CreateCourseInputSchema>;

// ---- Teacher: earnings --------------------------------------------------------

export const EarningsSchema = z.object({
  teacher_id: z.string(),
  from: z.string(),
  to: z.string(),
  gross: z.string(),
  teacher_share: z.string(),
  earnings: z.string(),
  platform_fee: z.string(),
  units: z.number(),
});
export type Earnings = z.infer<typeof EarningsSchema>;

export const CourseSalesSchema = z.object({
  course_id: z.string(),
  gross: z.string(),
  units: z.number(),
});
export type CourseSales = z.infer<typeof CourseSalesSchema>;
