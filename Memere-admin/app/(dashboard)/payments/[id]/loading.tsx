import { Skeleton } from "@/components/ui/skeleton";

export default function PaymentDetailLoading() {
  return (
    <div className="flex flex-col gap-6 max-w-2xl">
      <div className="flex items-start justify-between gap-4">
        <div className="flex flex-col gap-2">
          <Skeleton className="h-8 w-56" />
          <Skeleton className="h-4 w-40" />
        </div>
        <Skeleton className="h-8 w-20" />
      </div>
      <Skeleton className="h-56 rounded-lg" />
    </div>
  );
}
