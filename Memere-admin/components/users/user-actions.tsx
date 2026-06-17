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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { clientAction } from "@/lib/client-action";
import type { User } from "@/lib/api/schemas";

type ActiveDialog = null | "suspend" | "reactivate" | "role";

interface UserActionsProps {
  user: User;
}

export function UserActions({ user }: UserActionsProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [dialog, setDialog] = useState<ActiveDialog>(null);
  const [reason, setReason] = useState("");
  const [newRole, setNewRole] = useState(user.role);

  function close() {
    setDialog(null);
    setReason("");
  }

  function afterSuccess(message: string) {
    toast.success(message);
    close();
    router.refresh();
    queryClient.invalidateQueries({ queryKey: ["users"] });
  }

  const suspendMutation = useMutation({
    mutationFn: () =>
      clientAction(`/api/admin/users/${user.id}/suspend`, { reason }),
    onSuccess: () =>
      afterSuccess(`${user.first_name} ${user.last_name} has been suspended.`),
    onError: (err: Error) => toast.error(err.message),
  });

  const reactivateMutation = useMutation({
    mutationFn: () =>
      clientAction(`/api/admin/users/${user.id}/reactivate`),
    onSuccess: () =>
      afterSuccess(`${user.first_name} ${user.last_name} has been reactivated.`),
    onError: (err: Error) => toast.error(err.message),
  });

  const roleMutation = useMutation({
    mutationFn: () =>
      clientAction(`/api/admin/users/${user.id}/role`, { role: newRole }),
    onSuccess: () =>
      afterSuccess(`Role changed to ${newRole}.`),
    onError: (err: Error) => toast.error(err.message),
  });

  return (
    <>
      <div className="flex flex-wrap gap-2">
        {user.is_active ? (
          <Button
            variant="destructive"
            size="sm"
            onClick={() => setDialog("suspend")}
          >
            Suspend user
          </Button>
        ) : (
          <Button
            variant="outline"
            size="sm"
            onClick={() => setDialog("reactivate")}
          >
            Reactivate user
          </Button>
        )}
        <Button variant="outline" size="sm" onClick={() => setDialog("role")}>
          Change role
        </Button>
      </div>

      {/* Suspend dialog */}
      <Dialog open={dialog === "suspend"} onOpenChange={(o) => !o && close()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Suspend user</DialogTitle>
            <DialogDescription>
              Provide a reason for suspending{" "}
              <strong>
                {user.first_name} {user.last_name}
              </strong>
              . This is required.
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-2 py-2">
            <Label htmlFor="reason">Reason</Label>
            <textarea
              id="reason"
              className="min-h-[80px] w-full rounded-md border bg-background px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
              placeholder="Enter suspension reason…"
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
              disabled={!reason.trim() || suspendMutation.isPending}
              onClick={() => suspendMutation.mutate()}
            >
              {suspendMutation.isPending ? "Suspending…" : "Confirm suspension"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Reactivate dialog */}
      <Dialog open={dialog === "reactivate"} onOpenChange={(o) => !o && close()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reactivate user</DialogTitle>
            <DialogDescription>
              Reactivate{" "}
              <strong>
                {user.first_name} {user.last_name}
              </strong>
              ? They will regain full access to the platform.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={close}>
              Cancel
            </Button>
            <Button
              disabled={reactivateMutation.isPending}
              onClick={() => reactivateMutation.mutate()}
            >
              {reactivateMutation.isPending ? "Reactivating…" : "Confirm reactivation"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Change role dialog */}
      <Dialog open={dialog === "role"} onOpenChange={(o) => !o && close()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Change role</DialogTitle>
            <DialogDescription>
              Change the role for{" "}
              <strong>
                {user.first_name} {user.last_name}
              </strong>
              .
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-2 py-2">
            <Label>New role</Label>
            <Select value={newRole} onValueChange={setNewRole}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="student">Student</SelectItem>
                <SelectItem value="teacher">Teacher</SelectItem>
                <SelectItem value="admin">Admin</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={close}>
              Cancel
            </Button>
            <Button
              disabled={newRole === user.role || roleMutation.isPending}
              onClick={() => roleMutation.mutate()}
            >
              {roleMutation.isPending ? "Saving…" : "Save role"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
