package cmd

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "postgres"
)

type AutoGenerate struct {
	ID                    string  `json:"id"`
	Chain                 string  `json:"chain"`
	Name                  string  `json:"name"`
	SiteURL               string  `json:"site_url"`
	LogoURL               string  `json:"logo_url"`
	HasSupportedPortfolio bool    `json:"has_supported_portfolio"`
	Tvl                   float64 `json:"tvl"`
	netUsdValue           float64 `json:"net_usd_value"`
	AssetUsdValue         float64 `json:"asset_usd_value"`
	DebtUsdValue          int     `json:"debt_usd_value"`
}

func (t *AutoGenerate) Dump(indent string) {
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

	m := make(map[string]interface{})

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS "address"("usd" CHAR, "adresses" CHAR, "created_at" CHAR)`)
	if err != nil {
		return err
	}
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS "debank_api_results"("ID" CHAR, "Chain" CHAR, "Name" CHAR,
"SiteURL" CHAR, "LogoURL" CHAR, "HasSupportedPortfolio" boolean, "Tvl" FLOAT, "netUsdValue" FLOAT, "AssetUsdValue" FLOAT, "DebtUsdValue" INTEGER)`)
	if err != nil {
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

		sb := string(body)
		m[value] = sb

		b := new(bytes.Buffer)
		e := gob.NewEncoder(b)

		// Encoding the map
		err = e.Encode(m)
		if err != nil {
			panic(err)
		}

		var decodedMap map[string]interface{}
		d := gob.NewDecoder(b)

		// Decoding the serialized data
		err = d.Decode(&decodedMap)
		if err != nil {
			panic(err)
		}

		for key, value := range decodedMap {
			fmt.Println("Key:", key, "Value:", value)

			for k, v := range value {
				fmt.Print(k)
				fmt.Print(v)
			}
		}

	}

	return c.JSON(http.StatusOK, "ok")
}

/// ! use redis - 4 hours

/// подключиться к postgres, создав подключение к нему и таблицу
/// получить данные из внешнего api и сохранить в мапу
/// отправить данные в postgres
