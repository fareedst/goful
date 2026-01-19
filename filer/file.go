package filer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fareedst/goful/look"
	"github.com/fareedst/goful/message"
	"github.com/fareedst/goful/util"
	"github.com/fareedst/goful/widget"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

var statView = fileStatView{true, true, true}

type fileStatView struct {
	size       bool
	permission bool
	time       bool
}

// SetStatView sets the file state view.
func SetStatView(size, permission, time bool) { statView = fileStatView{size, permission, time} }

// ToggleSizeView toggles the file size view.
func ToggleSizeView() { statView.size = !statView.size }

// TogglePermView toggles the file permission view.
func TogglePermView() { statView.permission = !statView.permission }

// ToggleTimeView toggles the file time view.
func ToggleTimeView() { statView.time = !statView.time }

var timeFormat = "06-01-02 15:04"

// SetTimeFormat sets the time format of files.
func SetTimeFormat(format string) {
	timeFormat = format
}

// FileStat is file information.
type FileStat struct {
	os.FileInfo             // os.Lstat(path)
	stat        os.FileInfo // os.Stat(path)
	path        string      // full path of file
	name        string      // base name of path or ".." as upper directory
	display     string      // display name for draw
	marked      bool        // marked whether
}

// NewFileStat creates a new file stat of the file in the directory.
func NewFileStat(dir string, name string) *FileStat {
	path := filepath.Join(dir, name)

	lstat, err := os.Lstat(path)
	if err != nil {
		message.Error(err)
		return nil
	}
	stat, err := os.Stat(path)
	if err != nil {
		stat = lstat
	}

	var display string
	if stat.IsDir() {
		display = name
	} else {
		display = util.RemoveExt(name)
	}

	return &FileStat{
		FileInfo: lstat,
		stat:     stat,
		path:     path,
		name:     name,
		display:  display,
		marked:   false,
	}
}

// Name returns the file name.
func (f *FileStat) Name() string {
	return f.name
}

// SetDisplay sets the display name for drawing.
func (f *FileStat) SetDisplay(name string) {
	f.display = name
}

// ResetDisplay resets the display name to the file name.
func (f *FileStat) ResetDisplay() {
	if f.stat.IsDir() {
		f.display = f.name
	} else {
		f.display = util.RemoveExt(f.name)
	}
}

// Mark the file.
func (f *FileStat) Mark() {
	f.marked = true
}

// Markoff the file.
func (f *FileStat) Markoff() {
	f.marked = false
}

// ToggleMark toggles the file mark.
func (f *FileStat) ToggleMark() {
	f.marked = !f.marked
}

// Path returns the file path.
func (f *FileStat) Path() string {
	return f.path
}

// Ext retruns the file extension.
func (f *FileStat) Ext() string {
	if f.stat.IsDir() {
		return ""
	}
	if ext := filepath.Ext(f.Name()); ext != f.Name() {
		return ext
	}
	return ""
}

// IsLink reports whether the symlink.
func (f *FileStat) IsLink() bool {
	return f.Mode()&os.ModeSymlink != 0
}

// IsExec reports whether the executable file.
func (f *FileStat) IsExec() bool {
	return f.stat.Mode().Perm()&0111 != 0
}

// IsFIFO reports whether the named pipe file.
func (f *FileStat) IsFIFO() bool {
	return f.stat.Mode()&os.ModeNamedPipe != 0
}

// IsDevice reports whether the device file.
func (f *FileStat) IsDevice() bool {
	return f.stat.Mode()&os.ModeDevice != 0
}

// IsCharDevice reports whether the character device file.
func (f *FileStat) IsCharDevice() bool {
	return f.stat.Mode()&os.ModeCharDevice != 0
}

// IsSocket reports whether the socket file.
func (f *FileStat) IsSocket() bool {
	return f.stat.Mode()&os.ModeSocket != 0
}

// IsMarked reports whether the marked file.
func (f *FileStat) IsMarked() bool {
	return f.marked
}

func (f *FileStat) suffix() string {
	if f.IsLink() {
		link, _ := os.Readlink(f.Path())
		if f.stat.IsDir() {
			return "@ -> " + link + "/"
		}
		return "@ -> " + link
	} else if f.IsDir() {
		return "/"
	} else if f.IsFIFO() {
		return "|"
	} else if f.IsSocket() {
		return "="
	} else if f.IsExec() {
		return "*"
	}
	return ""
}

func (f *FileStat) states() string {
	ret := f.Ext()
	if statView.size {
		if f.stat.IsDir() {
			ret += fmt.Sprintf("%8s", "<DIR>")
		} else {
			ret += fmt.Sprintf("%8s", util.FormatSize(f.stat.Size()))
		}
	}
	if statView.permission {
		ret += " " + f.stat.Mode().String()
	}
	if statView.time {
		ret += " " + f.stat.ModTime().Format(timeFormat)
	}
	return ret
}

func (f *FileStat) look() tcell.Style {
	switch {
	case f.IsMarked():
		return look.Marked()
	case f.IsLink():
		if f.stat.IsDir() {
			return look.SymlinkDir()
		}
		return look.Symlink()
	case f.IsDir():
		return look.Directory()
	case f.IsExec():
		return look.Executable()
	default:
		return look.Default()
	}
}

// Draw the file name and file stats.
func (f *FileStat) Draw(x, y, width int, focus bool) {
	f.DrawWithComparison(x, y, width, focus, nil)
}

// DrawWithComparison draws the file with optional comparison coloring.
// [IMPL:COMPARISON_DRAW] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (f *FileStat) DrawWithComparison(x, y, width int, focus bool, cmp *CompareState) {
	baseStyle := f.look()
	if focus {
		baseStyle = baseStyle.Reverse(true)
	}

	// Determine if comparison colors should be applied
	useCompare := cmp != nil && cmp.NamePresent && look.ComparisonEnabled()

	// Calculate widths for each section
	ext := f.Ext()
	var sizeStr, timeStr string
	if statView.size {
		if f.stat.IsDir() {
			sizeStr = fmt.Sprintf("%8s", "<DIR>")
		} else {
			sizeStr = fmt.Sprintf("%8s", util.FormatSize(f.stat.Size()))
		}
	}
	permStr := ""
	if statView.permission {
		permStr = " " + f.stat.Mode().String()
	}
	if statView.time {
		timeStr = " " + f.stat.ModTime().Format(timeFormat)
	}

	// Build states string for width calculation
	states := ext + sizeStr + permStr + timeStr
	nameWidth := width - len(states)

	// Draw file name with prefix
	pre := " "
	if f.marked {
		pre = "*"
	}
	nameDisplay := pre + f.display + f.suffix()
	nameDisplay = runewidth.Truncate(nameDisplay, nameWidth, "~")
	nameDisplay = runewidth.FillRight(nameDisplay, nameWidth)

	// Apply name comparison style if applicable
	nameStyle := baseStyle
	if useCompare {
		nameStyle = look.CompareNamePresent()
		if focus {
			nameStyle = nameStyle.Reverse(true)
		}
	}
	x = widget.SetCells(x, y, nameDisplay, nameStyle)

	// Draw extension with base style
	if ext != "" {
		x = widget.SetCells(x, y, ext, baseStyle)
	}

	// Draw size with comparison style
	// [IMPL:DIGEST_COMPARISON] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
	if sizeStr != "" {
		sizeStyle := baseStyle
		if useCompare {
			switch cmp.SizeState {
			case SizeEqual:
				sizeStyle = look.CompareSizeEqual()
			case SizeSmallest:
				sizeStyle = look.CompareSizeSmallest()
			case SizeLargest:
				sizeStyle = look.CompareSizeLargest()
			}
			// Apply digest comparison attributes
			switch cmp.DigestState {
			case DigestEqual:
				sizeStyle = sizeStyle.Underline(true)
			case DigestDifferent:
				sizeStyle = sizeStyle.StrikeThrough(true)
			}
			if focus {
				sizeStyle = sizeStyle.Reverse(true)
			}
		}
		x = widget.SetCells(x, y, sizeStr, sizeStyle)
	}

	// Draw permission with base style
	if permStr != "" {
		x = widget.SetCells(x, y, permStr, baseStyle)
	}

	// Draw time with comparison style
	if timeStr != "" {
		timeStyle := baseStyle
		if useCompare {
			switch cmp.TimeState {
			case TimeEqual:
				timeStyle = look.CompareTimeEqual()
			case TimeEarliest:
				timeStyle = look.CompareTimeEarliest()
			case TimeLatest:
				timeStyle = look.CompareTimeLatest()
			}
			if focus {
				timeStyle = timeStyle.Reverse(true)
			}
		}
		widget.SetCells(x, y, timeStr, timeStyle)
	}
}
