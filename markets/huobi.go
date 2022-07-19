package markets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type huobiResponse struct {
	Data []struct {
		UID        int    `json:"uid"`
		UserName   string `json:"userName"`
		CoinID     int    `json:"coinId"`
		Currency   int    `json:"currency"`
		PayMethods []struct {
			Name string `json:"name"`
		} `json:"payMethods"`
		MinTradeLimit string `json:"minTradeLimit"`
		MaxTradeLimit string `json:"maxTradeLimit"`
		Price         string `json:"price"`
		TradeCount    string `json:"tradeCount"`
	} `json:"data"`
}

const (
	huobiDataUrl = "https://otc-api.bitderiv.com/v1/data/trade-market?coinId=2&currency=11&tradeType=%s&currPage=1&payMethod=0&acceptOrder=0&country=&blockType=general&online=1&range=0&amount=&onlyTradable=false"
)

func huobiRequest(url string, orderType string) (huobiResponse, error) {
	url = fmt.Sprintf(huobiDataUrl, orderType)
	resp, err := http.Get(url)
	if err != nil {
		return huobiResponse{}, err
	}

	var respJson huobiResponse
	json.NewDecoder(resp.Body).Decode(&respJson)

	return respJson, nil
}

func MonitorHuobiPrice(orders chan<- Order) {
	for {
		dataBuy, err := huobiRequest(huobiDataUrl, "sell")
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		bestOrderBuy := dataBuy.Data[0]
		buyPrice, _ := strconv.ParseFloat(bestOrderBuy.Price, 64)

		var paymentMethods string
		for _, method := range bestOrderBuy.PayMethods {
			if len(paymentMethods) != 0 {
				paymentMethods += ","
			}

			paymentMethods += method.Name
		}

		orders <- Order{
			"-",
			"Huobi",
			bestOrderBuy.UserName,
			buyPrice,
			0,
			bestOrderBuy.TradeCount,
			bestOrderBuy.MinTradeLimit,
			bestOrderBuy.MaxTradeLimit,
			paymentMethods,
			fmt.Sprintf("https://c2c.huobi.com/ru-ru/trader/%d", bestOrderBuy.UID),
		}
	}
}
