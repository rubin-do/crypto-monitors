package markets

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type request struct {
	Page          int      `json:"page"`
	Rows          int      `json:"rows"`
	PayTypes      []string `json:"payTypes"`
	Countries     *string  `json:"countries"`
	PublisherType *string  `json:"publisherType"`
	Asset         string   `json:"asset"`
	Fiat          string   `json:"fiat"`
	TradeType     string   `json:"tradeType"`
}

type second struct {
	Price string
}

type first struct {
	Adv second
}

type response struct {
	Data []first
}

func MonitorBinancePrice(prices chan float64) {
	values := request{1, 1, []string{"Tinkoff"}, nil, nil, "USDT", "RUB", "BUY"}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := http.Post("https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", "application/json", bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err)
		}

		var resp_json response

		json.NewDecoder(resp.Body).Decode(&resp_json)

		price, err := strconv.ParseFloat(resp_json.Data[0].Adv.Price, 64)

		if err != nil {
			log.Fatal(err)
		}

		prices <- price
	}
}
