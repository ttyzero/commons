package files

import "testing"

func TestParsePath(t *testing.T) {
	ev, err := Parse("  /tmp/foo.go  ")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Path != "/tmp/foo.go" || ev.Line != 0 {
		t.Fatalf("%+v", ev)
	}
}

func TestParseLineCol(t *testing.T) {
	ev, err := Parse("/tmp/foo.go:12:4")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Path != "/tmp/foo.go" || ev.Line != 12 || ev.Col != 4 {
		t.Fatalf("%+v", ev)
	}
	ev, err = Parse("file:///tmp/foo.go:3")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Path != "/tmp/foo.go" || ev.Line != 3 {
		t.Fatalf("%+v", ev)
	}
}

func TestParseWindowsishColon(t *testing.T) {
	// A path with no trailing number stays intact.
	ev, err := Parse("/tmp/foo:bar.go")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Path != "/tmp/foo:bar.go" || ev.Line != 0 {
		t.Fatalf("%+v", ev)
	}
}

func TestParseJSON(t *testing.T) {
	ev, err := Parse(`{"path":"/code/a.go","line":9,"origin":"vim","root":"/code"}`)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Path != "/code/a.go" || ev.Line != 9 || ev.Origin != "vim" || ev.Root != "/code" {
		t.Fatalf("%+v", ev)
	}
}

func TestParseEmpty(t *testing.T) {
	ev, err := Parse("")
	if err != nil || ev.Path != "" {
		t.Fatalf("%+v %v", ev, err)
	}
}
