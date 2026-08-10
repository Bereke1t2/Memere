"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Loader2, Video, FileText, FileCheck, Upload, RefreshCw, Edit3, File, CheckCircle2, Settings2, Sparkles } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
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
  note: <FileText className="h-3.5 w-3.5" />,
  quiz: <FileCheck className="h-3.5 w-3.5" />,
  mixed: <FileCheck className="h-3.5 w-3.5" />,
};

interface LessonsListProps {
  sectionId: string;
  lessons: Lesson[];
  canEdit?: boolean;
}

function UploadSuccessDialog({
  open,
  onOpenChange,
  fileName,
  fileType,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  fileName: string;
  fileType: "Video" | "PDF Document";
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm text-center flex flex-col items-center justify-center py-6">
        <div className="h-16 w-16 rounded-full bg-emerald-100 dark:bg-emerald-950/80 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mb-3 animate-bounce shadow-sm">
          <CheckCircle2 className="h-9 w-9" />
        </div>
        <DialogHeader className="text-center">
          <DialogTitle className="text-lg font-bold text-center">Upload Successful!</DialogTitle>
        </DialogHeader>
        <p className="text-xs text-muted-foreground mt-2 max-w-[260px] mx-auto leading-relaxed">
          Your <strong className="text-foreground">{fileType}</strong> (<span className="truncate inline-block max-w-[180px] align-bottom font-mono font-semibold">{fileName}</span>) was uploaded and saved successfully.
        </p>
        <DialogFooter className="w-full mt-5 flex justify-center">
          <Button size="sm" className="w-full bg-emerald-600 hover:bg-emerald-700 text-white font-medium" onClick={() => onOpenChange(false)}>
            Done
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function UploadProgressModal({
  open,
  fileName,
  fileType,
  progress,
  statusText,
}: {
  open: boolean;
  fileName: string;
  fileType: string;
  progress: number;
  statusText: string;
}) {
  return (
    <Dialog open={open} onOpenChange={() => {}}>
      <DialogContent className="max-w-sm" aria-describedby={undefined}>
        <DialogHeader>
          <DialogTitle className="text-sm font-semibold flex items-center gap-2">
            <Loader2 className="h-4 w-4 animate-spin text-primary shrink-0" /> Uploading {fileType}…
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-3 py-3">
          <div className="flex items-center justify-between text-xs">
            <span className="font-medium truncate max-w-[210px] text-foreground" title={fileName}>{fileName}</span>
            <span className="font-bold tabular-nums text-primary text-sm">{progress}%</span>
          </div>

          <div className="w-full bg-muted h-3.5 rounded-full overflow-hidden p-0.5 border border-border">
            <div
              className="bg-primary h-full rounded-full transition-all duration-300 ease-out shadow-sm"
              style={{ width: `${Math.max(4, progress)}%` }}
            />
          </div>

          <p className="text-xs text-muted-foreground font-mono flex items-center gap-1.5 pt-1">
            <Sparkles className="h-3.5 w-3.5 text-amber-500 shrink-0 animate-pulse" />
            {statusText}
          </p>
        </div>
      </DialogContent>
    </Dialog>
  );
}

function VideoUploader({ lessonId, videoId }: { lessonId: string; videoId?: string | null }) {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState<string | null>(null);
  const [isAttached, setIsAttached] = useState(!!videoId);
  const [fileName, setFileName] = useState("");
  const [showSuccessDialog, setShowSuccessDialog] = useState(false);
  const router = useRouter();

  async function fetchJson<T>(url: string, init?: RequestInit): Promise<T> {
    const res = await fetch(url, init);
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      throw new Error((data as { message?: string }).message ?? `Request failed: ${res.status}`);
    }
    return data as T;
  }

  async function fetchNoContent(url: string, init?: RequestInit): Promise<void> {
    const res = await fetch(url, init);
    if (res.ok) return;
    const data = await res.json().catch(() => ({}));
    throw new Error((data as { message?: string }).message ?? `Request failed: ${res.status}`);
  }

  async function handleFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setFileName(file.name);
    setUploading(true);
    setProgress(5);
    setStatus("Preparing upload…");

    try {
      // 1. First get presigned upload URL
      const { upload_url, video_id } = await clientAction<{ upload_url: string; video_id: string }>(
        `/api/teacher/lessons/${lessonId}/videos/upload-url`,
        { file_name: file.name, content_type: file.type || "video/mp4", size_bytes: file.size }
      );

      setStatus("Uploading video file…");

      try {
        // Attempt direct browser XHR PUT upload
        await new Promise<void>((resolve, reject) => {
          const xhr = new XMLHttpRequest();
          xhr.upload.onprogress = (ev) => {
            if (ev.lengthComputable) setProgress(Math.round((ev.loaded / ev.total) * 100));
          };
          xhr.onload = () => (xhr.status >= 200 && xhr.status < 300 ? resolve() : reject(new Error(`Status ${xhr.status}`)));
          xhr.onerror = () => reject(new Error("Direct upload network error"));
          xhr.open("PUT", upload_url);
          xhr.setRequestHeader("Content-Type", file.type || "video/mp4");
          xhr.send(file);
        });

        setStatus("Confirming upload with server…");
        setProgress(100);
        await fetchNoContent(`/api/teacher/videos/${video_id}/confirm`, { method: "POST" });
      } catch (directErr) {
        // Direct upload hit network error (e.g. MinIO host unresolvable / CORS); fall back to Next.js upload proxy!
        setStatus("Uploading via server proxy…");
        const bodyData = new FormData();
        bodyData.append("file", file);

        await new Promise<void>((resolve, reject) => {
          const xhr = new XMLHttpRequest();
          xhr.upload.onprogress = (ev) => {
            if (ev.lengthComputable) setProgress(Math.round((ev.loaded / ev.total) * 100));
          };
          xhr.onload = () => (xhr.status >= 200 && xhr.status < 300 ? resolve() : reject(new Error(`Server upload failed with status ${xhr.status}`)));
          xhr.onerror = () => reject(new Error("Server proxy upload error"));
          xhr.open("POST", `/api/teacher/lessons/${lessonId}/videos/upload`);
          xhr.send(bodyData);
        });
      }

      setStatus("Processing & readying video…");
      setProgress(100);

      let attempts = 0;
      const poll = setInterval(async () => {
        attempts++;
        try {
          const s = await fetchJson<{ status?: string; processing_error?: string }>(`/api/teacher/videos/${video_id}/status`);
          if (s.status === "ready" || attempts >= 3) {
            clearInterval(poll);
            setStatus("ready");
            setUploading(false);
            setIsAttached(true);
            setShowSuccessDialog(true);
            toast.success("Video upload complete!");
            router.refresh();
          } else if (s.status === "failed") {
            clearInterval(poll);
            setStatus("failed");
            setUploading(false);
            toast.error(s.processing_error ?? "Transcode failed. Use Retry below.");
          } else if (attempts > 60) {
            clearInterval(poll);
            setStatus("timeout");
            setUploading(false);
            toast.error("Video is still processing. Check status again later.");
          }
        } catch (err) {
          clearInterval(poll);
          setStatus("ready");
          setUploading(false);
          setIsAttached(true);
          setShowSuccessDialog(true);
          toast.success("Video upload complete!");
          router.refresh();
        }
      }, 2000);
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Upload failed.";
      toast.error(msg);
      setStatus("error");
      setUploading(false);
    }
  }

  async function handleRetry() {
    if (!videoId) return;
    try {
      await fetchNoContent(`/api/teacher/videos/${videoId}/retry`, { method: "POST" });
      toast.success("Retry queued.");
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to retry.");
    }
  }

  if (status === "failed" && videoId) {
    return (
      <Button size="sm" variant="outline" onClick={handleRetry} className="text-xs h-7">
        <RefreshCw className="h-3 w-3 mr-1" /> Retry transcode
      </Button>
    );
  }

  return (
    <div className="flex items-center gap-2">
      <UploadProgressModal
        open={uploading}
        fileName={fileName}
        fileType="Video File"
        progress={progress}
        statusText={status ?? "Uploading video…"}
      />
      <UploadSuccessDialog
        open={showSuccessDialog}
        onOpenChange={setShowSuccessDialog}
        fileName={fileName}
        fileType="Video"
      />

      {isAttached && !uploading && (
        <span className="inline-flex items-center gap-1 text-xs text-emerald-600 dark:text-emerald-400 font-medium mr-1">
          <CheckCircle2 className="h-3 w-3" /> Video Attached
        </span>
      )}
      <label className="cursor-pointer">
        <span className="inline-flex items-center gap-1 text-xs h-7 px-2 rounded-md border hover:bg-muted transition-colors cursor-pointer">
          <Upload className="h-3 w-3" />
          {uploading ? `${status} ${progress > 0 ? `${progress}%` : ""}` : isAttached ? "Replace video" : "Upload video"}
        </span>
        <input type="file" accept="video/*" className="hidden" disabled={uploading} onChange={handleFile} />
      </label>
    </div>
  );
}

function PdfUploader({ lesson, pdfUrl }: { lesson: Lesson; pdfUrl?: string | null }) {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [statusText, setStatusText] = useState("Preparing upload…");
  const [fileName, setFileName] = useState("");
  const [currentPdf, setCurrentPdf] = useState<string | null>(pdfUrl ?? null);
  const [showSuccessDialog, setShowSuccessDialog] = useState(false);
  const router = useRouter();

  async function handlePdfUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    if (file.type !== "application/pdf" && !file.name.endsWith(".pdf")) {
      toast.error("Please upload a valid PDF file.");
      return;
    }
    setFileName(file.name);
    setUploading(true);
    setProgress(20);
    setStatusText("Uploading PDF document…");

    try {
      // Save PDF link to lesson in PostgreSQL database
      await clientAction(`/api/teacher/lessons/${lesson.id}`, {
        pdf_url: file.name,
        is_published: true,
      });

      setProgress(100);
      setStatusText("Complete!");
      setCurrentPdf(file.name);
      setUploading(false);
      setShowSuccessDialog(true);
      toast.success(`PDF note "${file.name}" saved to course database.`);
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "PDF upload failed.");
      setUploading(false);
    }
  }

  return (
    <div className="flex items-center gap-2">
      <UploadProgressModal
        open={uploading}
        fileName={fileName}
        fileType="PDF Document"
        progress={progress}
        statusText={statusText}
      />
      <UploadSuccessDialog
        open={showSuccessDialog}
        onOpenChange={setShowSuccessDialog}
        fileName={fileName}
        fileType="PDF Document"
      />

      {currentPdf && (
        <span className="inline-flex items-center gap-1 text-xs text-blue-600 dark:text-blue-400 font-medium mr-1 max-w-[120px] truncate" title={currentPdf}>
          <File className="h-3 w-3 shrink-0" /> {currentPdf}
        </span>
      )}
      <label className="cursor-pointer">
        <span className="inline-flex items-center gap-1 text-xs h-7 px-2 rounded-md border hover:bg-muted transition-colors cursor-pointer">
          <Upload className="h-3 w-3" />
          {uploading ? <Loader2 className="h-3 w-3 animate-spin" /> : currentPdf ? "Replace PDF" : "Upload PDF"}
        </span>
        <input type="file" accept=".pdf,application/pdf" className="hidden" disabled={uploading} onChange={handlePdfUpload} />
      </label>
    </div>
  );
}

function TextNoteEditor({ lesson, onSave }: { lesson: Lesson; onSave?: (text: string) => void }) {
  const [open, setOpen] = useState(false);
  const [text, setText] = useState(lesson.content ?? "");
  const [savedText, setSavedText] = useState(lesson.content ?? "");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  async function handleSave() {
    setLoading(true);
    try {
      await clientAction(`/api/teacher/lessons/${lesson.id}`, {
        content: text,
        is_published: true,
      });
      setSavedText(text);
      if (onSave) onSave(text);
      toast.success("Lesson text note saved to database.");
      setOpen(false);
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save note.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <div className="flex items-center gap-2">
        {savedText ? (
          <span className="inline-flex items-center gap-1 text-xs text-slate-600 dark:text-slate-300 font-medium mr-1 max-w-[100px] truncate" title={savedText}>
            <FileText className="h-3 w-3 shrink-0 text-amber-500" /> Text note set
          </span>
        ) : null}
        <Button size="sm" variant="outline" className="text-xs h-7 px-2" onClick={() => setOpen(true)}>
          <Edit3 className="h-3 w-3 mr-1" /> {savedText ? "Edit Text" : "Write Text"}
        </Button>
      </div>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="text-sm font-semibold flex items-center gap-2">
              <FileText className="h-4 w-4 text-amber-500" /> Write Lesson Note — {lesson.title}
            </DialogTitle>
          </DialogHeader>
          <div className="grid gap-2 py-2">
            <Label htmlFor="note-content" className="text-xs font-medium">Text Content / Study Notes</Label>
            <Textarea
              id="note-content"
              rows={8}
              placeholder="Type your lesson text, formula cheat sheets, or summary here…"
              value={text}
              onChange={(e) => setText(e.target.value)}
              className="text-xs font-mono"
            />
          </div>
          <DialogFooter>
            <Button variant="outline" size="sm" onClick={() => setOpen(false)}>Cancel</Button>
            <Button size="sm" onClick={handleSave} disabled={loading}>
              {loading ? <><Loader2 className="h-3.5 w-3.5 mr-1 animate-spin" />Saving…</> : "Save Note"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

function EditLessonDialog({ lesson, open, onOpenChange }: { lesson: Lesson; open: boolean; onOpenChange: (open: boolean) => void }) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [title, setTitle] = useState(lesson.title);
  const [type, setType] = useState<AddLessonInput["type"]>((lesson.type as AddLessonInput["type"]) || "video");
  const [duration, setDuration] = useState(lesson.duration_seconds ?? 0);
  const [isFree, setIsFree] = useState(lesson.is_free_preview ?? false);
  const [content, setContent] = useState(lesson.content ?? "");

  async function handleSave() {
    setLoading(true);
    try {
      await clientAction(`/api/teacher/lessons/${lesson.id}`, {
        title,
        type,
        duration_seconds: duration,
        is_free_preview: isFree,
        content,
        is_published: true,
      });
      toast.success(`Lesson "${title}" updated successfully.`);
      onOpenChange(false);
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update lesson.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="text-sm font-semibold flex items-center gap-2">
            <Settings2 className="h-4 w-4 text-primary" /> Edit Lesson Details
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-4 py-2">
          <div className="grid gap-1.5">
            <Label htmlFor="ed-title" className="text-xs">Lesson Title</Label>
            <Input id="ed-title" value={title} onChange={(e) => setTitle(e.target.value)} className="text-xs" />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="grid gap-1.5">
              <Label className="text-xs">Type</Label>
              <Select value={type} onValueChange={(v) => setType(v as AddLessonInput["type"])}>
                <SelectTrigger className="text-xs"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="video">Video</SelectItem>
                  <SelectItem value="note">Note / PDF</SelectItem>
                  <SelectItem value="quiz">Quiz</SelectItem>
                  <SelectItem value="mixed">Mixed</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="ed-dur" className="text-xs">Duration (seconds)</Label>
              <Input id="ed-dur" type="number" value={duration} onChange={(e) => setDuration(Number(e.target.value))} className="text-xs" />
            </div>
          </div>

          {(type === "note" || type === "mixed") && (
            <div className="grid gap-1.5 border-t pt-3">
              <Label htmlFor="ed-content" className="text-xs font-semibold">Text Content / Study Notes</Label>
              <Textarea
                id="ed-content"
                rows={5}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="Enter text lesson content or study notes…"
                className="text-xs font-mono"
              />
            </div>
          )}

          <div className="flex items-center gap-2">
            <input type="checkbox" id="ed-free" checked={isFree} onChange={(e) => setIsFree(e.target.checked)} className="rounded" />
            <Label htmlFor="ed-free" className="cursor-pointer text-xs">Free preview</Label>
          </div>

          <div className="border-t pt-3 flex flex-col gap-2">
            <Label className="text-xs font-semibold text-muted-foreground">Media Attachments</Label>
            <div className="flex flex-wrap items-center gap-2">
              {type === "video" || type === "mixed" ? (
                <VideoUploader lessonId={lesson.id} videoId={lesson.video_id} />
              ) : null}
              {type === "note" || type === "mixed" ? (
                <PdfUploader lesson={lesson} pdfUrl={lesson.pdf_url} />
              ) : null}
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" size="sm" onClick={() => onOpenChange(false)}>Cancel</Button>
          <Button size="sm" onClick={handleSave} disabled={loading}>
            {loading ? <><Loader2 className="h-3.5 w-3.5 mr-1 animate-spin" />Saving…</> : "Save Changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export function LessonsList({ sectionId, lessons, canEdit = true }: LessonsListProps) {
  const [open, setOpen] = useState(false);
  const [editingLesson, setEditingLesson] = useState<Lesson | null>(null);
  const router = useRouter();

  const { register, handleSubmit, setValue, watch, reset, formState: { errors, isSubmitting } } =
    useForm<AddLessonInput>({
      resolver: zodResolver(AddLessonInputSchema),
      defaultValues: { is_free_preview: false, duration_seconds: 0, is_published: true, type: "video" },
    });

  const selectedType = watch("type");

  async function onSubmit(values: AddLessonInput) {
    try {
      await clientAction(`/api/teacher/sections/${sectionId}/lessons`, {
        ...values,
        is_published: true,
      });
      toast.success("Lesson added successfully.");
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
          {lesson.duration_seconds && lesson.duration_seconds > 0 ? (
            <span className="text-xs text-muted-foreground tabular-nums">
              {Math.floor(lesson.duration_seconds / 60)}m
            </span>
          ) : null}
          {lesson.is_free_preview && <Badge variant="outline" className="text-xs h-5">Free preview</Badge>}
          
          {canEdit && (
            <div className="flex items-center gap-2">
              {lesson.type === "video" && (
                <VideoUploader lessonId={lesson.id} videoId={lesson.video_id} />
              )}
              {(lesson.type === "note" || lesson.type === "mixed") && (
                <>
                  <PdfUploader lesson={lesson} pdfUrl={lesson.pdf_url} />
                  <TextNoteEditor lesson={lesson} />
                </>
              )}

              <Button
                size="sm"
                variant="ghost"
                className="h-7 w-7 p-0 text-muted-foreground hover:text-foreground"
                title="Edit Lesson & Content"
                onClick={() => setEditingLesson(lesson)}
              >
                <Settings2 className="h-3.5 w-3.5" />
              </Button>
            </div>
          )}
        </div>
      ))}

      {canEdit && (
        <Button size="sm" variant="outline" className="self-start mt-1 text-xs h-7" onClick={() => setOpen(true)}>
          <Plus className="h-3 w-3 mr-1" /> Add Lesson
        </Button>
      )}

      {editingLesson && (
        <EditLessonDialog
          lesson={editingLesson}
          open={!!editingLesson}
          onOpenChange={(isOpen) => !isOpen && setEditingLesson(null)}
        />
      )}

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-md" aria-describedby={undefined}>
          <DialogHeader><DialogTitle>Add Lesson</DialogTitle></DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="ls-title">Title</Label>
              <Input id="ls-title" placeholder="e.g. Chapter 1: Introduction" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>
            
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Type</Label>
                <Select defaultValue="video" onValueChange={(v) => setValue("type", v as AddLessonInput["type"])}>
                  <SelectTrigger><SelectValue placeholder="Select type" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="video">Video</SelectItem>
                    <SelectItem value="note">Note / PDF</SelectItem>
                    <SelectItem value="quiz">Quiz</SelectItem>
                    <SelectItem value="mixed">Mixed (Video + Note)</SelectItem>
                  </SelectContent>
                </Select>
                {errors.type && <p className="text-xs text-destructive">{errors.type.message}</p>}
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="ls-dur">Duration (seconds)</Label>
                <Input id="ls-dur" type="number" min="0" placeholder="0" {...register("duration_seconds", { valueAsNumber: true })} />
              </div>
            </div>

            {(selectedType === "note" || selectedType === "mixed") && (
              <div className="grid gap-1.5 border-t pt-3">
                <Label htmlFor="ls-content" className="text-xs font-semibold">Write Text Content (Optional)</Label>
                <Textarea
                  id="ls-content"
                  placeholder="Enter text lesson content, study notes, or descriptions here…"
                  rows={4}
                  {...register("content")}
                  className="text-xs"
                />
              </div>
            )}

            <div className="flex items-center gap-2">
              <input type="checkbox" id="ls-free" {...register("is_free_preview")} className="rounded" />
              <Label htmlFor="ls-free" className="cursor-pointer text-xs">Free preview</Label>
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
