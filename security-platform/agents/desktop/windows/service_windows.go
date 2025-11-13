//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

const serviceName = "SecurityPlatformAgent"

type windowsService struct {
	agent *agent.Agent
}

func (ws *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	changes <- svc.Status{State: svc.StartPending}

	// Start agent
	config := loadConfig("")
	ws.agent = agent.NewAgent(config)
	if err := ws.agent.Start(); err != nil {
		log.Printf("Failed to start agent: %v", err)
		changes <- svc.Status{State: svc.Stopped}
		return false, 1
	}

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				ws.agent.Stop()
				break loop
			case svc.Interrogate:
				changes <- c.CurrentStatus
			default:
				log.Printf("Unexpected control request: %v", c.Cmd)
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return false, 0
}

func installWindowsService() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", serviceName)
	}

	s, err = m.CreateService(serviceName, exePath, mgr.Config{
		DisplayName: "Security Platform Agent",
		Description: "Security Platform Desktop Agent for Windows",
		StartType:   mgr.StartAutomatic,
	})
	if err != nil {
		return err
	}
	defer s.Close()

	err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %v", err)
	}

	return nil
}

func uninstallWindowsService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", serviceName)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		return err
	}

	err = eventlog.Remove(serviceName)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %v", err)
	}

	return nil
}
