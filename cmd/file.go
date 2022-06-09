package cmd

import (
	"crypto/db"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func GetInfo(c echo.Context) error {
	_, err := db.ConnectDb()
	if err != nil {
		log.Printf("Error connecting to db: %v", err)
		return err
	}

	err = db.SaveAddress(c)
	err = db.SaveDebank(c)

	return c.HTML(http.StatusOK, "ok")
}

/// ! use redis - 4 hours
