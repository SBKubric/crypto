package db

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

//var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "postgres"
)

type UserStat struct {
	ItemId                int         `json:"itemId"`
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

var adresses = []string{
	"0xba8a8f39b2315d4bc725c026ce3898c2c7e74f57",
	"0x2bd4284509bf6626d5def7ef20d4ca38ce71792e",
	"0x3ea91c76b176779d10cc2a27fd2687888886f0c2",
	"0xe8e94110e568fd45c8eb578bef0f36b5f154b794",
	"0x21bce0768110b9a8c50942be257637a843a7eac6",
	"0x9429614ccabfb2b24f444f33ede29d4575ebcdd1",
	"0x12244c23101f66741dae553c8836a9b2fd4e413a",
	"0x8c2753ee27ba890fbb60653d156d92e1c334f528",
}

func ConnectDb() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, nil
}

func CreateAddressTable() (*sql.DB, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Printf("Error connecting to db: %v", err)
	}

	_, err = db.Query(`CREATE TABLE "address" (
    "id" SERIAL  UNIQUE NOT NULL,
    "addresses" CHAR(1024) UNIQUE NOT NULL,
    "created_at" date   NOT NULL,
    CONSTRAINT "pk_Address" PRIMARY KEY (
        "id"
     ) )`)

	if err != nil {
		log.Printf("Error creating table address: %v", err)
	}
	return db, nil
}

func InsertAddresses(db *sql.DB) error {
	/// какой именно из запросов сохранятеся много раз
	for _, value := range adresses {
		_, err := db.Query(`INSERT INTO "address"(
							"addresses",
							"created_at"
						 )
						 VALUES($1, $2)`, value, time.Now())
		if err != nil {
			log.Printf("Error inserting data to table address: %v", err)
		}
	}
	return nil

}

func SaveAddress(c echo.Context) error {
	_, err := ConnectDb()
	if err != nil {
		log.Printf("Error connecting to db: %v", err)
	}

	db, err := CreateAddressTable()
	err = InsertAddresses(db)

	return c.String(http.StatusOK, "Saved addresses")
}

func CreateDebankTable() (*sql.DB, error) {
	db, err := ConnectDb()
	if err != nil {
		log.Printf("Error connecting to db: %v", err)
	}

	_, err = db.Query(`CREATE TABLE "debank_api_results" (
		"ItemId"  SERIAL   NOT NULL,
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
		addressId  INTEGER,
		CONSTRAINT "pk_debank_api_results" PRIMARY KEY (
		"ItemId"
	)
	)`)

	if err != nil {
		log.Printf("Error creating table debank_api_results: %v", err)
	}

	_, err = db.Query(`ALTER TABLE "debank_api_results" ADD CONSTRAINT "fk_Debank_usd" FOREIGN KEY("addressid")
	REFERENCES "address" ("id");`)
	if err != nil {
		log.Printf("Error adding foreign key: %v", err)
	}

	return db, nil
}

func InsertDebankResults(db *sql.DB) error {

	var m = make(map[string][]*UserStat)

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
			for _, v := range value {
				_, err := db.Query(`INSERT INTO debank_api_results(
				   "ID",
				   "Chain",
				   "Name",
				   "SiteURL",
				   "LogoURL",
				   "HasSupportedPortfolio",
				   "Tvl",
				   "netUsdValue",
				   "AssetUsdValue",
				   "DebtUsdValue",
				   "addressid"
			   )
			   VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, (SELECT id from address WHERE "addresses"= $11) )`,
					v.ID, v.Chain, v.Name, v.SiteURL, v.LogoURL, v.HasSupportedPortfolio, v.Tvl, v.netUsdValue, v.AssetUsdValue, v.DebtUsdValue, key)
				if err != nil {
					log.Printf("Error inserting data to debank table: %v", err)
					log.Printf("Info: %v", v)
				}
			}
		}
	}
	return nil
}

func SaveDebank(c echo.Context) error {
	db, err := CreateDebankTable()
	if err != nil {
		log.Printf("Couldn't creata table debank: %v", err)
	}

	err = InsertDebankResults(db)
	if err != nil {
		log.Printf("Couldn't insert data to table debank: %v", err)
	}
	return c.HTML(http.StatusOK, "ok")
}

/// ! use redis - 4 hours
