# FIRST_PROMPT.md — How to Start with Antigravity
# Copy the prompt below word-for-word into Antigravity to kick off Phase 1

---

## YOUR FIRST PROMPT (copy-paste this exactly)

```
You are building a Flutter mobile app called ExamPrep (project folder: memere_mobile).

Before writing any code, read these files in this exact order:
1. SKILL.md               — master rules, architecture, naming, design system overview
2. phase1/SKILL_1.md      — project scaffold, pubspec.yaml, folder structure
3. phase1/SKILL_2.md      — complete design system (colors, typography, theme, shared widgets)
4. phase1/SKILL_3.md      — core infrastructure, auth domain, auth data, auth screens

The full architecture specification is in:
  memere_mobile/docs/memere_Design_Specification.md
Read that too before starting.

We are on PHASE 1 ONLY. Do not build anything from Phase 2 or later.

Phase 1 goal: A compilable Flutter app that shows:
  Splash screen → Onboarding (3 pages) → Login screen → Register screen
  Dark theme, Sora + DM Sans fonts, purple accent (#6C63FF)

Build in this order:
  Step 1: Create pubspec.yaml and all folders (SKILL_1)
  Step 2: Create all design system files (SKILL_2)
  Step 3: Create core infrastructure files (SKILL_3 Part A)
  Step 4: Create auth domain layer (SKILL_3 Part B)
  Step 5: Create auth data layer (SKILL_3 Part C)
  Step 6: Create auth presentation layer + all screens (SKILL_3 Part D)
  Step 7: Run flutter analyze and fix any issues
  Step 8: Confirm the Phase 1 checklist in SKILL_3 passes

Do not skip any step. Do not add features not in the skill files.
When done, tell me "Phase 1 complete" and show me the checklist results.
```

---

## HOW TO MOVE TO EACH PHASE

### → Phase 2 (Course Browsing)
```
Phase 1 is complete and all checks pass.
Read SKILL.md and phase2/SKILL_1.md.
We are starting Phase 2: Course Browsing.
Reference: memere_mobile/docs/memere_Design_Specification.md
Build in order as specified in the phase2 skill files.
```

### → Phase 3 (Video Player)
```
Phase 2 is complete.
Read SKILL.md and phase3/SKILL_1.md.
We are starting Phase 3: HLS Video Player + Offline Download.
```

### → Phase 4 (Quiz Engine)
```
Phase 3 is complete.
Read SKILL.md and phase4/SKILL_1.md.
We are starting Phase 4: Quiz Engine.
```

_(Continue this pattern for each phase)_

---

## TIPS FOR WORKING WITH ANTIGRAVITY

1. **One phase at a time** — never ask it to build Phase 3 features while in Phase 1
2. **If it drifts** — remind it: "Read SKILL.md again and follow the architecture rules"
3. **If it uses setState** — remind it: "No setState in feature screens. Use Riverpod only."
4. **If it skips a layer** — remind it: "Screens must not call repositories directly. Use use cases via providers."
5. **If it hardcodes a URL** — remind it: "Use Env.baseUrl from AppConstants, never hardcode."
6. **Start each new session** with: "Read SKILL.md and the current phase skill files before continuing."

---

## FILES TO UPLOAD TO ANTIGRAVITY PROJECT

Put these files in your Antigravity project (or reference folder):

```
SKILL.md
FIRST_PROMPT.md
phase1/
  SKILL_1.md
  SKILL_2.md
  SKILL_3.md
```

When Phase 2 skill files are ready, add:
```
phase2/
  SKILL_1.md
  SKILL_2.md
```

---

## WHAT EACH PHASE BUILDS

| Phase | Name | Key Screens / Features |
|-------|------|----------------------|
| 1 | Foundation | Splash, Onboarding, Login, Register — full auth flow |
| 2 | Courses | Home (bottom nav), Course list, Course detail, Search, Subject filters |
| 3 | Video | HLS video player, adaptive bitrate, offline download, progress saving |
| 4 | Quiz | Quiz screen, question cards, submission, results with breakdown |
| 5 | Mock Exam | Timed exam, server-synced timer, exam analytics, percentile rank |
| 6 | Payment | Chapa/Telebirr WebView flow, enrollment confirmation, purchase history |
| 7 | Dashboard | Student progress dashboard, streaks, completion %, weak areas |
| 8 | Notifications | FCM push setup, in-app notification center, notification settings |
| 9 | Teacher | Content upload screens, quiz builder, course management (teacher role only) |
| 10 | Polish | Performance, accessibility, Play Store assets, release build |

---

## IMPORTANT REMINDERS

- The **backend is Go** (not built here) — Flutter calls the REST API defined in the spec doc
- **Correct answers are NEVER in the client** — quiz grading is always server-side
- **Exam timer is server-side** — the Flutter timer is display-only
- **Video = pre-signed CDN URLs** — never direct S3 access
- **Payments use idempotency keys** — handled in the Payment service (Phase 6)
- **Offline first**: Hive cache + SyncQueue pattern from SKILL.md Section 10
