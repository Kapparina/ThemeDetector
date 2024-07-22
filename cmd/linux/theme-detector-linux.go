// /usr/bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

const (
	lightTheme = "Catppuccin Latte"
	darkTheme  = "Catppuccin Mocha"
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

var (
	darkThemeCmd CommandSet = CommandSet{
		BatThemeCmd: []string{"set", "-Ux", "BAT_THEME", darkTheme},
	}
)

func monitor(fn func(bool)) { // TODO: fill this

}

func main() { // TODO: fill this

}
