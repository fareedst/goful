package externalmenu

import (
	"testing"

	"github.com/anmitsu/goful/externalcmd"
)

func TestEnsureMenuSpecsAddsPlaceholder_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_BINDER]
	specs := ensureMenuSpecs(nil)
	if len(specs) != 1 || specs[0].Placeholder == false {
		t.Fatalf("expected placeholder spec, got %+v", specs)
	}
	if specs[0].Label != placeholderExternalCommandLabel {
		t.Fatalf("unexpected placeholder label: %q", specs[0].Label)
	}
}

func TestBuildMenuArgsInvokesShellWithOffset_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_BINDER]
	entries := []externalcmd.Entry{
		{Key: "c", Label: "copy", Command: "cp -vai %m %D2", Offset: -2},
	}
	specs := ensureMenuSpecs(buildMenuSpecs(entries))

	var (
		calledCmd    string
		calledOffset int
	)
	args := buildMenuArgs(specs, func(cmd string, offset ...int) {
		calledCmd = cmd
		if len(offset) > 0 {
			calledOffset = offset[0]
		}
	}, nil)

	menuArgs := args[externalcmd.MenuName]
	if len(menuArgs) != 3 {
		t.Fatalf("expected 3 menu args, got %d", len(menuArgs))
	}
	callback, ok := menuArgs[2].(func())
	if !ok {
		t.Fatalf("expected callback func, got %T", menuArgs[2])
	}
	callback()

	if calledCmd != "cp -vai %m %D2" || calledOffset != -2 {
		t.Fatalf("shell invoker mismatch: cmd=%q offset=%d", calledCmd, calledOffset)
	}
}

func TestBuildMenuArgsInvokesMenuLauncher_REQ_EXTERNAL_COMMAND_CONFIG(t *testing.T) {
	// [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:EXTERNAL_COMMAND_REGISTRY] [IMPL:EXTERNAL_COMMAND_BINDER]
	entries := []externalcmd.Entry{
		{Key: "A", Label: "archives", RunMenu: "archive"},
	}
	specs := ensureMenuSpecs(buildMenuSpecs(entries))
	args := buildMenuArgs(specs, nil, func(name string) {
		if name != "archive" {
			t.Fatalf("unexpected menu name %q", name)
		}
	})
	menuArgs := args[externalcmd.MenuName]
	if len(menuArgs) != 3 {
		t.Fatalf("expected 3 values for menu entry, got %d", len(menuArgs))
	}
	callback := menuArgs[2].(func())
	callback() // should invoke menu opener without panic
}
