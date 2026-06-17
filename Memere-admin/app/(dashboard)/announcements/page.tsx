import { AnnouncementsClient } from "./announcements-client";

export default function AnnouncementsPage() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Announcements</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Broadcast a push notification to a user segment.
        </p>
      </div>
      <AnnouncementsClient />
    </div>
  );
}
