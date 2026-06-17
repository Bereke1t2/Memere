import { redirect } from "next/navigation";
import { requireStaff } from "@/lib/auth/session";

export default async function MyCoursesPage() {
  const { user } = await requireStaff();
  if (user.role !== "teacher") redirect(user.role === "admin" ? "/courses" : "/");

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">My Courses</h1>
        <p className="text-sm text-muted-foreground mt-1">Manage your courses.</p>
      </div>
      <p className="text-sm text-muted-foreground">Coming soon — skill 3.</p>
    </div>
  );
}
