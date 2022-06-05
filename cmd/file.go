package cmd

import (
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

	_, err = db.Query(`CREATE TABLE "address" (
    "id" SERIAL NOT NULL,
    "addresses" CHAR(1024)   NOT NULL,
    "created_at" date   NOT NULL,
    CONSTRAINT "pk_Address" PRIMARY KEY (
        "id"
     )
)`)
	if err != nil {
		log.Printf("Error1: %v", err)
	}

	_, err = db.Query(`CREATE TABLE "debank_api_results" (
		"ItemId"  SERIAL  NOT NULL,
		"ID" CHAR(1024)   NOT NULL,
		"Chain" CHAR(1024)   NOT NULL,
		"Name" CHAR(1024)   NOT NULL,
		"SiteURL" CHAR(1024)   NOT NULL,
		"LogoURL" CHAR(1024)   NOT NULL,
		"HasSupportedPortfolio" bool   NOT NULL,
		"Tvl" float   NOT NULL,
		"netUsdValue" float   NOT NULL,
		"AssetUsdValue" float   NOT NULL,
		"DebtUsdValue" float   NOT NULL,
		CONSTRAINT "pk_debank_api_results" PRIMARY KEY (
		"ItemId"
	)
	)`)

	if err != nil {
		log.Printf("Error2: %v", err)
	}

	//	_, err = db.Query(`ALTER TABLE "address" ADD CONSTRAINT "fk_Address_usd" FOREIGN KEY("usd")
	//REFERENCES "debank_api_results" ("ItemId");`)
	//
	//	if err != nil {
	//		log.Printf("Error3: %v", err)
	//	}

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

			_, err = db.Query(`INSERT INTO "address"(
							"addresses",
							"created_at"
						 )
						 VALUES($1, $2)`, key, time.Now())
			if err != nil {
				log.Printf("Errorrr: %v", err)
				fmt.Println(url)
			}
			for _, v := range value {

				_, err = db.Query(`INSERT INTO "debank_api_results"(
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
					 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, v.ID, v.Chain, v.Name, v.SiteURL, v.LogoURL, v.HasSupportedPortfolio,
					v.Tvl, v.netUsdValue, v.AssetUsdValue, v.DebtUsdValue)

				if err != nil {
					log.Printf("Error: %v", err)
					fmt.Println(url)
				}
			}
		}
	}

	return c.HTML(http.StatusOK, "ok")
}

/// ! use redis - 4 hours

/// подключиться к postgres, создав подключение к нему и таблицу
/// получить данные из внешнего api и сохранить в мапу
/// отправить данные в postgres
