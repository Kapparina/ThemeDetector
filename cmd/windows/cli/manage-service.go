//go:build windows

package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc"

	"github.com/Kapparina/ThemeDetector/cmd/windows/constants"
	"github.com/Kapparina/ThemeDetector/cmd/windows/service"
)

var (
	manageCmd = &cobra.Command{
		Use:           "service { start | stop | pause | continue | install | uninstall }",
		Short:         "Manage service",
		Args:          cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: true,
		ValidArgs: []string{
			"start",
			"stop",
			"pause",
			"continue",
			"install",
			"uninstall",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "start":
				return service.Start(constants.ServiceName)
			case "stop":
				return service.ControlService(constants.ServiceName, svc.Stop, svc.Stopped)
			case "pause":
				return service.ControlService(constants.ServiceName, svc.Pause, svc.Paused)
			case "continue":
				return service.ControlService(constants.ServiceName, svc.Continue, svc.Running)
			case "install":
				return service.InstallService(constants.ServiceName, constants.ServiceDescription)
			case "uninstall":
				return service.RemoveService(constants.ServiceName)
			default:
				return errors.New("invalid argument")
			}
		},
	}
)
