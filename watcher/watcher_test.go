package watcher

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatcherDetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("key: value1"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := New(path, 100*time.Millisecond)
	var called atomic.Int32
	w.OnChange(func() {
		called.Add(1)
	})

	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	time.Sleep(200 * time.Millisecond)

	if err := os.WriteFile(path, []byte("key: value2"), 0o644); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	if called.Load() < 1 {
		t.Error("expected callback to be called at least once")
	}
}

func TestWatcherDebounce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("key: value1"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := New(path, 300*time.Millisecond)
	var called atomic.Int32
	w.OnChange(func() {
		called.Add(1)
	})

	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	time.Sleep(200 * time.Millisecond)

	// Rapid writes should be debounced
	for i := range 5 {
		if err := os.WriteFile(path, []byte("key: value"+string(rune('0'+i))), 0o644); err != nil {
			t.Fatal(err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(600 * time.Millisecond)

	c := called.Load()
	if c > 2 {
		t.Errorf("expected debounce to limit calls, got %d", c)
	}
}
