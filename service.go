package webdriver

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

type CheckStatusFunc func() bool

type Service struct {
	path string
	args []string
	cmd  *exec.Cmd
}

func NewService(path string, args []string) *Service {
	return &Service{
		path: path,
		args: args,
	}
}

func (s *Service) Start() error {
	if s.cmd != nil {
		return errors.New("webdriver already running")
	}

	s.cmd = exec.Command(s.path, s.args...) // nolint gosec

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := s.cmd.Start(); err != nil {
		return err
	}

	go func() {
		_, _ = io.Copy(os.Stdout, stdout)
	}()
	go func() {
		_, _ = io.Copy(os.Stderr, stderr)
	}()

	return nil
}

func (s *Service) Stop() error {
	if s.cmd == nil {
		return errors.New("webDriver not running")
	}

	if runtime.GOOS == "windows" {
		if err := s.cmd.Process.Kill(); err != nil {
			return err
		}
	} else {
		if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) WaitForBoot(timeout time.Duration, fn CheckStatusFunc) error {
	timeoutChan := time.After(timeout)

	failedChan := make(chan struct{}, 1)
	startedChan := make(chan struct{})

	go func() {
		up := fn()
		for !up {
			select {
			case <-failedChan:
				return
			default:
				time.Sleep(500 * time.Millisecond) // nolint gomnd

				up = fn()
			}
		}
		startedChan <- struct{}{}
	}()

	select {
	case <-timeoutChan:
		failedChan <- struct{}{}
		return errors.New("failed to start before timeout")
	case <-startedChan:
		return nil
	}
}
