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

type PayMethod struct {
	PayType string
}

type second struct {
	Price                string
	TradableQuantity     string
	TradeMethods         []PayMethod
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

func binancePostRequest(jsonData []byte) (response, float64, error) {
	respRaw, err := http.Post("https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return response{}, 0., err
	}

	var resp response
	json.NewDecoder(respRaw.Body).Decode(&resp)

	price, err := strconv.ParseFloat(resp.Data[0].Adv.Price, 64)

	if err != nil {
		return response{}, 0., err
	}

	return resp, price, nil
}

func parsePaymentMethods(order second) string {
	var paymentMethods string
	for _, method := range order.TradeMethods {
		paymentMethods += method.PayType
	}
	return paymentMethods
}

func MonitorBinancePrice(orders chan<- Order) {
	buyValues := request{1, 1, []string{"Tinkoff", "RosBank", "RaiffeisenBankRussia"}, nil, nil, "USDT", "RUB", "BUY"}
	jsonBuyData, err := json.Marshal(buyValues)

	if err != nil {
		log.Fatal(err)
	}

	sellValues := request{1, 1, []string{}, nil, nil, "USDT", "RUB", "SELL"}
	jsonSellData, err := json.Marshal(sellValues)

	if err != nil {
		log.Fatal(err)
	}

	for {
		respBuy, priceBuy, err := binancePostRequest(jsonBuyData)
		if err != nil {
			log.Fatal(err)
		}

		_, priceSell, err := binancePostRequest(jsonSellData)

		if err != nil {
			log.Fatal(err)
		}

		orders <- Order{
			"Null",
			"Binance",
			respBuy.Data[0].Advertiser["nickName"],
			priceBuy,
			priceSell,
			respBuy.Data[0].Adv.TradableQuantity,
			respBuy.Data[0].Adv.MinSingleTransAmount,
			respBuy.Data[0].Adv.MaxSingleTransAmount,
			parsePaymentMethods(respBuy.Data[0].Adv),
			fmt.Sprintf("https://p2p.binance.com/en/advertiserDetail?advertiserNo=%s", respBuy.Data[0].Advertiser["userNo"]),
		}
	}
}
