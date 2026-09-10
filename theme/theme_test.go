package theme

import (
	"path/filepath"
	"testing"
)

func isolateFile(t *testing.T) {
	t.Helper()
	t.Setenv(EnvFile, filepath.Join(t.TempDir(), "theme"))
}

func TestParseBus(t *testing.T) {
	cmd, ok := ParseBus("THEME dark")
	if !ok || cmd.Spec == nil || cmd.Spec.Polarity != "dark" {
		t.Fatalf("THEME dark: %+v %v", cmd, ok)
	}
	cmd, ok = ParseBus("BORDERS=0")
	if !ok || cmd.Borderless == nil || !*cmd.Borderless {
		t.Fatalf("BORDERS=0 should be borderless")
	}
	cmd, ok = ParseBus("BODERS=1")
	if !ok || cmd.Borderless == nil || *cmd.Borderless {
		t.Fatalf("BORDERS=1 should keep frame")
	}
	cmd, ok = ParseBus("THEME light BORDERLESS=1")
	if !ok || cmd.Spec == nil || cmd.Spec.Polarity != "light" || cmd.Borderless == nil || !*cmd.Borderless {
		t.Fatalf("combo %+v", cmd)
	}
	if _, ok := ParseBus("/tmp/foo.go:12"); ok {
		t.Fatal("files payload must not parse as theme")
	}
	cmd, ok = ParseBus("THEME nord")
	if !ok || cmd.Spec.Name != "nord" || cmd.Spec.Palette != "nord" {
		t.Fatalf("nord %+v", cmd)
	}
}

func TestForNamed(t *testing.T) {
	p := For(Parse("dracula"))
	if !p.Dark || p.Name != "dracula" {
		t.Fatalf("%+v", p)
	}
	p = For(Parse("catppuccin-latte"))
	if p.Dark || p.Name != "catppuccin" {
		t.Fatalf("latte %+v", p)
	}
	p = For(Parse("kanagawa"))
	if !p.Dark || p.Name != "kanagawa" {
		t.Fatalf("kanagawa %+v", p)
	}
}

func TestCatalogRoundtrip(t *testing.T) {
	for _, e := range Catalog() {
		got := Parse(e.Name)
		if got.Name != e.Spec.Name || got.Palette != e.Spec.Palette || got.Polarity != e.Spec.Polarity {
			t.Errorf("Parse(%q) = %+v, catalog %+v", e.Name, got, e.Spec)
		}
		p := For(e.Spec)
		if p.Name == "" {
			t.Errorf("empty palette for %q", e.Name)
		}
		if e.Name != "auto" && e.Spec.Polarity != "auto" && p.Dark != (e.Spec.Polarity == "dark") {
			t.Errorf("%s polarity %s pal.Dark=%v", e.Name, e.Spec.Polarity, p.Dark)
		}
	}
}

func TestBorderless(t *testing.T) {
	isolateFile(t)
	t.Setenv(EnvBorderless, "")
	if Borderless(false) {
		t.Fatal("default")
	}
	if !Borderless(true) {
		t.Fatal("flag")
	}
	t.Setenv(EnvBorderless, "yes")
	if !Borderless(false) {
		t.Fatal("env")
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		in, pol, pal string
	}{
		{"", "auto", "charm"},
		{"auto", "auto", "charm"},
		{"DARK", "dark", "charm"},
		{"light", "light", "charm"},
		{"charm", "auto", "charm"},
		{"charm-dark", "dark", "charm"},
		{"charm-light", "light", "charm"},
		{"dark:high-contrast", "dark", "charm"},
		{"mystery", "auto", "charm"},
		{"nord", "dark", "nord"},
		{"latte", "light", "catppuccin"},
		{"gruvbox-light", "light", "gruvbox"},
		{"tokyonight", "dark", "tokyonight"},
		{"dracula", "dark", "dracula"},
		{"rosepine-dawn", "light", "rosepine"},
		{"kanagawa", "dark", "kanagawa"},
		{"kanagawa-lotus", "light", "kanagawa"},
		{"mocha", "dark", "catppuccin"},
		{"onedark", "dark", "onedark"},
		{"solarized-light", "light", "solarized"},
	}
	for _, c := range cases {
		got := Parse(c.in)
		if got.Polarity != c.pol || got.Palette != c.pal {
			t.Errorf("Parse(%q) = %+v, want pol=%s pal=%s", c.in, got, c.pol, c.pal)
		}
	}
}

func TestParseCOLORFGBG(t *testing.T) {
	cases := []struct {
		in   string
		dark bool
		ok   bool
	}{
		{"", true, false},
		{"15;0", true, true},
		{"0;15", false, true},
		{"0;7", false, true},
		{"default;0", true, true},
		{"nope", true, false},
		{"15;99", true, false},
	}
	for _, c := range cases {
		dark, ok := parseCOLORFGBG(c.in)
		if dark != c.dark || ok != c.ok {
			t.Errorf("parseCOLORFGBG(%q) = %v,%v want %v,%v", c.in, dark, ok, c.dark, c.ok)
		}
	}
}

func TestApplyOSC(t *testing.T) {
	auto := Parse("auto")
	got := auto.ApplyOSC(false)
	if got.Polarity != "light" {
		t.Fatalf("auto + light bg → %q", got.Polarity)
	}
	forced := Parse("dark").ApplyOSC(false)
	if forced.Polarity != "dark" {
		t.Fatalf("forced dark should ignore OSC, got %q", forced.Polarity)
	}
}

func TestResolveFlagWins(t *testing.T) {
	isolateFile(t)
	t.Setenv(EnvTTY, "light")
	got := Resolve("dark")
	if got.Polarity != "dark" {
		t.Fatalf("flag should win, got %+v", got)
	}
	got = Resolve("")
	if got.Polarity != "light" {
		t.Fatalf("env should apply, got %+v", got)
	}
}

func TestFromEnvCLITHEME(t *testing.T) {
	isolateFile(t)
	t.Setenv(EnvTTY, "")
	t.Setenv(EnvCLI, "light")
	got := FromEnv()
	if got.Polarity != "light" {
		t.Fatalf("CLITHEME fallback: %+v", got)
	}
}
