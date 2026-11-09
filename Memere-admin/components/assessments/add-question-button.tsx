"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Loader2, Plus, X } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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

type Scope = { kind: "quiz"; quizId: string } | { kind: "exam"; examId: string };

const BLANK_OPTIONS = [
  { text: "", is_correct: false },
  { text: "", is_correct: false },
  { text: "", is_correct: false },
  { text: "", is_correct: false },
];

/** Adds a single question to a quiz or exam. Bulk sets come in via the importer. */
export function AddQuestionButton({
  scope,
  nextOrderIndex,
}: {
  scope: Scope;
  /** Appends the question after the existing ones. */
  nextOrderIndex: number;
}) {
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [text, setText] = useState("");
  const [type, setType] = useState("multiple_choice");
  const [marks, setMarks] = useState(scope.kind === "exam" ? 1 : 5);
  const [explanation, setExplanation] = useState("");
  const [subject, setSubject] = useState("");
  const [topic, setTopic] = useState("");
  const [answers, setAnswers] = useState(BLANK_OPTIONS);
  const router = useRouter();

  function reset() {
    setText("");
    setType("multiple_choice");
    setMarks(scope.kind === "exam" ? 1 : 5);
    setExplanation("");
    setSubject("");
    setTopic("");
    setAnswers(BLANK_OPTIONS);
  }

  function changeType(next: string) {
    setType(next);
    if (next === "true_false") {
      setAnswers([
        { text: "True", is_correct: true },
        { text: "False", is_correct: false },
      ]);
    } else {
      setAnswers(BLANK_OPTIONS);
    }
  }

  // True/False and short_answer are single-key by nature; multiple choice allows
  // several correct options.
  function toggleCorrect(index: number, value: boolean) {
    setAnswers((prev) =>
      prev.map((a, i) =>
        type === "multiple_choice"
          ? i === index
            ? { ...a, is_correct: value }
            : a
          : { ...a, is_correct: i === index ? value : false }
      )
    );
  }

  async function save() {
    const trimmed = answers.map((a) => ({ ...a, text: a.text.trim() })).filter((a) => a.text);
    if (!text.trim()) {
      toast.error("The question prompt cannot be empty.");
      return;
    }
    if (trimmed.length < 2) {
      toast.error("A question needs at least 2 answer options.");
      return;
    }
    if (!trimmed.some((a) => a.is_correct)) {
      toast.error("Mark at least one option as correct.");
      return;
    }

    const payload = {
      text: text.trim(),
      type,
      explanation: explanation.trim() || null,
      subject: subject.trim() || null,
      topic: topic.trim() || null,
      order_index: nextOrderIndex,
      answers: trimmed.map((a, i) => ({ text: a.text, is_correct: a.is_correct, order_index: i })),
      ...(scope.kind === "exam" ? { marks } : { points: marks }),
    };

    const url =
      scope.kind === "quiz"
        ? `/api/teacher/quizzes/${scope.quizId}/questions`
        : `/api/teacher/exams/${scope.examId}/questions`;

    setSaving(true);
    try {
      await clientAction(url, payload);
      toast.success("Question added.");
      setOpen(false);
      reset();
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add question.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <Button size="sm" variant="outline" onClick={() => setOpen(true)}>
        <Plus className="mr-1.5 h-4 w-4" /> Add question
      </Button>

      <Dialog
        open={open}
        onOpenChange={(o) => {
          if (saving) return;
          setOpen(o);
          if (!o) reset();
        }}
      >
        <DialogContent className="max-h-[88vh] max-w-xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Add question</DialogTitle>
            <DialogDescription className="text-xs">
              The ticked options are graded as correct.
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="aq-text">Question prompt</Label>
              <textarea
                id="aq-text"
                rows={3}
                placeholder="Type your question here…"
                className="flex w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                value={text}
                onChange={(e) => setText(e.target.value)}
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label>Type</Label>
                <Select value={type} onValueChange={changeType}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="multiple_choice">Multiple choice</SelectItem>
                    <SelectItem value="true_false">True / False</SelectItem>
                    <SelectItem value="short_answer">Short answer</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="aq-marks">{scope.kind === "exam" ? "Marks" : "Points"}</Label>
                <Input
                  id="aq-marks"
                  type="number"
                  min="1"
                  value={marks}
                  onChange={(e) => setMarks(Number(e.target.value))}
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="aq-subject">Subject (optional)</Label>
                <Input
                  id="aq-subject"
                  placeholder="e.g. Physics"
                  value={subject}
                  onChange={(e) => setSubject(e.target.value)}
                />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="aq-topic">Topic (optional)</Label>
                <Input
                  id="aq-topic"
                  placeholder="e.g. Kinematics"
                  value={topic}
                  onChange={(e) => setTopic(e.target.value)}
                />
              </div>
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="aq-exp">Explanation (shown after grading)</Label>
              <textarea
                id="aq-exp"
                rows={2}
                className="flex w-full resize-none rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                value={explanation}
                onChange={(e) => setExplanation(e.target.value)}
              />
            </div>

            <div className="grid gap-2 border-t pt-3">
              <div className="flex items-center justify-between">
                <Label className="text-xs font-semibold">
                  Answer options — tick the correct one
                </Label>
                {type !== "true_false" && answers.length < 8 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-6 px-2 text-xs text-primary"
                    onClick={() => setAnswers((prev) => [...prev, { text: "", is_correct: false }])}
                  >
                    <Plus className="mr-1 h-3 w-3" /> Add option
                  </Button>
                )}
              </div>

              {answers.map((a, i) => (
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
                    disabled={type === "true_false"}
                    value={a.text}
                    onChange={(e) =>
                      setAnswers((prev) =>
                        prev.map((x, j) => (j === i ? { ...x, text: e.target.value } : x))
                      )
                    }
                  />
                  {type !== "true_false" && answers.length > 2 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="h-8 w-8 shrink-0 p-0 text-muted-foreground hover:text-destructive"
                      onClick={() => setAnswers((prev) => prev.filter((_, j) => j !== i))}
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              ))}
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setOpen(false);
                reset();
              }}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button onClick={save} disabled={saving}>
              {saving ? (
                <>
                  <Loader2 className="mr-1.5 h-4 w-4 animate-spin" />
                  Adding…
                </>
              ) : (
                "Add question"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
