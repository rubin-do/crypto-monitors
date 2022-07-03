package markets

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Price                string
	TradableQuantity     string
	MinSingleTransAmount string
	MaxSingleTransAmount string
}

type first struct {
	Adv        second
	Advertiser map[string]string
}

type response struct {
	Data []first
}

func MonitorBinancePrice(orders chan Order) {
	buy_values := request{1, 1, []string{"Tinkoff", "RosBank"}, nil, nil, "USDT", "RUB", "BUY"}
	json_buy_data, err := json.Marshal(buy_values)

	if err != nil {
		log.Fatal(err)
	}

	sell_values := request{1, 1, []string{}, nil, nil, "USDT", "RUB", "SELL"}
	json_sell_data, err := json.Marshal(sell_values)

	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := http.Post("https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", "application/json", bytes.NewBuffer(json_buy_data))

		if err != nil {
			log.Fatal(err)
		}

		var resp_json response

		json.NewDecoder(resp.Body).Decode(&resp_json)

		buy_price, err := strconv.ParseFloat(resp_json.Data[0].Adv.Price, 64)

		if err != nil {
			log.Fatal(err)
		}

		resp, err = http.Post("https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", "application/json", bytes.NewBuffer(json_sell_data))

		if err != nil {
			log.Fatal(err)
		}

		var resp_json_sell response

		json.NewDecoder(resp.Body).Decode(&resp_json_sell)

		sell_price, err := strconv.ParseFloat(resp_json_sell.Data[0].Adv.Price, 64)

		if err != nil {
			log.Fatal(err)
		}

		orders <- Order{
			"Binance",
			resp_json.Data[0].Advertiser["nickName"],
			buy_price,
			sell_price,
			resp_json.Data[0].Adv.TradableQuantity,
			resp_json.Data[0].Adv.MaxSingleTransAmount,
			resp_json.Data[0].Adv.MinSingleTransAmount,
			fmt.Sprintf("https://p2p.binance.com/en/advertiserDetail?advertiserNo=%s", resp_json.Data[0].Advertiser["userNo"]),
		}
	}
}
