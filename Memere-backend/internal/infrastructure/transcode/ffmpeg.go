// Package transcode is the FFmpeg-backed implementation of the video transcoder
// boundary. The worker depends on the Transcoder interface, never on os/exec, so
// the local-ffmpeg implementation here can later be swapped for AWS MediaConvert
// without touching the worker. It is the only place os/exec and ffmpeg flags
// live.
package transcode

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// renditionNames are the variant stream names produced, matching the spec §8.2
// ladder. Only 480p/720p/1080p are encoded: the §4.2.4 videos table stores keys
// for exactly these three (the 360p row in the player table has no column), so
// it is intentionally omitted here.
var renditionNames = []string{"480p", "720p", "1080p"}

// Renditions is the local on-disk result of a transcode: a master playlist, one
// directory of playlist+segments per variant, a thumbnail, and the probed source
// duration. Paths are absolute under the worker's workDir; the worker uploads
// them to object storage under the Skill 1 key scheme.
type Renditions struct {
	MasterPath    string            // local path to master.m3u8
	VariantDirs   map[string]string // "480p" -> local dir holding playlist.m3u8 + seg_*.ts
	ThumbnailPath string
	DurationSec   int
}

// Transcoder turns a local source file into an HLS adaptive-bitrate ladder.
type Transcoder interface {
	// ToHLS reads srcPath and writes HLS outputs into workDir, returning the
	// produced paths. workDir must already exist; ToHLS creates per-variant
	// subdirectories under it.
	ToHLS(ctx context.Context, srcPath, workDir string) (*Renditions, error)
	// Probe returns the source media duration in whole seconds.
	Probe(ctx context.Context, srcPath string) (int, error)
}

// FFmpeg shells out to the ffmpeg/ffprobe binaries on PATH.
type FFmpeg struct {
	bin     string
	ffprobe string
}

// NewFFmpeg builds an FFmpeg runner using ffmpeg/ffprobe from PATH.
func NewFFmpeg() *FFmpeg { return &FFmpeg{bin: "ffmpeg", ffprobe: "ffprobe"} }

// thumbnailSeconds is where the poster frame is grabbed (spec §8.1: 5s mark).
const thumbnailSeconds = 5

// ToHLS encodes the three-rendition ladder in a single ffmpeg invocation using
// -var_stream_map, then extracts the thumbnail and probes the duration.
func (f *FFmpeg) ToHLS(ctx context.Context, src, workDir string) (*Renditions, error) {
	variantDirs := make(map[string]string, len(renditionNames))
	for _, name := range renditionNames {
		dir := filepath.Join(workDir, name)
		// The HLS muxer does not always create the %v subdirectories itself, so
		// make them up front.
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return nil, fmt.Errorf("mkdir variant %s: %w", name, err)
		}
		variantDirs[name] = dir
	}

	args := []string{
		"-y", "-i", src,
		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];" +
			"[v1]scale=w=854:h=480[v480];" +
			"[v2]scale=w=1280:h=720[v720];" +
			"[v3]scale=w=1920:h=1080[v1080]",
		// 480p
		"-map", "[v480]", "-map", "0:a:0?", "-c:v:0", "libx264", "-b:v:0", "800k", "-maxrate:v:0", "856k", "-bufsize:v:0", "1200k",
		// 720p
		"-map", "[v720]", "-map", "0:a:0?", "-c:v:1", "libx264", "-b:v:1", "1500k", "-maxrate:v:1", "1600k", "-bufsize:v:1", "2100k",
		// 1080p
		"-map", "[v1080]", "-map", "0:a:0?", "-c:v:2", "libx264", "-b:v:2", "3000k", "-maxrate:v:2", "3200k", "-bufsize:v:2", "4200k",
		"-c:a", "aac", "-b:a:0", "96k", "-b:a:1", "128k", "-b:a:2", "128k",
		"-preset", "veryfast",
		"-f", "hls",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(workDir, "%v", "seg_%05d.ts"),
		"-master_pl_name", "master.m3u8",
		"-var_stream_map", "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p",
		filepath.Join(workDir, "%v", "playlist.m3u8"),
	}

	if err := f.run(ctx, f.bin, args...); err != nil {
		return nil, fmt.Errorf("ffmpeg hls: %w", err)
	}

	thumb := filepath.Join(workDir, "thumb.jpg")
	if err := f.run(ctx, f.bin, "-y", "-ss", strconv.Itoa(thumbnailSeconds), "-i", src, "-frames:v", "1", thumb); err != nil {
		// A clip shorter than the seek point: retry from the first frame so a
		// short video still gets a thumbnail.
		if err2 := f.run(ctx, f.bin, "-y", "-i", src, "-frames:v", "1", thumb); err2 != nil {
			return nil, fmt.Errorf("ffmpeg thumbnail: %w", err)
		}
	}

	dur, err := f.Probe(ctx, src)
	if err != nil {
		return nil, err
	}

	return &Renditions{
		MasterPath:    filepath.Join(workDir, "master.m3u8"),
		VariantDirs:   variantDirs,
		ThumbnailPath: thumb,
		DurationSec:   dur,
	}, nil
}

// Probe returns the source duration in whole seconds via ffprobe.
func (f *FFmpeg) Probe(ctx context.Context, src string) (int, error) {
	out, err := f.output(ctx, f.ffprobe,
		"-v", "error", "-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1", src)
	if err != nil {
		return 0, fmt.Errorf("ffprobe: %w", err)
	}
	secs, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
	if err != nil {
		return 0, fmt.Errorf("ffprobe parse duration %q: %w", out, err)
	}
	return int(secs + 0.5), nil
}

// run executes a command, surfacing the tail of stderr on failure for
// diagnostics. stderr is never logged on success (it is verbose and noisy).
func (f *FFmpeg) run(ctx context.Context, bin string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %s", err, tail(stderr.String(), 500))
	}
	return nil
}

// output runs a command and returns its stdout.
func (f *FFmpeg) output(ctx context.Context, bin string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, tail(stderr.String(), 500))
	}
	return stdout.String(), nil
}

// tail returns the last n bytes of s, for bounded error messages.
func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
