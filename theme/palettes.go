package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

func hx(s string) color.Color { return lipgloss.Color(s) }

func mk(name string, dark bool, accent, accent2, text, muted, subtle, success, danger, warn, border string, heat [5]string) Palette {
	return Palette{
		Name:    name,
		Dark:    dark,
		Accent:  hx(accent),
		Accent2: hx(accent2),
		Text:    hx(text),
		Muted:   hx(muted),
		Subtle:  hx(subtle),
		Success: hx(success),
		Danger:  hx(danger),
		Warn:    hx(warn),
		Border:  hx(border),
		Heat: [5]color.Color{
			hx(heat[0]), hx(heat[1]), hx(heat[2]), hx(heat[3]), hx(heat[4]),
		},
	}
}

var familyPalettes = map[string]func(dark bool) Palette{
	PaletteID:    charmPal,
	"dracula":    draculaPal,
	"nord":       nordPal,
	"gruvbox":    gruvboxPal,
	"catppuccin": catppuccinPal,
	"tokyonight": tokyonightPal,
	"solarized":  solarizedPal,
	"onedark":    onedarkPal,
	"rosepine":   rosepinePal,
	"everforest": everforestPal,
	"kanagawa":   kanagawaPal,
}

func charmPal(dark bool) Palette {
	if dark {
		return mk("charm", true,
			"#FF5F87", "#7571F9", "#E8E8F0", "#8B8BA3", "#5C5C72",
			"#02BF87", "#FE5F86", "#FFB86C", "#7571F9",
			[5]string{"#2A2A38", "#0D3B2E", "#1A6B45", "#1FA86A", "#02BF87"})
	}
	return mk("charm", false,
		"#D6409F", "#5A56E0", "#1A1A2E", "#6B6B80", "#9A9AAE",
		"#02BA84", "#D6406A", "#C97800", "#5A56E0",
		[5]string{"#E4E4EE", "#9BE9A8", "#40C463", "#30A14E", "#02BA84"})
}

func draculaPal(dark bool) Palette {
	// Official Dracula is dark-only.
	_ = dark
	return mk("dracula", true,
		"#FF79C6", "#BD93F9", "#F8F8F2", "#6272A4", "#44475A",
		"#50FA7B", "#FF5555", "#FFB86C", "#BD93F9",
		[5]string{"#282A36", "#1D4A32", "#2E8B57", "#50FA7B", "#8AFF80"})
}

func nordPal(dark bool) Palette {
	if !dark {
		return mk("nord", false,
			"#5E81AC", "#81A1C1", "#2E3440", "#4C566A", "#7B88A1",
			"#A3BE8C", "#BF616A", "#D08770", "#5E81AC",
			[5]string{"#ECEFF4", "#D8DEE9", "#A3BE8C", "#8FBCBB", "#5E81AC"})
	}
	return mk("nord", true,
		"#88C0D0", "#81A1C1", "#ECEFF4", "#D8DEE9", "#4C566A",
		"#A3BE8C", "#BF616A", "#EBCB8B", "#5E81AC",
		[5]string{"#2E3440", "#3B4252", "#A3BE8C", "#8FBCBB", "#88C0D0"})
}

func gruvboxPal(dark bool) Palette {
	if !dark {
		return mk("gruvbox", false,
			"#AF3A03", "#427B58", "#3C3836", "#7C6F64", "#928374",
			"#79740E", "#9D0006", "#B57614", "#076678",
			[5]string{"#FBF1C7", "#D5C4A1", "#B8BB26", "#98971A", "#79740E"})
	}
	return mk("gruvbox", true,
		"#FE8019", "#8EC07C", "#EBDBB2", "#A89984", "#928374",
		"#B8BB26", "#FB4934", "#FABD2F", "#83A598",
		[5]string{"#282828", "#3C3836", "#98971A", "#B8BB26", "#B8BB26"})
}

func catppuccinPal(dark bool) Palette {
	if !dark {
		return mk("catppuccin", false,
			"#EA76CB", "#8839EF", "#4C4F69", "#6C6F85", "#9CA0B0",
			"#40A02B", "#D20F39", "#DF8E1D", "#7287FD",
			[5]string{"#EFF1F5", "#BCC0CC", "#8BD5CA", "#40A02B", "#179299"})
	}
	return mk("catppuccin", true,
		"#F5C2E7", "#CBA6F7", "#CDD6F4", "#A6ADC8", "#6C7086",
		"#A6E3A1", "#F38BA8", "#F9E2AF", "#B4BEFE",
		[5]string{"#1E1E2E", "#313244", "#74C7A0", "#A6E3A1", "#94E2D5"})
}

func tokyonightPal(dark bool) Palette {
	if !dark {
		return mk("tokyonight", false,
			"#F52A65", "#9854F1", "#3760BF", "#6172B0", "#848CB5",
			"#587539", "#F52A65", "#8C6C3E", "#2E7DE9",
			[5]string{"#E1E2E7", "#C4C8DA", "#9ECE6A", "#587539", "#2E7DE9"})
	}
	return mk("tokyonight", true,
		"#BB9AF7", "#7AA2F7", "#C0CAF5", "#A9B1D6", "#565F89",
		"#9ECE6A", "#F7768E", "#E0AF68", "#7DCFFF",
		[5]string{"#1A1B26", "#24283B", "#73DACA", "#9ECE6A", "#7DCFFF"})
}

func solarizedPal(dark bool) Palette {
	if !dark {
		return mk("solarized", false,
			"#D33682", "#6C71C4", "#657B83", "#93A1A1", "#839496",
			"#859900", "#DC322F", "#B58900", "#268BD2",
			[5]string{"#FDF6E3", "#EEE8D5", "#859900", "#2AA198", "#268BD2"})
	}
	return mk("solarized", true,
		"#D33682", "#6C71C4", "#839496", "#657B83", "#586E75",
		"#859900", "#DC322F", "#B58900", "#268BD2",
		[5]string{"#002B36", "#073642", "#859900", "#2AA198", "#268BD2"})
}

func onedarkPal(dark bool) Palette {
	if !dark {
		return mk("onedark", false,
			"#A626A4", "#4078F2", "#383A42", "#696C77", "#A0A1A7",
			"#50A14F", "#E45649", "#C18401", "#4078F2",
			[5]string{"#FAFAFA", "#E5E5E6", "#50A14F", "#0184BC", "#4078F2"})
	}
	return mk("onedark", true,
		"#C678DD", "#61AFEF", "#ABB2BF", "#5C6370", "#4B5263",
		"#98C379", "#E06C75", "#E5C07B", "#56B6C2",
		[5]string{"#282C34", "#3E4451", "#98C379", "#56B6C2", "#61AFEF"})
}

func rosepinePal(dark bool) Palette {
	if !dark {
		return mk("rosepine", false,
			"#B4637A", "#907AA9", "#575279", "#797593", "#9893A5",
			"#286983", "#B4637A", "#EA9D34", "#56949F",
			[5]string{"#FAF4ED", "#FFF7ED", "#56949F", "#286983", "#907AA9"})
	}
	return mk("rosepine", true,
		"#EBBCBA", "#C4A7E7", "#E0DEF4", "#908CAA", "#6E6A86",
		"#31748F", "#EB6F92", "#F6C177", "#9CCFD8",
		[5]string{"#191724", "#1F1D2E", "#31748F", "#9CCFD8", "#C4A7E7"})
}

func everforestPal(dark bool) Palette {
	if !dark {
		return mk("everforest", false,
			"#8DA101", "#3A94C5", "#5C6A72", "#7A8478", "#939F91",
			"#8DA101", "#F85552", "#DFA000", "#35A77C",
			[5]string{"#F3EAD3", "#E5DFC5", "#93C470", "#8DA101", "#35A77C"})
	}
	return mk("everforest", true,
		"#A7C080", "#7FBBB3", "#D3C6AA", "#9DA9A0", "#7A8478",
		"#A7C080", "#E67E80", "#DBBC7F", "#83C092",
		[5]string{"#2B3339", "#323C41", "#83C092", "#A7C080", "#7FBBB3"})
}

func kanagawaPal(dark bool) Palette {
	if !dark {
		return mk("kanagawa", false,
			"#B35B79", "#624C83", "#545464", "#766B90", "#77713F",
			"#6F894E", "#C84053", "#CC6D00", "#4D699B",
			[5]string{"#F2ECBC", "#D5CEA3", "#B7D0AE", "#6F894E", "#5E857A"})
	}
	return mk("kanagawa", true,
		"#D27E99", "#957FB8", "#DCD7BA", "#C8C093", "#727169",
		"#98BB6C", "#FF5D62", "#E6C384", "#7E9CD8",
		[5]string{"#1F1F28", "#2D4F67", "#76946A", "#98BB6C", "#7AA89F"})
}
