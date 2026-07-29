package convert

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMP4Path(t *testing.T) {
	if got := MP4Path("/tmp/Movie.2024.1080p.mkv"); got != "/tmp/Movie.2024.1080p.mp4" {
		t.Fatalf("MP4Path = %q", got)
	}
}

func TestNeedsRemuxSkipsFreshMP4(t *testing.T) {
	dir := t.TempDir()
	mkv := filepath.Join(dir, "clip.mkv")
	mp4 := filepath.Join(dir, "clip.mp4")
	if err := os.WriteFile(mkv, []byte("mkv"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mp4, []byte("mp4-data"), 0o644); err != nil {
		t.Fatal(err)
	}
	if NeedsRemux(mkv) {
		t.Fatal("should skip when mp4 already exists")
	}
}

func TestFindMKVInTree(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "pack")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	mkv := filepath.Join(sub, "a.mkv")
	if err := os.WriteFile(mkv, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "b.mp4"), []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := FindMKV(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || filepath.Base(got[0]) != "a.mkv" {
		t.Fatalf("FindMKV = %#v", got)
	}
}

func TestMKVToMP4Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}
	dir := t.TempDir()
	mkv := filepath.Join(dir, "sample.mkv")
	cmd := exec.Command("ffmpeg",
		"-hide_banner", "-loglevel", "error",
		"-f", "lavfi", "-i", "color=c=black:s=32x32:d=0.2",
		"-f", "lavfi", "-i", "sine=f=440:d=0.2",
		"-c:v", "libx264", "-pix_fmt", "yuv420p",
		"-c:a", "aac", "-shortest", "-y", mkv,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("create mkv: %v\n%s", err, out)
	}
	out, err := MKVToMP4(context.Background(), mkv, Options{RemoveSource: true})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(out) != ".mp4" {
		t.Fatalf("output = %q", out)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(mkv); !os.IsNotExist(err) {
		t.Fatal("source mkv should be removed")
	}
}
