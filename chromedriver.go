package webdriver

import (
	"fmt"
	"time"
)

type ChromeDriver struct {
	webDriver
	path    string
	port    int
	service *Service
	timeout time.Duration
}

func NewChromeDriver(path string, optFns ...func(o *Options)) (WebDriver, error) {
	opts := Options{
		Port:        0,
		BootTimeout: 10 * time.Second, //nolint gomnd
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Port == 0 {
		port, err := GetFreePort()
		if err != nil {
			return nil, err
		}

		opts.Port = port
	}

	service := NewService(path, []string{fmt.Sprintf("--port=%d", opts.Port)})

	cd := &ChromeDriver{
		path:    path,
		port:    opts.Port,
		service: service,
		timeout: opts.BootTimeout,
	}

	cd.client = NewRestClient(fmt.Sprintf("http://127.0.0.1:%d", opts.Port))

	return cd, nil
}

func (d *ChromeDriver) Start() error {
	if err := d.service.Start(); err != nil {
		return err
	}

	if err := d.service.WaitForBoot(d.timeout, func() bool {
		status, err := d.Status()
		if err != nil {
			return false
		}

		return status.Ready
	}); err != nil {
		_ = d.service.Stop()
		return err
	}

	return nil
}

func (d *ChromeDriver) Stop() error {
	if err := d.service.Stop(); err != nil {
		return err
	}

	return nil
}

func (d *ChromeDriver) NewSession(optFns ...func(o *SessionOptions)) (*Session, error) {
	opts := SessionOptions{
		AlwaysMatch: newDefaultChromeDriverCapabilities(),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return d.newSession(opts)
}

func newDefaultChromeDriverCapabilities() Capabilities {
	caps := NewCapabilities()
	caps.BrowserName("chrome")

	return caps
}
