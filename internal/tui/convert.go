package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/maxiguillermo1/pork/internal/convert"
	"github.com/maxiguillermo1/pork/internal/engine"
)

type convertJob struct {
	magnet    string
	total     int
	converted int
	failed    int
	done      bool
	lastErr   string
}

type convertDoneMsg struct {
	magnet    string
	converted int
	failed    int
	err       error
}

func (a *App) initConvertTracker() {
	if a.convertJobs == nil {
		a.convertJobs = make(map[string]*convertJob)
	}
}

func (a *App) remuxEnabled() bool {
	return a.cfg != nil && a.cfg.VideoRemuxEnabled()
}

func (a *App) conversionNote(magnet string) string {
	if !a.remuxEnabled() {
		return ""
	}
	job, ok := a.convertJobs[magnet]
	if !ok {
		return ""
	}
	if !job.done {
		if job.total <= 1 {
			return "remuxing → mp4"
		}
		return fmt.Sprintf("remuxing → mp4 (%d files)", job.total)
	}
	if job.converted > 0 && job.failed == 0 {
		if job.converted == 1 {
			return "remuxed to mp4"
		}
		return fmt.Sprintf("remuxed %d files to mp4", job.converted)
	}
	if job.converted > 0 {
		return fmt.Sprintf("remuxed %d to mp4 · %d failed", job.converted, job.failed)
	}
	if job.lastErr != "" {
		return "remux failed: " + truncate(job.lastErr, 48)
	}
	return ""
}

func (a *App) conversionsActive() bool {
	for _, job := range a.convertJobs {
		if job != nil && !job.done {
			return true
		}
	}
	return false
}

func (a *App) queueRemuxForCompleted() tea.Cmd {
	if !a.remuxEnabled() || !convert.Available() {
		return nil
	}
	a.initConvertTracker()
	var cmds []tea.Cmd
	for _, s := range a.downloads.snaps {
		if s.State != engine.StateDone && s.State != engine.StateSeeding {
			continue
		}
		if cmd := a.startRemuxIfNeeded(s.Magnet, s.Hash, s.DownloadDir, s.DataPath); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	for i := range a.st.Entries {
		e := &a.st.Entries[i]
		if !e.Done {
			continue
		}
		if cmd := a.startRemuxIfNeeded(e.Magnet, metainfo.Hash{}, e.DownloadDir, e.DataPath); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (a *App) startRemuxIfNeeded(magnet string, h metainfo.Hash, downloadDir, dataPath string) tea.Cmd {
	if magnet == "" {
		return nil
	}
	a.initConvertTracker()
	if job := a.convertJobs[magnet]; job != nil {
		return nil
	}
	paths := a.mkvPathsForRemux(h, downloadDir, dataPath)
	if len(paths) == 0 {
		return nil
	}
	a.convertJobs[magnet] = &convertJob{magnet: magnet, total: len(paths)}
	return a.runRemux(magnet, paths)
}

func (a *App) mkvPathsForRemux(h metainfo.Hash, downloadDir, dataPath string) []string {
	if h != (metainfo.Hash{}) {
		if paths := a.eng.MKVPaths(h); len(paths) > 0 {
			return paths
		}
	}
	root := dataPath
	if root == "" {
		root = downloadDir
	}
	paths, err := convert.FindMKV(root)
	if err != nil {
		return nil
	}
	return paths
}

func (a *App) runRemux(magnet string, paths []string) tea.Cmd {
	removeSource := true
	if a.cfg != nil {
		removeSource = a.cfg.RemoveMKVAfterRemux
	}
	opts := convert.Options{RemoveSource: removeSource}
	return func() tea.Msg {
		ctx := context.Background()
		converted, failed := 0, 0
		var lastErr error
		for _, path := range paths {
			if _, err := convert.MKVToMP4(ctx, path, opts); err != nil {
				failed++
				lastErr = err
				continue
			}
			converted++
		}
		return convertDoneMsg{magnet: magnet, converted: converted, failed: failed, err: lastErr}
	}
}

func (a *App) onConvertDone(msg convertDoneMsg) tea.Cmd {
	job := a.convertJobs[msg.magnet]
	if job != nil {
		job.done = true
		job.converted = msg.converted
		job.failed = msg.failed
		if msg.err != nil {
			job.lastErr = msg.err.Error()
		}
	}
	if msg.failed > 0 && msg.err != nil {
		if msg.converted == 0 {
			a.errText = "remux failed: " + msg.err.Error()
		} else {
			a.errText = fmt.Sprintf("remux partial: %d ok, %d failed (%s)", msg.converted, msg.failed, truncate(msg.err.Error(), 40))
		}
		return clearErrCmd()
	}
	return nil
}

func (a *App) previewRemuxHint(p *previewModel) string {
	if !a.remuxEnabled() || !convert.Available() {
		return ""
	}
	n := 0
	for _, f := range p.files {
		if p.excluded[f.Index] {
			continue
		}
		if convert.IsMKV(f.Path) {
			n++
		}
	}
	if n == 0 {
		return ""
	}
	if n == 1 {
		return styleFaint.Render(" · mkv will remux to mp4 after download")
	}
	return styleFaint.Render(fmt.Sprintf(" · %d mkv files will remux to mp4 after download", n))
}
