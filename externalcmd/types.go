package externalcmd

// MenuName is the canonical menu that hosts external command bindings.
const MenuName = "external-command"

// Entry represents a single external command binding definition.
// [IMPL:EXTERNAL_COMMAND_LOADER] [ARCH:EXTERNAL_COMMAND_REGISTRY] [REQ:EXTERNAL_COMMAND_CONFIG]
type Entry struct {
	Menu      string   `json:"menu" yaml:"menu"`
	Key       string   `json:"key" yaml:"key"`
	Label     string   `json:"label" yaml:"label"`
	Command   string   `json:"command" yaml:"command"`
	RunMenu   string   `json:"runMenu" yaml:"runMenu"`
	Offset    int      `json:"offset" yaml:"offset"`
	Platforms []string `json:"platforms" yaml:"platforms"`
	Disabled  bool     `json:"disabled" yaml:"disabled"`
}
