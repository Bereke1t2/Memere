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
import { Label } from "@/components/ui/label";
import { clientAction } from "@/lib/client-action";
import type { Course } from "@/lib/api/schemas";

interface CourseActionsProps {
  course: Course;
}

export function CourseActions({ course }: CourseActionsProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [reason, setReason] = useState("");

  function close() {
    setOpen(false);
    setReason("");
  }

  const unpublishMutation = useMutation({
    mutationFn: () =>
      clientAction(`/api/admin/courses/${course.id}/unpublish`, { reason }),
    onSuccess: () => {
      toast.success(`"${course.title}" has been unpublished.`);
      close();
      router.refresh();
      queryClient.invalidateQueries({ queryKey: ["courses"] });
    },
    onError: (err: Error) => toast.error(err.message),
  });

  if (!course.is_published) return null;

  return (
    <>
      <Button variant="destructive" size="sm" onClick={() => setOpen(true)}>
        Unpublish
      </Button>

      <Dialog open={open} onOpenChange={(o) => !o && close()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Unpublish course</DialogTitle>
            <DialogDescription>
              Provide a reason for unpublishing{" "}
              <strong>&ldquo;{course.title}&rdquo;</strong>. This is required.
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-2 py-2">
            <Label htmlFor="reason">Reason</Label>
            <textarea
              id="reason"
              className="min-h-[80px] w-full rounded-md border bg-background px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="Enter reason for unpublishing…"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={close}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              disabled={!reason.trim() || unpublishMutation.isPending}
              onClick={() => unpublishMutation.mutate()}
            >
              {unpublishMutation.isPending ? "Unpublishing…" : "Confirm unpublish"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
