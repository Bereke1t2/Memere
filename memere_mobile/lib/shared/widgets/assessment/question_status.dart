import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';

/// Where a student stands on one question.
///
/// Ordered by how much attention the question still needs, so a palette can
/// sort or prioritise by it without re-deriving the meaning.
enum QuestionStatus {
  /// Answered and not flagged — nothing left to do.
  answered,

  /// Answered but flagged to revisit.
  answeredFlagged,

  /// Flagged with no answer yet.
  flagged,

  /// Neither answered nor flagged. The state that loses marks silently.
  blank;

  static QuestionStatus of({required bool answered, required bool flagged}) {
    if (answered) return flagged ? QuestionStatus.answeredFlagged : QuestionStatus.answered;
    return flagged ? QuestionStatus.flagged : QuestionStatus.blank;
  }

  bool get isAnswered =>
      this == QuestionStatus.answered || this == QuestionStatus.answeredFlagged;

  bool get isFlagged =>
      this == QuestionStatus.flagged || this == QuestionStatus.answeredFlagged;

  /// Fill for the palette cell. Blank is deliberately the quietest: the grid
  /// should read as "what is left" at a glance, not as a wall of colour.
  Color get fill => switch (this) {
        QuestionStatus.answered => AppColors.accentPrimary,
        QuestionStatus.answeredFlagged => AppColors.accentPrimary,
        QuestionStatus.flagged => AppColors.warningSurface,
        QuestionStatus.blank => AppColors.bgTertiary,
      };

  Color get border => switch (this) {
        QuestionStatus.answered => AppColors.accentPrimary,
        QuestionStatus.answeredFlagged => AppColors.warning,
        QuestionStatus.flagged => AppColors.warning,
        QuestionStatus.blank => AppColors.borderStrong,
      };

  Color get foreground => switch (this) {
        QuestionStatus.answered => AppColors.textInverse,
        QuestionStatus.answeredFlagged => AppColors.textInverse,
        QuestionStatus.flagged => AppColors.warning,
        QuestionStatus.blank => AppColors.textSecondary,
      };

  String get label => switch (this) {
        QuestionStatus.answered => 'Answered',
        QuestionStatus.answeredFlagged => 'Answered, flagged',
        QuestionStatus.flagged => 'Flagged',
        QuestionStatus.blank => 'Not answered',
      };
}
