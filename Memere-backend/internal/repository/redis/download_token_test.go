package redis

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDownloadToken_SingleUseConsume(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	store := NewDownloadTokenStore(c)
	ctx := context.Background()

	token := "tok_" + uuid.NewString()
	videoID := uuid.NewString()
	userID := uuid.NewString()
	defer c.Del(ctx, downloadTokenKey(token))

	if err := store.Issue(ctx, token, videoID, userID, time.Minute); err != nil {
		t.Fatalf("Issue: %v", err)
	}

	// First consume returns the bound pair.
	gotVid, gotUID, ok, err := store.Consume(ctx, token)
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	if !ok {
		t.Fatalf("first Consume missed")
	}
	if gotVid != videoID || gotUID != userID {
		t.Fatalf("Consume = (%s,%s), want (%s,%s)", gotVid, gotUID, videoID, userID)
	}

	// Second consume misses (GETDEL removed it) — single use.
	_, _, ok, err = store.Consume(ctx, token)
	if err != nil {
		t.Fatalf("second Consume err: %v", err)
	}
	if ok {
		t.Fatalf("second Consume should miss (token is single-use)")
	}
}

func TestDownloadToken_UnknownTokenMisses(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	store := NewDownloadTokenStore(c)

	_, _, ok, err := store.Consume(context.Background(), "tok_"+uuid.NewString())
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	if ok {
		t.Fatalf("unknown token should miss")
	}
}

func TestDownloadToken_Expires(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	store := NewDownloadTokenStore(c)
	ctx := context.Background()

	token := "tok_" + uuid.NewString()
	defer c.Del(ctx, downloadTokenKey(token))
	if err := store.Issue(ctx, token, uuid.NewString(), uuid.NewString(), 50*time.Millisecond); err != nil {
		t.Fatalf("Issue: %v", err)
	}
	time.Sleep(120 * time.Millisecond)
	_, _, ok, err := store.Consume(ctx, token)
	if err != nil {
		t.Fatalf("Consume: %v", err)
	}
	if ok {
		t.Fatalf("expired token should miss")
	}
}
