package main

import (
	"crypto/cmd"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/fetch", cmd.FetchData)
	e.GET("/deletables", cmd.DeleteTables)
	e.GET("/usd", cmd.GetUsd)
	err := e.Start(":9999")
	if err != nil {
		return
	}

}
