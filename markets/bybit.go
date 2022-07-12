package markets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type orderByBit struct {
	Id           string
	Price        string
	LastQuantity string
	MinAmount    string
	MaxAmount    string
	UserId       string
	NickName     string
	Payments     []int
}

type results struct {
	Count int
	Items []orderByBit
}

type responseBybit struct {
	RetCode int
	Result  results
}

func findIndexWithPaymentMethods(orders []orderByBit, methods map[int]string) int {
	for i, order := range orders {
		for _, methodId := range order.Payments {
			if _, ok := methods[methodId]; ok {
				return i
			}
		}
	}
	return 0
}

func parsePaymentMethodsByBit(order orderByBit, methods map[int]string) string {
	var paymentsParsed string
	for _, paymentId := range order.Payments {
		if len(paymentsParsed) != 0 {
			paymentsParsed += ","
		}
		paymentName, ok := methods[paymentId]
		if ok {
			paymentsParsed += paymentName
		}
	}
	return paymentsParsed
}

func MonitorByBitPrice(orders chan<- Order) {
	values := url.Values{
		"userId":     {""},
		"tokenId":    {"USDT"},
		"currencyId": {"RUB"},
		"payment":    {""},
		"side":       {"1"},
		"size":       {"15"},
		"page":       {"1"},
		"amount":     {""},
	}
	payments := map[int]string{75: "Tinkoff", 185: "Rosbank", 51: "Payeer",
		14: "Bank Transfer", 27: "FPS", 64: "Raiffeisenbank", 44: "MTS-Bank"}

	for {

		// fetch buy prices
		values["side"] = []string{"1"}
		resp, err := http.PostForm("https://api2.bybit.com/spot/api/otc/item/list", values)

		if err != nil || resp.StatusCode != 200 {
			log.Fatal(resp, err)
		}

		var responseJson responseBybit

		json.NewDecoder(resp.Body).Decode(&responseJson)

		firstValid := findIndexWithPaymentMethods(responseJson.Result.Items, payments)

		buyPrice, err := strconv.ParseFloat(responseJson.Result.Items[firstValid].Price, 64)

		if err != nil {
			log.Fatal(err)
		}

		// fetch sell prices
		values["side"] = []string{"0"}
		resp, err = http.PostForm("https://api2.bybit.com/spot/api/otc/item/list", values)
		if err != nil {
			log.Fatal(err)
		}

		var responseJsonSell responseBybit

		json.NewDecoder(resp.Body).Decode(&responseJsonSell)

		sellValidIndex := findIndexWithPaymentMethods(responseJsonSell.Result.Items, payments)
		sellPrice, err := strconv.ParseFloat(responseJsonSell.Result.Items[sellValidIndex].Price, 64)
		if err != nil {
			log.Fatal(err)
		}

		orders <- Order{
			responseJson.Result.Items[firstValid].Id,
			"ByBit",
			responseJson.Result.Items[firstValid].NickName,
			buyPrice,
			sellPrice,
			responseJson.Result.Items[firstValid].LastQuantity,
			responseJson.Result.Items[firstValid].MinAmount,
			responseJson.Result.Items[firstValid].MaxAmount,
			parsePaymentMethodsByBit(responseJson.Result.Items[firstValid], payments),
			fmt.Sprintf("https://www.bybit.com/fiat/trade/otc/profile/%s/USDT/RUB", responseJson.Result.Items[firstValid].UserId),
		}
	}

}
