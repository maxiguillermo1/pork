package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maxiguillermo1/pork/internal/config"
	"github.com/maxiguillermo1/pork/internal/engine"
)

func TestManualRemuxDownloadRequiresFinishedState(t *testing.T) {
	cfg := config.Default(t.TempDir())
	app := &App{cfg: cfg}

	cmd := app.manualRemuxDownload(downloadItem{
		Magnet: "magnet:?xt=urn:btih:abc",
		State:  engine.StateDownloading,
	})
	if cmd == nil {
		t.Fatal("expected clearErrCmd")
	}
	if app.errText == "" {
		t.Fatal("expected error text for active download")
	}
}

func TestManualRemuxDownloadFindsMKV(t *testing.T) {
	dir := t.TempDir()
	mkv := filepath.Join(dir, "clip.mkv")
	if err := os.WriteFile(mkv, []byte("mkv"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default(t.TempDir())
	app := &App{cfg: cfg}

	cmd := app.manualRemuxDownload(downloadItem{
		Magnet:      "magnet:?xt=urn:btih:def",
		DownloadDir: dir,
		DataPath:    dir,
		State:       engine.StateDone,
	})
	if cmd == nil {
		t.Fatal("expected remux command")
	}
	if job := app.convertJobs["magnet:?xt=urn:btih:def"]; job == nil || job.total != 1 {
		t.Fatalf("convert job = %#v, want 1 mkv queued", app.convertJobs["magnet:?xt=urn:btih:def"])
	}
}

func TestHasMKVForRemux(t *testing.T) {
	dir := t.TempDir()
	mkv := filepath.Join(dir, "movie.mkv")
	if err := os.WriteFile(mkv, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default(t.TempDir())
	app := &App{cfg: cfg}
	it := downloadItem{DownloadDir: dir, DataPath: dir, State: engine.StateDone}
	if !app.hasMKVForRemux(it) {
		t.Fatal("expected mkv to be detected")
	}
}
