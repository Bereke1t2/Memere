"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Check, X, ShieldAlert, KeyRound } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/format";
import type { CourseAccessRequest } from "@/lib/api/schemas";

interface CourseRequestsClientProps {
  role: "admin" | "teacher";
}

export function CourseRequestsClient({ role }: CourseRequestsClientProps) {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState<string>("pending");
  const [rejectDialogItem, setRejectDialogItem] = useState<CourseAccessRequest | null>(null);
  const [rejectReason, setRejectReason] = useState("");

  const apiBase = role === "teacher" ? "/api/teacher/course-requests" : "/api/admin/course-requests";

  const { data, isLoading, isError, error } = useQuery<{
    requests: CourseAccessRequest[];
    total: number;
  }>({
    queryKey: ["course-requests", role, statusFilter],
    queryFn: async () => {
      const qs = new URLSearchParams();
      if (statusFilter && statusFilter !== "all") {
        qs.set("status", statusFilter);
      }
      const res = await fetch(`${apiBase}?${qs}`);
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to load requests");
      }
      return res.json();
    },
  });

  const approveMutation = useMutation({
    mutationFn: async (requestId: string) => {
      const res = await fetch(`${apiBase}/${requestId}/approve`, { method: "POST" });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to approve request");
      }
      return res.json();
    },
    onSuccess: () => {
      toast.success("Course access request approved.");
      queryClient.invalidateQueries({ queryKey: ["course-requests"] });
    },
    onError: (err: Error) => toast.error(err.message),
  });

  const rejectMutation = useMutation({
    mutationFn: async ({ requestId, reason }: { requestId: string; reason: string }) => {
      const res = await fetch(`${apiBase}/${requestId}/reject`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ reason: reason || undefined }),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to reject request");
      }
      return res.json();
    },
    onSuccess: () => {
      toast.success("Course access request rejected.");
      setRejectDialogItem(null);
      setRejectReason("");
      queryClient.invalidateQueries({ queryKey: ["course-requests"] });
    },
    onError: (err: Error) => toast.error(err.message),
  });

  const requests = data?.requests ?? [];

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between gap-2 flex-wrap">
        <div className="flex items-center gap-2">
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Requests</SelectItem>
              <SelectItem value="pending">Pending</SelectItem>
              <SelectItem value="approved">Approved</SelectItem>
              <SelectItem value="rejected">Rejected</SelectItem>
            </SelectContent>
          </Select>
        </div>
        {data && (
          <p className="text-xs text-muted-foreground">
            Showing {requests.length} {data.total !== undefined ? `of ${data.total} ` : ""}request{requests.length !== 1 ? "s" : ""}
          </p>
        )}
      </div>

      {isLoading ? (
        <div className="space-y-2">
          <Skeleton className="h-12 w-full" />
          <Skeleton className="h-12 w-full" />
          <Skeleton className="h-12 w-full" />
        </div>
      ) : isError ? (
        <div className="rounded-lg border border-destructive/30 p-4 text-sm text-destructive">
          Failed to load requests: {(error as Error).message}
        </div>
      ) : requests.length === 0 ? (
        <div className="rounded-lg border p-8 text-center text-sm text-muted-foreground">
          No course access requests found for the selected filter.
        </div>
      ) : (
        <div className="rounded-lg border overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-muted/50 border-b text-xs text-muted-foreground font-medium text-left">
              <tr>
                <th className="p-3">Student</th>
                <th className="p-3">Course</th>
                <th className="p-3">Status</th>
                <th className="p-3">Requested At</th>
                <th className="p-3 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {requests.map((req) => (
                <tr key={req.id} className="hover:bg-muted/30 transition-colors">
                  <td className="p-3">
                    <div className="font-medium">{req.student_name || "Student"}</div>
                    <div className="text-xs text-muted-foreground">{req.student_email}</div>
                  </td>
                  <td className="p-3 font-medium">{req.course_title || req.course_id}</td>
                  <td className="p-3">
                    {req.status === "pending" ? (
                      <Badge variant="outline" className="text-amber-700 bg-amber-50 border-amber-200">
                        Pending
                      </Badge>
                    ) : req.status === "approved" ? (
                      <Badge variant="secondary" className="text-green-700 bg-green-100">
                        Approved
                      </Badge>
                    ) : (
                      <Badge variant="destructive">
                        Rejected
                      </Badge>
                    )}
                  </td>
                  <td className="p-3 text-xs text-muted-foreground tabular-nums">
                    {formatDate(req.created_at)}
                  </td>
                  <td className="p-3 text-right">
                    {req.status === "pending" ? (
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          size="sm"
                          variant="default"
                          className="h-8 bg-green-600 hover:bg-green-700 text-white"
                          disabled={approveMutation.isPending}
                          onClick={() => approveMutation.mutate(req.id)}
                        >
                          <Check className="h-3.5 w-3.5 mr-1" />
                          Approve
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          className="h-8 text-destructive hover:bg-destructive/10"
                          disabled={rejectMutation.isPending}
                          onClick={() => setRejectDialogItem(req)}
                        >
                          <X className="h-3.5 w-3.5 mr-1" />
                          Reject
                        </Button>
                      </div>
                    ) : (
                      <span className="text-xs text-muted-foreground">
                        {req.status === "approved" ? "Access active" : req.rejection_reason || "Rejected"}
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Rejection Dialog */}
      <Dialog open={!!rejectDialogItem} onOpenChange={(o) => !o && setRejectDialogItem(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reject Course Access Request</DialogTitle>
            <DialogDescription>
              Reject access request from <strong>{rejectDialogItem?.student_name || "student"}</strong> for course{" "}
              <strong>{rejectDialogItem?.course_title}</strong>.
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-2 py-2">
            <Label htmlFor="req-reason">Rejection Reason (optional)</Label>
            <textarea
              id="req-reason"
              className="min-h-[80px] w-full rounded-md border bg-background px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="e.g. Prerequisites not met, please complete Introduction first."
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setRejectDialogItem(null)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              disabled={rejectMutation.isPending}
              onClick={() => {
                if (rejectDialogItem) {
                  rejectMutation.mutate({ requestId: rejectDialogItem.id, reason: rejectReason });
                }
              }}
            >
              {rejectMutation.isPending ? "Rejecting…" : "Confirm Rejection"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
