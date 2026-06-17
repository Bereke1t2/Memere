"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Loader2, Video, FileText, FileCheck, Upload, RefreshCw } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter,
} from "@/components/ui/dialog";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select";
import { clientAction } from "@/lib/client-action";
import { AddLessonInputSchema, type AddLessonInput, type Lesson } from "@/lib/api/schemas";

const TYPE_ICONS: Record<string, React.ReactNode> = {
  video: <Video className="h-3.5 w-3.5" />,
  text: <FileText className="h-3.5 w-3.5" />,
  pdf: <FileCheck className="h-3.5 w-3.5" />,
};

interface LessonsListProps {
  sectionId: string;
  lessons: Lesson[];
}

function VideoUploader({ lessonId, videoId }: { lessonId: string; videoId?: string | null }) {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState<string | null>(null);
  const router = useRouter();

  async function handleFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setUploading(true);
    setProgress(0);
    setStatus("Getting upload URL…");
    try {
      const { upload_url, video_id } = await clientAction<{ upload_url: string; video_id: string }>(
        `/api/teacher/lessons/${lessonId}/videos/upload-url`,
        { file_name: file.name, content_type: file.type, size_bytes: file.size }
      );

      setStatus("Uploading…");
      await new Promise<void>((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.upload.onprogress = (ev) => {
          if (ev.lengthComputable) setProgress(Math.round((ev.loaded / ev.total) * 100));
        };
        xhr.onload = () => (xhr.status >= 200 && xhr.status < 300 ? resolve() : reject(new Error(`Upload failed: ${xhr.status}`)));
        xhr.onerror = () => reject(new Error("Network error during upload"));
        xhr.open("PUT", upload_url);
        xhr.setRequestHeader("Content-Type", file.type);
        xhr.send(file);
      });

      setStatus("Confirming…");
      await fetch(`/api/teacher/videos/${video_id}/confirm`, { method: "POST" });
      setStatus("Processing…");

      let attempts = 0;
      const poll = setInterval(async () => {
        attempts++;
        try {
          const s = await fetch(`/api/teacher/videos/${video_id}/status`).then((r) => r.json());
          if (s.status === "ready") {
            clearInterval(poll);
            setStatus("ready");
            setUploading(false);
            toast.success("Video ready.");
            router.refresh();
          } else if (s.status === "failed") {
            clearInterval(poll);
            setStatus("failed");
            setUploading(false);
            toast.error("Transcode failed. Use Retry below.");
          } else if (attempts > 60) {
            clearInterval(poll);
            setStatus("timeout");
            setUploading(false);
          }
        } catch { /* continue polling */ }
      }, 5000);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Upload failed.");
      setStatus("error");
      setUploading(false);
    }
  }

  async function handleRetry() {
    if (!videoId) return;
    try {
      await fetch(`/api/teacher/videos/${videoId}/retry`, { method: "POST" });
      toast.success("Retry queued.");
      router.refresh();
    } catch {
      toast.error("Failed to retry.");
    }
  }

  if (status === "failed" && videoId) {
    return (
      <Button size="sm" variant="outline" onClick={handleRetry} className="text-xs h-7">
        <RefreshCw className="h-3 w-3 mr-1" /> Retry transcode
      </Button>
    );
  }

  if (videoId && !uploading) {
    return <span className="text-xs text-muted-foreground">Video attached</span>;
  }

  return (
    <div className="flex items-center gap-2">
      <label className="cursor-pointer">
        <span className="inline-flex items-center gap-1 text-xs h-7 px-2 rounded-md border hover:bg-muted transition-colors cursor-pointer">
          <Upload className="h-3 w-3" />
          {uploading ? `${status} ${progress > 0 ? `${progress}%` : ""}` : "Upload video"}
        </span>
        <input type="file" accept="video/*" className="hidden" disabled={uploading} onChange={handleFile} />
      </label>
    </div>
  );
}

export function LessonsList({ sectionId, lessons }: LessonsListProps) {
  const [open, setOpen] = useState(false);
  const router = useRouter();

  const { register, handleSubmit, setValue, reset, formState: { errors, isSubmitting } } =
    useForm<AddLessonInput>({
      resolver: zodResolver(AddLessonInputSchema),
      defaultValues: { is_free_preview: false, duration_seconds: 0, is_published: false },
    });

  async function onSubmit(values: AddLessonInput) {
    try {
      await clientAction(`/api/teacher/sections/${sectionId}/lessons`, values);
      toast.success("Lesson added.");
      setOpen(false);
      reset();
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add lesson.");
    }
  }

  return (
    <div className="flex flex-col gap-2 mt-2">
      {lessons.map((lesson, idx) => (
        <div key={lesson.id} className="flex items-center gap-3 px-3 py-2 rounded-lg bg-muted/40 text-sm">
          <span className="text-muted-foreground w-5 text-center text-xs">{lesson.order_index ?? idx + 1}</span>
          <span className="flex items-center gap-1 text-muted-foreground">
            {TYPE_ICONS[lesson.type] ?? <FileText className="h-3.5 w-3.5" />}
          </span>
          <span className="flex-1 font-medium">{lesson.title}</span>
          {lesson.duration_seconds && lesson.duration_seconds > 0 && (
            <span className="text-xs text-muted-foreground tabular-nums">
              {Math.floor(lesson.duration_seconds / 60)}m
            </span>
          )}
          {lesson.is_free_preview && <Badge variant="outline" className="text-xs h-5">Free preview</Badge>}
          {lesson.type === "video" && (
            <VideoUploader lessonId={lesson.id} videoId={lesson.video_id} />
          )}
        </div>
      ))}

      <Button size="sm" variant="outline" className="self-start mt-1 text-xs h-7" onClick={() => setOpen(true)}>
        <Plus className="h-3 w-3 mr-1" /> Add Lesson
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-sm">
          <DialogHeader><DialogTitle>Add Lesson</DialogTitle></DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="ls-title">Title</Label>
              <Input id="ls-title" placeholder="e.g. Introduction" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Type</Label>
                <Select onValueChange={(v) => setValue("type", v as AddLessonInput["type"])}>
                  <SelectTrigger><SelectValue placeholder="Select type" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="video">Video</SelectItem>
                    <SelectItem value="text">Text</SelectItem>
                    <SelectItem value="pdf">PDF</SelectItem>
                  </SelectContent>
                </Select>
                {errors.type && <p className="text-xs text-destructive">{errors.type.message}</p>}
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="ls-dur">Duration (seconds)</Label>
                <Input id="ls-dur" type="number" min="0" placeholder="0" {...register("duration_seconds", { valueAsNumber: true })} />
              </div>
            </div>
            <div className="flex items-center gap-2">
              <input type="checkbox" id="ls-free" {...register("is_free_preview")} className="rounded" />
              <Label htmlFor="ls-free" className="cursor-pointer">Free preview</Label>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setOpen(false); reset(); }}>Cancel</Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Adding…</> : "Add Lesson"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
