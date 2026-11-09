"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Loader2, Plus } from "lucide-react";
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
import type { Quiz } from "@/lib/api/schemas";

/**
 * Creates a quiz from the top-level Quizzes page.
 *
 * Leaving the course unset creates a standalone quiz, which only an admin may do
 * — a course-less quiz has no owning teacher to scope editing to. Teachers are
 * therefore given the course picker with no "none" option.
 */
export function CreateQuizDialog({
  courses,
  canCreateStandalone,
}: {
  courses: { id: string; title: string }[];
  canCreateStandalone: boolean;
}) {
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [title, setTitle] = useState("");
  const [courseId, setCourseId] = useState<string>(canCreateStandalone ? "none" : "");
  const [minutes, setMinutes] = useState(30);
  const [passPercentage, setPassPercentage] = useState(60);
  const [maxAttempts, setMaxAttempts] = useState(3);
  const [randomize, setRandomize] = useState(false);
  const router = useRouter();

  async function create() {
    if (!title.trim()) {
      toast.error("A title is required.");
      return;
    }
    if (!canCreateStandalone && !courseId) {
      toast.error("Pick a course for this quiz.");
      return;
    }

    setSaving(true);
    try {
      const quiz = await clientAction<Quiz>("/api/teacher/quizzes", {
        title: title.trim(),
        course_id: courseId && courseId !== "none" ? courseId : null,
        time_limit_seconds: minutes > 0 ? minutes * 60 : null,
        pass_percentage: passPercentage,
        max_attempts: maxAttempts > 0 ? maxAttempts : null,
        randomize_questions: randomize,
      });
      toast.success("Quiz created. Add its questions next.");
      setOpen(false);
      setTitle("");
      router.push(`/quizzes/${quiz.id}`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create quiz.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <Button size="sm" onClick={() => setOpen(true)}>
        <Plus className="mr-1.5 h-4 w-4" /> New quiz
      </Button>

      <Dialog open={open} onOpenChange={(o) => !saving && setOpen(o)}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Create quiz</DialogTitle>
            <DialogDescription className="text-xs">
              {canCreateStandalone
                ? "Leave the course unset to create a standalone quiz that is not tied to any course."
                : "Choose which of your courses this quiz belongs to."}
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="nq-title">Title</Label>
              <Input
                id="nq-title"
                placeholder="e.g. Unit 1 Practice Quiz"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>

            <div className="grid gap-1.5">
              <Label>Course</Label>
              <Select value={courseId} onValueChange={setCourseId}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a course" />
                </SelectTrigger>
                <SelectContent>
                  {canCreateStandalone && (
                    <SelectItem value="none">Standalone — not tied to a course</SelectItem>
                  )}
                  {courses.map((c) => (
                    <SelectItem key={c.id} value={c.id}>
                      {c.title}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="nq-time">Time limit (minutes, 0 = none)</Label>
                <Input
                  id="nq-time"
                  type="number"
                  min="0"
                  value={minutes}
                  onChange={(e) => setMinutes(Number(e.target.value))}
                />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="nq-pass">Pass %</Label>
                <Input
                  id="nq-pass"
                  type="number"
                  min="0"
                  max="100"
                  value={passPercentage}
                  onChange={(e) => setPassPercentage(Number(e.target.value))}
                />
              </div>
            </div>

            <div className="grid grid-cols-2 items-center gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="nq-att">Max attempts (0 = unlimited)</Label>
                <Input
                  id="nq-att"
                  type="number"
                  min="0"
                  value={maxAttempts}
                  onChange={(e) => setMaxAttempts(Number(e.target.value))}
                />
              </div>
              <div className="flex items-center gap-2 pt-5">
                <input
                  id="nq-rand"
                  type="checkbox"
                  className="rounded"
                  checked={randomize}
                  onChange={(e) => setRandomize(e.target.checked)}
                />
                <Label htmlFor="nq-rand" className="cursor-pointer text-sm">
                  Randomize questions
                </Label>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)} disabled={saving}>
              Cancel
            </Button>
            <Button onClick={create} disabled={saving}>
              {saving ? (
                <>
                  <Loader2 className="mr-1.5 h-4 w-4 animate-spin" />
                  Creating…
                </>
              ) : (
                "Create quiz"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
