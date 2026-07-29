// Package convert remuxes downloaded video files without re-encoding.
package convert

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrFFmpegMissing = errors.New("ffmpeg not found in PATH")

// Options tune MKV→MP4 remux behavior.
type Options struct {
	RemoveSource bool
}

// Available reports whether ffmpeg is on PATH.
func Available() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// IsMKV reports whether path has a .mkv extension.
func IsMKV(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".mkv")
}

// MP4Path returns the sibling .mp4 path for an input file.
func MP4Path(input string) string {
	ext := filepath.Ext(input)
	if ext == "" {
		return input + ".mp4"
	}
	return strings.TrimSuffix(input, ext) + ".mp4"
}

// NeedsRemux reports whether input should be remuxed to MP4.
func NeedsRemux(input string) bool {
	if !IsMKV(input) {
		return false
	}
	out := MP4Path(input)
	st, err := os.Stat(out)
	if err != nil {
		return true
	}
	in, err := os.Stat(input)
	if err != nil {
		return false
	}
	// Skip when an existing MP4 is at least as new and non-empty.
	return st.Size() == 0 || in.ModTime().After(st.ModTime())
}

// FindMKV walks root (file or directory) and returns absolute MKV paths that
// still need remuxing.
func FindMKV(root string) ([]string, error) {
	if root == "" {
		return nil, nil
	}
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	if !info.IsDir() {
		if NeedsRemux(root) {
			abs, err := filepath.Abs(root)
			if err != nil {
				return nil, err
			}
			return []string{abs}, nil
		}
		return nil, nil
	}
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if NeedsRemux(path) {
			abs, aerr := filepath.Abs(path)
			if aerr != nil {
				return aerr
			}
			out = append(out, abs)
		}
		return nil
	})
	return out, err
}

// MKVToMP4 remuxes input to MP4 using stream copy (no re-encode).
func MKVToMP4(ctx context.Context, input string, opts Options) (string, error) {
	if !IsMKV(input) {
		return "", fmt.Errorf("not an mkv file: %s", input)
	}
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", ErrFFmpegMissing
	}
	in, err := os.Stat(input)
	if err != nil {
		return "", err
	}
	if in.IsDir() {
		return "", fmt.Errorf("input is a directory: %s", input)
	}

	output := MP4Path(input)
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return "", err
	}
	tmp := output + ".pork-remux.part"
	_ = os.Remove(tmp)

	args := []string{
		"-hide_banner", "-loglevel", "error", "-nostdin", "-y",
		"-i", input,
		"-map", "0",
		"-c", "copy",
		"-movflags", "+faststart",
		tmp,
	}
	cmd := exec.CommandContext(ctx, ffmpeg, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.Remove(tmp)
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			return "", fmt.Errorf("ffmpeg remux failed: %w", err)
		}
		return "", fmt.Errorf("ffmpeg remux failed: %s", msg)
	}
	if err := os.Rename(tmp, output); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	if opts.RemoveSource {
		if err := os.Remove(input); err != nil && !os.IsNotExist(err) {
			return output, fmt.Errorf("remux ok but could not remove source: %w", err)
		}
	}
	return output, nil
}
