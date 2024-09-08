//go:build windows

package service

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

type themeService struct{}

func (t *themeService) Execute(
	args []string,
	req <-chan svc.ChangeRequest,
	changes chan<- svc.Status,
) (
	ssec bool,
	errno uint32,
) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-req:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				output := strings.Join(args, "-")
				output += fmt.Sprintf("-%d", c.Context)
				_ = elog.Info(1, output)
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				_ = elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		if elog, err = eventlog.Open(name); err != nil {
			return
		}
	}
	defer func(elog debug.Log) {
		_ = elog.Close()
	}(elog)

	_ = elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	if err = run(name, &themeService{}); err != nil {
		_ = elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	_ = elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
