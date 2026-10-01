import { redirect } from "next/navigation";
import { requireStaff } from "@/lib/auth/session";
import { listAnnouncementsFromDB } from "@/lib/api/announcements-db";
import { AnnouncementsClient } from "./announcements-client";
import type { Announcement } from "@/lib/api/schemas";

export const dynamic = "force-dynamic";

export default async function AnnouncementsPage() {
  const { user } = await requireStaff();
  if (user.role !== "admin") redirect("/");

  let initialAnnouncements: Announcement[] = [];
  try {
    const res = await listAnnouncementsFromDB(50);
    initialAnnouncements = res?.announcements ?? [];
  } catch (err) {
    console.error("[AnnouncementsPage] failed to fetch initial announcements:", err);
  }

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Announcements</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Broadcast and manage push notifications and in-app announcements.
        </p>
      </div>
      <AnnouncementsClient initialData={initialAnnouncements} />
    </div>
  );
}
