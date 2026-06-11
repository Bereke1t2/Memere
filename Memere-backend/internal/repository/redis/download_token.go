package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// DownloadTokenStore backs the single-use download flow (spec §8.3 step 2). A
// CloudFront/S3 signed URL is time-limited but *replayable* inside its window;
// to honor "single-use" we gate the download behind a server-side token that is
// atomically consumed (GETDEL) on first redeem, so a second redemption misses.
//
// Streaming is deliberately NOT single-use — HLS fetches many segments within
// the window — so only the download flow uses this store.
type DownloadTokenStore struct {
	client *goredis.Client
}

// NewDownloadTokenStore builds a DownloadTokenStore over a go-redis client.
func NewDownloadTokenStore(client *goredis.Client) *DownloadTokenStore {
	return &DownloadTokenStore{client: client}
}

var _ service.DownloadTokens = (*DownloadTokenStore)(nil)

// downloadTokenKey namespaces tokens so they cannot collide with other keys.
func downloadTokenKey(token string) string {
	return "dl:" + token
}

// the stored value is "videoID|userID"; the separator never appears in a UUID.
const downloadTokenSep = "|"

// Issue stores token -> videoID|userID with the given TTL (the signed-URL
// window). The token must be unguessable; the caller mints it with crypto/rand.
func (s *DownloadTokenStore) Issue(ctx context.Context, token, videoID, userID string, ttl time.Duration) error {
	val := videoID + downloadTokenSep + userID
	if err := s.client.Set(ctx, downloadTokenKey(token), val, ttl).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// Consume atomically fetches and deletes the token (GETDEL). The first call
// returns the bound (videoID, userID, true); any later call — or an expired /
// unknown token — returns ok=false with no error (a miss is a normal outcome the
// caller maps to 404/410).
func (s *DownloadTokenStore) Consume(ctx context.Context, token string) (videoID, userID string, ok bool, err error) {
	val, gerr := s.client.GetDel(ctx, downloadTokenKey(token)).Result()
	if errors.Is(gerr, goredis.Nil) {
		return "", "", false, nil
	}
	if gerr != nil {
		return "", "", false, apperror.Internal(gerr)
	}
	vid, uid, found := strings.Cut(val, downloadTokenSep)
	if !found {
		// Corrupt value — treat as a miss rather than leaking a malformed pair.
		return "", "", false, nil
	}
	return vid, uid, true, nil
}
