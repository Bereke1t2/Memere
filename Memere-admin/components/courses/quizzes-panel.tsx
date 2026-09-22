"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Loader2, Trash2, Pencil, HelpCircle, X } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
  CreateQuizInputSchema,
  AddQuizQuestionInputSchema,
  type CreateQuizInput,
  type AddQuizQuestionInput,
  type Quiz,
  type Lesson,
} from "@/lib/api/schemas";

interface QuizzesPanelProps {
  courseId: string;
  quizzes: Quiz[];
  lessons?: Lesson[];
  canEdit?: boolean;
}

function AddQuestionDialog({
  quizId,
  onDone,
}: {
  quizId: string;
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
  } = useForm<AddQuizQuestionInput>({
    resolver: zodResolver(AddQuizQuestionInputSchema),
    defaultValues: {
      type: "multiple_choice",
      points: 5,
      order_index: 0,
      subject: "",
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

  async function onSubmit(values: AddQuizQuestionInput) {
    const hasCorrect = values.answers.some((a) => a.is_correct);
    if (!hasCorrect) {
      toast.error("Please mark at least one answer as correct.");
      return;
    }
    try {
      await clientAction(`/api/teacher/quizzes/${quizId}/questions`, values);
      toast.success("Question added.");
      setOpen(false);
      reset();
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
            <DialogTitle>Add Quiz Question</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="q-text">Question Prompt</Label>
              <textarea
                id="q-text"
                rows={3}
                placeholder="Type your question here…"
                className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
                {...register("text")}
              />
              {errors.text && <p className="text-xs text-destructive">{errors.text.message}</p>}
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Question Type</Label>
                <Select
                  value={selectedType}
                  onValueChange={(v) => handleTypeChange(v as "multiple_choice" | "true_false")}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="multiple_choice">Multiple Choice</SelectItem>
                    <SelectItem value="true_false">True / False</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="q-pts">Points (Marks)</Label>
                <Input
                  id="q-pts"
                  type="number"
                  min="1"
                  {...register("points", { valueAsNumber: true })}
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="q-sub">Subject (optional)</Label>
                <Input id="q-sub" placeholder="e.g. Physics" {...register("subject")} />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="q-top">Topic (optional)</Label>
                <Input id="q-top" placeholder="e.g. Kinematics" {...register("topic")} />
              </div>
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="q-exp">Explanation (optional)</Label>
              <Input id="q-exp" placeholder="Shown after student submits" {...register("explanation")} />
            </div>

            <div className="grid gap-2 border-t pt-3">
              <div className="flex items-center justify-between">
                <Label className="text-xs font-semibold">Answer Choices (Check the correct answer)</Label>
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
                        // For True/False, ensure single selection
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
                    placeholder={`Answer Choice ${i + 1}`}
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
                  {Array.isArray(errors.answers) ? "Please fill all answer choices." : errors.answers.message}
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

function EditQuizDialog({
  quiz,
  lessons,
  onDone,
}: {
  quiz: Quiz;
  lessons?: Lesson[];
  onDone: () => void;
}) {
  const [open, setOpen] = useState(false);
  const [minutes, setMinutes] = useState(
    quiz.time_limit_seconds ? Math.round(quiz.time_limit_seconds / 60) : 30
  );

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<CreateQuizInput>({
    resolver: zodResolver(CreateQuizInputSchema),
    defaultValues: {
      title: quiz.title,
      lesson_id: quiz.lesson_id ?? null,
      time_limit_seconds: quiz.time_limit_seconds ?? 1800,
      pass_percentage: quiz.pass_percentage ?? 60,
      randomize_questions: quiz.randomize_questions ?? false,
      max_attempts: quiz.max_attempts ?? 3,
    },
  });

  async function onSubmit(values: CreateQuizInput) {
    try {
      await clientAction(`/api/teacher/quizzes/${quiz.id}`, {
        ...values,
        time_limit_seconds: minutes > 0 ? minutes * 60 : null,
        lesson_id: values.lesson_id === "none" ? null : values.lesson_id,
      }, "PUT");
      toast.success("Quiz updated.");
      setOpen(false);
      onDone();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update quiz.");
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
            <DialogTitle>Edit Quiz</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="eqz-title">Title</Label>
              <Input id="eqz-title" placeholder="e.g. Chapter 1 Quiz" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>

            {lessons && lessons.length > 0 && (
              <div className="grid gap-1.5">
                <Label>Associated Lesson (optional)</Label>
                <Select
                  defaultValue={quiz.lesson_id ?? "none"}
                  onValueChange={(v) => setValue("lesson_id", v === "none" ? null : v)}
                >
                  <SelectTrigger><SelectValue placeholder="Whole course" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">Whole course (Not linked to specific lesson)</SelectItem>
                    {lessons.map((l) => (
                      <SelectItem key={l.id} value={l.id}>
                        {l.title}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="eqz-time">Time limit (minutes)</Label>
                <Input
                  id="eqz-time"
                  type="number"
                  min="0"
                  value={minutes}
                  onChange={(e) => setMinutes(Number(e.target.value))}
                />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="eqz-pass">Pass % (0–100)</Label>
                <Input
                  id="eqz-pass"
                  type="number"
                  min="0"
                  max="100"
                  {...register("pass_percentage", { valueAsNumber: true })}
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3 items-center">
              <div className="grid gap-1.5">
                <Label htmlFor="eqz-att">Max attempts</Label>
                <Input
                  id="eqz-att"
                  type="number"
                  min="1"
                  {...register("max_attempts", { valueAsNumber: true })}
                />
              </div>
              <div className="flex items-center gap-2 pt-5">
                <input
                  type="checkbox"
                  id="eqz-rand"
                  {...register("randomize_questions")}
                  className="rounded"
                />
                <Label htmlFor="eqz-rand" className="cursor-pointer text-sm">
                  Randomize questions
                </Label>
              </div>
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

function DeleteQuizDialog({
  quiz,
  open,
  onOpenChange,
  onDone,
}: {
  quiz: Quiz;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDone: () => void;
}) {
  const [deleting, setDeleting] = useState(false);

  async function handleDelete() {
    setDeleting(true);
    try {
      await clientAction(`/api/teacher/quizzes/${quiz.id}`, undefined, "DELETE");
      toast.success(`Quiz "${quiz.title}" deleted.`);
      onOpenChange(false);
      onDone();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete quiz.");
    } finally {
      setDeleting(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={(o) => !deleting && onOpenChange(o)}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle className="text-sm font-semibold flex items-center gap-2 text-destructive">
            <Trash2 className="h-4 w-4" /> Delete Quiz?
          </DialogTitle>
          <DialogDescription className="text-xs pt-1">
            This will permanently delete <strong>{quiz.title}</strong> and all questions inside it.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="mt-2">
          <Button variant="outline" size="sm" onClick={() => onOpenChange(false)} disabled={deleting}>
            Cancel
          </Button>
          <Button variant="destructive" size="sm" onClick={handleDelete} disabled={deleting}>
            {deleting ? <><Loader2 className="h-3.5 w-3.5 mr-1 animate-spin" />Deleting…</> : "Delete Quiz"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export function QuizzesPanel({
  courseId,
  quizzes,
  lessons,
  canEdit = true,
}: QuizzesPanelProps) {
  const [createOpen, setCreateOpen] = useState(false);
  const [deleteQuizItem, setDeleteQuizItem] = useState<Quiz | null>(null);
  const [minutes, setMinutes] = useState(30);
  const router = useRouter();

  const {
    register,
    handleSubmit,
    setValue,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CreateQuizInput>({
    resolver: zodResolver(CreateQuizInputSchema),
    defaultValues: {
      time_limit_seconds: 1800,
      pass_percentage: 60,
      randomize_questions: false,
      max_attempts: 3,
    },
  });

  async function onCreate(values: CreateQuizInput) {
    try {
      await clientAction(`/api/teacher/courses/${courseId}/quizzes`, {
        ...values,
        time_limit_seconds: minutes > 0 ? minutes * 60 : null,
        lesson_id: values.lesson_id === "none" ? null : values.lesson_id,
      });
      toast.success("Quiz created.");
      setCreateOpen(false);
      reset();
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create quiz.");
    }
  }

  return (
    <div className="flex flex-col gap-4">
      {canEdit && (
        <div className="flex justify-end">
          <Button size="sm" onClick={() => setCreateOpen(true)}>
            <Plus className="h-4 w-4 mr-1.5" /> New Quiz
          </Button>
        </div>
      )}

      {quizzes.length === 0 ? (
        <div className="py-10 text-center text-sm text-muted-foreground">
          No quizzes yet. Create a quiz to test your students.
        </div>
      ) : (
        quizzes.map((quiz) => (
          <Card key={quiz.id} className="border shadow-none">
            <CardContent className="p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="font-medium text-sm">{quiz.title}</p>
                  <div className="flex items-center gap-3 mt-1 text-xs text-muted-foreground flex-wrap">
                    {quiz.time_limit_seconds != null && (
                      <span>{Math.round(quiz.time_limit_seconds / 60)}m limit</span>
                    )}
                    {quiz.pass_percentage != null && <span>{quiz.pass_percentage}% to pass</span>}
                    {quiz.max_attempts != null && <span>Max {quiz.max_attempts} attempts</span>}
                    {quiz.randomize_questions && <span>Randomized</span>}
                    {quiz.question_count != null && (
                      <Badge variant="secondary" className="text-xs">
                        {quiz.question_count} questions
                      </Badge>
                    )}
                  </div>
                </div>
                {canEdit && (
                  <div className="flex items-center gap-1.5">
                    <EditQuizDialog
                      quiz={quiz}
                      lessons={lessons}
                      onDone={() => router.refresh()}
                    />
                    <AddQuestionDialog quizId={quiz.id} onDone={() => router.refresh()} />
                    <Button
                      size="sm"
                      variant="ghost"
                      className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                      title="Delete Quiz"
                      onClick={() => setDeleteQuizItem(quiz)}
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

      {deleteQuizItem && (
        <DeleteQuizDialog
          quiz={deleteQuizItem}
          open={!!deleteQuizItem}
          onOpenChange={(open) => !open && setDeleteQuizItem(null)}
          onDone={() => router.refresh()}
        />
      )}

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-md" aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Create Quiz</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onCreate)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="qz-title">Title</Label>
              <Input id="qz-title" placeholder="e.g. Chapter 1 Quiz" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>

            {lessons && lessons.length > 0 && (
              <div className="grid gap-1.5">
                <Label>Associated Lesson (optional)</Label>
                <Select
                  defaultValue="none"
                  onValueChange={(v) => setValue("lesson_id", v === "none" ? null : v)}
                >
                  <SelectTrigger><SelectValue placeholder="Whole course" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">Whole course (Not linked to specific lesson)</SelectItem>
                    {lessons.map((l) => (
                      <SelectItem key={l.id} value={l.id}>
                        {l.title}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="qz-time">Time limit (minutes)</Label>
                <Input
                  id="qz-time"
                  type="number"
                  min="0"
                  value={minutes}
                  onChange={(e) => setMinutes(Number(e.target.value))}
                />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="qz-pass">Pass % (0–100)</Label>
                <Input
                  id="qz-pass"
                  type="number"
                  min="0"
                  max="100"
                  {...register("pass_percentage", { valueAsNumber: true })}
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3 items-center">
              <div className="grid gap-1.5">
                <Label htmlFor="qz-att">Max attempts</Label>
                <Input
                  id="qz-att"
                  type="number"
                  min="1"
                  {...register("max_attempts", { valueAsNumber: true })}
                />
              </div>
              <div className="flex items-center gap-2 pt-5">
                <input
                  type="checkbox"
                  id="qz-rand"
                  {...register("randomize_questions")}
                  className="rounded"
                />
                <Label htmlFor="qz-rand" className="cursor-pointer text-sm">
                  Randomize questions
                </Label>
              </div>
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setCreateOpen(false); reset(); }}>
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? <><Loader2 className="h-4 w-4 mr-1.5 animate-spin" />Creating…</> : "Create Quiz"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
