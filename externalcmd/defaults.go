package externalcmd

import "runtime"

var posixDefaults = []Entry{
	{Key: "c", Label: "copy %m to %D2    ", Command: "cp -vai %m %D2"},
	{Key: "m", Label: "move %m to %D2    ", Command: "mv -vi %m %D2"},
	{Key: "D", Label: "remove %m files   ", Command: "rm -vR %m"},
	{Key: "k", Label: "make directory    ", Command: "mkdir -vp ./"},
	{Key: "n", Label: "create newfile    ", Command: "touch ./"},
	{Key: "T", Label: "time copy %f to %m", Command: "touch -r %f %m"},
	{Key: "M", Label: "change mode %m    ", Command: "chmod 644 %m", Offset: -3},
	{Key: "r", Label: "move (rename) %f  ", Command: "mv -vi %f %~f"},
	{Key: "R", Label: "bulk rename %m    ", Command: `rename -v "s///" %m`, Offset: -6},
	{Key: "f", Label: "find . -name      ", Command: `find . -name "*"`, Offset: -1},
	{Key: "A", Label: "archives menu     ", RunMenu: "archive"},
}

var windowsDefaults = []Entry{
	{Key: "c", Label: "copy %~f to %~D2 ", Command: "robocopy /e %~f %~D2"},
	{Key: "m", Label: "move %~f to %~D2 ", Command: "move /-y %~f %~D2"},
	{Key: "d", Label: "del /s %~m       ", Command: "del /s %~m"},
	{Key: "D", Label: "rd /s /q %~m     ", Command: "rd /s /q %~m"},
	{Key: "k", Label: "make directory   ", Command: "mkdir "},
	{Key: "n", Label: "create newfile   ", Command: "copy nul "},
	{Key: "r", Label: "move (rename) %f ", Command: "move /-y %~f ./"},
	{Key: "w", Label: "where . *        ", Command: "where . *"},
}

// Defaults returns a deep copy of the built-in entries for the current GOOS.
// [IMPL:EXTERNAL_COMMAND_LOADER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
func Defaults(goos string) []Entry {
	if goos == "" {
		goos = runtime.GOOS
	}

	if goos == "windows" {
		return cloneEntries(windowsDefaults)
	}
	return cloneEntries(posixDefaults)
}

func cloneEntries(in []Entry) []Entry {
	out := make([]Entry, len(in))
	for i, entry := range in {
		entry.Menu = MenuName
		// copy slice to avoid sharing backing arrays
		if len(entry.Platforms) > 0 {
			entry.Platforms = append([]string(nil), entry.Platforms...)
		}
		out[i] = entry
	}
	return out
}
