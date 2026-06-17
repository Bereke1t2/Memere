"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Loader2, Trash2, CheckCircle2 } from "lucide-react";
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
import {
  CreateQuizInputSchema, AddQuizQuestionInputSchema,
  type CreateQuizInput, type AddQuizQuestionInput, type Quiz,
} from "@/lib/api/schemas";

interface QuizzesPanelProps {
  courseId: string;
  quizzes: Quiz[];
}

function AddQuestionDialog({ quizId, onDone }: { quizId: string; onDone: () => void }) {
  const [open, setOpen] = useState(false);

  const { register, handleSubmit, control, setValue, watch, reset, formState: { errors, isSubmitting } } =
    useForm<AddQuizQuestionInput>({
      resolver: zodResolver(AddQuizQuestionInputSchema),
      defaultValues: {
        type: "multiple_choice",
        points: 5,
        order_index: 0,
        answers: [
          { text: "", is_correct: false },
          { text: "", is_correct: false },
          { text: "", is_correct: false },
          { text: "", is_correct: false },
        ],
      },
    });

  const { fields } = useFieldArray({ control, name: "answers" });
  const answers = watch("answers");

  async function onSubmit(values: AddQuizQuestionInput) {
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
        <DialogContent className="max-w-lg max-h-[80vh] overflow-y-auto">
          <DialogHeader><DialogTitle>Add Question</DialogTitle></DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="q-text">Question</Label>
              <textarea
                id="q-text"
                rows={2}
                className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-none"
                {...register("text")}
              />
              {errors.text && <p className="text-xs text-destructive">{errors.text.message}</p>}
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Type</Label>
                <Select defaultValue="multiple_choice" onValueChange={(v) => setValue("type", v as AddQuizQuestionInput["type"])}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="multiple_choice">Multiple choice</SelectItem>
                    <SelectItem value="true_false">True / False</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="q-pts">Points</Label>
                <Input id="q-pts" type="number" min="1" {...register("points", { valueAsNumber: true })} />
              </div>
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="q-exp">Explanation (optional)</Label>
              <Input id="q-exp" placeholder="Shown after attempt" {...register("explanation")} />
            </div>
            <div className="grid gap-2">
              <Label>Answers — tick the correct one(s)</Label>
              {fields.map((field, i) => (
                <div key={field.id} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={answers?.[i]?.is_correct ?? false}
                    onChange={(e) => setValue(`answers.${i}.is_correct`, e.target.checked)}
                    className="rounded"
                  />
                  <Input placeholder={`Answer ${i + 1}`} {...register(`answers.${i}.text`)} />
                </div>
              ))}
              {errors.answers && <p className="text-xs text-destructive">{Array.isArray(errors.answers) ? "Fix answers" : errors.answers.message}</p>}
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setOpen(false); reset(); }}>Cancel</Button>
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

export function QuizzesPanel({ courseId, quizzes }: QuizzesPanelProps) {
  const [createOpen, setCreateOpen] = useState(false);
  const router = useRouter();

  const { register, handleSubmit, setValue, reset, formState: { errors, isSubmitting } } =
    useForm<CreateQuizInput>({
      resolver: zodResolver(CreateQuizInputSchema),
      defaultValues: { time_limit_seconds: 1800, pass_percentage: 60, randomize_questions: false, max_attempts: 3 },
    });

  async function onCreate(values: CreateQuizInput) {
    try {
      await clientAction(`/api/teacher/courses/${courseId}/quizzes`, values);
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
      <div className="flex justify-end">
        <Button size="sm" onClick={() => setCreateOpen(true)}>
          <Plus className="h-4 w-4 mr-1.5" /> New Quiz
        </Button>
      </div>

      {quizzes.length === 0 ? (
        <div className="py-10 text-center text-sm text-muted-foreground">No quizzes yet.</div>
      ) : (
        quizzes.map((quiz) => (
          <Card key={quiz.id} className="border shadow-none">
            <CardContent className="p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="font-medium text-sm">{quiz.title}</p>
                  <div className="flex items-center gap-3 mt-1 text-xs text-muted-foreground">
                    {quiz.time_limit_seconds != null && (
                      <span>{Math.round(quiz.time_limit_seconds / 60)}m limit</span>
                    )}
                    {quiz.pass_percentage != null && <span>{quiz.pass_percentage}% to pass</span>}
                    {quiz.max_attempts != null && <span>Max {quiz.max_attempts} attempts</span>}
                    {quiz.question_count != null && (
                      <Badge variant="secondary" className="text-xs">{quiz.question_count} questions</Badge>
                    )}
                  </div>
                </div>
                <AddQuestionDialog quizId={quiz.id} onDone={() => router.refresh()} />
              </div>
            </CardContent>
          </Card>
        ))
      )}

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-sm">
          <DialogHeader><DialogTitle>Create Quiz</DialogTitle></DialogHeader>
          <form onSubmit={handleSubmit(onCreate)} className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="qz-title">Title</Label>
              <Input id="qz-title" placeholder="e.g. Chapter 1 Quiz" {...register("title")} />
              {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="qz-time">Time limit (seconds)</Label>
                <Input id="qz-time" type="number" min="0" {...register("time_limit_seconds", { valueAsNumber: true })} />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="qz-pass">Pass % (0–100)</Label>
                <Input id="qz-pass" type="number" min="0" max="100" {...register("pass_percentage", { valueAsNumber: true })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="qz-att">Max attempts</Label>
                <Input id="qz-att" type="number" min="1" {...register("max_attempts", { valueAsNumber: true })} />
              </div>
              <div className="flex items-center gap-2 pt-5">
                <input type="checkbox" id="qz-rand" {...register("randomize_questions")} className="rounded" />
                <Label htmlFor="qz-rand" className="cursor-pointer text-sm">Randomize</Label>
              </div>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setCreateOpen(false); reset(); }}>Cancel</Button>
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
