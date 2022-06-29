package markets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type request_by_bit struct {
	UserId     string `json:"userId"`
	TokenId    string `json:"tokenId"`
	CurrencyId string `json:"currencyId"`
	Payment    string `json:"payment"`
	Side       string `json:"side"`
	Size       string `json:"size"`
	Page       string `json:"page"`
	Amount     string `json:"amount"`
}

type order_by_bit struct {
	Price string
}

type results struct {
	Items []order_by_bit
}

type response_by_bit struct {
	Result results
}

func MonitorByBitPrice(prices chan float64) {
	values := request_by_bit{"", "USDT", "RUB", "75", "13", "10", "1", ""}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := http.Post("https://api2.bybit.com/spot/api/otc/item/list", "application/json", bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp)

		var response_json response_by_bit

		json.NewDecoder(resp.Body).Decode(&response_json)

		fmt.Println(response_json.Result.Items[0].Price)
	}

}
