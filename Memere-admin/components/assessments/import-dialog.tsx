"use client";

import { useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Upload, Loader2, AlertTriangle, Check, X, FileJson, FileCode } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { clientAction } from "@/lib/client-action";
import { parseImportFile, resolvePreviewAnswers, ParseError, type ParseResult } from "@/lib/import/parse";
import type { ImportResult } from "@/lib/api/schemas";

/**
 * Upload-and-import flow for a whole quiz or exam.
 *
 * The preview step is deliberately not skippable. A file decides the answer key
 * for every question in a graded assessment, and HTML in particular cannot be
 * interpreted with certainty, so the operator sees each resolved key and
 * confirms before anything is written. The import itself is atomic server-side:
 * either the whole file lands or none of it does.
 */
interface HeaderDraft {
  title: string;
  subject: string;
  grade: string;
  durationMinutes: string;
  passMarks: string;
  passPercentage: string;
  timeLimitMinutes: string;
}

const emptyHeader: HeaderDraft = {
  title: "",
  subject: "",
  grade: "12",
  durationMinutes: "",
  passMarks: "",
  passPercentage: "60",
  timeLimitMinutes: "",
};

/**
 * Prefills the header form from whatever the file supplied, falling back to the
 * file name for a title so an untitled bank file still arrives with something
 * recognizable rather than "(untitled)".
 */
function seedHeader(result: ParseResult, fileName: string): HeaderDraft {
  const doc = result.document;
  const fallbackTitle = fileName.replace(/\.(html?|json)$/i, "").trim();

  if (doc.kind === "exam") {
    const e = doc.exam!;
    return {
      ...emptyHeader,
      title: e.title || fallbackTitle,
      subject: e.subject || "",
      grade: e.grade ? String(e.grade) : "12",
      durationMinutes: e.duration_minutes ? String(e.duration_minutes) : "",
      passMarks: e.pass_marks ? String(e.pass_marks) : "",
    };
  }
  const q = doc.quiz!;
  return {
    ...emptyHeader,
    title: q.title || fallbackTitle,
    passPercentage: q.pass_percentage != null ? String(q.pass_percentage) : "60",
    timeLimitMinutes: q.time_limit_seconds ? String(Math.round(q.time_limit_seconds / 60)) : "",
  };
}

export function ImportAssessmentDialog({
  kind,
  courseId,
}: {
  /** Restricts which document kind is accepted, so an exam file cannot be imported from the Quizzes page. */
  kind: "quiz" | "exam";
  /** When set, a document without its own course_id is imported under this course. */
  courseId?: string;
}) {
  const [open, setOpen] = useState(false);
  const [parsed, setParsed] = useState<ParseResult | null>(null);
  const [error, setError] = useState<{ message: string; issues: string[] } | null>(null);
  const [fileName, setFileName] = useState("");
  const [importing, setImporting] = useState(false);
  // Header values the operator can confirm or correct before importing. A
  // question-bank file carries a title at best, so these are editable for every
  // upload rather than only when something is missing.
  const [header, setHeader] = useState<HeaderDraft>(emptyHeader);
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  function reset() {
    setParsed(null);
    setError(null);
    setFileName("");
    setHeader(emptyHeader);
    setImporting(false);
    if (inputRef.current) inputRef.current.value = "";
  }

  async function handleFile(file: File) {
    setError(null);
    setParsed(null);
    setFileName(file.name);
    try {
      const text = await file.text();
      const result = parseImportFile(file.name, text, kind);
      if (result.document.kind !== kind) {
        setError({
          message: `This file describes a ${result.document.kind}, but this page imports ${kind}s.`,
          issues: [`Change "kind" to "${kind}", or use the ${result.document.kind} page.`],
        });
        return;
      }
      setParsed(result);
      setHeader(seedHeader(result, file.name));
    } catch (err) {
      if (err instanceof ParseError) {
        setError({ message: err.message, issues: err.issues });
        return;
      }
      setError({ message: err instanceof Error ? err.message : "Could not read this file.", issues: [] });
    }
  }

  async function handleImport() {
    if (!parsed) return;
    setImporting(true);

    // The header the operator confirmed wins over whatever the file carried.
    const doc = { ...parsed.document };
    if (doc.kind === "exam" && doc.exam) {
      doc.exam = {
        ...doc.exam,
        title: header.title.trim(),
        subject: header.subject.trim(),
        grade: Number(header.grade) || 12,
        duration_minutes: Number(header.durationMinutes) || 0,
        pass_marks: Number(header.passMarks) || 0,
      };
    } else if (doc.kind === "quiz" && doc.quiz) {
      const minutes = Number(header.timeLimitMinutes) || 0;
      doc.quiz = {
        ...doc.quiz,
        title: header.title.trim(),
        pass_percentage: Number(header.passPercentage) || 0,
        time_limit_seconds: minutes > 0 ? minutes * 60 : null,
      };
    }

    // A document with no course of its own is filed under the course being
    // viewed, if any. An explicit course_id in the file always wins.
    if (courseId) {
      if (doc.kind === "quiz" && doc.quiz && !doc.quiz.course_id) {
        doc.quiz = { ...doc.quiz, course_id: courseId };
      } else if (doc.kind === "exam" && doc.exam && !doc.exam.course_id) {
        doc.exam = { ...doc.exam, course_id: courseId };
      }
    }

    try {
      const result = await clientAction<ImportResult>("/api/teacher/assessments/import", doc);
      toast.success(
        `Imported "${result.title}" with ${result.question_count} question${result.question_count === 1 ? "" : "s"}.`
      );
      setOpen(false);
      reset();
      router.refresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Import failed.");
    } finally {
      setImporting(false);
    }
  }

  const questions = parsed?.document.questions ?? [];
  const totalMarks =
    parsed?.document.kind === "exam"
      ? questions.reduce((sum, q) => sum + (q.marks ?? q.points ?? 1), 0)
      : undefined;

  // Blocking problems with the header the operator has entered. Import stays
  // disabled while any of these hold, so a required field is never defaulted
  // behind their back on a graded assessment.
  const headerProblems: string[] = [];
  if (parsed) {
    if (!header.title.trim()) headerProblems.push("A title is required.");
    if (kind === "exam") {
      if (!header.subject.trim()) headerProblems.push("A subject is required.");
      const grade = Number(header.grade);
      if (!Number.isInteger(grade) || grade < 1 || grade > 13) {
        headerProblems.push("Grade must be between 1 and 13.");
      }
      if (!(Number(header.durationMinutes) > 0)) {
        headerProblems.push("A duration in minutes is required.");
      }
      if (totalMarks !== undefined && Number(header.passMarks) > totalMarks) {
        headerProblems.push(`Pass marks cannot exceed the ${totalMarks} marks on offer.`);
      }
    }
  }

  return (
    <>
      <Button
        size="sm"
        variant="outline"
        onClick={() => {
          reset();
          setOpen(true);
        }}
      >
        <Upload className="h-4 w-4 mr-1.5" /> Import from file
      </Button>

      <Dialog
        open={open}
        onOpenChange={(o) => {
          if (importing) return;
          setOpen(o);
          if (!o) reset();
        }}
      >
        <DialogContent className="max-w-3xl max-h-[88vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Import {kind === "exam" ? "an exam" : "a quiz"} from a file</DialogTitle>
            <DialogDescription className="text-xs">
              Upload a <code>.json</code> or <code>.html</code> file. Every question and its correct
              answer is shown below for review — nothing is saved until you confirm. Start from a
              template:{" "}
              <a
                href={`/templates/${kind}-import-template.json`}
                target="_blank"
                rel="noreferrer"
                className="underline hover:text-foreground"
              >
                {kind}-import-template.json
              </a>
              {kind === "quiz" && (
                <>
                  {" · "}
                  <a
                    href="/templates/quiz-import-template.html"
                    target="_blank"
                    rel="noreferrer"
                    className="underline hover:text-foreground"
                  >
                    quiz-import-template.html
                  </a>
                </>
              )}
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col gap-4 py-2">
            <div className="flex items-center gap-3">
              <input
                ref={inputRef}
                type="file"
                accept=".json,.html,.htm,application/json,text/html"
                className="hidden"
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) void handleFile(file);
                }}
              />
              <Button type="button" variant="secondary" size="sm" onClick={() => inputRef.current?.click()}>
                Choose file
              </Button>
              {fileName && (
                <span className="flex items-center gap-1.5 text-xs text-muted-foreground">
                  {fileName.toLowerCase().endsWith(".json") ? (
                    <FileJson className="h-3.5 w-3.5" />
                  ) : (
                    <FileCode className="h-3.5 w-3.5" />
                  )}
                  {fileName}
                </span>
              )}
            </div>

            {error && (
              <div className="rounded-md border border-destructive/40 bg-destructive/5 p-3">
                <p className="flex items-start gap-2 text-sm font-medium text-destructive">
                  <X className="h-4 w-4 mt-0.5 shrink-0" />
                  {error.message}
                </p>
                {error.issues.length > 0 && (
                  <ul className="mt-2 ml-6 list-disc text-xs text-destructive/90 space-y-0.5">
                    {error.issues.slice(0, 12).map((issue, i) => (
                      <li key={i}>{issue}</li>
                    ))}
                    {error.issues.length > 12 && <li>…and {error.issues.length - 12} more.</li>}
                  </ul>
                )}
              </div>
            )}

            {parsed && (
              <>
                <div className="rounded-md border bg-muted/30 p-3 text-sm">
                  <div className="flex flex-wrap items-center gap-x-4 gap-y-1">
                    <span className="font-medium">{header.title || "(untitled)"}</span>
                    <Badge variant="secondary" className="text-xs">
                      {questions.length} question{questions.length === 1 ? "" : "s"}
                    </Badge>
                    {totalMarks !== undefined && (
                      <Badge variant="secondary" className="text-xs">
                        {totalMarks} marks total
                      </Badge>
                    )}
                    <Badge variant="outline" className="text-xs">
                      {parsed.dialect === "canonical"
                        ? "import format"
                        : parsed.dialect === "markup"
                          ? "read from markup"
                          : `question bank (${parsed.dialect})`}
                    </Badge>
                  </div>
                </div>

                <div className="flex flex-col gap-3 rounded-md border p-3">
                  <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                    {kind === "exam" ? "Exam details" : "Quiz details"}
                    {parsed.missingHeader && parsed.missingHeader.length > 0 && (
                      <span className="ml-2 font-normal normal-case tracking-normal text-amber-700 dark:text-amber-400">
                        — this file does not specify{" "}
                        {parsed.missingHeader.join(", ").replace(/_/g, " ")}
                      </span>
                    )}
                  </p>

                  <div className="grid gap-1.5">
                    <Label htmlFor="imp-title" className="text-xs">Title</Label>
                    <Input
                      id="imp-title"
                      value={header.title}
                      onChange={(e) => setHeader((h) => ({ ...h, title: e.target.value }))}
                    />
                  </div>

                  {kind === "exam" ? (
                    <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
                      <div className="grid gap-1.5">
                        <Label htmlFor="imp-subject" className="text-xs">Subject</Label>
                        <Input
                          id="imp-subject"
                          placeholder="e.g. Biology"
                          value={header.subject}
                          onChange={(e) => setHeader((h) => ({ ...h, subject: e.target.value }))}
                        />
                      </div>
                      <div className="grid gap-1.5">
                        <Label htmlFor="imp-grade" className="text-xs">Grade</Label>
                        <Input
                          id="imp-grade"
                          type="number"
                          min="1"
                          max="13"
                          value={header.grade}
                          onChange={(e) => setHeader((h) => ({ ...h, grade: e.target.value }))}
                        />
                      </div>
                      <div className="grid gap-1.5">
                        <Label htmlFor="imp-dur" className="text-xs">Duration (min)</Label>
                        <Input
                          id="imp-dur"
                          type="number"
                          min="1"
                          placeholder="e.g. 150"
                          value={header.durationMinutes}
                          onChange={(e) =>
                            setHeader((h) => ({ ...h, durationMinutes: e.target.value }))
                          }
                        />
                      </div>
                      <div className="grid gap-1.5">
                        <Label htmlFor="imp-pass" className="text-xs">
                          Pass marks{totalMarks !== undefined ? ` / ${totalMarks}` : ""}
                        </Label>
                        <Input
                          id="imp-pass"
                          type="number"
                          min="0"
                          value={header.passMarks}
                          onChange={(e) => setHeader((h) => ({ ...h, passMarks: e.target.value }))}
                        />
                      </div>
                    </div>
                  ) : (
                    <div className="grid grid-cols-2 gap-3">
                      <div className="grid gap-1.5">
                        <Label htmlFor="imp-pct" className="text-xs">Pass %</Label>
                        <Input
                          id="imp-pct"
                          type="number"
                          min="0"
                          max="100"
                          value={header.passPercentage}
                          onChange={(e) =>
                            setHeader((h) => ({ ...h, passPercentage: e.target.value }))
                          }
                        />
                      </div>
                      <div className="grid gap-1.5">
                        <Label htmlFor="imp-limit" className="text-xs">
                          Time limit (min, blank = none)
                        </Label>
                        <Input
                          id="imp-limit"
                          type="number"
                          min="0"
                          value={header.timeLimitMinutes}
                          onChange={(e) =>
                            setHeader((h) => ({ ...h, timeLimitMinutes: e.target.value }))
                          }
                        />
                      </div>
                    </div>
                  )}

                  {headerProblems.length > 0 && (
                    <ul className="list-disc pl-5 text-xs text-destructive">
                      {headerProblems.map((p, i) => (
                        <li key={i}>{p}</li>
                      ))}
                    </ul>
                  )}
                </div>

                {parsed.warnings.length > 0 && (
                  <div className="rounded-md border border-amber-500/40 bg-amber-500/5 p-3">
                    <p className="flex items-center gap-2 text-sm font-medium text-amber-700 dark:text-amber-400">
                      <AlertTriangle className="h-4 w-4 shrink-0" />
                      Review before importing
                    </p>
                    <ul className="mt-2 ml-6 list-disc text-xs text-amber-700/90 dark:text-amber-400/90 space-y-0.5">
                      {parsed.warnings.slice(0, 10).map((w, i) => (
                        <li key={i}>{w}</li>
                      ))}
                      {parsed.warnings.length > 10 && (
                        <li>…and {parsed.warnings.length - 10} more.</li>
                      )}
                    </ul>
                  </div>
                )}

                <div className="flex flex-col gap-3">
                  <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                    Questions — the ticked option is what will be graded as correct
                  </p>
                  {questions.map((q, i) => {
                    const answers = resolvePreviewAnswers(q);
                    return (
                      <div key={i} className="rounded-md border p-3">
                        <div className="flex items-start justify-between gap-3">
                          <p className="text-sm font-medium">
                            {i + 1}. {q.text}
                          </p>
                          <Badge variant="outline" className="text-xs shrink-0">
                            {q.marks ?? q.points ?? 1} {kind === "exam" ? "marks" : "pts"}
                          </Badge>
                        </div>
                        <ul className="mt-2 space-y-1">
                          {answers.map((a, j) => (
                            <li
                              key={j}
                              className={
                                a.isCorrect
                                  ? "flex items-center gap-2 text-sm font-medium text-emerald-700 dark:text-emerald-400"
                                  : "flex items-center gap-2 text-sm text-muted-foreground"
                              }
                            >
                              {a.isCorrect ? (
                                <Check className="h-3.5 w-3.5 shrink-0" />
                              ) : (
                                <span className="h-3.5 w-3.5 shrink-0" />
                              )}
                              <span>{String.fromCharCode(65 + j)}. {a.text}</span>
                            </li>
                          ))}
                        </ul>
                        {q.explanation && (
                          <p className="mt-2 text-xs text-muted-foreground border-t pt-2">
                            <span className="font-medium">Explanation: </span>
                            {q.explanation}
                          </p>
                        )}
                      </div>
                    );
                  })}
                </div>
              </>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              disabled={importing}
              onClick={() => {
                setOpen(false);
                reset();
              }}
            >
              Cancel
            </Button>
            <Button
              type="button"
              onClick={handleImport}
              disabled={!parsed || importing || headerProblems.length > 0}
            >
              {importing ? (
                <>
                  <Loader2 className="h-4 w-4 mr-1.5 animate-spin" />
                  Importing…
                </>
              ) : (
                `Import ${questions.length} question${questions.length === 1 ? "" : "s"}`
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
