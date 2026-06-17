import { Skeleton } from "@/components/ui/skeleton";

export default function CourseDetailLoading() {
  return (
    <div className="flex flex-col gap-6 max-w-2xl">
      <div className="flex items-start justify-between gap-4">
        <div className="flex flex-col gap-2">
          <Skeleton className="h-8 w-64" />
          <Skeleton className="h-4 w-80" />
        </div>
        <Skeleton className="h-8 w-24" />
      </div>
      <Skeleton className="h-56 rounded-lg" />
      <Skeleton className="h-24 rounded-lg" />
      <div className="flex flex-col gap-2">
        <Skeleton className="h-6 w-24" />
        {[0, 1, 2].map((i) => (
          <Skeleton key={i} className="h-11 rounded-md" />
        ))}
      </div>
    </div>
  );
}
