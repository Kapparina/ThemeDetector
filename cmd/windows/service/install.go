//go:build windows

package service

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func getExecPath() (exePath string, err error) {
	exePath, err = os.Executable()
	if err != nil {
		return
	}
	return filepath.EvalSymlinks(exePath)
}

func InstallService(name, desc string) error {
	var (
		m   *mgr.Mgr
		s   *mgr.Service
		err error
	)
	exePath, err := getExecPath()
	if err != nil {
		return err
	}
	if m, err = mgr.Connect(); err != nil {
		return err
	}
	defer func(m *mgr.Mgr) {
		_ = m.Disconnect()
	}(m)
	if s, err = m.OpenService(name); err == nil {
		_ = s.Close()
		return fmt.Errorf("service %s already exists", name)
	}
	s, err = m.CreateService(name, exePath, mgr.Config{DisplayName: desc, StartType: mgr.StartAutomatic})
	if err != nil {
		return err
	}
	defer func(s *mgr.Service) {
		_ = s.Close()
	}(s)
	if err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {
		return err
	}
	return nil
}

func RemoveService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer func(m *mgr.Mgr) {
		_ = m.Disconnect()
	}(m)

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s not found", name)
	}
	defer func(s *mgr.Service) {
		_ = s.Close()
	}(s)

	if err = s.Delete(); err != nil {
		return err
	}
	if err = eventlog.Remove(name); err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %v", err)
	}
	return nil
}
