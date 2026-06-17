import { redirect } from "next/navigation";
import { requireStaff } from "@/lib/auth/session";
import { listMyCourses } from "@/lib/api/endpoints";
import { MyCoursesClient } from "./my-courses-client";

export default async function MyCoursesPage() {
  const { user } = await requireStaff();
  if (user.role !== "teacher") redirect(user.role === "admin" ? "/courses" : "/");

  const result = await listMyCourses({ teacherId: user.id, limit: 20 });

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">My Courses</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Create and manage your courses.
        </p>
      </div>
      <MyCoursesClient initialData={{ items: result.data, next: result.next_cursor }} />
    </div>
  );
}
