package webdriver

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/hupe1980/gowebdriver/bidi"
)

type WebDriver interface {
	// Start webDriver service
	Start() error

	// Stop webDriver service
	Stop() error

	// Query the server status
	Status() (*Status, error)

	Port() int

	// Create a new session
	NewSession(optFns ...func(o *SessionOptions)) (*Session, error)

	// Delete a session
	DeleteSession(id string) error
}

type Options struct {
	Port        int
	BootTimeout time.Duration
}

type webDriver struct {
	port    int
	service *Service
	timeout time.Duration
	client  *RestClient
}

func (w *webDriver) Start() error {
	if err := w.service.Start(); err != nil {
		return err
	}

	if err := w.service.WaitForBoot(w.timeout, func() bool {
		status, err := w.Status()
		if err != nil {
			return false
		}

		return status.Ready
	}); err != nil {
		_ = w.service.Stop()
		return err
	}

	return nil
}

func (w *webDriver) Stop() error {
	if err := w.service.Stop(); err != nil {
		return err
	}

	return nil
}

type Status struct {
	Build struct {
		// Version of driver
		Version string `json:"version"`
	} `json:"build"`
	Message string `json:"message"`
	OS      struct {
		// Operating system architecture
		Arch string `json:"arch"`
		// Name of operating system
		Name string `json:"name"`
		// Version of operating system
		Version string `json:"version"`
	} `json:"os"`
	Ready bool `json:"ready"`
}

func (w *webDriver) Status() (*Status, error) {
	data, err := w.client.Get("/status")
	if err != nil {
		return nil, err
	}

	status := &Status{}
	if err := json.Unmarshal(data, status); err != nil {
		return nil, err
	}

	return status, nil
}

func (w *webDriver) DeleteSession(id string) error {
	_, err := w.client.Delete(fmt.Sprintf("/session/%s", id))
	return err
}

func (w *webDriver) Port() int {
	return w.port
}

type SessionOptions struct {
	AlwaysMatch Capabilities
	FirstMatch  []Capabilities
}

func (w *webDriver) newSession(opts SessionOptions) (*Session, error) {
	params := Params{
		"alwaysMatch": opts.AlwaysMatch,
	}

	if opts.FirstMatch != nil {
		params["firstMatch"] = opts.FirstMatch
	}

	data, err := w.client.Post("/session", &Params{
		"capabilities": params,
	})
	if err != nil {
		return nil, err
	}

	session := &Session{}
	if err := json.Unmarshal(data, session); err != nil {
		return nil, err
	}

	session.client = w.client

	if session.IsBiDiSession() {
		biDiSession, err := bidi.New(session.Capabilities.WebSocketURL(), nil)
		if err != nil {
			return nil, err
		}

		session.biDiSession = biDiSession
	}

	return session, nil
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
