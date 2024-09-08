//go:build windows

package app

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/charmbracelet/log"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	targetKey   = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	targetValue = "AppsUseLightTheme"

	lightTheme = "Catppuccin Latte"
	darkTheme  = "Catppuccin Mocha"
)

func IsDark() bool {
	key, err := GetKey(
		registry.CURRENT_USER,
		targetKey,
	)
	if err != nil {
		log.Fatal("Key retrieval error", "error", err)
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	value, err := GetIntValue(key, targetValue)
	if err != nil {
		log.Fatal("Value retrieval error", "error", err)
	}
	return strings.Contains(string(GetAppearance(value)), "Dark")
}

func GetKey(key registry.Key, path string) (k registry.Key, err error) {
	k, err = registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return k, errors.New("error opening key")
	}
	return
}

func GetIntValue(key registry.Key, name string) (v uint64, err error) {
	v, vType, err := key.GetIntegerValue(name)
	if err != nil {
		switch {
		case errors.Is(err, registry.ErrNotExist):
			return 0, errors.New("value does not exist")
		case errors.Is(err, registry.ErrUnexpectedType):
			return 0, fmt.Errorf("unexpected value type: %v", vType)
		default:
			return 0, errors.New("error getting value")
		}
	}
	return
}

func monitor(fn func(bool)) {
	changed := make(chan bool, 1)

	go func() {
		k, err := registry.OpenKey(
			registry.CURRENT_USER,
			targetKey,
			syscall.KEY_NOTIFY|registry.QUERY_VALUE,
		)
		if err != nil {
			log.Fatal("Key retrieval error", "error", err)
		}
		var wasDark = IsDark()

		for {
			err := windows.RegNotifyChangeKeyValue(
				windows.Handle(k),
				false,
				windows.CERT_CLOSE_STORE_FORCE_FLAG|windows.REG_NOTIFY_CHANGE_LAST_SET,
				windows.Handle(0),
				false,
			)
			if err != nil {
				log.Fatal("Error monitoring registry key", "error", err)
			}
			isDark := IsDark()
			if isDark != wasDark {
				wasDark = isDark
				changed <- isDark
			}
		}
	}()
	for {
		val := <-changed
		fn(val)
	}
}

func RunApp(isDebug bool) {
	if isDebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}
	log.Info("Monitoring appearance")
	monitor(func(isDark bool) {
		if isDark {
			log.Info("Dark mode enabled")
			if err := darkThemeCmd.SetBatTheme(); err != nil {
				log.Error("Error setting themes", "error", err)
			}
		} else {
			log.Info("Light mode enabled")
			if err := lightThemeCmd.SetBatTheme(); err != nil {
				log.Error("Error setting themes", "error", err)
			}
		}
	})
}
