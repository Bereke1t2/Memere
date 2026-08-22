// Package media holds small, stateless helpers shared by the content usecases
// (course, video) that operate on stored media artifacts. It depends only on the
// domain ports (service.ObjectStore) and entities, so it introduces no
// usecase-to-usecase coupling.
package media

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// PurgeVideo permanently deletes every object-storage artifact for v: the
// original upload, the HLS renditions (master + variant playlists + segments)
// and the thumbnail. Callers invoke it when a video is replaced or its lesson is
// deleted, so the bucket doesn't accumulate orphaned files.
//
// It is best-effort: failures are logged, never returned, so a storage hiccup
// never blocks the caller's database soft-delete (the row is the source of
// truth; a leaked object can be swept later). The scalar keys recorded on the
// row are deleted directly. The HLS .ts segments are NOT tracked individually,
// so when the store supports prefix deletion the whole hls/<id>/ and
// thumbnails/<id>/ trees are swept too. Single-file backends (Google Drive)
// implement no PrefixDeleter and need none — their one object is the scalar
// HLSMasterKey deleted above.
func PurgeVideo(ctx context.Context, store service.ObjectStore, v *entity.Video) {
	if store == nil || v == nil {
		return
	}
	for _, k := range []*string{
		v.OriginalFileKey, v.HLSMasterKey,
		v.Res480pKey, v.Res720pKey, v.Res1080pKey, v.ThumbnailKey,
	} {
		if k == nil || *k == "" {
			continue
		}
		if err := store.Delete(ctx, *k); err != nil {
			slog.WarnContext(ctx, "purge video object failed", "video_id", v.ID, "key", *k, "err", err)
		}
	}
	if pd, ok := store.(service.PrefixDeleter); ok {
		for _, prefix := range []string{
			fmt.Sprintf("hls/%s/", v.ID),
			fmt.Sprintf("thumbnails/%s/", v.ID),
		} {
			if err := pd.DeletePrefix(ctx, prefix); err != nil {
				slog.WarnContext(ctx, "purge video prefix failed", "video_id", v.ID, "prefix", prefix, "err", err)
			}
		}
	}
}
