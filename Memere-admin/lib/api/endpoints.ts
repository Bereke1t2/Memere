import { apiFetch } from "./server";
import {
  AuthResponseSchema,
  UserSchema,
  UserDetailResponseSchema,
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
  AddSectionInputSchema,
  TeacherCourseListSchema,
  SectionListResponseSchema,
  LessonSchema,
  LessonListResponseSchema,
  QuizSchema,
  QuizListResponseSchema,
  QuizQuestionSchema,
  AuthoredQuizSchema,
  AuthoredExamSchema,
  AuthoredQuestionSchema,
  ImportResultSchema,
  ExamSchema,
  ExamListResponseSchema,
  ExamStatsSchema,
  ExamQuestionInputSchema,
  VideoUploadResponseSchema,
  VideoStatusSchema,
  EarningsSchema,
  CourseSalesSchema,
  UserCourseAccessListResponseSchema,
  CourseAccessRequestSchema,
  CourseAccessRequestListResponseSchema,
  AnnouncementListResponseSchema,
  AnnouncementSchema,
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
  type AddSectionInput,
  type TeacherCourseList,
  type CreateCourseInput,
  type Lesson,
  type AddLessonInput,
  type Quiz,
  type CreateQuizInput,
  type QuizQuestion,
  type AddQuizQuestionInput,
  type AuthoredQuiz,
  type AuthoredExam,
  type AuthoredQuestion,
  type UpdateQuestionInput,
  type UpdateExamQuestionInput,
  type ImportDocument,
  type ImportResult,
  type Exam,
  type CreateExamInput,
  type ExamStats,
  type ExamQuestionInput,
  type VideoUploadResponse,
  type VideoStatus,
  type Earnings,
  type CourseSales,
  type UserCourseAccessItem,
  type UserCourseAccessListResponse,
  type CourseAccessRequest,
  type CourseAccessRequestListResponse,
  type Announcement,
  type AnnouncementListResponse,
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
  approval_status?: string;
}): Promise<PaginatedUsers> {
  const qs = new URLSearchParams();
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.after) qs.set("after", params.after);
  if (params.role) qs.set("role", params.role);
  if (params.approval_status) qs.set("approval_status", params.approval_status);
  const data = await apiFetch(`/admin/users?${qs}`, {
    schema: PaginatedUsersSchema,
  });
  return data!;
}

export async function getUser(id: string): Promise<User> {
  const data = await apiFetch(`/admin/users/${id}`, { schema: UserSchema });
  return data!;
}

export async function approveUser(id: string): Promise<void> {
  await apiFetch(`/admin/users/${id}/approve`, {
    method: "POST",
  });
}

export async function rejectUser(id: string, reason?: string): Promise<void> {
  await apiFetch(`/admin/users/${id}/reject`, {
    method: "POST",
    body: { reason },
  });
}

export async function getUserCourseAccess(id: string): Promise<UserCourseAccessListResponse> {
  const data = await apiFetch(`/admin/users/${id}/courses`, {
    schema: UserCourseAccessListResponseSchema,
  });
  return data!;
}

export async function grantCourseAccess(id: string, courseIds: string[]): Promise<void> {
  await apiFetch(`/admin/users/${id}/grant-access`, {
    method: "POST",
    body: { course_ids: courseIds },
  });
}

export async function grantAllCoursesAccess(id: string): Promise<void> {
  await apiFetch(`/admin/users/${id}/grant-access`, {
    method: "POST",
    body: { grant_all: true },
  });
}

export async function revokeCourseAccess(id: string, courseId: string): Promise<void> {
  await apiFetch(`/admin/users/${id}/revoke-access`, {
    method: "POST",
    body: { course_id: courseId },
  });
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
    schema: SectionListResponseSchema,
  });
  return data?.data ?? [];
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

export async function listAnnouncements(limit = 50): Promise<AnnouncementListResponse> {
  const data = await apiFetch(`/admin/announcements?limit=${limit}`, {
    schema: AnnouncementListResponseSchema,
  });
  return data ?? { announcements: [] };
}

export async function deleteAnnouncement(id: string): Promise<void> {
  await apiFetch(`/admin/announcements/${id}`, {
    method: "DELETE",
  });
}

// ---- Teacher: My Courses (Phase 5) --------------------------------------------

export async function listMyCourses(params: {
  teacherId: string;
  limit?: number;
  after?: string;
}): Promise<TeacherCourseList> {
  const qs = new URLSearchParams();
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.after) qs.set("after", params.after);
  const data = await apiFetch(`/courses?${qs}`, { schema: TeacherCourseListSchema });
  return data!;
}

export async function createCourse(input: CreateCourseInput): Promise<Course> {
  const data = await apiFetch("/courses", { method: "POST", body: input, schema: CourseSchema });
  return data!;
}

export async function updateCourse(id: string, input: Partial<CreateCourseInput>): Promise<Course> {
  const data = await apiFetch(`/courses/${id}`, { method: "PUT", body: input, schema: CourseSchema });
  return data!;
}

export async function deleteCourse(id: string): Promise<void> {
  await apiFetch(`/courses/${id}`, { method: "DELETE" });
}

export async function publishCourse(id: string): Promise<void> {
  await apiFetch(`/courses/${id}/publish`, { method: "POST" });
}

// ---- Teacher: Earnings (Phase 5) ----------------------------------------------

export async function getMyEarnings(from: string, to: string): Promise<Earnings> {
  // Backend requires RFC3339; date-only strings (YYYY-MM-DD) are rejected.
  const fromRFC = from.includes("T") ? from : `${from}T00:00:00Z`;
  const toRFC = to.includes("T") ? to : `${to}T23:59:59Z`;
  const data = await apiFetch(
    `/me/earnings?from=${encodeURIComponent(fromRFC)}&to=${encodeURIComponent(toRFC)}`,
    { schema: EarningsSchema }
  );
  return data!;
}

export async function getCourseSales(courseId: string): Promise<CourseSales> {
  const data = await apiFetch(`/courses/${courseId}/sales`, { schema: CourseSalesSchema });
  return data!;
}

// ---- Teacher: Sections --------------------------------------------------------

export async function createSection(courseId: string, input: AddSectionInput): Promise<Section> {
  const data = await apiFetch(`/courses/${courseId}/sections`, {
    method: "POST",
    body: input,
    schema: SectionSchema,
  });
  return data!;
}

export async function updateSection(sectionId: string, input: Partial<AddSectionInput>): Promise<Section> {
  const data = await apiFetch(`/sections/${sectionId}`, {
    method: "PUT",
    body: input,
    schema: SectionSchema,
  });
  return data!;
}

export async function deleteSection(sectionId: string): Promise<void> {
  await apiFetch(`/sections/${sectionId}`, { method: "DELETE" });
}

// ---- Teacher: Lessons (Phase 5 Skill 6) ----------------------------------------

export async function listLessons(sectionId: string): Promise<Lesson[]> {
  const data = await apiFetch(`/sections/${sectionId}/lessons`, { schema: LessonListResponseSchema });
  return data?.data ?? [];
}

export async function addLesson(sectionId: string, input: AddLessonInput): Promise<Lesson> {
  const data = await apiFetch(`/sections/${sectionId}/lessons`, {
    method: "POST", body: input, schema: LessonSchema,
  });
  return data!;
}

export async function updateLesson(lessonId: string, input: Partial<AddLessonInput>): Promise<Lesson> {
  const data = await apiFetch(`/lessons/${lessonId}`, {
    method: "PUT", body: input, schema: LessonSchema,
  });
  return data!;
}

export async function deleteLesson(lessonId: string): Promise<void> {
  await apiFetch(`/lessons/${lessonId}`, { method: "DELETE" });
}

// ---- Teacher: Quizzes (Phase 5 Skill 6) ----------------------------------------

export async function listQuizzes(courseId: string): Promise<Quiz[]> {
  try {
    const raw = await apiFetch(`/courses/${courseId}/quizzes`);
    if (!raw) return [];
    const parsed = QuizListResponseSchema.safeParse(raw);
    return parsed.success ? parsed.data.data : [];
  } catch {
    return [];
  }
}

export async function createQuiz(courseId: string, input: CreateQuizInput): Promise<Quiz> {
  const data = await apiFetch(`/courses/${courseId}/quizzes`, {
    method: "POST", body: input, schema: QuizSchema,
  });
  return data!;
}

export async function updateQuiz(quizId: string, input: Partial<CreateQuizInput>): Promise<Quiz> {
  const data = await apiFetch(`/quizzes/${quizId}`, {
    method: "PUT", body: input, schema: QuizSchema,
  });
  return data!;
}

export async function deleteQuiz(quizId: string): Promise<void> {
  await apiFetch(`/quizzes/${quizId}`, { method: "DELETE" });
}

export async function addQuizQuestion(quizId: string, input: AddQuizQuestionInput): Promise<QuizQuestion> {
  const data = await apiFetch(`/quizzes/${quizId}/questions`, {
    method: "POST", body: input, schema: QuizQuestionSchema,
  });
  return data!;
}

// ---- Teacher: Exams (Phase 5 Skill 6) ------------------------------------------

export async function listExams(courseId: string): Promise<Exam[]> {
  try {
    const raw = await apiFetch(`/courses/${courseId}/exams`);
    if (!raw) return [];
    const parsed = ExamListResponseSchema.safeParse(raw);
    return parsed.success ? parsed.data.data : [];
  } catch {
    return [];
  }
}

export async function createExam(courseId: string, input: CreateExamInput): Promise<Exam> {
  const data = await apiFetch(`/courses/${courseId}/exams`, {
    method: "POST", body: input, schema: ExamSchema,
  });
  return data!;
}

export async function updateExam(examId: string, input: Partial<CreateExamInput>): Promise<Exam> {
  const data = await apiFetch(`/exams/${examId}`, {
    method: "PUT", body: input, schema: ExamSchema,
  });
  return data!;
}

export async function deleteExam(examId: string): Promise<void> {
  await apiFetch(`/exams/${examId}`, { method: "DELETE" });
}

export async function addExamQuestion(examId: string, input: ExamQuestionInput): Promise<void> {
  await apiFetch(`/exams/${examId}/questions`, { method: "POST", body: input });
}

export async function publishExamEndpoint(examId: string): Promise<void> {
  await apiFetch(`/exams/${examId}/publish`, { method: "POST" });
}

export async function getExamStats(examId: string): Promise<ExamStats> {
  const data = await apiFetch(`/exams/${examId}/stats`, { schema: ExamStatsSchema });
  return data!;
}

// ---- Assessments: standalone create, listing, per-question editing -------------

/** A page of quizzes across all courses, for the top-level Quizzes page. */
export async function listAllQuizzes(params: {
  limit?: number;
  after?: string;
  course_id?: string;
  standalone?: boolean;
  published?: boolean;
}): Promise<{ data: Quiz[]; next_cursor?: string | null }> {
  const qs = new URLSearchParams();
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.after) qs.set("after", params.after);
  if (params.course_id) qs.set("course_id", params.course_id);
  if (params.standalone !== undefined) qs.set("standalone", String(params.standalone));
  if (params.published !== undefined) qs.set("published", String(params.published));
  const data = await apiFetch(`/quizzes?${qs}`, { schema: QuizListResponseSchema });
  return data ?? { data: [] };
}

/** Creates a quiz with no course (admin) or under the course_id in the body. */
export async function createStandaloneQuiz(input: CreateQuizInput): Promise<Quiz> {
  const data = await apiFetch("/quizzes", { method: "POST", body: input, schema: QuizSchema });
  return data!;
}

/** The quiz and its questions WITH answer keys, for the editor. */
export async function getQuizForAuthoring(quizId: string): Promise<AuthoredQuiz> {
  const data = await apiFetch(`/quizzes/${quizId}/questions`, { schema: AuthoredQuizSchema });
  return data!;
}

export async function updateQuizQuestion(
  questionId: string,
  input: UpdateQuestionInput
): Promise<AuthoredQuestion> {
  const data = await apiFetch(`/questions/${questionId}`, {
    method: "PUT",
    body: input,
    schema: AuthoredQuestionSchema,
  });
  return data!;
}

export async function deleteQuizQuestion(questionId: string): Promise<void> {
  await apiFetch(`/questions/${questionId}`, { method: "DELETE" });
}

export async function reorderQuizQuestions(quizId: string, questionIds: string[]): Promise<void> {
  await apiFetch(`/quizzes/${quizId}/questions/order`, {
    method: "PUT",
    body: { question_ids: questionIds },
  });
}

/** A page of exams including the caller's unpublished drafts. */
export async function listAllExams(params: {
  limit?: number;
  after?: string;
  subject?: string;
  grade?: number;
  published?: boolean;
}): Promise<{ data: Exam[]; next_cursor?: string | null }> {
  const qs = new URLSearchParams();
  if (params.limit) qs.set("limit", String(params.limit));
  if (params.after) qs.set("after", params.after);
  if (params.subject) qs.set("subject", params.subject);
  if (params.grade) qs.set("grade", String(params.grade));
  if (params.published !== undefined) qs.set("published", String(params.published));
  const data = await apiFetch(`/exams?${qs}`, { schema: ExamListResponseSchema });
  return data ?? { data: [] };
}

/** Creates an exam with no course (admin) or under the course_id in the body. */
export async function createStandaloneExam(input: CreateExamInput): Promise<Exam> {
  const data = await apiFetch("/exams", { method: "POST", body: input, schema: ExamSchema });
  return data!;
}

/** The exam and its questions WITH answer keys, for the editor. */
export async function getExamForAuthoring(examId: string): Promise<AuthoredExam> {
  const data = await apiFetch(`/exams/${examId}/questions`, { schema: AuthoredExamSchema });
  return data!;
}

export async function updateExamQuestion(
  examId: string,
  questionId: string,
  input: UpdateExamQuestionInput
): Promise<AuthoredQuestion> {
  const data = await apiFetch(`/exams/${examId}/questions/${questionId}`, {
    method: "PUT",
    body: input,
    schema: AuthoredQuestionSchema,
  });
  return data!;
}

export async function removeExamQuestion(examId: string, questionId: string): Promise<Exam> {
  const data = await apiFetch(`/exams/${examId}/questions/${questionId}`, {
    method: "DELETE",
    schema: ExamSchema,
  });
  return data!;
}

export async function reorderExamQuestions(examId: string, questionIds: string[]): Promise<void> {
  await apiFetch(`/exams/${examId}/questions/order`, {
    method: "PUT",
    body: { question_ids: questionIds },
  });
}

export async function unpublishExamEndpoint(examId: string): Promise<void> {
  await apiFetch(`/exams/${examId}/unpublish`, { method: "POST" });
}

/** Bulk-creates a whole quiz or exam from one import document, atomically. */
export async function importAssessment(doc: ImportDocument): Promise<ImportResult> {
  const path = doc.kind === "exam" ? "/exams/import" : "/quizzes/import";
  const data = await apiFetch(path, { method: "POST", body: doc, schema: ImportResultSchema });
  return data!;
}

// ---- Teacher: Video pipeline (Phase 5 Skill 6) ----------------------------------

export async function requestVideoUpload(lessonId: string, input: {
  file_name: string; content_type: string; size_bytes: number;
}): Promise<VideoUploadResponse> {
  const data = await apiFetch(`/lessons/${lessonId}/videos/upload-url`, {
    method: "POST", body: input, schema: VideoUploadResponseSchema,
  });
  return data!;
}

export async function getVideoStatus(videoId: string): Promise<VideoStatus> {
  const data = await apiFetch(`/videos/${videoId}/status`, { schema: VideoStatusSchema });
  return data!;
}

export async function confirmVideoUpload(videoId: string): Promise<void> {
  await apiFetch(`/videos/${videoId}/confirm`, { method: "POST" });
}

export async function retryTranscode(videoId: string): Promise<void> {
  await apiFetch(`/videos/${videoId}/retry`, { method: "POST" });
}

// ---- Course Access Requests (Admin & Teacher) -------------------------------

export async function listAdminCourseRequests(params?: {
  status?: string;
  page?: number;
  limit?: number;
}): Promise<CourseAccessRequestListResponse> {
  const qs = new URLSearchParams();
  if (params?.status) qs.set("status", params.status);
  if (params?.page) qs.set("page", String(params.page));
  if (params?.limit) qs.set("limit", String(params.limit));
  const data = await apiFetch(`/admin/course-requests?${qs}`, {
    schema: CourseAccessRequestListResponseSchema,
  });
  return data!;
}

export async function approveAdminCourseRequest(requestId: string): Promise<void> {
  await apiFetch(`/admin/course-requests/${requestId}/approve`, {
    method: "POST",
  });
}

export async function rejectAdminCourseRequest(requestId: string, reason?: string): Promise<void> {
  await apiFetch(`/admin/course-requests/${requestId}/reject`, {
    method: "POST",
    body: { reason },
  });
}

export async function listTeacherCourseRequests(params?: {
  course_id?: string;
  status?: string;
  page?: number;
  limit?: number;
}): Promise<CourseAccessRequestListResponse> {
  const qs = new URLSearchParams();
  if (params?.course_id) qs.set("course_id", params.course_id);
  if (params?.status) qs.set("status", params.status);
  if (params?.page) qs.set("page", String(params.page));
  if (params?.limit) qs.set("limit", String(params.limit));
  const data = await apiFetch(`/teacher/course-requests?${qs}`, {
    schema: CourseAccessRequestListResponseSchema,
  });
  return data!;
}

export async function approveTeacherCourseRequest(requestId: string): Promise<void> {
  await apiFetch(`/teacher/course-requests/${requestId}/approve`, {
    method: "POST",
  });
}

export async function rejectTeacherCourseRequest(requestId: string, reason?: string): Promise<void> {
  await apiFetch(`/teacher/course-requests/${requestId}/reject`, {
    method: "POST",
    body: { reason },
  });
}

