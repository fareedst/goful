package configpaths

import (
	"testing"

	"github.com/anmitsu/goful/util"
)

func stubLookup(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		val, ok := values[key]
		return val, ok
	}
}

func TestResolvePathsPrefersFlags_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER]
	resolver := Resolver{
		LookupEnv: stubLookup(map[string]string{
			EnvStateKey:    "/env/state.json",
			EnvHistoryKey:  "/env/history",
			EnvCommandsKey: "/env/commands.json",
		}),
	}

	paths := resolver.Resolve("/flag/state.json", "/flag/history", "/flag/commands.json")
	if paths.State != "/flag/state.json" || paths.StateSource != flagStateSourceLabel {
		t.Fatalf("flags must override env/default, got state=%q src=%q", paths.State, paths.StateSource)
	}
	if paths.History != "/flag/history" || paths.HistorySource != flagHistorySourceLabel {
		t.Fatalf("flags must override env/default, got history=%q src=%q", paths.History, paths.HistorySource)
	}
	if paths.Commands != "/flag/commands.json" || paths.CommandsSource != flagCommandsSourceLabel {
		t.Fatalf("flags must override env/default for commands, got %q (%q)", paths.Commands, paths.CommandsSource)
	}
}

func TestResolvePathsFallsBackToEnv_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER]
	stateEnv := "/env/overridden-state.json"
	historyEnv := "/env/overridden-history"
	commandsEnv := "/env/commands.json"
	resolver := Resolver{
		LookupEnv: stubLookup(map[string]string{
			EnvStateKey:    stateEnv,
			EnvHistoryKey:  historyEnv,
			EnvCommandsKey: commandsEnv,
		}),
	}

	paths := resolver.Resolve("", "", "")
	if paths.State != stateEnv || paths.StateSource != "env:"+EnvStateKey {
		t.Fatalf("env should supply state path, got %q (%q)", paths.State, paths.StateSource)
	}
	if paths.History != historyEnv || paths.HistorySource != "env:"+EnvHistoryKey {
		t.Fatalf("env should supply history path, got %q (%q)", paths.History, paths.HistorySource)
	}
	if paths.Commands != commandsEnv || paths.CommandsSource != "env:"+EnvCommandsKey {
		t.Fatalf("env should supply commands path, got %q (%q)", paths.Commands, paths.CommandsSource)
	}
}

func TestResolvePathsDefaults_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER]
	resolver := Resolver{}
	paths := resolver.Resolve("", "", "")

	wantState := util.ExpandPath(DefaultStatePath)
	wantHistory := util.ExpandPath(DefaultHistoryPath)
	wantCommands := util.ExpandPath(DefaultCommandsPath)
	if paths.State != wantState || paths.StateSource != defaultSourceLabel {
		t.Fatalf("default state mismatch: got %q (%q), want %q", paths.State, paths.StateSource, wantState)
	}
	if paths.History != wantHistory || paths.HistorySource != defaultSourceLabel {
		t.Fatalf("default history mismatch: got %q (%q), want %q", paths.History, paths.HistorySource, wantHistory)
	}
	if paths.Commands != wantCommands || paths.CommandsSource != defaultSourceLabel {
		t.Fatalf("default commands mismatch: got %q (%q), want %q", paths.Commands, paths.CommandsSource, wantCommands)
	}
}

func TestResolvePathsIgnoresEmptyEnv_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER]
	resolver := Resolver{
		LookupEnv: stubLookup(map[string]string{
			EnvStateKey:    "",
			EnvHistoryKey:  "",
			EnvCommandsKey: "",
		}),
	}
	paths := resolver.Resolve("", "", "")
	if paths.StateSource != defaultSourceLabel || paths.HistorySource != defaultSourceLabel || paths.CommandsSource != defaultSourceLabel {
		t.Fatalf("empty env values should fall back to defaults, got stateSrc=%q historySrc=%q commandsSrc=%q", paths.StateSource, paths.HistorySource, paths.CommandsSource)
	}
}
