import { notFound } from "next/navigation";
import Link from "next/link";
import { getPayment } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/errors";
import { Badge } from "@/components/ui/badge";
import { PaymentActions } from "@/components/payments/payment-actions";
import { BreadcrumbSetter } from "@/lib/breadcrumb-context";
import { formatDate, formatMoney } from "@/lib/format";

function Row({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center gap-1 py-3 border-b last:border-0">
      <dt className="w-44 shrink-0 text-sm text-muted-foreground">{label}</dt>
      <dd className="text-sm">{children}</dd>
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
  switch (status) {
    case "completed":
      return (
        <Badge variant="secondary" className="text-green-700 bg-green-100 dark:text-green-400 dark:bg-green-950">
          Completed
        </Badge>
      );
    case "pending":
      return (
        <Badge variant="secondary" className="text-amber-700 bg-amber-100 dark:text-amber-400 dark:bg-amber-950">
          Pending
        </Badge>
      );
    case "failed":
      return <Badge variant="destructive">Failed</Badge>;
    case "refunded":
      return <Badge variant="secondary">Refunded</Badge>;
    default:
      return <Badge variant="outline" className="capitalize">{status}</Badge>;
  }
}

export default async function PaymentDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let payment;
  try {
    payment = await getPayment(id);
  } catch (err) {
    if (
      err instanceof ApiError &&
      (err.status === 404 || err.code === "NOT_FOUND" || err.code === "RESOURCE_NOT_FOUND")
    ) {
      notFound();
    }
    throw err;
  }

  const shortId = payment.id.slice(0, 8);

  return (
    <>
      <BreadcrumbSetter label={`Payment ${shortId}…`} />

      <div className="flex flex-col gap-6 max-w-2xl">
        <div className="flex items-start justify-between gap-4 flex-wrap">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">
              Payment <span className="font-mono text-lg">{shortId}…</span>
            </h1>
            <p className="text-sm text-muted-foreground mt-0.5">
              {formatMoney(payment.amount, payment.currency)} via{" "}
              <span className="capitalize">{payment.provider}</span>
            </p>
          </div>
          <PaymentActions payment={payment} />
        </div>

        <div className="rounded-lg border">
          <dl>
            <Row label="Status">
              <StatusBadge status={payment.status} />
            </Row>
            <Row label="Amount">
              <span className="font-medium tabular-nums">
                {formatMoney(payment.amount, payment.currency)}
              </span>
            </Row>
            <Row label="Currency">{payment.currency}</Row>
            <Row label="Provider">
              <span className="capitalize">{payment.provider}</span>
            </Row>
            {payment.provider_txn_id && (
              <Row label="Provider Txn ID">
                <span className="font-mono text-xs break-all">
                  {payment.provider_txn_id}
                </span>
              </Row>
            )}
            <Row label="Student">
              <Link
                href={`/users/${payment.student_id}`}
                className="font-mono text-xs hover:underline"
              >
                {payment.student_id.slice(0, 8)}&hellip;
              </Link>
            </Row>
            <Row label="Payment ID">
              <span className="font-mono text-xs break-all">{payment.id}</span>
            </Row>
            <Row label="Date">{formatDate(payment.created_at)}</Row>
          </dl>
        </div>
      </div>
    </>
  );
}
