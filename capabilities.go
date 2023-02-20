package webdriver

type Capabilities map[string]interface{}

// NewCapabilities returns a Capabilities instance with any provided features enabled.
func NewCapabilities(features ...string) Capabilities {
	c := Capabilities{}

	for _, feature := range features {
		c.With(feature)
	}

	return c
}

// BrowserName sets the desired browser name.
func (c Capabilities) BrowserName(name string) Capabilities {
	c["browserName"] = name
	return c
}

// A ProxyConfig instance defines the desired proxy configuration the WebDriver
// should use to proxy a Page.
//
// See: https://www.w3.org/TR/webdriver/#proxy
type ProxyConfig struct {
	ProxyType          string `json:"proxyType"`
	ProxyAutoconfigURL string `json:"proxyAutoconfigUrl,omitempty"`
	FTPProxy           string `json:"ftpProxy,omitempty"`
	HTTPProxy          string `json:"httpProxy,omitempty"`
	SSLProxy           string `json:"sslProxy,omitempty"`
	SOCKSProxy         string `json:"socksProxy,omitempty"`
	SOCKSUsername      string `json:"socksUsername,omitempty"`
	SOCKSPassword      string `json:"socksPassword,omitempty"`
	NoProxy            string `json:"noProxy,omitempty"`
}

// Proxy sets the desired proxy configuration.
func (c Capabilities) Proxy(p ProxyConfig) Capabilities {
	c["proxy"] = p
	return c
}

// BrowserVersion sets the desired browser version.
func (c Capabilities) BrowserVersion(version string) Capabilities {
	c["browserVersion"] = version
	return c
}

// PlatformName sets the desired browser platform.
func (c Capabilities) PlatformName(platform string) Capabilities {
	c["platformName"] = platform
	return c
}

// With enables the provided feature (ex. "acceptInsecureCerts").
func (c Capabilities) With(feature string) Capabilities {
	c[feature] = true
	return c
}

// Without disables the provided feature (ex. "acceptInsecureCerts").
func (c Capabilities) Without(feature string) Capabilities {
	c[feature] = false
	return c
}

// Sets an arbitrary key-value pair
func (c Capabilities) Set(key string, value string) Capabilities {
	c[key] = value
	return c
}
