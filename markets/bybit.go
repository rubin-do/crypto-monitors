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
}

type results struct {
	Count int
	Items []orderByBit
}

type responseBybit struct {
	RetCode int
	Result  results
}

func MonitorByBitPrice(orders chan<- Order) {
	values := url.Values{
		"userId":     {""},
		"tokenId":    {"USDT"},
		"currencyId": {"RUB"},
		"payment":    {"75"},
		"side":       {"1"},
		"size":       {"10"},
		"page":       {"1"},
		"amount":     {""},
	}

	for {

		// fetch buy prices
		values["side"] = []string{"1"}
		resp, err := http.PostForm("https://api2.bybit.com/spot/api/otc/item/list", values)

		if err != nil || resp.StatusCode != 200 {
			log.Fatal(resp, err)
		}

		var responseJson responseBybit

		json.NewDecoder(resp.Body).Decode(&responseJson)

		buyPrice, err := strconv.ParseFloat(responseJson.Result.Items[0].Price, 64)

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

		sellPrice, err := strconv.ParseFloat(responseJsonSell.Result.Items[0].Price, 64)
		if err != nil {
			log.Fatal(err)
		}

		orders <- Order{
			responseJson.Result.Items[0].Id,
			"ByBit",
			responseJson.Result.Items[0].NickName,
			buyPrice,
			sellPrice,
			responseJson.Result.Items[0].LastQuantity,
			responseJson.Result.Items[0].MinAmount,
			responseJson.Result.Items[0].MaxAmount,
			"-",
			fmt.Sprintf("https://www.bybit.com/fiat/trade/otc/profile/%s/USDT/RUB", responseJson.Result.Items[0].UserId),
		}
	}

}
