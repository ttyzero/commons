package theme

// Entry is one named look for switchers such as ttythemer.
type Entry struct {
	Name string
	Desc string
	Spec Spec
}

// Catalog is the set we advertise on the theme bus.
func Catalog() []Entry {
	names := []struct{ name, desc string }{
		{"auto", "follow the terminal (OSC 11)"},
		{"dark", "charm · dark"},
		{"light", "charm · light"},
		{"catppuccin", "Catppuccin Mocha — the ricing default"},
		{"catppuccin-latte", "Catppuccin Latte"},
		{"dracula", "Dracula"},
		{"nord", "Nord"},
		{"nord-light", "Nord Snow Storm"},
		{"gruvbox", "Gruvbox dark"},
		{"gruvbox-light", "Gruvbox light"},
		{"tokyonight", "Tokyo Night"},
		{"tokyonight-day", "Tokyo Night Day"},
		{"solarized", "Solarized dark"},
		{"solarized-light", "Solarized light"},
		{"onedark", "One Dark"},
		{"onelight", "One Light"},
		{"rosepine", "Rosé Pine"},
		{"rosepine-dawn", "Rosé Pine Dawn"},
		{"everforest", "Everforest"},
		{"everforest-light", "Everforest light"},
		{"kanagawa", "Kanagawa Wave"},
		{"kanagawa-lotus", "Kanagawa Lotus"},
	}
	out := make([]Entry, 0, len(names))
	for _, n := range names {
		out = append(out, Entry{Name: n.name, Desc: n.desc, Spec: Parse(n.name)})
	}
	return out
}
