// Package app provides synchronized command operations across all panes.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anmitsu/goful/cmdline"
	"github.com/anmitsu/goful/message"
)

// SyncResult holds the result of a sync operation.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
type SyncResult struct {
	Succeeded int
	Skipped   int
	Failures  []SyncFailure
}

// SyncFailure records a failure in a specific pane.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
type SyncFailure struct {
	PaneIndex int
	Error     error
}

// SyncMode starts the sync command mode for synchronized operations across all panes.
// If ignoreFailures is true, operations continue even if some panes fail.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) SyncMode(ignoreFailures bool) {
	g.next = cmdline.New(&syncMode{g, ignoreFailures}, g)
}

// syncMode is the prefix mode that waits for an operation key.
// Press '!' to toggle ignore-failures mode, then c/d/r to execute.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
type syncMode struct {
	*Goful
	ignoreFailures bool
}

func (m *syncMode) String() string { return "sync" }

func (m *syncMode) Prompt() string {
	if m.ignoreFailures {
		return "Sync! [c]opy [d]elete [r]ename (ignore failures): "
	}
	return "Sync [c]opy [d]elete [r]ename [!]ignore: "
}

func (m *syncMode) Draw(c *cmdline.Cmdline) { c.DrawLine() }

func (m *syncMode) Run(c *cmdline.Cmdline) {
	key := c.String()

	// Handle '!' toggle for ignore-failures mode
	if key == "!" {
		m.ignoreFailures = !m.ignoreFailures
		c.SetText("")
		return // Don't exit, wait for operation key
	}

	c.Exit()

	// Get the filename from cursor in the focused pane
	file := m.File()
	if file == nil || file.Name() == ".." {
		message.Errorf("[REQ:SYNC_COMMANDS] no file selected")
		return
	}
	filename := file.Name()

	switch key {
	case "c":
		m.StartSyncCopy(filename, m.ignoreFailures)
	case "d":
		m.StartSyncDelete(filename, m.ignoreFailures)
	case "r":
		m.StartSyncRename(filename, m.ignoreFailures)
	default:
		// Any other key exits the mode silently
	}
}

// StartSyncCopy initiates the sync copy operation.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS] [IMPL:TOOLBAR_SYNC_COPY] [REQ:TOOLBAR_SYNC_BUTTONS]
func (g *Goful) StartSyncCopy(filename string, ignoreFailures bool) {
	c := cmdline.New(&syncCopyMode{g, filename, ignoreFailures}, g)
	// Default to the same filename - user must change it to a different name
	c.SetText(filename)
	g.next = c
}

// StartSyncDelete initiates the sync delete operation.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS] [IMPL:TOOLBAR_SYNC_DELETE] [REQ:TOOLBAR_SYNC_BUTTONS]
func (g *Goful) StartSyncDelete(filename string, ignoreFailures bool) {
	g.next = cmdline.New(&syncDeleteMode{g, filename, ignoreFailures}, g)
}

// StartSyncRename initiates the sync rename operation.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS] [IMPL:TOOLBAR_SYNC_RENAME] [REQ:TOOLBAR_SYNC_BUTTONS]
func (g *Goful) StartSyncRename(filename string, ignoreFailures bool) {
	c := cmdline.New(&syncRenameMode{g, filename, ignoreFailures}, g)
	c.SetText(filename)
	c.MoveCursor(-len(filepath.Ext(filename)))
	g.next = c
}

// executeSyncCopy executes copy for a file across all panes.
// Copies the source file to a new filename in each pane's directory.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) executeSyncCopy(filename, newName string, ignoreFailures bool) {
	ws := g.Workspace()
	result := SyncResult{}

	for i := 0; i < len(ws.Dirs); i++ {
		idx := (ws.Focus + i) % len(ws.Dirs)
		dir := ws.Dirs[idx]

		file := dir.FindFileByName(filename)
		if file == nil {
			result.Skipped++
			continue
		}

		// Execute copy to new filename in the same directory
		srcPath := file.Path()
		dstPath := filepath.Join(dir.Path, newName)

		err := g.copyFileForSync(srcPath, dstPath)
		if err != nil {
			result.Failures = append(result.Failures, SyncFailure{idx, err})
			if !ignoreFailures {
				g.reportSyncResult("copy", filename, result)
				return
			}
		} else {
			result.Succeeded++
		}
	}

	g.reportSyncResult("copy", filename, result)
	g.Workspace().ReloadAll()
}

// executeSyncDelete executes delete for a file across all panes.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) executeSyncDelete(filename string, ignoreFailures bool) {
	ws := g.Workspace()
	result := SyncResult{}

	for i := 0; i < len(ws.Dirs); i++ {
		idx := (ws.Focus + i) % len(ws.Dirs)
		dir := ws.Dirs[idx]

		file := dir.FindFileByName(filename)
		if file == nil {
			result.Skipped++
			continue
		}

		// Execute delete for this file
		err := os.RemoveAll(file.Path())
		if err != nil {
			result.Failures = append(result.Failures, SyncFailure{idx, err})
			if !ignoreFailures {
				g.reportSyncResult("delete", filename, result)
				g.Workspace().ReloadAll()
				return
			}
		} else {
			result.Succeeded++
		}
	}

	g.reportSyncResult("delete", filename, result)
	g.Workspace().ReloadAll()
}

// executeSyncRename executes rename for a file across all panes.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) executeSyncRename(oldName, newName string, ignoreFailures bool) {
	ws := g.Workspace()
	result := SyncResult{}

	for i := 0; i < len(ws.Dirs); i++ {
		idx := (ws.Focus + i) % len(ws.Dirs)
		dir := ws.Dirs[idx]

		file := dir.FindFileByName(oldName)
		if file == nil {
			result.Skipped++
			continue
		}

		// Execute rename for this file
		oldPath := file.Path()
		newPath := filepath.Join(dir.Path, newName)

		err := os.Rename(oldPath, newPath)
		if err != nil {
			result.Failures = append(result.Failures, SyncFailure{idx, err})
			if !ignoreFailures {
				g.reportSyncResult("rename", oldName, result)
				g.Workspace().ReloadAll()
				return
			}
		} else {
			result.Succeeded++
		}
	}

	g.reportSyncResult("rename", oldName, result)
	g.Workspace().ReloadAll()
}

// copyFileForSync copies a single file for sync operations.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) copyFileForSync(src, dst string) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Check if destination exists
	if _, err := os.Lstat(dst); err == nil {
		// Destination exists - overwrite without asking in sync mode
		// since user already confirmed the operation
	}

	if srcInfo.IsDir() {
		return copyDirRecursive(src, dst)
	}
	return copyFileSimple(src, dst)
}

// copyDirRecursive copies a directory recursively.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func copyDirRecursive(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDirRecursive(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFileSimple(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFileSimple copies a single file.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func copyFileSimple(src, dst string) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Handle symlinks
	if srcInfo.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(link, dst)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 4096)
	for {
		n, err := srcFile.Read(buf)
		if n > 0 {
			if _, wErr := dstFile.Write(buf[:n]); wErr != nil {
				return wErr
			}
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
	}

	return nil
}

// reportSyncResult displays the result of a sync operation.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) reportSyncResult(op, filename string, result SyncResult) {
	if len(result.Failures) == 0 {
		if result.Skipped > 0 {
			message.Infof("[REQ:SYNC_COMMANDS] %s '%s': %d succeeded, %d skipped (file not found)",
				op, filename, result.Succeeded, result.Skipped)
		} else {
			message.Infof("[REQ:SYNC_COMMANDS] %s '%s': %d succeeded",
				op, filename, result.Succeeded)
		}
	} else {
		// Report failures
		failedPanes := make([]int, len(result.Failures))
		for i, f := range result.Failures {
			failedPanes[i] = f.PaneIndex + 1 // 1-based for user display
		}
		message.Errorf("[REQ:SYNC_COMMANDS] %s '%s': %d succeeded, %d failed (panes %v), %d skipped",
			op, filename, result.Succeeded, len(result.Failures), failedPanes, result.Skipped)
	}
}

// getPaneCount returns the number of panes for display in prompts.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) getPaneCount() int {
	return len(g.Workspace().Dirs)
}

// countFilesNamed counts how many panes have a file with the given name.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) countFilesNamed(filename string) int {
	count := 0
	for _, dir := range g.Workspace().Dirs {
		if dir.FindFileByName(filename) != nil {
			count++
		}
	}
	return count
}

// panesWithFile returns a string describing which panes have the file.
// [IMPL:SYNC_EXECUTE] [ARCH:SYNC_MODE] [REQ:SYNC_COMMANDS]
func (g *Goful) panesWithFile(filename string) string {
	count := g.countFilesNamed(filename)
	total := g.getPaneCount()
	if count == total {
		return fmt.Sprintf("all %d panes", total)
	}
	return fmt.Sprintf("%d of %d panes", count, total)
}
