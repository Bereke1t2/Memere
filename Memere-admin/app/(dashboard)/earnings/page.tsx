import { redirect } from "next/navigation";
import { requireStaff } from "@/lib/auth/session";

export default async function EarningsPage() {
  const { user } = await requireStaff();
  if (user.role !== "teacher") redirect(user.role === "admin" ? "/revenue" : "/");

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Earnings</h1>
        <p className="text-sm text-muted-foreground mt-1">Your earnings overview.</p>
      </div>
      <p className="text-sm text-muted-foreground">Coming soon — skill 4.</p>
    </div>
  );
}
