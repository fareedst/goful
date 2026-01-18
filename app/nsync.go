//go:build darwin
// +build darwin

// Package app provides nsync SDK integration for multi-target copy/move operations.
// [IMPL:NSYNC_OBSERVER] [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/progress"
	"github.com/anmitsu/goful/widget"
	"github.com/fareedst/nsync/pkg/nsync"
)

// gofulObserver implements nsync.Observer to bridge nsync progress events
// to goful's progress widget and message system.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
type gofulObserver struct {
	mu           sync.Mutex
	lastBytes    int64
	totalBytes   int64
	totalItems   int
	startTime    time.Time
	destinations int
}

// newGofulObserver creates a new observer for bridging nsync to goful progress.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func newGofulObserver() *gofulObserver {
	return &gofulObserver{
		startTime: time.Now(),
	}
}

// OnStart is called when sync begins with the plan of work.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (o *gofulObserver) OnStart(plan nsync.Plan) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.totalBytes = plan.TotalBytes
	o.totalItems = plan.TotalItems
	o.destinations = plan.TotalDestinations
	o.startTime = time.Now()

	// Initialize progress display
	progress.Start(float64(plan.TotalBytes))
	progress.StartTaskCount(plan.TotalItems)

	message.Infof("Syncing %d items (%s) to %d destinations",
		plan.TotalItems, formatBytes(plan.TotalBytes), plan.TotalDestinations)
}

// OnItemStart is called when a source item begins processing.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (o *gofulObserver) OnItemStart(item nsync.ItemInfo) {
	// Create a fake FileInfo for the progress display
	fi := &fakeFileInfo{
		name: filepath.Base(item.SourcePath),
		size: item.Size,
	}
	progress.StartTask(fi)
}

// OnItemProgress is called with streaming byte updates during copy.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (o *gofulObserver) OnItemProgress(item nsync.ItemInfo, bytesCopied int64) {
	o.mu.Lock()
	delta := bytesCopied - o.lastBytes
	o.lastBytes = bytesCopied
	o.mu.Unlock()

	if delta > 0 {
		progress.Update(float64(delta))
	}
}

// OnItemComplete is called when a source item finishes processing.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (o *gofulObserver) OnItemComplete(item nsync.ItemInfo, result nsync.ItemResult) {
	progress.FinishTask()

	// Report per-item errors
	if result.Error != nil {
		message.Errorf("Failed: %s: %v", filepath.Base(item.SourcePath), result.Error)
	}

	// Report per-destination errors
	for _, destResult := range result.DestResults {
		if destResult.Error != nil {
			message.Errorf("Failed to %s: %v", filepath.Base(destResult.DestPath), destResult.Error)
		}
	}

	// Reset per-item byte counter for next item
	o.mu.Lock()
	o.lastBytes = 0
	o.mu.Unlock()
}

// OnProgress is called periodically with overall sync statistics.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (o *gofulObserver) OnProgress(stats nsync.Stats) {
	// Trigger a UI refresh
	progress.Draw()
	widget.Show()
}

// OnFinish is called when sync completes with final results.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (o *gofulObserver) OnFinish(result nsync.Result) {
	progress.Finish()
}

// fakeFileInfo implements os.FileInfo for progress display.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
type fakeFileInfo struct {
	name string
	size int64
}

func (f *fakeFileInfo) Name() string       { return f.name }
func (f *fakeFileInfo) Size() int64        { return f.size }
func (f *fakeFileInfo) Mode() os.FileMode  { return 0 }
func (f *fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f *fakeFileInfo) IsDir() bool        { return false }
func (f *fakeFileInfo) Sys() interface{}   { return nil }

// formatBytes formats bytes as a human-readable string.
// [IMPL:NSYNC_OBSERVER] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// syncCopy performs a multi-destination copy using nsync.
// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (g *Goful) syncCopy(sources []string, destinations []string) {
	srcAbs := make([]string, len(sources))
	for i, src := range sources {
		srcAbs[i], _ = filepath.Abs(src)
	}

	destAbs := make([]string, len(destinations))
	for i, dst := range destinations {
		destAbs[i], _ = filepath.Abs(dst)
	}

	g.asyncFilectrl(func() {
		observer := newGofulObserver()
		cfg := nsync.Config{
			Sources:         srcAbs,
			Destinations:    destAbs,
			Recursive:       true,
			Move:            false,
			Jobs:            4,
			DestParallelism: len(destAbs), // Parallel to all destinations
			Mode:            nsync.ModeCompare,
			CompareMethod:   nsync.CompareSizeMtime, // Fast comparison
		}

		syncer, err := nsync.New(cfg, nsync.WithObserver(observer))
		if err != nil {
			message.Error(err)
			return
		}

		ctx := context.Background()
		result, err := syncer.Sync(ctx)

		if err != nil {
			message.Errorf("Sync failed: %v", err)
		} else if result.Cancelled {
			message.Info("Sync cancelled")
		} else if result.ItemsFailed > 0 {
			message.Errorf("Copied with %d failures: %d items to %d destinations",
				result.ItemsFailed, result.ItemsCompleted, len(destAbs))
		} else {
			message.Infof("Copied %d items (%s) to %d destinations in %s",
				result.ItemsCompleted, formatBytes(result.BytesCopied),
				len(destAbs), result.Duration.Round(time.Millisecond))
		}
	})
}

// syncMove performs a multi-destination move using nsync.
// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
func (g *Goful) syncMove(sources []string, destinations []string) {
	srcAbs := make([]string, len(sources))
	for i, src := range sources {
		srcAbs[i], _ = filepath.Abs(src)
	}

	destAbs := make([]string, len(destinations))
	for i, dst := range destinations {
		destAbs[i], _ = filepath.Abs(dst)
	}

	g.asyncFilectrl(func() {
		observer := newGofulObserver()
		cfg := nsync.Config{
			Sources:         srcAbs,
			Destinations:    destAbs,
			Recursive:       true,
			Move:            true, // Enable move semantics
			Jobs:            4,
			DestParallelism: len(destAbs),
			Mode:            nsync.ModeCompare,
			CompareMethod:   nsync.CompareSizeMtime,
		}

		syncer, err := nsync.New(cfg, nsync.WithObserver(observer))
		if err != nil {
			message.Error(err)
			return
		}

		ctx := context.Background()
		result, err := syncer.Sync(ctx)

		if err != nil {
			message.Errorf("Sync failed: %v", err)
		} else if result.Cancelled {
			message.Info("Sync cancelled")
		} else if result.ItemsFailed > 0 {
			message.Errorf("Moved with %d failures: %d items to %d destinations",
				result.ItemsFailed, result.ItemsCompleted, len(destAbs))
		} else {
			message.Infof("Moved %d items (%s) to %d destinations in %s",
				result.ItemsCompleted, formatBytes(result.BytesCopied),
				len(destAbs), result.Duration.Round(time.Millisecond))
		}
	})
}

// CopyAll prompts for confirmation then copies selected files to all other visible workspace directories.
// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func (g *Goful) CopyAll() {
	destinations := otherWindowDirPaths(g.Workspace())

	if len(destinations) == 0 {
		// Fall back to regular copy when only one pane
		message.Info("Only one pane visible - use regular copy (c)")
		g.Copy()
		return
	}

	// Collect sources from marks or cursor
	var sources []string
	if g.Dir().IsMark() {
		sources = g.Dir().MarkfilePaths()
		g.Dir().MarkClear()
	} else {
		file := g.File()
		if file == nil {
			message.Error(fmt.Errorf("no file selected"))
			return
		}
		sources = []string{file.Path()}
	}

	if len(sources) == 0 {
		message.Error(fmt.Errorf("no files to copy"))
		return
	}

	// Start confirmation mode instead of executing immediately
	// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
	g.next = cmdline.New(&copyAllMode{g, sources, destinations}, g)
}

// doCopyAll executes the multi-target copy after confirmation.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func (g *Goful) doCopyAll(sources, destinations []string) {
	g.syncCopy(sources, destinations)
}

// MoveAll prompts for confirmation then moves selected files to all other visible workspace directories.
// [IMPL:NSYNC_COPY_MOVE] [ARCH:NSYNC_INTEGRATION] [REQ:NSYNC_MULTI_TARGET]
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func (g *Goful) MoveAll() {
	destinations := otherWindowDirPaths(g.Workspace())

	if len(destinations) == 0 {
		// Fall back to regular move when only one pane
		message.Info("Only one pane visible - use regular move (m)")
		g.Move()
		return
	}

	// Collect sources from marks or cursor
	var sources []string
	if g.Dir().IsMark() {
		sources = g.Dir().MarkfilePaths()
		g.Dir().MarkClear()
	} else {
		file := g.File()
		if file == nil {
			message.Error(fmt.Errorf("no file selected"))
			return
		}
		sources = []string{file.Path()}
	}

	if len(sources) == 0 {
		message.Error(fmt.Errorf("no files to move"))
		return
	}

	// Start confirmation mode instead of executing immediately
	// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
	g.next = cmdline.New(&moveAllMode{g, sources, destinations}, g)
}

// doMoveAll executes the multi-target move after confirmation.
// [IMPL:NSYNC_CONFIRMATION] [ARCH:NSYNC_CONFIRMATION] [REQ:NSYNC_CONFIRMATION]
func (g *Goful) doMoveAll(sources, destinations []string) {
	g.syncMove(sources, destinations)
}
