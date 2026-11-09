import {
  ImportDocumentSchema,
  IMPORT_FORMAT_VERSION,
  type ImportDocument,
  type ImportQuestion,
} from "@/lib/api/schemas";

/**
 * Parsing of uploaded assessment files into the canonical import document.
 *
 * Two input formats are supported:
 *
 *   .json — the canonical format, validated against ImportDocumentSchema.
 *   .html — hand-written markup, read structurally by parseHtml below.
 *
 * HTML parsing is inherently lossy: markup that looks right to a human can still
 * yield a different answer key than intended. Every parse therefore returns
 * `warnings` alongside the document, and the import UI requires the operator to
 * review the resolved key for every question before anything is written. Nothing
 * here writes data — a parse result is a proposal, not an import.
 */

export interface ParseResult {
  document: ImportDocument;
  /** Non-fatal observations the operator should read before confirming. */
  warnings: string[];
  /** Which input format produced the document. */
  source: "json" | "html";
}

export class ParseError extends Error {
  constructor(
    message: string,
    /** Per-field problems, when the failure came from schema validation. */
    public issues: string[] = []
  ) {
    super(message);
    this.name = "ParseError";
  }
}

/** Dispatches on file extension, falling back to sniffing the content. */
export function parseImportFile(fileName: string, text: string): ParseResult {
  const lower = fileName.toLowerCase();
  if (lower.endsWith(".json")) return parseJson(text);
  if (lower.endsWith(".html") || lower.endsWith(".htm")) return parseHtml(text);

  // Unknown extension: a leading "{" means JSON, anything else is treated as
  // markup. Guessing beats rejecting a correctly-formatted file over its name.
  return text.trimStart().startsWith("{") ? parseJson(text) : parseHtml(text);
}

// ---- JSON --------------------------------------------------------------------

export function parseJson(text: string): ParseResult {
  let raw: unknown;
  try {
    raw = JSON.parse(text);
  } catch (err) {
    throw new ParseError(
      `This file is not valid JSON: ${err instanceof Error ? err.message : "parse failed"}`
    );
  }

  const parsed = ImportDocumentSchema.safeParse(raw);
  if (!parsed.success) {
    throw new ParseError(
      "The JSON does not match the import format.",
      parsed.error.issues.map((i) => `${i.path.join(".") || "(root)"}: ${i.message}`)
    );
  }

  return { document: parsed.data, warnings: collectWarnings(parsed.data), source: "json" };
}

// ---- HTML --------------------------------------------------------------------

/**
 * Reads a question set out of HTML markup.
 *
 * Recognized structure — each question is an element carrying
 * `data-question` (or a `.question` class), containing:
 *
 *   - the prompt: `[data-text]`, `.question-text`, or the first heading/`<p>`
 *   - the options: `<li>` elements, or `[data-option]` elements
 *   - the correct option(s): marked with `data-correct`, a `correct` class,
 *     `<strong>`/`<b>` wrapping, or a leading "*"
 *   - optionally `[data-explanation]` / `.explanation`
 *   - optionally `data-marks` / `data-points`, `data-subject`, `data-topic`
 *
 * A document that embeds a canonical JSON block in
 * `<script type="application/json">` is read from that block instead, which is
 * exact and warning-free.
 */
export function parseHtml(text: string): ParseResult {
  if (typeof DOMParser === "undefined") {
    throw new ParseError("HTML files can only be parsed in the browser.");
  }

  const doc = new DOMParser().parseFromString(text, "text/html");

  // An embedded canonical block is unambiguous, so prefer it over scraping.
  const embedded = doc.querySelector('script[type="application/json"]');
  if (embedded?.textContent?.trim()) {
    const result = parseJson(embedded.textContent);
    return { ...result, source: "html" };
  }

  const warnings: string[] = [
    "This file was read from HTML markup, which cannot be interpreted with certainty. Check the correct answer on every question below before importing.",
  ];

  const nodes = Array.from(
    doc.querySelectorAll<HTMLElement>("[data-question], .question")
  );
  if (nodes.length === 0) {
    throw new ParseError(
      "No questions found in this HTML. Each question needs a data-question attribute or a question class."
    );
  }

  const questions: ImportQuestion[] = [];
  nodes.forEach((node, i) => {
    const number = i + 1;
    const prompt = readPrompt(node);
    if (!prompt) {
      throw new ParseError(`Question ${number} has no readable prompt text.`);
    }

    const optionEls = readOptionElements(node);
    if (optionEls.length < 2) {
      throw new ParseError(
        `Question ${number} ("${truncate(prompt)}") has ${optionEls.length} option(s); at least 2 are required.`
      );
    }

    const answers = optionEls.map((el) => {
      const marked = isMarkedCorrect(el);
      return { text: stripCorrectMarker(textOf(el)), is_correct: marked };
    });

    const correctCount = answers.filter((a) => a.is_correct).length;
    if (correctCount === 0) {
      throw new ParseError(
        `Question ${number} ("${truncate(prompt)}") has no option marked correct. Mark it with data-correct, a correct class, <strong>, or a leading "*".`
      );
    }
    if (correctCount > 1) {
      // Legitimate for multi-select, but far more often a stray <strong>, so it
      // is surfaced rather than silently accepted.
      warnings.push(
        `Question ${number} has ${correctCount} options marked correct. Confirm this question really accepts multiple answers.`
      );
    }

    const marks = readNumber(node, ["marks", "points"]);
    const explanation = readOptional(node, "[data-explanation], .explanation", "explanation");
    const type = answers.length === 2 && looksBoolean(answers) ? "true_false" : "multiple_choice";

    questions.push({
      text: prompt,
      type,
      ...(marks !== undefined ? { points: marks, marks } : {}),
      ...(explanation ? { explanation } : {}),
      ...(node.dataset.subject ? { subject: node.dataset.subject } : {}),
      ...(node.dataset.topic ? { topic: node.dataset.topic } : {}),
      answers,
    });
  });

  // The header comes from data attributes on <body> or a [data-assessment]
  // wrapper; anything missing is filled in by the operator in the dialog.
  const header = doc.querySelector<HTMLElement>("[data-assessment]") ?? doc.body;
  const kind = header?.dataset.kind === "exam" ? "exam" : "quiz";
  const title =
    header?.dataset.title ??
    doc.querySelector("title")?.textContent?.trim() ??
    doc.querySelector("h1")?.textContent?.trim() ??
    "";

  const document: ImportDocument =
    kind === "exam"
      ? {
          version: IMPORT_FORMAT_VERSION,
          kind: "exam",
          exam: {
            title,
            subject: header?.dataset.subject ?? "",
            grade: Number(header?.dataset.grade) || 12,
            duration_minutes: Number(header?.dataset.durationMinutes) || 60,
            pass_marks: Number(header?.dataset.passMarks) || 0,
          },
          questions,
        }
      : {
          version: IMPORT_FORMAT_VERSION,
          kind: "quiz",
          quiz: {
            title,
            pass_percentage: Number(header?.dataset.passPercentage) || 60,
            ...(header?.dataset.timeLimitSeconds
              ? { time_limit_seconds: Number(header.dataset.timeLimitSeconds) }
              : {}),
          },
          questions,
        };

  return { document, warnings: [...warnings, ...collectWarnings(document)], source: "html" };
}

// ---- helpers -----------------------------------------------------------------

function textOf(el: Element): string {
  return (el.textContent ?? "").replace(/\s+/g, " ").trim();
}

function truncate(s: string, max = 48): string {
  return s.length > max ? `${s.slice(0, max)}…` : s;
}

function readPrompt(node: HTMLElement): string {
  const explicit = node.querySelector("[data-text], .question-text");
  if (explicit) return textOf(explicit);
  if (node.dataset.text) return node.dataset.text.trim();

  const heading = node.querySelector("h1, h2, h3, h4, h5, h6, p");
  if (heading) return textOf(heading);

  // Last resort: the node's own text minus its options, so the prompt does not
  // swallow the answer list.
  const clone = node.cloneNode(true) as HTMLElement;
  clone.querySelectorAll("li, [data-option], [data-explanation], .explanation").forEach((n) => n.remove());
  return textOf(clone);
}

function readOptionElements(node: HTMLElement): HTMLElement[] {
  const tagged = Array.from(node.querySelectorAll<HTMLElement>("[data-option]"));
  if (tagged.length > 0) return tagged;
  return Array.from(node.querySelectorAll<HTMLElement>("li"));
}

function isMarkedCorrect(el: HTMLElement): boolean {
  if (el.dataset.correct !== undefined && el.dataset.correct !== "false") return true;
  if (el.classList.contains("correct")) return true;
  if (el.querySelector("strong, b")) return true;
  if (textOf(el).startsWith("*")) return true;
  return false;
}

function stripCorrectMarker(text: string): string {
  return text.replace(/^\*\s*/, "").trim();
}

function readOptional(node: HTMLElement, selector: string, dataKey: string): string | undefined {
  const el = node.querySelector(selector);
  if (el) {
    const t = textOf(el);
    if (t) return t;
  }
  const raw = node.dataset[dataKey];
  return raw?.trim() || undefined;
}

function readNumber(node: HTMLElement, keys: string[]): number | undefined {
  for (const key of keys) {
    const raw = node.dataset[key];
    if (raw !== undefined) {
      const n = Number(raw);
      if (Number.isFinite(n) && n > 0) return Math.trunc(n);
    }
  }
  return undefined;
}

function looksBoolean(answers: { text: string }[]): boolean {
  const set = answers.map((a) => a.text.toLowerCase().trim()).sort();
  return set.length === 2 && set[0] === "false" && set[1] === "true";
}

/**
 * Observations worth showing the operator for any document, however it was
 * parsed. These never block an import — they flag the shapes that are legal but
 * usually a mistake.
 */
function collectWarnings(doc: ImportDocument): string[] {
  const warnings: string[] = [];

  doc.questions.forEach((q, i) => {
    const number = i + 1;
    const optionCount = q.answers?.length ?? q.options?.length ?? 0;
    if (optionCount < 2) {
      warnings.push(`Question ${number} has only ${optionCount} option(s).`);
    }
    if (q.answers) {
      const correct = q.answers.filter((a) => a.is_correct).length;
      if (correct > 1) {
        warnings.push(
          `Question ${number} marks ${correct} options correct — confirm it accepts multiple answers.`
        );
      }
    }
    if (!q.explanation) {
      warnings.push(`Question ${number} has no explanation, so students see no feedback after grading.`);
    }
  });

  if (doc.kind === "exam" && doc.exam) {
    const total = doc.questions.reduce((sum, q) => sum + (q.marks ?? q.points ?? 1), 0);
    const pass = doc.exam.pass_marks ?? 0;
    if (pass > total) {
      warnings.push(
        `Pass marks (${pass}) exceed the total this question set produces (${total}); the import will be rejected.`
      );
    }
  }

  return warnings;
}

/** Resolves a question's options to display text + correctness for the preview. */
export function resolvePreviewAnswers(
  q: ImportQuestion
): { text: string; isCorrect: boolean }[] {
  if (q.answers?.length) {
    return q.answers.map((a) => ({ text: a.text, isCorrect: a.is_correct }));
  }
  const options = q.options ?? [];
  const selectors = Array.isArray(q.correct) ? q.correct : [q.correct];
  const correct = new Set<number>();

  for (const sel of selectors) {
    if (typeof sel === "number" && Number.isInteger(sel) && sel >= 0 && sel < options.length) {
      correct.add(sel);
      continue;
    }
    if (typeof sel === "string") {
      const needle = sel.trim().toLowerCase();
      // Mirror the backend's resolution order exactly: option text first, then
      // the letter label, so the preview cannot disagree with what is stored.
      const byText = options.findIndex((o) => o.trim().toLowerCase() === needle);
      if (byText >= 0) {
        correct.add(byText);
        continue;
      }
      if (needle.length === 1) {
        const idx = needle.toUpperCase().charCodeAt(0) - 65;
        if (idx >= 0 && idx < options.length) correct.add(idx);
      }
    }
  }

  return options.map((text, i) => ({ text, isCorrect: correct.has(i) }));
}
