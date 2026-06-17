"use client";

import { useState } from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { toast } from "sonner";
import { clientAction } from "@/lib/client-action";
import { Plus, Trash2 } from "lucide-react";
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
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

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

interface HistoryEntry {
  title: string;
  segment: string;
  sentAt: string;
}

export function AnnouncementsClient() {
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [allConfirmed, setAllConfirmed] = useState(false);
  const [sending, setSending] = useState(false);
  const [history, setHistory] = useState<HistoryEntry[]>([]);
  const [pending, setPending] = useState<FormValues | null>(null);

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

  function onSubmit(values: FormValues) {
    setPending(values);
    setAllConfirmed(false);
    setConfirmOpen(true);
  }

  async function confirmSend() {
    if (!pending) return;
    setSending(true);

    const data: Record<string, string> = {};
    for (const pair of pending.kvPairs ?? []) {
      if (pair.key.trim()) data[pair.key.trim()] = pair.value;
    }

    try {
      await clientAction("/api/admin/announcements", {
        title: pending.title,
        body: pending.body,
        segment: pending.segment,
        data: Object.keys(data).length ? data : undefined,
      });

      const segmentLabel =
        SEGMENTS.find((s) => s.value === pending.segment)?.label ?? pending.segment;

      toast.success(`Announcement sent to ${segmentLabel}.`);
      setHistory((h) => [
        { title: pending.title, segment: segmentLabel, sentAt: new Date().toISOString() },
        ...h,
      ]);
      reset();
      setConfirmOpen(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to send.");
    } finally {
      setSending(false);
    }
  }

  const isAll = pending?.segment === "all";
  const canConfirm = !isAll || allConfirmed;

  return (
    <div className="flex flex-col gap-8 max-w-xl">
      {/* Composer */}
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="title">Title</Label>
          <Input id="title" placeholder="Announcement title" {...register("title")} />
          {errors.title && (
            <p className="text-xs text-destructive">{errors.title.message}</p>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <Label htmlFor="body">Message</Label>
          <textarea
            id="body"
            className="min-h-[120px] w-full rounded-md border bg-background px-3 py-2 text-sm resize-y focus:outline-none focus:ring-2 focus:ring-ring"
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
              setValue("segment", v as FormValues["segment"], { shouldValidate: true })
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
        <div className="flex flex-col gap-2">
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
                className="h-8 w-8 shrink-0"
                onClick={() => remove(i)}
                aria-label="Remove field"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </Button>
            </div>
          ))}
        </div>

        <Button type="submit" className="self-start">
          Preview &amp; send
        </Button>
      </form>

      {/* Confirm dialog */}
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

      {/* Session history */}
      {history.length > 0 && (
        <>
          <Separator />
          <div className="flex flex-col gap-3">
            <div>
              <h2 className="text-sm font-semibold">Sent this session</h2>
              <p className="text-xs text-muted-foreground">
                Session-only — the backend stores no retrievable announcement list.
              </p>
            </div>
            <ol className="flex flex-col gap-2">
              {history.map((entry, i) => (
                <li
                  key={i}
                  className="flex items-center gap-3 rounded-md border px-4 py-2.5 text-sm"
                >
                  <span className="flex-1 font-medium truncate">{entry.title}</span>
                  <Badge variant="outline" className="capitalize shrink-0">
                    {entry.segment}
                  </Badge>
                  <span className="text-xs text-muted-foreground shrink-0">
                    {new Date(entry.sentAt).toLocaleTimeString()}
                  </span>
                </li>
              ))}
            </ol>
          </div>
        </>
      )}
    </div>
  );
}
