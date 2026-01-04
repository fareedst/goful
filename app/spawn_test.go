package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/util"
)

func TestExpandMacro(t *testing.T) {
	tmpHome := t.TempDir()

	prevHome, hadHome := os.LookupEnv("HOME")
	if err := os.Setenv("HOME", tmpHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if hadHome {
			_ = os.Setenv("HOME", prevHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
	}()

	g := NewGoful("")
	g.Workspace().ReloadAll() // in home directory

	home, _ := os.UserHomeDir()
	if got := g.File().Name(); got != ".." {
		t.Fatalf("unexpected cursor file %q in temp home %q", got, home)
	}
	macros := []struct {
		in  string
		out string
	}{
		{`%f`, `".."`},
		{`%F`, fmt.Sprintf(`"%s"`, filepath.Dir(home))},
		{`%x`, `"."`},
		{`%X`, fmt.Sprintf(`"%s"`, filepath.Dir(home))},
		{`%m`, `".."`},
		{`%M`, fmt.Sprintf(`"%s"`, filepath.Dir(home))},
		{`%d`, fmt.Sprintf(`"%s"`, filepath.Base(home))},
		{`%D`, fmt.Sprintf(`"%s"`, home)},
		{`%d2`, fmt.Sprintf(`"%s"`, filepath.Base(home))},
		{`%D2`, fmt.Sprintf(`"%s"`, home)},
		{`%~f`, ".."},
		{`%~F`, filepath.Dir(home)},
		{`%~x`, "."},
		{`%~X`, filepath.Dir(home)},
		{`%~m`, ".."},
		{`%~M`, filepath.Dir(home)},
		{`%~d`, filepath.Base(home)},
		{`%~D`, home},
		{`%~d2`, filepath.Base(home)},
		{`%~D2`, home},
		{`%%%f`, `%%".."`},
		{`%%%~f`, `%%..`},
		{`%~~f`, `%~~f`},
		{`\%f%f`, `%f".."`},
		{`\%~f%~f`, `%~f..`},
		{`%\f%f`, `%f".."`},
		{`%\~f%~f`, `%~f..`},
		{"%AA%ff", `%AA".."f`},
		{"%~A~A%~ff", `%~A~A..f`},
		{"%m %f", `".." ".."`},
		{"%~f %f %~m", `.. ".." ..`},
	}

	for _, macro := range macros {
		ret, _ := g.expandMacro(macro.in)
		if ret != macro.out {
			t.Errorf("%s -> %s result %s\n", macro.in, macro.out, ret)
		}
	}
}

func TestOtherWindowDirPaths_REQ_WINDOW_MACRO_ENUMERATION(t *testing.T) {
	ws := &filer.Workspace{
		Dirs: []*filer.Directory{
			stubDirectory("/one"),
			stubDirectory("/two"),
			stubDirectory("/three"),
			stubDirectory("/four"),
		},
		Focus: 0,
	}

	assertOrder := func(focus int, want []string) {
		ws.Focus = focus
		if got := otherWindowDirPaths(ws); fmt.Sprint(got) != fmt.Sprint(want) {
			t.Fatalf("focus %d: got %v want %v", focus, got, want)
		}
	}

	assertOrder(0, []string{"/two", "/three", "/four"})
	assertOrder(2, []string{"/four", "/one", "/two"})

	ws.Dirs = ws.Dirs[:1]
	ws.Focus = 0
	if got := otherWindowDirPaths(ws); len(got) != 0 {
		t.Fatalf("expected empty slice for single window, got %v", got)
	}
}

func TestFormatDirListForMacro_REQ_WINDOW_MACRO_ENUMERATION(t *testing.T) {
	paths := []string{"/one", "/path with space"}

	gotQuoted := formatDirListForMacro(paths, true)
	wantQuoted := fmt.Sprintf("%s %s", util.Quote(paths[0]), util.Quote(paths[1]))
	if gotQuoted != wantQuoted {
		t.Fatalf("quoted format mismatch: got %q want %q", gotQuoted, wantQuoted)
	}

	gotRaw := formatDirListForMacro(paths, false)
	wantRaw := "/one /path with space"
	if gotRaw != wantRaw {
		t.Fatalf("raw format mismatch: got %q want %q", gotRaw, wantRaw)
	}

	if got := formatDirListForMacro(nil, true); got != "" {
		t.Fatalf("expected empty string for nil input, got %q", got)
	}
}

func TestExpandMacroWindowEnumeration_REQ_WINDOW_MACRO_ENUMERATION(t *testing.T) {
	g := NewGoful("")
	ws := g.Workspace()
	ws.Dirs = []*filer.Directory{
		stubDirectory("/alpha"),
		stubDirectory("/beta"),
		stubDirectory("/gamma"),
	}

	ws.Focus = 1 // current: /beta

	got, _ := g.expandMacro("echo %D %D@")
	want := fmt.Sprintf("echo %s %s %s",
		util.Quote("/beta"),
		util.Quote("/gamma"),
		util.Quote("/alpha"),
	)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}

	raw, _ := g.expandMacro("echo %~D@")
	wantRaw := fmt.Sprintf("echo %s %s", util.Quote("/gamma"), util.Quote("/alpha"))
	if raw != wantRaw {
		t.Fatalf("expected escaped list %q, got %q", wantRaw, raw)
	}

	ws.Dirs = ws.Dirs[:1]
	ws.Focus = 0
	empty, _ := g.expandMacro("echo %D@")
	if empty != "echo " {
		t.Fatalf("expected empty expansion for single window, got %q", empty)
	}
}

func stubDirectory(path string) *filer.Directory {
	return &filer.Directory{Path: path}
}
