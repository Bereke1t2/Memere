import { notFound } from "next/navigation";
import { getUser } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { Badge } from "@/components/ui/badge";
import { UserActions } from "@/components/users/user-actions";
import { UserCourseAccess } from "@/components/users/user-course-access";
import { BreadcrumbSetter } from "@/lib/breadcrumb-context";
import { formatDate } from "@/lib/format";

function ProfileRow({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center gap-1 py-3 border-b last:border-0">
      <dt className="w-40 shrink-0 text-sm text-muted-foreground">{label}</dt>
      <dd className="text-sm">{children}</dd>
    </div>
  );
}

export default async function UserDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let user;
  try {
    user = await getUser(id);
  } catch (err) {
    if (err instanceof ApiError && (err.status === 404 || err.code === "NOT_FOUND" || err.code === "RESOURCE_NOT_FOUND")) {
      notFound();
    }
    throw err;
  }

  return (
    <>
    <BreadcrumbSetter label={`${user.first_name} ${user.last_name}`} />
    <div className="flex flex-col gap-6 max-w-4xl">
      <div className="flex items-start justify-between gap-4 flex-wrap">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">
            {user.first_name} {user.last_name}
          </h1>
          <p className="text-sm text-muted-foreground mt-0.5">{user.email}</p>
        </div>
        <UserActions user={user} />
      </div>

      <div className="rounded-lg border bg-card p-4">
        <dl>
          <ProfileRow label="Approval Status">
            {user.approval_status === "pending" ? (
              <Badge variant="outline" className="text-amber-700 bg-amber-50 border-amber-200 dark:text-amber-400 dark:bg-amber-950 dark:border-amber-800">
                Pending Approval
              </Badge>
            ) : user.approval_status === "rejected" ? (
              <Badge variant="destructive">Rejected</Badge>
            ) : (
              <Badge variant="secondary" className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950">
                Approved
              </Badge>
            )}
          </ProfileRow>
          <ProfileRow label="Account Status">
            {user.is_active ? (
              <Badge variant="secondary" className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950">
                Active
              </Badge>
            ) : (
              <Badge variant="destructive">Suspended</Badge>
            )}
          </ProfileRow>
          <ProfileRow label="Role">
            <Badge
              variant={
                user.role === "admin"
                  ? "default"
                  : user.role === "teacher"
                    ? "secondary"
                    : "outline"
              }
              className="capitalize"
            >
              {user.role}
            </Badge>
          </ProfileRow>
          <ProfileRow label="Email verified">
            {user.is_email_verified ? (
              <span className="text-green-700 dark:text-green-400">Yes</span>
            ) : (
              <span className="text-muted-foreground">No</span>
            )}
          </ProfileRow>
          <ProfileRow label="Phone">
            {user.phone ?? <span className="text-muted-foreground">—</span>}
          </ProfileRow>
          <ProfileRow label="Joined">
            {formatDate(user.created_at)}
          </ProfileRow>
          <ProfileRow label="Last login">
            {formatDate(user.last_login_at)}
          </ProfileRow>
        </dl>
      </div>

      {/* Course Access Management (Admins can grant/revoke courses directly) */}
      <UserCourseAccess userId={user.id} />
    </div>
    </>
  );
}
