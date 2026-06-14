package progress

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// sameDay reports whether two UTC dates fall on the same calendar day in the
// configured streak timezone.
func sameDay(a, b time.Time, loc *time.Location) bool {
	ay, am, ad := a.In(loc).Date()
	by, bm, bd := b.In(loc).Date()
	return ay == by && am == bm && ad == bd
}

// isYesterday reports whether a is exactly one calendar day before b in loc.
func isYesterday(a, b time.Time, loc *time.Location) bool {
	return sameDay(a.AddDate(0, 0, 1), b, loc)
}

// truncateToDay returns midnight of t's calendar day in loc.
func truncateToDay(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

// bumpStreak updates the study streak for studentID based on the current time.
// Rules:
//   - First ever study: streak = 1
//   - Same calendar day: no change (already counted)
//   - Yesterday: streak++
//   - Any gap: reset to 1
//
// longest_streak is updated if current > longest. The timezone comes from
// s.loc (STREAK_TZ env, default Africa/Addis_Ababa per spec).
func (s *Service) bumpStreak(ctx context.Context, studentID uuid.UUID) error {
	st, err := s.repo.GetStreak(ctx, studentID)
	if err != nil {
		return err
	}

	today := truncateToDay(s.clock(), s.loc)

	if st.LastStudyDate != nil {
		last := truncateToDay(*st.LastStudyDate, s.loc)
		switch {
		case sameDay(last, today, s.loc):
			return nil // already counted today
		case isYesterday(last, today, s.loc):
			st.CurrentStreak++
		default:
			st.CurrentStreak = 1 // gap → reset
		}
	} else {
		st.CurrentStreak = 1
	}

	if st.CurrentStreak > st.LongestStreak {
		st.LongestStreak = st.CurrentStreak
	}
	st.LastStudyDate = &today
	return s.repo.UpsertStreak(ctx, st)
}
