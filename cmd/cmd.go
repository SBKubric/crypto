package cmd

import (
	"crypto/pkg/db"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func FetchData(c echo.Context) error {
	_, err := db.ConnectDb()
	if err != nil {
		log.Printf("Error connecting to db: %v", err)
		return err
	}

	err = db.SaveAddress(c)
	err = db.SaveDebank(c)

	return c.HTML(http.StatusOK, "Saved data")
}

func GetUsd(c echo.Context) error {
	_, err := db.ConnectDb()
	if err != nil {
		log.Printf("Error connecting to db: %v", err)
		return err
	}

	addresses, err := db.GetUsd()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, addresses)
}

func DeleteTables(c echo.Context) error {
	_, err := db.ConnectDb()
	if err != nil {
		log.Printf("Error connecting to db: %v", err)
		return err
	}

	err = db.DeleteTables()

	return c.HTML(http.StatusOK, "Deleted table")

}

/// ! use redis - 4 hours
