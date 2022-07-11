package markets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type order_by_bit struct {
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
	Items []order_by_bit
}

type response_bybit struct {
	Ret_code int
	Result   results
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

		var response_json response_bybit

		json.NewDecoder(resp.Body).Decode(&response_json)

		buy_price, err := strconv.ParseFloat(response_json.Result.Items[0].Price, 64)

		if err != nil {
			log.Fatal(err)
		}

		// fetch sell prices
		values["side"] = []string{"0"}
		resp, err = http.PostForm("https://api2.bybit.com/spot/api/otc/item/list", values)
		if err != nil {
			log.Fatal(err)
		}

		var response_json_sell response_bybit

		json.NewDecoder(resp.Body).Decode(&response_json_sell)

		sell_price, err := strconv.ParseFloat(response_json_sell.Result.Items[0].Price, 64)
		if err != nil {
			log.Fatal(err)
		}

		orders <- Order{
			response_json.Result.Items[0].Id,
			"ByBit",
			response_json.Result.Items[0].NickName,
			buy_price,
			sell_price,
			response_json.Result.Items[0].LastQuantity,
			response_json.Result.Items[0].MinAmount,
			response_json.Result.Items[0].MaxAmount,
			fmt.Sprintf("https://www.bybit.com/fiat/trade/otc/profile/%s/USDT/RUB", response_json.Result.Items[0].UserId),
		}
	}

}
