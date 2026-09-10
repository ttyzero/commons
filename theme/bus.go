package theme

import "strings"

const Channel = "theme"

// Command is a live look change from the ttybus `theme` channel.
type Command struct {
	Spec       *Spec
	Borderless *bool
}

// ParseBus understands one-line look commands, e.g.
//
//	THEME dark
//	THEME light BORDERS=0
//	BORDERS=1
//	BORDERLESS=1
//
// BODERS is accepted as a typo for BORDERS.
func ParseBus(payload string) (Command, bool) {
	fields := strings.Fields(payload)
	if len(fields) == 0 {
		return Command{}, false
	}
	var cmd Command
	ok := false
	for i := 0; i < len(fields); i++ {
		tok := fields[i]
		key, val, kv := strings.Cut(tok, "=")
		upper := strings.ToUpper(key)
		if !kv && (upper == "THEME" || upper == "TTYTHEME") {
			if i+1 >= len(fields) {
				continue
			}
			i++
			spec := Parse(fields[i])
			cmd.Spec = &spec
			ok = true
			continue
		}
		if !kv {
			continue
		}
		switch upper {
		case "THEME", "TTYTHEME":
			spec := Parse(val)
			cmd.Spec = &spec
			ok = true
		case "BORDERS", "BODERS":
			if on, parsed := parseOn(val); parsed {
				b := !on
				cmd.Borderless = &b
				ok = true
			}
		case "BORDERLESS":
			if on, parsed := parseOn(val); parsed {
				cmd.Borderless = &on
				ok = true
			}
		}
	}
	return cmd, ok
}

func parseOn(s string) (on bool, ok bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	default:
		return false, false
	}
}

// Apply mutates spec/palette/borderless from a bus command.
// osc requests a fresh terminal background query when THEME is auto.
func Apply(cmd Command, spec *Spec, pal *Palette, borderless *bool) (osc bool) {
	if cmd.Spec != nil {
		*spec = *cmd.Spec
		*pal = For(*spec)
		osc = spec.Polarity == "auto"
	}
	if cmd.Borderless != nil {
		*borderless = *cmd.Borderless
	}
	return osc
}
