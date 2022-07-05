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

func BinancePostRequest(json_data []byte) (response, float64, error) {
	resp_raw, err := http.Post("https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		return response{}, 0., err
	}

	var resp response
	json.NewDecoder(resp_raw.Body).Decode(&resp)

	price, err := strconv.ParseFloat(resp.Data[0].Adv.Price, 64)

	if err != nil {
		return response{}, 0., err
	}

	return resp, price, nil
}

func MonitorBinancePrice(orders chan<- Order) {
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
		resp_buy, price_buy, err := BinancePostRequest(json_buy_data)
		if err != nil {
			log.Fatal(err)
		}

		_, price_sell, err := BinancePostRequest(json_sell_data)

		if err != nil {
			log.Fatal(err)
		}

		orders <- Order{
			"Binance",
			resp_buy.Data[0].Advertiser["nickName"],
			price_buy,
			price_sell,
			resp_buy.Data[0].Adv.TradableQuantity,
			resp_buy.Data[0].Adv.MinSingleTransAmount,
			resp_buy.Data[0].Adv.MaxSingleTransAmount,
			fmt.Sprintf("https://p2p.binance.com/en/advertiserDetail?advertiserNo=%s", resp_buy.Data[0].Advertiser["userNo"]),
		}
	}
}
