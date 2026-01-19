package externalmenu

import (
	"strings"

	"github.com/fareedst/goful/app"
	"github.com/fareedst/goful/externalcmd"
	"github.com/fareedst/goful/menu"
	"github.com/fareedst/goful/message"
)

const placeholderExternalCommandLabel = "no external commands configured"

type shellInvoker func(cmd string, offset ...int)
type menuOpener func(name string)

type menuSpec struct {
	Menu        string
	Key         string
	Label       string
	Command     string
	RunMenu     string
	Offset      int
	Placeholder bool
}

// Register wires menu entries produced by the loader into goful.
// [IMPL:EXTERNAL_COMMAND_BINDER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
func Register(g *app.Goful, entries []externalcmd.Entry) {
	specs := ensureMenuSpecs(buildMenuSpecs(entries))
	argsByMenu := buildMenuArgs(specs,
		func(cmd string, offset ...int) { g.Shell(cmd, offset...) },
		func(name string) { g.Menu(name) },
	)
	for name, args := range argsByMenu {
		if len(args) == 0 {
			continue
		}
		menu.Add(name, args...)
	}
}

func buildMenuSpecs(entries []externalcmd.Entry) []menuSpec {
	specs := make([]menuSpec, 0, len(entries))
	for _, entry := range entries {
		menuName := entry.Menu
		if strings.TrimSpace(menuName) == "" {
			menuName = externalcmd.MenuName
		}
		if strings.TrimSpace(entry.Key) == "" || strings.TrimSpace(entry.Label) == "" {
			continue
		}
		specs = append(specs, menuSpec{
			Menu:    menuName,
			Key:     entry.Key,
			Label:   entry.Label,
			Command: entry.Command,
			RunMenu: entry.RunMenu,
			Offset:  entry.Offset,
		})
	}
	return specs
}

func ensureMenuSpecs(specs []menuSpec) []menuSpec {
	if len(specs) > 0 {
		return specs
	}
	return []menuSpec{
		{
			Menu:        externalcmd.MenuName,
			Key:         "-",
			Label:       placeholderExternalCommandLabel,
			Placeholder: true,
		},
	}
}

func buildMenuArgs(specs []menuSpec, shell shellInvoker, openMenu menuOpener) map[string][]interface{} {
	argsByMenu := make(map[string][]interface{})
	for _, spec := range specs {
		callback := makeCommandCallback(spec, shell, openMenu)
		argsByMenu[spec.Menu] = append(argsByMenu[spec.Menu], spec.Key, spec.Label, callback)
	}
	return argsByMenu
}

func makeCommandCallback(spec menuSpec, shell shellInvoker, openMenu menuOpener) func() {
	if spec.Placeholder {
		return func() {
			message.Info("[REQ:EXTERNAL_COMMAND_CONFIG] No external commands configured. Provide -commands, set GOFUL_COMMANDS_FILE, or create " + externalcmd.MenuName + " entries to replace the defaults.")
		}
	}
	if spec.RunMenu != "" {
		return func() {
			if openMenu != nil {
				openMenu(spec.RunMenu)
			}
		}
	}
	return func() {
		if shell == nil || spec.Command == "" {
			message.Errorf("[REQ:EXTERNAL_COMMAND_CONFIG] command is empty; skipping entry %s", spec.Key)
			return
		}
		if spec.Offset != 0 {
			shell(spec.Command, spec.Offset)
			return
		}
		shell(spec.Command)
	}
}
