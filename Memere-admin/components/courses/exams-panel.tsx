"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Loader2, Globe, BarChart3, Pencil, Trash2, X } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { clientAction } from "@/lib/client-action";
import {
  CreateExamInputSchema,
  ExamQuestionInputSchema,
  type CreateExamInput,
  type ExamQuestionInput,
  type Exam,
  type ExamStats,
} from "@/lib/api/schemas";

interface ExamsPanelProps {
  courseId: string;
  exams: Exam[];
  canEdit?: boolean;
}

function AddExamQuestionDialog({
  exam,
  onDone,
}: {
  exam: Exam;
  onDone: () => void;
}) {
  const [open, setOpen] = useState(false);

  const {
    register,
    handleSubmit,
    control,
    setValue,
    watch,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ExamQuestionInput>({
    resolver: zodResolver(ExamQuestionInputSchema),
    defaultValues: {
      type: "multiple_choice",
      marks: 5,
      order_index: 0,
      subject: exam.subject ?? "",
      topic: "",
      explanation: "",
      answers: [
        { text: "", is_correct: false },
        { text: "", is_correct: false },
        { text: "", is_correct: false },
        { text: "", is_correct: false },
      ],
    },
  });

  const { fields, append, remove, replace } = useFieldArray({ control, name: "answers" });
  const selectedType = watch("type");
  const answers = watch("answers");

  function handleTypeChange(val: "multiple_choice" | "true_false") {
    setValue("type", val);
    if (val === "true_false") {
      replace([
        { text: "True", is_correct: true },
        { text: "False", is_correct: false },
      ]);
    } else {
      replace([
        { text: "", is_correct: false },
        { text: "", is_correct: false },
        { text: "", is_correct: false },
        { text: "", is_correct: false },
      ]);
    }
  }

  async function onSubmit(values: ExamQuestionInput) {
    const hasCorrect = values.answers?.some((a) => a.is_correct);
    if (!hasCorrect) {
      toast.error("Please mark at least one answer as correct.");
      return;
    }
    try {
      await clientAction(`/api/teacher/exams/${exam.id}/questions`, values as unknown as Record<string, unknown>);
      toast.success("Question added.");
      setOpen(false);
      reset({
        type: "multiple_choice",
        marks: 5,
        order_index: 0,
        subject: exam.subject ?? "",
        topic: "",
        explanation: "",
        answers: [
          { text: "", is_correct: false },
          { text: "", is_correct: false },
          { text: "", is_correct: false },
          { text: "", is_correct: false },
        ],
      });
      onDone();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add question.");
    }
  }

  return (
    <>
      <Button size="sm" variant="outline" className="text-xs h-7" onClick={() => setOpen(true)}>
        <Plus className="h-3 w-3 mr-1" /> Add Question
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-lg max-h-[85vh] overflow-y-auto" aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Add Exam Question</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="eq-text">Question Prompt</Label>
              <textarea
                id="eq-text"
                rows={3}
                placeholder="Type your question here…"
                className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
                {...register("text")}
              />
              {errors.text && <p className="text-xs text-destructive">{errors.text.message}</p>}
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Type</Label>
                <Select
                  value={selectedType}
                  onValueChange={(v) => handleTypeChange(v as "multiple_choice" | "true_false")}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="multiple_choice">Multiple choice</SelectItem>
                    <SelectItem value="true_false">True / False</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="eq-marks">Marks</Label>
                <Input
                  id="eq-marks"
                  type="number"
                  min="1"
                  {...register("marks", { valueAsNumber: true })}
                />
                {errors.marks && <p className="text-xs text-destructive">{errors.marks.message}</p>}
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="eq-sub">Subject (optional)</Label>
                <Input id="eq-sub" placeholder="Subject" {...register("subject")} />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="eq-top">Topic (optional)</Label>
                <Input id="eq-top" placeholder="Topic" {...register("topic")} />
              </div>
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="eq-exp">Explanation (optional)</Label>
              <Input id="eq-exp" placeholder="Shown after attempt" {...register("explanation")} />
            </div>

            <div className="grid gap-2 border-t pt-3">
              <div className="flex items-center justify-between">
                <Label className="text-xs font-semibold">Answers — tick the correct one(s)</Label>
                {selectedType === "multiple_choice" && fields.length < 6 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-6 text-xs px-2 text-primary"
                    onClick={() => append({ text: "", is_correct: false })}
                  >
                    <Plus className="h-3 w-3 mr-1" /> Add choice
                  </Button>
                )}
              </div>

              {fields.map((field, i) => (
                <div key={field.id} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={answers?.[i]?.is_correct ?? false}
                    onChange={(e) => {
                      if (selectedType === "true_false") {
                        fields.forEach((_, idx) => {
                          setValue(`answers.${idx}.is_correct`, idx === i ? e.target.checked : false);
                        });
                      } else {
                        setValue(`answers.${i}.is_correct`, e.target.checked);
                      }
                    }}
                    className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                    title="Mark as correct"
                  />
                  <Input
                    placeholder={`Answer ${i + 1}`}
                    disabled={selectedType === "true_false"}
                    {...register(`answers.${i}.text`)}
                  />
                  {selectedType === "multiple_choice" && fields.length > 2 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive shrink-0"
                      onClick={() => remove(i)}
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              ))}
              {errors.answers && (
                <p className="text-xs text-destructive">
                  {Array.isArray(errors.answers) ? "Fix answers" : errors.answers.message}
                </p>
              )}
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setOpen(false); reset(); }}>
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Adding…</> : "Add Question"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}

function EditExamDialog({
  exam,
  onDone,
}: {
  exam: Exam;
  onDone: () => void;
}) {
  const [open, setOpen] = useState(false);

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<CreateExamInput>({
    resolver: zodResolver(CreateExamInputSchema),
    defaultValues: {
      title: exam.title,
      subject: exam.subject ?? "",
      grade: exam.grade ?? 12,
      duration_minutes: exam.duration_minutes ?? 60,
      pass_marks: exam.pass_marks ?? 50,
      instructions: exam.instructions ?? "",
    },
  });

  async function onSubmit(values: CreateExamInput) {
    try {
      await clientAction(`/api/teacher/exams/${exam.id}`, values as unknown as Record<string, unknown>, "PUT");
      toast.success("Exam updated.");
      setOpen(false);
      onDone();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update exam.");
    }
  }

  return (
    <>
      <Button size="sm" variant="ghost" className="text-xs h-7" onClick={() => setOpen(true)}>
        <Pencil className="h-3 w-3 mr-1" /> Edit
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-md" aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Edit Exam</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="ed-ex-title">Title</Label>
              <Input id="ed-ex-title" placeholder="e.g. EUEE 2024 Mock Exam" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="ed-ex-sub">Subject</Label>
                <Input id="ed-ex-sub" placeholder="Mathematics" {...register("subject")} />
                {errors.subject && <p className="text-xs text-destructive">{errors.subject.message}</p>}
              </div>
              <div className="grid gap-1.5">
                <Label>Grade</Label>
                <Select
                  defaultValue={String(exam.grade ?? 12)}
                  onValueChange={(v) => setValue("grade", Number(v))}
                >
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
                <Label htmlFor="ed-ex-dur">Duration (minutes)</Label>
                <Input
                  id="ed-ex-dur"
                  type="number"
                  min="1"
                  {...register("duration_minutes", { valueAsNumber: true })}
                />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="ed-ex-pass">Pass marks</Label>
                <Input
                  id="ed-ex-pass"
                  type="number"
                  min="0"
                  {...register("pass_marks", { valueAsNumber: true })}
                />
              </div>
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="ed-ex-inst">Instructions (optional)</Label>
              <Textarea
                id="ed-ex-inst"
                rows={3}
                placeholder="Guidelines, rules, or instructions for students…"
                {...register("instructions")}
              />
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Saving…</> : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}

function DeleteExamDialog({
  exam,
  open,
  onOpenChange,
  onDone,
}: {
  exam: Exam;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDone: () => void;
}) {
  const [deleting, setDeleting] = useState(false);

  async function handleDelete() {
    setDeleting(true);
    try {
      await clientAction(`/api/teacher/exams/${exam.id}`, undefined, "DELETE");
      toast.success(`Exam "${exam.title}" deleted.`);
      onOpenChange(false);
      onDone();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete exam.");
    } finally {
      setDeleting(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={(o) => !deleting && onOpenChange(o)}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle className="text-sm font-semibold flex items-center gap-2 text-destructive">
            <Trash2 className="h-4 w-4" /> Delete Exam?
          </DialogTitle>
          <DialogDescription className="text-xs pt-1">
            This will permanently delete <strong>{exam.title}</strong> and all associated questions and attempts.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="mt-2">
          <Button variant="outline" size="sm" onClick={() => onOpenChange(false)} disabled={deleting}>
            Cancel
          </Button>
          <Button variant="destructive" size="sm" onClick={handleDelete} disabled={deleting}>
            {deleting ? <><Loader2 className="h-3.5 w-3.5 mr-1 animate-spin" />Deleting…</> : "Delete Exam"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function ExamStatsCard({ examId }: { examId: string }) {
  const [stats, setStats] = useState<ExamStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [loaded, setLoaded] = useState(false);

  async function loadStats() {
    setLoading(true);
    try {
      const res = await fetch(`/api/teacher/exams/${examId}/stats`);
      if (res.ok) {
        setStats(await res.json());
        setLoaded(true);
      } else {
        toast.error("Failed to load stats.");
      }
    } catch {
      toast.error("Failed to load stats.");
    } finally {
      setLoading(false);
    }
  }

  if (!loaded) {
    return (
      <Button
        size="sm"
        variant="ghost"
        className="text-xs h-7 text-muted-foreground"
        onClick={loadStats}
        disabled={loading}
      >
        {loading ? (
          <Loader2 className="h-3 w-3 animate-spin mr-1" />
        ) : (
          <BarChart3 className="h-3 w-3 mr-1" />
        )}
        Stats
      </Button>
    );
  }

  if (!stats) return null;

  return (
    <div className="mt-3 grid grid-cols-2 sm:grid-cols-3 gap-2 text-xs border rounded-md p-2 bg-muted/20">
      <div>
        <p className="text-muted-foreground">Attempts</p>
        <p className="font-medium">{stats.total_attempts}</p>
      </div>
      {stats.avg_score != null && (
        <div>
          <p className="text-muted-foreground">Avg score</p>
          <p className="font-medium">{stats.avg_score.toFixed(1)}</p>
        </div>
      )}
      {stats.pass_rate != null && (
        <div>
          <p className="text-muted-foreground">Pass rate</p>
          <p className="font-medium">{(stats.pass_rate * 100).toFixed(0)}%</p>
        </div>
      )}
      {stats.highest_score != null && (
        <div>
          <p className="text-muted-foreground">Highest</p>
          <p className="font-medium">{stats.highest_score}</p>
        </div>
      )}
      {stats.lowest_score != null && (
        <div>
          <p className="text-muted-foreground">Lowest</p>
          <p className="font-medium">{stats.lowest_score}</p>
        </div>
      )}
    </div>
  );
}

export function ExamsPanel({ courseId, exams, canEdit = true }: ExamsPanelProps) {
  const [createOpen, setCreateOpen] = useState(false);
  const [publishId, setPublishId] = useState<string | null>(null);
  const [deleteExamItem, setDeleteExamItem] = useState<Exam | null>(null);
  const [busy, setBusy] = useState(false);
  const router = useRouter();

  const {
    register,
    handleSubmit,
    setValue,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CreateExamInput>({
    resolver: zodResolver(CreateExamInputSchema),
    defaultValues: { duration_minutes: 60, pass_marks: 50, grade: 12 },
  });

  async function onCreate(values: CreateExamInput) {
    try {
      await clientAction(`/api/teacher/courses/${courseId}/exams`, values as unknown as Record<string, unknown>);
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
      {canEdit && (
        <div className="flex justify-end">
          <Button size="sm" onClick={() => setCreateOpen(true)}>
            <Plus className="h-4 w-4 mr-1.5" /> New Exam
          </Button>
        </div>
      )}

      {exams.length === 0 ? (
        <div className="py-10 text-center text-sm text-muted-foreground">
          No exams yet. Create an exam to evaluate your students.
        </div>
      ) : (
        exams.map((exam) => (
          <Card key={exam.id} className="border shadow-none">
            <CardContent className="p-4">
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <p className="font-medium text-sm">{exam.title}</p>
                    <Badge variant={exam.is_published ? "default" : "secondary"} className="text-xs">
                      {exam.is_published ? "Published" : "Draft"}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-3 mt-1 text-xs text-muted-foreground flex-wrap">
                    {exam.subject && <span>{exam.subject}</span>}
                    {exam.grade && <span>Grade {exam.grade}</span>}
                    {exam.duration_minutes && <span>{exam.duration_minutes}min</span>}
                    {exam.pass_marks != null && <span>Pass: {exam.pass_marks} marks</span>}
                    {exam.total_marks != null && <span>Total: {exam.total_marks} marks</span>}
                  </div>
                  {exam.instructions && (
                    <p className="text-xs text-muted-foreground mt-2 line-clamp-2">
                      {exam.instructions}
                    </p>
                  )}
                  {exam.is_published && <ExamStatsCard examId={exam.id} />}
                </div>
                {canEdit && (
                  <div className="flex items-center gap-1.5 shrink-0">
                    <EditExamDialog exam={exam} onDone={() => router.refresh()} />
                    <AddExamQuestionDialog exam={exam} onDone={() => router.refresh()} />
                    {!exam.is_published && (
                      <Button
                        size="sm"
                        variant="outline"
                        className="text-xs h-7"
                        onClick={() => setPublishId(exam.id)}
                      >
                        <Globe className="h-3 w-3 mr-1" /> Publish
                      </Button>
                    )}
                    <Button
                      size="sm"
                      variant="ghost"
                      className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                      title="Delete Exam"
                      onClick={() => setDeleteExamItem(exam)}
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        ))
      )}

      {deleteExamItem && (
        <DeleteExamDialog
          exam={deleteExamItem}
          open={!!deleteExamItem}
          onOpenChange={(open) => !open && setDeleteExamItem(null)}
          onDone={() => router.refresh()}
        />
      )}

      {/* Create dialog */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-md" aria-describedby={undefined}>
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
              <Textarea
                id="ex-inst"
                rows={3}
                placeholder="Guidelines, rules, or instructions for students…"
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
