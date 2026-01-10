// Package comparecolors configuration tests.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
package comparecolors

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gdamore/tcell/v2"
)

// TestDefaultConfig tests that DefaultConfig returns sensible defaults.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestDefaultConfig_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Name.Present == "" {
		t.Error("expected non-empty name.present color")
	}
	if cfg.Size.Equal == "" {
		t.Error("expected non-empty size.equal color")
	}
	if cfg.Size.Smallest == "" {
		t.Error("expected non-empty size.smallest color")
	}
	if cfg.Size.Largest == "" {
		t.Error("expected non-empty size.largest color")
	}
	if cfg.Time.Equal == "" {
		t.Error("expected non-empty time.equal color")
	}
	if cfg.Time.Earliest == "" {
		t.Error("expected non-empty time.earliest color")
	}
	if cfg.Time.Latest == "" {
		t.Error("expected non-empty time.latest color")
	}
}

// TestLoad_MissingFile tests that Load returns defaults for missing files.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestLoad_MissingFile_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}

	// Should return defaults
	def := DefaultConfig()
	if cfg.Name.Present != def.Name.Present {
		t.Errorf("expected default name.present, got %s", cfg.Name.Present)
	}
}

// TestLoad_EmptyPath tests that Load returns defaults for empty path.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestLoad_EmptyPath_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}

	def := DefaultConfig()
	if cfg.Name.Present != def.Name.Present {
		t.Errorf("expected default name.present, got %s", cfg.Name.Present)
	}
}

// TestLoad_ValidYAML tests that Load parses valid YAML correctly.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestLoad_ValidYAML_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	// Create temp file with YAML content
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "compare_colors.yaml")

	yamlContent := `
name:
  present: "blue"
size:
  equal: "magenta"
  smallest: "white"
  largest: "black"
time:
  equal: "magenta"
  earliest: "white"
  latest: "black"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if cfg.Name.Present != "blue" {
		t.Errorf("expected name.present=blue, got %s", cfg.Name.Present)
	}
	if cfg.Size.Smallest != "white" {
		t.Errorf("expected size.smallest=white, got %s", cfg.Size.Smallest)
	}
	if cfg.Time.Latest != "black" {
		t.Errorf("expected time.latest=black, got %s", cfg.Time.Latest)
	}
}

// TestLoad_PartialYAML tests that Load fills in defaults for partial configs.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestLoad_PartialYAML_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "compare_colors.yaml")

	// Only specify some values
	yamlContent := `
name:
  present: "blue"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if cfg.Name.Present != "blue" {
		t.Errorf("expected name.present=blue, got %s", cfg.Name.Present)
	}

	// Other values should be defaults
	def := DefaultConfig()
	if cfg.Size.Smallest != def.Size.Smallest {
		t.Errorf("expected default size.smallest, got %s", cfg.Size.Smallest)
	}
}

// TestLoad_InvalidYAML tests that Load returns defaults for invalid YAML.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestLoad_InvalidYAML_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "compare_colors.yaml")

	// Invalid YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}

	// Should still return usable defaults
	def := DefaultConfig()
	if cfg.Name.Present != def.Name.Present {
		t.Errorf("expected default config on error, got %s", cfg.Name.Present)
	}
}

// TestParse_NamedColors tests that Parse correctly converts named colors.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestParse_NamedColors_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg := &Config{
		Name: NameColors{Present: "yellow"},
		Size: SizeColors{
			Equal:    "cyan",
			Smallest: "red",
			Largest:  "green",
		},
		Time: TimeColors{
			Equal:    "cyan",
			Earliest: "red",
			Latest:   "green",
		},
	}

	parsed := cfg.Parse()

	// Check that styles have foreground colors set
	fg, _, _ := parsed.NamePresent.Decompose()
	if fg != tcell.ColorYellow {
		t.Errorf("expected yellow for name.present, got %v", fg)
	}

	fg, _, _ = parsed.SizeSmallest.Decompose()
	if fg != tcell.ColorRed {
		t.Errorf("expected red for size.smallest, got %v", fg)
	}

	fg, _, _ = parsed.SizeLargest.Decompose()
	if fg != tcell.ColorGreen {
		t.Errorf("expected green for size.largest, got %v", fg)
	}
}

// TestParse_HexColors tests that Parse correctly converts hex colors.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestParse_HexColors_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg := &Config{
		Name: NameColors{Present: "#ff0000"}, // Red
		Size: SizeColors{
			Equal:    "#00ff00", // Green
			Smallest: "#0000ff", // Blue
			Largest:  "#ffffff", // White
		},
		Time: TimeColors{
			Equal:    "#00ff00",
			Earliest: "#0000ff",
			Latest:   "#ffffff",
		},
	}

	parsed := cfg.Parse()

	fg, _, _ := parsed.NamePresent.Decompose()
	expected := tcell.NewHexColor(0xff0000)
	if fg != expected {
		t.Errorf("expected #ff0000, got %v", fg)
	}
}

// TestParse_EmptyColor tests that Parse handles empty color strings.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestParse_EmptyColor_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg := &Config{
		Name: NameColors{Present: ""},
		Size: SizeColors{
			Equal:    "",
			Smallest: "",
			Largest:  "",
		},
		Time: TimeColors{
			Equal:    "",
			Earliest: "",
			Latest:   "",
		},
	}

	parsed := cfg.Parse()

	fg, _, _ := parsed.NamePresent.Decompose()
	if fg != tcell.ColorDefault {
		t.Errorf("expected default color for empty string, got %v", fg)
	}
}

// TestParse_InvalidColor tests that Parse handles invalid color names.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestParse_InvalidColor_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg := &Config{
		Name: NameColors{Present: "notacolor"},
		Size: SizeColors{
			Equal:    "invalid",
			Smallest: "nope",
			Largest:  "wrong",
		},
		Time: TimeColors{
			Equal:    "invalid",
			Earliest: "nope",
			Latest:   "wrong",
		},
	}

	parsed := cfg.Parse()

	fg, _, _ := parsed.NamePresent.Decompose()
	if fg != tcell.ColorDefault {
		t.Errorf("expected default color for invalid name, got %v", fg)
	}
}

// TestParse_CaseInsensitive tests that color names are case-insensitive.
// [REQ:FILE_COMPARISON_COLORS] [IMPL:COMPARE_COLOR_CONFIG]
func TestParse_CaseInsensitive_REQ_FILE_COMPARISON_COLORS(t *testing.T) {
	cfg := &Config{
		Name: NameColors{Present: "YELLOW"},
		Size: SizeColors{
			Equal:    "Cyan",
			Smallest: "RED",
			Largest:  "Green",
		},
		Time: TimeColors{
			Equal:    "CYAN",
			Earliest: "Red",
			Latest:   "GREEN",
		},
	}

	parsed := cfg.Parse()

	fg, _, _ := parsed.NamePresent.Decompose()
	if fg != tcell.ColorYellow {
		t.Errorf("expected yellow (case-insensitive), got %v", fg)
	}
}
