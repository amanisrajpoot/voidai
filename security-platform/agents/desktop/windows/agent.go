package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/security-platform/go-agent"
)

var elog debug.Log

type windowsService struct {
	agent *agent.Agent
}

func (m *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	// Start agent
	config := &agent.Config{
		ControlPlaneURL: os.Getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL"),
		AuthKey:         os.Getenv("SECURITY_PLATFORM_AUTH_KEY"),
		ServiceName:     "security-platform-agent",
		Environment:     "production",
	}

	ag, err := agent.NewAgent(config)
	if err != nil {
		elog.Error(1, fmt.Sprintf("Failed to create agent: %v", err))
		return false, 1
	}

	if err := ag.Start(); err != nil {
		elog.Error(1, fmt.Sprintf("Failed to start agent: %v", err))
		return false, 1
	}

	m.agent = ag
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				elog.Info(1, "Stopping service")
				ag.Stop()
				break loop
			default:
				elog.Error(1, fmt.Sprintf("Unexpected control request #%d", c))
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return false, 0
}

func runService(name string, isDebug bool) error {
	elog, err := eventlog.Open(name)
	if err != nil {
		return err
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("Starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &windowsService{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return err
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
	return nil
}

func installService(name, desc string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}
	s, err = m.CreateService(name, exePath, mgr.Config{
		DisplayName:      desc,
		StartType:       mgr.StartAutomatic,
		ServiceControls: mgr.ServiceControls{},
	}, "is", "auto-started")
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func removeService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

func main() {
	var (
		install   = flag.Bool("install", false, "Install the service")
		remove    = flag.Bool("remove", false, "Remove the service")
		debugFlag = flag.Bool("debug", false, "Run in debug mode")
	)
	flag.Parse()

	var err error
	if *install {
		err = installService("SecurityPlatformAgent", "Security Platform Agent")
	} else if *remove {
		err = removeService("SecurityPlatformAgent")
	} else {
		err = runService("SecurityPlatformAgent", *debugFlag)
	}
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
}
