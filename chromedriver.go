package webdriver

import (
	"fmt"
	"time"
)

type chromeDriver struct {
	webDriver
	path string
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

	cd := &chromeDriver{
		path: path,
	}

	cd.port = opts.Port
	cd.service = NewService(path, []string{fmt.Sprintf("--port=%d", opts.Port)})
	cd.timeout = opts.BootTimeout
	cd.client = NewRestClient(fmt.Sprintf("http://127.0.0.1:%d", opts.Port))

	return cd, nil
}

func (d *chromeDriver) NewSession(optFns ...func(o *SessionOptions)) (*Session, error) {
	opts := SessionOptions{
		AlwaysMatch: newDefaultChromeDriverCapabilities(),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return d.newSession(opts)
}

func newDefaultChromeDriverCapabilities() Capabilities {
	caps := Capabilities{}
	caps.SetBrowserName("chrome")
	caps.SetWebSocketURL(true)

	return caps
}
