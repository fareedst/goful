package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fareedst/goful/filer"
	"github.com/fareedst/goful/message"
	"github.com/fareedst/goful/util"
)

const startupDebugPrefix = "DEBUG: [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]"

// ParseStartupDirs normalizes positional CLI arguments into absolute directory paths while collecting warnings.
// [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]
func ParseStartupDirs(args []string) ([]string, []string) {
	var dirs []string
	var warnings []string

	for _, raw := range args {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}

		expanded := util.ExpandPath(trimmed)
		absPath, err := filepath.Abs(expanded)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s: %v", trimmed, err))
			continue
		}

		info, statErr := os.Stat(absPath)
		if statErr != nil {
			warnings = append(warnings, fmt.Sprintf("%s: %v", absPath, statErr))
			continue
		}
		if !info.IsDir() {
			warnings = append(warnings, fmt.Sprintf("%s is not a directory", absPath))
			continue
		}

		dirs = append(dirs, filepath.Clean(absPath))
	}

	return dirs, warnings
}

// SeedStartupWorkspaces resizes the active workspace's directory windows so each entry maps to a positional CLI argument.
// Returns true when seeding occurs so callers can detect fallback scenarios.
// [IMPL:WORKSPACE_START_DIRS] [ARCH:WORKSPACE_BOOTSTRAP] [REQ:WORKSPACE_START_DIRS]
func SeedStartupWorkspaces(g *Goful, dirs []string, debug bool) bool {
	if len(dirs) == 0 {
		if debug {
			message.Infof("%s no startup directories provided; using persisted workspace state", startupDebugPrefix)
		}
		return false
	}

	ws := g.Workspace()
	existing := len(ws.Dirs)
	nextDirs := make([]*filer.Directory, len(dirs))

	for idx, path := range dirs {
		var directory *filer.Directory
		if idx < existing {
			directory = ws.Dirs[idx]
		} else {
			directory = filer.NewDirectory(0, 0, 0, 0)
		}
		directory.Chdir(path)
		nextDirs[idx] = directory
		if debug {
			message.Infof("%s window=%d path=%s", startupDebugPrefix, idx+1, path)
		}
	}

	ws.Dirs = nextDirs
	ws.SetFocus(0)
	ws.ReloadAll()

	width, height := ws.Width(), ws.Height()
	x, y := ws.LeftTop()
	ws.Resize(x, y, width, height)

	return true
}
