module github.com/anmitsu/goful

// [IMPL:GO_MOD_UPDATE] [ARCH:GO_RUNTIME_STRATEGY] [REQ:GO_TOOLCHAIN_LTS]
// Adopt the Go 1.24 LTS baseline for consistent local + CI builds.
go 1.24.0

toolchain go1.24.3

// [IMPL:DEP_BUMP] [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]
require (
	github.com/gdamore/tcell/v2 v2.13.5
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/mattn/go-runewidth v0.0.19
	gopkg.in/yaml.v3 v3.0.1
)

// [IMPL:DEP_BUMP] [ARCH:DEPENDENCY_POLICY] [REQ:DEPENDENCY_REFRESH]
require (
	github.com/clipperhouse/uax29/v2 v2.2.0 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/term v0.38.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

require github.com/cespare/xxhash/v2 v2.3.0 // indirect
