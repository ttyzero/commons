// Package theme is the shared palette for ttyzero companion panes
// (gitwing, peek, buscope).
//
// Resolution, highest precedence first:
//
//  1. --theme flag (passed in by the binary)
//  2. $TTYTHEME — family override for every companion pane
//  3. $CLITHEME — proposed cross-tool standard (dark|light|auto)
//  4. OSC 11 via Bubble Tea's RequestBackgroundColor (live terminal bg)
//  5. $COLORFGBG — cheap, static, set by some terminals (rxvt, Konsole, iTerm)
//  6. dark
//
// Named values for TTYTHEME and --theme:
//
//	auto         follow the terminal (default)
//	dark, light  force polarity; Charm palette
//	charm        alias for auto
//	charm-dark, charm-light
//
// Peek and buscope should honor the same $TTYTHEME so a tmux workspace
// of companion panes stays visually one family. Do not paint a panel
// background — the terminal's own theme shows through; we only tint
// type, borders, and accents.
package theme

import (
	"image/color"
	"os"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	EnvTTY        = "TTYTHEME"
	EnvCLI        = "CLITHEME"
	EnvFGBG       = "COLORFGBG"
	EnvBorderless = "TTYBORDERLESS"
	Default       = "auto"
	PaletteID     = "charm"
)

// Spec is a resolved theme request before the terminal has answered OSC 11.
type Spec struct {
	// Polarity is "auto", "dark", or "light".
	Polarity string
	// Palette is a named color family. Currently only "charm".
	Palette string
}

// Parse interprets a --theme / $TTYTHEME / $CLITHEME value.
func Parse(s string) Spec {
	s = strings.TrimSpace(strings.ToLower(s))
	if i := strings.IndexByte(s, ':'); i >= 0 {
		// CLITHEME allows "dark:modifier"; ignore unknown modifiers.
		s = s[:i]
	}
	spec := Spec{Polarity: "auto", Palette: PaletteID}
	switch s {
	case "", "auto", "charm":
		return spec
	case "dark", "light":
		spec.Polarity = s
		return spec
	case "charm-dark":
		spec.Polarity = "dark"
		return spec
	case "charm-light":
		spec.Polarity = "light"
		return spec
	default:
		return spec
	}
}

// FromEnv reads TTYTHEME, then CLITHEME. Empty means auto.
func FromEnv() Spec {
	if v := os.Getenv(EnvTTY); v != "" {
		return Parse(v)
	}
	if v := os.Getenv(EnvCLI); v != "" {
		return Parse(v)
	}
	return Parse(Default)
}

// Borderless is true when --borderless is set or $TTYBORDERLESS is
// 1/true/yes. Shared by gitwing, peek, and buscope.
func Borderless(flag bool) bool {
	if flag {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvBorderless))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// Resolve combines an explicit flag with the environment. The flag wins
// when non-empty.
func Resolve(flag string) Spec {
	if strings.TrimSpace(flag) != "" {
		return Parse(flag)
	}
	return FromEnv()
}

// COLORFGBG reports a polarity guess from $COLORFGBG ("fg;bg" ANSI
// indices). ok is false when the variable is unset or unparseable.
// Background 7 or 15 is treated as light; everything else as dark.
func COLORFGBG() (dark bool, ok bool) {
	return parseCOLORFGBG(os.Getenv(EnvFGBG))
}

func parseCOLORFGBG(v string) (dark bool, ok bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return true, false
	}
	parts := strings.Split(v, ";")
	bg := parts[len(parts)-1]
	n, err := strconv.Atoi(bg)
	if err != nil || n < 0 || n > 15 {
		return true, false
	}
	light := n == 7 || n == 15
	return !light, true
}

// InitialDark is the polarity to use before OSC 11 returns.
func (s Spec) InitialDark() bool {
	switch s.Polarity {
	case "light":
		return false
	case "dark":
		return true
	default:
		if dark, ok := COLORFGBG(); ok {
			return dark
		}
		return true
	}
}

// ApplyOSC updates polarity from a terminal background query.
// Forced dark/light specs ignore the query.
func (s Spec) ApplyOSC(bgIsDark bool) Spec {
	if s.Polarity == "auto" {
		out := s
		if bgIsDark {
			out.Polarity = "dark"
		} else {
			out.Polarity = "light"
		}
		return out
	}
	return s
}

// Palette is a resolved set of colors. Background is left unset so the
// terminal's own theme shows through.
type Palette struct {
	Name    string
	Dark    bool
	Accent  color.Color
	Accent2 color.Color
	Text    color.Color
	Muted   color.Color
	Subtle  color.Color
	Success color.Color
	Danger  color.Color
	Warn    color.Color
	Border  color.Color
	// Heat is a 5-stop GitHub-style activity ramp (0 empty … 4 hottest).
	Heat [5]color.Color
}

// New builds the Charm palette for a polarity.
func New(dark bool) Palette {
	ld := lipgloss.LightDark(dark)
	return Palette{
		Name:    PaletteID,
		Dark:    dark,
		Accent:  ld(lipgloss.Color("#D6409F"), lipgloss.Color("#FF5F87")),
		Accent2: ld(lipgloss.Color("#5A56E0"), lipgloss.Color("#7571F9")),
		Text:    ld(lipgloss.Color("#1A1A2E"), lipgloss.Color("#E8E8F0")),
		Muted:   ld(lipgloss.Color("#6B6B80"), lipgloss.Color("#8B8BA3")),
		Subtle:  ld(lipgloss.Color("#9A9AAE"), lipgloss.Color("#5C5C72")),
		Success: ld(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87")),
		Danger:  ld(lipgloss.Color("#D6406A"), lipgloss.Color("#FE5F86")),
		Warn:    ld(lipgloss.Color("#C97800"), lipgloss.Color("#FFB86C")),
		Border:  ld(lipgloss.Color("#5A56E0"), lipgloss.Color("#7571F9")),
		Heat: [5]color.Color{
			ld(lipgloss.Color("#E4E4EE"), lipgloss.Color("#2A2A38")),
			ld(lipgloss.Color("#9BE9A8"), lipgloss.Color("#0D3B2E")),
			ld(lipgloss.Color("#40C463"), lipgloss.Color("#1A6B45")),
			ld(lipgloss.Color("#30A14E"), lipgloss.Color("#1FA86A")),
			ld(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87")),
		},
	}
}
