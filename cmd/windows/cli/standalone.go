//go:build windows

package cli

import (
	"github.com/spf13/cobra"

	"github.com/Kapparina/ThemeDetector/cmd/windows/app"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run ThemeDetector standalone",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunApp(debug)
	},
}
