package webdriver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapabilities(t *testing.T) {
	caps := Capabilities{}
	caps.SetBrowserName("chrome")
	assert.Equal(t, "chrome", caps["browserName"])
}

func TestChromeOptions(t *testing.T) {
	co := ChromeOptions{}
	co.AddArg("--headless")
	assert.ElementsMatch(t, co["args"], []string{"--headless"})

	co.AddArg("--disable-blink-features=AutomationControlled")
	assert.ElementsMatch(t, co["args"], []string{"--headless", "--disable-blink-features=AutomationControlled"})
}
