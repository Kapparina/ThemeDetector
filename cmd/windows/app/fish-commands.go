//go:build windows

package app

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type CommandSet struct {
	BatThemeCmd  []string
	FishThemeCmd []string
}

func (c *CommandSet) SetBatTheme() error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "wsl", c.BatThemeCmd...).Run(); err != nil {
		return fmt.Errorf("error setting variable: %w", err)
	}
	return nil
}

func (c *CommandSet) SetFishTheme() error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "wsl", c.FishThemeCmd...).Run(); err != nil {
		return fmt.Errorf("error setting shell theme: %w", err)
	}
	return nil
}

func (c *CommandSet) SetThemes() error {
	if err := c.SetBatTheme(); err != nil {
		return err
	}
	if err := c.SetFishTheme(); err != nil {
		return err
	}
	return nil

}

var (
	darkThemeCmd = CommandSet{
		BatThemeCmd:  []string{"set", "-Ux", "BAT_THEME", darkTheme},
		FishThemeCmd: []string{"fish_config", "theme", "choose", darkTheme},
	}
	lightThemeCmd = CommandSet{
		BatThemeCmd:  []string{"set", "-Ux", "BAT_THEME", lightTheme},
		FishThemeCmd: []string{"fish_config", "theme", "choose", lightTheme},
	}
)
