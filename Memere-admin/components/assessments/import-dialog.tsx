"use client";

import { useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Upload, Loader2, AlertTriangle, Check, X, FileJson, FileCode } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
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
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  function reset() {
    setParsed(null);
    setError(null);
    setFileName("");
    setImporting(false);
    if (inputRef.current) inputRef.current.value = "";
  }

  async function handleFile(file: File) {
    setError(null);
    setParsed(null);
    setFileName(file.name);
    try {
      const text = await file.text();
      const result = parseImportFile(file.name, text);
      if (result.document.kind !== kind) {
        setError({
          message: `This file describes a ${result.document.kind}, but this page imports ${kind}s.`,
          issues: [`Change "kind" to "${kind}", or use the ${result.document.kind} page.`],
        });
        return;
      }
      setParsed(result);
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

    // A document with no course of its own is filed under the course being
    // viewed, if any. An explicit course_id in the file always wins.
    const doc = { ...parsed.document };
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
  const header =
    parsed?.document.kind === "exam" ? parsed.document.exam : parsed?.document.quiz;
  const totalMarks =
    parsed?.document.kind === "exam"
      ? questions.reduce((sum, q) => sum + (q.marks ?? q.points ?? 1), 0)
      : undefined;

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
                    <span className="font-medium">{header?.title || "(untitled)"}</span>
                    <Badge variant="secondary" className="text-xs">
                      {questions.length} question{questions.length === 1 ? "" : "s"}
                    </Badge>
                    {totalMarks !== undefined && (
                      <Badge variant="secondary" className="text-xs">
                        {totalMarks} marks total
                      </Badge>
                    )}
                    <Badge variant="outline" className="text-xs">
                      parsed from {parsed.source.toUpperCase()}
                    </Badge>
                  </div>
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
            <Button type="button" onClick={handleImport} disabled={!parsed || importing}>
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
