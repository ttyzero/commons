package theme

import (
	"os"
	"path/filepath"
	"strings"
)

// Saved is the last look ttythemer (or a hand-edited file) persisted.
type Saved struct {
	Spec       Spec
	HasSpec    bool
	Borderless bool
	HasBorders bool
}

// ConfigPath is $TTYTHEME_FILE, else $XDG_CONFIG_HOME/ttyzero/theme,
// else ~/.config/ttyzero/theme. Same XDG rule as ttybus (not UserConfigDir).
func ConfigPath() string {
	if p := os.Getenv(EnvFile); p != "" {
		return p
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "ttyzero", "theme")
}

// Load reads the theme file. ok is false when missing or empty of commands.
func Load() (Saved, bool) {
	path := ConfigPath()
	if path == "" {
		return Saved{}, false
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Saved{}, false
	}
	var saved Saved
	ok := false
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cmd, parsed := ParseBus(line)
		if !parsed {
			continue
		}
		if cmd.Spec != nil {
			saved.Spec = *cmd.Spec
			saved.HasSpec = true
			ok = true
		}
		if cmd.Borderless != nil {
			saved.Borderless = *cmd.Borderless
			saved.HasBorders = true
			ok = true
		}
	}
	return saved, ok
}

// Save writes Encode(spec, borderless) atomically. ttythemer calls this
// whenever it broadcasts so the next pane launch matches.
func Save(spec Spec, borderless bool) error {
	path := ConfigPath()
	if path == "" {
		return os.ErrInvalid
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body := "# ttyzero look — written by ttythemer\n" + Encode(spec, borderless) + "\n"
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(body), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
