"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { clientAction } from "@/lib/client-action";
import { formatMoney } from "@/lib/format";
import type { AdminPayment } from "@/lib/api/schemas";

interface PaymentActionsProps {
  payment: AdminPayment;
}

export function PaymentActions({ payment }: PaymentActionsProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);

  const refundMutation = useMutation({
    mutationFn: () => clientAction(`/api/admin/payments/${payment.id}/refund`),
    onSuccess: () => {
      toast.success("Refund processed successfully.");
      setOpen(false);
      router.refresh();
      queryClient.invalidateQueries({ queryKey: ["payments"] });
    },
    onError: (err: Error) => toast.error(err.message),
  });

  if (payment.status !== "completed") return null;

  const formatted = formatMoney(payment.amount, payment.currency);

  return (
    <>
      <Button variant="destructive" size="sm" onClick={() => setOpen(true)}>
        Refund
      </Button>

      <Dialog open={open} onOpenChange={(o) => !o && setOpen(false)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm refund</DialogTitle>
            <DialogDescription>
              This refunds <strong>{formatted}</strong> and{" "}
              <strong>cannot be undone</strong>. The payment will be marked as
              refunded and the student will be notified.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={refundMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              disabled={refundMutation.isPending}
              onClick={() => refundMutation.mutate()}
            >
              {refundMutation.isPending ? "Processing…" : `Refund ${formatted}`}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
