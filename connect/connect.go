// Package connect dials ttybus and subscribes to files + theme for
// companion panes.
package connect

import (
	"os"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/ttyzero/commons/theme"
	ttybus "github.com/ttyzero/ttybus/client"
)

// Session is a live ttybus connection with files and theme subscriptions.
type Session struct {
	Client *ttybus.Client
	Files  <-chan string
	Theme  <-chan string
}

// Dial connects to ttybus, starting the user daemon if needed.
func Dial() (*ttybus.Client, error) {
	c, err := ttybus.Dial("")
	if err == nil {
		return c, nil
	}
	if _, look := exec.LookPath("ttybus"); look != nil {
		return nil, err
	}
	_ = exec.Command("ttybus", "serve", "--daemon").Run()
	var last error
	d := 20 * time.Millisecond
	for i := 0; i < 10; i++ {
		c, last = ttybus.Dial("")
		if last == nil {
			return c, nil
		}
		time.Sleep(d)
		d *= 2
		if d > 200*time.Millisecond {
			d = 200 * time.Millisecond
		}
	}
	return nil, last
}

// Open hellos as name and subscribes to filesChannel (default "files")
// plus the shared theme channel.
func Open(name, filesChannel string) (*Session, error) {
	if filesChannel == "" {
		filesChannel = "files"
	}
	c, err := Dial()
	if err != nil {
		return nil, err
	}
	cwd, _ := os.Getwd()
	if err := c.Hello(ttybus.Info{Name: name, CWD: cwd}); err != nil {
		_ = c.Close()
		return nil, err
	}
	filesCh, err := c.Sub(filesChannel)
	if err != nil {
		_ = c.Close()
		return nil, err
	}
	themeCh, _ := c.Sub(theme.Channel)
	return &Session{Client: c, Files: filesCh, Theme: themeCh}, nil
}

func (s *Session) Close() error {
	if s == nil || s.Client == nil {
		return nil
	}
	return s.Client.Close()
}

func (s *Session) Chan(name string) <-chan string {
	if s == nil {
		return nil
	}
	if name == theme.Channel {
		return s.Theme
	}
	return s.Files
}

// Line is a payload from a named subscription.
type Line struct {
	Channel string
	Body    string
}

// Closed means that subscription's channel was closed.
type Closed struct {
	Channel string
}

// Wait reads the next line from ch as a tea.Msg.
func Wait(channel string, ch <-chan string) tea.Cmd {
	if ch == nil {
		return nil
	}
	return func() tea.Msg {
		s, ok := <-ch
		if !ok {
			return Closed{Channel: channel}
		}
		return Line{Channel: channel, Body: s}
	}
}

// ApplyLook applies a theme-bus command. Returns RequestBackgroundColor
// when THEME auto so the pane can re-query OSC 11.
func ApplyLook(cmd theme.Command, spec *theme.Spec, pal *theme.Palette, borderless *bool) tea.Cmd {
	if theme.Apply(cmd, spec, pal, borderless) {
		return tea.RequestBackgroundColor
	}
	return nil
}
