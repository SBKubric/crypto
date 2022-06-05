package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "postgres"
)

type UserStat struct {
	ID                    string      `json:"id"`
	Chain                 string      `json:"chain"`
	Name                  string      `json:"name"`
	SiteURL               string      `json:"site_url"`
	LogoURL               string      `json:"logo_url"`
	HasSupportedPortfolio bool        `json:"has_supported_portfolio"`
	Tvl                   float64     `json:"tvl"`
	netUsdValue           float64     `json:"net_usd_value"`
	AssetUsdValue         float64     `json:"asset_usd_value"`
	DebtUsdValue          json.Number `json:"debt_usd_value"`
}

func (t *UserStat) Dump(indent string) {
	fmt.Println(indent+"id:", t.ID)
	fmt.Print(indent+"chain: ", t.Chain)
	fmt.Print(indent+"Name: ", t.Name)

}

func GetInfo(c echo.Context) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	adresses := []string{
		"0xba8a8f39b2315d4bc725c026ce3898c2c7e74f57",
		"0x2bd4284509bf6626d5def7ef20d4ca38ce71792e",
		"0x3ea91c76b176779d10cc2a27fd2687888886f0c2",
		"0xe8e94110e568fd45c8eb578bef0f36b5f154b794",
		"0x21bce0768110b9a8c50942be257637a843a7eac6",
		"0x9429614ccabfb2b24f444f33ede29d4575ebcdd1",
		"0x12244c23101f66741dae553c8836a9b2fd4e413a",
		"0x8c2753ee27ba890fbb60653d156d92e1c334f528",
	}

	m := make(map[string][]*UserStat)

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS "address"("usd" CHAR(1024), "adresses" CHAR(1024), "created_at" CHAR(1024), CONSTRAINT "pk_Address" PRIMARY KEY ("usd"))`)
	if err != nil {
		return err
	}
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS "debank_api_results"("ID" CHAR(1024), "Chain" CHAR(1024), "Name" CHAR(1024),
"SiteURL" CHAR(1024), "LogoURL" CHAR(1024), "HasSupportedPortfolio" boolean, "Tvl" CHAR(1024), "netUsdValue" CHAR(1024), "AssetUsdValue" CHAR(1024), "DebtUsdValue" CHAR(1024),
    CONSTRAINT "pk_Debank" PRIMARY KEY ("AssetUsdValue"))`)
	if err != nil {
		return err
	}

	_, err = db.Query(`ALTER TABLE "address" ADD CONSTRAINT "fk_Address" FOREIGN KEY ("usd") REFERENCES "debank_api_results" ("AssetUsdValue")`)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	fmt.Println("Table created")

	for _, value := range adresses {
		url := fmt.Sprintf("https://openapi.debank.com/v1/user/simple_protocol_list?id=%s", value)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var userStat []*UserStat

		err = json.Unmarshal(body, &userStat)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("--------\n", userStat)
		}

		m[value] = userStat

		for key, value := range m {

			var AssetUsdValue string

			cur := time.Now()
			for _, v := range value {

				ctx := context.Background()

				err := db.QueryRowContext(ctx, `INSERT INTO "debank_api_results"(
							"ID",
							"Chain",
							"Name",
							"SiteURL",
						    "LogoURL",
							"HasSupportedPortfolio",
						    "Tvl",
						    "netUsdValue",
							"AssetUsdValue",
		                   	"DebtUsdValue"
						 )
						 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING "netUsdValue"`, v.ID, v.Chain, v.Name, v.SiteURL, v.LogoURL, v.HasSupportedPortfolio,
					v.Tvl, v.netUsdValue, v.AssetUsdValue, v.DebtUsdValue).Scan(&AssetUsdValue)
				fmt.Println(AssetUsdValue)

				if err != nil {
					log.Printf("Error: %v", err)
					fmt.Println(url)
					return err
				}
			}
			_, err = db.Query(`INSERT INTO "address"(
							"usd",
							"adresses",
							"created_at"
						 )
						 VALUES($1, $2, $3)`, AssetUsdValue, key, cur)
			if err != nil {
				log.Printf("Errorrr: %v", err)
				return err
			}
		}
	}

	return c.JSON(http.StatusOK, "ok")
}

/// ! use redis - 4 hours

/// подключиться к postgres, создав подключение к нему и таблицу
/// получить данные из внешнего api и сохранить в мапу
/// отправить данные в postgres
