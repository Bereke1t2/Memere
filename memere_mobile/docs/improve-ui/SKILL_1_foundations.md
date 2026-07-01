# SKILL_1 — Foundations (Tokens & Theme)

> Read `README.md` first. This file upgrades the **design tokens**: shadows,
> motion, radius helpers, and a few color refinements — then wires them into the
> theme. It introduces **no new screens** and **no new packages**. Everything
> here is additive; existing `AppColors` / `AppSizes` / `AppTextStyles` values
> are kept so no screen breaks.

Each block is headed with its exact target path. Copy the block in. Run
`flutter analyze` after each file — keep it at **0 errors**.

---

## OBJECTIVE

Give the app a physical, premium feel before touching a single screen:

1. **Elevation system** — soft, layered shadows instead of 1px borders.
2. **Motion system** — one set of durations + curves so every animation in the
   app agrees with every other.
3. **Color refinements** — a real scrim, a hairline, and an "on-glow" so cards
   can lift without looking flat.
4. **Theme wiring** — push the new tokens into `ThemeData` so defaults improve
   for free.

---

## 1. Elevation / shadows

### FILE — lib/core/constants/app_shadows.dart  *(new)*

```dart
import 'package:flutter/material.dart';

/// Layered shadow tokens for a dark UI.
///
/// On a near-black background a single hard shadow reads as a smudge. Each
/// elevation here stacks a tight contact shadow (definition) over a wide,
/// very soft ambient shadow (lift). Use these instead of `Border.all` to
/// separate a surface from its background.
abstract class AppShadows {
  /// Resting cards, list tiles, chips.
  static const List<BoxShadow> sm = [
    BoxShadow(color: Color(0x33000000), blurRadius: 2,  offset: Offset(0, 1)),
    BoxShadow(color: Color(0x1F000000), blurRadius: 8,  offset: Offset(0, 4)),
  ];

  /// Raised cards, the pressed/active card target, popovers.
  static const List<BoxShadow> md = [
    BoxShadow(color: Color(0x40000000), blurRadius: 4,  offset: Offset(0, 2)),
    BoxShadow(color: Color(0x29000000), blurRadius: 16, offset: Offset(0, 8)),
  ];

  /// Bottom sheets, dialogs, the floating checkout bar.
  static const List<BoxShadow> lg = [
    BoxShadow(color: Color(0x4D000000), blurRadius: 8,  offset: Offset(0, 4)),
    BoxShadow(color: Color(0x33000000), blurRadius: 32, offset: Offset(0, 16)),
  ];

  /// Indigo "lift" — a colored glow under a primary CTA or selected item.
  /// Use sparingly; this is the one place color leaks into elevation.
  static const List<BoxShadow> accentGlow = [
    BoxShadow(color: Color(0x336366F1), blurRadius: 20, offset: Offset(0, 8)),
  ];
}
```

---

## 2. Motion

### FILE — lib/core/constants/app_motion.dart  *(new)*

```dart
import 'package:flutter/animation.dart';

/// One motion language for the whole app. If a duration or curve isn't here,
/// it shouldn't be in a widget. Keeping these centralized is what makes the
/// app feel "designed" rather than "animated in ten different ways".
abstract class AppMotion {
  // ── Durations ────────────────────────────────────────────────────────────
  /// Tap / press feedback. Must feel instant.
  static const Duration fast = Duration(milliseconds: 120);

  /// The default for almost everything: fades, slides, expands.
  static const Duration base = Duration(milliseconds: 220);

  /// Page transitions, hero flights, larger reveals.
  static const Duration slow = Duration(milliseconds: 320);

  /// Staggered list — delay added per item index (cap the total; see SKILL_3).
  static const Duration stagger = Duration(milliseconds: 40);

  // ── Curves ───────────────────────────────────────────────────────────────
  /// Standard ease for entering + moving elements.
  static const Curve standard = Cubic(0.2, 0.0, 0.0, 1.0); // emphasized-decelerate

  /// Elements leaving the screen.
  static const Curve exit = Cubic(0.3, 0.0, 0.8, 0.15);

  /// Press-down / spring-y feedback.
  static const Curve emphasized = Curves.easeOutBack;
}
```

---

## 3. Color refinements (additive)

Add these to the **bottom** of `AppColors` — do not change existing values.

### FILE — lib/core/constants/app_colors.dart  *(append inside the class)*

```dart
  // ── Refinements (SKILL_1) ────────────────────────────────────────────────
  /// True modal scrim. Darker than bgOverlay so sheets read as "on top".
  static const Color scrim = Color(0xCC07080A);

  /// Hairline used on top of a shadowed card for a crisp top edge (optional).
  static const Color hairline = Color(0x14FFFFFF);

  /// Pressed-state tint laid over a surface (8% white).
  static const Color pressedOverlay = Color(0x14FFFFFF);

  /// Skeleton base + highlight for shimmer (see SKILL_2).
  static const Color skeletonBase      = Color(0xFF1A1D23);
  static const Color skeletonHighlight = Color(0xFF262A32);
```

---

## 4. Radius helper (optional but convenient)

### FILE — lib/core/constants/app_sizes.dart  *(append inside the class)*

```dart
  // ── Radius helpers (SKILL_1) ─────────────────────────────────────────────
  static const Duration _unused = Duration.zero; // keep imports stable if added later
```

> Skip the block above if it adds an unused warning — it's only a placeholder
> reminder. The real win is just using the existing `radius*` constants
> consistently (`radiusMd` for cards, `radiusFull` for chips/badges).

---

## 5. Theme wiring

Two small changes to `app_theme.dart`: drop the hard border on cards in favor of
shadow-based elevation, and give the default page transition the app's motion.

### FILE — lib/core/theme/app_theme.dart  *(edit `cardTheme` + add `pageTransitionsTheme`)*

```dart
      // ── Cards (SKILL_1: shadow-based, no hard border) ────────────────────
      cardTheme: CardThemeData(
        color: AppColors.bgSecondary,
        elevation: 0, // we draw our own shadows via AppShadows
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSizes.radiusLg),
        ),
      ),

      // ── Page transitions (SKILL_1) ───────────────────────────────────────
      pageTransitionsTheme: const PageTransitionsTheme(
        builders: {
          TargetPlatform.android: _SharedAxisFadeTransitionBuilder(),
          TargetPlatform.iOS: CupertinoPageTransitionsBuilder(),
        },
      ),
```

> `_SharedAxisFadeTransitionBuilder` is defined in **SKILL_2** (page
> transitions). If you wire the theme before adding it, temporarily use
> `FadeUpwardsPageTransitionsBuilder()` so the file still compiles.

---

## Checklist — SKILL_1

- [ ] `app_shadows.dart` created; `AppShadows.sm/md/lg/accentGlow` resolve.
- [ ] `app_motion.dart` created; durations + curves resolve.
- [ ] `AppColors` refinements appended (no existing value changed).
- [ ] `cardTheme` no longer paints a hard border.
- [ ] `flutter analyze` → **0 errors** (warnings about unused placeholder OK to remove).
- [ ] App still builds and every existing screen renders unchanged.

Once green, go to **SKILL_2** to build the components that consume these tokens.
