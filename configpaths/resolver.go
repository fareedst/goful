package configpaths

import (
	"os"

	"github.com/anmitsu/goful/util"
)

const (
	// DefaultStatePath is the legacy location for the persisted UI state.
	DefaultStatePath = "~/.goful/state.json"
	// DefaultHistoryPath is the legacy location for cmdline history.
	DefaultHistoryPath = "~/.goful/history/shell"
	// DefaultCommandsPath is the default location for external command configs.
	DefaultCommandsPath = "~/.goful/external_commands.yaml"
	// DefaultExcludesPath is the default list of filenames to hide.
	DefaultExcludesPath = "~/.goful/excludes"

	// EnvStateKey configures the state path when flags are not provided.
	EnvStateKey = "GOFUL_STATE_PATH"
	// EnvHistoryKey configures the history path when flags are not provided.
	EnvHistoryKey = "GOFUL_HISTORY_PATH"
	// EnvCommandsKey configures the external command config path.
	EnvCommandsKey = "GOFUL_COMMANDS_FILE"
	// EnvExcludesKey configures the exclude list path.
	EnvExcludesKey = "GOFUL_EXCLUDES_FILE"

	flagStateSourceLabel    = "flag:-state"
	flagHistorySourceLabel  = "flag:-history"
	flagCommandsSourceLabel = "flag:-commands"
	flagExcludesSourceLabel = "flag:-exclude-names"
	defaultSourceLabel      = "default"
)

// Paths captures the resolved persistence locations plus their provenance.
type Paths struct {
	State          string
	History        string
	Commands       string
	Excludes       string
	StateSource    string
	HistorySource  string
	CommandsSource string
	ExcludesSource string
}

// Resolver enforces the [REQ:CONFIGURABLE_STATE_PATHS] precedence contract:
// CLI flags override environment variables, which override defaults.
// [IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS]
type Resolver struct {
	LookupEnv func(string) (string, bool)
}

// Resolve returns the final state/history/commands/excludes paths plus provenance metadata.
// [IMPL:STATE_PATH_RESOLVER] [ARCH:STATE_PATH_SELECTION] [REQ:CONFIGURABLE_STATE_PATHS] [REQ:EXTERNAL_COMMAND_CONFIG] [REQ:FILER_EXCLUDE_NAMES]
func (r Resolver) Resolve(flagState, flagHistory, flagCommands, flagExcludes string) Paths {
	state, stateSource := r.resolveOne(flagState, EnvStateKey, DefaultStatePath, flagStateSourceLabel)
	history, historySource := r.resolveOne(flagHistory, EnvHistoryKey, DefaultHistoryPath, flagHistorySourceLabel)
	commands, commandsSource := r.resolveOne(flagCommands, EnvCommandsKey, DefaultCommandsPath, flagCommandsSourceLabel)
	excludes, excludesSource := r.resolveOne(flagExcludes, EnvExcludesKey, DefaultExcludesPath, flagExcludesSourceLabel)

	return Paths{
		State:          state,
		History:        history,
		Commands:       commands,
		Excludes:       excludes,
		StateSource:    stateSource,
		HistorySource:  historySource,
		CommandsSource: commandsSource,
		ExcludesSource: excludesSource,
	}
}

func (r Resolver) resolveOne(flagValue, envKey, defaultValue, flagLabel string) (string, string) {
	if flagValue != "" {
		return util.ExpandPath(flagValue), flagLabel
	}
	if envValue, ok := r.lookupEnv(envKey); ok && envValue != "" {
		return util.ExpandPath(envValue), "env:" + envKey
	}
	return util.ExpandPath(defaultValue), defaultSourceLabel
}

func (r Resolver) lookupEnv(key string) (string, bool) {
	if r.LookupEnv != nil {
		return r.LookupEnv(key)
	}
	return os.LookupEnv(key)
}
