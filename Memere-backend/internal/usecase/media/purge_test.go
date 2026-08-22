package media

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// scalarStore implements only Delete (embedding the ObjectStore interface for the
// unused methods) and deliberately does NOT implement PrefixDeleter — it models a
// single-file backend like Google Drive.
type scalarStore struct {
	service.ObjectStore
	deleted map[string]bool
}

func newScalarStore() *scalarStore { return &scalarStore{deleted: map[string]bool{}} }

func (s *scalarStore) Delete(_ context.Context, key string) error {
	s.deleted[key] = true
	return nil
}

// prefixStore adds DeletePrefix, modelling S3/MinIO/B2.
type prefixStore struct {
	*scalarStore
	prefixes map[string]bool
}

func newPrefixStore() *prefixStore {
	return &prefixStore{scalarStore: newScalarStore(), prefixes: map[string]bool{}}
}

func (s *prefixStore) DeletePrefix(_ context.Context, prefix string) error {
	s.prefixes[prefix] = true
	return nil
}

var _ service.PrefixDeleter = (*prefixStore)(nil)

func sampleVideo() *entity.Video {
	id := uuid.New()
	orig := "originals/c/l/" + id.String() + "/source.mp4"
	master := "hls/" + id.String() + "/master.m3u8"
	v480 := "hls/" + id.String() + "/480p/playlist.m3u8"
	thumb := "thumbnails/" + id.String() + "/thumb.jpg"
	return &entity.Video{
		ID: id, OriginalFileKey: &orig, HLSMasterKey: &master,
		Res480pKey: &v480, ThumbnailKey: &thumb,
	}
}

func TestPurgeVideo_ScalarStoreDeletesKnownKeys(t *testing.T) {
	st := newScalarStore()
	v := sampleVideo()

	PurgeVideo(context.Background(), st, v)

	for _, k := range []string{*v.OriginalFileKey, *v.HLSMasterKey, *v.Res480pKey, *v.ThumbnailKey} {
		if !st.deleted[k] {
			t.Errorf("scalar key %q was not deleted", k)
		}
	}
}

func TestPurgeVideo_PrefixStoreAlsoSweepsSegments(t *testing.T) {
	st := newPrefixStore()
	v := sampleVideo()

	PurgeVideo(context.Background(), st, v)

	if !st.deleted[*v.OriginalFileKey] {
		t.Error("original object was not deleted")
	}
	if !st.prefixes["hls/"+v.ID.String()+"/"] {
		t.Error("hls/ prefix (segments) was not swept")
	}
	if !st.prefixes["thumbnails/"+v.ID.String()+"/"] {
		t.Error("thumbnails/ prefix was not swept")
	}
}

func TestPurgeVideo_SkipsNilAndEmptyKeys(t *testing.T) {
	st := newScalarStore()
	only := "originals/only.mp4"
	// A video with just one key set: nil/empty keys must be skipped, not deleted.
	PurgeVideo(context.Background(), st, &entity.Video{ID: uuid.New(), OriginalFileKey: &only})
	if !st.deleted[only] {
		t.Error("the one set key was not deleted")
	}
	if st.deleted[""] {
		t.Error("empty key should never be passed to Delete")
	}
}

func TestPurgeVideo_NilSafe(t *testing.T) {
	// Must not panic on a nil store or nil video.
	PurgeVideo(context.Background(), nil, sampleVideo())
	PurgeVideo(context.Background(), newScalarStore(), nil)
}
