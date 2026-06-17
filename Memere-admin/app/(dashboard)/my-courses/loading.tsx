import { Skeleton } from "@/components/ui/skeleton";

export default function MyCoursesLoading() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <Skeleton className="h-8 w-40" />
        <Skeleton className="h-9 w-32" />
      </div>
      <div className="rounded-xl border overflow-hidden">
        <Skeleton className="h-10 w-full" />
        {[0, 1, 2, 3, 4].map((i) => (
          <Skeleton key={i} className="h-14 w-full mt-px" />
        ))}
      </div>
    </div>
  );
}
