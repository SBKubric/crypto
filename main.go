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
	//e.GET("/saveaddress", db.SaveAddress)
	//e.GET("/savedebank", db.SaveDebank)
	err := e.Start(":1323")
	if err != nil {
		return
	}

}
