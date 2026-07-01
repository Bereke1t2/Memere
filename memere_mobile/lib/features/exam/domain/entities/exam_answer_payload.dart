/// Maps a question ID to its answer value.
///
/// - single choice / true-false: `questionId -> answerId` (String)
/// - multi select: `questionId -> List<String>`
/// - short answer: `questionId -> String`
///
/// Backend accepts `map[string]any`.
typedef ExamAnswerPayload = Map<String, Object>;
