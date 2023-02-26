# WebDriver client for Golang
![Build Status](https://github.com/hupe1980/gowebdriver/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/gowebdriver.svg)](https://pkg.go.dev/github.com/hupe1980/gowebdriver)
> Golang bindings that conform to the [W3C WebDriver](https://www.w3.org/TR/webdriver/) and [W3C WebDriver BiDi](https://w3c.github.io/webdriver-bidi/) standard for controlling web browsers.

## How to use
```go
chromeDriver, err := webdriver.NewChromeDriver("/path/to/chromedriver")
if err != nil {
	panic(err)
}

if err := chromeDriver.Start(); err != nil {
	panic(err)
}
defer chromeDriver.Stop()

session, err := chromeDriver.NewSession()
if err != nil {
    panic(err)
}
defer session.Close()

if err = session.NavigateTo("https://golang.org"); err != nil {
	 panic(err)
}
```

## Take Screenshots
```go
data, err := session.TakeScreenshot()
if err != nil {
	panic(err)
}

if err := os.WriteFile("./screenshot.png", data, 0600); err != nil {
	panic(err)
}
```

## BiDi Session
```go
biDiSession := session.BiDiSession()

bc, err := biDiSession.NewBrowsingContext(bidi.BrowsingContextTypeWindow, nil)
if err != nil {
	panic(err)
}

defer bc.Close()

navigation, err := bc.Navigate("https://golang.org", bidi.BrowsingContextReadinessStateComplete)
if err != nil {
	panic(err)
}
```

## License
[MIT](LICENCE)