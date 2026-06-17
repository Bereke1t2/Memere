import { Skeleton } from "@/components/ui/skeleton";

export default function AnnouncementsLoading() {
  return (
    <div className="flex flex-col gap-6 max-w-xl">
      <div className="flex flex-col gap-1">
        <Skeleton className="h-8 w-40" />
        <Skeleton className="h-4 w-72" />
      </div>
      <Skeleton className="h-10 w-full" />
      <Skeleton className="h-32 w-full" />
      <Skeleton className="h-10 w-48" />
      <Skeleton className="h-10 w-36" />
    </div>
  );
}
