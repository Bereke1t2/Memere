# Assessment import format

One uploaded file creates one whole quiz or exam — header, questions, options,
answer key and explanations. Upload it from **Quizzes → Import from file** or
**Exams → Import from file** (also available inside a course, on its Quizzes and
Exams tabs, which files the import under that course).

Templates to copy from:

- [`quiz-import-template.json`](../public/templates/quiz-import-template.json)
- [`exam-import-template.json`](../public/templates/exam-import-template.json)
- [`quiz-import-template.html`](../public/templates/quiz-import-template.html)

Two things hold for every import:

- **Nothing is saved until you confirm.** The import screen lists every question
  with the option it resolved as correct. Review that list — it is the answer key
  students will be graded against.
- **An import is all-or-nothing.** If question 40 is invalid, the 39 before it are
  not created either. A half-imported assessment would look complete to a student
  while silently mis-scoring, so the whole file is rejected with the offending
  question named by its position.

## JSON (recommended)

```json
{
  "version": 1,
  "kind": "quiz",
  "quiz": { "title": "Cell Biology — Practice Quiz", "pass_percentage": 60 },
  "questions": [
    {
      "text": "Which organelle is the main site of ATP production?",
      "points": 2,
      "explanation": "Mitochondria carry out oxidative phosphorylation.",
      "options": ["Nucleus", "Mitochondrion", "Ribosome", "Golgi apparatus"],
      "correct": "Mitochondrion"
    }
  ]
}
```

### Top level

| Field | Required | Notes |
|---|---|---|
| `version` | yes | Always `1`. |
| `kind` | yes | `"quiz"` or `"exam"`. Must match the page you upload from. |
| `quiz` | for quizzes | Header block, below. |
| `exam` | for exams | Header block, below. |
| `questions` | yes | At least 1, at most 500 per file. |

### `quiz` header

| Field | Required | Default | Notes |
|---|---|---|---|
| `title` | yes | — | Max 200 characters. |
| `course_id` | no | `null` | A course UUID, or `null` for a standalone quiz. Only an admin may create standalone quizzes; a standalone quiz stays invisible to students until published. |
| `lesson_id` | no | `null` | Attaches the quiz to one lesson. |
| `time_limit_seconds` | no | none | Omit or `null` for untimed. |
| `pass_percentage` | no | `0` | 0–100. |
| `randomize_questions` | no | `false` | |
| `max_attempts` | no | unlimited | |

### `exam` header

| Field | Required | Default | Notes |
|---|---|---|---|
| `title` | yes | — | Max 200 characters. |
| `subject` | yes | — | e.g. `"Physics"`. |
| `grade` | yes | — | 1–13 (13 = Freshman). |
| `duration_minutes` | yes | — | Server-enforced; the client timer is display-only. |
| `course_id` | no | `null` | A course UUID, or `null` for a standalone mock exam (admin only). |
| `pass_marks` | no | `0` | Must not exceed the total the questions add up to. |
| `instructions` | no | none | Max 5000 characters. |

`total_marks` is **not** a field — it is calculated from the question marks, so the
header total can never disagree with the parts.

### Questions

| Field | Required | Default | Notes |
|---|---|---|---|
| `text` | yes | — | The prompt. Max 5000 characters. |
| `type` | no | `multiple_choice` | `multiple_choice`, `true_false`, or `short_answer`. |
| `points` / `marks` | no | `1` | The same value. Use `marks` in exams; if both appear, `marks` wins. |
| `explanation` | no | none | Shown to the student **after** grading. Max 5000 characters. |
| `subject` | no | none | Feeds the subject-breakdown analytics. |
| `topic` | no | none | Feeds the weak-areas-by-topic analytics. |
| `answers` | one of | — | Explicit form, below. |
| `options` + `correct` | one of | — | Shorthand form, below. |

Question order is the order they appear in the file. Any `order_index` you supply
is ignored.

#### Options: the explicit form

Spells out each option and whether it is correct. Unambiguous, and the form the
exporter would produce:

```json
{
  "text": "Ribosomes are enclosed in a membrane.",
  "type": "true_false",
  "explanation": "Ribosomes are non-membrane-bound complexes of RNA and protein.",
  "answers": [
    { "text": "True",  "is_correct": false },
    { "text": "False", "is_correct": true }
  ]
}
```

#### Options: the shorthand form

Lists the options and points at the correct one, which is far less tedious for a
large file:

```json
{ "text": "Capital of France?", "options": ["London", "Berlin", "Paris"], "correct": "Paris" }
```

`correct` accepts any of:

| Form | Example | Meaning |
|---|---|---|
| Number | `"correct": 2` | 0-based index, so `2` is the **third** option. |
| Option text | `"correct": "Paris"` | Case-insensitive exact match. |
| Letter | `"correct": "C"` | `A` = first option, `B` = second, … |
| Array | `"correct": [0, 2]` | Multiple correct options; forms may be mixed. |

Resolution order is **option text first, then the letter label**. So if an option's
text genuinely is `"A"`, `"correct": "A"` selects that option rather than meaning
"the first one".

Use exactly one of `answers` or `options` per question. Supplying both, neither, or
a `correct` that matches nothing is rejected — a question with no correct answer
cannot be graded, so it is never imported silently.

## HTML

HTML is read structurally. It cannot be interpreted with certainty, so the import
screen always flags it and shows every resolved answer key for confirmation. For
an exact import, prefer JSON — or embed a JSON block in the HTML, which is used
verbatim and skips the guesswork entirely:

```html
<script type="application/json">
  { "version": 1, "kind": "quiz", "quiz": { "title": "…" }, "questions": [ … ] }
</script>
```

Otherwise the recognized structure is:

```html
<body data-assessment data-kind="quiz"
      data-title="Cell Biology — Practice Quiz" data-pass-percentage="60">

  <div class="question" data-marks="2" data-subject="Biology" data-topic="Cell structure">
    <p class="question-text">Which organelle is the main site of ATP production?</p>
    <ul>
      <li>Nucleus</li>
      <li data-correct>Mitochondrion</li>
      <li>Ribosome</li>
    </ul>
    <p class="explanation">Mitochondria carry out oxidative phosphorylation.</p>
  </div>

</body>
```

- **Header** — data attributes on `<body>` or any `[data-assessment]` element:
  `data-kind`, `data-title`, `data-pass-percentage`, `data-time-limit-seconds`
  for quizzes; plus `data-subject`, `data-grade`, `data-duration-minutes`,
  `data-pass-marks` for exams. The title falls back to `<title>` then `<h1>`.
- **Question** — any element with `data-question` or `class="question"`. Optional
  `data-marks` / `data-points`, `data-subject`, `data-topic`.
- **Prompt** — `[data-text]` or `.question-text`, else the first heading or `<p>`.
- **Options** — `<li>` elements, or any `[data-option]` elements.
- **Correct option** — marked in any one of four ways: `data-correct` on the
  option, `class="correct"`, wrapping it in `<strong>`/`<b>`, or a leading `*`
  (the `*` is stripped from the stored text).
- **Explanation** — `[data-explanation]` or `.explanation`.

Two options reading exactly "True" and "False" are detected as a `true_false`
question. More than one option marked correct is accepted but flagged, since it is
far more often a stray `<strong>` than a genuine multi-answer question.

A question with no marked option, or with fewer than two options, aborts the
import and names the question.

## API

The admin UI is the intended entry point, but the endpoints are plain JSON and
take the document above as the request body:

```
POST /api/v1/quizzes/import
POST /api/v1/exams/import
```

Both require a teacher or admin bearer token, and both are transactional.
