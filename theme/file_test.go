package theme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLoadRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "theme")
	t.Setenv(EnvFile, path)
	t.Setenv(EnvTTY, "")
	t.Setenv(EnvCLI, "")
	t.Setenv(EnvBorderless, "")

	spec := Parse("nord")
	if err := Save(spec, true); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "THEME nord BORDERS=0") {
		t.Fatalf("body %q", body)
	}

	got, ok := Load()
	if !ok || !got.HasSpec || got.Spec.Name != "nord" || !got.HasBorders || !got.Borderless {
		t.Fatalf("%+v ok=%v", got, ok)
	}
	if Resolve("").Name != "nord" {
		t.Fatalf("Resolve from file: %+v", Resolve(""))
	}
	if !Borderless(false) {
		t.Fatal("file BORDERS=0 should be borderless")
	}
}

func TestFileComments(t *testing.T) {
	path := filepath.Join(t.TempDir(), "theme")
	t.Setenv(EnvFile, path)
	if err := os.WriteFile(path, []byte("# hi\n\nTHEME catppuccin\nBORDERS=1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, ok := Load()
	if !ok || got.Spec.Name != "catppuccin" || !got.HasBorders || got.Borderless {
		t.Fatalf("%+v", got)
	}
}

func TestEnvBeatsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "theme")
	t.Setenv(EnvFile, path)
	if err := Save(Parse("nord"), true); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvTTY, "dracula")
	t.Setenv(EnvBorderless, "0")
	if Resolve("").Name != "dracula" {
		t.Fatalf("%+v", Resolve(""))
	}
	if Borderless(false) {
		t.Fatal("TTYBORDERLESS=0 should keep frame")
	}
}

func TestMissingFile(t *testing.T) {
	t.Setenv(EnvFile, filepath.Join(t.TempDir(), "nope"))
	t.Setenv(EnvTTY, "")
	t.Setenv(EnvCLI, "")
	if _, ok := Load(); ok {
		t.Fatal("missing")
	}
	if Resolve("").Polarity != "auto" {
		t.Fatalf("%+v", Resolve(""))
	}
}

func TestEncode(t *testing.T) {
	if Encode(Parse("nord"), false) != "THEME nord BORDERS=1" {
		t.Fatal(Encode(Parse("nord"), false))
	}
	if Encode(Parse("auto"), true) != "THEME auto BORDERS=0" {
		t.Fatal(Encode(Parse("auto"), true))
	}
}
