import { redirect } from "next/navigation";
import { requireStaff } from "@/lib/auth/session";
import { getMyEarnings, listMyCourses, getCourseSales } from "@/lib/api/endpoints";
import { EarningsDashboard } from "@/components/earnings/earnings-dashboard";

function toISO(d: Date) {
  return d.toISOString().split("T")[0];
}

function resolveRange(p: { from?: string; to?: string }) {
  return {
    to: p.to ?? toISO(new Date()),
    from: p.from ?? toISO(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)),
  };
}

export default async function EarningsPage({
  searchParams,
}: {
  searchParams: Promise<{ from?: string; to?: string }>;
}) {
  const { user } = await requireStaff();
  if (user.role !== "teacher") redirect(user.role === "admin" ? "/revenue" : "/");

  const params = await searchParams;
  const { from, to } = resolveRange(params);

  const [earnings, courseList] = await Promise.all([
    getMyEarnings(from, to),
    listMyCourses({ teacherId: user.id, limit: 50 }),
  ]);

  const courses = courseList.data;

  const sales = await Promise.all(
    courses.map((c) => getCourseSales(c.id).catch(() => ({ course_id: c.id, gross: "0", units: 0 })))
  );

  return <EarningsDashboard earnings={earnings} courses={courses} sales={sales} from={from} to={to} />;
}
