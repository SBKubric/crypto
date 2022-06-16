package models

import (
	"encoding/json"
)

type NullString string

type UserStat struct {
	ItemId                int         `json:"itemId"`
	ID                    string      `json:"id"`
	Chain                 string      `json:"chain"`
	Name                  string      `json:"name"`
	SiteURL               string      `json:"site_url"`
	LogoURL               string      `json:"logo_url"`
	HasSupportedPortfolio bool        `json:"has_supported_portfolio"`
	Tvl                   float64     `json:"tvl"`
	NetUsdValue           float64     `json:"net_usd_value"`
	AssetUsdValue         json.Number `json:"asset_usd_value"`
	DebtUsdValue          json.Number `json:"debt_usd_value"`
}

type Address struct {
	Id         int    `json:"id"`
	Usd        string `json:"usd"`
	Addresses  string `json:"address"`
	Created_at string `json:"created_at"`
}
