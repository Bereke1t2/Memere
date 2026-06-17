"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Loader2, Globe } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription,
} from "@/components/ui/dialog";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select";
import { clientAction } from "@/lib/client-action";
import { CreateExamInputSchema, type CreateExamInput, type Exam } from "@/lib/api/schemas";

interface ExamsPanelProps {
  courseId: string;
  exams: Exam[];
}

export function ExamsPanel({ courseId, exams }: ExamsPanelProps) {
  const [createOpen, setCreateOpen] = useState(false);
  const [publishId, setPublishId] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const router = useRouter();

  const { register, handleSubmit, setValue, reset, formState: { errors, isSubmitting } } =
    useForm<CreateExamInput>({
      resolver: zodResolver(CreateExamInputSchema),
      defaultValues: { duration_minutes: 60, pass_marks: 50, grade: 12 },
    });

  async function onCreate(values: CreateExamInput) {
    try {
      await clientAction(`/api/teacher/courses/${courseId}/exams`, values);
      toast.success("Exam created.");
      setCreateOpen(false);
      reset();
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create exam.");
    }
  }

  async function handlePublish() {
    if (!publishId) return;
    setBusy(true);
    try {
      await fetch(`/api/teacher/exams/${publishId}/publish`, { method: "POST" });
      toast.success("Exam published.");
      setPublishId(null);
      router.refresh();
    } catch {
      toast.error("Failed to publish exam.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-end">
        <Button size="sm" onClick={() => setCreateOpen(true)}>
          <Plus className="h-4 w-4 mr-1.5" /> New Exam
        </Button>
      </div>

      {exams.length === 0 ? (
        <div className="py-10 text-center text-sm text-muted-foreground">No exams yet.</div>
      ) : (
        exams.map((exam) => (
          <Card key={exam.id} className="border shadow-none">
            <CardContent className="p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <div className="flex items-center gap-2">
                    <p className="font-medium text-sm">{exam.title}</p>
                    <Badge variant={exam.is_published ? "default" : "secondary"} className="text-xs">
                      {exam.is_published ? "Published" : "Draft"}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-3 mt-1 text-xs text-muted-foreground">
                    {exam.subject && <span>{exam.subject}</span>}
                    {exam.grade && <span>Grade {exam.grade}</span>}
                    {exam.duration_minutes && <span>{exam.duration_minutes}min</span>}
                    {exam.pass_marks != null && <span>Pass: {exam.pass_marks} marks</span>}
                  </div>
                </div>
                {!exam.is_published && (
                  <Button size="sm" variant="outline" className="text-xs h-7 shrink-0"
                    onClick={() => setPublishId(exam.id)}>
                    <Globe className="h-3 w-3 mr-1" /> Publish
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        ))
      )}

      {/* Create dialog */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-sm">
          <DialogHeader><DialogTitle>Create Exam</DialogTitle></DialogHeader>
          <form onSubmit={handleSubmit(onCreate)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="ex-title">Title</Label>
              <Input id="ex-title" placeholder="e.g. EUEE 2024 Mock Exam" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="ex-sub">Subject</Label>
                <Input id="ex-sub" placeholder="Mathematics" {...register("subject")} />
                {errors.subject && <p className="text-xs text-destructive">{errors.subject.message}</p>}
              </div>
              <div className="grid gap-1.5">
                <Label>Grade</Label>
                <Select defaultValue="12" onValueChange={(v) => setValue("grade", Number(v))}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {Array.from({ length: 12 }, (_, i) => i + 1).map((g) => (
                      <SelectItem key={g} value={String(g)}>Grade {g}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="ex-dur">Duration (minutes)</Label>
                <Input id="ex-dur" type="number" min="1" {...register("duration_minutes", { valueAsNumber: true })} />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="ex-pass">Pass marks</Label>
                <Input id="ex-pass" type="number" min="0" {...register("pass_marks", { valueAsNumber: true })} />
              </div>
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="ex-inst">Instructions (optional)</Label>
              <textarea
                id="ex-inst"
                rows={2}
                className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
                {...register("instructions")}
              />
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setCreateOpen(false); reset(); }}>Cancel</Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Creating…</> : "Create Exam"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Publish confirm */}
      <Dialog open={!!publishId} onOpenChange={(v) => !v && setPublishId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Publish exam?</DialogTitle>
            <DialogDescription>
              This exam will become visible to students enrolled in the course.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPublishId(null)}>Cancel</Button>
            <Button onClick={handlePublish} disabled={busy}>
              {busy ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Publishing…</> : "Publish"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
