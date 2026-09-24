import { requireStaff } from "@/lib/auth/session";
import { CourseRequestsClient } from "./course-requests-client";

export default async function CourseRequestsPage() {
  const { user } = await requireStaff();
  const role = user.role === "admin" ? "admin" : "teacher";

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Course Access Requests</h1>
        <p className="text-sm text-muted-foreground mt-1">
          {role === "admin"
            ? "Review and approve/reject student access requests across all courses."
            : "Review and approve/reject student access requests for your courses."}
        </p>
      </div>
      <CourseRequestsClient role={role} />
    </div>
  );
}
