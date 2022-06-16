package test

import (
	"crypto/pkg/models"
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

var (
	url2 = "127.0.0.1:9999"
)

func TestFetch(t *testing.T) {
	_, err := Fetch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUsd(t *testing.T) {
	_, err := GetUsd()
	if err != nil {
		t.Fatal(err)
	}
}

func Fetch() (string, error) {
	url := fmt.Sprintf("http://%v/fetch", url2)

	resp, err := resty.New().R().
		Get(url)
	if err != nil {
		return "", fmt.Errorf("Fetch data failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("SetSetting failed. status: %v. message: %v", resp.StatusCode(), resp.String())
	}
	return "Saved data", nil
}

func GetUsd() ([]*models.Address, error) {
	url := fmt.Sprintf("http://%v/usd", url2)

	resp, err := resty.New().R().
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("Get usd data failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("SetSetting failed. status: %v. message: %v", resp.StatusCode(), resp.String())
	}

	return
}
