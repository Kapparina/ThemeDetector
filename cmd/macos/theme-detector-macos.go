// /usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

const (
	lightTheme string = "Catppuccin Latte"
	darkTheme  string = "Catppuccin Mocha"
)

var (
	darkThemeCmd CommandSet = CommandSet{
		BatThemeCmd: []string{"set", "-Ux", "BAT_THEME", darkTheme},
	}
	lightThemeCmd CommandSet = CommandSet{
		BatThemeCmd: []string{"set", "-Ux", "BAT_THEME", lightTheme},
	}
	pList           string = filepath.Join(os.Getenv("HOME"), `/Library/Preferences/.GlobalPreferences.plist`)
	isDark, wasDark bool
)

type CommandSet struct {
	BatThemeCmd []string
}

func (c *CommandSet) SetBatTheme() error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "wsl", c.BatThemeCmd...).Run(); err != nil {
		return fmt.Errorf("error setting variable: %w", err)
	}
	return nil
}

func checkDarkMode() bool {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return false
		}
	}
	return true
}

func monitor(fn func(bool)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op.Has(fsnotify.Create) {
					isDark = checkDarkMode()
					if isDark && !wasDark {
						fn(isDark)
						wasDark = isDark
					}
					if !isDark && wasDark {
						fn(isDark)
						wasDark = isDark
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("A watcher error occurred", "error", err)
			}
		}
	}()
	err = watcher.Add(pList)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func setTheme(theme string) {
	switch theme {
	case "dark":
		if err := darkThemeCmd.SetBatTheme(); err != nil {
			log.Error("Error setting dark theme", "error", err)
		}
	case "light":
		if err := lightThemeCmd.SetBatTheme(); err != nil {
			log.Error("Error setting light theme", "error", err)
		}
	}
}

func main() {
	monitor(func(isDark bool) {
		if isDark {
			setTheme("dark")
		} else {
			setTheme("light")
		}
	})
}
