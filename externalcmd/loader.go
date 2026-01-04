package externalcmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/anmitsu/goful/util"
	"gopkg.in/yaml.v3"
)

// EnvDebugCommands enables DEBUG output for command loading.
const EnvDebugCommands = "GOFUL_DEBUG_COMMANDS"

// Options controls how the loader resolves and parses the config file.
type Options struct {
	Path     string
	GOOS     string
	ReadFile func(string) ([]byte, error)
	Logf     func(format string, args ...interface{})
	Debug    bool
}

type fileConfig struct {
	entries         []Entry
	inheritDefaults bool
}

type configWrapper struct {
	Commands        []Entry `json:"commands" yaml:"commands"`
	InheritDefaults *bool   `json:"inheritDefaults" yaml:"inheritDefaults"`
}

// Load resolves and parses the external command configuration file.
// Falls back to baked-in defaults if the file is missing or invalid.
// [IMPL:EXTERNAL_COMMAND_LOADER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
func Load(opts Options) ([]Entry, error) {
	defaults := Defaults(opts.GOOS)
	resolvedPath := util.ExpandPath(strings.TrimSpace(opts.Path))
	if resolvedPath == "" {
		debugf(&opts, "using defaults for external commands; no path provided")
		return defaults, nil
	}

	readFn := opts.ReadFile
	if readFn == nil {
		readFn = os.ReadFile
	}

	data, err := readFn(resolvedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			debugf(&opts, "config %s not found; falling back to defaults", resolvedPath)
			return defaults, nil
		}
		return defaults, fmt.Errorf("read external commands %s: %w", resolvedPath, err)
	}

	cfg, err := parseConfig(data)
	if err != nil {
		return defaults, fmt.Errorf("parse external commands %s: %w", resolvedPath, err)
	}

	sanitized, err := sanitizeEntries(cfg.entries, opts.GOOS, &opts, resolvedPath)
	if err != nil {
		return defaults, err
	}

	// [IMPL:EXTERNAL_COMMAND_APPEND] preserve compiled defaults unless configs opt out.
	if cfg.inheritDefaults {
		merged := mergeWithDefaults(defaults, sanitized)
		debugf(&opts, "[IMPL:EXTERNAL_COMMAND_APPEND] inheritDefaults=true; prepended %d command(s) ahead of %d default(s) from %s", len(sanitized), len(defaults), resolvedPath)
		return merged, nil
	}

	debugf(&opts, "[IMPL:EXTERNAL_COMMAND_APPEND] inheritDefaults=false; replacing %d default(s) with %d command(s) from %s", len(defaults), len(sanitized), resolvedPath)
	return sanitized, nil
}

func parseConfig(data []byte) (fileConfig, error) {
	if cfg, err := parseJSONConfig(data); err == nil {
		return cfg, nil
	}
	if cfg, err := parseYAMLConfig(data); err == nil {
		return cfg, nil
	}
	return fileConfig{}, errors.New("expected JSON or YAML array or object with `commands` field")
}

func parseJSONConfig(data []byte) (fileConfig, error) {
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err == nil {
		return fileConfig{entries: entries, inheritDefaults: true}, nil
	}

	var wrapper configWrapper
	if err := json.Unmarshal(data, &wrapper); err == nil {
		return fileConfig{
			entries:         wrapper.Commands,
			inheritDefaults: inheritDefaultsOrTrue(wrapper.InheritDefaults),
		}, nil
	}
	return fileConfig{}, errors.New("json decode failed")
}

func parseYAMLConfig(data []byte) (fileConfig, error) {
	var entries []Entry
	if err := yaml.Unmarshal(data, &entries); err == nil {
		return fileConfig{entries: entries, inheritDefaults: true}, nil
	}

	var wrapper configWrapper
	if err := yaml.Unmarshal(data, &wrapper); err == nil {
		return fileConfig{
			entries:         wrapper.Commands,
			inheritDefaults: inheritDefaultsOrTrue(wrapper.InheritDefaults),
		}, nil
	}
	return fileConfig{}, errors.New("yaml decode failed")
}

func inheritDefaultsOrTrue(flag *bool) bool {
	if flag == nil {
		return true
	}
	return *flag
}

func mergeWithDefaults(defaults, overrides []Entry) []Entry {
	if len(overrides) == 0 {
		return defaults
	}
	merged := make([]Entry, 0, len(defaults)+len(overrides))
	merged = append(merged, overrides...)
	merged = append(merged, defaults...)
	return merged
}

func sanitizeEntries(entries []Entry, goos string, opts *Options, path string) ([]Entry, error) {
	normalized := make([]Entry, 0, len(entries))
	seen := make(map[string]int)
	goos = strings.ToLower(goos)

	for idx, entry := range entries {
		entry.Menu = strings.TrimSpace(entry.Menu)
		if entry.Menu == "" {
			entry.Menu = MenuName
		}
		entry.Key = strings.TrimSpace(entry.Key)
		if entry.Key == "" {
			return nil, fmt.Errorf("entry %d missing `key`", idx)
		}

		if strings.TrimSpace(entry.Label) == "" {
			return nil, fmt.Errorf("entry %q missing `label`", entry.Key)
		}

		commandTrimmed := strings.TrimSpace(entry.Command)
		runMenuTrimmed := strings.TrimSpace(entry.RunMenu)
		if commandTrimmed == "" && runMenuTrimmed == "" {
			return nil, fmt.Errorf("entry %q must provide `command` or `runMenu`", entry.Key)
		}
		if commandTrimmed != "" && runMenuTrimmed != "" {
			return nil, fmt.Errorf("entry %q cannot set both `command` and `runMenu`", entry.Key)
		}
		entry.RunMenu = runMenuTrimmed

		if len(entry.Platforms) > 0 {
			match := false
			for _, platform := range entry.Platforms {
				if strings.EqualFold(strings.TrimSpace(platform), goos) {
					match = true
					break
				}
			}
			if !match {
				debugf(opts, "skipping %s/%s from %s: GOOS=%s not in %+v", entry.Menu, entry.Key, path, goos, entry.Platforms)
				continue
			}
		}

		if entry.Disabled {
			debugf(opts, "skipping disabled command %s/%s from %s", entry.Menu, entry.Key, path)
			continue
		}

		dedupeKey := entry.Menu + "|" + entry.Key
		if prevIdx, exists := seen[dedupeKey]; exists {
			return nil, fmt.Errorf("duplicate shortcut %s/%s (entries %d and %d)", entry.Menu, entry.Key, prevIdx, idx)
		}
		seen[dedupeKey] = idx

		normalized = append(normalized, entry)
	}

	return normalized, nil
}

func debugf(opts *Options, format string, args ...interface{}) {
	if opts == nil || !opts.Debug {
		return
	}
	if opts.Logf != nil {
		opts.Logf(format, args...)
		return
	}
	fmt.Fprintf(os.Stderr, "DEBUG: [IMPL:EXTERNAL_COMMAND_LOADER] "+format+"\n", args...)
}
