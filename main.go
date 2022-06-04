package main

import (
	"crypto/cmd"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/fetch", cmd.GetInfo)
	err := e.Start(":1323")
	if err != nil {
		return
	}

}
