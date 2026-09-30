"use client";

import { useState } from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { clientAction } from "@/lib/client-action";
import { Plus, Trash2, Megaphone, Users, RefreshCw, AlertTriangle, Send } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import type { Announcement, AnnouncementListResponse } from "@/lib/api/schemas";

const SEGMENTS = [
  { value: "students", label: "Students" },
  { value: "teachers", label: "Teachers" },
  { value: "subscribers", label: "Subscribers" },
  { value: "all", label: "All users" },
] as const;

const schema = z.object({
  title: z.string().min(1, "Title is required"),
  body: z.string().min(1, "Body is required"),
  segment: z.enum(["all", "students", "teachers", "subscribers"]),
  kvPairs: z.array(z.object({ key: z.string(), value: z.string() })).optional(),
});

type FormValues = z.infer<typeof schema>;

function formatDateTime(iso: string | null | undefined): string {
  if (!iso) return "—";
  try {
    const d = new Date(iso);
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "numeric",
      minute: "2-digit",
    }).format(d);
  } catch {
    return iso;
  }
}

interface AnnouncementsClientProps {
  initialData?: Announcement[];
}

export function AnnouncementsClient({ initialData = [] }: AnnouncementsClientProps) {
  const queryClient = useQueryClient();
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [allConfirmed, setAllConfirmed] = useState(false);
  const [sending, setSending] = useState(false);
  const [pending, setPending] = useState<FormValues | null>(null);

  // Deletion dialog state
  const [deleteItem, setDeleteItem] = useState<Announcement | null>(null);

  const {
    register,
    handleSubmit,
    control,
    reset,
    watch,
    setValue,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { segment: "students", kvPairs: [] },
  });

  const { fields, append, remove } = useFieldArray({ control, name: "kvPairs" });
  // eslint-disable-next-line react-hooks/incompatible-library
  const segment = watch("segment");

  // Fetch announcements
  const {
    data,
    isLoading,
    isError,
    error,
    refetch,
    isFetching,
  } = useQuery<AnnouncementListResponse>({
    queryKey: ["admin-announcements"],
    initialData: { announcements: initialData },
    queryFn: async () => {
      const res = await fetch("/api/admin/announcements");
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to load announcements");
      }
      return res.json();
    },
  });

  const announcements = Array.isArray(data?.announcements)
    ? data.announcements
    : [];

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await clientAction(`/api/admin/announcements/${id}`, undefined, "DELETE");
    },
    onSuccess: () => {
      toast.success("Announcement deleted from all user notification feeds.");
      setDeleteItem(null);
      queryClient.invalidateQueries({ queryKey: ["admin-announcements"] });
    },
    onError: (err: Error) => {
      toast.error(err.message || "Failed to delete announcement.");
    },
  });

  function onSubmit(values: FormValues) {
    setPending(values);
    setAllConfirmed(false);
    setConfirmOpen(true);
  }

  async function confirmSend() {
    if (!pending) return;
    setSending(true);

    const dataPairs: Record<string, string> = {};
    for (const pair of pending.kvPairs ?? []) {
      if (pair.key.trim()) dataPairs[pair.key.trim()] = pair.value;
    }

    try {
      await clientAction("/api/admin/announcements", {
        title: pending.title,
        body: pending.body,
        segment: pending.segment,
        data: Object.keys(dataPairs).length ? dataPairs : undefined,
      });

      const segmentLabel =
        SEGMENTS.find((s) => s.value === pending.segment)?.label ?? pending.segment;

      toast.success(`Announcement sent to ${segmentLabel}.`);
      reset();
      setConfirmOpen(false);
      queryClient.invalidateQueries({ queryKey: ["admin-announcements"] });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to send.");
    } finally {
      setSending(false);
    }
  }

  const isAll = pending?.segment === "all";
  const canConfirm = !isAll || allConfirmed;

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      {/* Composer Card */}
      <Card className="lg:col-span-5">
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <Megaphone className="h-5 w-5 text-primary" />
            New Announcement
          </CardTitle>
          <CardDescription>
            Compose and broadcast a push notification to users.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="title">Title</Label>
              <Input
                id="title"
                placeholder="Announcement title"
                {...register("title")}
              />
              {errors.title && (
                <p className="text-xs text-destructive">{errors.title.message}</p>
              )}
            </div>

            <div className="flex flex-col gap-1.5">
              <Label htmlFor="body">Message</Label>
              <textarea
                id="body"
                className="min-h-[110px] w-full rounded-md border bg-background px-3 py-2 text-sm resize-y focus:outline-none focus:ring-2 focus:ring-ring"
                placeholder="Write your announcement…"
                {...register("body")}
              />
              {errors.body && (
                <p className="text-xs text-destructive">{errors.body.message}</p>
              )}
            </div>

            <div className="flex flex-col gap-1.5">
              <Label>Audience</Label>
              <Select
                value={segment}
                onValueChange={(v) =>
                  setValue("segment", v as FormValues["segment"], {
                    shouldValidate: true,
                  })
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {SEGMENTS.map((s) => (
                    <SelectItem key={s.value} value={s.value}>
                      {s.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Optional key/value data pairs */}
            <div className="flex flex-col gap-2 pt-1">
              <div className="flex items-center justify-between">
                <Label className="text-muted-foreground text-xs">
                  Extra data (optional key/value pairs)
                </Label>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="h-7 text-xs"
                  onClick={() => append({ key: "", value: "" })}
                >
                  <Plus className="h-3 w-3 mr-1" />
                  Add field
                </Button>
              </div>
              {fields.map((field, i) => (
                <div key={field.id} className="flex gap-2 items-center">
                  <Input
                    placeholder="key"
                    className="h-8 text-xs"
                    {...register(`kvPairs.${i}.key`)}
                  />
                  <Input
                    placeholder="value"
                    className="h-8 text-xs"
                    {...register(`kvPairs.${i}.value`)}
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 shrink-0 text-muted-foreground hover:text-destructive"
                    onClick={() => remove(i)}
                    aria-label="Remove field"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
              ))}
            </div>

            <Button type="submit" className="mt-2 w-full sm:w-auto self-start">
              <Send className="h-4 w-4 mr-2" />
              Preview &amp; send
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* Announcements List */}
      <Card className="lg:col-span-7">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
          <div>
            <CardTitle className="text-lg">Sent Announcements</CardTitle>
            <CardDescription>
              View and remove active announcements from student notification feeds.
            </CardDescription>
          </div>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => refetch()}
            disabled={isFetching}
            title="Refresh announcements"
          >
            <RefreshCw className={`h-4 w-4 ${isFetching ? "animate-spin" : ""}`} />
          </Button>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              <Skeleton className="h-20 w-full rounded-md" />
              <Skeleton className="h-20 w-full rounded-md" />
              <Skeleton className="h-20 w-full rounded-md" />
            </div>
          ) : isError ? (
            <div className="rounded-lg border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive flex items-center gap-2">
              <AlertTriangle className="h-4 w-4 shrink-0" />
              <span>Failed to load announcements: {(error as Error).message}</span>
            </div>
          ) : announcements.length === 0 ? (
            <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">
              <Megaphone className="h-8 w-8 mx-auto mb-2 opacity-40" />
              <p className="font-medium">No announcements yet</p>
              <p className="text-xs text-muted-foreground mt-1">
                Broadcasted announcements will appear here and can be deleted at any time.
              </p>
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              {announcements.map((item) => (
                <div
                  key={item.id}
                  className="flex items-start justify-between gap-4 rounded-lg border p-4 hover:bg-muted/30 transition-colors"
                >
                  <div className="flex flex-col gap-1.5 flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <h4 className="font-semibold text-sm leading-tight text-foreground truncate">
                        {item.title}
                      </h4>
                      {item.recipient_count !== undefined && item.recipient_count > 0 && (
                        <Badge variant="secondary" className="text-xs font-normal gap-1">
                          <Users className="h-3 w-3" />
                          {item.recipient_count} {item.recipient_count === 1 ? "recipient" : "recipients"}
                        </Badge>
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground whitespace-pre-wrap line-clamp-3">
                      {item.body}
                    </p>
                    <span className="text-[11px] text-muted-foreground/80 mt-1 tabular-nums">
                      {formatDateTime(item.created_at)}
                    </span>
                  </div>

                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-muted-foreground hover:text-destructive hover:bg-destructive/10 shrink-0"
                    onClick={() => setDeleteItem(item)}
                    title="Delete announcement"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Broadcast confirmation dialog */}
      <Dialog open={confirmOpen} onOpenChange={(o) => !o && setConfirmOpen(false)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm broadcast</DialogTitle>
            <DialogDescription>
              You are about to send{" "}
              <strong>&ldquo;{pending?.title}&rdquo;</strong> to{" "}
              <strong>
                {SEGMENTS.find((s) => s.value === pending?.segment)?.label}
              </strong>
              .
            </DialogDescription>
          </DialogHeader>

          {isAll && (
            <div className="rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive flex flex-col gap-2">
              <p className="font-medium">⚠ This sends to ALL users on the platform.</p>
              <label className="flex items-center gap-2 cursor-pointer select-none">
                <input
                  type="checkbox"
                  checked={allConfirmed}
                  onChange={(e) => setAllConfirmed(e.target.checked)}
                  className="accent-destructive"
                />
                I understand — send to every user
              </label>
            </div>
          )}

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setConfirmOpen(false)}
              disabled={sending}
            >
              Cancel
            </Button>
            <Button
              variant={isAll ? "destructive" : "default"}
              disabled={!canConfirm || sending}
              onClick={confirmSend}
            >
              {sending ? "Sending…" : "Send announcement"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete confirmation dialog */}
      <Dialog open={!!deleteItem} onOpenChange={(o) => !o && setDeleteItem(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Announcement</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete{" "}
              <strong>&ldquo;{deleteItem?.title}&rdquo;</strong>?
            </DialogDescription>
          </DialogHeader>

          <div className="rounded-md border border-amber-500/30 bg-amber-500/10 p-3 text-xs text-amber-800 dark:text-amber-300 flex items-start gap-2">
            <AlertTriangle className="h-4 w-4 shrink-0 mt-0.5" />
            <span>
              This will permanently purge this announcement and remove it from the in-app notification feed of all{" "}
              {deleteItem?.recipient_count ? `${deleteItem.recipient_count} ` : ""}recipient users.
            </span>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteItem(null)}
              disabled={deleteMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              disabled={deleteMutation.isPending}
              onClick={() => {
                if (deleteItem) {
                  deleteMutation.mutate(deleteItem.id);
                }
              }}
            >
              {deleteMutation.isPending ? "Deleting…" : "Delete Announcement"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
