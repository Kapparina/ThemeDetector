//go:build windows

package constants

const (
	ServiceName        string = "ThemeDetector"
	ServiceDescription string = "Detect whether the system theme is light or dark, and adjust $BAT_THEME accordingly."
)
