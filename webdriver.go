package webdriver

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type WebDriver interface {
	// Start webDriver service
	Start() error

	// Stop webDriver service
	Stop() error

	// Query the server status
	Status() (*Status, error)

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
	client *RestClient
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
