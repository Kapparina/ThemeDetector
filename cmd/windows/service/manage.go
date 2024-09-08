//go:build windows

package service

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func Start(name string) error {
	var m *mgr.Mgr
	var s *mgr.Service
	var err error

	if m, err = mgr.Connect(); err != nil {
		return err
	}
	defer func(m *mgr.Mgr) {
		_ = m.Disconnect()
	}(m)

	if s, err = m.OpenService(name); err != nil {
		return fmt.Errorf("could not access service: %w", err)
	}
	defer func(s *mgr.Service) {
		_ = s.Close()
	}(s)

	if err = s.Start("is", "manual-started"); err != nil {
		return fmt.Errorf("could not start service: %w", err)
	}
	return nil
}

func ControlService(name string, command svc.Cmd, newState svc.State) error {
	var m *mgr.Mgr
	var s *mgr.Service
	var status svc.Status
	var err error

	if m, err = mgr.Connect(); err != nil {
		return err
	}
	defer func(m *mgr.Mgr) {
		_ = m.Disconnect()
	}(m)

	if s, err = m.OpenService(name); err != nil {
		return fmt.Errorf("could not access service: %w", err)
	}
	defer func(s *mgr.Service) {
		_ = s.Close()
	}(s)
	if status, err = s.Control(command); err != nil {
		return fmt.Errorf("could not send control=%d: %v", command, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != newState {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to move to state=%d", newState)
		}
		time.Sleep(300 * time.Millisecond)
		if status, err = s.Query(); err != nil {
			return fmt.Errorf("could not query service: %w", err)
		}
	}
	return nil
}
