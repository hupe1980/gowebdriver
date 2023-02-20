# WebDriver client for Golang
![Build Status](https://github.com/hupe1980/gowebdriver/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/gowebdriver.svg)](https://pkg.go.dev/github.com/hupe1980/gowebdriver)
> Golang bindings that conform to the [W3C WebDriver standard](https://www.w3.org/TR/webdriver) for controlling web browsers.

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

## License
[MIT](LICENCE)