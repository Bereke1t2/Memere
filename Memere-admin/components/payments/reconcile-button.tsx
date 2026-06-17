"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { clientAction } from "@/lib/client-action";
import { RefreshCw } from "lucide-react";

export function ReconcileButton() {
  const queryClient = useQueryClient();

  const reconcileMutation = useMutation({
    mutationFn: () =>
      clientAction<{ reconciled: number }>("/api/admin/payments/reconcile"),
    onSuccess: (data) => {
      toast.success(
        data.reconciled === 0
          ? "No pending payments to reconcile."
          : `Reconciled ${data.reconciled} pending payment${data.reconciled !== 1 ? "s" : ""}.`
      );
      queryClient.invalidateQueries({ queryKey: ["payments"] });
    },
    onError: (err: Error) => toast.error(err.message),
  });

  return (
    <Button
      variant="outline"
      size="sm"
      disabled={reconcileMutation.isPending}
      onClick={() => reconcileMutation.mutate()}
    >
      <RefreshCw className={`h-3.5 w-3.5 mr-1.5 ${reconcileMutation.isPending ? "animate-spin" : ""}`} />
      {reconcileMutation.isPending ? "Reconciling…" : "Reconcile pending"}
    </Button>
  );
}
