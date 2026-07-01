# Improve UI/UX — Memere Mobile

A focused design-system upgrade that takes the app from "functional but generic"
to **professional, cohesive, and quietly animated**. It builds on the existing
tokens (`AppColors`, `AppTextStyles`, `AppSizes`, `AppTheme`) rather than
replacing them, so it can be applied incrementally without breaking screens.

## Design principles

1. **One accent, used sparingly.** Indigo (`#6366F1`) is the only brand hue.
   Color earns attention; everything else is neutral. No rainbow.
2. **Depth through soft shadows, not borders everywhere.** Cards lift off the
   background with subtle shadows instead of hard 1px outlines.
3. **Motion with intent.** Every transition has a reason: feedback (press),
   continuity (hero/shared-axis), or arrival (staggered entrance). Durations are
   short (180–260ms) and consistent. Never animate just to decorate.
4. **Hierarchy first.** Generous spacing, a tight type scale, and clear section
   headers make scanning effortless.
5. **Every state is designed.** Loading (shimmer), empty, error, and success all
   look intentional — not a spinner in the middle of a blank screen.

## How to use this skill

Read the files in order. Each is self-contained and ends with a checklist.

| File | Scope |
|------|-------|
| `SKILL_1_foundations.md` | Tokens: refined color, elevation/shadows, motion, radius, typography; theme wiring |
| `SKILL_2_components.md` | Reusable animated building blocks: tap feedback, cards, shimmer, badges, section headers, page transitions |
| `SKILL_3_screens_and_motion.md` | Applying it per screen: catalog, course detail (hero), checkout, result (success animation), My Learning, Profile; plus the QA checklist |

Each template block is headed with the exact target path
(`### FILE — lib/...`). Copy the block into that file. Code is plain Flutter +
the packages already in `pubspec.yaml` (`shimmer`, `lottie`,
`cached_network_image`, `flutter_svg`) — **no new dependencies required** unless
explicitly noted.

## Constraints

- Dark-first. All values target the dark theme.
- No new packages unless a section says otherwise (one optional: `lottie` is
  already present, so success animations use it or a pure-Flutter fallback).
- Keep `flutter analyze` at 0 errors after each file.
- Don't regress accessibility: minimum 44×44 tap targets, AA text contrast.

## Order of impact (if you only do some of it)

1. **Shadows + AnimatedTap** (`SKILL_2`) — instantly makes the app feel premium.
2. **Page transitions + staggered lists** (`SKILL_2`/`SKILL_3`) — flow feels designed.
3. **Hero course image + success animation** (`SKILL_3`) — signature moments.
4. **Token cleanup** (`SKILL_1`) — consistency that compounds.
