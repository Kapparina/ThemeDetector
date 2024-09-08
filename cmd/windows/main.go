//go:build windows

package main

import (
	"github.com/charmbracelet/log"
	"golang.org/x/sys/windows/svc"

	"github.com/Kapparina/ThemeDetector/cmd/windows/cli"
	"github.com/Kapparina/ThemeDetector/cmd/windows/constants"
	"github.com/Kapparina/ThemeDetector/cmd/windows/service"
)

func main() {
	if inService, err := svc.IsWindowsService(); err != nil {
		log.Fatalf("Error checking if running as service: %v", err)
	} else if inService {
		service.RunService(constants.ServiceName, false)
	} else {
		cli.Execute()
	}
}
