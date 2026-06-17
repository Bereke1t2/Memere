import { listPayments } from "@/lib/api/endpoints";
import { PaymentsClient } from "./payments-client";

export default async function PaymentsPage({
  searchParams,
}: {
  searchParams: Promise<{ after?: string; status?: string; limit?: string }>;
}) {
  const params = await searchParams;
  const limit = Math.min(parseInt(params.limit ?? "20", 10), 100);
  const after = params.after;
  const status = params.status;

  const result = await listPayments({ limit, after, status });

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Payments</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Browse and audit platform transactions.
        </p>
      </div>
      <PaymentsClient initialData={{ items: result.payments, next: result.next }} />
    </div>
  );
}
