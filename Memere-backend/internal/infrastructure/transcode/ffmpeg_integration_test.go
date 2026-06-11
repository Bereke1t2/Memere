//go:build ffmpeg

// These tests shell out to real ffmpeg/ffprobe and are excluded from normal CI.
// Run locally where ffmpeg is installed:
//
//	go test -tags ffmpeg ./internal/infrastructure/transcode/... -v
package transcode

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// makeSampleClip synthesizes a ~6s test video (color bars + tone) so the
// thumbnail seek at 5s and all three renditions have content to encode.
func makeSampleClip(t *testing.T, dir string) string {
	t.Helper()
	src := filepath.Join(dir, "sample.mp4")
	cmd := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "testsrc=duration=6:size=1920x1080:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=440:duration=6",
		"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p",
		"-c:a", "aac", "-shortest", src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("make sample clip: %v\n%s", err, out)
	}
	return src
}

func TestFFmpegToHLS_ProducesLadder(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}
	work := t.TempDir()
	src := makeSampleClip(t, work)

	out, err := NewFFmpeg().ToHLS(context.Background(), src, work)
	if err != nil {
		t.Fatalf("ToHLS: %v", err)
	}

	if _, err := os.Stat(out.MasterPath); err != nil {
		t.Fatalf("master playlist missing: %v", err)
	}
	for _, name := range renditionNames {
		pl := filepath.Join(out.VariantDirs[name], "playlist.m3u8")
		if _, err := os.Stat(pl); err != nil {
			t.Fatalf("variant %s playlist missing: %v", name, err)
		}
		segs, _ := filepath.Glob(filepath.Join(out.VariantDirs[name], "seg_*.ts"))
		if len(segs) == 0 {
			t.Fatalf("variant %s produced no segments", name)
		}
	}
	if _, err := os.Stat(out.ThumbnailPath); err != nil {
		t.Fatalf("thumbnail missing: %v", err)
	}
	if out.DurationSec < 5 || out.DurationSec > 7 {
		t.Fatalf("duration = %ds, want ~6s", out.DurationSec)
	}
}
