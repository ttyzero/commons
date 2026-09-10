// Package theme is the shared palette for ttyzero companion panes
// (gitwing, peek, buscope).
//
// Resolution, highest precedence first:
//
//  1. --theme flag (passed in by the binary)
//  2. $TTYTHEME — family override for every companion pane
//  3. $CLITHEME — proposed cross-tool standard (dark|light|auto)
//  4. ~/.config/ttyzero/theme (or $XDG_CONFIG_HOME/ttyzero/theme) —
//     written by ttythemer; same one-line THEME/BORDERS as the bus
//  5. OSC 11 via Bubble Tea's RequestBackgroundColor (live terminal bg)
//  6. $COLORFGBG — cheap, static, set by some terminals (rxvt, Konsole, iTerm)
//  7. dark
//
// Named values for TTYTHEME and --theme:
//
//	auto, dark, light, charm, charm-dark, charm-light
//	catppuccin (mocha), catppuccin-latte
//	dracula
//	nord, nord-light
//	gruvbox, gruvbox-light
//	tokyonight, tokyonight-day
//	solarized, solarized-light
//	onedark, onelight
//	rosepine, rosepine-dawn
//	everforest, everforest-light
//	kanagawa, kanagawa-lotus
//
// Peek and buscope should honor the same $TTYTHEME so a tmux workspace
// of companion panes stays visually one family. Do not paint a panel
// background — the terminal's own theme shows through; we only tint
// type, borders, and accents. Catalog() is the advertised set.
package theme

import (
	"image/color"
	"os"
	"strconv"
	"strings"
)

const (
	EnvTTY        = "TTYTHEME"
	EnvCLI        = "CLITHEME"
	EnvFGBG       = "COLORFGBG"
	EnvBorderless = "TTYBORDERLESS"
	EnvFile       = "TTYTHEME_FILE"
	Default       = "auto"
	PaletteID     = "charm"
)

// Spec is a resolved theme request before the terminal has answered OSC 11.
type Spec struct {
	// Name is the wire name on the theme bus (nord, gruvbox-light, auto…).
	Name string
	// Polarity is "auto", "dark", or "light".
	Polarity string
	// Palette is the family: charm, catppuccin, dracula, nord, …
	Palette string
}

type named struct {
	family, pol, wire string
}

var namedLooks = map[string]named{
	"auto":             {PaletteID, "auto", "auto"},
	"charm":            {PaletteID, "auto", "auto"},
	"dark":             {PaletteID, "dark", "dark"},
	"light":            {PaletteID, "light", "light"},
	"charm-dark":       {PaletteID, "dark", "charm-dark"},
	"charm-light":      {PaletteID, "light", "charm-light"},
	"dracula":          {"dracula", "dark", "dracula"},
	"nord":             {"nord", "dark", "nord"},
	"nord-light":       {"nord", "light", "nord-light"},
	"nord-snow":        {"nord", "light", "nord-light"},
	"gruvbox":          {"gruvbox", "dark", "gruvbox"},
	"gruvbox-dark":     {"gruvbox", "dark", "gruvbox"},
	"gruvbox-light":    {"gruvbox", "light", "gruvbox-light"},
	"catppuccin":       {"catppuccin", "dark", "catppuccin"},
	"catppuccin-mocha": {"catppuccin", "dark", "catppuccin"},
	"mocha":            {"catppuccin", "dark", "catppuccin"},
	"catppuccin-latte": {"catppuccin", "light", "catppuccin-latte"},
	"latte":            {"catppuccin", "light", "catppuccin-latte"},
	"tokyonight":       {"tokyonight", "dark", "tokyonight"},
	"tokyo":            {"tokyonight", "dark", "tokyonight"},
	"tokyo-night":      {"tokyonight", "dark", "tokyonight"},
	"tokyonight-day":   {"tokyonight", "light", "tokyonight-day"},
	"solarized":        {"solarized", "dark", "solarized"},
	"solarized-dark":   {"solarized", "dark", "solarized"},
	"solarized-light":  {"solarized", "light", "solarized-light"},
	"onedark":          {"onedark", "dark", "onedark"},
	"one-dark":         {"onedark", "dark", "onedark"},
	"onelight":         {"onedark", "light", "onelight"},
	"one-light":        {"onedark", "light", "onelight"},
	"rosepine":         {"rosepine", "dark", "rosepine"},
	"rose-pine":        {"rosepine", "dark", "rosepine"},
	"rosepine-dawn":    {"rosepine", "light", "rosepine-dawn"},
	"rose-pine-dawn":   {"rosepine", "light", "rosepine-dawn"},
	"everforest":       {"everforest", "dark", "everforest"},
	"everforest-light": {"everforest", "light", "everforest-light"},
	"kanagawa":         {"kanagawa", "dark", "kanagawa"},
	"kanagawa-wave":    {"kanagawa", "dark", "kanagawa"},
	"kanagawa-lotus":   {"kanagawa", "light", "kanagawa-lotus"},
	"kanagawa-light":   {"kanagawa", "light", "kanagawa-lotus"},
}

// Parse interprets a --theme / $TTYTHEME / $CLITHEME / theme-bus value.
func Parse(s string) Spec {
	s = strings.TrimSpace(strings.ToLower(s))
	if i := strings.IndexByte(s, ':'); i >= 0 {
		// CLITHEME allows "dark:modifier"; ignore unknown modifiers.
		s = s[:i]
	}
	if n, ok := namedLooks[s]; ok {
		return Spec{Name: n.wire, Polarity: n.pol, Palette: n.family}
	}
	return Spec{Name: "auto", Polarity: "auto", Palette: PaletteID}
}

// FromEnv reads TTYTHEME, then CLITHEME, then the ttyzero theme file.
// Empty means auto.
func FromEnv() Spec {
	if v := os.Getenv(EnvTTY); v != "" {
		return Parse(v)
	}
	if v := os.Getenv(EnvCLI); v != "" {
		return Parse(v)
	}
	if s, ok := Load(); ok && s.HasSpec {
		return s.Spec
	}
	return Parse(Default)
}

// Borderless is true when --borderless is set, $TTYBORDERLESS is
// 1/true/yes, or the theme file says BORDERS=0.
func Borderless(flag bool) bool {
	if flag {
		return true
	}
	if v, ok := os.LookupEnv(EnvBorderless); ok && strings.TrimSpace(v) != "" {
		on, parsed := parseOn(v)
		return parsed && on
	}
	if s, ok := Load(); ok && s.HasBorders {
		return s.Borderless
	}
	return false
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

// New builds the Charm palette for a polarity (tests and fallbacks).
func New(dark bool) Palette {
	pol := "light"
	if dark {
		pol = "dark"
	}
	return For(Spec{Name: pol, Polarity: pol, Palette: PaletteID})
}

// For resolves a named spec to colors.
func For(s Spec) Palette {
	fam := s.Palette
	if fam == "" {
		fam = PaletteID
	}
	fn, ok := familyPalettes[fam]
	if !ok {
		fn = familyPalettes[PaletteID]
	}
	return fn(s.InitialDark())
}
