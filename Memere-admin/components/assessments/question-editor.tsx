"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  Check,
  ChevronDown,
  ChevronUp,
  Loader2,
  Pencil,
  Plus,
  Trash2,
  X,
} from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { clientAction } from "@/lib/client-action";
import type { AuthoredQuestion } from "@/lib/api/schemas";

/**
 * Question-by-question editing for a quiz or exam.
 *
 * The answer key is visible here by design — an author cannot correct a key they
 * cannot see. The data arrives from the teacher/admin-gated authoring endpoints,
 * which the backend re-checks for ownership; the student-facing API never returns
 * is_correct at all.
 *
 * `scope` decides which endpoint an edit goes to. A quiz question is addressed by
 * its own id; an exam question is addressed through its exam, because the
 * per-exam marks and order live on the join row and the exam's total_marks has to
 * be re-derived whenever they change.
 */
type Scope = { kind: "quiz"; quizId: string } | { kind: "exam"; examId: string };

const QUESTION_TYPES = [
  { value: "multiple_choice", label: "Multiple choice" },
  { value: "true_false", label: "True / False" },
  { value: "short_answer", label: "Short answer" },
] as const;

interface DraftAnswer {
  text: string;
  is_correct: boolean;
}

interface Draft {
  text: string;
  type: string;
  marks: number;
  explanation: string;
  subject: string;
  topic: string;
  answers: DraftAnswer[];
}

function toDraft(q: AuthoredQuestion): Draft {
  return {
    text: q.text,
    type: q.type || "multiple_choice",
    marks: q.marks ?? q.points ?? 1,
    explanation: q.explanation ?? "",
    subject: q.subject ?? "",
    topic: q.topic ?? "",
    answers: q.answers.map((a) => ({ text: a.text, is_correct: a.is_correct })),
  };
}

export function QuestionList({
  scope,
  questions,
}: {
  scope: Scope;
  questions: AuthoredQuestion[];
}) {
  const router = useRouter();
  const [editing, setEditing] = useState<AuthoredQuestion | null>(null);
  const [deleting, setDeleting] = useState<AuthoredQuestion | null>(null);
  const [movingId, setMovingId] = useState<string | null>(null);

  // Questions arrive ordered; reordering sends the whole id list because the
  // backend requires a complete permutation (a partial list would silently drop
  // a question from the order).
  async function move(index: number, direction: -1 | 1) {
    const target = index + direction;
    if (target < 0 || target >= questions.length) return;

    const ids = questions.map((q) => q.id);
    [ids[index], ids[target]] = [ids[target], ids[index]];

    setMovingId(questions[index].id);
    const url =
      scope.kind === "quiz"
        ? `/api/teacher/quizzes/${scope.quizId}/questions/order`
        : `/api/teacher/exams/${scope.examId}/questions/order`;
    try {
      await clientAction(url, { question_ids: ids }, "PUT");
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to reorder.");
    } finally {
      setMovingId(null);
    }
  }

  if (questions.length === 0) {
    return (
      <div className="rounded-md border border-dashed py-10 text-center text-sm text-muted-foreground">
        No questions yet. Add one, or import a whole set from a file.
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-3">
      {questions.map((q, i) => (
        <div key={q.id} className="rounded-md border p-4">
          <div className="flex items-start justify-between gap-3">
            <div className="min-w-0">
              <p className="text-sm font-medium">
                {i + 1}. {q.text}
              </p>
              <div className="mt-1 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                <Badge variant="secondary" className="text-xs">
                  {q.marks ?? q.points ?? 1} {scope.kind === "exam" ? "marks" : "pts"}
                </Badge>
                <span>{QUESTION_TYPES.find((t) => t.value === q.type)?.label ?? q.type}</span>
                {q.subject && <span>· {q.subject}</span>}
                {q.topic && <span>· {q.topic}</span>}
              </div>
            </div>

            <div className="flex shrink-0 items-center gap-1">
              <Button
                size="sm"
                variant="ghost"
                className="h-7 w-7 p-0"
                title="Move up"
                disabled={i === 0 || movingId !== null}
                onClick={() => move(i, -1)}
              >
                {movingId === q.id ? (
                  <Loader2 className="h-3.5 w-3.5 animate-spin" />
                ) : (
                  <ChevronUp className="h-3.5 w-3.5" />
                )}
              </Button>
              <Button
                size="sm"
                variant="ghost"
                className="h-7 w-7 p-0"
                title="Move down"
                disabled={i === questions.length - 1 || movingId !== null}
                onClick={() => move(i, 1)}
              >
                <ChevronDown className="h-3.5 w-3.5" />
              </Button>
              <Button size="sm" variant="ghost" className="h-7 text-xs" onClick={() => setEditing(q)}>
                <Pencil className="mr-1 h-3 w-3" /> Edit
              </Button>
              <Button
                size="sm"
                variant="ghost"
                className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                title={scope.kind === "exam" ? "Remove from exam" : "Delete question"}
                onClick={() => setDeleting(q)}
              >
                <Trash2 className="h-3.5 w-3.5" />
              </Button>
            </div>
          </div>

          <ul className="mt-3 space-y-1">
            {q.answers.map((a, j) => (
              <li
                key={a.id ?? j}
                className={
                  a.is_correct
                    ? "flex items-center gap-2 text-sm font-medium text-emerald-700 dark:text-emerald-400"
                    : "flex items-center gap-2 text-sm text-muted-foreground"
                }
              >
                {a.is_correct ? (
                  <Check className="h-3.5 w-3.5 shrink-0" />
                ) : (
                  <span className="h-3.5 w-3.5 shrink-0" />
                )}
                <span>
                  {String.fromCharCode(65 + j)}. {a.text}
                </span>
              </li>
            ))}
          </ul>

          {q.explanation && (
            <p className="mt-3 border-t pt-2 text-xs text-muted-foreground">
              <span className="font-medium">Explanation: </span>
              {q.explanation}
            </p>
          )}
        </div>
      ))}

      {editing && (
        <EditQuestionDialog
          scope={scope}
          question={editing}
          onClose={() => setEditing(null)}
          onSaved={() => {
            setEditing(null);
            router.refresh();
          }}
        />
      )}

      {deleting && (
        <DeleteQuestionDialog
          scope={scope}
          question={deleting}
          onClose={() => setDeleting(null)}
          onDeleted={() => {
            setDeleting(null);
            router.refresh();
          }}
        />
      )}
    </div>
  );
}

function EditQuestionDialog({
  scope,
  question,
  onClose,
  onSaved,
}: {
  scope: Scope;
  question: AuthoredQuestion;
  onClose: () => void;
  onSaved: () => void;
}) {
  const [draft, setDraft] = useState<Draft>(() => toDraft(question));
  const [saving, setSaving] = useState(false);

  function setAnswer(index: number, patch: Partial<DraftAnswer>) {
    setDraft((d) => ({
      ...d,
      answers: d.answers.map((a, i) => (i === index ? { ...a, ...patch } : a)),
    }));
  }

  // True/False and short_answer are single-key by nature, so ticking one option
  // clears the rest. Multiple choice allows several correct options.
  function toggleCorrect(index: number, value: boolean) {
    if (draft.type === "multiple_choice") {
      setAnswer(index, { is_correct: value });
      return;
    }
    setDraft((d) => ({
      ...d,
      answers: d.answers.map((a, i) => ({ ...a, is_correct: i === index ? value : false })),
    }));
  }

  function changeType(type: string) {
    setDraft((d) => {
      if (type === "true_false") {
        return {
          ...d,
          type,
          answers: [
            { text: "True", is_correct: true },
            { text: "False", is_correct: false },
          ],
        };
      }
      return { ...d, type };
    });
  }

  async function save() {
    const trimmed = draft.answers.map((a) => ({ ...a, text: a.text.trim() }));
    if (!draft.text.trim()) {
      toast.error("The question prompt cannot be empty.");
      return;
    }
    if (trimmed.length < 2) {
      toast.error("A question needs at least 2 answer options.");
      return;
    }
    if (trimmed.some((a) => !a.text)) {
      toast.error("Every answer option needs text.");
      return;
    }
    // Without a correct option the question cannot be graded at all, so this is
    // blocked here as well as server-side.
    if (!trimmed.some((a) => a.is_correct)) {
      toast.error("Mark at least one option as correct.");
      return;
    }

    const answers = trimmed.map((a, i) => ({
      text: a.text,
      is_correct: a.is_correct,
      order_index: i,
    }));

    const common = {
      text: draft.text.trim(),
      type: draft.type,
      explanation: draft.explanation.trim() || null,
      subject: draft.subject.trim() || null,
      topic: draft.topic.trim() || null,
      answers,
    };

    const [url, body] =
      scope.kind === "quiz"
        ? [`/api/teacher/questions/${question.id}`, { ...common, points: draft.marks }]
        : [
            `/api/teacher/exams/${scope.examId}/questions/${question.id}`,
            { ...common, marks: draft.marks },
          ];

    setSaving(true);
    try {
      await clientAction(url, body, "PUT");
      toast.success("Question saved.");
      onSaved();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save question.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <Dialog open onOpenChange={(o) => !o && !saving && onClose()}>
      <DialogContent className="max-w-xl max-h-[88vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit question</DialogTitle>
          <DialogDescription className="text-xs">
            The ticked options are graded as correct.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col gap-4 py-2">
          <div className="grid gap-1.5">
            <Label htmlFor="eq-text">Question prompt</Label>
            <textarea
              id="eq-text"
              rows={3}
              className="flex w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              value={draft.text}
              onChange={(e) => setDraft((d) => ({ ...d, text: e.target.value }))}
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="grid gap-1.5">
              <Label>Type</Label>
              <Select value={draft.type} onValueChange={changeType}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {QUESTION_TYPES.map((t) => (
                    <SelectItem key={t.value} value={t.value}>
                      {t.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="eq-marks">{scope.kind === "exam" ? "Marks" : "Points"}</Label>
              <Input
                id="eq-marks"
                type="number"
                min="1"
                value={draft.marks}
                onChange={(e) => setDraft((d) => ({ ...d, marks: Number(e.target.value) }))}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="grid gap-1.5">
              <Label htmlFor="eq-subject">Subject (optional)</Label>
              <Input
                id="eq-subject"
                placeholder="e.g. Physics"
                value={draft.subject}
                onChange={(e) => setDraft((d) => ({ ...d, subject: e.target.value }))}
              />
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="eq-topic">Topic (optional)</Label>
              <Input
                id="eq-topic"
                placeholder="e.g. Kinematics"
                value={draft.topic}
                onChange={(e) => setDraft((d) => ({ ...d, topic: e.target.value }))}
              />
            </div>
          </div>

          <div className="grid gap-1.5">
            <Label htmlFor="eq-exp">Explanation (shown after grading)</Label>
            <textarea
              id="eq-exp"
              rows={2}
              className="flex w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              value={draft.explanation}
              onChange={(e) => setDraft((d) => ({ ...d, explanation: e.target.value }))}
            />
          </div>

          <div className="grid gap-2 border-t pt-3">
            <div className="flex items-center justify-between">
              <Label className="text-xs font-semibold">
                Answer options — tick the correct one
              </Label>
              {draft.type !== "true_false" && draft.answers.length < 8 && (
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="h-6 px-2 text-xs text-primary"
                  onClick={() =>
                    setDraft((d) => ({ ...d, answers: [...d.answers, { text: "", is_correct: false }] }))
                  }
                >
                  <Plus className="mr-1 h-3 w-3" /> Add option
                </Button>
              )}
            </div>

            {draft.answers.map((a, i) => (
              <div key={i} className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={a.is_correct}
                  onChange={(e) => toggleCorrect(i, e.target.checked)}
                  className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                  title="Mark as correct"
                />
                <span className="w-4 text-xs text-muted-foreground">
                  {String.fromCharCode(65 + i)}
                </span>
                <Input
                  placeholder={`Option ${String.fromCharCode(65 + i)}`}
                  disabled={draft.type === "true_false"}
                  value={a.text}
                  onChange={(e) => setAnswer(i, { text: e.target.value })}
                />
                {draft.type !== "true_false" && draft.answers.length > 2 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 shrink-0 p-0 text-muted-foreground hover:text-destructive"
                    onClick={() =>
                      setDraft((d) => ({ ...d, answers: d.answers.filter((_, j) => j !== i) }))
                    }
                  >
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </div>
            ))}
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={onClose} disabled={saving}>
            Cancel
          </Button>
          <Button type="button" onClick={save} disabled={saving}>
            {saving ? (
              <>
                <Loader2 className="mr-1.5 h-4 w-4 animate-spin" />
                Saving…
              </>
            ) : (
              "Save question"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function DeleteQuestionDialog({
  scope,
  question,
  onClose,
  onDeleted,
}: {
  scope: Scope;
  question: AuthoredQuestion;
  onClose: () => void;
  onDeleted: () => void;
}) {
  const [deleting, setDeleting] = useState(false);

  async function handleDelete() {
    const url =
      scope.kind === "quiz"
        ? `/api/teacher/questions/${question.id}`
        : `/api/teacher/exams/${scope.examId}/questions/${question.id}`;
    setDeleting(true);
    try {
      await clientAction(url, undefined, "DELETE");
      toast.success(scope.kind === "exam" ? "Question removed from exam." : "Question deleted.");
      onDeleted();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete question.");
    } finally {
      setDeleting(false);
    }
  }

  return (
    <Dialog open onOpenChange={(o) => !o && !deleting && onClose()}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-sm font-semibold text-destructive">
            <Trash2 className="h-4 w-4" />
            {scope.kind === "exam" ? "Remove from exam?" : "Delete question?"}
          </DialogTitle>
          <DialogDescription className="pt-1 text-xs">
            {scope.kind === "exam" ? (
              <>
                &ldquo;{question.text}&rdquo; will be detached from this exam and the total marks
                recalculated. The question stays in the question bank, so any other exam using it is
                unaffected.
              </>
            ) : (
              <>&ldquo;{question.text}&rdquo; and all of its answer options will be deleted.</>
            )}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="mt-2">
          <Button variant="outline" size="sm" onClick={onClose} disabled={deleting}>
            Cancel
          </Button>
          <Button variant="destructive" size="sm" onClick={handleDelete} disabled={deleting}>
            {deleting ? (
              <>
                <Loader2 className="mr-1 h-3.5 w-3.5 animate-spin" />
                Working…
              </>
            ) : scope.kind === "exam" ? (
              "Remove"
            ) : (
              "Delete"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
