import { IMPORT_FORMAT_VERSION, type ImportDocument, type ImportQuestion } from "@/lib/api/schemas";

/**
 * Readers for the "question bank" HTML files the content team produces.
 *
 * Those files are self-contained exam viewers: the questions live in a machine-
 * readable block inside the page, but under their own field names
 * (`questionText`, `correctAnswer`) rather than this system's. Four shapes are in
 * circulation, and all four carry the same question fields:
 *
 *   1. `<script type="application/json">{ "metadata": {…}, "questions": [ … ] }`
 *   2. `<script type="application/json">[ … ]`            — a bare array
 *   3. `window.__EXAM_DATA__  = [ … ]`                    — a JS array literal
 *   4. `window.__EXAM_STORE__ = { subject: [ … ] }`       — keyed JS literal
 *
 * Shapes 3 and 4 are JavaScript, not JSON: the keys are unquoted. They are read
 * with a scanner (see relaxedJsonParse) rather than `eval`, so nothing in an
 * uploaded file is ever executed.
 *
 * Normalizing happens here and nowhere else, so the preview, the browser-side
 * validation and the server all see one canonical shape.
 */

/** One question as the bank files write it. Every field is treated as untrusted. */
interface BankQuestion {
  id?: unknown;
  points?: unknown;
  marks?: unknown;
  sourceRef?: unknown;
  questionText?: unknown;
  text?: unknown;
  options?: unknown;
  correctAnswer?: unknown;
  explanation?: unknown;
  subject?: unknown;
  topic?: unknown;
}

export interface DialectResult {
  /** Normalized questions, ready for the canonical document. */
  questions: ImportQuestion[];
  /** Header values the file supplied; anything absent is left undefined. */
  header: { title?: string; subject?: string; durationMinutes?: number };
  warnings: string[];
  /** Which bank shape matched, for the preview to report. */
  dialect: string;
}

/**
 * Reads a question bank out of an HTML document, or returns null when the
 * document carries none of the known shapes.
 */
export function readQuestionBank(doc: Document, rawHtml: string): DialectResult | null {
  const found = locateBank(doc, rawHtml);
  if (!found) return null;

  const { rows, dialect, metadata } = found;
  const warnings: string[] = [];

  if (found.repairs) {
    warnings.push(
      `This file's embedded JSON had ${found.repairs} invalid escape sequence(s) — LaTeX commands such as \\{ written with a single backslash. They were read as literal backslashes. The source file is malformed and worth fixing at the generator.`
    );
  }
  const questions: ImportQuestion[] = [];

  // A single-backslash LaTeX command is the quietest way a bank file can go
  // wrong. `\times` is a VALID JSON escape (tab + "imes"), so JSON.parse accepts
  // it without complaint and the maths is destroyed with no error anywhere. The
  // fingerprint is a control character in the raw value — legitimate content
  // never contains one — and it has to be checked BEFORE htmlToText, which
  // collapses tabs into spaces and would hide the evidence.
  const controlChar = /[\t\f\v\b\r]/;
  let mangled = 0;

  rows.forEach((row, i) => {
    const number = i + 1;
    const q = row as BankQuestion;

    const rawText = firstString(q.questionText, q.text);
    const rawOptions = Array.isArray(q.options) ? q.options.map((o) => String(o ?? "")) : [];
    if (
      controlChar.test(rawText) ||
      controlChar.test(asString(q.explanation) ?? "") ||
      rawOptions.some((o) => controlChar.test(o))
    ) {
      mangled++;
    }

    const text = htmlToText(rawText);
    if (!text) {
      throw new BankFormatError(`Question ${number} has no question text.`);
    }

    if (!Array.isArray(q.options) || q.options.length < 2) {
      throw new BankFormatError(
        `Question ${number} ("${clip(text)}") has ${
          Array.isArray(q.options) ? q.options.length : 0
        } option(s); at least 2 are required.`
      );
    }
    const options = q.options.map((o) => htmlToText(String(o ?? "")));
    const blank = options.findIndex((o) => !o);
    if (blank >= 0) {
      throw new BankFormatError(
        `Question ${number} ("${clip(text)}") has an empty option at position ${blank + 1}.`
      );
    }

    // correctAnswer is a 0-based index into options. It is the answer key, so an
    // out-of-range or non-numeric value aborts rather than defaulting to 0 —
    // defaulting would silently mark the first option correct for every bad row.
    const correct = q.correctAnswer;
    if (typeof correct !== "number" || !Number.isInteger(correct)) {
      throw new BankFormatError(
        `Question ${number} ("${clip(text)}") has correctAnswer ${JSON.stringify(
          correct
        )}, which is not a whole number.`
      );
    }
    if (correct < 0 || correct >= options.length) {
      throw new BankFormatError(
        `Question ${number} ("${clip(text)}") has correctAnswer ${correct}, out of range for ${
          options.length
        } options.`
      );
    }

    const rawExplanation = asString(q.explanation);
    const topic = categoryLabel(rawExplanation) ?? asString(q.topic) ?? undefined;

    // sourceRef is the provenance of the question ("Q129 [2012/13 E.C. | Booklet
    // 46]"). There is no field for it, and it is genuinely useful to a student
    // reviewing an answer, so it is kept as a trailing line of the explanation
    // rather than dropped.
    const sourceRef = asString(q.sourceRef);
    let explanation = htmlToText(rawExplanation ?? "");
    if (sourceRef) {
      explanation = explanation ? `${explanation}\n\nSource: ${sourceRef}` : `Source: ${sourceRef}`;
    }

    const points = positiveInt(q.marks) ?? positiveInt(q.points) ?? 1;

    questions.push({
      text,
      type: looksTrueFalse(options) ? "true_false" : "multiple_choice",
      points,
      marks: points,
      ...(explanation ? { explanation } : {}),
      ...(asString(q.subject) ? { subject: asString(q.subject)! } : {}),
      ...(topic ? { topic } : {}),
      options,
      correct,
    });
  });

  if (questions.length === 0) {
    throw new BankFormatError("The question bank in this file is empty.");
  }

  // Bank explanations often restate the key in prose ("Correct Choice: B"). That
  // is an independent assertion of the same fact as correctAnswer, so comparing
  // them catches a generator that mismatched the two. When they disagree one of
  // them is wrong and importing either is a coin flip on a graded question, so
  // the operator is told which questions to check rather than it passing
  // silently.
  const disagreements: string[] = [];
  questions.forEach((q, i) => {
    const stated = /correct\s*(?:choice|answer|option)\s*[:\-]?\s*\(?([A-D])\)?/i.exec(
      q.explanation ?? ""
    );
    if (!stated) return;
    const wantIndex = stated[1].toUpperCase().charCodeAt(0) - 65;
    const keyIndex = (q.options ?? []).findIndex((_, j) => j === q.correct);
    if (typeof q.correct === "number" && keyIndex >= 0 && wantIndex !== q.correct) {
      disagreements.push(
        `Question ${i + 1}: the answer key marks ${String.fromCharCode(
          65 + q.correct
        )} but its explanation says ${stated[1].toUpperCase()}.`
      );
    }
  });
  if (disagreements.length > 0) {
    warnings.push(
      `${disagreements.length} question(s) have an answer key that contradicts their own explanation — check these before importing: ${disagreements
        .slice(0, 5)
        .join(" ")}${disagreements.length > 5 ? " …" : ""}`
    );
  }

  if (mangled > 0) {
    warnings.push(
      `${mangled} question(s) contain stray control characters, which usually means the source file wrote a LaTeX command with a single backslash (\\times instead of \\\\times). Those commands have already been lost and cannot be recovered here — fix the generator and re-export rather than importing this file.`
    );
  }

  // Content fidelity: these files are built for a MathJax-rendering viewer. The
  // LaTeX source is preserved verbatim, but whether a student sees rendered maths
  // depends on the app, so it is called out rather than passed over.
  const withMath = questions.filter(
    (q) => hasMath(q.text) || (q.options ?? []).some(hasMath) || hasMath(q.explanation ?? "")
  ).length;
  if (withMath > 0) {
    warnings.push(
      `${withMath} question(s) contain LaTeX maths ($…$). The notation is imported exactly as written, but it only displays as maths if the student app renders LaTeX.`
    );
  }

  const missingExplanations = questions.filter((q) => !q.explanation).length;
  if (missingExplanations > 0) {
    warnings.push(
      `${missingExplanations} question(s) have no explanation, so students get no feedback after grading.`
    );
  }

  return {
    questions,
    header: {
      title: documentTitle(doc),
      subject: metadata.subject,
      durationMinutes: metadata.durationMinutes,
    },
    warnings,
    dialect,
  };
}

/** Raised when a file matches a bank shape but its contents cannot be trusted. */
export class BankFormatError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "BankFormatError";
  }
}

// ---- locating the bank -------------------------------------------------------

interface Located {
  rows: unknown[];
  dialect: string;
  metadata: { subject?: string; durationMinutes?: number };
  /** Count of invalid JSON escapes repaired while reading, for reporting. */
  repairs?: number;
}

function locateBank(doc: Document, rawHtml: string): Located | null {
  const metaSubject =
    doc.querySelector('meta[name="exam:subject"]')?.getAttribute("content")?.trim() || undefined;

  // Shapes 1 and 2 — a JSON block. Several files ship an EMPTY block alongside a
  // populated JS literal, so an empty one is skipped rather than accepted as an
  // answer; otherwise those files would import zero questions and look fine.
  for (const script of Array.from(doc.querySelectorAll('script[type="application/json"]'))) {
    const body = script.textContent?.trim();
    if (!body || body === "[]" || body === "{}") continue;
    let parsed: unknown;
    let repairs = 0;
    try {
      const read = parseJsonLenient(body);
      parsed = read.value;
      repairs = read.repairs;
    } catch {
      continue;
    }
    if (Array.isArray(parsed) && parsed.length > 0 && isBankRow(parsed[0])) {
      return { rows: parsed, dialect: "json-array", metadata: { subject: metaSubject }, repairs };
    }
    if (parsed && typeof parsed === "object") {
      const obj = parsed as Record<string, unknown>;
      if (Array.isArray(obj.questions) && obj.questions.length > 0 && isBankRow(obj.questions[0])) {
        const md = (obj.metadata ?? {}) as Record<string, unknown>;
        return {
          rows: obj.questions,
          dialect: "json-metadata",
          repairs,
          metadata: {
            subject: asString(md.subject) ?? metaSubject,
            durationMinutes: positiveInt(md.duration_minutes),
          },
        };
      }
    }
  }

  // Shapes 3 and 4 — a JS literal assigned to a known global.
  for (const global of ["__EXAM_DATA__", "__EXAM_STORE__"]) {
    const literal = extractAssignment(rawHtml, global);
    if (!literal) continue;
    let parsed: unknown;
    try {
      parsed = relaxedJsonParse(literal);
    } catch (err) {
      throw new BankFormatError(
        `Found window.${global} in this file but could not read it: ${
          err instanceof Error ? err.message : "unknown error"
        }`
      );
    }

    if (Array.isArray(parsed) && parsed.length > 0 && isBankRow(parsed[0])) {
      return { rows: parsed, dialect: `js-${global}`, metadata: { subject: metaSubject } };
    }
    // __EXAM_STORE__ keys questions by subject. Every group is imported as one
    // assessment, and the key names the subject when the page does not.
    if (parsed && typeof parsed === "object") {
      const groups = Object.entries(parsed as Record<string, unknown>).filter(
        ([, v]) => Array.isArray(v) && v.length > 0 && isBankRow((v as unknown[])[0])
      );
      if (groups.length > 0) {
        const rows = groups.flatMap(([, v]) => v as unknown[]);
        return {
          rows,
          dialect: `js-${global}`,
          metadata: {
            subject: metaSubject ?? (groups.length === 1 ? titleCase(groups[0][0]) : undefined),
          },
        };
      }
    }
  }

  return null;
}

/** A row is a bank question when it has the text + options + key triple. */
function isBankRow(value: unknown): boolean {
  if (!value || typeof value !== "object") return false;
  const row = value as BankQuestion;
  const hasText = typeof row.questionText === "string" || typeof row.text === "string";
  return hasText && Array.isArray(row.options) && "correctAnswer" in row;
}

/**
 * Pulls the literal out of `window.<name> = <literal>;` by scanning for the
 * matching bracket, so a `}` or `]` inside a string cannot end it early.
 */
function extractAssignment(src: string, globalName: string): string | null {
  const assign = new RegExp(`window\\s*\\.\\s*${globalName}\\s*=\\s*`).exec(src);
  if (!assign) return null;

  let i = assign.index + assign[0].length;
  while (i < src.length && /\s/.test(src[i])) i++;
  const open = src[i];
  if (open !== "[" && open !== "{") return null;
  const close = open === "[" ? "]" : "}";

  let depth = 0;
  let inString: string | null = null;
  for (let j = i; j < src.length; j++) {
    const c = src[j];
    if (inString) {
      if (c === "\\") j++;
      else if (c === inString) inString = null;
      continue;
    }
    if (c === '"' || c === "'") inString = c;
    else if (c === open) depth++;
    else if (c === close) {
      depth--;
      if (depth === 0) return src.slice(i, j + 1);
    }
  }
  return null;
}

/**
 * Parses JSON, retrying once with lone backslashes escaped when the strict parse
 * fails on an invalid escape sequence.
 *
 * Question banks are full of LaTeX, and a generator that writes `$\{x\}$` into a
 * JSON string emits `\{`, which JSON rejects — the backslash should have been
 * doubled. Repairing it is safe precisely because it only ever runs on input
 * JSON.parse has already refused: a file that parses strictly is never touched,
 * so the repair cannot change the meaning of a valid document. Every repair is
 * counted and reported to the operator.
 */
export function parseJsonLenient(text: string): { value: unknown; repairs: number } {
  try {
    return { value: JSON.parse(text), repairs: 0 };
  } catch (strictError) {
    let repairs = 0;
    let out = "";
    let inString = false;

    for (let i = 0; i < text.length; i++) {
      const c = text[i];
      if (!inString) {
        if (c === '"') inString = true;
        out += c;
        continue;
      }
      if (c === "\\") {
        const next = text[i + 1] ?? "";
        if ('\\"/bfnrtu'.includes(next)) {
          out += c + next;
          i++;
        } else {
          // A LaTeX command, not an escape: keep the backslash literally.
          out += "\\\\";
          repairs++;
        }
        continue;
      }
      if (c === '"') inString = false;
      out += c;
    }

    if (repairs === 0) throw strictError;
    return { value: JSON.parse(out), repairs };
  }
}

/**
 * Parses a JavaScript object literal by rewriting it into JSON, then handing it
 * to JSON.parse. Nothing is evaluated.
 *
 * The rewrite is done with a scanner rather than a regex because both
 * transformations are only valid outside string literals: question text really
 * does contain fragments like `"Student 1: blind folded"`, and a regex that
 * quoted every `identifier:` would corrupt the content it was meant to preserve.
 */
export function relaxedJsonParse(src: string): unknown {
  let out = "";
  let inString: string | null = null;
  // The last structural character seen outside a string; an unquoted key is only
  // legal directly after `{` or `,`.
  let lastSignificant = "";

  for (let i = 0; i < src.length; i++) {
    const c = src[i];

    if (inString) {
      out += c;
      if (c === "\\") {
        out += src[++i] ?? "";
      } else if (c === inString) {
        inString = null;
      }
      continue;
    }

    if (c === '"') {
      inString = '"';
      out += c;
      lastSignificant = c;
      continue;
    }
    if (c === "'") {
      // Single-quoted string: re-emit as a double-quoted JSON string.
      let body = "";
      i++;
      for (; i < src.length && src[i] !== "'"; i++) {
        if (src[i] === "\\") {
          body += src[i] + (src[++i] ?? "");
        } else if (src[i] === '"') {
          body += '\\"';
        } else {
          body += src[i];
        }
      }
      out += `"${body}"`;
      lastSignificant = '"';
      continue;
    }

    // Line and block comments are not JSON; drop them.
    if (c === "/" && src[i + 1] === "/") {
      while (i < src.length && src[i] !== "\n") i++;
      continue;
    }
    if (c === "/" && src[i + 1] === "*") {
      i += 2;
      while (i < src.length && !(src[i] === "*" && src[i + 1] === "/")) i++;
      i++;
      continue;
    }

    if (/\s/.test(c)) {
      out += c;
      continue;
    }

    // An unquoted key: an identifier directly after `{` or `,`, followed by `:`.
    if (/[A-Za-z_$]/.test(c) && (lastSignificant === "{" || lastSignificant === ",")) {
      let j = i;
      let ident = "";
      while (j < src.length && /[\w$]/.test(src[j])) ident += src[j++];
      let k = j;
      while (k < src.length && /\s/.test(src[k])) k++;
      if (src[k] === ":") {
        out += `"${ident}"`;
        i = j - 1;
        lastSignificant = '"';
        continue;
      }
    }

    // Trailing comma before a closing bracket is legal JS, not JSON.
    if (c === "," ) {
      let k = i + 1;
      while (k < src.length && /\s/.test(src[k])) k++;
      if (src[k] === "}" || src[k] === "]") continue;
    }

    out += c;
    lastSignificant = c;
  }

  return JSON.parse(out);
}

// ---- normalizing values ------------------------------------------------------

/**
 * Converts the bank files' rich text to the plain text this system stores.
 * `<br>` and block boundaries become newlines, tags are dropped and entities
 * decoded. LaTeX delimiters are left untouched: `$x^2$` is content, not markup.
 */
export function htmlToText(input: string): string {
  if (!input) return "";
  if (typeof DOMParser === "undefined") {
    return input.replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
  }

  const withBreaks = input
    .replace(/<\s*br\s*\/?\s*>/gi, "\n")
    .replace(/<\s*\/\s*(p|div|li|tr)\s*>/gi, "\n");

  // Parsing as a document (rather than assigning innerHTML) keeps the markup
  // inert; textContent then yields decoded entities without the tags.
  const parsed = new DOMParser().parseFromString(`<body>${withBreaks}</body>`, "text/html");
  const text = parsed.body?.textContent ?? "";

  return text
    .split("\n")
    .map((line) => line.replace(/[ \t ]+/g, " ").trim())
    .join("\n")
    .replace(/\n{3,}/g, "\n\n")
    .trim();
}

/**
 * Many bank explanations open with a bold category label — "Antonyms:",
 * "Analogy (Tool & Purpose):" — which is exactly the per-question topic the
 * weak-areas analytics wants. Generic lead-ins like "Step-by-step derivation:"
 * describe the explanation's form rather than the question's topic, so they are
 * ignored instead of becoming a meaningless topic on every question.
 */
export function categoryLabel(explanationHtml: string | undefined): string | undefined {
  if (!explanationHtml) return undefined;
  const m = /^\s*<\s*(b|strong)\s*>(.*?)<\s*\/\s*\1\s*>/i.exec(explanationHtml);
  if (!m) return undefined;

  const label = htmlToText(m[2]).replace(/[:\s]+$/, "").trim();
  if (!label || label.length > 60) return undefined;
  if (/^(step[-\s]?by[-\s]?step|detailed|derivation|explanation|rationale|solution|answer)\b/i.test(label)) {
    return undefined;
  }
  return label;
}

function documentTitle(doc: Document): string | undefined {
  const title = doc.querySelector("title")?.textContent?.trim();
  if (title) return title;
  return doc.querySelector("h1")?.textContent?.trim() || undefined;
}

function looksTrueFalse(options: string[]): boolean {
  if (options.length !== 2) return false;
  const set = options.map((o) => o.toLowerCase().trim()).sort();
  return set[0] === "false" && set[1] === "true";
}

function hasMath(s: string): boolean {
  return /\$[^$]+\$|\\\(|\\\[|\\text\{|\\frac/.test(s);
}

function asString(v: unknown): string | undefined {
  if (typeof v !== "string") return undefined;
  const t = v.trim();
  return t || undefined;
}

function firstString(...vals: unknown[]): string {
  for (const v of vals) {
    if (typeof v === "string" && v.trim()) return v;
  }
  return "";
}

function positiveInt(v: unknown): number | undefined {
  const n = typeof v === "number" ? v : typeof v === "string" ? Number(v) : NaN;
  return Number.isInteger(n) && n > 0 ? n : undefined;
}

function titleCase(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

function clip(s: string, max = 48): string {
  const flat = s.replace(/\s+/g, " ");
  return flat.length > max ? `${flat.slice(0, max)}…` : flat;
}

/** Builds the canonical document for a dialect result and a chosen kind. */
export function bankToDocument(
  result: DialectResult,
  kind: "quiz" | "exam"
): ImportDocument {
  const title = result.header.title ?? "";
  if (kind === "exam") {
    return {
      version: IMPORT_FORMAT_VERSION,
      kind: "exam",
      exam: {
        title,
        subject: result.header.subject ?? "",
        // Placeholders the dialog requires the operator to confirm: a bank file
        // carries neither, and guessing a duration or grade for a graded exam
        // silently is worse than asking.
        grade: 12,
        duration_minutes: result.header.durationMinutes ?? 0,
        pass_marks: 0,
      },
      questions: result.questions,
    };
  }
  return {
    version: IMPORT_FORMAT_VERSION,
    kind: "quiz",
    quiz: { title, pass_percentage: 60 },
    questions: result.questions,
  };
}
