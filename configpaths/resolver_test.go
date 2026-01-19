package configpaths

import (
	"testing"

	"github.com/fareedst/goful/util"
)

func stubLookup(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		val, ok := values[key]
		return val, ok
	}
}

func TestResolvePathsPrefersFlags_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER] [REQ:FILE_COMPARISON_COLORS]
	resolver := Resolver{
		LookupEnv: stubLookup(map[string]string{
			EnvStateKey:         "/env/state.json",
			EnvHistoryKey:       "/env/history",
			EnvCommandsKey:      "/env/commands.json",
			EnvExcludesKey:      "/env/excludes.txt",
			EnvCompareColorsKey: "/env/compare_colors.yaml",
		}),
	}

	paths := resolver.Resolve("/flag/state.json", "/flag/history", "/flag/commands.json", "/flag/excludes.txt", "/flag/compare_colors.yaml")
	if paths.State != "/flag/state.json" || paths.StateSource != flagStateSourceLabel {
		t.Fatalf("flags must override env/default, got state=%q src=%q", paths.State, paths.StateSource)
	}
	if paths.History != "/flag/history" || paths.HistorySource != flagHistorySourceLabel {
		t.Fatalf("flags must override env/default, got history=%q src=%q", paths.History, paths.HistorySource)
	}
	if paths.Commands != "/flag/commands.json" || paths.CommandsSource != flagCommandsSourceLabel {
		t.Fatalf("flags must override env/default for commands, got %q (%q)", paths.Commands, paths.CommandsSource)
	}
	if paths.Excludes != "/flag/excludes.txt" || paths.ExcludesSource != flagExcludesSourceLabel {
		t.Fatalf("flags must override env/default for excludes, got %q (%q)", paths.Excludes, paths.ExcludesSource)
	}
	if paths.CompareColors != "/flag/compare_colors.yaml" || paths.CompareColorsSource != flagCompareColorsSourceLabel {
		t.Fatalf("flags must override env/default for compare-colors, got %q (%q)", paths.CompareColors, paths.CompareColorsSource)
	}
}

func TestResolvePathsFallsBackToEnv_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER] [REQ:FILE_COMPARISON_COLORS]
	stateEnv := "/env/overridden-state.json"
	historyEnv := "/env/overridden-history"
	commandsEnv := "/env/commands.json"
	excludesEnv := "/env/excludes.txt"
	compareColorsEnv := "/env/compare_colors.yaml"
	resolver := Resolver{
		LookupEnv: stubLookup(map[string]string{
			EnvStateKey:         stateEnv,
			EnvHistoryKey:       historyEnv,
			EnvCommandsKey:      commandsEnv,
			EnvExcludesKey:      excludesEnv,
			EnvCompareColorsKey: compareColorsEnv,
		}),
	}

	paths := resolver.Resolve("", "", "", "", "")
	if paths.State != stateEnv || paths.StateSource != "env:"+EnvStateKey {
		t.Fatalf("env should supply state path, got %q (%q)", paths.State, paths.StateSource)
	}
	if paths.History != historyEnv || paths.HistorySource != "env:"+EnvHistoryKey {
		t.Fatalf("env should supply history path, got %q (%q)", paths.History, paths.HistorySource)
	}
	if paths.Commands != commandsEnv || paths.CommandsSource != "env:"+EnvCommandsKey {
		t.Fatalf("env should supply commands path, got %q (%q)", paths.Commands, paths.CommandsSource)
	}
	if paths.Excludes != excludesEnv || paths.ExcludesSource != "env:"+EnvExcludesKey {
		t.Fatalf("env should supply excludes path, got %q (%q)", paths.Excludes, paths.ExcludesSource)
	}
	if paths.CompareColors != compareColorsEnv || paths.CompareColorsSource != "env:"+EnvCompareColorsKey {
		t.Fatalf("env should supply compare-colors path, got %q (%q)", paths.CompareColors, paths.CompareColorsSource)
	}
}

func TestResolvePathsDefaults_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER] [REQ:FILE_COMPARISON_COLORS]
	resolver := Resolver{}
	paths := resolver.Resolve("", "", "", "", "")

	wantState := util.ExpandPath(DefaultStatePath)
	wantHistory := util.ExpandPath(DefaultHistoryPath)
	wantCommands := util.ExpandPath(DefaultCommandsPath)
	wantExcludes := util.ExpandPath(DefaultExcludesPath)
	wantCompareColors := util.ExpandPath(DefaultCompareColorsPath)
	if paths.State != wantState || paths.StateSource != defaultSourceLabel {
		t.Fatalf("default state mismatch: got %q (%q), want %q", paths.State, paths.StateSource, wantState)
	}
	if paths.History != wantHistory || paths.HistorySource != defaultSourceLabel {
		t.Fatalf("default history mismatch: got %q (%q), want %q", paths.History, paths.HistorySource, wantHistory)
	}
	if paths.Commands != wantCommands || paths.CommandsSource != defaultSourceLabel {
		t.Fatalf("default commands mismatch: got %q (%q), want %q", paths.Commands, paths.CommandsSource, wantCommands)
	}
	if paths.Excludes != wantExcludes || paths.ExcludesSource != defaultSourceLabel {
		t.Fatalf("default excludes mismatch: got %q (%q), want %q", paths.Excludes, paths.ExcludesSource, wantExcludes)
	}
	if paths.CompareColors != wantCompareColors || paths.CompareColorsSource != defaultSourceLabel {
		t.Fatalf("default compare-colors mismatch: got %q (%q), want %q", paths.CompareColors, paths.CompareColorsSource, wantCompareColors)
	}
}

func TestResolvePathsIgnoresEmptyEnv_REQ_CONFIGURABLE_STATE_PATHS(t *testing.T) {
	// [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [ARCH:STATE_PATH_SELECTION] [IMPL:STATE_PATH_RESOLVER] [REQ:FILE_COMPARISON_COLORS]
	resolver := Resolver{
		LookupEnv: stubLookup(map[string]string{
			EnvStateKey:         "",
			EnvHistoryKey:       "",
			EnvCommandsKey:      "",
			EnvExcludesKey:      "",
			EnvCompareColorsKey: "",
		}),
	}
	paths := resolver.Resolve("", "", "", "", "")
	if paths.StateSource != defaultSourceLabel || paths.HistorySource != defaultSourceLabel || paths.CommandsSource != defaultSourceLabel || paths.ExcludesSource != defaultSourceLabel || paths.CompareColorsSource != defaultSourceLabel {
		t.Fatalf("empty env values should fall back to defaults, got stateSrc=%q historySrc=%q commandsSrc=%q excludesSrc=%q compareColorsSrc=%q", paths.StateSource, paths.HistorySource, paths.CommandsSource, paths.ExcludesSource, paths.CompareColorsSource)
	}
}
