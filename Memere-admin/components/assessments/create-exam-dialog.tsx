"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Loader2, Plus } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
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
import type { Exam } from "@/lib/api/schemas";

/**
 * Creates an exam from the top-level Exams page.
 *
 * Leaving the course unset creates a standalone mock exam, which only an admin
 * may do — a course-less exam has no owning teacher to scope editing to.
 *
 * Total marks is deliberately absent: it is derived from the question set as
 * questions are added, so the header total can never drift from the parts.
 */
export function CreateExamDialog({
  courses,
  canCreateStandalone,
}: {
  courses: { id: string; title: string }[];
  canCreateStandalone: boolean;
}) {
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [title, setTitle] = useState("");
  const [subject, setSubject] = useState("");
  const [grade, setGrade] = useState(12);
  const [duration, setDuration] = useState(120);
  const [passMarks, setPassMarks] = useState(50);
  const [instructions, setInstructions] = useState("");
  const [courseId, setCourseId] = useState<string>(canCreateStandalone ? "none" : "");
  const router = useRouter();

  async function create() {
    if (!title.trim()) {
      toast.error("A title is required.");
      return;
    }
    if (!subject.trim()) {
      toast.error("A subject is required.");
      return;
    }
    if (!canCreateStandalone && !courseId) {
      toast.error("Pick a course for this exam.");
      return;
    }

    setSaving(true);
    try {
      const exam = await clientAction<Exam>("/api/teacher/exams", {
        title: title.trim(),
        subject: subject.trim(),
        grade,
        duration_minutes: duration,
        pass_marks: passMarks,
        instructions: instructions.trim() || null,
        course_id: courseId && courseId !== "none" ? courseId : null,
      });
      toast.success("Exam created. Add its questions next.");
      setOpen(false);
      setTitle("");
      router.push(`/exams/${exam.id}`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create exam.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <Button size="sm" onClick={() => setOpen(true)}>
        <Plus className="mr-1.5 h-4 w-4" /> New exam
      </Button>

      <Dialog open={open} onOpenChange={(o) => !saving && setOpen(o)}>
        <DialogContent className="max-h-[88vh] max-w-md overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create exam</DialogTitle>
            <DialogDescription className="text-xs">
              {canCreateStandalone
                ? "Leave the course unset to create a standalone mock exam. Total marks is calculated from the questions you add."
                : "Total marks is calculated from the questions you add."}
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col gap-4 py-2">
            <div className="grid gap-1.5">
              <Label htmlFor="ne-title">Title</Label>
              <Input
                id="ne-title"
                placeholder="e.g. Physics Mock Exam 2026"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="grid gap-1.5">
                <Label htmlFor="ne-subject">Subject</Label>
                <Input
                  id="ne-subject"
                  placeholder="e.g. Physics"
                  value={subject}
                  onChange={(e) => setSubject(e.target.value)}
                />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="ne-grade">Grade (1–13)</Label>
                <Input
                  id="ne-grade"
                  type="number"
                  min="1"
                  max="13"
                  value={grade}
                  onChange={(e) => setGrade(Number(e.target.value))}
                />
              </div>
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
                <Label htmlFor="ne-dur">Duration (minutes)</Label>
                <Input
                  id="ne-dur"
                  type="number"
                  min="1"
                  value={duration}
                  onChange={(e) => setDuration(Number(e.target.value))}
                />
              </div>
              <div className="grid gap-1.5">
                <Label htmlFor="ne-pass">Pass marks</Label>
                <Input
                  id="ne-pass"
                  type="number"
                  min="0"
                  value={passMarks}
                  onChange={(e) => setPassMarks(Number(e.target.value))}
                />
              </div>
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="ne-inst">Instructions (optional)</Label>
              <Textarea
                id="ne-inst"
                rows={3}
                placeholder="Shown to students before they start."
                value={instructions}
                onChange={(e) => setInstructions(e.target.value)}
              />
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
                "Create exam"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
