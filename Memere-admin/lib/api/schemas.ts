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
  approval_status: z.enum(["pending", "approved", "rejected"]).optional().default("approved"),
  is_active: z.boolean(),
  is_email_verified: z.boolean(),
  last_login_at: z.string().optional().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type User = z.infer<typeof UserSchema>;

export const UserDetailResponseSchema = z.object({
  user: UserSchema,
});
export type UserDetailResponse = z.infer<typeof UserDetailResponseSchema>;

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
  language: z.string().optional().nullable().default("en"),
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

export const AddSectionInputSchema = z.object({
  title: z.string().min(1, "Title is required").max(200),
  description: z.string().max(1000).optional().nullable(),
  order_index: z.number().int().min(0).optional(),
  order: z.number().int().min(0).optional(),
  is_published: z.boolean().optional(),
});
export type AddSectionInput = z.infer<typeof AddSectionInputSchema>;

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
  quiz_id: z.string().optional().nullable(),
  content: z.string().optional().nullable(),
  pdf_url: z.string().optional().nullable(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
}).passthrough();
export type Lesson = z.infer<typeof LessonSchema>;

export const LessonListResponseSchema = z.object({ data: LessonSchema.array() });

export const AddLessonInputSchema = z.object({
  title: z.string().min(1, "Title is required").max(200),
  type: z.enum(["video", "note", "quiz", "mixed"]),
  is_free_preview: z.boolean().optional(),
  duration_seconds: z.number().min(0).optional(),
  is_published: z.boolean().optional(),
  order_index: z.number().int().min(0).optional(),
  quiz_id: z.string().optional().nullable(),
  content: z.string().optional().nullable(),
  pdf_url: z.string().optional().nullable(),
});
export type AddLessonInput = z.infer<typeof AddLessonInputSchema>;

// ---- Quiz ---------------------------------------------------------------------

export const QuizSchema = z.object({
  id: z.string(),
  // Null for a standalone quiz — one authored on its own rather than under a
  // course. Only an admin may create or edit those.
  course_id: z.string().optional().nullable(),
  lesson_id: z.string().optional().nullable(),
  title: z.string(),
  time_limit_seconds: z.number().optional().nullable(),
  pass_percentage: z.number().optional().default(60),
  randomize_questions: z.boolean().optional().default(false),
  max_attempts: z.number().optional().nullable(),
  // Gates standalone quizzes only; a course-attached quiz is gated by access to
  // its course and ignores this flag.
  is_published: z.boolean().optional().default(false),
  question_count: z.number().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});
export type Quiz = z.infer<typeof QuizSchema>;

export const QuizListResponseSchema = z.object({
  data: QuizSchema.array(),
  next_cursor: z.string().optional().nullable(),
});

export const CreateQuizInputSchema = z.object({
  title: z.string().min(1, "Title is required").max(200),
  // Omit or null to create a standalone quiz (admin only).
  course_id: z.string().optional().nullable(),
  lesson_id: z.string().optional().nullable(),
  time_limit_seconds: z.number().min(0).optional().nullable(),
  pass_percentage: z.number().min(0).max(100).optional(),
  randomize_questions: z.boolean().optional(),
  max_attempts: z.number().min(1).optional().nullable(),
  is_published: z.boolean().optional(),
});
export type CreateQuizInput = z.infer<typeof CreateQuizInputSchema>;

export const AnswerSchema = z.object({
  text: z.string().min(1, "Answer text is required"),
  is_correct: z.boolean().optional(),
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
  type: z.enum(["multiple_choice", "true_false"]).optional(),
  points: z.number().min(1).optional(),
  explanation: z.string().optional().nullable(),
  order_index: z.number().int().min(0).optional(),
  subject: z.string().optional().nullable(),
  topic: z.string().optional().nullable(),
  answers: AnswerSchema.array().min(2, "At least 2 answers required"),
});
export type AddQuizQuestionInput = z.infer<typeof AddQuizQuestionInputSchema>;

// ---- Exam ---------------------------------------------------------------------

export const ExamSchema = z.object({
  id: z.string(),
  course_id: z.string().optional().nullable(),
  title: z.string(),
  subject: z.string().optional(),
  grade: z.number().optional(),
  duration_minutes: z.number().optional(),
  total_marks: z.number().optional(),
  pass_marks: z.number().optional(),
  instructions: z.string().optional().nullable(),
  is_published: z.boolean().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});
export type Exam = z.infer<typeof ExamSchema>;

export const ExamListResponseSchema = z.object({
  data: ExamSchema.array(),
  next_cursor: z.string().optional().nullable(),
});

export const CreateExamInputSchema = z.object({
  title: z.string().min(1, "Title is required").max(200),
  // Omit or null to create a standalone exam (admin only).
  course_id: z.string().optional().nullable(),
  subject: z.string().min(1, "Subject is required").max(100),
  grade: z.number().int().min(1).max(13),
  duration_minutes: z.number().min(1).optional(),
  pass_marks: z.number().min(0).optional(),
  instructions: z.string().max(5000).optional().nullable(),
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
  question_id: z.string().optional(),
  marks: z.number().min(1).optional(),
  order_index: z.number().int().min(0).optional(),
  text: z.string().min(1, "Question text is required").optional(),
  type: z.enum(["multiple_choice", "true_false"]).optional(),
  explanation: z.string().optional().nullable(),
  subject: z.string().optional().nullable(),
  topic: z.string().optional().nullable(),
  answers: AnswerSchema.array().min(2, "At least 2 answers required").optional(),
});
export type ExamQuestionInput = z.infer<typeof ExamQuestionInputSchema>;

// ---- Assessment authoring (per-question editing) -------------------------------

// The authoring schemas below carry is_correct, which the student-facing API
// never returns. They are only ever parsed from the teacher/admin-gated
// authoring routes (GET /quizzes/:id/questions, GET /exams/:id/questions).

export const AuthoredAnswerSchema = z.object({
  id: z.string().optional(),
  text: z.string(),
  is_correct: z.boolean().optional().default(false),
  order_index: z.number().optional().default(0),
});
export type AuthoredAnswer = z.infer<typeof AuthoredAnswerSchema>;

export const AuthoredQuestionSchema = z.object({
  id: z.string(),
  text: z.string(),
  type: z.string(),
  points: z.number().optional().default(1),
  marks: z.number().optional().default(1),
  explanation: z.string().optional().nullable(),
  order_index: z.number().optional().default(0),
  subject: z.string().optional().nullable(),
  topic: z.string().optional().nullable(),
  answers: AuthoredAnswerSchema.array().default([]),
});
export type AuthoredQuestion = z.infer<typeof AuthoredQuestionSchema>;

export const AuthoredQuizSchema = z.object({
  quiz: QuizSchema,
  questions: AuthoredQuestionSchema.array().default([]),
});
export type AuthoredQuiz = z.infer<typeof AuthoredQuizSchema>;

export const AuthoredExamSchema = z.object({
  exam: ExamSchema,
  questions: AuthoredQuestionSchema.array().default([]),
});
export type AuthoredExam = z.infer<typeof AuthoredExamSchema>;

// An answer option as submitted by the editor. At least one must be correct,
// which the editor enforces before sending and the backend re-checks.
export const EditAnswerInputSchema = z.object({
  text: z.string().min(1, "Answer text is required"),
  is_correct: z.boolean(),
  order_index: z.number().int().min(0).optional(),
});

export const UpdateQuestionInputSchema = z.object({
  text: z.string().min(1, "Question text is required").optional(),
  type: z.enum(["multiple_choice", "true_false", "short_answer"]).optional(),
  points: z.number().int().min(1).optional(),
  explanation: z.string().optional().nullable(),
  order_index: z.number().int().min(0).optional(),
  subject: z.string().optional().nullable(),
  topic: z.string().optional().nullable(),
  answers: EditAnswerInputSchema.array().min(2, "At least 2 answers required").optional(),
});
export type UpdateQuestionInput = z.infer<typeof UpdateQuestionInputSchema>;

export const UpdateExamQuestionInputSchema = UpdateQuestionInputSchema.omit({
  points: true,
}).extend({
  marks: z.number().int().min(1).optional(),
});
export type UpdateExamQuestionInput = z.infer<typeof UpdateExamQuestionInputSchema>;

export const ReorderQuestionsInputSchema = z.object({
  question_ids: z.string().array().min(1),
});
export type ReorderQuestionsInput = z.infer<typeof ReorderQuestionsInputSchema>;

// ---- Assessment import --------------------------------------------------------

// Mirrors the backend wire format in dto/assessment_import.go. The browser-side
// parser normalizes uploads (JSON or HTML) into this shape, the preview renders
// it, and only after the operator confirms is it POSTed.

export const IMPORT_FORMAT_VERSION = 1;

export const ImportAnswerSchema = z.object({
  text: z.string().min(1, "Answer text is required"),
  is_correct: z.boolean(),
});

export const ImportQuestionSchema = z
  .object({
    text: z.string().min(1, "Question text is required"),
    type: z.enum(["multiple_choice", "true_false", "short_answer"]).optional(),
    points: z.number().int().min(1).optional(),
    marks: z.number().int().min(1).optional(),
    explanation: z.string().optional().nullable(),
    subject: z.string().optional().nullable(),
    topic: z.string().optional().nullable(),
    answers: ImportAnswerSchema.array().optional(),
    options: z.string().array().optional(),
    correct: z.unknown().optional(),
  })
  .refine((q) => (q.answers?.length ?? 0) > 0 || (q.options?.length ?? 0) > 0, {
    message: 'Provide either "answers" or "options" + "correct"',
  })
  .refine((q) => !((q.answers?.length ?? 0) > 0 && (q.options?.length ?? 0) > 0), {
    message: 'Use either "answers" or "options" + "correct", not both',
  })
  .refine((q) => (q.options?.length ?? 0) === 0 || q.correct !== undefined, {
    message: '"correct" is required when using "options"',
  });
export type ImportQuestion = z.infer<typeof ImportQuestionSchema>;

export const ImportQuizHeaderSchema = z.object({
  title: z.string().min(1, "Title is required").max(200),
  course_id: z.string().optional().nullable(),
  lesson_id: z.string().optional().nullable(),
  time_limit_seconds: z.number().int().min(0).optional().nullable(),
  pass_percentage: z.number().min(0).max(100).optional(),
  randomize_questions: z.boolean().optional(),
  max_attempts: z.number().int().min(1).optional().nullable(),
});

export const ImportExamHeaderSchema = z.object({
  title: z.string().min(1, "Title is required").max(200),
  course_id: z.string().optional().nullable(),
  subject: z.string().min(1, "Subject is required").max(100),
  grade: z.number().int().min(1).max(13),
  duration_minutes: z.number().int().min(1),
  pass_marks: z.number().int().min(0).optional(),
  instructions: z.string().max(5000).optional().nullable(),
});

export const ImportDocumentSchema = z
  .object({
    version: z.literal(IMPORT_FORMAT_VERSION),
    kind: z.enum(["quiz", "exam"]),
    quiz: ImportQuizHeaderSchema.optional(),
    exam: ImportExamHeaderSchema.optional(),
    questions: ImportQuestionSchema.array().min(1, "At least one question is required"),
  })
  .refine((d) => d.kind !== "quiz" || d.quiz !== undefined, {
    message: 'A "quiz" block is required when kind is "quiz"',
    path: ["quiz"],
  })
  .refine((d) => d.kind !== "exam" || d.exam !== undefined, {
    message: 'An "exam" block is required when kind is "exam"',
    path: ["exam"],
  });
export type ImportDocument = z.infer<typeof ImportDocumentSchema>;

export const ImportResultSchema = z.object({
  kind: z.string(),
  id: z.string(),
  title: z.string(),
  question_count: z.number(),
  total_marks: z.number().optional().nullable(),
});
export type ImportResult = z.infer<typeof ImportResultSchema>;

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
  short_description: z.string().max(300).optional().nullable(),
  subject: z.string().min(1, "Subject is required"),
  grade: z.number().int().min(1, "Grade is required").max(13),
  level: z.enum(["beginner", "intermediate", "advanced"]),
  price: z.number().min(0).optional(),
  currency: z.string().optional(),
  is_free: z.boolean().optional(),
  language: z.string().optional(),
  thumbnail_url: z.string().optional().nullable(),
  metadata: z.record(z.string(), z.unknown()).optional().nullable(),
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

// ---- Course Access & Requests ------------------------------------------------

export const UserCourseAccessItemSchema = z.object({
  course_id: z.string(),
  title: z.string().optional(),
  course_title: z.string().optional(),
  subject: z.string().optional(),
  grade: z.number().optional(),
  has_access: z.boolean(),
  access_source: z.string().optional().nullable(),
  is_free: z.boolean(),
  price: z.union([z.string(), z.number()]),
  currency: z.string(),
  request_status: z.string().optional().nullable(),
});
export type UserCourseAccessItem = z.infer<typeof UserCourseAccessItemSchema>;

export const UserCourseAccessListResponseSchema = z.object({
  student_id: z.string().optional(),
  courses: z.array(UserCourseAccessItemSchema),
});
export type UserCourseAccessListResponse = z.infer<typeof UserCourseAccessListResponseSchema>;

export const CourseAccessRequestSchema = z.object({
  id: z.string(),
  student_id: z.string(),
  student_name: z.string().optional(),
  student_email: z.string().optional(),
  course_id: z.string(),
  course_title: z.string().optional(),
  status: z.enum(["pending", "approved", "rejected"]),
  reviewed_by: z.string().optional().nullable(),
  reviewed_at: z.string().optional().nullable(),
  rejection_reason: z.string().optional().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});
export type CourseAccessRequest = z.infer<typeof CourseAccessRequestSchema>;

export const CourseAccessRequestListResponseSchema = z.object({
  requests: z.array(CourseAccessRequestSchema),
  next: z.string().optional().nullable(),
  total: z.number().optional(),
});
export type CourseAccessRequestListResponse = z.infer<typeof CourseAccessRequestListResponseSchema>;

// ---- Announcements -----------------------------------------------------------

export const AnnouncementSchema = z.object({
  id: z.string(),
  title: z.string().nullish().transform((v) => v || "Announcement"),
  body: z.string().nullish().transform((v) => v || ""),
  recipient_count: z.number().nullish().transform((v) => v ?? 0),
  created_at: z.string().nullish().transform((v) => v || new Date().toISOString()),
});
export type Announcement = z.infer<typeof AnnouncementSchema>;

export const AnnouncementListResponseSchema = z.object({
  announcements: z.array(AnnouncementSchema).default([]),
});
export type AnnouncementListResponse = z.infer<typeof AnnouncementListResponseSchema>;


