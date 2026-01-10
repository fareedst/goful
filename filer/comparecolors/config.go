// Package comparecolors provides configurable color schemes for cross-directory file comparison.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
package comparecolors

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
)

// Config holds the color scheme for file comparison highlighting.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type Config struct {
	Name NameColors `yaml:"name"`
	Size SizeColors `yaml:"size"`
	Time TimeColors `yaml:"time"`
}

// NameColors defines colors for file name comparison states.
type NameColors struct {
	Present string `yaml:"present"` // File name appears in multiple directories
}

// SizeColors defines colors for file size comparison states.
type SizeColors struct {
	Equal    string `yaml:"equal"`    // Same size across directories
	Smallest string `yaml:"smallest"` // Smallest among same-named files
	Largest  string `yaml:"largest"`  // Largest among same-named files
}

// TimeColors defines colors for modification time comparison states.
type TimeColors struct {
	Equal    string `yaml:"equal"`    // Same modification time
	Earliest string `yaml:"earliest"` // Oldest among same-named files
	Latest   string `yaml:"latest"`   // Newest among same-named files
}

// DefaultConfig returns sensible default colors for comparison highlighting.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func DefaultConfig() *Config {
	return &Config{
		Name: NameColors{
			Present: "yellow",
		},
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
}

// Load reads the comparison color configuration from a YAML file.
// Returns DefaultConfig if the file does not exist or is invalid.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func Load(path string) (*Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), fmt.Errorf("read compare colors config: %w", err)
	}

	cfg := DefaultConfig() // Start with defaults so partial configs work
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), fmt.Errorf("parse compare colors config: %w", err)
	}

	return cfg, nil
}

// ParsedConfig holds the resolved tcell.Style values for each comparison state.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
type ParsedConfig struct {
	NamePresent  tcell.Style
	SizeEqual    tcell.Style
	SizeSmallest tcell.Style
	SizeLargest  tcell.Style
	TimeEqual    tcell.Style
	TimeEarliest tcell.Style
	TimeLatest   tcell.Style
}

// Parse converts the string-based Config to resolved tcell.Style values.
// [IMPL:COMPARE_COLOR_CONFIG] [ARCH:FILE_COMPARISON_ENGINE] [REQ:FILE_COMPARISON_COLORS]
func (c *Config) Parse() *ParsedConfig {
	base := tcell.StyleDefault
	return &ParsedConfig{
		NamePresent:  base.Foreground(parseColor(c.Name.Present)).Bold(true),
		SizeEqual:    base.Foreground(parseColor(c.Size.Equal)),
		SizeSmallest: base.Foreground(parseColor(c.Size.Smallest)),
		SizeLargest:  base.Foreground(parseColor(c.Size.Largest)),
		TimeEqual:    base.Foreground(parseColor(c.Time.Equal)),
		TimeEarliest: base.Foreground(parseColor(c.Time.Earliest)),
		TimeLatest:   base.Foreground(parseColor(c.Time.Latest)),
	}
}

// parseColor converts a color name or hex code to a tcell.Color.
// Supports named colors (red, green, blue, etc.) and hex codes (#RRGGBB).
func parseColor(s string) tcell.Color {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return tcell.ColorDefault
	}

	// Check for hex color
	if strings.HasPrefix(s, "#") {
		if len(s) == 7 {
			if val, err := strconv.ParseInt(s[1:], 16, 32); err == nil {
				return tcell.NewHexColor(int32(val))
			}
		}
		return tcell.ColorDefault
	}

	// Named colors
	switch s {
	case "black":
		return tcell.ColorBlack
	case "red":
		return tcell.ColorRed
	case "green":
		return tcell.ColorGreen
	case "yellow":
		return tcell.ColorYellow
	case "blue":
		return tcell.ColorBlue
	case "magenta", "fuchsia":
		return tcell.ColorFuchsia
	case "cyan", "aqua":
		return tcell.ColorAqua
	case "white":
		return tcell.ColorWhite
	case "lime":
		return tcell.ColorLime
	case "navy":
		return tcell.ColorNavy
	case "olive":
		return tcell.ColorOlive
	case "purple":
		return tcell.ColorPurple
	case "teal":
		return tcell.ColorTeal
	case "silver":
		return tcell.ColorSilver
	case "gray", "grey":
		return tcell.ColorGray
	case "maroon":
		return tcell.ColorMaroon
	case "orange":
		return tcell.ColorOrange
	default:
		return tcell.ColorDefault
	}
}
